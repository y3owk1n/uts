package info

import (
	"context"
	"os/exec"
	"strconv"
	"strings"

	"github.com/y3owk1n/uts/internal/util"
)

// pdfPages returns the page count of a PDF using pdfinfo (poppler) or, when
// that is missing, Ghostscript. It returns 0 when neither tool is available
// or the file cannot be read.
func pdfPages(file string) int {
	if util.HasTool("pdfinfo") {
		out, err := exec.CommandContext(context.Background(), "pdfinfo", file).Output()
		if err == nil {
			for line := range strings.SplitSeq(string(out), "\n") {
				if after, ok := strings.CutPrefix(line, "Pages:"); ok {
					pages, _ := strconv.Atoi(strings.TrimSpace(after))

					return pages
				}
			}
		}
	}

	if util.HasTool("gs") {
		// PostScript string literal: escape backslashes and parentheses.
		escaped := strings.NewReplacer(`\`, `\\`, `(`, `\(`, `)`, `\)`).Replace(file)

		out, err := exec.CommandContext(context.Background(), "gs",
			"-q", "-dNODISPLAY", "-dNOSAFER",
			"-c", "("+escaped+") (r) file runpdfbegin pdfpagecount = quit").Output()
		if err == nil {
			lines := strings.Fields(string(out))
			if len(lines) > 0 {
				pages, _ := strconv.Atoi(lines[len(lines)-1])

				return pages
			}
		}
	}

	return 0
}
