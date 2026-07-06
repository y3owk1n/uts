//nolint:mnd
package compress

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	derrors "github.com/y3owk1n/uts/internal/core/errors"
	"github.com/y3owk1n/uts/internal/ui"
	"github.com/y3owk1n/uts/internal/util"
)

var errImageMagick = derrors.New(derrors.CodeToolNotFound, "imagemagick not found")

// ImageOptions represents options for image compression.
type ImageOptions struct {
	Files     []string
	Quality   string
	OutputDir string
	InPlace   bool
	DryRun    bool
}

// Image compresses image files using various tools.
func Image(opts ImageOptions) error {
	qualityVal, err := util.PresetVal(opts.Quality, 60, 80, 90)
	if err != nil {
		return err
	}

	ui.Message.Infof("Image compression at %s quality (value=%d)", opts.Quality, qualityVal)
	total := len(opts.Files)

	for idx, file := range opts.Files {
		if !util.FileExists(file) {
			ui.Message.Warnf("File not found: %s", file)

			continue
		}

		ext := strings.ToLower(strings.TrimPrefix(filepath.Ext(file), "."))
		out := util.CalcOutputPath(file, "small", opts.OutputDir)
		origSize := util.FileSize(file)

		ui.Message.Stepf("[%d/%d] %s (%s)", idx+1, total, file, util.HumanSize(origSize))

		if opts.DryRun {
			ui.Message.Infof(
				"[dry-run] Would compress %s -> %s (format=%s, quality=%d)%s",
				file,
				out,
				ext,
				qualityVal,
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

		var compressOK bool
		switch ext {
		case "png":
			compressOK = compressPNG(file, out, qualityVal) == nil
		case "jpg", "jpeg":
			compressOK = compressJPEG(file, out, qualityVal) == nil
		case "webp":
			compressOK = compressWebP(file, out, qualityVal) == nil
		case "gif":
			compressOK = compressGIF(file, out) == nil
		case "bmp", "tiff", "tif":
			compressOK = compressGeneric(file, out, qualityVal) == nil
		case "heic", "heif":
			out = strings.TrimSuffix(out, filepath.Ext(out)) + ".jpg"
			compressOK = compressHEIC(file, out, qualityVal) == nil
		case "avif", "avifs":
			out = strings.TrimSuffix(out, filepath.Ext(out)) + ".avif"
			compressOK = compressAVIF(file, out, qualityVal) == nil
		default:
			ui.Message.Warnf("Unsupported image format: %s — skipping %s", ext, file)
			spinner.Stop()

			continue
		}

		spinner.Stop()

		if compressOK && util.FileExists(out) {
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
		}
	}

	if total > 1 {
		ui.Message.Successf("Compressed %d image files", total)
	}

	return nil
}

func compressPNG(file, out string, quality int) error {
	if util.HasTool("pngquant") {
		ui.Message.Mutedf("Using pngquant")

		cmd := exec.CommandContext(context.Background(), "pngquant",
			fmt.Sprintf("--quality=%d-%d", quality-10, quality),
			"--speed", "1", "--strip", "--output", out, "--", file)

		err := cmd.Run()
		if err != nil {
			return err
		}

		if util.HasTool("optipng") && util.FileExists(out) {
			ui.Message.Mutedf("Optimizing with optipng")

			_ = exec.CommandContext(context.Background(), "optipng", "-quiet", "-o2", out).Run()
		}

		return nil
	}

	if util.HasTool("optipng") {
		ui.Message.Mutedf("Using optipng")

		err := copyFile(file, out)
		if err != nil {
			return err
		}

		return exec.CommandContext(context.Background(), "optipng", "-quiet", "-o2", out).Run()
	}

	return magickCmd(file, out, quality)
}

func compressJPEG(file, out string, quality int) error {
	if util.HasTool("jpegoptim") {
		ui.Message.Mutedf("Using jpegoptim")

		err := copyFile(file, out)
		if err != nil {
			return err
		}

		return exec.CommandContext(context.Background(), "jpegoptim", fmt.Sprintf("--max=%d", quality), "--strip-all", "--quiet", out).
			Run()
	}

	return magickCmd(file, out, quality)
}

func compressWebP(file, out string, quality int) error {
	if util.HasTool("cwebp") {
		ui.Message.Mutedf("Using cwebp")

		return exec.CommandContext(context.Background(), "cwebp", "-q", strconv.Itoa(quality), "-m", "6", file, "-o", out).
			Run()
	}

	return magickCmd(file, out, quality)
}

func compressGIF(file, out string) error {
	if util.HasTool("gifsicle") {
		ui.Message.Mutedf("Using gifsicle")

		return exec.CommandContext(context.Background(), "gifsicle", "-O3", "--lossy=80", file, "-o", out).
			Run()
	}

	return magickCmd(file, out, 80)
}

func compressGeneric(file, out string, quality int) error {
	return magickCmd(file, out, quality)
}

func compressHEIC(file, out string, quality int) error {
	if util.HasTool("heif-convert") {
		ui.Message.Mutedf("Using heif-convert")

		return exec.CommandContext(context.Background(), "heif-convert", "-q", strconv.Itoa(quality), file, out).
			Run()
	}

	return magickCmd(file, out, quality)
}

func compressAVIF(file, out string, quality int) error {
	if util.HasTool("cavif") {
		ui.Message.Mutedf("Using cavif")

		return exec.CommandContext(context.Background(), "cavif", "-q", strconv.Itoa(quality), "-s", "6", "-o", out, file).
			Run()
	}

	if util.HasTool("avifenc") {
		ui.Message.Mutedf("Using avifenc")

		quantizer := (100 - quality) * 63 / 100

		return exec.CommandContext(context.Background(), "avifenc", "--min", "0", "--max", strconv.Itoa(quantizer), "-s", "6", file, out).
			Run()
	}

	return magickCmd(file, out, quality)
}

func magickCmd(input, output string, quality int) error {
	if util.HasTool("magick") {
		ui.Message.Mutedf("Using ImageMagick")

		return exec.CommandContext(context.Background(), "magick", input, "-quality", strconv.Itoa(quality), "-strip", output).
			Run()
	}

	if util.HasTool("convert") {
		ui.Message.Mutedf("Using ImageMagick (convert)")

		return exec.CommandContext(context.Background(), "convert", input, "-quality", strconv.Itoa(quality), "-strip", output).
			Run()
	}

	ui.Message.Errorf("ImageMagick not found — install: brew install imagemagick")

	return errImageMagick
}

func copyFile(src, dst string) error {
	input, err := os.ReadFile(src)
	if err != nil {
		return err
	}

	return os.WriteFile(dst, input, 0o644)
}
