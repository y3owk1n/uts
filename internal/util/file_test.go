package util

import (
	"os"
	"path/filepath"
	"testing"
)

func TestHumanSize(t *testing.T) {
	tests := []struct {
		input int64
		want  string
	}{
		{0, "0 B"},
		{512, "512 B"},
		{1024, "1.0 KB"},
		{1536, "1.5 KB"},
		{1048576, "1.0 MB"},
		{1572864, "1.5 MB"},
		{1073741824, "1.0 GB"},
	}
	for _, tt := range tests {
		got := HumanSize(tt.input)
		if got != tt.want {
			t.Errorf("HumanSize(%d) = %q; want %q", tt.input, got, tt.want)
		}
	}
}

func TestCompressionRatio(t *testing.T) {
	tests := []struct {
		orig, comp int64
		want       string
	}{
		{1000, 700, "(-30.0%)"},
		{1000, 1000, "(0.0%)"},
		{1000, 1300, "(+30.0%)"},
		{0, 100, "0%"},
	}
	for _, tt := range tests {
		got := CompressionRatio(tt.orig, tt.comp)
		if got != tt.want {
			t.Errorf("CompressionRatio(%d, %d) = %q; want %q", tt.orig, tt.comp, got, tt.want)
		}
	}
}

func TestOutputPath(t *testing.T) {
	tests := []struct {
		input, suffix, want string
	}{
		{"/dir/video.mp4", "small", "/dir/video-small.mp4"},
		{"photo.png", "small", "photo-small.png"},
		{".hidden", "small", ".hidden-small"},
	}
	for _, tt := range tests {
		got := OutputPath(tt.input, tt.suffix)
		if got != tt.want {
			t.Errorf("OutputPath(%q, %q) = %q; want %q", tt.input, tt.suffix, got, tt.want)
		}
	}
}

func TestOutputPathExt(t *testing.T) {
	got := OutputPathExt("/dir/track.wav", "small", "m4a")
	want := "/dir/track-small.m4a"
	if got != want {
		t.Errorf("OutputPathExt = %q; want %q", got, want)
	}
}

func TestConvertOutputPath(t *testing.T) {
	got := ConvertOutputPath("photo.heic", "jpg")
	want := "photo.jpg"
	if got != want {
		t.Errorf("ConvertOutputPath = %q; want %q", got, want)
	}
}

func TestFileSize(t *testing.T) {
	f := t.TempDir() + "/test.bin"
	if err := os.WriteFile(f, []byte("hello world"), 0644); err != nil {
		t.Fatal(err)
	}
	if got := FileSize(f); got != 11 {
		t.Errorf("FileSize = %d; want 11", got)
	}
	if got := FileSize("/nonexistent"); got != 0 {
		t.Errorf("FileSize(nonexistent) = %d; want 0", got)
	}
}

func TestFileExists(t *testing.T) {
	f := t.TempDir() + "/exists.txt"
	os.WriteFile(f, []byte("x"), 0644)
	if !FileExists(f) {
		t.Errorf("FileExists(%q) = false; want true", f)
	}
	if FileExists("/nonexistent") {
		t.Errorf("FileExists(nonexistent) = true; want false")
	}
}

func TestResolveGlobs(t *testing.T) {
	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, "a.txt"), nil, 0644)
	os.WriteFile(filepath.Join(dir, "b.txt"), nil, 0644)
	os.WriteFile(filepath.Join(dir, "c.go"), nil, 0644)

	tests := []struct {
		args    []string
		want    int
	}{
		{[]string{filepath.Join(dir, "*.txt")}, 2},
		{[]string{filepath.Join(dir, "*.go")}, 1},
		{[]string{filepath.Join(dir, "*.py")}, 0},
		{[]string{filepath.Join(dir, "a.txt")}, 1},
	}
	for _, tt := range tests {
		got := ResolveGlobs(tt.args, false)
		if len(got) != tt.want {
			t.Errorf("ResolveGlobs(%v) = %d results; want %d", tt.args, len(got), tt.want)
		}
	}
}

func TestEnsureDir(t *testing.T) {
	base := t.TempDir()
	path := filepath.Join(base, "a", "b", "c", "file.txt")
	if err := EnsureDir(path); err != nil {
		t.Fatalf("EnsureDir: %v", err)
	}
	if _, err := os.Stat(filepath.Dir(path)); os.IsNotExist(err) {
		t.Error("EnsureDir did not create parent directory")
	}
}

func TestMaybeInPlace(t *testing.T) {
	dir := t.TempDir()
	orig := filepath.Join(dir, "original.txt")
	comp := filepath.Join(dir, "compressed.txt")
	os.WriteFile(orig, nil, 0644)
	MaybeInPlace(comp, orig)
	if _, err := os.Stat(orig); os.IsNotExist(err) {
		t.Error("MaybeInPlace should not rename when compressed doesn't exist")
	}
}


