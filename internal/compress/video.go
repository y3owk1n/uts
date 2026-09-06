//nolint:goconst
package compress

import (
	"fmt"

	derrors "github.com/y3owk1n/uts/internal/core/errors"
	"github.com/y3owk1n/uts/internal/ffmpeg"
	"github.com/y3owk1n/uts/internal/format"
	"github.com/y3owk1n/uts/internal/job"
	"github.com/y3owk1n/uts/internal/ui"
	"github.com/y3owk1n/uts/internal/util"
)

// VideoOptions represents options for video compression.
type VideoOptions struct {
	Files     []string
	Quality   string
	OutputDir string
	InPlace   bool
	DryRun    bool
}

// Video compresses video files using ffmpeg, keeping the input container.
func Video(opts VideoOptions) error {
	crf, preset, err := util.VideoQuality(opts.Quality)
	if err != nil {
		return err
	}

	err = ffmpeg.Check()
	if err != nil {
		return err
	}

	ui.Message.Infof(
		"Video compression at %s quality (crf=%d, preset=%s)",
		opts.Quality,
		crf,
		preset,
	)

	return job.Run(opts.Files, job.Options{
		Verb:    "Compressing",
		Done:    "Compressed",
		Noun:    "video file",
		Compare: true,
		InPlace: opts.InPlace,
		DryRun:  opts.DryRun,
		Code:    derrors.CodeCompressionFailed,
	}, func(file string) (*job.Job, error) {
		ext := format.Ext(file)
		if format.Classify(ext) != format.Video {
			return nil, derrors.Newf(
				derrors.CodeUnsupportedFormat,
				"unsupported video format .%s",
				ext,
			)
		}

		out := util.CalcOutputPath(file, "small", opts.OutputDir)
		vcodec, acodec := format.VideoCodecs(ext)

		return &job.Job{
			Input:  file,
			Output: out,
			Steps: []job.Step{
				job.Exec("ffmpeg", ffmpeg.EncodeArgs(file, out, ext, crf, preset)...),
			},
			Note: fmt.Sprintf("%s/%s crf=%d", vcodec, acodec, crf),
		}, nil
	})
}
