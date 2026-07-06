# CLI Reference

`uts` is a unified command-line media tool that wraps multiple specialized utilities under a single, consistent syntax. This guide provides a detailed reference for all commands, parameters, and flags.

---

## Usage

```bash
uts <category> <action> <input...> [options]
```

Where:

- `<category>` defines the type of media file (e.g. `video`, `image`, `pdf`).
- `<action>` specifies what to do (e.g. `compress`, `convert`, `extract`).
- `<input...>` is one or more file paths, directories, or glob patterns.
- `[options]` are flags modifying the default behavior.

---

## Categories & Actions

| Category      | Actions                       | Aliases     | Description                                                  |
| ------------- | ----------------------------- | ----------- | ------------------------------------------------------------ |
| **`video`**   | `compress`, `convert`         | `v`         | Video processing (mp4, mov, mkv, webm, avi, etc.)            |
| **`image`**   | `compress`, `convert`         | `img`, `i`  | Image processing (png, jpg, webp, gif, heic, avif, etc.)     |
| **`pdf`**     | `compress`, `convert`         | `p`         | PDF document compression and format conversions              |
| **`audio`**   | `compress`, `convert`         | `a`         | Audio processing (wav, flac, mp3, m4a, opus, etc.)           |
| **`archive`** | `compress`, `extract`, `list` | `arc`, `ar` | Creating, unpacking, and inspecting compression archives     |
| **`info`**    | _N/A_ (Direct action)         | _None_      | Inspects file metadata and suggests action wrappers          |
| **`convert`** | _N/A_ (Direct action)         | `x`         | Top-level shortcut to convert images, videos, audio, or PDFs |

---

## Global Options

These options apply to commands across all categories:

| Flag | Long Flag     | Description                                                                      | Default            |
| ---- | ------------- | -------------------------------------------------------------------------------- | ------------------ |
| `-q` | `--quality`   | Compression preset quality level (`low`, `medium`, `high`) or raw numeric value. | `medium`           |
| `-o` | `--output`    | Destination directory.                                                           | Same as input file |
| `-i` | `--in-place`  | Replace the original source file with the processed version.                     | `false`            |
| `-n` | `--dry-run`   | Print the compiled command without executing it.                                 | `false`            |
| `-v` | `--verbose`   | Output raw debug logs and tool output commands.                                  | `false`            |
| `-r` | `--recursive` | Evaluate recursive glob patterns (e.g., `**/*.png`).                             | `false`            |
|      | `--algorithm` | Archive type algorithm selection (e.g., `gzip`, `zstd`, `xz`, `brotli`, `zip`).  | `auto`             |
|      | `--to`        | Target format file extension when performing conversion.                         | _None_             |
| `-h` | `--help`      | Display syntax, actions, and options helper info.                                |                    |
|      | `--version`   | Display version, commit hash, and build timestamp.                               |                    |

---

## Quality Presets

The `-q, --quality` flag converts high-level presets (`low`, `medium`, `high`) to category-specific codec configurations under the hood:

| Level        | Video (CRF)                | Audio (Bitrate)          | Image (Quality)        | PDF (DPI / Preset)         |
| ------------ | -------------------------- | ------------------------ | ---------------------- | -------------------------- |
| **`high`**   | `23` (Slow, best quality)  | `192k`                   | `90%`                  | `400 DPI` (`/printer`)     |
| **`medium`** | `28` (Balanced speed/size) | `128k`                   | `80%`                  | `300 DPI` (`/ebook`)       |
| **`low`**    | `32` (Fast, smallest size) | `96k`                    | `60%`                  | `150 DPI` (`/screen`)      |
| **`<num>`**  | Raw CRF (`0` to `51`)      | Raw bitrate (e.g. `256`) | Quality (`1` to `100`) | Raw DPI (e.g. `72`, `600`) |

---

## Detailed Examples

### Video Commands

Compression uses `ffmpeg` with `libx265` to output highly compressed HEVC files.

```bash
# Compress using low quality preset (smaller file)
uts video compress screen-recording.mp4 -q low

# Compress in place (replaces original vacation.mov with compressed file)
uts video compress vacation.mov -q high -i

# Preview compression commands without editing files
uts video compress lecture.mkv --dry-run

# Compress multiple input files concurrently
uts video compress clip1.mp4 clip2.mp4 clip3.mp4 -q medium

# Compress all MP4 files in subdirectories recursively
uts video compress '*.mp4' -r -q medium
```

Conversion changes format extensions utilizing the proper encoder wrappers:

```bash
# Convert QuickTime MOV to MP4
uts video convert clip.mov --to mp4

# Convert MKV to WebM with quality presets
uts video convert recording.mkv --to webm -q medium

# Convert in-place and replace original files
uts video convert clip1.mov clip2.mov --to mp4 -i
```

### Image Commands

Leverages `pngquant`, `jpegoptim`, `cwebp`, `gifsicle`, `optipng`, or `imagemagick` based on format support and tool availability:

```bash
# Compress PNG using pngquant/optipng
uts image compress screenshot.png -q medium

# Compress image and overwrite the original file
uts image compress logo.jpg -q high -i

# Find and compress all JPG files recursively
uts image compress '**/*.jpg' -r

# Compress HEIC photos down to a smaller profile size
uts image compress photo.heic -q low
```

Conversion handles format shifting (e.g. modern formatting like `avif` or `webp`):

```bash
# Convert HEIC to JPG format
uts image convert photo.heic --to jpg

# Convert PNG to WebP with custom quality configuration
uts image convert screenshot.png --to webp -q high

# Batch convert files matching a glob pattern
uts image convert '*.heic' --to jpg
```

### PDF Commands

Leverages `ghostscript` for compression and `poppler` (`pdftoppm`) / `imagemagick` for format conversions:

```bash
# Compress PDF down to a low-res web profile
uts pdf compress thesis.pdf -q low

# Output compressed PDF to a custom folder
uts pdf compress report.pdf -q medium -o ./web/

# Convert PDF document pages to individual JPG images
uts pdf convert report.pdf --to jpg

# Convert multiple pages to high-res PNGs
uts pdf convert slides.pdf --to png -q high

# Stitch/Convert multiple images into a single PDF
uts pdf convert page1.png page2.png page3.jpg --to pdf
```

### Audio Commands

Processes file compression and bitrates through `ffmpeg`:

```bash
# Compress WAV file down to small audio profile
uts audio compress podcast.wav -q low

# Convert audio format to MP3
uts audio convert track.wav --to mp3

# Convert FLAC lossless format to M4A with high preset quality
uts audio convert song.flac --to m4a -q high
```

### Archive Commands

Builds, extracts, and lists tar, zip, gzip, zstd, xz, and brotli packages:

```bash
# Compress a folder to zstd tarball
uts archive compress ./project/ --algorithm zstd

# Compress folder to standard zip file
uts archive compress ./photos/ --algorithm zip

# Extract a zip or tar archive
uts archive extract backup.zip

# List contents of a compressed tarball
uts archive list project.tar.gz
```

### File Info

Inspects a media file's details (duration, codecs, resolution, file size) and displays context-aware `uts` suggestions:

```bash
uts info video.mp4
uts info screenshot.png
```

### Top-Level Shortcut

Provides a fast routing proxy to the format converter:

```bash
uts convert image photo.heic --to jpg
uts convert video clip.mov --to mp4
uts convert audio track.wav --to mp3 -q 96
uts convert pdf report.pdf --to jpg
```

---

## Output Behavior

- **Default suffix**: Unless `--in-place` is specified, files are written as `<name>-small.<ext>` in the source directory.
- **Output directories**: When `--output <dir>` is set, target files are written to the destination directory. Suffixes are omitted if output name collisions do not occur.
- **Priority**: If both `--output` and `--in-place` are defined, `--in-place` is ignored to prevent accidental source deletions.

---

## Exit Codes

`uts` returns these exit codes:

| Code    | Meaning                                                                            |
| ------- | ---------------------------------------------------------------------------------- |
| **`0`** | Success. All input files processed cleanly.                                        |
| **`1`** | Failure. One or more operations failed, or required dependency tools were missing. |

---

> [!TIP]
> To configure missing dependencies, check the [Installation Guide](INSTALLATION.md). If you'd like to check tasks or build target actions, see the [Development Guide](DEVELOPMENT.md).
