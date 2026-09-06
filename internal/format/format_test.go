//nolint:goconst
package format_test

import (
	"testing"

	"github.com/y3owk1n/uts/internal/format"
)

func TestClassify(t *testing.T) {
	tests := []struct {
		ext  string
		want format.Category
	}{
		{"png", format.Image},
		{".JPG", format.Image},
		{"mov", format.Video},
		{"flac", format.Audio},
		{"pdf", format.PDF},
		{"tgz", format.Archive},
		{"br", format.Archive},
		{"txt", format.Unknown},
		{"", format.Unknown},
	}
	for _, testCase := range tests {
		if got := format.Classify(testCase.ext); got != testCase.want {
			t.Errorf("Classify(%q) = %q; want %q", testCase.ext, got, testCase.want)
		}
	}
}

func TestExt(t *testing.T) {
	tests := []struct{ path, want string }{
		{"photo.HEIC", "heic"},
		{"/a/b/clip.tar.gz", "gz"},
		{"noext", ""},
		{".hidden", "hidden"},
	}
	for _, testCase := range tests {
		if got := format.Ext(testCase.path); got != testCase.want {
			t.Errorf("Ext(%q) = %q; want %q", testCase.path, got, testCase.want)
		}
	}
}

func TestSame(t *testing.T) {
	if !format.Same("jpeg", "jpg") || !format.Same("TIF", "tiff") || format.Same("png", "jpg") {
		t.Error("Same: alias handling is wrong")
	}
}

func TestAudioCompressTarget(t *testing.T) {
	tests := []struct{ ext, codec, out string }{
		{"mp3", "libmp3lame", "mp3"},
		{"ogg", "libvorbis", "ogg"},
		{"opus", "libopus", "opus"},
		{"wav", "aac", "m4a"},
		{"flac", "aac", "m4a"},
		{"m4a", "aac", "m4a"},
		{"wma", "aac", "m4a"},
	}
	for _, testCase := range tests {
		codec, out := format.AudioCompressTarget(testCase.ext)
		if codec != testCase.codec || out != testCase.out {
			t.Errorf(
				"AudioCompressTarget(%q) = (%q, %q); want (%q, %q)",
				testCase.ext,
				codec,
				out,
				testCase.codec,
				testCase.out,
			)
		}
	}
}

func TestContainerAccepts(t *testing.T) {
	tests := []struct {
		container, codec string
		want             bool
	}{
		{"mp4", "h264", true},
		{"mp4", "aac", true},
		{"mp4", "vp9", false},
		{"mp4", "pcm_s16le", false},
		{"webm", "vp9", true},
		{"webm", "h264", false},
		{"mkv", "h264", true},
		{"mkv", "opus", true},
		{"flv", "hevc", false},
		{"gif", "h264", false},
	}
	for _, testCase := range tests {
		if got := format.ContainerAccepts(
			testCase.container,
			testCase.codec,
		); got != testCase.want {
			t.Errorf(
				"ContainerAccepts(%q, %q) = %v; want %v",
				testCase.container,
				testCase.codec,
				got,
				testCase.want,
			)
		}
	}
}

func TestVideoCodecs(t *testing.T) {
	tests := []struct{ ext, v, a string }{
		{".mp4", "libx264", "aac"},
		{"mov", "libx264", "aac"},
		{"mkv", "libx265", "aac"},
		{"webm", "libvpx-vp9", "libopus"},
		{"avi", "libx264", "mp3"},
		{"flv", "libx264", "aac"},
		{"unknown", "libx264", "aac"},
	}
	for _, testCase := range tests {
		video, audio := format.VideoCodecs(testCase.ext)
		if video != testCase.v || audio != testCase.a {
			t.Errorf(
				"VideoCodecs(%q) = (%q, %q); want (%q, %q)",
				testCase.ext,
				video,
				audio,
				testCase.v,
				testCase.a,
			)
		}
	}
}
