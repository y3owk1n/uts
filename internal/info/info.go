// Package info provides file information display functionality.
//
//nolint:mnd
package info

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
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
	// JSON prints one machine-readable array instead of styled panels.
	JSON bool
}

// Report is everything uts knows about one file. It is the JSON schema of
// "uts info --json"; fields that do not apply to a file are omitted.
type Report struct {
	Path      string          `json:"path"`
	Name      string          `json:"name"`
	SizeBytes int64           `json:"sizeBytes"`
	Size      string          `json:"size"`
	Ext       string          `json:"ext"`
	Category  format.Category `json:"category"`
	Tool      string          `json:"tool"`

	DurationSeconds float64 `json:"durationSeconds,omitempty"`
	Duration        string  `json:"duration,omitempty"`
	Width           int     `json:"width,omitempty"`
	Height          int     `json:"height,omitempty"`
	FrameRate       string  `json:"frameRate,omitempty"`
	VideoCodec      string  `json:"videoCodec,omitempty"`
	AudioCodec      string  `json:"audioCodec,omitempty"`
	SampleRate      string  `json:"sampleRate,omitempty"`
	Channels        int     `json:"channels,omitempty"`
	BitRateKbps     int64   `json:"bitRateKbps,omitempty"`
	Pages           int     `json:"pages,omitempty"`

	Error string `json:"error,omitempty"`
}

// Show displays information about the given files. It returns an error when
// any file could not be read so the CLI exits non-zero.
func Show(opts Options) error {
	reports := make([]Report, 0, len(opts.Files))
	failed := 0

	for _, file := range opts.Files {
		report := Collect(file)
		if report.Error != "" {
			failed++
		}

		reports = append(reports, report)
	}

	if opts.JSON {
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")

		err := enc.Encode(reports)
		if err != nil {
			return derrors.Wrap(err, derrors.CodeInvalidInput, "encode JSON")
		}
	} else {
		ui.PrintBanner(opts.Version)

		for _, report := range reports {
			render(report)
		}
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

// Collect gathers the report for one file. Media details come from ffprobe
// and PDF page counts from pdfinfo or Ghostscript; both are best effort.
func Collect(file string) Report {
	report := Report{Path: file, Name: filepath.Base(file)}

	fileInfo, err := os.Stat(file)
	if err != nil {
		report.Error = "cannot access"

		return report
	}

	if fileInfo.IsDir() {
		report.Error = "not a file"

		return report
	}

	report.Ext = format.Ext(file)
	report.Category = format.Classify(report.Ext)
	report.Tool = toolHint(report.Ext)
	report.SizeBytes = fileInfo.Size()
	report.Size = util.HumanSize(fileInfo.Size())

	switch report.Category {
	case format.Video, format.Audio, format.Image:
		probeInto(&report, file)
	case format.PDF:
		report.Pages = pdfPages(file)
	case format.Archive, format.Unknown:
	}

	return report
}

func probeInto(report *Report, file string) {
	probe, err := ffmpeg.Probe(file)
	if err != nil {
		return
	}

	if dur := probe.Duration(); dur > 0 && report.Category != format.Image {
		report.DurationSeconds = dur
		report.Duration = formatDuration(dur)
	}

	if video := probe.Video(); video != nil {
		report.VideoCodec = video.Codec
		report.Width = video.Width
		report.Height = video.Height

		if report.Category == format.Video {
			report.FrameRate = frameRate(video.FrameRate)
		}
	}

	if audio := probe.Audio(); audio != nil {
		report.AudioCodec = audio.Codec
		report.SampleRate = audio.SampleRate
		report.Channels = audio.Channels
	}

	if rate := probe.BitRate(); rate > 0 && report.Category != format.Image {
		report.BitRateKbps = rate / 1000
	}
}

type row struct{ key, value string }

func render(report Report) {
	if report.Error != "" {
		switch report.Error {
		case "not a file":
			ui.Message.Warnf("Not a file: %s", report.Path)
		default:
			ui.Message.Warnf("Cannot access: %s", report.Path)
		}

		return
	}

	palette := ui.Style.Palette()

	catColor := palette.Accent
	if report.Category == format.Unknown {
		catColor = palette.Warning
	}

	rows := make([]row, 0, 9)
	rows = append(rows,
		row{"Size:", ui.Message.Accent(report.Size)},
		row{"Type:", ui.Message.Accent("." + report.Ext)},
		row{"Category:", lipgloss.NewStyle().Foreground(catColor).Render(string(report.Category))},
		row{"Tool:", ui.Message.Accent(report.Tool)},
	)
	rows = append(rows, detailRows(report)...)

	var body strings.Builder
	for _, r := range rows {
		body.WriteString("  " + keyStyle(palette, r.key) + "  " + r.value + "\n")
	}

	if suggestions := suggestActions(report.Category, report.Path); suggestions != "" {
		body.WriteString(
			"\n" + lipgloss.NewStyle().
				Foreground(palette.Subtle).
				Render("Suggestions") +
				"\n" + suggestions,
		)
	}

	_, _ = lipgloss.Fprint(
		os.Stdout,
		ui.Panel.Section(ui.Message.Highlight(report.Name), body.String()),
	)
}

func detailRows(report Report) []row {
	rows := make([]row, 0, 5)

	if report.Duration != "" {
		rows = append(rows, row{"Duration:", ui.Message.Accent(report.Duration)})
	}

	if report.Width > 0 {
		res := fmt.Sprintf("%dx%d", report.Width, report.Height)
		if report.FrameRate != "" {
			res += " @ " + report.FrameRate + " fps"
		}

		rows = append(rows, row{"Resolution:", ui.Message.Accent(res)})
	}

	var codecs []string

	if report.VideoCodec != "" {
		codecs = append(codecs, report.VideoCodec)
	}

	if report.AudioCodec != "" {
		desc := report.AudioCodec
		if report.SampleRate != "" {
			desc += " " + report.SampleRate + " Hz"
		}

		if report.Channels > 0 {
			desc += fmt.Sprintf(" %dch", report.Channels)
		}

		codecs = append(codecs, desc)
	}

	if len(codecs) > 0 {
		rows = append(rows, row{"Codec:", ui.Message.Accent(strings.Join(codecs, ", "))})
	}

	if report.BitRateKbps > 0 {
		rows = append(
			rows,
			row{"Bitrate:", ui.Message.Accent(fmt.Sprintf("%d kb/s", report.BitRateKbps))},
		)
	}

	if report.Pages > 0 {
		rows = append(rows, row{"Pages:", ui.Message.Accent(strconv.Itoa(report.Pages))})
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
