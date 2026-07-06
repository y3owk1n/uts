# uts

**One CLI to rule them all.** Compress, convert, and inspect any media file without remembering a dozen different command-line tools.

[![CI](https://github.com/y3owk1n/uts/actions/workflows/ci.yml/badge.svg)](https://github.com/y3owk1n/uts/actions/workflows/ci.yml)
[![License](https://img.shields.io/github/license/y3owk1n/uts.svg)](LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/y3owk1n/uts)](https://goreportcard.com/report/github.com/y3owk1n/uts)

---

Media files are everywhere, and every format has its own toolchain — `ffmpeg` for video, `pngquant` for PNGs, `jpegoptim` for JPEGs, `ghostscript` for PDFs, `tar` + a half-dozen compressors for archives. Each one has a *different* CLI syntax that you have to look up every time.

**uts** wraps all of them behind a single, predictable API:

```
uts <category> <action> <input...> [flags]
```

That's it. Same pattern for video, images, audio, PDFs, and archives. Install the tools, install uts, and you're done — no more Googling flags.

---

## Quick start

```bash
# Images
uts image compress screenshot.png -q low
uts image convert photo.heic --to jpg

# Video
uts video compress recording.mp4 -i
uts video convert clip.mov --to mp4

# Audio
uts audio compress podcast.wav -q high
uts audio convert track.flac --to mp3

# PDF
uts pdf compress thesis.pdf -q low
uts pdf convert report.pdf --to jpg

# Archives
uts archive compress ./project/ --algorithm zstd
uts archive extract backup.zip
uts archive list project.tar.gz

# Info
uts info video.mp4
```

---

## Install

```bash
# go install
go install github.com/y3owk1n/uts@latest

# nix flake
nix profile install github:y3owk1n/uts

# binary
# download from github.com/y3owk1n/uts/releases

# homebrew (coming soon)
brew install y3owk1n/tap/uts
```

### Dependencies

uts auto-detects whatever tools you already have installed and picks the best one for each format. Install what you need:

| Category | Recommended                                  |
|----------|----------------------------------------------|
| image    | `brew install imagemagick pngquant jpegoptim cwebp gifsicle` |
| video    | `brew install ffmpeg`                        |
| audio    | `brew install ffmpeg`                        |
| pdf      | `brew install ghostscript poppler imagemagick` |
| archive  | `brew install zstd xz brotli`                |

---

## The uts pattern

```
uts <category> <action> <files...> [flags]
```

| Category  | Actions                     | Formats                                                               |
|-----------|-----------------------------|-----------------------------------------------------------------------|
| `image`   | compress, convert           | png, jpg, webp, gif, bmp, tiff, heic, avif                           |
| `video`   | compress, convert           | mp4, mov, mkv, avi, webm, m4v, flv, wmv                              |
| `audio`   | compress, convert           | wav, flac, aac, mp3, m4a, opus, ogg, wma                             |
| `pdf`     | compress, convert           | PDF ↔ jpg/png, images → PDF                                          |
| `archive` | compress, extract, list     | zip, tar, tar.gz, tar.zst, tar.xz, tar.bz2                           |

Quality is controlled the same way across every category via `-q`:

| Preset  | Video (CRF) | Image (%) | Audio (kbps) | PDF (DPI) |
|---------|-------------|-----------|--------------|-----------|
| `high`  | 23          | 90        | 192          | 300/400   |
| `medium`| 28          | 80        | 128          | 150/300   |
| `low`   | 32          | 60        | 96           | 72/150    |
| `<n>`   | raw CRF     | raw %     | raw kbps     | raw DPI   |

---

## Flags

| Flag                    | Description                                           |
|-------------------------|-------------------------------------------------------|
| `-q, --quality`         | Quality preset or raw value                           |
| `-o, --output`          | Output directory                                      |
| `-i, --in-place`        | Replace original file                                 |
| `-n, --dry-run`         | Preview without executing                             |
| `-v, --verbose`         | Debug logging                                         |
| `-r, --recursive`       | Recursive glob patterns                               |
| `--algorithm`           | Archive algorithm (auto, gzip, zstd, xz, brotli, zip) |
| `--to`                  | Target format for conversion                          |

---

## How it works

uts detects the tools on your system and dispatches to the best one for each format:

| Format  | Compress                               | Convert                        |
|---------|----------------------------------------|--------------------------------|
| PNG     | pngquant + optipng                     | ImageMagick / sips             |
| JPEG    | jpegoptim                              | ImageMagick / sips             |
| WebP    | cwebp                                  | ImageMagick / sips             |
| GIF     | gifsicle                               | ImageMagick / sips             |
| HEIC    | heif-convert → JPEG                    | ImageMagick / sips             |
| AVIF    | cavif / avifenc                        | ImageMagick / sips             |
| Video   | ffmpeg (libx265)                       | ffmpeg                         |
| Audio   | ffmpeg (aac)                           | ffmpeg                         |
| PDF     | ghostscript                            | pdftoppm / ImageMagick         |
| Archive | tar + zstd/xz/gzip/brotli/zip          | tar / unzip / zstd / xz / bz2  |

No config files, no XML, no YAML — just the tools you already have, wrapped in a consistent CLI.

---

## Development

```bash
just build       # bin/uts
just test        # run tests
just lint        # golangci-lint
just genman      # man pages
```

---

## License

[MIT](LICENSE)
