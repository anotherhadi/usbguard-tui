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
	rules := listRules()
	var devices []Device
	for _, line := range strings.Split(strings.TrimSpace(string(out)), "\n") {
		if line == "" {
			continue
		}
		d, err := parseLine(line)
		if err == nil {
			d.Permanent = rules[d.Hash] == d.Status
			devices = append(devices, d)
		}
	}
	return devices, nil
}

func listRules() map[string]Status {
	out, err := exec.Command("usbguard", "list-rules").Output()
	if err != nil {
		return nil
	}
	rules := make(map[string]Status)
	for _, line := range strings.Split(strings.TrimSpace(string(out)), "\n") {
		if line == "" {
			continue
		}
		d, err := parseLine(line)
		if err == nil && d.Hash != "" {
			rules[d.Hash] = d.Status
		}
	}
	return rules
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
	out, err := exec.Command("systemctl", "cat", "usbguard").Output()
	if err != nil {
		return false
	}
	configPath := extractConfigPath(string(out))
	if configPath == "" {
		return false
	}
	ruleFile := parseRuleFilePath(configPath)
	return strings.HasPrefix(ruleFile, "/nix/store/")
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
	case strings.Contains(lower, "permission denied"), strings.Contains(lower, "not authorized"):
		return ErrPermission
	case strings.Contains(lower, "read-only"), strings.Contains(lower, "immutable"):
		return ErrReadOnly
	default:
		return errors.New(strings.TrimSpace(output))
	}
}
