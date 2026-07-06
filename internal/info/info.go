package info

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/y3owk1n/uts/internal/ui"
	"github.com/y3owk1n/uts/internal/util"
)

type Options struct {
	Files []string
}

func Show(opts Options) {
	for _, file := range opts.Files {
		info, err := os.Stat(file)
		if err != nil {
			ui.Message.Warnf("Cannot access: %s", file)
			continue
		}
		if info.IsDir() {
			ui.Message.Warnf("Not a file: %s", file)
			continue
		}

		ext := strings.ToLower(strings.TrimPrefix(filepath.Ext(file), "."))
		size := util.FileSize(file)

		fmt.Println()

		primary := ui.Message.Highlight
		dim := ui.Message.Dim

		fmt.Printf("  %s\n", primary(filepath.Base(file)))
		ui.Message.Pair("Size", util.HumanSize(size))
		ui.Message.Pair("Type", "."+ext)

		switch {
		case isVideo(ext):
			ui.Message.Pair("Category", ui.Message.Accent("video"))
			ui.Message.Mutedf("  Tool:     ffmpeg (libx265)")
			ui.Message.Mutedf("  Compress: %s", dim(fmt.Sprintf("uts video compress \"%s\" [-q low|medium|high|<0-51>]", file)))
			ui.Message.Mutedf("  Convert:  %s", dim(fmt.Sprintf("uts video convert \"%s\" --to mkv", file)))
		case isImage(ext):
			ui.Message.Pair("Category", ui.Message.Accent("image"))
			tool := "ImageMagick"
			if ext == "png" {
				tool = "pngquant + optipng"
			} else if ext == "jpg" || ext == "jpeg" {
				tool = "jpegoptim"
			} else if ext == "webp" {
				tool = "cwebp"
			} else if ext == "gif" {
				tool = "gifsicle"
			}
			ui.Message.Mutedf("  Tool:     %s", tool)
			ui.Message.Mutedf("  Compress: %s", dim(fmt.Sprintf("uts image compress \"%s\" [-q low|medium|high|<1-100>]", file)))
			ui.Message.Mutedf("  Convert:  %s", dim(fmt.Sprintf("uts image convert \"%s\" --to webp", file)))
		case isAudio(ext):
			ui.Message.Pair("Category", ui.Message.Accent("audio"))
			ui.Message.Mutedf("  Tool:     ffmpeg (aac)")
			ui.Message.Mutedf("  Compress: %s", dim(fmt.Sprintf("uts audio compress \"%s\" [-q low|medium|high|<kbps>]", file)))
			ui.Message.Mutedf("  Convert:  %s", dim(fmt.Sprintf("uts audio convert \"%s\" --to mp3", file)))
		case ext == "pdf":
			ui.Message.Pair("Category", ui.Message.Accent("pdf"))
			ui.Message.Mutedf("  Tool:     ghostscript")
			ui.Message.Mutedf("  Compress: %s", dim(fmt.Sprintf("uts pdf compress \"%s\" [-q low|medium|high|<dpi>]", file)))
			ui.Message.Mutedf("  Convert:  %s", dim(fmt.Sprintf("uts pdf convert \"%s\" --to jpg", file)))
		case isArchive(ext):
			ui.Message.Pair("Category", ui.Message.Accent("archive"))
			switch ext {
			case "zip":
				ui.Message.Mutedf("  Tool:     unzip")
				ui.Message.Mutedf("  Extract:  %s", dim(fmt.Sprintf("uts archive extract \"%s\"", file)))
				ui.Message.Mutedf("  List:     %s", dim(fmt.Sprintf("uts archive list \"%s\"", file)))
			default:
				ui.Message.Mutedf("  Tool:     tar")
				ui.Message.Mutedf("  Extract:  %s", dim(fmt.Sprintf("uts archive extract \"%s\"", file)))
				ui.Message.Mutedf("  List:     %s", dim(fmt.Sprintf("uts archive list \"%s\"", file)))
			}
		default:
			ui.Message.Pair("Category", ui.Message.Warn("unknown"))
			ui.Message.Mutedf("  No uts strategy for .%s", ext)
		}
	}
	fmt.Println()
}

func isVideo(ext string) bool {
	switch ext {
	case "mp4", "mov", "mkv", "avi", "webm", "m4v", "flv", "wmv":
		return true
	}
	return false
}

func isImage(ext string) bool {
	switch ext {
	case "png", "jpg", "jpeg", "webp", "gif", "bmp", "tiff", "tif", "heic", "heif", "avif", "avifs":
		return true
	}
	return false
}

func isAudio(ext string) bool {
	switch ext {
	case "wav", "flac", "aac", "mp3", "m4a", "opus", "ogg", "wma":
		return true
	}
	return false
}

func isArchive(ext string) bool {
	switch ext {
	case "zip", "tar", "gz", "tgz", "zst", "zstd", "xz", "txz", "bz2", "tbz2", "7z":
		return true
	}
	return false
}
