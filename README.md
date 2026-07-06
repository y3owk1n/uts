<div align="center">

# uts

**The universal remote for your media files.**

_Stop memorizing complex commands. Compress, convert, and inspect any media file using the same simple pattern._

![Go Version](https://img.shields.io/github/go-mod/go-version/y3owk1n/uts?style=flat-square&logo=go)
[![Latest Release](https://img.shields.io/github/v/release/y3owk1n/uts?style=flat-square)](https://github.com/y3owk1n/uts/releases)
[![License](https://img.shields.io/github/license/y3owk1n/uts?style=flat-square)](LICENSE)

**100% free, open-source, and zero telemetry. Everything runs locally on your machine.**

|      macOS       |      Linux       |
| :--------------: | :--------------: |
| Full featured ✅ | Full featured ✅ |

</div>

https://github.com/user-attachments/assets/2732009f-8096-4a52-820f-f5bc735adb19

---

## Why `uts`?

Do you remember the exact `ffmpeg` flags to compress a video? What about `pngquant` for PNGs, or `ghostscript` for PDFs?

Every time you need to do something with a media file, you end up Googling, copy‑pasting, and praying the command works. It's a massive waste of time.

`uts` ends that chaos. It wraps all your everyday media operations into one dead‑simple command pattern:

```bash
uts <category> <action> <input...> [flags]
```

That’s it. The same pattern works for images, video, audio, PDFs, and archives. Install `uts`, install your encoders, and you'll never look up a media‑related command again.

---

## See It In Action

Here are a few real‑world examples that show how `uts` saves you time:

```bash
# Compress image with a low-quality (small size) preset
uts image compress screenshot.png -q low

# Convert HEIC photo to JPG
uts image convert photo.heic --to jpg

# Compress video in-place (replaces original recording)
uts video compress recording.mp4 -i

# Convert MOV video to MP4
uts video convert clip.mov --to mp4

# Convert multiple pages to high-res PNGs
uts pdf convert slides.pdf --to png -q high

# Stitch/Convert multiple images into a single PDF
uts pdf convert page1.png page2.png page3.jpg --to pdf

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

## Say Goodbye to Tool Anxiety

**Forget the command salad**
`uts` replaces the complicated syntax of `ffmpeg`, `pngquant`, `jpegoptim`, `ghostscript`, `tar`, and many others with a single, memorable pattern.

**Quality, simplified**
Stop fiddling with CRF values, DPI, and bitrates. `uts` gives you three intuitive presets via the `-q` (or `--quality`) flag:

| Preset       | When to use                                       |
| :----------- | :------------------------------------------------ |
| **`high`**   | Archival, printing, or when quality is paramount  |
| **`medium`** | Everyday use – a good balance of size and quality |
| **`low`**    | Quick sharing, tiny files, maximum compression    |

You can also pass a raw numeric value (e.g., `-q 75`) for fine‑grained control.

**Your data, your machine**
No cloud, no telemetry, no backdoors. `uts` is just a smart dispatcher that runs the right local tools (like `ffmpeg`) on your own hardware.

---

## Getting Started

```bash
# 1. Install uts (Homebrew is recommended)
brew install y3owk1n/tap/uts

# 2. Install the underlying tools (example for macOS)
brew install ffmpeg imagemagick pngquant jpegoptim ghostscript

# 3. Start using it!
uts image compress my-photo.jpg -q low
```

For complete installation instructions and dependency details, see the [Installation Guide](docs/INSTALLATION.md).

---

## Powerful Abilities, Simple Commands

`uts` organizes all operations into clear categories and actions:

| Category      | Actions                       | Supported formats (partial)                |
| :------------ | :---------------------------- | :----------------------------------------- |
| **`image`**   | `compress`, `convert`         | png, jpg, webp, gif, heic, avif, tiff, bmp |
| **`video`**   | `compress`, `convert`         | mp4, mov, mkv, avi, webm, m4v, flv, wmv    |
| **`audio`**   | `compress`, `convert`         | mp3, flac, aac, wav, m4a, opus, ogg, wma   |
| **`pdf`**     | `compress`, `convert`         | pdf, jpg, png                              |
| **`archive`** | `compress`, `extract`, `list` | zip, tar, tar.gz, tar.zst, tar.xz, tar.bz2 |

---

## How Does It Work?

`uts` is not a re‑encoding engine – it's an intelligent orchestrator. It detects which command‑line tools you already have installed and chooses the best one for each job.

```
                  ┌───► [pngquant / optipng] (PNG compression)
                  ├───► [jpegoptim]          (JPEG optimisation)
uts image ───────►├───► [cwebp]              (WebP conversion)
                  └───► [imagemagick]        (other formats)

uts video ───────────► [ffmpeg]              (video/audio processing)
uts pdf ─────────────► [ghostscript]          (PDF compression)
uts archive ─────────► [tar / zstd / zip]    (archiving)
```

You don't need to install all of them – `uts` will intelligently work with what you have.

---

## Advanced Usage

- **Batch processing**: Combine with the `-r` (or `--recursive`) flag to process entire directory trees.
- **Dry‑run preview**: Use `-n` (or `--dry-run`) to see exactly what commands would be executed, without making any changes.
- **In‑place replacement**: Use `-i` (or `--in-place`) to overwrite the original file.
- **Custom output directory**: Use `-o` (or `--output`) to specify where results are saved.

For a complete reference, check out the [CLI Guide](docs/CLI.md).

---

### Color Customization

`uts` supports full color palette customization via environment variables. See the [Configuration Guide](docs/CONFIGURATION.md) for details.

---

### Documentation & Contributions

- **[Installation](docs/INSTALLATION.md)** – Detailed setup for all platforms.
- **[CLI Reference](docs/CLI.md)** – Every command, flag, and exit code explained.
- **[Configuration](docs/CONFIGURATION.md)** – Color customization and environment variables.
- **[Development](docs/DEVELOPMENT.md)** – Build tasks, testing, and release process.
- **[Contributing](CONTRIBUTING.md)** – Code of conduct, commit guidelines, and pull request workflow.

`uts` is written in pure Go with a clean, modular architecture. Contributions of all kinds are warmly welcomed.

If `uts` has made your life easier, please consider starring the repository on GitHub ⭐ – it means a lot to us.

---

## License

MIT — see [LICENSE](LICENSE).

<div align="center">
<br/>

**Stop memorising. Start doing.**

```bash
brew install y3owk1n/tap/uts && uts image compress screenshot.png -q low
```

<br/>

Made with ❤️ by <a href="https://github.com/y3owk1n">y3owk1n</a>

</div>
