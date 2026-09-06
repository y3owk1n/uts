//nolint:goconst
package ffmpeg_test

import (
	"slices"
	"strings"
	"testing"

	"github.com/y3owk1n/uts/internal/ffmpeg"
)

func has(args []string, want ...string) bool {
	joined := " " + strings.Join(args, " ") + " "

	return strings.Contains(joined, " "+strings.Join(want, " ")+" ")
}

func TestEncodeArgsVP9(t *testing.T) {
	args := ffmpeg.EncodeArgs("in.mp4", "out.webm", "webm", 30, "fast")

	if !has(args, "-c:v", "libvpx-vp9") || !has(args, "-b:v", "0") || !has(args, "-cpu-used", "4") {
		t.Errorf("VP9 args missing constant-quality flags: %v", args)
	}

	if slices.Contains(args, "-preset") || slices.Contains(args, "-movflags") {
		t.Errorf("VP9 args must not carry x264-only flags: %v", args)
	}
}

func TestEncodeArgsMP4(t *testing.T) {
	args := ffmpeg.EncodeArgs("in.mov", "out.mp4", ".mp4", 23, "slow")

	for _, want := range [][]string{
		{"-c:v", "libx264"},
		{"-crf", "23"},
		{"-preset", "slow"},
		{"-pix_fmt", "yuv420p"},
		{"-movflags", "+faststart"},
		{"-c:a", "aac"},
	} {
		if !has(args, want...) {
			t.Errorf("mp4 args missing %v: %v", want, args)
		}
	}

	if args[len(args)-1] != "out.mp4" {
		t.Errorf("output must be last: %v", args)
	}
}

func TestEncodeArgsHEVCTag(t *testing.T) {
	if !has(ffmpeg.EncodeArgs("a", "b", "mkv", 28, "medium"), "-c:v", "libx265") {
		t.Error("mkv should use libx265")
	}

	if has(ffmpeg.EncodeArgs("a", "b", "mkv", 28, "medium"), "-tag:v", "hvc1") {
		t.Error("hvc1 tag only applies to mp4/mov")
	}
}

func TestRemuxArgs(t *testing.T) {
	args := ffmpeg.RemuxArgs("in.mov", "out.mp4", "mp4")
	if !has(args, "-c", "copy") || !has(args, "-movflags", "+faststart") {
		t.Errorf("remux args wrong: %v", args)
	}

	if has(ffmpeg.RemuxArgs("in.mp4", "out.mkv", "mkv"), "-movflags", "+faststart") {
		t.Error("faststart is mp4/mov only")
	}
}

func TestCanRemux(t *testing.T) {
	info := &ffmpeg.Info{Streams: []ffmpeg.Stream{
		{Type: "video", Codec: "h264"},
		{Type: "audio", Codec: "aac"},
		{Type: "data", Codec: "tmcd"},
	}}

	if !ffmpeg.CanRemux(info, "mp4") {
		t.Error("h264+aac should remux into mp4")
	}

	if ffmpeg.CanRemux(info, "webm") {
		t.Error("h264 must not remux into webm")
	}

	cover := &ffmpeg.Info{
		Streams: []ffmpeg.Stream{{Type: "video", Codec: "mjpeg"}, {Type: "audio", Codec: "aac"}},
	}
	if ffmpeg.CanRemux(cover, "mp4") {
		t.Error("cover art streams must force a re-encode")
	}

	if ffmpeg.CanRemux(&ffmpeg.Info{}, "mp4") {
		t.Error("no streams means nothing to remux")
	}
}
