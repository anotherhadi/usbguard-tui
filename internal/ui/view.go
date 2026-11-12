package ui

import (
	"fmt"
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/x/ansi"

	"github.com/anotherhadi/foam/draw"
	"github.com/anotherhadi/foam/minsize"
	"github.com/anotherhadi/foam/theme"

	"github.com/anotherhadi/usbguard-tui/internal/guard"
)

var minSize = minsize.Size{W: 40, H: 15}

func (m *Model) View() tea.View {
	v := tea.NewView("")
	v.AltScreen = true
	v.WindowTitle = "usbguard-tui"
	if theme.Mouse() {
		v.MouseMode = tea.MouseModeCellMotion
	}
	if m.w == 0 || m.h == 0 {
		return v
	}
	if !minSize.Fits(m.w, m.h) {
		v.SetContent(minSize.View(m.w, m.h))
		return v
	}

	header := m.renderHeader()
	top, bottom := len(header), 0
	if theme.HelpBar() {
		bottom = 1
	}
	if theme.Spacious() {
		top++
		if bottom > 0 {
			bottom++
		}
	}
	listH := max(m.h-top-bottom, 0)

	lines := append([]string(nil), header...)
	if theme.Spacious() {
		lines = append(lines, "")
	}
	lines = append(lines, m.devices.View(m.w, listH, m.overlays.Top() == nil))
	if theme.Spacious() && theme.HelpBar() {
		lines = append(lines, "")
	}
	if theme.HelpBar() {
		lines = append(lines, m.helpBar())
	}
	content := draw.Fit(strings.Join(lines, "\n"), m.w, m.h)

	c := lipgloss.NewCanvas(m.w, m.h)
	c.Compose(lipgloss.NewCompositor(lipgloss.NewLayer(content)))
	draw.FillRow(c, 0, theme.Surface())
	if theme.HelpBar() {
		draw.FillRow(c, m.h-1, theme.Surface())
	}
	m.overlays.Draw(c)
	if t := m.toasts(); len(t) > 0 {
		c.Compose(lipgloss.NewCompositor(t...))
	}
	v.SetContent(c.Render())
	return v
}

func (m *Model) headerHeight() int { return len(m.renderHeader()) }

// bodyTop is where the device list starts: under the header, and under the
// blank line foam's spacious setting adds there.
func (m *Model) bodyTop() int {
	top := m.headerHeight()
	if theme.Spacious() {
		top++
	}
	return top
}

func (m *Model) renderHeader() []string {
	title := theme.Accent().Style().Bold(true).Render(" " + theme.Icon("󱊟 ", "") + "usbguard-tui ")
	lines := []string{
		title,
		m.renderServiceLine(),
		m.renderPolicyLine(),
		m.renderDevicesLine(),
		m.renderRulesLine(),
	}
	if pending := m.renderPendingRulesLine(); pending != "" {
		lines = append(lines, pending)
	}
	for i, l := range lines {
		lines[i] = ansi.Truncate(l, max(m.w, 0), "…")
	}
	return lines
}

func (m *Model) renderServiceLine() string {
	label := infoLabel("Service")
	switch m.daemonStatus {
	case "active":
		return label + theme.Success().Render("active")
	case "":
		return label + theme.Muted().Render("checking…")
	default:
		return label + theme.Muted().Render(m.daemonStatus)
	}
}

func (m *Model) renderPolicyLine() string {
	label := infoLabel("Default policy")
	if m.defaultPolicy == "" {
		return label + theme.Muted().Render("unknown")
	}
	return label + statusColor(m.defaultPolicy).Render(string(m.defaultPolicy))
}

func (m *Model) renderDevicesLine() string {
	label := infoLabel("Devices")
	if m.counts.total == 0 {
		return label + theme.Muted().Render("0")
	}
	parts := []string{theme.Muted().Render(fmt.Sprintf("%d total", m.counts.total))}
	if m.counts.allowed > 0 {
		parts = append(parts, statusColor(guard.Allowed).Render(fmt.Sprintf("%d allow", m.counts.allowed)))
	}
	if m.counts.blocked > 0 {
		parts = append(parts, statusColor(guard.Blocked).Render(fmt.Sprintf("%d block", m.counts.blocked)))
	}
	if m.counts.rejected > 0 {
		parts = append(parts, statusColor(guard.Rejected).Render(fmt.Sprintf("%d reject", m.counts.rejected)))
	}
	return label + strings.Join(parts, theme.Muted().Render("  ·  "))
}

func (m *Model) renderRulesLine() string {
	label := infoLabel("Rules")
	switch {
	case m.rulesManaged:
		return label + theme.Warning().Render("read-only (NixOS managed)")
	case !m.rulesWritable:
		return label + theme.Warning().Render("read-only")
	default:
		return label + theme.Success().Render("writable")
	}
}

func (m *Model) renderPendingRulesLine() string {
	n := len(m.pendingRules)
	if n == 0 {
		return ""
	}
	label := infoLabel("Pending rules")
	return label + theme.Warning().Render(fmt.Sprintf("%s queued (%s to copy)", plural(n, "rule", "rules"), m.keyHint("device.copy-rules")))
}

// toasts draws the live notifications at the top right, newest first.
func (m *Model) toasts() []*lipgloss.Layer {
	var out []*lipgloss.Layer
	w := min(48, max(m.w/3, 24), m.w)
	y := 0
	m.toastRects = m.toastRects[:0]
	live := m.notes.Live()
	for i := len(live) - 1; i >= 0 && i >= len(live)-5; i-- {
		s := live[i].Toast(w)
		h := lipgloss.Height(s)
		if len(out) > 0 && y+h > m.h {
			break // the older ones would be cut
		}
		m.toastRects = append(m.toastRects, toastRect{X: max(m.w-w, 0), Y: y, W: w, H: h, id: live[i].ID})
		out = append(out, lipgloss.NewLayer(s).X(max(m.w-w, 0)).Y(y).Z(20))
		y += h
	}
	return out
}

func (m *Model) helpBar() string {
	ks, ls := theme.Primary().Style(), theme.Subtle().Style()
	var parts []string
	for _, h := range m.reg.Hints(m.path()) {
		first, _, _ := strings.Cut(h.Keys, " / ")
		parts = append(parts, ks.Render(first)+" "+ls.Render(h.Label))
	}
	return draw.Pad(ansi.Truncate(" "+strings.Join(parts, ls.Render("  ·  ")), m.w, "…"), m.w)
}
