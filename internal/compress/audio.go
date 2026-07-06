package compress

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/y3owk1n/uts/internal/ui"
	"github.com/y3owk1n/uts/internal/util"
)

type AudioOptions struct {
	Files     []string
	Quality   string
	OutputDir string
	InPlace   bool
	DryRun    bool
}

func Audio(opts AudioOptions) error {
	bitrate, err := util.AudioBitrate(opts.Quality)
	if err != nil {
		return err
	}

	ui.Message.Infof("Audio compression at %s quality (bitrate=%s, codec=aac)", opts.Quality, bitrate)
	total := len(opts.Files)

	for i, file := range opts.Files {
		if !util.FileExists(file) {
			ui.Message.Warnf("File not found: %s", file)
			continue
		}

		out := util.OutputPathExt(file, "small", "m4a")
		origSize := util.FileSize(file)

		ui.Message.Stepf("[%d/%d] %s (%s)", i+1, total, file, util.HumanSize(origSize))

		if opts.DryRun {
			ui.Message.Infof("[dry-run] Would compress %s -> %s (bitrate=%s)", file, out, bitrate)
			continue
		}

		util.EnsureDir(out)
		sp := ui.NewSpinner(nil, 0)
		sp.SetSuffix(fmt.Sprintf("Compressing %s...", file))
		sp.Start()

		output, err := exec.Command("ffmpeg",
			"-i", file,
			"-c:a", "aac",
			"-b:a", bitrate,
			"-y", out,
		).CombinedOutput()
		sp.Stop()

		if err == nil && util.FileExists(out) {
			newSize := util.FileSize(out)
			ratio := util.CompressionRatio(origSize, newSize)
			ui.Message.Successf("%s: %s → %s %s", file, util.HumanSize(origSize), util.HumanSize(newSize), ratio)
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

func isAudio(ext string) bool {
	switch strings.ToLower(ext) {
	case "wav", "flac", "aac", "mp3", "m4a", "opus", "ogg", "wma":
		return true
	}
	return false
}

func isVideo(ext string) bool {
	switch strings.ToLower(ext) {
	case "mp4", "mov", "mkv", "avi", "webm", "m4v", "flv", "wmv":
		return true
	}
	return false
}

func isImage(ext string) bool {
	switch strings.ToLower(filepath.Ext(ext)) {
	case ".png", ".jpg", ".jpeg", ".webp", ".gif", ".bmp", ".tiff", ".tif", ".heic", ".heif", ".avif", ".avifs":
		return true
	}
	return false
}
