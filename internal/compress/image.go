//nolint:mnd,goconst
package compress

import (
	"fmt"
	"strconv"

	derrors "github.com/y3owk1n/uts/internal/core/errors"
	"github.com/y3owk1n/uts/internal/format"
	"github.com/y3owk1n/uts/internal/job"
	"github.com/y3owk1n/uts/internal/ui"
	"github.com/y3owk1n/uts/internal/util"
)

var errImageMagick = derrors.New(
	derrors.CodeToolNotFound,
	"ImageMagick not found — install: brew install imagemagick",
)

// ImageOptions represents options for image compression.
type ImageOptions struct {
	Files     []string
	Quality   string
	OutputDir string
	InPlace   bool
	DryRun    bool
	// MaxEdge caps the longest edge in pixels; 0 keeps the dimensions.
	MaxEdge int
}

// maxNote describes a --max setting for the intro line.
func maxNote(maxEdge int) string {
	if maxEdge > 0 {
		return fmt.Sprintf(", longest edge ≤ %dpx", maxEdge)
	}

	return ""
}

// Image compresses image files using the best tool available per format.
func Image(opts ImageOptions) error {
	quality, err := util.ImageQuality(opts.Quality)
	if err != nil {
		return err
	}

	ui.Message.Infof(
		"Image compression at %s quality (value=%d)%s",
		opts.Quality,
		quality,
		maxNote(opts.MaxEdge),
	)

	return job.Run(opts.Files, job.Options{
		Verb:    "Compressing",
		Done:    "Compressed",
		Noun:    "image file",
		Compare: true,
		InPlace: opts.InPlace,
		DryRun:  opts.DryRun,
		Code:    derrors.CodeCompressionFailed,
	}, func(file string) (*job.Job, error) {
		return planImage(file, quality, opts.MaxEdge, opts.OutputDir)
	})
}

func planImage(file string, quality, maxEdge int, outputDir string) (*job.Job, error) {
	ext := format.Ext(file)
	out := util.CalcOutputPath(file, "small", outputDir)

	var (
		tool  string
		steps []job.Step
		err   error
	)

	switch ext {
	case "png":
		tool, steps, err = pngSteps(file, out, quality, maxEdge)
	case "jpg", "jpeg":
		tool, steps, err = jpegSteps(file, out, quality, maxEdge)
	case "webp":
		tool, steps, err = webpSteps(file, out, quality, maxEdge)
	case "gif":
		tool, steps, err = gifSteps(file, out, maxEdge)
	case "heic", "heif":
		// HEIC has no lossy re-encoder in common toolchains; JPEG is the
		// portable result users expect from "compress this photo".
		out = util.CalcOutputPathExt(file, "small", "jpg", outputDir)
		tool, steps, err = heicSteps(file, out, quality, maxEdge)
	case "bmp", "tiff", "tif", "avif", "avifs":
		tool, steps, err = magickSteps(file, out, quality, maxEdge)
	default:
		return nil, derrors.Newf(derrors.CodeUnsupportedFormat, "unsupported image format .%s", ext)
	}

	if err != nil {
		return nil, err
	}

	return &job.Job{Input: file, Output: out, Steps: steps, Note: "via " + tool}, nil
}

// pngSteps compresses a PNG. With a size cap ImageMagick writes the resized
// file first and pngquant then quantizes it in place.
func pngSteps(file, out string, quality, maxEdge int) (string, []job.Step, error) {
	switch {
	case util.HasTool("pngquant"):
		quant := fmt.Sprintf("--quality=%d-%d", max(quality-10, 0), quality)

		var (
			steps []job.Step
			tool  string
		)

		if maxEdge > 0 {
			resize, err := resizeStep(file, out, maxEdge)
			if err != nil {
				return "", nil, err
			}

			steps = append(
				steps,
				resize,
				job.Exec(
					"pngquant",
					quant,
					"--speed",
					"1",
					"--strip",
					"--force",
					"--ext",
					".png",
					out,
				),
			)
			tool = "ImageMagick + pngquant"
		} else {
			steps = append(
				steps,
				job.Exec(
					"pngquant",
					quant,
					"--speed",
					"1",
					"--strip",
					"--force",
					"--output",
					out,
					"--",
					file,
				),
			)
			tool = "pngquant"
		}

		if util.HasTool("optipng") {
			steps = append(steps, job.Exec("optipng", "-quiet", "-o2", out))
			tool += " + optipng"
		}

		return tool, steps, nil
	case util.HasTool("optipng"):
		first, err := stageStep(file, out, maxEdge)
		if err != nil {
			return "", nil, err
		}

		return "optipng", []job.Step{first, job.Exec("optipng", "-quiet", "-o2", out)}, nil
	}

	return magickSteps(file, out, quality, maxEdge)
}

func jpegSteps(file, out string, quality, maxEdge int) (string, []job.Step, error) {
	if util.HasTool("jpegoptim") {
		first, err := stageStep(file, out, maxEdge)
		if err != nil {
			return "", nil, err
		}

		return "jpegoptim", []job.Step{
			first,
			job.Exec("jpegoptim", fmt.Sprintf("--max=%d", quality), "--strip-all", "--quiet", out),
		}, nil
	}

	return magickSteps(file, out, quality, maxEdge)
}

func webpSteps(file, out string, quality, maxEdge int) (string, []job.Step, error) {
	if util.HasTool("cwebp") && maxEdge == 0 {
		return "cwebp", []job.Step{
			job.Exec("cwebp", "-quiet", "-q", strconv.Itoa(quality), "-m", "6", file, "-o", out),
		}, nil
	}

	return magickSteps(file, out, quality, maxEdge)
}

func gifSteps(file, out string, maxEdge int) (string, []job.Step, error) {
	if util.HasTool("gifsicle") {
		args := []string{"-O3", "--lossy=80"}
		if maxEdge > 0 {
			args = append(args, "--resize-fit", fmt.Sprintf("%dx%d", maxEdge, maxEdge))
		}

		return "gifsicle", []job.Step{job.Exec("gifsicle", append(args, file, "-o", out)...)}, nil
	}

	return magickSteps(file, out, 80, maxEdge)
}

func heicSteps(file, out string, quality, maxEdge int) (string, []job.Step, error) {
	if util.HasTool("heif-convert") && maxEdge == 0 {
		return "heif-convert", []job.Step{
			job.Exec("heif-convert", "-q", strconv.Itoa(quality), file, out),
		}, nil
	}

	return magickSteps(file, out, quality, maxEdge)
}

func magickSteps(file, out string, quality, maxEdge int) (string, []job.Step, error) {
	bin := util.MagickBin()
	if bin == "" {
		return "", nil, errImageMagick
	}

	args := make([]string, 0, 7)
	args = append(args, file)
	args = append(args, resizeArgs(maxEdge)...)
	args = append(args, "-quality", strconv.Itoa(quality), "-strip", out)

	return "ImageMagick", []job.Step{job.Exec(bin, args...)}, nil
}

// resizeArgs returns ImageMagick arguments that shrink the longest edge to
// maxEdge, never enlarging. Empty when no cap is set.
func resizeArgs(maxEdge int) []string {
	if maxEdge <= 0 {
		return nil
	}

	return []string{"-resize", fmt.Sprintf("%dx%d>", maxEdge, maxEdge)}
}

// resizeStep writes a resized copy of file to out with ImageMagick.
func resizeStep(file, out string, maxEdge int) (job.Step, error) {
	bin := util.MagickBin()
	if bin == "" {
		return job.Step{}, derrors.New(derrors.CodeToolNotFound,
			"--max needs ImageMagick for this format — install: brew install imagemagick")
	}

	args := make([]string, 0, 4)
	args = append(args, file)
	args = append(args, resizeArgs(maxEdge)...)
	args = append(args, out)

	return job.Exec(bin, args...), nil
}

// stageStep places the working copy at out: a plain copy, or a resized copy
// when a size cap is set. In-place tools then operate on out.
func stageStep(file, out string, maxEdge int) (job.Step, error) {
	if maxEdge > 0 {
		return resizeStep(file, out, maxEdge)
	}

	return copyStep(file, out), nil
}

func copyStep(src, dst string) job.Step {
	return job.Do(fmt.Sprintf("copy %s -> %s", src, dst), func() error {
		return util.CopyFile(src, dst)
	})
}
