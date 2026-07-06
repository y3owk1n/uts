package util

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func FileSize(path string) int64 {
	info, err := os.Stat(path)
	if err != nil {
		return 0
	}
	return info.Size()
}

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

func OutputPathExt(input, suffix, newExt string) string {
	dir := filepath.Dir(input)
	base := filepath.Base(input)
	name := strings.TrimSuffix(base, filepath.Ext(base))

	if name == "" {
		return filepath.Join(dir, "."+newExt)
	}
	return filepath.Join(dir, name+"-"+suffix+"."+newExt)
}

func ConvertOutputPath(input, targetExt string) string {
	dir := filepath.Dir(input)
	base := filepath.Base(input)
	name := strings.TrimSuffix(base, filepath.Ext(base))

	return filepath.Join(dir, name+"."+targetExt)
}

func MaybeInPlace(compressed, original string) {
	if inPlace, err := os.Stat(compressed); err == nil && inPlace != nil {
		os.Rename(compressed, original)
	}
}

func FileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

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

func EnsureDir(path string) error {
	dir := filepath.Dir(path)
	if dir != "." && dir != "/" {
		return os.MkdirAll(dir, 0755)
	}
	return nil
}
