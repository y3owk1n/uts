//nolint:goconst,testpackage
package util

import (
	"testing"
)

// TestPresetVal tests PresetVal.
func TestPresetVal(t *testing.T) {
	tests := []struct {
		level  string
		low    int
		medium int
		high   int
		want   int
		err    bool
	}{
		{"low", 60, 80, 90, 60, false},
		{"medium", 60, 80, 90, 80, false},
		{"high", 60, 80, 90, 90, false},
		{"75", 60, 80, 90, 75, false},
		{"invalid", 60, 80, 90, 0, true},
	}
	for _, testCase := range tests {
		got, err := PresetVal(testCase.level, testCase.low, testCase.medium, testCase.high)
		if (err != nil) != testCase.err {
			t.Errorf("PresetVal(%q) error = %v; want error=%v", testCase.level, err, testCase.err)
		}

		if got != testCase.want {
			t.Errorf("PresetVal(%q) = %d; want %d", testCase.level, got, testCase.want)
		}
	}
}

// TestPresetValEmpty tests PresetVal with an empty string.
func TestPresetValEmpty(t *testing.T) {
	_, err := PresetVal("", 60, 80, 90)
	if err == nil {
		t.Error("PresetVal('') expected error")
	}
}

// TestVideoQuality tests VideoQuality.
func TestVideoQuality(t *testing.T) {
	tests := []struct {
		level    string
		wantCRF  int
		wantPres string
		err      bool
	}{
		{"low", 32, "fast", false},
		{"medium", 28, "medium", false},
		{"high", 23, "slow", false},
		{"20", 20, "medium", false},
		{"15", 15, "slow", false},
		{"30", 30, "fast", false},
		{"invalid", 0, "", true},
	}
	for _, testCase := range tests {
		crf, preset, err := VideoQuality(testCase.level)
		if (err != nil) != testCase.err {
			t.Errorf(
				"VideoQuality(%q) error = %v; want error=%v",
				testCase.level,
				err,
				testCase.err,
			)
		}

		if crf != testCase.wantCRF || preset != testCase.wantPres {
			t.Errorf(
				"VideoQuality(%q) = (%d, %q); want (%d, %q)",
				testCase.level,
				crf,
				preset,
				testCase.wantCRF,
				testCase.wantPres,
			)
		}
	}
}

// TestAudioBitrate tests AudioBitrate.
func TestAudioBitrate(t *testing.T) {
	tests := []struct {
		level string
		want  string
		err   bool
	}{
		{"low", "96k", false},
		{"medium", "128k", false},
		{"high", "192k", false},
		{"256", "256k", false},
		{"invalid", "", true},
	}
	for _, testCase := range tests {
		got, err := AudioBitrate(testCase.level)
		if (err != nil) != testCase.err {
			t.Errorf(
				"AudioBitrate(%q) error = %v; want error=%v",
				testCase.level,
				err,
				testCase.err,
			)
		}

		if got != testCase.want {
			t.Errorf("AudioBitrate(%q) = %q; want %q", testCase.level, got, testCase.want)
		}
	}
}

// TestPDFDPI tests PDFDPI.
func TestPDFDPI(t *testing.T) {
	tests := []struct {
		level string
		want  int
		wantS string
		err   bool
	}{
		{"low", 150, "/screen", false},
		{"medium", 300, "/ebook", false},
		{"high", 400, "/printer", false},
		{"200", 200, "", false},
		{"invalid", 0, "", true},
	}
	for _, testCase := range tests {
		dpi, settings, err := PDFDPI(testCase.level)
		if (err != nil) != testCase.err {
			t.Errorf("PDFDPI(%q) error = %v; want error=%v", testCase.level, err, testCase.err)
		}

		if dpi != testCase.want || settings != testCase.wantS {
			t.Errorf(
				"PDFDPI(%q) = (%d, %q); want (%d, %q)",
				testCase.level,
				dpi,
				settings,
				testCase.want,
				testCase.wantS,
			)
		}
	}
}

// TestQualityRanges tests that out-of-range numeric values are rejected with
// a clear error instead of being handed to the tool.
func TestQualityRanges(t *testing.T) {
	bad := []struct{ kind, level string }{
		{"image", "0"},
		{"image", "101"},
		{"video", "52"},
		{"audio", "20"},
		{"audio", "9999"},
		{"pdf", "10"},
		{"pdf", "5000"},
	}

	for _, testCase := range bad {
		var err error

		switch testCase.kind {
		case "image":
			_, err = ImageQuality(testCase.level)
		case "video":
			_, _, err = VideoQuality(testCase.level)
		case "audio":
			_, err = AudioBitrate(testCase.level)
		default:
			_, _, err = PDFDPI(testCase.level)
		}

		if err == nil {
			t.Errorf("%s quality %s: expected an out-of-range error", testCase.kind, testCase.level)
		}
	}

	val, err := ImageQuality("100")
	if err != nil || val != 100 {
		t.Errorf("ImageQuality(100) = %d, %v; want 100, nil", val, err)
	}

	crf, _, err := VideoQuality("0")
	if err != nil || crf != 0 {
		t.Errorf("VideoQuality(0) = %d, %v; want 0, nil", crf, err)
	}
}
