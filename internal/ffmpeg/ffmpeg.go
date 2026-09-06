// Package ffmpeg builds ffmpeg argument lists and reads ffprobe metadata so
// the video, audio and info commands agree on codecs and flags.
//
//nolint:goconst,tagliatelle // ffprobe dictates snake_case JSON keys.
package ffmpeg

import (
	"context"
	"encoding/json"
	"errors"
	"os/exec"
	"strconv"
	"strings"

	derrors "github.com/y3owk1n/uts/internal/core/errors"
	"github.com/y3owk1n/uts/internal/format"
	"github.com/y3owk1n/uts/internal/job"
	"github.com/y3owk1n/uts/internal/util"
)

const microsPerSecond = 1_000_000

// Step builds an ffmpeg job step that reports progress. The input is probed
// for its duration so the runner can show a percentage and time remaining;
// without ffprobe the step still runs, just without the numbers.
func Step(input string, args ...string) job.Step {
	full := append([]string{"-nostats", "-loglevel", "error", "-progress", "pipe:1"}, args...)
	step := job.Exec("ffmpeg", full...)
	step.Progress = &job.Progress{Total: DurationOf(input), Parse: ParseProgress}

	return step
}

// DurationOf returns the media duration in seconds, or 0 when unknown.
func DurationOf(file string) float64 {
	info, err := Probe(file)
	if err != nil {
		return 0
	}

	return info.Duration()
}

// ParseProgress reads one line of "ffmpeg -progress" output and returns the
// seconds of output written so far.
func ParseProgress(line string) (float64, bool) {
	value, ok := strings.CutPrefix(line, "out_time_us=")
	if !ok {
		// Older ffmpeg only prints out_time_ms, which is also microseconds.
		value, ok = strings.CutPrefix(line, "out_time_ms=")
		if !ok {
			return 0, false
		}
	}

	micros, err := strconv.ParseInt(strings.TrimSpace(value), 10, 64)
	if err != nil || micros < 0 {
		return 0, false
	}

	return float64(micros) / microsPerSecond, true
}

// ErrNoProbe is returned when ffprobe is not installed.
var ErrNoProbe = errors.New("ffprobe not found")

// Check returns an error when ffmpeg is not installed.
func Check() error {
	if util.HasTool("ffmpeg") {
		return nil
	}

	return derrors.New(derrors.CodeToolNotFound, "ffmpeg not found — install: brew install ffmpeg")
}

// EncodeArgs returns ffmpeg arguments that re-encode in into out for the given
// container at a CRF quality. Codec-specific flags are chosen so every
// supported container produces a playable file:
//
//   - x264/x265 get -preset, yuv420p for player compatibility, and hvc1
//     tagging so Apple players recognize HEVC in mp4/mov
//   - VP9 needs -b:v 0 for CRF to mean constant quality, and speed is set via
//     -cpu-used because it has no -preset
//   - mp4/mov get +faststart so playback can begin before download ends
func EncodeArgs(input, out, container string, crf int, preset string) []string {
	container = format.Normalize(strings.TrimPrefix(container, "."))
	vcodec, acodec := format.VideoCodecs(container)

	args := []string{
		"-i",
		input,
		"-map",
		"0:v:0",
		"-map",
		"0:a?",
		"-c:v",
		vcodec,
		"-crf",
		strconv.Itoa(crf),
	}

	switch vcodec {
	case "libvpx-vp9":
		args = append(args, "-b:v", "0", "-row-mt", "1", "-cpu-used", vp9Speed(preset))
	default:
		args = append(args, "-preset", preset, "-pix_fmt", "yuv420p")

		if vcodec == "libx265" && (container == "mp4" || container == "mov") {
			args = append(args, "-tag:v", "hvc1")
		}
	}

	args = append(args, "-c:a", acodec, "-b:a", "128k")

	if container == "mp4" || container == "mov" || container == "m4v" {
		args = append(args, "-movflags", "+faststart")
	}

	return append(args, "-y", out)
}

func vp9Speed(preset string) string {
	switch preset {
	case "slow":
		return "1"
	case "fast":
		return "4"
	}

	return "2"
}

// RemuxArgs returns ffmpeg arguments that copy every video and audio stream
// into a new container without re-encoding. Subtitle and data streams are
// dropped because their support varies per container.
func RemuxArgs(input, out, container string) []string {
	args := []string{"-i", input, "-map", "0:v", "-map", "0:a?", "-c", "copy"}

	switch format.Normalize(strings.TrimPrefix(container, ".")) {
	case "mp4", "mov", "m4v":
		args = append(args, "-movflags", "+faststart")
	}

	return append(args, "-y", out)
}

// Stream is one ffprobe stream.
type Stream struct {
	Type       string `json:"codec_type"`
	Codec      string `json:"codec_name"`
	Width      int    `json:"width"`
	Height     int    `json:"height"`
	SampleRate string `json:"sample_rate"`
	Channels   int    `json:"channels"`
	FrameRate  string `json:"r_frame_rate"`
}

// Info is the subset of ffprobe output uts cares about.
type Info struct {
	Streams []Stream `json:"streams"`
	Format  struct {
		Duration string `json:"duration"`
		BitRate  string `json:"bit_rate"`
		Name     string `json:"format_name"`
	} `json:"format"`
}

// Video returns the first video stream, or nil.
func (i Info) Video() *Stream {
	for idx := range i.Streams {
		if i.Streams[idx].Type == "video" {
			return &i.Streams[idx]
		}
	}

	return nil
}

// Audio returns the first audio stream, or nil.
func (i Info) Audio() *Stream {
	for idx := range i.Streams {
		if i.Streams[idx].Type == "audio" {
			return &i.Streams[idx]
		}
	}

	return nil
}

// Duration returns the container duration in seconds, or 0 when unknown.
func (i Info) Duration() float64 {
	dur, err := strconv.ParseFloat(i.Format.Duration, 64)
	if err != nil {
		return 0
	}

	return dur
}

// BitRate returns the overall bit rate in bits per second, or 0 when unknown.
func (i Info) BitRate() int64 {
	rate, err := strconv.ParseInt(i.Format.BitRate, 10, 64)
	if err != nil {
		return 0
	}

	return rate
}

// Probe reads stream and container metadata with ffprobe.
func Probe(file string) (*Info, error) {
	if !util.HasTool("ffprobe") {
		return nil, ErrNoProbe
	}

	out, err := exec.CommandContext(context.Background(), "ffprobe",
		"-v", "error", "-show_streams", "-show_format", "-of", "json", file).Output()
	if err != nil {
		return nil, err
	}

	var info Info

	err = json.Unmarshal(out, &info)
	if err != nil {
		return nil, err
	}

	return &info, nil
}

// CanRemux reports whether every video and audio stream of a probed file can
// be copied as-is into the target container.
func CanRemux(info *Info, container string) bool {
	seen := false

	for _, stream := range info.Streams {
		if stream.Type != "video" && stream.Type != "audio" {
			continue
		}

		// Embedded cover art is a still image, not a video track.
		if stream.Type == "video" && (stream.Codec == "mjpeg" || stream.Codec == "png") {
			return false
		}

		if !format.ContainerAccepts(container, stream.Codec) {
			return false
		}

		seen = true
	}

	return seen
}
