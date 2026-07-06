//nolint:goconst
package convert

import (
	"context"
	"fmt"
	"os/exec"
	"strings"

	"github.com/y3owk1n/uts/internal/ui"
	"github.com/y3owk1n/uts/internal/util"
)

// VideoOptions represents options for video conversion.
type VideoOptions struct {
	Files     []string
	Target    string
	OutputDir string
	InPlace   bool
	DryRun    bool
}

// Video converts video files to the target format.
func Video(opts VideoOptions) error {
	target := strings.ToLower(opts.Target)
	switch target {
	case "mp4", "mkv", "webm", "mov", "avi", "flv":
	default:
		ui.Message.Errorf("Unsupported target format: .%s", target)

		return nil
	}

	vcodec, acodec := videoCodecs(target)
	ui.Message.Infof("Converting video to .%s (%s/%s)", target, vcodec, acodec)

	total := len(opts.Files)
	for idx, file := range opts.Files {
		if !util.FileExists(file) {
			ui.Message.Warnf("File not found: %s", file)

			continue
		}

		if strings.HasSuffix(strings.ToLower(file), "."+target) {
			ui.Message.Warnf("Already .%s, skipping: %s", target, file)

			continue
		}

		out := util.CalcConvertOutputPath(file, target, opts.OutputDir)
		origSize := util.FileSize(file)

		ui.Message.Stepf(
			"[%d/%d] %s → .%s (%s)",
			idx+1,
			total,
			file,
			target,
			util.HumanSize(origSize),
		)

		if opts.DryRun {
			ui.Message.Infof(
				"[dry-run] Would convert %s -> %s%s",
				file,
				out,
				util.InPlaceHint(opts.InPlace),
			)

			continue
		}

		_ = util.EnsureDir(out)

		spinner := ui.NewSpinner(nil, 0)
		spinner.SetSuffix(fmt.Sprintf("Converting %s...", file))
		spinner.Start()

		output, err := exec.CommandContext(
			context.Background(), "ffmpeg",
			"-i", file,
			"-vcodec", vcodec,
			"-acodec", acodec,
			"-y", out,
		).CombinedOutput()

		spinner.Stop()

		if err == nil && util.FileExists(out) {
			ui.Message.Successf(
				"%s: %s → %s",
				file,
				util.HumanSize(origSize),
				util.HumanSize(util.FileSize(out)),
			)

			if opts.InPlace {
				util.MaybeReplaceOrRemove(out, file)
			}
		} else {
			ui.Message.Errorf("Conversion failed: %s", file)
			ui.Message.Mutedf("ffmpeg: %s", string(output))
		}
	}

	if total > 1 {
		ui.Message.Successf("Converted %d video files", total)
	}

	return nil
}

func videoCodecs(target string) (string, string) {
	switch target {
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
