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
}

// Video converts video files to the target container.
func Video(opts VideoOptions) error {
	target := format.Normalize(opts.Target)
	if !slices.Contains(format.VideoTargets, target) {
		return unsupportedTarget(target, format.VideoTargets)
	}

	crf, preset, err := util.VideoQuality(opts.Quality)
	if err != nil {
		return err
	}

	err = ffmpeg.Check()
	if err != nil {
		return err
	}

	vcodec, acodec := format.VideoCodecs(target)
	ui.Message.Infof("Converting video to .%s (%s/%s, crf=%d)", target, vcodec, acodec, crf)

	return job.Run(opts.Files, job.Options{
		Verb:    "Converting",
		Done:    "Converted",
		Noun:    "video files",
		InPlace: opts.InPlace,
		DryRun:  opts.DryRun,
		Code:    derrors.CodeConversionFailed,
	}, func(file string) (*job.Job, error) {
		ext := format.Ext(file)
		if format.Classify(ext) != format.Video {
			return nil, derrors.Newf(
				derrors.CodeUnsupportedFormat,
				"unsupported video format .%s",
				ext,
			)
		}

		if format.Same(ext, target) {
			return &job.Job{Input: file, Skip: "Already ." + target}, nil
		}

		out := util.CalcConvertOutputPath(file, target, opts.OutputDir)

		if !opts.QualitySet {
			info, probeErr := ffmpeg.Probe(file)
			if probeErr == nil && ffmpeg.CanRemux(info, target) {
				return &job.Job{
					Input:  file,
					Output: out,
					Steps:  []job.Step{job.Exec("ffmpeg", ffmpeg.RemuxArgs(file, out, target)...)},
					Note:   "stream copy, no re-encode (pass -q to force re-encoding)",
				}, nil
			}
		}

		return &job.Job{
			Input:  file,
			Output: out,
			Steps: []job.Step{
				job.Exec("ffmpeg", ffmpeg.EncodeArgs(file, out, target, crf, preset)...),
			},
			Note: fmt.Sprintf("%s/%s crf=%d preset=%s", vcodec, acodec, crf, preset),
		}, nil
	})
}
