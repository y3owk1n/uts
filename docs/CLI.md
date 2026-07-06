# uts CLI Reference

## Usage

```
uts <category> <action> <input...> [options]
```

## Categories

- `video` - Video files (mp4, mov, mkv, avi, webm)
- `image` - Images (png, jpg, webp, gif, bmp, tiff, heic, avif)
- `pdf` - PDF documents
- `audio` - Audio files (wav, flac, aac, mp3, m4a, opus)
- `archive` - Compress/extract/list archives

## Actions

- `compress` - Compress files
- `convert` - Convert between formats
- `extract` - Extract archives (archive only)
- `list` - List archive contents (archive only)
- `info` - Show file info and suggestions (top-level)

## Global Options

- `-q, --quality` - Quality: low, medium, high, or number
- `-o, --output` - Output directory
- `-i, --in-place` - Replace original
- `-n, --dry-run` - Preview without changes
- `-v, --verbose` - Verbose output
- `-r, --recursive` - Recursive glob
- `--algorithm` - Archive algorithm
- `--to` - Target format for conversion
- `--version` - Show version
