package guard

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type Status string

const (
	Allowed  Status = "allow"
	Blocked  Status = "block"
	Rejected Status = "reject"
)

type RuleState string

const (
	RuleDefault   RuleState = "default"
	RuleTemporary RuleState = "temporary"
	RulePermanent RuleState = "permanent"
)

type Device struct {
	ID        int
	Name      string
	Status    Status
	VidPid    string
	Hash      string
	RuleState RuleState
}

func (d Device) Title() string       { return d.Name }
func (d Device) Description() string { return fmt.Sprintf("id:%-3d  %s", d.ID, d.VidPid) }
func (d Device) FilterValue() string { return d.Name + " " + d.VidPid }

func parseLine(line string) (Device, error) {
	colonIdx := strings.Index(line, ":")
	if colonIdx < 0 {
		return Device{}, errors.New("invalid format")
	}

	id, err := strconv.Atoi(strings.TrimSpace(line[:colonIdx]))
	if err != nil {
		return Device{}, err
	}

	rest := strings.TrimSpace(line[colonIdx+1:])
	parts := strings.Fields(rest)
	if len(parts) < 1 {
		return Device{}, errors.New("missing status")
	}

	status := Status(parts[0])
	name := extractField(rest, "name")
	if name == "" {
		name = fmt.Sprintf("Unknown Device #%d", id)
	}

	return Device{
		ID:     id,
		Name:   name,
		Status: status,
		VidPid: extractUnquoted(rest, "id"),
		Hash:   extractField(rest, "hash"),
	}, nil
}

func extractField(rule, field string) string {
	prefix := field + ` "`
	idx := strings.Index(rule, prefix)
	if idx < 0 {
		return ""
	}
	rest := rule[idx+len(prefix):]
	end := strings.Index(rest, `"`)
	if end < 0 {
		return ""
	}
	return rest[:end]
}

func NixOSRule(dev Device, status Status) string {
	return fmt.Sprintf("%s id %s name \"%s\"", status, dev.VidPid, dev.Name)
}

func extractUnquoted(rule, field string) string {
	prefix := field + " "
	idx := strings.Index(rule, prefix)
	if idx < 0 {
		return ""
	}
	rest := rule[idx+len(prefix):]
	end := strings.IndexAny(rest, " \t\n")
	if end < 0 {
		return rest
	}
	return rest[:end]
}
