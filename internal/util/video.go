//nolint:goconst
package util

import "strings"

// VideoCodecs returns the recommended video and audio codec for a given
// container format extension (e.g. ".mp4", ".mov", ".webm").
func VideoCodecs(ext string) (string, string) {
	switch strings.ToLower(strings.TrimPrefix(ext, ".")) {
	case "mp4", "mov":
		return "libx264", "aac"
	case "mkv":
		return "libx265", "aac"
	case "webm":
		return "libvpx-vp9", "libopus"
	case "avi":
		return "libx264", "mp3"
	case "flv":
		return "libx264", "aac"
	}

	return "libx264", "aac"
}
