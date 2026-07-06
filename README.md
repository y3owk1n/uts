<div align="center">

# uts

**One CLI to rule them all.**

_Compress, convert, and inspect any media file without remembering a dozen different command-line tools._

![Go Version](https://img.shields.io/github/go-mod/go-version/y3owk1n/uts?style=flat-square&logo=go)
[![Latest Release](https://img.shields.io/github/v/release/y3owk1n/uts?style=flat-square)](https://github.com/y3owk1n/uts/releases)
[![License](https://img.shields.io/github/license/y3owk1n/uts?style=flat-square)](LICENSE)

**Free and open-source.** Zero telemetry. Runs entirely local.

|      macOS       |      Linux       |
| :--------------: | :--------------: |
| Full featured ✅ | Full featured ✅ |

</div>

---

Media files are everywhere, but every format has its own specialized toolchain: `ffmpeg` for video/audio, `pngquant` for PNGs, `jpegoptim` for JPEGs, `ghostscript` for PDFs, and `tar` + various compressors for archives. Each tool has a different, complex command syntax that you have to look up every single time.

**uts** wraps all of them behind a single, predictable CLI pattern:

```bash
uts <category> <action> <input...> [flags]
```

That's it. The same pattern works for video, images, audio, PDFs, and archives. Install your encoders, install `uts`, and you're done — no more Googling flags.

---

## Quick start

```bash
# Compress image with a low-quality (small size) preset
uts image compress screenshot.png -q low

# Convert HEIC photo to JPG
uts image convert photo.heic --to jpg

# Compress video in-place (replaces original recording)
uts video compress recording.mp4 -i

# Convert MOV video to MP4
uts video convert clip.mov --to mp4

# Compress WAV audio to MP3 using high preset quality
uts audio convert track.wav --to mp3 -q high

# Compress PDF document to a medium-res profile
uts pdf compress draft.pdf -q medium

# Compress folder to a zstd tarball archive
uts archive compress ./project/ --algorithm zstd

# Inspect media file details and get format-specific recommendations
uts info recording.mp4
```

---

## Install

```bash
# Homebrew (Recommended)
brew install y3owk1n/tap/uts

# Via Go Toolchain
go install github.com/y3owk1n/uts@latest

# Pre-compiled Binary
# Download for your platform from github.com/y3owk1n/uts/releases
```

For complete instructions and environment details, see the [Installation Guide](docs/INSTALLATION.md).

---

## The uts Pattern

```bash
uts <category> <action> <files...> [flags]
```

### Supported Actions

| Category      | Actions                       | Formats                                          |
| ------------- | ----------------------------- | ------------------------------------------------ |
| **`image`**   | `compress`, `convert`         | png, jpg, jpeg, webp, gif, bmp, tiff, heic, avif |
| **`video`**   | `compress`, `convert`         | mp4, mov, mkv, avi, webm, m4v, flv, wmv          |
| **`audio`**   | `compress`, `convert`         | wav, flac, aac, mp3, m4a, opus, ogg, wma         |
| **`pdf`**     | `compress`, `convert`         | pdf, jpg, png                                    |
| **`archive`** | `compress`, `extract`, `list` | zip, tar, tar.gz, tar.zst, tar.xz, tar.bz2       |

### Quality Presets (`-q, --quality`)

Quality is mapped consistently across media formats using high-level presets:

| Level        | Video (CRF)        | Image (%)             | Audio (kbps)          | PDF (DPI / Setting)    |
| ------------ | ------------------ | --------------------- | --------------------- | ---------------------- |
| **`high`**   | `23` (Slow, best)  | `90%`                 | `192k`                | `400 DPI` (`/printer`) |
| **`medium`** | `28` (Default)     | `80%`                 | `128k`                | `300 DPI` (`/ebook`)   |
| **`low`**    | `32` (Fast, small) | `60%`                 | `96k`                 | `150 DPI` (`/screen`)  |
| **`<num>`**  | Raw CRF (`0-51`)   | Raw quality (`1-100`) | Raw kbps (e.g. `256`) | Raw DPI (e.g. `72`)    |

---

## Global Options

| Flag | Long Flag     | Description                                                              |
| ---- | ------------- | ------------------------------------------------------------------------ |
| `-q` | `--quality`   | Compression preset quality level or raw numeric value                    |
| `-o` | `--output`    | Destination directory (defaults to same folder as input)                 |
| `-i` | `--in-place`  | Overwrite the input file with the output file                            |
| `-n` | `--dry-run`   | Show commands that would be executed without running them                |
| `-v` | `--verbose`   | Output raw debug logs and CLI tool commands                              |
| `-r` | `--recursive` | Evaluate recursive glob patterns (e.g., `**/*.png`)                      |
|      | `--algorithm` | Archive type algorithm selection (`gzip`, `zstd`, `xz`, `brotli`, `zip`) |
|      | `--to`        | Target format file extension when performing conversion                  |

---

## How It Works

`uts` auto-detects whatever command-line tools you already have installed and dispatches the execution logic to the most suitable tool:

```
                  ┌───► [pngquant / optipng] (PNG)
                  ├───► [jpegoptim]          (JPEG)
uts image ───────►├───► [cwebp]              (WebP)
                  └───► [imagemagick]        (AVIF/HEIC)

uts video ───────────► [ffmpeg]              (MP4/MKV/WebM)
uts pdf ─────────────► [ghostscript]          (PDF Compression)
uts archive ─────────► [tar / zstd / zip]    (Archives)
```

For complete mapping behaviors and tool fallback mechanics, check out the [CLI Guide](docs/CLI.md).

---

## Documentation

Everything you need to go deep:

| Guide                                | What's in it                                                        |
| ------------------------------------ | ------------------------------------------------------------------- |
| [Installation](docs/INSTALLATION.md) | Go installation, Nix setup, bin releases, dependency setup          |
| [CLI Reference](docs/CLI.md)         | Details on commands, presets, exit codes, and extensive examples    |
| [Development](docs/DEVELOPMENT.md)   | Build tasks, testing tiers, code style, and release process         |
| [Contributing](CONTRIBUTING.md)      | Code of conduct, making changes, commit messages, and PR guidelines |

---

## Support & Contributing

Contributions are always welcome. `uts` is built in pure Go with a clean, modular structure. Check out [CONTRIBUTING.md](CONTRIBUTING.md) to get started.

If `uts` has earned a place in your toolbelt, please consider starring the repository! ⭐

---

## License

MIT — see [LICENSE](LICENSE).

<div align="center">
<br/>

**One command to rule them all.**

```bash
brew install y3owk1n/tap/uts && uts image compress screenshot.png -q low
```

<br/>

Made with ❤️ by <a href="https://github.com/y3owk1n">y3owk1n</a>

</div>
