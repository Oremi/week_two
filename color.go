package main

import (
	"flag"
	"fmt"
	"math"
	"os"
	"strings"
)

type colorFlags []string

func (i *colorFlags) String() string     { return strings.Join(*i, ",") }
func (i *colorFlags) Set(v string) error { *i = append(*i, v); return nil }

const charHeight = 8

func main() {
	var colors colorFlags
	flag.Var(&colors, "color", "Color (name, #hex, rgb(r,g,b), or hsl(h,s,l))")
	flag.Parse()

	args := flag.Args()
	if len(args) == 0 {
		fmt.Println("Usage: go run . --color=red [substring] \"text\"")
		return
	}

	text, sub := args[len(args)-1], ""
	if len(args) > 1 {
		sub = args[0]
	}

	data, _ := os.ReadFile("standard.txt")
	font := strings.Split(strings.ReplaceAll(string(data), "\r\n", "\n"), "\n")

	var ansiCodes []string
	for _, c := range colors {
		if code := parseColor(c); code != "" {
			ansiCodes = append(ansiCodes, code)
		}
	}

	for _, line := range strings.Split(strings.ReplaceAll(text, "\\n", "\n"), "\n") {
		if line == "" {
			fmt.Println()
			continue
		}

		for row := 0; row < charHeight; row++ {
			colorIdx := 0
			for i := 0; i < len(line); i++ {
				isColored := len(ansiCodes) > 0 && (sub == "" || isInside(line, sub, i))

				if isColored {
					fmt.Print(ansiCodes[colorIdx%len(ansiCodes)])
					colorIdx++
				}

				fIdx := int(line[i]-32)*9 + 1 + row
				if fIdx < len(font) {
					fmt.Print(font[fIdx])
				}
				if isColored {
					fmt.Print("\033[0m")
				}
			}
			fmt.Println()
		}
	}
}

func parseColor(c string) string {
	c = strings.ToLower(strings.TrimSpace(c))
	var r, g, b int

	// Named Colors
	names := map[string]string{"red": "31", "green": "32", "yellow": "33", "blue": "34", "magenta": "35", "cyan": "36", "orange": "38;5;208"}
	if code, ok := names[c]; ok {
		return "\033[" + code + "m"
	}

	// HEX
	if _, err := fmt.Sscanf(c, "#%02x%02x%02x", &r, &g, &b); err == nil {
		return fmt.Sprintf("\033[38;2;%d;%d;%dm", r, g, b)
	}
	// RGB
	if _, err := fmt.Sscanf(c, "rgb(%d,%d,%d)", &r, &g, &b); err == nil {
		return fmt.Sprintf("\033[38;2;%d;%d;%dm", r, g, b)
	}
	// HSL
	var h, s, l float64
	if _, err := fmt.Sscanf(strings.ReplaceAll(c, "%", ""), "hsl(%f,%f,%f)", &h, &s, &l); err == nil {
		r, g, b := hslToRgb(h/360, s/100, l/100)
		return fmt.Sprintf("\033[38;2;%d;%d;%dm", r, g, b)
	}
	return ""
}

func hslToRgb(h, s, l float64) (int, int, int) {
	f := func(n float64) float64 {
		k := math.Mod(n+h*12, 12)
		a := s * math.Min(l, 1-l)
		return l - a*math.Max(-1, math.Min(math.Min(k-3, 9-k), 1))
	}
	return int(f(0) * 255), int(f(8) * 255), int(f(4) * 255)
}

func isInside(text, sub string, i int) bool {
	for start := 0; ; {
		idx := strings.Index(text[start:], sub)
		if idx == -1 {
			break
		}
		if i >= start+idx && i < start+idx+len(sub) {
			return true
		}
		start += idx + 1
	}
	return false
}
