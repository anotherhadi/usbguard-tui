package ui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/ansi"
)

var dimStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("238"))

// placeOverlay renders fg centered over bg, with bg stripped and rendered dim gray.
func placeOverlay(bg, fg string, width, height int) string {
	fgLines := strings.Split(fg, "\n")
	fgH := len(fgLines)
	fgW := 0
	for _, l := range fgLines {
		if w := lipgloss.Width(l); w > fgW {
			fgW = w
		}
	}

	bgLines := strings.Split(bg, "\n")

	x0 := (width - fgW) / 2
	y0 := (height - fgH) / 2

	result := make([]string, height)
	for i := 0; i < height; i++ {
		raw := ""
		if i < len(bgLines) {
			raw = ansi.Strip(bgLines[i])
		}
		if w := lipgloss.Width(raw); w < width {
			raw += strings.Repeat(" ", width-w)
		}

		fgIdx := i - y0
		if fgIdx < 0 || fgIdx >= fgH {
			result[i] = dimStyle.Render(raw)
		} else {
			fgLine := fgLines[fgIdx]
			fgLineW := lipgloss.Width(fgLine)
			left := ansi.Truncate(raw, x0, "")
			right := ansi.Cut(raw, x0+fgLineW, width)
			result[i] = dimStyle.Render(left) + fgLine + dimStyle.Render(right)
		}
	}

	return strings.Join(result, "\n")
}
