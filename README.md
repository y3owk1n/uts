# uts

**All-in-one utility toolkit** – Compress, convert, and manage media files from the command line.

[![CI](https://github.com/y3owk1n/uts/actions/workflows/ci.yml/badge.svg)](https://github.com/y3owk1n/uts/actions/workflows/ci.yml)
[![License](https://img.shields.io/github/license/y3owk1n/uts.svg)](LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/y3owk1n/uts)](https://goreportcard.com/report/github.com/y3owk1n/uts)

---

## Installation

### Homebrew (coming soon)

```bash
brew install y3owk1n/tap/uts
```

### Go install

```bash
go install github.com/y3owk1n/uts@latest
```

### Binary releases

Download the latest binary from the [releases page](https://github.com/y3owk1n/uts/releases).

### Nix flake

```bash
nix profile install github:y3owk1n/uts
```

### Build from source

```bash
git clone https://github.com/y3owk1n/uts.git
cd uts
nix develop  # or: devbox shell
just build
./bin/uts --help
```

---

## Quick start

```bash
# Image operations
uts image compress screenshot.png -q low
uts image convert photo.heic --to jpg

# Video operations
uts video compress recording.mp4 -i
uts video convert clip.mov --to mp4

# Audio operations
uts audio compress podcast.wav -q high
uts audio convert track.flac --to mp3

# PDF operations
uts pdf compress thesis.pdf -q low
uts pdf convert report.pdf --to jpg

# Archive operations
uts archive compress ./project/ --algorithm zstd
uts archive extract backup.zip
uts archive list project.tar.gz

# File info
uts info video.mp4
```

---

## Categories

| Category  | Operations                          | Formats                                                                 |
|-----------|-------------------------------------|-------------------------------------------------------------------------|
| `image`   | compress, convert                   | png, jpg, webp, gif, bmp, tiff, heic, avif                             |
| `video`   | compress, convert                   | mp4, mov, mkv, avi, webm, m4v, flv, wmv                                |
| `audio`   | compress, convert                   | wav, flac, aac, mp3, m4a, opus, ogg, wma                               |
| `pdf`     | compress, convert                   | PDF → jpg/png, images → PDF                                            |
| `archive` | compress, extract, list             | zip, tar, tar.gz, tar.zst, tar.xz, tar.bz2                             |

---

## Quality presets

All compress and convert commands accept `-q` / `--quality` with presets or raw values:

| Preset  | Video (CRF) | Image (%) | Audio (kbps) | PDF (DPI) |
|---------|-------------|-----------|--------------|-----------|
| `high`  | 23 (slow)   | 90        | 192          | 300/400   |
| `medium`| 28 (medium) | 80        | 128          | 150/300   |
| `low`   | 32 (fast)   | 60        | 96           | 72/150    |
| `<n>`   | raw CRF     | raw %     | raw kbps     | raw DPI   |

---

## Tools

uts detects and uses the best available tool for each format:

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

---

## Flags

| Flag                    | Description                                           |
|-------------------------|-------------------------------------------------------|
| `-q, --quality`         | Quality preset: low, medium, high, or numeric         |
| `-o, --output`          | Output directory (default: same as input)             |
| `-i, --in-place`        | Replace original file with compressed version         |
| `-n, --dry-run`         | Show what would be done without doing it              |
| `-v, --verbose`         | Enable verbose/debug logging                          |
| `-r, --recursive`       | Enable recursive glob patterns                        |
| `--algorithm`           | Archive algorithm: auto, gzip, zstd, xz, brotli, zip |
| `--to`                  | Target format for conversion                          |

---

## Development

### Prerequisites

- Go 1.26+
- [just](https://github.com/casey/just) command runner
- [devbox](https://www.jetify.com/devbox/) (optional, for reproducible shell)

### Commands

```bash
just build       # Build binary to bin/uts
just test        # Run all tests
just test-unit   # Run unit tests
just lint        # Run golangci-lint
just fmt         # auto-format code
just genman      # Generate man pages (build/man/)
just clean       # Remove build artifacts
```

### Release

```bash
just release-ci v1.2.3
```

Builds cross-platform binaries for darwin/linux/windows × arm64/amd64.

---

## Similar projects

- [nvs](https://github.com/y3owk1n/nvs) – Neovim Version Switcher (same UI toolkit)

---

## License

[MIT](LICENSE)
