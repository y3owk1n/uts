//nolint:goconst,mnd
package util

import (
	"fmt"

	derrors "github.com/y3owk1n/uts/internal/core/errors"
)

// PresetVal converts a quality preset name to its numeric value.
func PresetVal(level string, low, med, high int) (int, error) {
	if isNumeric(level) {
		return parseInt(level), nil
	}

	switch level {
	case "low":
		return low, nil
	case "medium":
		return med, nil
	case "high":
		return high, nil
	default:
		return 0, derrors.Newf(
			derrors.CodeInvalidInput,
			"invalid quality: %s (use low, medium, high, or a number)",
			level,
		)
	}
}

// VideoQuality converts a quality level to CRF and bitrate values.
func VideoQuality(level string) (int, string, error) {
	if isNumeric(level) {
		crf := parseInt(level)

		var preset string
		switch {
		case crf < 18:
			preset = "slow"
		case crf < 28:
			preset = "medium"
		default:
			preset = "fast"
		}

		return crf, preset, nil
	}

	switch level {
	case "low":
		return 32, "fast", nil
	case "medium":
		return 28, "medium", nil
	case "high":
		return 23, "slow", nil
	default:
		return 0, "", derrors.Newf(
			derrors.CodeInvalidInput,
			"invalid quality: %s (use low, medium, high, or CRF 0-51)",
			level,
		)
	}
}

// AudioBitrate converts a quality level to an audio bitrate.
func AudioBitrate(level string) (string, error) {
	if isNumeric(level) {
		return level + "k", nil
	}

	v, err := PresetVal(level, 96, 128, 192)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%dk", v), nil
}

// PDFDPI converts a quality level to PDF DPI.
func PDFDPI(level string) (int, string, error) {
	if isNumeric(level) {
		return parseInt(level), "", nil
	}

	switch level {
	case "low":
		return 150, "/screen", nil
	case "medium":
		return 300, "/ebook", nil
	case "high":
		return 400, "/printer", nil
	default:
		return 0, "", derrors.Newf(
			derrors.CodeInvalidInput,
			"invalid quality: %s (use low, medium, high, or DPI)",
			level,
		)
	}
}

func isNumeric(str string) bool {
	if str == "" {
		return false
	}

	for _, c := range str {
		if c < '0' || c > '9' {
			return false
		}
	}

	return true
}

func parseInt(str string) int {
	num := 0
	for _, c := range str {
		num = num*10 + int(c-'0')
	}

	return num
}
