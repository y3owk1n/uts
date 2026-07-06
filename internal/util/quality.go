package util

import "fmt"

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
		return 0, fmt.Errorf("invalid quality: %s (use low, medium, high, or a number)", level)
	}
}

func VideoQuality(level string) (crf int, preset string, err error) {
	if isNumeric(level) {
		crf = parseInt(level)
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
		return 0, "", fmt.Errorf("invalid quality: %s (use low, medium, high, or CRF 0-51)", level)
	}
}

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
		return 0, "", fmt.Errorf("invalid quality: %s (use low, medium, high, or DPI)", level)
	}
}

func isNumeric(s string) bool {
	if s == "" {
		return false
	}
	for _, c := range s {
		if c < '0' || c > '9' {
			return false
		}
	}
	return true
}

func parseInt(s string) int {
	n := 0
	for _, c := range s {
		n = n*10 + int(c-'0')
	}
	return n
}
