//nolint:goconst
package convert

import (
	"fmt"
	"slices"

	derrors "github.com/y3owk1n/uts/internal/core/errors"
	"github.com/y3owk1n/uts/internal/ffmpeg"
	"github.com/y3owk1n/uts/internal/format"
	"github.com/y3owk1n/uts/internal/job"
	"github.com/y3owk1n/uts/internal/ui"
	"github.com/y3owk1n/uts/internal/util"
)

// VideoOptions represents options for video conversion.
type VideoOptions struct {
	Files   []string
	Target  string
	Quality string
	// QualitySet is true when the user passed -q explicitly. Without it a
	// conversion prefers a lossless stream copy whenever the codecs fit the
	// target container.
	QualitySet bool
	OutputDir  string
	InPlace    bool
	DryRun     bool
	// Jobs, SkipExisting and Backup are passed through to job.Options.
	Jobs         int
	SkipExisting bool
	Backup       bool
	// MaxEdge caps the longest edge in pixels; it always implies a re-encode.
	MaxEdge int
}

// Video converts video files to the target container.
func Video(opts VideoOptions) error {
	target := format.Normalize(opts.Target)
	if !slices.Contains(format.VideoTargets, target) {
		return unsupportedTarget(target, format.VideoTargets)
	}

	err := ffmpeg.Check()
	if err != nil {
		return err
	}

	if target == "gif" {
		return videoToGif(opts)
	}

	crf, preset, err := util.VideoQuality(opts.Quality)
	if err != nil {
		return err
	}

	vcodec, acodec := format.VideoCodecs(target)
	ui.Message.Infof(
		"Converting video to .%s (%s/%s, crf=%d)%s",
		target,
		vcodec,
		acodec,
		crf,
		maxNote(opts.MaxEdge),
	)

	return job.Run(opts.Files, job.Options{
		Verb:         "Converting",
		Done:         "Converted",
		Noun:         "video file",
		InPlace:      opts.InPlace,
		DryRun:       opts.DryRun,
		Jobs:         opts.Jobs,
		SkipExisting: opts.SkipExisting,
		Backup:       opts.Backup,
		Code:         derrors.CodeConversionFailed,
	}, func(file string) (*job.Job, error) {
		plan, err := videoJob(file, target, opts.OutputDir)
		if err != nil || plan.Skip != "" {
			return plan, err
		}

		if !opts.QualitySet && opts.MaxEdge == 0 {
			info, probeErr := ffmpeg.Probe(file)
			if probeErr == nil && ffmpeg.CanRemux(info, target) {
				plan.Steps = []job.Step{
					ffmpeg.Step(file, ffmpeg.RemuxArgs(file, plan.Output, target)...),
				}
				plan.Note = "stream copy, no re-encode (pass -q to force re-encoding)"

				return plan, nil
			}
		}

		plan.Steps = []job.Step{
			ffmpeg.Step(
				file,
				ffmpeg.EncodeArgs(file, plan.Output, target, crf, preset, opts.MaxEdge)...),
		}
		plan.Note = fmt.Sprintf("%s/%s crf=%d preset=%s", vcodec, acodec, crf, preset)

		return plan, nil
	})
}

// videoJob checks a convert input and returns a job with its output path
// set, or a skip job when the file is already in the target format.
func videoJob(file, target, outputDir string) (*job.Job, error) {
	ext := format.Ext(file)
	if !slices.Contains(format.VideoConvertExts, ext) {
		return nil, derrors.Newf(
			derrors.CodeUnsupportedFormat,
			"unsupported video format .%s",
			ext,
		)
	}

	if format.Same(ext, target) {
		return &job.Job{Input: file, Skip: "Already ." + target}, nil
	}

	return &job.Job{
		Input:  file,
		Output: util.CalcConvertOutputPath(file, target, outputDir),
	}, nil
}

// videoToGif renders each video as an animated GIF with ffmpeg, then runs
// gifsicle on the result when it is installed.
func videoToGif(opts VideoOptions) error {
	fps, dither, err := util.GifQuality(opts.Quality)
	if err != nil {
		return err
	}

	ui.Message.Infof(
		"Converting video to .gif (%d fps, %s dither)%s",
		fps,
		dither,
		maxNote(opts.MaxEdge),
	)

	return job.Run(opts.Files, job.Options{
		Verb:         "Converting",
		Done:         "Converted",
		Noun:         "video file",
		InPlace:      opts.InPlace,
		DryRun:       opts.DryRun,
		Jobs:         opts.Jobs,
		SkipExisting: opts.SkipExisting,
		Backup:       opts.Backup,
		Code:         derrors.CodeConversionFailed,
	}, func(file string) (*job.Job, error) {
		plan, err := videoJob(file, "gif", opts.OutputDir)
		if err != nil || plan.Skip != "" {
			return plan, err
		}

		plan.Steps = []job.Step{
			ffmpeg.Step(file, ffmpeg.GifArgs(file, plan.Output, fps, dither, opts.MaxEdge)...),
		}
		plan.Note = fmt.Sprintf("palettegen, %d fps", fps)

		if util.HasTool("gifsicle") {
			plan.Steps = append(
				plan.Steps,
				job.Exec("gifsicle", "-b", "-O3", "--lossy=80", plan.Output),
			)
			plan.Note += ", optimized with gifsicle"
		}

		return plan, nil
	})
}
