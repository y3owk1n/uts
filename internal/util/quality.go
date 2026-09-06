//nolint:goconst,mnd
package util

import (
	"fmt"

	derrors "github.com/y3owk1n/uts/internal/core/errors"
)

// Numeric quality ranges accepted per category.
const (
	ImageQualityMin = 1
	ImageQualityMax = 100
	VideoCRFMin     = 0
	VideoCRFMax     = 51
	AudioKbpsMin    = 32
	AudioKbpsMax    = 512
	PDFDPIMin       = 36
	PDFDPIMax       = 1200
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

// ImageQuality converts a quality level to an image quality percentage.
func ImageQuality(level string) (int, error) {
	val, err := PresetVal(level, 60, 80, 90)
	if err != nil {
		return 0, err
	}

	return val, checkRange(val, ImageQualityMin, ImageQualityMax, "image quality")
}

// VideoQuality converts a quality level to CRF and bitrate values.
func VideoQuality(level string) (int, string, error) {
	if isNumeric(level) {
		crf := parseInt(level)

		err := checkRange(crf, VideoCRFMin, VideoCRFMax, "video CRF")
		if err != nil {
			return 0, "", err
		}

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
	kbps, err := PresetVal(level, 96, 128, 192)
	if err != nil {
		return "", err
	}

	err = checkRange(kbps, AudioKbpsMin, AudioKbpsMax, "audio bitrate (kbps)")
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%dk", kbps), nil
}

// PDFDPI converts a quality level to PDF DPI.
func PDFDPI(level string) (int, string, error) {
	if isNumeric(level) {
		dpi := parseInt(level)

		return dpi, "", checkRange(dpi, PDFDPIMin, PDFDPIMax, "PDF DPI")
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

func checkRange(val, low, high int, what string) error {
	if val < low || val > high {
		return derrors.Newf(
			derrors.CodeInvalidInput,
			"%s must be between %d and %d, got %d",
			what,
			low,
			high,
			val,
		)
	}

	return nil
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
