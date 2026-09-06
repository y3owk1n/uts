//nolint:goconst
package convert

import (
	"fmt"

	derrors "github.com/y3owk1n/uts/internal/core/errors"
	"github.com/y3owk1n/uts/internal/ffmpeg"
	"github.com/y3owk1n/uts/internal/format"
	"github.com/y3owk1n/uts/internal/job"
	"github.com/y3owk1n/uts/internal/ui"
	"github.com/y3owk1n/uts/internal/util"
)

// AudioOptions represents options for audio conversion.
type AudioOptions struct {
	Files     []string
	Target    string
	Quality   string
	OutputDir string
	InPlace   bool
	DryRun    bool
}

// Audio converts audio files, or extracts the audio track of video files, to
// the target format.
func Audio(opts AudioOptions) error {
	codec, outExt := format.AudioCodec(opts.Target)
	if codec == "" {
		return unsupportedTarget(opts.Target, format.AudioTargets)
	}

	bitrate, err := util.AudioBitrate(opts.Quality)
	if err != nil {
		return err
	}

	err = ffmpeg.Check()
	if err != nil {
		return err
	}

	if format.Lossless(outExt) {
		ui.Message.Infof("Converting audio to .%s (%s, lossless)", outExt, codec)
	} else {
		ui.Message.Infof("Converting audio to .%s (%s, %s)", outExt, codec, bitrate)
	}

	return job.Run(opts.Files, job.Options{
		Verb:    "Converting",
		Done:    "Converted",
		Noun:    "audio files",
		InPlace: opts.InPlace,
		DryRun:  opts.DryRun,
		Code:    derrors.CodeConversionFailed,
	}, func(file string) (*job.Job, error) {
		ext := format.Ext(file)
		if cat := format.Classify(ext); cat != format.Audio && cat != format.Video {
			return nil, derrors.Newf(
				derrors.CodeUnsupportedFormat,
				"unsupported audio source .%s",
				ext,
			)
		}

		if format.Same(ext, outExt) {
			return &job.Job{Input: file, Skip: "Already ." + outExt}, nil
		}

		out := util.CalcConvertOutputPath(file, outExt, opts.OutputDir)

		args := []string{"-i", file, "-vn", "-c:a", codec}
		if !format.Lossless(outExt) {
			args = append(args, "-b:a", bitrate)
		}

		args = append(args, "-y", out)

		return &job.Job{
			Input:  file,
			Output: out,
			Steps:  []job.Step{job.Exec("ffmpeg", args...)},
			Note:   fmt.Sprintf(".%s → .%s via %s", ext, outExt, codec),
		}, nil
	})
}
