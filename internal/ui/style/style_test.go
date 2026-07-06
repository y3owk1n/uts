//nolint:testpackage
package style

import (
	"bytes"
	"os"
	"strings"
	"sync"
	"testing"
)

// unsetColorEnv clears every UTS_COLOR_<NAME> variable so
// that tests start from a known state.
func unsetColorEnv(t *testing.T) {
	t.Helper()

	slots := []string{
		"PRIMARY", "TEXT", "MUTED", "SUBTLE",
		"BORDER", "ACCENT", "SUCCESS", "WARNING", "ERROR",
	}
	for _, slot := range slots {
		t.Setenv("UTS_COLOR_"+slot, "")
		t.Setenv("UTS_COLOR_"+slot+"_LIGHT", "")
		t.Setenv("UTS_COLOR_"+slot+"_DARK", "")
	}
}

func TestIsValidColor(t *testing.T) {
	tests := []struct {
		val  string
		want bool
	}{
		// Hex with #
		{"#abc", true},
		{"#abcdef", true},
		{"#abcdef12", true},
		// Hex without # (should still be valid)
		{"abc", true},
		{"abcdef", true},
		{"abcdef12", true},
		// Named colors
		{"red", true},
		{"BLACK", true},
		{"cyan", true},
		// ANSI 256 range
		{"0", true},
		{"255", true},
		// Invalid
		{"256", false},
		{"-1", false},
		{"invalid", false},
		{"#abcd", false}, // 4 digits
		{"", false},
	}

	for _, tt := range tests {
		if got := isValidColor(tt.val); got != tt.want {
			t.Errorf("isValidColor(%q) = %v, want %v", tt.val, got, tt.want)
		}
	}
}

func TestEnvColor_PrependHash(t *testing.T) {
	t.Setenv("UTS_COLOR_TEST", "c595d4")

	val, ok := envColor("UTS_COLOR_TEST")
	if !ok || val != "#c595d4" {
		t.Errorf("envColor with missing # got (%q, %v), want (\"#c595d4\", true)", val, ok)
	}
}

func TestEnvColor_ANSI256NotPrepended(t *testing.T) {
	t.Setenv("UTS_COLOR_TEST", "255")

	val, ok := envColor("UTS_COLOR_TEST")
	if !ok || val != "255" {
		t.Errorf("envColor with ANSI 255 got (%q, %v), want (\"255\", true)", val, ok)
	}
}

func TestEnvColor_ValidHexUnchanged(t *testing.T) {
	t.Setenv("UTS_COLOR_TEST", "#112233")

	val, ok := envColor("UTS_COLOR_TEST")
	if !ok || val != "#112233" {
		t.Errorf("envColor with valid hex got (%q, %v), want (\"#112233\", true)", val, ok)
	}
}

func TestEnvColor_InvalidReturnsEmpty(t *testing.T) {
	t.Setenv("UTS_COLOR_TEST", "garbage")

	oldStderr := os.Stderr
	read, write, _ := os.Pipe()
	os.Stderr = write
	warnedInvalidColors = sync.Map{}

	val, valOk := envColor("UTS_COLOR_TEST")

	_ = write.Close()

	os.Stderr = oldStderr

	var buf bytes.Buffer

	_, _ = buf.ReadFrom(read)
	warning := buf.String()

	if valOk || val != "" {
		t.Errorf("envColor with invalid color got (%q, %v), want (\"\", false)", val, valOk)
	}

	if !strings.Contains(warning, "UTS_COLOR_TEST") ||
		!strings.Contains(warning, "is not a valid color") {
		t.Errorf("expected warning on stderr, got %q", warning)
	}
}

func TestEnvColor_EmptyReturnsEmpty(t *testing.T) {
	t.Setenv("UTS_COLOR_TEST", "")

	val, ok := envColor("UTS_COLOR_TEST")
	if ok || val != "" {
		t.Errorf("envColor with empty got (%q, %v), want (\"\", false)", val, ok)
	}
}

func TestEnvColor_WhitespaceOnlyReturnsEmpty(t *testing.T) {
	t.Setenv("UTS_COLOR_TEST", "   ")

	val, ok := envColor("UTS_COLOR_TEST")
	if ok || val != "" {
		t.Errorf("envColor with spaces got (%q, %v), want (\"\", false)", val, ok)
	}
}

func TestDefault_OverrideBothVariants(t *testing.T) {
	unsetColorEnv(t)
	t.Setenv("UTS_COLOR_PRIMARY", "#FF00FF")

	pal := Default()

	r, g, b, _ := pal.Primary.Light.RGBA() //nolint:varnamelen
	if r>>8 != 255 || g>>8 != 0 || b>>8 != 255 {
		t.Errorf("Primary.Light = (%d,%d,%d), want (255,0,255)", r>>8, g>>8, b>>8)
	}

	r, g, b, _ = pal.Primary.Dark.RGBA()
	if r>>8 != 255 || g>>8 != 0 || b>>8 != 255 {
		t.Errorf("Primary.Dark = (%d,%d,%d), want (255,0,255)", r>>8, g>>8, b>>8)
	}
}

func TestDefault_OverrideLightDarkSeparately(t *testing.T) {
	unsetColorEnv(t)
	t.Setenv("UTS_COLOR_ACCENT_LIGHT", "#111111")
	t.Setenv("UTS_COLOR_ACCENT_DARK", "#EEEEEE")

	pal := Default()

	read, _, _, _ := pal.Accent.Light.RGBA() //nolint:dogsled
	if read>>8 != 0x11 {
		t.Errorf("Accent.Light R = %d, want 0x11", read>>8)
	}

	read, _, _, _ = pal.Accent.Dark.RGBA() //nolint:dogsled
	if read>>8 != 0xEE {
		t.Errorf("Accent.Dark R = %d, want 0xEE", read>>8)
	}
}

func TestDefault_InvalidOverrideKeepsDefault(t *testing.T) {
	unsetColorEnv(t)

	oldStderr := os.Stderr
	_, write, _ := os.Pipe()
	os.Stderr = write
	warnedInvalidColors = sync.Map{}

	t.Setenv("UTS_COLOR_ERROR", "garbage")

	pal := Default()

	_ = write.Close()

	os.Stderr = oldStderr

	// Should keep the base palette defaults.
	base := basePalette()
	wantR, wantG, wantB, _ := base.Error.Light.RGBA()

	gotR, gotG, gotB, _ := pal.Error.Light.RGBA()
	if gotR != wantR || gotG != wantG || gotB != wantB {
		t.Error("Error.Light should keep default when UTS_COLOR_ERROR is invalid")
	}
}
