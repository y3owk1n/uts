//nolint:testpackage,goconst
package util

import (
	"os"
	"path/filepath"
	"testing"
)

// TestHumanSize tests HumanSize.
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
	for _, tc := range tests {
		got := HumanSize(tc.input)
		if got != tc.want {
			t.Errorf("HumanSize(%d) = %q; want %q", tc.input, got, tc.want)
		}
	}
}

// TestCompressionRatio tests CompressionRatio.
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
	for _, tc := range tests {
		got := CompressionRatio(tc.orig, tc.comp)
		if got != tc.want {
			t.Errorf("CompressionRatio(%d, %d) = %q; want %q", tc.orig, tc.comp, got, tc.want)
		}
	}
}

// TestOutputPath tests OutputPath.
func TestOutputPath(t *testing.T) {
	tests := []struct {
		input, suffix, want string
	}{
		{"/dir/video.mp4", "small", "/dir/video-small.mp4"},
		{"photo.png", "small", "photo-small.png"},
		{".hidden", "small", ".hidden-small"},
	}
	for _, tc := range tests {
		got := OutputPath(tc.input, tc.suffix)
		if got != tc.want {
			t.Errorf("OutputPath(%q, %q) = %q; want %q", tc.input, tc.suffix, got, tc.want)
		}
	}
}

// TestOutputPathExt tests OutputPathExt.
func TestOutputPathExt(t *testing.T) {
	got := OutputPathExt("/dir/track.wav", "small", "m4a")

	want := "/dir/track-small.m4a"
	if got != want {
		t.Errorf("OutputPathExt = %q; want %q", got, want)
	}
}

// TestConvertOutputPath tests ConvertOutputPath.
func TestConvertOutputPath(t *testing.T) {
	got := ConvertOutputPath("photo.heic", "jpg")

	want := "photo.jpg"
	if got != want {
		t.Errorf("ConvertOutputPath = %q; want %q", got, want)
	}
}

// TestFileSize tests FileSize.
func TestFileSize(t *testing.T) {
	file := t.TempDir() + "/test.bin"

	err := os.WriteFile(file, []byte("hello world"), 0o644)
	if err != nil {
		t.Fatal(err)
	}

	if got := FileSize(file); got != 11 {
		t.Errorf("FileSize = %d; want 11", got)
	}

	if got := FileSize("/nonexistent"); got != 0 {
		t.Errorf("FileSize(nonexistent) = %d; want 0", got)
	}
}

// TestFileExists tests FileExists.
func TestFileExists(t *testing.T) {
	file := t.TempDir() + "/exists.txt"
	//nolint:errcheck
	os.WriteFile(file, []byte("x"), 0o644)

	if !FileExists(file) {
		t.Errorf("FileExists(%q) = false; want true", file)
	}

	if FileExists("/nonexistent") {
		t.Errorf("FileExists(nonexistent) = true; want false")
	}
}

// TestResolveGlobs tests ResolveGlobs.
func TestResolveGlobs(t *testing.T) {
	dir := t.TempDir()
	//nolint:errcheck
	os.WriteFile(filepath.Join(dir, "a.txt"), nil, 0o644)
	//nolint:errcheck
	os.WriteFile(filepath.Join(dir, "b.txt"), nil, 0o644)
	//nolint:errcheck
	os.WriteFile(filepath.Join(dir, "c.go"), nil, 0o644)

	tests := []struct {
		args []string
		want int
	}{
		{[]string{filepath.Join(dir, "*.txt")}, 2},
		{[]string{filepath.Join(dir, "*.go")}, 1},
		{[]string{filepath.Join(dir, "*.py")}, 0},
		{[]string{filepath.Join(dir, "a.txt")}, 1},
	}
	for _, tc := range tests {
		got := ResolveGlobs(tc.args, false)
		if len(got) != tc.want {
			t.Errorf("ResolveGlobs(%v) = %d results; want %d", tc.args, len(got), tc.want)
		}
	}
}

// TestEnsureDir tests EnsureDir.
func TestEnsureDir(t *testing.T) {
	base := t.TempDir()

	path := filepath.Join(base, "a", "b", "c", "file.txt")

	err := EnsureDir(path)
	if err != nil {
		t.Fatalf("EnsureDir: %v", err)
	}

	_, err = os.Stat(filepath.Dir(path))

	if os.IsNotExist(err) {
		t.Error("EnsureDir did not create parent directory")
	}
}

// TestMaybeInPlace tests MaybeInPlace.
func TestMaybeInPlace(t *testing.T) {
	// When compressed doesn't exist, original should remain untouched.
	dir := t.TempDir()
	orig := filepath.Join(dir, "original.txt")
	comp := filepath.Join(dir, "compressed.txt")

	//nolint:errcheck
	os.WriteFile(orig, []byte("original"), 0o644)
	MaybeInPlace(comp, orig)

	data, err := os.ReadFile(orig)
	if err != nil {
		t.Fatal(err)
	}

	if string(data) != "original" {
		t.Errorf("MaybeInPlace without compressed: got %q; want \"original\"", string(data))
	}
}

// TestMaybeInPlaceRename tests that MaybeInPlace renames the compressed file to the original.
func TestMaybeInPlaceRename(t *testing.T) {
	dir := t.TempDir()
	orig := filepath.Join(dir, "video.mp4")
	comp := filepath.Join(dir, "video-small.mp4")

	//nolint:errcheck
	os.WriteFile(orig, []byte("old"), 0o644)
	//nolint:errcheck
	os.WriteFile(comp, []byte("compressed"), 0o644)

	MaybeInPlace(comp, orig)

	if FileExists(comp) {
		t.Error("MaybeInPlace should remove the compressed file after rename")
	}

	data, err := os.ReadFile(orig)
	if err != nil {
		t.Fatal(err)
	}

	if string(data) != "compressed" {
		t.Errorf("MaybeInPlace: original should contain compressed content; got %q", string(data))
	}
}

// TestRemoveInPlace tests RemoveInPlace.
func TestRemoveInPlace(t *testing.T) {
	dir := t.TempDir()
	orig := filepath.Join(dir, "video.mov")

	//nolint:errcheck
	os.WriteFile(orig, []byte("old"), 0o644)

	RemoveInPlace(orig)

	if FileExists(orig) {
		t.Error("RemoveInPlace should delete the original file")
	}
}

// TestRemoveInPlaceNonexistent tests that RemoveInPlace doesn't panic on missing files.
func TestRemoveInPlaceNonexistent(t *testing.T) {
	// Should not panic even if the file doesn't exist.
	RemoveInPlace("/nonexistent/file.txt")
}

// TestInPlaceHint tests InPlaceHint.
func TestInPlaceHint(t *testing.T) {
	if got := InPlaceHint(true); got != " (in-place)" {
		t.Errorf("InPlaceHint(true) = %q; want \" (in-place)\"", got)
	}

	if got := InPlaceHint(false); got != "" {
		t.Errorf("InPlaceHint(false) = %q; want \"\"", got)
	}
}
