//nolint:mnd
package util

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// FileSize returns the size of the file at the given path.
func FileSize(path string) int64 {
	info, err := os.Stat(path)
	if err != nil {
		return 0
	}

	return info.Size()
}

// HumanSize formats a byte count as a human-readable string.
func HumanSize(bytes int64) string {
	switch {
	case bytes < 1024:
		return fmt.Sprintf("%d B", bytes)
	case bytes < 1048576:
		kb := bytes / 1024
		frac := (bytes % 1024) * 10 / 1024

		return fmt.Sprintf("%d.%d KB", kb, frac)
	case bytes < 1073741824:
		mb := bytes / 1048576
		frac := (bytes % 1048576) * 10 / 1048576

		return fmt.Sprintf("%d.%d MB", mb, frac)
	default:
		gb := bytes / 1073741824
		rem := (bytes % 1073741824) * 100 / 1073741824

		return fmt.Sprintf("%d.%d GB", gb, rem)
	}
}

// CompressionRatio returns the compression ratio as a formatted string.
func CompressionRatio(orig, compressed int64) string {
	if orig == 0 {
		return "0%"
	}

	pct := (orig - compressed) * 1000 / orig
	whole := pct / 10

	frac := pct % 10
	if pct == 0 {
		return "(0.0%)"
	}

	if pct < 0 {
		whole = -whole
		frac = -frac

		return fmt.Sprintf("(+%d.%d%%)", whole, frac)
	}

	return fmt.Sprintf("(-%d.%d%%)", whole, frac)
}

// OutputPath generates an output path by inserting a suffix before the extension.
func OutputPath(input, suffix string) string {
	dir := filepath.Dir(input)
	base := filepath.Base(input)
	ext := filepath.Ext(base)
	name := strings.TrimSuffix(base, ext)

	if name == "" {
		name = "." + ext[1:]

		return filepath.Join(dir, name+"-"+suffix)
	}

	return filepath.Join(dir, name+"-"+suffix+ext)
}

// OutputPathExt generates an output path with a new extension.
func OutputPathExt(input, suffix, newExt string) string {
	dir := filepath.Dir(input)
	base := filepath.Base(input)
	name := strings.TrimSuffix(base, filepath.Ext(base))

	if name == "" {
		return filepath.Join(dir, "."+newExt)
	}

	return filepath.Join(dir, name+"-"+suffix+"."+newExt)
}

// ConvertOutputPath converts a file path to a new extension.
func ConvertOutputPath(input, targetExt string) string {
	dir := filepath.Dir(input)
	base := filepath.Base(input)
	name := strings.TrimSuffix(base, filepath.Ext(base))

	return filepath.Join(dir, name+"."+targetExt)
}

// MaybeInPlace renames the compressed file to the original if compression succeeded.
func MaybeInPlace(compressed, original string) {
	inPlace, err := os.Stat(compressed)
	if err == nil && inPlace != nil {
		_ = os.Rename(compressed, original)
	}
}

// RemoveInPlace deletes the original file after a successful in-place conversion
// where the output has a different extension than the input.
func RemoveInPlace(original string) {
	_ = os.Remove(original)
}

// FileExists reports whether the given path exists.
func FileExists(path string) bool {
	_, err := os.Stat(path)

	return err == nil
}

// ResolveGlobs resolves glob patterns and returns matching file paths.
func ResolveGlobs(args []string, recursive bool) []string {
	var result []string
	for _, arg := range args {
		if strings.ContainsAny(arg, "*?[") {
			matches, err := filepath.Glob(arg)
			if err != nil || len(matches) == 0 {
				continue
			}

			result = append(result, matches...)
		} else {
			result = append(result, arg)
		}
	}

	return result
}

// EnsureDir ensures the parent directory of the given path exists.
func EnsureDir(path string) error {
	dir := filepath.Dir(path)
	if dir != "." && dir != "/" {
		return os.MkdirAll(dir, 0o755)
	}

	return nil
}

// InPlaceHint returns " (in-place)" when inPlace is true, or "" otherwise.
// Use it to append an in-place indicator to dry-run messages.
func InPlaceHint(inPlace bool) string {
	if inPlace {
		return " (in-place)"
	}

	return ""
}
