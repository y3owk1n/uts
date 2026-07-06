package convert

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/y3owk1n/uts/internal/ui"
	"github.com/y3owk1n/uts/internal/util"
)

type ImageOptions struct {
	Files     []string
	Target    string
	Quality   string
	OutputDir string
	InPlace   bool
	DryRun    bool
}

func Image(opts ImageOptions) error {
	target := strings.ToLower(opts.Target)
	switch target {
	case "jpg", "jpeg", "png", "webp", "gif", "bmp", "tiff", "tif", "avif":
	default:
		ui.Message.Errorf("Unsupported target format: .%s", target)
		return nil
	}

	qualityVal, err := util.PresetVal(opts.Quality, 60, 80, 90)
	if err != nil {
		return err
	}

	if target == "jpeg" {
		target = "jpg"
	}

	ui.Message.Infof("Converting images to .%s (quality=%d)", target, qualityVal)

	total := len(opts.Files)
	for i, file := range opts.Files {
		if !util.FileExists(file) {
			ui.Message.Warnf("File not found: %s", file)
			continue
		}

		ext := strings.ToLower(strings.TrimPrefix(filepath.Ext(file), "."))
		if ext == target || (ext == "jpeg" && target == "jpg") {
			ui.Message.Warnf("Already .%s, skipping: %s", target, file)
			continue
		}

		out := util.ConvertOutputPath(file, target)
		origSize := util.FileSize(file)

		ui.Message.Stepf("[%d/%d] .%s → .%s (%s)", i+1, total, ext, target, util.HumanSize(origSize))

		if opts.DryRun {
			ui.Message.Infof("[dry-run] Would convert %s -> %s", file, out)
			continue
		}

		util.EnsureDir(out)
		sp := ui.NewSpinner(nil, 0)
		sp.SetSuffix(fmt.Sprintf("Converting %s...", file))
		sp.Start()

		var convertErr error
		if hasMagick() {
			if hasTool("magick") {
				convertErr = exec.Command("magick", file, "-quality", fmt.Sprintf("%d", qualityVal), "-strip", out).Run()
			} else {
				convertErr = exec.Command("convert", file, "-quality", fmt.Sprintf("%d", qualityVal), "-strip", out).Run()
			}
		} else if hasTool("sips") {
			sipsFmt := target
			if target == "jpg" {
				sipsFmt = "jpeg"
			}
			convertErr = exec.Command("sips", "-s", "format", sipsFmt, file, "--out", out).Run()
		} else {
			sp.Stop()
			ui.Message.Errorf("No conversion tool found. Install: brew install imagemagick")
			return nil
		}
		sp.Stop()

		if convertErr == nil && util.FileExists(out) {
			ui.Message.Successf("%s: %s → %s", file, util.HumanSize(origSize), util.HumanSize(util.FileSize(out)))
		} else {
			ui.Message.Errorf("Conversion failed: %s", file)
		}
	}

	if total > 1 {
		ui.Message.Successf("Converted %d image files", total)
	}
	return nil
}

func hasMagick() bool {
	return hasTool("magick") || hasTool("convert")
}

func hasTool(name string) bool {
	_, err := exec.LookPath(name)
	return err == nil
}
