package convert

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/y3owk1n/uts/internal/ui"
	"github.com/y3owk1n/uts/internal/util"
)

type AudioOptions struct {
	Files     []string
	Target    string
	Quality   string
	OutputDir string
	DryRun    bool
}

func Audio(opts AudioOptions) error {
	target := strings.ToLower(opts.Target)
	codec, extHint := audioCodec(target)
	if codec == "" {
		ui.Message.Errorf("Unsupported target format: .%s", target)
		return nil
	}

	bitrate, err := util.AudioBitrate(opts.Quality)
	if err != nil {
		return err
	}

	ui.Message.Infof("Converting audio to .%s (%s, %s)", extHint, codec, bitrate)

	total := len(opts.Files)
	for i, file := range opts.Files {
		if !util.FileExists(file) {
			ui.Message.Warnf("File not found: %s", file)
			continue
		}

		if strings.HasSuffix(strings.ToLower(file), "."+extHint) {
			ui.Message.Warnf("Already .%s, skipping: %s", extHint, file)
			continue
		}

		out := util.ConvertOutputPath(file, extHint)
		origSize := util.FileSize(file)

		ui.Message.Stepf("[%d/%d] %s → .%s (%s)", i+1, total, file, extHint, util.HumanSize(origSize))

		if opts.DryRun {
			ui.Message.Infof("[dry-run] Would convert %s -> %s", file, out)
			continue
		}

		util.EnsureDir(out)
		sp := ui.NewSpinner(nil, 0)
		sp.SetSuffix(fmt.Sprintf("Converting %s...", file))
		sp.Start()

		output, err := exec.Command("ffmpeg",
			"-i", file,
			"-c:a", codec,
			"-b:a", bitrate,
			"-y", out,
		).CombinedOutput()
		sp.Stop()

		if err == nil && util.FileExists(out) {
			ui.Message.Successf("%s: %s → %s", file, util.HumanSize(origSize), util.HumanSize(util.FileSize(out)))
		} else {
			ui.Message.Errorf("Conversion failed: %s", file)
			ui.Message.Mutedf("ffmpeg: %s", string(output))
		}
	}

	if total > 1 {
		ui.Message.Successf("Converted %d audio files", total)
	}
	return nil
}

func audioCodec(target string) (string, string) {
	switch target {
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
