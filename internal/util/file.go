//nolint:mnd
package util

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
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
// When recursive is true, patterns containing "**" are expanded recursively.
func ResolveGlobs(args []string, recursive bool) []string {
	var result []string
	for _, arg := range args {
		if strings.ContainsAny(arg, "*?[") {
			if recursive && strings.Contains(arg, "**") {
				matches := resolveRecursiveGlob(arg)
				if len(matches) > 0 {
					result = append(result, matches...)
				}

				continue
			}

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

// resolveRecursiveGlob handles "**" glob patterns by walking the directory tree.
func resolveRecursiveGlob(pattern string) []string {
	parts := strings.SplitN(pattern, "**", 2)
	if len(parts) != 2 {
		return nil
	}

	root := parts[0]
	suffix := parts[1]

	// Clean leading separator from suffix so it matches like "/*.png".
	suffix = strings.TrimPrefix(suffix, "/")

	if root == "" {
		root = "."
	}

	var matches []string

	//nolint:nilerr // Walk errors are non-fatal; empty result is fine for globs.
	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil
		}

		if d.IsDir() {
			return nil
		}

		match, _ := filepath.Match(suffix, filepath.Base(path))
		if match {
			matches = append(matches, path)
		}

		return nil
	})
	if err != nil {
		return nil
	}

	return matches
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

// HasTool reports whether the named executable is in PATH.
func HasTool(name string) bool {
	_, err := exec.LookPath(name)

	return err == nil
}

// CalcOutputPath computes an output path respecting the output directory.
// When outputDir is non-empty the file is placed there under its own name;
// otherwise it stays next to the input with the suffix inserted before the
// extension.
func CalcOutputPath(input, suffix, outputDir string) string {
	if outputDir != "" {
		return filepath.Join(outputDir, filepath.Base(input))
	}

	return OutputPath(input, suffix)
}

// CalcOutputPathExt is CalcOutputPath for operations that change the
// extension while keeping the "-suffix" naming (e.g. audio compress to .m4a).
func CalcOutputPathExt(input, suffix, newExt, outputDir string) string {
	if outputDir != "" {
		return filepath.Join(outputDir, ConvertOutputPath(filepath.Base(input), newExt))
	}

	return OutputPathExt(input, suffix, newExt)
}

// CalcConvertOutputPath computes a conversion output path respecting the
// output directory. The result always carries the target extension.
func CalcConvertOutputPath(input, targetExt, outputDir string) string {
	if outputDir != "" {
		return filepath.Join(outputDir, ConvertOutputPath(filepath.Base(input), targetExt))
	}

	return ConvertOutputPath(input, targetExt)
}

// MaybeReplaceOrRemove handles in-place replacement after processing.
// If compressed has the same extension as original it is renamed in place;
// otherwise the original is deleted (used for conversions to different formats).
func MaybeReplaceOrRemove(compressed, original string) {
	if filepath.Ext(compressed) == filepath.Ext(original) {
		MaybeInPlace(compressed, original)
	} else {
		RemoveInPlace(original)
	}
}

// ExpandInputs turns CLI arguments into a flat list of file paths.
//
// Glob patterns are expanded (recursive "**" patterns when recursive is set),
// and directories are replaced by the files they contain: their immediate
// children by default, or the whole tree when recursive is set. When exts is
// non-empty only files with one of those extensions are kept from expanded
// directories, so a stray .DS_Store never becomes an "unsupported format"
// warning. Explicit file arguments are passed through untouched so a missing
// file is still reported by the caller.
func ExpandInputs(args []string, recursive bool, exts []string) []string {
	var result []string

	for _, path := range ResolveGlobs(args, recursive) {
		info, err := os.Stat(path)
		if err != nil || !info.IsDir() {
			result = append(result, path)

			continue
		}

		result = append(result, listDir(path, recursive, exts)...)
	}

	return result
}

func listDir(dir string, recursive bool, exts []string) []string {
	var files []string

	keep := func(path string) {
		if len(exts) == 0 ||
			slices.Contains(exts, strings.ToLower(strings.TrimPrefix(filepath.Ext(path), "."))) {
			files = append(files, path)
		}
	}

	if !recursive {
		entries, err := os.ReadDir(dir)
		if err != nil {
			return nil
		}

		for _, entry := range entries {
			if !entry.IsDir() {
				keep(filepath.Join(dir, entry.Name()))
			}
		}

		return files
	}

	//nolint:nilerr // Unreadable entries are skipped, not fatal.
	_ = filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil
		}

		if !d.IsDir() {
			keep(path)
		}

		return nil
	})

	return files
}

// CopyFile copies src to dst, creating or truncating dst.
func CopyFile(src, dst string) error {
	source, err := os.Open(src)
	if err != nil {
		return err
	}
	defer source.Close() //nolint:errcheck

	out, err := os.Create(dst)
	if err != nil {
		return err
	}

	_, err = io.Copy(out, source)
	if err != nil {
		_ = out.Close()

		return err
	}

	return out.Close()
}

// Describe renders a command the way a user would type it: program name
// followed by its arguments, quoting any argument that contains whitespace.
func Describe(cmd *exec.Cmd) string {
	parts := make([]string, 0, len(cmd.Args))

	for idx, arg := range cmd.Args {
		if idx == 0 {
			arg = filepath.Base(arg)
		}

		if strings.ContainsAny(arg, " \t'\"") {
			arg = fmt.Sprintf("%q", arg)
		}

		parts = append(parts, arg)
	}

	return strings.Join(parts, " ")
}

// MagickBin returns the ImageMagick executable to use: "magick" (v7),
// "convert" (v6), or "" when neither is installed.
func MagickBin() string {
	if HasTool("magick") {
		return "magick"
	}

	if HasTool("convert") {
		return "convert"
	}

	return ""
}
