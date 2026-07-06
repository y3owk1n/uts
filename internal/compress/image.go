package compress

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/y3owk1n/uts/internal/ui"
	"github.com/y3owk1n/uts/internal/util"
)

type ImageOptions struct {
	Files     []string
	Quality   string
	OutputDir string
	InPlace   bool
	DryRun    bool
}

func Image(opts ImageOptions) error {
	qualityVal, err := util.PresetVal(opts.Quality, 60, 80, 90)
	if err != nil {
		return err
	}

	ui.Message.Infof("Image compression at %s quality (value=%d)", opts.Quality, qualityVal)
	total := len(opts.Files)

	for i, file := range opts.Files {
		if !util.FileExists(file) {
			ui.Message.Warnf("File not found: %s", file)
			continue
		}

		ext := strings.ToLower(strings.TrimPrefix(filepath.Ext(file), "."))
		out := util.OutputPath(file, "small")
		origSize := util.FileSize(file)

		ui.Message.Stepf("[%d/%d] %s (%s)", i+1, total, file, util.HumanSize(origSize))

		if opts.DryRun {
			ui.Message.Infof("[dry-run] Would compress %s -> %s (format=%s, quality=%d)", file, out, ext, qualityVal)
			continue
		}

		util.EnsureDir(out)
		sp := ui.NewSpinner(nil, 0)
		sp.SetSuffix(fmt.Sprintf("Compressing %s...", file))
		sp.Start()

		compressOK := true
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
			sp.Stop()
			continue
		}
		sp.Stop()

		if compressOK && util.FileExists(out) {
			newSize := util.FileSize(out)
			ratio := util.CompressionRatio(origSize, newSize)
			ui.Message.Successf("%s: %s → %s %s", file, util.HumanSize(origSize), util.HumanSize(newSize), ratio)
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
	if hasTool("pngquant") {
		ui.Message.Mutedf("Using pngquant")
		cmd := exec.Command("pngquant",
			fmt.Sprintf("--quality=%d-%d", quality-10, quality),
			"--speed", "1", "--strip", "--output", out, "--", file)
		if err := cmd.Run(); err != nil {
			return err
		}
		if hasTool("optipng") && util.FileExists(out) {
			ui.Message.Mutedf("Optimizing with optipng")
			exec.Command("optipng", "-quiet", "-o2", out).Run()
		}
		return nil
	}
	if hasTool("optipng") {
		ui.Message.Mutedf("Using optipng")
		if err := copyFile(file, out); err != nil {
			return err
		}
		return exec.Command("optipng", "-quiet", "-o2", out).Run()
	}
	return magickCmd(file, out, quality)
}

func compressJPEG(file, out string, quality int) error {
	if hasTool("jpegoptim") {
		ui.Message.Mutedf("Using jpegoptim")
		if err := copyFile(file, out); err != nil {
			return err
		}
		return exec.Command("jpegoptim", fmt.Sprintf("--max=%d", quality), "--strip-all", "--quiet", out).Run()
	}
	return magickCmd(file, out, quality)
}

func compressWebP(file, out string, quality int) error {
	if hasTool("cwebp") {
		ui.Message.Mutedf("Using cwebp")
		return exec.Command("cwebp", "-q", fmt.Sprintf("%d", quality), "-m", "6", file, "-o", out).Run()
	}
	return magickCmd(file, out, quality)
}

func compressGIF(file, out string) error {
	if hasTool("gifsicle") {
		ui.Message.Mutedf("Using gifsicle")
		return exec.Command("gifsicle", "-O3", "--lossy=80", file, "-o", out).Run()
	}
	return magickCmd(file, out, 80)
}

func compressGeneric(file, out string, quality int) error {
	return magickCmd(file, out, quality)
}

func compressHEIC(file, out string, quality int) error {
	if hasTool("heif-convert") {
		ui.Message.Mutedf("Using heif-convert")
		return exec.Command("heif-convert", "-q", fmt.Sprintf("%d", quality), file, out).Run()
	}
	return magickCmd(file, out, quality)
}

func compressAVIF(file, out string, quality int) error {
	if hasTool("cavif") {
		ui.Message.Mutedf("Using cavif")
		return exec.Command("cavif", "-q", fmt.Sprintf("%d", quality), "-s", "6", "-o", out, file).Run()
	}
	if hasTool("avifenc") {
		ui.Message.Mutedf("Using avifenc")
		quantizer := (100 - quality) * 63 / 100
		return exec.Command("avifenc", "--min", "0", "--max", fmt.Sprintf("%d", quantizer), "-s", "6", file, out).Run()
	}
	return magickCmd(file, out, quality)
}

func magickCmd(input, output string, quality int) error {
	if hasTool("magick") {
		ui.Message.Mutedf("Using ImageMagick")
		return exec.Command("magick", input, "-quality", fmt.Sprintf("%d", quality), "-strip", output).Run()
	}
	if hasTool("convert") {
		ui.Message.Mutedf("Using ImageMagick (convert)")
		return exec.Command("convert", input, "-quality", fmt.Sprintf("%d", quality), "-strip", output).Run()
	}
	ui.Message.Errorf("No ImageMagick found. Install: brew install imagemagick")
	return fmt.Errorf("imagemagick not found")
}

func copyFile(src, dst string) error {
	input, err := os.ReadFile(src)
	if err != nil {
		return err
	}
	return os.WriteFile(dst, input, 0644)
}

func hasTool(name string) bool {
	_, err := exec.LookPath(name)
	return err == nil
}
