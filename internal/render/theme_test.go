package render

import (
	"runtime"
	"testing"
)

func TestDetectThemeEnvOverride(t *testing.T) {
	t.Setenv("GLAMOUR_STYLE", "light")
	if DetectTheme() != ThemeLight {
		t.Fatal("expected light from env")
	}
	t.Setenv("GLAMOUR_STYLE", "dark")
	if DetectTheme() != ThemeDark {
		t.Fatal("expected dark from env")
	}
}

func TestThemeString(t *testing.T) {
	if ThemeLight.String() != "light" {
		t.Fatal(ThemeLight.String())
	}
	if ThemeDark.String() != "dark" {
		t.Fatal(ThemeDark.String())
	}
}

func TestColorFGBGBackground(t *testing.T) {
	t.Setenv("COLORFGBG", "0;15")
	bg, ok := colorFGBGBackground()
	if !ok || bg != 15 {
		t.Fatalf("got %d %v", bg, ok)
	}
}

func TestMacOSAppearanceLightWhenUnset(t *testing.T) {
	if runtime.GOOS != "darwin" {
		t.Skip("macOS only")
	}
	dark, ok := macOSAppearanceDark()
	if !ok {
		t.Fatal("expected macOS appearance detection to succeed")
	}
	// This machine is in Light mode (AppleInterfaceStyle absent).
	if dark {
		t.Log("system is currently Dark — skipping light assertion")
		return
	}
	t.Setenv("GLAMOUR_STYLE", "")
	if DetectTheme() != ThemeLight {
		t.Fatalf("expected ThemeLight on macOS Light appearance, got %s", DetectTheme())
	}
}
