package compress

import (
	"context"
	"fmt"
	"os/exec"

	"github.com/y3owk1n/uts/internal/ui"
	"github.com/y3owk1n/uts/internal/util"
)

// AudioOptions represents options for audio compression.
type AudioOptions struct {
	Files     []string
	Quality   string
	OutputDir string
	InPlace   bool
	DryRun    bool
}

// Audio compresses audio files using ffmpeg.
func Audio(opts AudioOptions) error {
	bitrate, err := util.AudioBitrate(opts.Quality)
	if err != nil {
		return err
	}

	ui.Message.Infof(
		"Audio compression at %s quality (bitrate=%s, codec=aac)",
		opts.Quality,
		bitrate,
	)
	total := len(opts.Files)

	for idx, file := range opts.Files {
		if !util.FileExists(file) {
			ui.Message.Warnf("File not found: %s", file)

			continue
		}

		out := util.OutputPathExt(file, "small", "m4a")
		origSize := util.FileSize(file)

		ui.Message.Stepf("[%d/%d] %s (%s)", idx+1, total, file, util.HumanSize(origSize))

		if opts.DryRun {
			ui.Message.Infof(
				"[dry-run] Would compress %s -> %s (bitrate=%s)%s",
				file,
				out,
				bitrate,
				util.InPlaceHint(opts.InPlace),
			)

			continue
		}

		err := util.EnsureDir(out)
		if err != nil {
			ui.Message.Errorf("Failed to create output directory: %v", err)

			continue
		}

		spinner := ui.NewSpinner(nil, 0)
		spinner.SetSuffix(fmt.Sprintf("Compressing %s...", file))
		spinner.Start()

		cmd := exec.CommandContext(
			context.Background(), "ffmpeg",
			"-i", file,
			"-c:a", "aac",
			"-b:a", bitrate,
			"-y", out,
		)
		output, err := cmd.CombinedOutput()

		spinner.Stop()

		if err == nil && util.FileExists(out) {
			newSize := util.FileSize(out)
			ratio := util.CompressionRatio(origSize, newSize)
			ui.Message.Successf(
				"%s: %s → %s %s",
				file,
				util.HumanSize(origSize),
				util.HumanSize(newSize),
				ratio,
			)

			if opts.InPlace {
				util.MaybeInPlace(out, file)
			}
		} else {
			ui.Message.Errorf("Compression failed: %s", file)
			ui.Message.Mutedf("ffmpeg: %s", string(output))
		}
	}

	if total > 1 {
		ui.Message.Successf("Compressed %d audio files", total)
	}

	return nil
}
