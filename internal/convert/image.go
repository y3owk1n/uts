//nolint:mnd,goconst
package convert

import (
	"fmt"
	"slices"
	"strconv"
	"strings"

	derrors "github.com/y3owk1n/uts/internal/core/errors"
	"github.com/y3owk1n/uts/internal/format"
	"github.com/y3owk1n/uts/internal/job"
	"github.com/y3owk1n/uts/internal/ui"
	"github.com/y3owk1n/uts/internal/util"
)

// ImageOptions represents options for image conversion.
type ImageOptions struct {
	Files     []string
	Target    string
	Quality   string
	OutputDir string
	InPlace   bool
	DryRun    bool
	// MaxEdge caps the longest edge in pixels; 0 keeps the dimensions.
	MaxEdge int
}

// sipsTargets are the formats macOS sips can write, used when ImageMagick is
// missing.
var sipsTargets = []string{"jpg", "png", "tiff", "gif", "bmp"}

// Image converts image files to the target format.
func Image(opts ImageOptions) error {
	target := format.Normalize(opts.Target)
	if !slices.Contains(format.ImageTargets, target) {
		return unsupportedTarget(target, format.ImageTargets)
	}

	quality, err := util.ImageQuality(opts.Quality)
	if err != nil {
		return err
	}

	ui.Message.Infof(
		"Converting images to .%s (quality=%d)%s",
		target,
		quality,
		maxNote(opts.MaxEdge),
	)

	return job.Run(opts.Files, job.Options{
		Verb:    "Converting",
		Done:    "Converted",
		Noun:    "image file",
		InPlace: opts.InPlace,
		DryRun:  opts.DryRun,
		Code:    derrors.CodeConversionFailed,
	}, func(file string) (*job.Job, error) {
		ext := format.Ext(file)
		if format.Classify(ext) != format.Image {
			return nil, derrors.Newf(
				derrors.CodeUnsupportedFormat,
				"unsupported image format .%s",
				ext,
			)
		}

		if format.Same(ext, target) {
			return &job.Job{Input: file, Skip: "Already ." + target}, nil
		}

		out := util.CalcConvertOutputPath(file, target, opts.OutputDir)

		tool, step, err := imageStep(file, out, ext, target, quality, opts.MaxEdge)
		if err != nil {
			return nil, err
		}

		return &job.Job{
			Input:  file,
			Output: out,
			Steps:  []job.Step{step},
			Note:   fmt.Sprintf(".%s → .%s via %s", ext, target, tool),
		}, nil
	})
}

func imageStep(file, out, ext, target string, quality, maxEdge int) (string, job.Step, error) {
	// cavif and avifenc cannot resize, so a size cap goes through ImageMagick.
	if target == "avif" && maxEdge == 0 && (ext == "png" || ext == "jpg" || ext == "jpeg") {
		switch {
		case util.HasTool("cavif"):
			return "cavif", job.Exec(
				"cavif",
				"-q",
				strconv.Itoa(quality),
				"-s",
				"6",
				"-o",
				out,
				file,
			), nil
		case util.HasTool("avifenc"):
			quantizer := (100 - quality) * 63 / 100

			return "avifenc", job.Exec(
				"avifenc",
				"--min",
				"0",
				"--max",
				strconv.Itoa(quantizer),
				"-s",
				"6",
				file,
				out,
			), nil
		}
	}

	if bin := util.MagickBin(); bin != "" {
		args := []string{file}
		if maxEdge > 0 {
			args = append(args, "-resize", fmt.Sprintf("%dx%d>", maxEdge, maxEdge))
		}

		args = append(args, "-quality", strconv.Itoa(quality), "-strip", out)

		return "ImageMagick", job.Exec(bin, args...), nil
	}

	if util.HasTool("sips") && slices.Contains(sipsTargets, target) {
		sipsFmt := target
		if target == "jpg" {
			sipsFmt = "jpeg"
		}

		args := []string{}
		if maxEdge > 0 {
			args = append(args, "-Z", strconv.Itoa(maxEdge))
		}

		args = append(args, "-s", "format", sipsFmt, file, "--out", out)

		return "sips", job.Exec("sips", args...), nil
	}

	return "", job.Step{}, derrors.New(derrors.CodeToolNotFound,
		"ImageMagick not found — install: brew install imagemagick")
}

// maxNote describes a --max setting for the intro line.
func maxNote(maxEdge int) string {
	if maxEdge > 0 {
		return fmt.Sprintf(", longest edge ≤ %dpx", maxEdge)
	}

	return ""
}

func unsupportedTarget(target string, valid []string) error {
	return derrors.Newf(derrors.CodeUnsupportedFormat,
		"unsupported target format .%s (use one of: %s)", target, strings.Join(valid, ", "))
}
