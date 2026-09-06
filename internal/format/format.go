// Package format is the single source of truth for the file formats uts
// understands: which extensions belong to which category and which codecs
// each container expects.
//
//nolint:goconst
package format

import (
	"path/filepath"
	"slices"
	"strings"
)

// Category is a media category handled by a uts command group.
type Category string

// Categories.
const (
	Image   Category = "image"
	Video   Category = "video"
	Audio   Category = "audio"
	PDF     Category = "pdf"
	Archive Category = "archive"
	Unknown Category = "unknown"
)

// Extension sets per category, lowercase without the leading dot.
var (
	ImageExts = []string{
		"png",
		"jpg",
		"jpeg",
		"webp",
		"gif",
		"bmp",
		"tiff",
		"tif",
		"heic",
		"heif",
		"avif",
		"avifs",
	}
	VideoExts   = []string{"mp4", "mov", "mkv", "avi", "webm", "m4v", "flv", "wmv"}
	AudioExts   = []string{"wav", "flac", "aac", "mp3", "m4a", "opus", "ogg", "wma"}
	PDFExts     = []string{"pdf"}
	ArchiveExts = []string{
		"zip",
		"tar",
		"gz",
		"tgz",
		"zst",
		"zstd",
		"xz",
		"txz",
		"bz2",
		"tbz2",
		"br",
		"7z",
	}

	// ImageTargets, VideoTargets, AudioTargets and PDFTargets are the
	// conversion targets accepted by each convert command.
	ImageTargets = []string{"jpg", "jpeg", "png", "webp", "gif", "bmp", "tiff", "tif", "avif"}
	VideoTargets = []string{"mp4", "mkv", "webm", "mov", "avi", "flv"}
	AudioTargets = []string{"mp3", "aac", "m4a", "wav", "flac", "opus", "ogg"}
	PDFTargets   = []string{"jpg", "jpeg", "png", "pdf"}
)

// Ext returns the lowercase extension of path without the leading dot.
func Ext(path string) string {
	return strings.ToLower(strings.TrimPrefix(filepath.Ext(path), "."))
}

// Classify returns the category an extension belongs to.
func Classify(ext string) Category {
	ext = strings.ToLower(strings.TrimPrefix(ext, "."))

	switch {
	case slices.Contains(ImageExts, ext):
		return Image
	case slices.Contains(VideoExts, ext):
		return Video
	case slices.Contains(AudioExts, ext):
		return Audio
	case slices.Contains(PDFExts, ext):
		return PDF
	case slices.Contains(ArchiveExts, ext):
		return Archive
	}

	return Unknown
}

// Exts returns the input extensions a category accepts.
func Exts(cat Category) []string {
	switch cat {
	case Image:
		return ImageExts
	case Video:
		return VideoExts
	case Audio:
		return AudioExts
	case PDF:
		return PDFExts
	case Archive:
		return ArchiveExts
	case Unknown:
		return nil
	}

	return nil
}

// Normalize maps alias extensions to their canonical form (jpeg -> jpg, tif -> tiff).
func Normalize(ext string) string {
	switch strings.ToLower(ext) {
	case "jpeg":
		return "jpg"
	case "tif":
		return "tiff"
	case "heif":
		return "heic"
	case "zstd":
		return "zst"
	}

	return strings.ToLower(ext)
}

// Same reports whether two extensions denote the same format.
func Same(a, b string) bool {
	return Normalize(a) == Normalize(b)
}

// VideoCodecs returns the encoder pair (video, audio) uts uses for a container.
func VideoCodecs(ext string) (string, string) {
	switch Normalize(strings.TrimPrefix(ext, ".")) {
	case "mkv":
		return "libx265", "aac"
	case "webm":
		return "libvpx-vp9", "libopus"
	case "avi":
		return "libx264", "mp3"
	}

	return "libx264", "aac"
}

// ContainerAccepts reports whether a container can hold the given ffprobe
// codec name without re-encoding. It is deliberately conservative: only
// well-supported pairings return true so a stream copy never yields an
// unplayable file.
func ContainerAccepts(container, codec string) bool {
	codec = strings.ToLower(codec)

	switch Normalize(strings.TrimPrefix(container, ".")) {
	case "mp4", "mov", "m4v":
		return slices.Contains(
			[]string{"h264", "hevc", "av1", "mpeg4", "aac", "mp3", "alac", "ac3"},
			codec,
		)
	case "mkv":
		return slices.Contains([]string{
			"h264", "hevc", "av1", "vp8", "vp9", "mpeg4", "mpeg2video",
			"aac", "mp3", "opus", "vorbis", "flac", "ac3", "eac3", "dts", "pcm_s16le",
		}, codec)
	case "webm":
		return slices.Contains([]string{"vp8", "vp9", "av1", "opus", "vorbis"}, codec)
	case "avi":
		return slices.Contains([]string{"h264", "mpeg4", "mp3", "ac3", "pcm_s16le"}, codec)
	case "flv":
		return slices.Contains([]string{"h264", "aac", "mp3"}, codec)
	}

	return false
}

// AudioCodec returns the ffmpeg encoder and output extension for an audio
// target. The second return is empty when the target is unsupported.
func AudioCodec(target string) (string, string) {
	switch strings.ToLower(target) {
	case "mp3":
		return "libmp3lame", "mp3"
	case "aac", "m4a":
		return "aac", "m4a"
	case "wav":
		return "pcm_s16le", "wav"
	case "flac":
		return "flac", "flac"
	case "opus":
		return "libopus", "opus"
	case "ogg":
		return "libvorbis", "ogg"
	}

	return "", ""
}

// Lossless reports whether an audio extension is a lossless format, where a
// bitrate setting is meaningless.
func Lossless(ext string) bool {
	switch strings.ToLower(strings.TrimPrefix(ext, ".")) {
	case "wav", "flac", "alac", "aiff", "aif":
		return true
	}

	return false
}

// AudioCompressTarget picks the output extension and encoder used when
// compressing an audio file. Lossy inputs keep their container so the result
// is a smaller file of the same kind. Lossless and exotic inputs become AAC in
// an .m4a container, the most widely playable lossy format.
func AudioCompressTarget(ext string) (string, string) {
	switch strings.ToLower(strings.TrimPrefix(ext, ".")) {
	case "mp3":
		return "libmp3lame", "mp3"
	case "ogg":
		return "libvorbis", "ogg"
	case "opus":
		return "libopus", "opus"
	}

	return "aac", "m4a"
}
