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
}

// Image compresses image files using the best tool available per format.
func Image(opts ImageOptions) error {
	quality, err := util.ImageQuality(opts.Quality)
	if err != nil {
		return err
	}

	ui.Message.Infof("Image compression at %s quality (value=%d)", opts.Quality, quality)

	return job.Run(opts.Files, job.Options{
		Verb:    "Compressing",
		Done:    "Compressed",
		Noun:    "image file",
		Compare: true,
		InPlace: opts.InPlace,
		DryRun:  opts.DryRun,
		Code:    derrors.CodeCompressionFailed,
	}, func(file string) (*job.Job, error) {
		return planImage(file, quality, opts.OutputDir)
	})
}

func planImage(file string, quality int, outputDir string) (*job.Job, error) {
	ext := format.Ext(file)
	out := util.CalcOutputPath(file, "small", outputDir)

	var (
		tool  string
		steps []job.Step
		err   error
	)

	switch ext {
	case "png":
		tool, steps, err = pngSteps(file, out, quality)
	case "jpg", "jpeg":
		tool, steps, err = jpegSteps(file, out, quality)
	case "webp":
		tool, steps, err = webpSteps(file, out, quality)
	case "gif":
		tool, steps, err = gifSteps(file, out)
	case "heic", "heif":
		// HEIC has no lossy re-encoder in common toolchains; JPEG is the
		// portable result users expect from "compress this photo".
		out = util.CalcOutputPathExt(file, "small", "jpg", outputDir)
		tool, steps, err = heicSteps(file, out, quality)
	case "bmp", "tiff", "tif", "avif", "avifs":
		tool, steps, err = magickSteps(file, out, quality)
	default:
		return nil, derrors.Newf(derrors.CodeUnsupportedFormat, "unsupported image format .%s", ext)
	}

	if err != nil {
		return nil, err
	}

	return &job.Job{Input: file, Output: out, Steps: steps, Note: "via " + tool}, nil
}

func pngSteps(file, out string, quality int) (string, []job.Step, error) {
	switch {
	case util.HasTool("pngquant"):
		steps := []job.Step{job.Exec("pngquant",
			fmt.Sprintf("--quality=%d-%d", max(quality-10, 0), quality),
			"--speed", "1", "--strip", "--force", "--output", out, "--", file)}
		tool := "pngquant"

		if util.HasTool("optipng") {
			steps = append(steps, job.Exec("optipng", "-quiet", "-o2", out))
			tool += " + optipng"
		}

		return tool, steps, nil
	case util.HasTool("optipng"):
		return "optipng", []job.Step{
			copyStep(file, out),
			job.Exec("optipng", "-quiet", "-o2", out),
		}, nil
	}

	return magickSteps(file, out, quality)
}

func jpegSteps(file, out string, quality int) (string, []job.Step, error) {
	if util.HasTool("jpegoptim") {
		return "jpegoptim", []job.Step{
			copyStep(file, out),
			job.Exec("jpegoptim", fmt.Sprintf("--max=%d", quality), "--strip-all", "--quiet", out),
		}, nil
	}

	return magickSteps(file, out, quality)
}

func webpSteps(file, out string, quality int) (string, []job.Step, error) {
	if util.HasTool("cwebp") {
		return "cwebp", []job.Step{
			job.Exec("cwebp", "-quiet", "-q", strconv.Itoa(quality), "-m", "6", file, "-o", out),
		}, nil
	}

	return magickSteps(file, out, quality)
}

func gifSteps(file, out string) (string, []job.Step, error) {
	if util.HasTool("gifsicle") {
		return "gifsicle", []job.Step{
			job.Exec("gifsicle", "-O3", "--lossy=80", file, "-o", out),
		}, nil
	}

	return magickSteps(file, out, 80)
}

func heicSteps(file, out string, quality int) (string, []job.Step, error) {
	if util.HasTool("heif-convert") {
		return "heif-convert", []job.Step{
			job.Exec("heif-convert", "-q", strconv.Itoa(quality), file, out),
		}, nil
	}

	return magickSteps(file, out, quality)
}

func magickSteps(file, out string, quality int) (string, []job.Step, error) {
	bin := util.MagickBin()
	if bin == "" {
		return "", nil, errImageMagick
	}

	return "ImageMagick", []job.Step{
		job.Exec(bin, file, "-quality", strconv.Itoa(quality), "-strip", out),
	}, nil
}

func copyStep(src, dst string) job.Step {
	return job.Do(fmt.Sprintf("copy %s -> %s", src, dst), func() error {
		return util.CopyFile(src, dst)
	})
}
