package guard

import (
	"errors"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func Check() error {
	_, err := exec.LookPath("usbguard")
	if err != nil {
		return ErrNotFound
	}
	return nil
}

func ListDevices() ([]Device, error) {
	out, err := exec.Command("usbguard", "list-devices").Output()
	if err != nil {
		return nil, wrapExecError(err)
	}

	deviceRuleText := devicePolicyRuleText()

	permanentTexts, havePermanentTexts := permanentRuleTexts()

	implicitTarget, haveImplicitTarget := implicitPolicyTarget()

	var devices []Device
	for _, line := range strings.Split(strings.TrimSpace(string(out)), "\n") {
		if line == "" {
			continue
		}
		d, err := parseLine(line)
		if err != nil {
			continue
		}
		d.RuleState = resolveRuleState(d, deviceRuleText, permanentTexts, havePermanentTexts, implicitTarget, haveImplicitTarget)
		devices = append(devices, d)
	}
	return devices, nil
}

func resolveRuleState(d Device, deviceRuleText map[int]string, permanentTexts map[string]bool, havePermanentTexts bool, implicitTarget Status, haveImplicitTarget bool) RuleState {
	text, matched := deviceRuleText[d.ID]
	if matched {
		if !havePermanentTexts {

			return RulePermanent
		}
		if permanentTexts[text] {
			return RulePermanent
		}
		return RuleTemporary
	}

	if !haveImplicitTarget || d.Status != implicitTarget {
		return RuleTemporary
	}
	return RuleDefault
}

func implicitPolicyTarget() (Status, bool) {
	out, err := exec.Command("usbguard", "get-parameter", "ImplicitPolicyTarget").Output()
	if err != nil {
		return "", false
	}
	return Status(strings.TrimSpace(string(out))), true
}

func DefaultPolicy() Status {
	target, _ := implicitPolicyTarget()
	return target
}

func devicePolicyRuleText() map[int]string {
	out, err := exec.Command("usbguard", "list-rules", "-d").Output()
	if err != nil {
		return nil
	}
	result := make(map[int]string)
	var currentRuleText string
	for _, line := range strings.Split(string(out), "\n") {
		if line == "" {
			continue
		}
		if line[0] == ' ' || line[0] == '\t' {
			trimmed := strings.TrimLeft(line, " \t")
			colonIdx := strings.Index(trimmed, ":")
			if colonIdx < 0 {
				continue
			}
			devID, err := strconv.Atoi(strings.TrimSpace(trimmed[:colonIdx]))
			if err != nil {
				continue
			}
			result[devID] = currentRuleText
			continue
		}
		colonIdx := strings.Index(line, ":")
		if colonIdx < 0 {
			continue
		}
		currentRuleText = normalizeRuleText(line[colonIdx+1:])
	}
	return result
}

func permanentRuleTexts() (map[string]bool, bool) {
	path := ruleFilePath()
	if path == "" {
		return nil, false
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, false
	}
	texts := make(map[string]bool)
	for _, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		texts[normalizeRuleText(line)] = true
	}
	return texts, true
}

func normalizeRuleText(s string) string {
	return strings.Join(strings.Fields(s), " ")
}

func AllowDevice(id int, permanent bool) error  { return applyPolicy("allow-device", id, permanent) }
func BlockDevice(id int, permanent bool) error  { return applyPolicy("block-device", id, permanent) }
func RejectDevice(id int, permanent bool) error { return applyPolicy("reject-device", id, permanent) }

func DaemonStatus() string {
	out, err := exec.Command("systemctl", "is-active", "usbguard").Output()
	if err != nil {
		return "unknown"
	}
	return strings.TrimSpace(string(out))
}

func applyPolicy(cmd string, id int, permanent bool) error {
	args := []string{cmd}
	if permanent {
		args = append(args, "-p")
	}
	args = append(args, strconv.Itoa(id))
	out, err := exec.Command("usbguard", args...).CombinedOutput()
	if err != nil {
		return classifyError(string(out))
	}
	return nil
}

func wrapExecError(err error) error {
	var exitErr *exec.ExitError
	if errors.As(err, &exitErr) {
		return classifyError(string(exitErr.Stderr))
	}
	return err
}

func IsRulesManaged() bool {
	return strings.HasPrefix(ruleFilePath(), "/nix/store/")
}

func RulesWritable() (bool, bool) {
	path := ruleFilePath()
	if path == "" {
		return false, false
	}
	f, err := os.OpenFile(path, os.O_WRONLY, 0)
	if err != nil {
		return false, true
	}
	f.Close()
	return true, true
}

func ruleFilePath() string {
	out, err := exec.Command("systemctl", "cat", "usbguard").Output()
	if err != nil {
		return ""
	}
	configPath := extractConfigPath(string(out))
	if configPath == "" {
		return ""
	}
	return parseRuleFilePath(configPath)
}

func extractConfigPath(s string) string {
	fields := strings.Fields(s)
	for i, f := range fields {
		if f == "-c" && i+1 < len(fields) {
			return fields[i+1]
		}
	}
	return ""
}

func parseRuleFilePath(configPath string) string {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return ""
	}
	for _, line := range strings.Split(string(data), "\n") {
		if after, ok := strings.CutPrefix(line, "RuleFile="); ok {
			return strings.TrimSpace(after)
		}
	}
	return ""
}

func classifyError(output string) error {
	lower := strings.ToLower(output)
	switch {
	case strings.Contains(lower, "permission denied"), strings.Contains(lower, "not authorized"),
		strings.Contains(lower, "operation not permitted"):
		return ErrPermission
	case strings.Contains(lower, "read-only"), strings.Contains(lower, "immutable"):
		return ErrReadOnly
	default:
		return errors.New(strings.TrimSpace(output))
	}
}
