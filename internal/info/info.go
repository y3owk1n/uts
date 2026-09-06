// Package info provides file information display functionality.
//
//nolint:mnd
package info

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"charm.land/lipgloss/v2"
	derrors "github.com/y3owk1n/uts/internal/core/errors"
	"github.com/y3owk1n/uts/internal/ffmpeg"
	"github.com/y3owk1n/uts/internal/format"
	"github.com/y3owk1n/uts/internal/ui"
	"github.com/y3owk1n/uts/internal/ui/style"
	"github.com/y3owk1n/uts/internal/util"
)

// Options represents options for displaying file information.
type Options struct {
	Files   []string
	Version string
}

// Show displays information about the given files. It returns an error when
// any file could not be read so the CLI exits non-zero.
func Show(opts Options) error {
	ui.PrintBanner(opts.Version)

	palette := ui.Style.Palette()
	failed := 0

	for _, file := range opts.Files {
		fileInfo, err := os.Stat(file)
		if err != nil {
			ui.Message.Warnf("Cannot access: %s", file)

			failed++

			continue
		}

		if fileInfo.IsDir() {
			ui.Message.Warnf("Not a file: %s", file)

			failed++

			continue
		}

		ext := format.Ext(file)
		cat := format.Classify(ext)

		catColor := palette.Accent
		if cat == format.Unknown {
			catColor = palette.Warning
		}

		rows := make([]row, 0, 8)
		rows = append(rows, []row{
			{"Size:", ui.Message.Accent(util.HumanSize(fileInfo.Size()))},
			{"Type:", ui.Message.Accent("." + ext)},
			{"Category:", lipgloss.NewStyle().Foreground(catColor).Render(string(cat))},
			{"Tool:", ui.Message.Accent(toolHint(ext))},
		}...)
		rows = append(rows, mediaRows(file, cat)...)

		var body strings.Builder
		for _, r := range rows {
			body.WriteString("  " + keyStyle(palette, r.key) + "  " + r.value + "\n")
		}

		if suggestions := suggestActions(cat, file); suggestions != "" {
			body.WriteString(
				"\n" + lipgloss.NewStyle().
					Foreground(palette.Subtle).
					Render("Suggestions") +
					"\n" + suggestions,
			)
		}

		_, _ = lipgloss.Fprint(
			os.Stdout,
			ui.Panel.Section(ui.Message.Highlight(filepath.Base(file)), body.String()),
		)
	}

	if failed > 0 {
		return derrors.Newf(
			derrors.CodeFileNotFound,
			"%d of %d files could not be read",
			failed,
			len(opts.Files),
		)
	}

	return nil
}

type row struct{ key, value string }

// mediaRows adds ffprobe-derived details for video, audio and image files.
// It is best effort: no ffprobe, or a file ffprobe cannot read, adds nothing.
func mediaRows(file string, cat format.Category) []row {
	if cat != format.Video && cat != format.Audio && cat != format.Image {
		return nil
	}

	probe, err := ffmpeg.Probe(file)
	if err != nil {
		return nil
	}

	rows := make([]row, 0, 4)

	if dur := probe.Duration(); dur > 0 && cat != format.Image {
		rows = append(rows, row{"Duration:", ui.Message.Accent(formatDuration(dur))})
	}

	if video := probe.Video(); video != nil && video.Width > 0 {
		res := fmt.Sprintf("%dx%d", video.Width, video.Height)
		if cat == format.Video {
			if fps := frameRate(video.FrameRate); fps != "" {
				res += " @ " + fps + " fps"
			}
		}

		rows = append(rows, row{"Resolution:", ui.Message.Accent(res)})
	}

	var codecs []string

	if video := probe.Video(); video != nil {
		codecs = append(codecs, video.Codec)
	}

	if audio := probe.Audio(); audio != nil {
		desc := audio.Codec
		if audio.SampleRate != "" {
			desc += " " + audio.SampleRate + " Hz"
		}

		if audio.Channels > 0 {
			desc += fmt.Sprintf(" %dch", audio.Channels)
		}

		codecs = append(codecs, desc)
	}

	if len(codecs) > 0 {
		rows = append(rows, row{"Codec:", ui.Message.Accent(strings.Join(codecs, ", "))})
	}

	if rate := probe.BitRate(); rate > 0 && cat != format.Image {
		rows = append(rows, row{"Bitrate:", ui.Message.Accent(fmt.Sprintf("%d kb/s", rate/1000))})
	}

	return rows
}

func formatDuration(seconds float64) string {
	dur := time.Duration(seconds * float64(time.Second)).Round(time.Second)

	hours := int(dur.Hours())
	minutes := int(dur.Minutes()) % 60
	secs := int(dur.Seconds()) % 60

	if hours > 0 {
		return fmt.Sprintf("%d:%02d:%02d", hours, minutes, secs)
	}

	return fmt.Sprintf("%d:%02d", minutes, secs)
}

// frameRate turns ffprobe's "30000/1001" into "29.97".
func frameRate(ratio string) string {
	num, den := 0.0, 0.0

	_, err := fmt.Sscanf(ratio, "%f/%f", &num, &den)
	if err != nil || den == 0 || num == 0 {
		return ""
	}

	return strings.TrimSuffix(strings.TrimSuffix(fmt.Sprintf("%.2f", num/den), "0"), ".0")
}

func keyStyle(palette style.Palette, key string) string {
	return lipgloss.NewStyle().
		Foreground(palette.Muted).
		Width(12).
		Align(lipgloss.Right).
		Render(key)
}

func toolHint(ext string) string {
	switch format.Classify(ext) {
	case format.Video:
		vcodec, _ := format.VideoCodecs(ext)

		return "ffmpeg (" + vcodec + ")"
	case format.Image:
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
			return "heif-convert (HEIC → JPEG)"
		default:
			return "ImageMagick"
		}
	case format.Audio:
		codec, _ := format.AudioCompressTarget(ext)

		return "ffmpeg (" + codec + ")"
	case format.PDF:
		return "ghostscript"
	case format.Archive:
		switch ext {
		case "zip":
			return "zip / unzip"
		case "zst", "zstd":
			return "tar + zstd"
		case "br":
			return "tar + brotli"
		default:
			return "tar"
		}
	case format.Unknown:
		return "—"
	}

	return "—"
}

func suggestActions(cat format.Category, file string) string {
	var lines []string

	switch cat {
	case format.Video:
		lines = append(
			lines,
			detail(
				"Compress",
				fmt.Sprintf("uts video compress %q [-q low|medium|high|<0-51>]", file),
			),
			detail("Convert", fmt.Sprintf("uts video convert %q --to mp4", file)),
			detail("Audio", fmt.Sprintf("uts audio convert %q --to mp3", file)),
		)
	case format.Image:
		lines = append(
			lines,
			detail(
				"Compress",
				fmt.Sprintf("uts image compress %q [-q low|medium|high|<1-100>]", file),
			),
			detail("Convert", fmt.Sprintf("uts image convert %q --to webp", file)),
		)
	case format.Audio:
		lines = append(
			lines,
			detail(
				"Compress",
				fmt.Sprintf("uts audio compress %q [-q low|medium|high|<kbps>]", file),
			),
			detail("Convert", fmt.Sprintf("uts audio convert %q --to mp3", file)),
		)
	case format.PDF:
		lines = append(lines,
			detail("Compress", fmt.Sprintf("uts pdf compress %q [-q low|medium|high|<dpi>]", file)),
			detail("Convert", fmt.Sprintf("uts pdf convert %q --to jpg", file)))
	case format.Archive:
		lines = append(lines,
			detail("Extract", fmt.Sprintf("uts archive extract %q", file)),
			detail("List", fmt.Sprintf("uts archive list %q", file)))
	case format.Unknown:
	}

	return strings.Join(lines, "")
}

func detail(label, cmd string) string {
	palette := ui.Style.Palette()
	labelStyle := lipgloss.NewStyle().
		Foreground(palette.Accent).
		Width(12).
		Align(lipgloss.Right).
		Render
	cmdStyle := lipgloss.NewStyle().Foreground(palette.Subtle).Render

	return "    " + labelStyle(label+":") + "  " + cmdStyle(cmd) + "\n"
}
