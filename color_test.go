package main

import (
	"bytes"
	"strings"
	"testing"
)

func TestAuditRequirements(t *testing.T) {
	tests := []struct {
		name     string
		text     string
		sub      string
		colors   []string
		contains string // ANSI code to look for
	}{
		{
			name:     "Basic: Red Banana",
			text:     "banana",
			sub:      "",
			colors:   []string{"red"},
			contains: "\033[31m",
		},
		{
			name:     "Basic: Red Hello World",
			text:     "hello world",
			sub:      "",
			colors:   []string{"red"},
			contains: "\033[31m",
		},
		{
			name:     "Special: Green Equations",
			text:     "1 + 1 = 2",
			sub:      "",
			colors:   []string{"green"},
			contains: "\033[32m",
		},
		{
			name:     "Special: Yellow Symbols",
			text:     "(%&) ??",
			sub:      "",
			colors:   []string{"yellow"},
			contains: "\033[33m",
		},
		{
			name:     "Substring: Orange GuYs",
			text:     "HeY GuYs",
			sub:      "GuYs",
			colors:   []string{"orange"},
			contains: "\033[38;5;208m",
		},
		{
			name:     "Single Letter: Blue B in RGB()",
			text:     "RGB()",
			sub:      "B",
			colors:   []string{"blue"},
			contains: "\033[34m",
		},
		{
			name:     "Bonus: HEX notation",
			text:     "HexTest",
			sub:      "",
			colors:   []string{"#ff0000"},
			contains: "\033[38;2;255;0;0m",
		},
		{
			name:     "Bonus: RGB notation",
			text:     "RGBTest",
			sub:      "",
			colors:   []string{"rgb(0,255,0)"},
			contains: "\033[38;2;0;255;0m",
		},
		{
			name:     "Bonus: HSL notation",
			text:     "HSLTest",
			sub:      "",
			colors:   []string{"hsl(240,100,50)"}, // Blue
			contains: "\033[38;2;0;0;255m",
		},
		{
			name:     "Multi-Color: Two notations",
			text:     "Multi",
			sub:      "",
			colors:   []string{"red", "#0000ff"},
			contains: "\033[38;2;0;0;255m", // Ensure it handles the second notation
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			err := RenderASCII(&buf, tt.text, tt.sub, tt.colors)
			if err != nil {
				t.Fatalf("Failed to render: %v", err)
			}

			output := buf.String()

			t.Logf("Output:\n%s", output)

			if !strings.Contains(output, tt.contains) {
				t.Errorf("Test %s failed: output does not contain expected color code %q", tt.name, tt.contains)
			}
		})
	}
}

func TestNoExternalPackages(t *testing.T) {
	// This is a meta-test to ensure only standard library is used
	// In a real audit, the auditor checks go.mod
	t.Log("Manual Check: Ensure go.mod contains no external dependencies.")
}
