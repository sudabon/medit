package render

import (
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"

	"github.com/muesli/termenv"
	"golang.org/x/term"
)

// Theme is the viewer color scheme.
type Theme int

const (
	// ThemeDark is for dark backgrounds (light text).
	ThemeDark Theme = iota
	// ThemeLight is for light backgrounds (dark text).
	ThemeLight
)

func (t Theme) String() string {
	if t == ThemeLight {
		return "light"
	}
	return "dark"
}

// DetectTheme chooses a theme from env, OS appearance, then terminal background.
//
// Order:
//  1. GLAMOUR_STYLE=dark|light
//  2. OS appearance (macOS Light/Dark, GNOME color-scheme)
//  3. Terminal background (COLORFGBG / OSC query via termenv)
//  4. ThemeDark
func DetectTheme() Theme {
	switch strings.ToLower(strings.TrimSpace(os.Getenv("GLAMOUR_STYLE"))) {
	case "light":
		return ThemeLight
	case "dark":
		return ThemeDark
	}

	if dark, ok := osAppearanceDark(); ok {
		if dark {
			return ThemeDark
		}
		return ThemeLight
	}

	if dark, ok := terminalBackgroundDark(); ok {
		if dark {
			return ThemeDark
		}
		return ThemeLight
	}

	return ThemeDark
}

// osAppearanceDark reports the desktop Light/Dark preference when available.
func osAppearanceDark() (dark bool, ok bool) {
	switch runtime.GOOS {
	case "darwin":
		return macOSAppearanceDark()
	case "linux":
		return linuxAppearanceDark()
	default:
		return false, false
	}
}

// macOSAppearanceDark uses AppleInterfaceStyle.
// The key is absent in Light mode; present as "Dark" in Dark mode.
func macOSAppearanceDark() (dark bool, ok bool) {
	out, err := exec.Command("defaults", "read", "-g", "AppleInterfaceStyle").CombinedOutput()
	if err != nil {
		// Missing key ⇒ Light.
		if strings.Contains(string(out), "does not exist") || isExitError(err) {
			return false, true
		}
		return false, false
	}
	return strings.EqualFold(strings.TrimSpace(string(out)), "Dark"), true
}

func isExitError(err error) bool {
	_, ok := err.(*exec.ExitError)
	return ok
}

func linuxAppearanceDark() (dark bool, ok bool) {
	out, err := exec.Command("gsettings", "get", "org.gnome.desktop.interface", "color-scheme").Output()
	if err != nil {
		out, err = exec.Command("gsettings", "get", "org.gnome.desktop.interface", "gtk-theme").Output()
		if err != nil {
			return false, false
		}
		name := strings.ToLower(string(out))
		return strings.Contains(name, "dark"), true
	}
	v := strings.ToLower(strings.TrimSpace(string(out)))
	switch {
	case strings.Contains(v, "prefer-dark"):
		return true, true
	case strings.Contains(v, "prefer-light"), strings.Contains(v, "default"):
		return false, true
	default:
		return false, false
	}
}

// terminalBackgroundDark uses COLORFGBG or termenv when the value is trustworthy.
func terminalBackgroundDark() (dark bool, ok bool) {
	if bg, ok := colorFGBGBackground(); ok {
		// 7 and 15 are white / bright white → light background.
		return bg != 7 && bg != 15, true
	}

	if !term.IsTerminal(int(os.Stdout.Fd())) {
		return false, false
	}

	// termenv defaults to black when the OSC query fails, which falsely
	// reports "dark". Only trust it when COLORTERM/TERM suggests a real query
	// path is likely — still prefer OS appearance above. Here we call it as a
	// last hint when OS detection was unavailable.
	return termenv.HasDarkBackground(), true
}

func colorFGBGBackground() (bg int, ok bool) {
	v := os.Getenv("COLORFGBG")
	if !strings.Contains(v, ";") {
		return 0, false
	}
	parts := strings.Split(v, ";")
	n, err := strconv.Atoi(parts[len(parts)-1])
	if err != nil {
		return 0, false
	}
	return n, true
}
