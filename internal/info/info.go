// Package info provides file information display functionality.
//
//nolint:goconst,mnd
package info

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"charm.land/lipgloss/v2"
	"github.com/y3owk1n/uts/internal/ui"
	"github.com/y3owk1n/uts/internal/ui/style"
	"github.com/y3owk1n/uts/internal/util"
)

// Options represents options for displaying file information.
type Options struct {
	Files   []string
	Version string
}

// Show displays information about the given files.
func Show(opts Options) {
	_, _ = lipgloss.Fprint(os.Stdout, ui.Banner.Logo(opts.Version))

	palette := ui.Style.Palette()

	for _, file := range opts.Files {
		fileInfo, err := os.Stat(file)
		if err != nil {
			ui.Message.Warnf("Cannot access: %s", file)

			continue
		}

		if fileInfo.IsDir() {
			ui.Message.Warnf("Not a file: %s", file)

			continue
		}

		ext := strings.ToLower(strings.TrimPrefix(filepath.Ext(file), "."))
		size := util.FileSize(file)
		base := filepath.Base(file)

		catColor := palette.Accent

		category := classify(ext)
		if category == "unknown" {
			catColor = palette.Warning
		}

		body := fmt.Sprintf(
			"  %s  %s\n  %s  %s\n  %s  %s\n  %s  %s\n",
			keyStyle(palette, "Size:"),
			ui.Message.Accent(util.HumanSize(size)),
			keyStyle(palette, "Type:"),
			ui.Message.Accent("."+ext),
			keyStyle(palette, "Category:"),
			lipgloss.NewStyle().Foreground(catColor).Render(category),
			keyStyle(palette, "Tool:"),
			ui.Message.Accent(toolHint(ext)),
		)

		suggestions := suggestActions(ext, file)
		if suggestions != "" {
			body += "\n" + lipgloss.NewStyle().
				Foreground(palette.Subtle).
				Render("Suggestions") +
				"\n" + suggestions
		}

		_, _ = lipgloss.Fprint(os.Stdout, ui.Panel.Section(ui.Message.Highlight(base), body))
	}
}

func keyStyle(palette style.Palette, key string) string {
	return lipgloss.NewStyle().
		Foreground(palette.Muted).
		Width(10).
		Align(lipgloss.Right).
		Render(key)
}

func toolHint(ext string) string {
	switch {
	case isVideo(ext):
		return "ffmpeg (libx265)"
	case isImage(ext):
		switch ext {
		case "png":
			return "pngquant + optipng"
		case "jpg", "jpeg":
			return "jpegoptim"
		case "webp":
			return "cwebp"
		case "gif":
			return "gifsicle"
		case "heic", "heif":
			return "ImageMagick (HEIC \u2192 JPEG)"
		case "avif", "avifs":
			return "cavif / libavif"
		default:
			return "ImageMagick"
		}
	case isAudio(ext):
		return "ffmpeg (aac)"
	case ext == "pdf":
		return "ghostscript"
	case isArchive(ext):
		switch ext {
		case "zip":
			return "unzip"
		case "gz", "tgz":
			return "tar (gzip)"
		case "zst", "zstd":
			return "tar (zstd)"
		case "xz", "txz":
			return "tar (xz)"
		case "bz2", "tbz2":
			return "tar (bzip2)"
		default:
			return "tar"
		}
	}

	return "\u2014"
}

func classify(ext string) string {
	switch {
	case isVideo(ext):
		return "video"
	case isImage(ext):
		return "image"
	case isAudio(ext):
		return "audio"
	case ext == "pdf":
		return "pdf"
	case isArchive(ext):
		return "archive"
	}

	return "unknown"
}

func suggestActions(ext, file string) string {
	var lines []string

	switch {
	case isVideo(ext):
		lines = append(
			lines,
			detail(
				"Compress",
				fmt.Sprintf("uts video compress %q [-q low|medium|high|<0-51>]", file),
			),
		)
		lines = append(lines, detail("Convert", fmt.Sprintf("uts video convert %q --to mkv", file)))
	case isImage(ext):
		lines = append(
			lines,
			detail(
				"Compress",
				fmt.Sprintf("uts image compress %q [-q low|medium|high|<1-100>]", file),
			),
		)
		lines = append(
			lines,
			detail("Convert", fmt.Sprintf("uts image convert %q --to webp", file)),
		)
	case isAudio(ext):
		lines = append(
			lines,
			detail(
				"Compress",
				fmt.Sprintf("uts audio compress %q [-q low|medium|high|<kbps>]", file),
			),
		)
		lines = append(lines, detail("Convert", fmt.Sprintf("uts audio convert %q --to mp3", file)))
	case ext == "pdf":
		lines = append(
			lines,
			detail("Compress", fmt.Sprintf("uts pdf compress %q [-q low|medium|high|<dpi>]", file)),
		)
		lines = append(lines, detail("Convert", fmt.Sprintf("uts pdf convert %q --to jpg", file)))
	case isArchive(ext):
		lines = append(lines, detail("Extract", fmt.Sprintf("uts archive extract %q", file)))
		lines = append(lines, detail("List", fmt.Sprintf("uts archive list %q", file)))
	}

	return strings.Join(lines, "")
}

func detail(label, cmd string) string {
	palette := ui.Style.Palette()
	labelStyle := lipgloss.NewStyle().
		Foreground(palette.Accent).
		Width(10).
		Align(lipgloss.Right).
		Render
	cmdStyle := lipgloss.NewStyle().Foreground(palette.Subtle).Render

	return "    " + labelStyle(label+":") + "  " + cmdStyle(cmd) + "\n"
}

func isVideo(ext string) bool {
	switch ext {
	case "mp4", "mov", "mkv", "avi", "webm", "m4v", "flv", "wmv":
		return true
	}

	return false
}

func isImage(ext string) bool {
	switch ext {
	case "png", "jpg", "jpeg", "webp", "gif", "bmp", "tiff", "tif", "heic", "heif", "avif", "avifs":
		return true
	}

	return false
}

func isAudio(ext string) bool {
	switch ext {
	case "wav", "flac", "aac", "mp3", "m4a", "opus", "ogg", "wma":
		return true
	}

	return false
}

func isArchive(ext string) bool {
	switch ext {
	case "zip", "tar", "gz", "tgz", "zst", "zstd", "xz", "txz", "bz2", "tbz2", "7z":
		return true
	}

	return false
}
