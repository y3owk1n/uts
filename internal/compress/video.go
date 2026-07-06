package compress

import (
	"context"
	"fmt"
	"os/exec"
	"strconv"

	"github.com/y3owk1n/uts/internal/ui"
	"github.com/y3owk1n/uts/internal/util"
)

// VideoOptions represents options for video compression.
type VideoOptions struct {
	Files     []string
	Quality   string
	OutputDir string
	InPlace   bool
	DryRun    bool
}

// Video compresses video files using ffmpeg.
func Video(opts VideoOptions) error {
	crf, preset, err := util.VideoQuality(opts.Quality)
	if err != nil {
		return err
	}

	ui.Message.Infof(
		"Video compression at %s quality (crf=%d, preset=%s)",
		opts.Quality,
		crf,
		preset,
	)

	total := len(opts.Files)
	for idx, file := range opts.Files {
		if !util.FileExists(file) {
			ui.Message.Warnf("File not found: %s", file)

			continue
		}

		out := util.OutputPath(file, "small")
		origSize := util.FileSize(file)

		ui.Message.Stepf("[%d/%d] %s (%s)", idx+1, total, file, util.HumanSize(origSize))

		if opts.DryRun {
			ui.Message.Infof(
				"[dry-run] Would compress %s -> %s (crf=%d, preset=%s)",
				file,
				out,
				crf,
				preset,
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
			"-vcodec", "libx265",
			"-crf", strconv.Itoa(crf),
			"-preset", preset,
			"-acodec", "aac",
			"-b:a", "128k",
			"-movflags", "+faststart",
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
		} else {
			ui.Message.Errorf("Compression failed: %s", file)
			ui.Message.Mutedf("ffmpeg: %s", string(output))
		}
	}

	if total > 1 {
		ui.Message.Successf("Compressed %d video files", total)
	}

	return nil
}
