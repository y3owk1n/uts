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

Inputs are expanded before anything runs:

- Quoted glob patterns (`'*.png'`, `'**/*.jpg'` with `-r`) are resolved by `uts`, so they work the same in every shell.
- A directory expands to the files it contains that the command supports (an image command ignores `.txt` files, for example). With `-r` the whole tree is walked.
- Explicit file paths are used as given. A missing file is reported and counted as a failure.

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
| `-v` | `--verbose`   | Log every external command that runs, plus debug output.                         | `false`            |
|      | `--quiet`     | Print only warnings and errors. Spinners and progress are suppressed too.        | `false`            |
| `-r` | `--recursive` | Walk directories recursively and expand `**` glob patterns.                      | `false`            |
| `-h` | `--help`      | Display syntax, actions, and options helper info.                                |                    |
|      | `--version`   | Display version, commit hash, and build timestamp.                               |                    |

When both `--output` and `--in-place` are given, `--in-place` is ignored with a warning so originals are never deleted.

### Command Options

| Flag          | Commands                | Description                                               | Default |
| ------------- | ----------------------- | --------------------------------------------------------- | ------- |
| `--to`        | every `convert` command | Target format. Tab completion lists the valid values.     | _None_  |
| `--algorithm` | `archive compress`      | Archive algorithm: `zip`, `gzip`, `zstd`, `xz`, `brotli`. | `zip`   |
| `--max`       | image and video `compress` and `convert` | Shrink so the longest edge is at most this many pixels. Never enlarges, keeps the aspect ratio. Forces a re-encode for video. | _off_   |

---

## Quality Presets

The `-q, --quality` flag converts high-level presets (`low`, `medium`, `high`) to category-specific codec configurations under the hood:

| Level        | Video (CRF)                | Audio (Bitrate)          | Image (Quality)        | PDF (DPI / Preset)         |
| ------------ | -------------------------- | ------------------------ | ---------------------- | -------------------------- |
| **`high`**   | `23` (Slow, best quality)  | `192k`                   | `90%`                  | `400 DPI` (`/printer`)     |
| **`medium`** | `28` (Balanced speed/size) | `128k`                   | `80%`                  | `300 DPI` (`/ebook`)       |
| **`low`**    | `32` (Fast, smallest size) | `96k`                    | `60%`                  | `150 DPI` (`/screen`)      |
| **`<num>`**  | Raw CRF (`0` to `51`)      | Raw kbps (`32` to `512`) | Quality (`1` to `100`) | Raw DPI (`36` to `1200`)   |

A number outside the range for its category is rejected before any tool runs, with a message naming the valid range.

---

## Detailed Examples

### Video Commands

Compression uses `ffmpeg` and keeps the input container: mp4/mov/avi/flv are encoded with H.264, mkv with H.265, and webm with VP9 (constant quality mode). While a file encodes, the spinner shows the percentage done and the time remaining (video and audio, when `ffprobe` is available).

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

# Downscale a 4K screen recording to 1080p while compressing
uts video compress recording.mov --max 1920
```

Conversion changes the container. When the source codecs already fit the target (for example H.264 + AAC from a `.mov` into `.mp4`) the streams are copied without re-encoding, which is instant and lossless. Otherwise the video is re-encoded at the `-q` quality. Passing `-q` explicitly always re-encodes.

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

# Shrink photos so the longest edge is 2000px, then compress
uts image compress ./photos --max 2000 -r
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

Processes file compression and bitrates through `ffmpeg`. Lossy inputs (`mp3`, `ogg`, `opus`, `m4a`) are re-encoded in their own format, while lossless inputs (`wav`, `flac`) and everything else become AAC in an `.m4a` container:

```bash
# Compress WAV file down to small audio profile (writes podcast-small.m4a)
uts audio compress podcast.wav -q low

# Convert audio format to MP3
uts audio convert track.wav --to mp3

# Extract the audio track of a video
uts audio convert lecture.mp4 --to mp3

# Convert FLAC lossless format to M4A with high preset quality
uts audio convert song.flac --to m4a -q high
```

### Archive Commands

Builds, extracts, and lists zip, tar, tar.gz, tar.zst, tar.xz, tar.bz2 and tar.br archives. Compressed tarballs are streamed through a pipe, so no temporary files are written next to the archive:

```bash
# Compress a folder to zstd tarball
uts archive compress ./project/ --algorithm zstd

# Compress folder to standard zip file
uts archive compress ./photos/ --algorithm zip

# Extract a zip or tar archive into a directory
uts archive extract backup.zip -o ./restored/

# List contents of a compressed tarball
uts archive list project.tar.gz
```

### File Info

Inspects a media file's details and displays context-aware `uts` suggestions. When `ffprobe` is installed, video, audio and image files also show duration, resolution, frame rate, codecs and bit rate. PDFs show their page count via `pdfinfo` or Ghostscript:

```bash
uts info video.mp4
uts info screenshot.png
uts info ./downloads -r
```

`--json` prints one array with the same details, for scripts:

```bash
uts info clip.mp4 --json | jq '.[0] | {durationSeconds, width, height, videoCodec}'
```

Fields: `path`, `name`, `sizeBytes`, `size`, `ext`, `category`, `tool`, and when available `durationSeconds`, `duration`, `width`, `height`, `frameRate`, `videoCodec`, `audioCodec`, `sampleRate`, `channels`, `bitRateKbps`, `pages`. An unreadable file has an `error` field and makes the command exit `1`.

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

- **Default suffix**: Compression writes `<name>-small.<ext>` next to the input. Conversion writes `<name>.<target>` next to the input.
- **Output directories**: With `--output <dir>`, compression writes `<dir>/<name>.<ext>` and conversion writes `<dir>/<name>.<target>`.
- **In-place**: `--in-place` replaces the original with the result. For conversions the original is removed and the result takes its base name.
- **Not smaller**: When a compressed result is not smaller than the original, `uts` says so. With `--in-place` the original is kept and the result discarded.
- **Priority**: If both `--output` and `--in-place` are defined, `--in-place` is ignored to prevent accidental source deletions.
- **Failures**: A failed step removes its partial output, every remaining file is still processed, and `uts` exits `1`.
- **Dry run**: `--dry-run` prints the exact external commands that would run for each file.

---

## Exit Codes

`uts` returns these exit codes:

| Code    | Meaning                                                                            |
| ------- | ---------------------------------------------------------------------------------- |
| **`0`** | Success. All input files processed cleanly.                                        |
| **`1`** | Failure. One or more files failed, no input matched, or a required tool was missing. |

---

> [!TIP]
> To configure missing dependencies, check the [Installation Guide](INSTALLATION.md). If you'd like to check tasks or build target actions, see the [Development Guide](DEVELOPMENT.md).
