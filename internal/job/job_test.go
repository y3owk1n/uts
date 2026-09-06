//nolint:testpackage
package job

import (
	"context"
	"os/exec"
	"strconv"
	"strings"
	"testing"
	"time"
)

func parseFake(line string) (float64, bool) {
	value, ok := strings.CutPrefix(line, "done=")
	if !ok {
		return 0, false
	}

	secs, err := strconv.ParseFloat(value, 64)

	return secs, err == nil
}

// TestRunStepReportsProgress tests that a step with a Progress parser streams
// stdout through it and reports each parsed value against the total.
func TestRunStepReportsProgress(t *testing.T) {
	step := Step{
		Cmd: exec.CommandContext(
			context.Background(),
			"sh",
			"-c",
			`printf 'done=0.5\nnoise\ndone=1.5\n'`,
		),
		Progress: &Progress{Total: 3, Parse: parseFake},
	}

	var got []float64

	err := runStep(step, func(done, total float64) {
		if total != 3 {
			t.Errorf("total = %v; want 3", total)
		}

		got = append(got, done)
	})
	if err != nil {
		t.Fatalf("runStep: %v", err)
	}

	if len(got) != 2 || got[0] != 0.5 || got[1] != 1.5 {
		t.Errorf("reported %v; want [0.5 1.5]", got)
	}
}

// TestRunStepProgressFailure tests that a failing progress step surfaces its
// stderr in the error.
func TestRunStepProgressFailure(t *testing.T) {
	step := Step{
		Cmd: exec.CommandContext(
			context.Background(),
			"sh",
			"-c",
			`echo "boom happened" >&2; exit 3`,
		),
		Progress: &Progress{Total: 1, Parse: parseFake},
	}

	err := runStep(step, func(float64, float64) {})
	if err == nil || !strings.Contains(err.Error(), "boom happened") {
		t.Fatalf("runStep error = %v; want stderr included", err)
	}
}

// TestProgressText tests the percentage and ETA rendering.
func TestProgressText(t *testing.T) {
	tests := []struct {
		done, total float64
		elapsed     time.Duration
		want        string
	}{
		{50, 100, 10 * time.Second, " 50% · 0:10 left"},
		{25, 100, 30 * time.Second, " 25% · 1:30 left"},
		{100, 100, time.Minute, "100%"},
		{1, 100, time.Second, "  1%"},
		{150, 100, time.Second, "100%"},
	}
	for _, testCase := range tests {
		got := progressText(testCase.done, testCase.total, testCase.elapsed)
		if got != testCase.want {
			t.Errorf("progressText(%v, %v, %v) = %q; want %q",
				testCase.done, testCase.total, testCase.elapsed, got, testCase.want)
		}
	}
}

// TestClock tests h:mm:ss rendering.
func TestClock(t *testing.T) {
	if got := clock(3725 * time.Second); got != "1:02:05" {
		t.Errorf("clock(1h2m5s) = %q; want 1:02:05", got)
	}

	if got := clock(65 * time.Second); got != "1:05" {
		t.Errorf("clock(65s) = %q; want 1:05", got)
	}
}
