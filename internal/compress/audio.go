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

// AudioOptions represents options for audio compression.
type AudioOptions struct {
	Files     []string
	Quality   string
	OutputDir string
	InPlace   bool
	DryRun    bool
}

// Audio compresses audio files using ffmpeg. Lossy inputs keep their format;
// lossless inputs (wav, flac) become AAC in an .m4a container.
func Audio(opts AudioOptions) error {
	bitrate, err := util.AudioBitrate(opts.Quality)
	if err != nil {
		return err
	}

	err = ffmpeg.Check()
	if err != nil {
		return err
	}

	ui.Message.Infof("Audio compression at %s quality (bitrate=%s)", opts.Quality, bitrate)

	return job.Run(opts.Files, job.Options{
		Verb:    "Compressing",
		Done:    "Compressed",
		Noun:    "audio files",
		Compare: true,
		InPlace: opts.InPlace,
		DryRun:  opts.DryRun,
		Code:    derrors.CodeCompressionFailed,
	}, func(file string) (*job.Job, error) {
		ext := format.Ext(file)
		if format.Classify(ext) != format.Audio {
			return nil, derrors.Newf(
				derrors.CodeUnsupportedFormat,
				"unsupported audio format .%s",
				ext,
			)
		}

		codec, outExt := format.AudioCompressTarget(ext)
		out := util.CalcOutputPathExt(file, "small", outExt, opts.OutputDir)

		return &job.Job{
			Input:  file,
			Output: out,
			Steps: []job.Step{job.Exec("ffmpeg",
				"-i", file, "-vn", "-c:a", codec, "-b:a", bitrate, "-y", out)},
			Note: fmt.Sprintf("%s %s → .%s", codec, bitrate, outExt),
		}, nil
	})
}
