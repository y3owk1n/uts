package compress

import (
	"context"
	"fmt"
	"os/exec"
	"path/filepath"
	"strconv"

	derrors "github.com/y3owk1n/uts/internal/core/errors"
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

	if !util.HasTool("ffmpeg") {
		return derrors.New(
			derrors.CodeToolNotFound,
			"ffmpeg not found — install: brew install ffmpeg",
		)
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

		out := util.CalcOutputPath(file, "small", opts.OutputDir)
		origSize := util.FileSize(file)

		ui.Message.Stepf("[%d/%d] %s (%s)", idx+1, total, file, util.HumanSize(origSize))

		if opts.DryRun {
			ui.Message.Infof(
				"[dry-run] Would compress %s -> %s (crf=%d, preset=%s)%s",
				file,
				out,
				crf,
				preset,
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

		ext := filepath.Ext(file)
		vcodec, acodec := util.VideoCodecs(ext)

		cmd := exec.CommandContext(
			context.Background(), "ffmpeg",
			"-i", file,
			"-vcodec", vcodec,
			"-crf", strconv.Itoa(crf),
			"-preset", preset,
			"-acodec", acodec,
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

			if opts.InPlace {
				util.MaybeReplaceOrRemove(out, file)
			}
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
