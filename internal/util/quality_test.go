package util

import (
	"testing"
)

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
	for _, tt := range tests {
		got, err := PresetVal(tt.level, tt.low, tt.medium, tt.high)
		if (err != nil) != tt.err {
			t.Errorf("PresetVal(%q) error = %v; want error=%v", tt.level, err, tt.err)
		}
		if got != tt.want {
			t.Errorf("PresetVal(%q) = %d; want %d", tt.level, got, tt.want)
		}
	}
}

func TestPresetValEmpty(t *testing.T) {
	_, err := PresetVal("", 60, 80, 90)
	if err == nil {
		t.Error("PresetVal('') expected error")
	}
}

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
	for _, tt := range tests {
		crf, preset, err := VideoQuality(tt.level)
		if (err != nil) != tt.err {
			t.Errorf("VideoQuality(%q) error = %v; want error=%v", tt.level, err, tt.err)
		}
		if crf != tt.wantCRF || preset != tt.wantPres {
			t.Errorf("VideoQuality(%q) = (%d, %q); want (%d, %q)", tt.level, crf, preset, tt.wantCRF, tt.wantPres)
		}
	}
}

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
	for _, tt := range tests {
		got, err := AudioBitrate(tt.level)
		if (err != nil) != tt.err {
			t.Errorf("AudioBitrate(%q) error = %v; want error=%v", tt.level, err, tt.err)
		}
		if got != tt.want {
			t.Errorf("AudioBitrate(%q) = %q; want %q", tt.level, got, tt.want)
		}
	}
}

func TestPDFDPI(t *testing.T) {
	tests := []struct {
		level  string
		want   int
		wantS  string
		err    bool
	}{
		{"low", 150, "/screen", false},
		{"medium", 300, "/ebook", false},
		{"high", 400, "/printer", false},
		{"200", 200, "", false},
		{"invalid", 0, "", true},
	}
	for _, tt := range tests {
		dpi, settings, err := PDFDPI(tt.level)
		if (err != nil) != tt.err {
			t.Errorf("PDFDPI(%q) error = %v; want error=%v", tt.level, err, tt.err)
		}
		if dpi != tt.want || settings != tt.wantS {
			t.Errorf("PDFDPI(%q) = (%d, %q); want (%d, %q)", tt.level, dpi, settings, tt.want, tt.wantS)
		}
	}
}
