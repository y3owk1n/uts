# uts CLI Reference

`uts` is an all-in-one utility toolkit for compressing, converting, and managing media files.

## Usage

```
uts <category> <action> <input...> [options]
```

## Categories & Actions

| Category  | Actions                          | Description                        |
|-----------|----------------------------------|------------------------------------|
| `video`   | `compress`, `convert`            | Video files (mp4, mov, mkv, ...)   |
| `image`   | `compress`, `convert`            | Images (png, jpg, webp, heic, ...) |
| `pdf`     | `compress`, `convert`            | PDF documents                      |
| `audio`   | `compress`, `convert`            | Audio files (wav, flac, mp3, ...)  |
| `archive` | `compress`, `extract`, `list`    | Compress/extract/list archives     |
| `info`    | (top-level)                      | Show file info and suggestions     |
| `convert` | (top-level, shortcut)            | Direct format conversion           |

## Options

| Flag                  | Description                                      |
|-----------------------|--------------------------------------------------|
| `-q, --quality`       | Quality: low, medium, high, or numeric value     |
| `-o, --output <dir>`  | Output directory (default: same as input)        |
| `-i, --in-place`      | Replace original file with compressed version    |
| `-n, --dry-run`       | Show what would be done without doing it         |
| `-v, --verbose`       | Verbose output                                   |
| `-r, --recursive`     | Enable recursive glob patterns (`**/*.ext`)      |
| `--algorithm <name>`  | Archive algorithm (auto, gzip, zstd, xz, brotli, zip) |
| `--to <format>`       | Target format for conversion                     |
| `-h, --help`          | Show help                                        |
| `--version`           | Show version                                     |

## Quality Presets

| Level    | Video (CRF) | Audio (bitrate) | Image (quality) | PDF (DPI) |
|----------|-------------|-----------------|-----------------|-----------|
| `high`   | crf=23      | 192k            | 90              | 400       |
| `medium` | crf=28      | 128k            | 80              | 300       |
| `low`    | crf=32      | 96k             | 60              | 150       |
| `<num>`  | 0-51 (CRF)  | kbps (e.g. 256) | 1-100           | DPI value |

## Examples

### Video Compression
```
uts video compress screen-recording.mp4 -q low
uts video compress vacation.mov -q high -i
uts video compress lecture.mkv --dry-run
uts video compress clip1.mp4 clip2.mp4 clip3.mp4 -q medium
uts video compress '*.mp4' -r -q medium
uts video compress lecture.mkv -q 25 --dry-run -v
```

### Video Conversion
```
uts video convert clip.mov --to mp4
uts video convert recording.mkv --to webm -q medium
uts video convert presentation.avi --to mp4 -q 18
uts video convert clip1.mov clip2.mov --to mp4 -i
uts video convert '*.mov' --to mp4
```

### Image Compression
```
uts image compress screenshot.png -q medium
uts image compress logo.jpg -q high -i
uts image compress photo1.png photo2.png photo3.png -q low
uts image compress '*.png' -r
uts image compress '**/*.jpg' -r
uts image compress photo.heic -q low
uts image compress animation.gif
uts image compress photo.webp -q 75 --dry-run -v
```

### Image Conversion
```
uts image convert photo.heic --to jpg
uts image convert screenshot.png --to webp -q high
uts image convert photo.jpg --to avif -q 70
uts image convert photo1.heic photo2.heic --to jpg
uts image convert '*.heic' --to jpg
```

### PDF Compression
```
uts pdf compress thesis.pdf -q low
uts pdf compress report.pdf -q medium -o ./web/
uts pdf compress doc1.pdf doc2.pdf doc3.pdf -q low
uts pdf compress slides.pdf -q 300 --dry-run
uts pdf compress '*.pdf' -r
```

### PDF Conversion
```
uts pdf convert report.pdf --to jpg
uts pdf convert slides.pdf --to png -q high
uts pdf convert document.pdf --to jpg -q 200
uts pdf convert '*.jpg' '*.png' --to pdf
uts pdf convert images/*.png --to pdf
```

### Audio Compression
```
uts audio compress podcast.wav -q low
uts audio compress voice-memo.m4a -q high
uts audio compress track1.wav track2.flac track3.m4a -q medium
uts audio compress voice.wav -q 256 --dry-run
uts audio compress '*.wav' -r
```

### Audio Conversion
```
uts audio convert track.wav --to mp3
uts audio convert song.flac --to m4a -q high
uts audio convert track1.wav track2.flac --to mp3
uts audio convert '*.wav' --to mp3 -q 96
uts audio convert lecture.wav --to mp3 -q 256
```

### Archive Compression
```
uts archive compress ./project/ --algorithm zstd
uts archive compress ./data/ --algorithm gzip
uts archive compress ./src/ --dry-run
uts archive compress ./docs/ --algorithm brotli
uts archive compress ./photos/ --algorithm zip
```

### Archive Extraction
```
uts archive extract backup.zip
uts archive extract project.tar.gz
uts archive extract '*.tar.zst'
uts archive extract backup.zip -o ./output/
uts archive extract backup.zip --dry-run
```

### Archive Listing
```
uts archive list backup.zip
uts archive list project.tar.gz
uts archive list '*.tar.zst'
```

### File Info
```
uts info video.mp4
uts info '*.png'
uts info photo.heic
```

### Top-Level Convert Shortcuts
```
uts convert image photo.heic --to jpg
uts convert video clip.mov --to mp4
uts convert audio track.wav --to mp3 -q 96
uts convert pdf report.pdf --to jpg
```

## Output Behavior

- Files are saved as `<name>-small.<ext>` in the same directory by default.
- Use `-o, --output <dir>` to specify a different output directory.
- Use `-i, --in-place` to replace the original file with the compressed version.
- `--in-place` is ignored when used with `--output`.

## Exit Codes

| Code | Meaning          |
|------|------------------|
| 0    | Success          |
| 1    | General error    |
