package multimedia

import (
	"context"
	"testing"

	"image/color"

	"github.com/disintegration/imaging"
	"github.com/stretchr/testify/assert"
	"github.com/sweet-go/stdlib/helper"
)

func TestMultimedia_SliceImage(t *testing.T) {
	t.Run("input not found", func(t *testing.T) {
		input := &SliceImageInput{
			SourcePath:   "testdata/notfound.jpg",
			OutputFormat: imaging.JPEG,
		}

		err := SliceImage(input)
		assert.Error(t, err)
	})

	t.Run("unknown output format", func(t *testing.T) {
		input := &SliceImageInput{
			SourcePath:   "testdata/overly_high.jpg",
			OutputFormat: imaging.Format(1000),
		}

		err := SliceImage(input)
		assert.Error(t, err)
	})

	t.Run("ok", func(t *testing.T) {
		input := &SliceImageInput{
			SourcePath:   "testdata/overly_high.jpg",
			OutputDir:    "testdata/",
			MaxHeight:    1920,
			MinHeight:    1000,
			AspectRatio:  16.0 / 9.0,
			OutputFormat: imaging.JPEG,
		}

		err := SliceImage(input)
		assert.NoError(t, err)

		assert.True(t, true, len(input.OutputFiles) > 0)

		for _, file := range input.OutputFiles {
			assert.FileExists(t, file)
			assert.NoError(t, helper.DeleteFile(file))
		}
	})
}

func TestMultimedia_ScaleDownImageByWidth(t *testing.T) {
	t.Run("file not found", func(t *testing.T) {
		input := &ScaleDownImageByWidthInput{
			SourcePath: "testdata/notfound.jpg",
		}

		err := ScaleDownImageByWidth(input)
		assert.Error(t, err)
	})

	t.Run("desired width is larger than original width", func(t *testing.T) {
		input := &ScaleDownImageByWidthInput{
			SourcePath: "testdata/overly_high.jpg",
			Width:      5000,
		}

		err := ScaleDownImageByWidth(input)
		assert.NoError(t, err)
	})

	t.Run("err unsupported format", func(t *testing.T) {
		input := &ScaleDownImageByWidthInput{
			SourcePath: "testdata/overly_high.jpg",
			OutputPath: "testdata/unknown_format.uknown_format_must_begone",
			Width:      100,
			Filter:     imaging.NearestNeighbor,
		}

		err := ScaleDownImageByWidth(input)
		assert.Error(t, err)
	})

	t.Run("ok", func(t *testing.T) {
		input := &ScaleDownImageByWidthInput{
			SourcePath: "testdata/overly_width.png",
			OutputPath: "testdata/test_output_ScaleDownImageByWidth_1080.jpg",
			Width:      1080,
			Filter:     imaging.NearestNeighbor,
		}

		err := ScaleDownImageByWidth(input)
		assert.NoError(t, err)

		assert.FileExists(t, input.OutputPath)

		assert.NoError(t, helper.DeleteFile(input.OutputPath))
	})
}

func TestMultimedia_ConvertImage(t *testing.T) {
	t.Run("file not found", func(t *testing.T) {
		input := &ConvertImageInput{
			SourcePath: "testdata/notfound_lah.jpg",
			OutputPath: "testdata/unknown_format.uknown_format_must_begone",
		}

		err := ConvertImage(input)
		assert.Error(t, err)
	})

	t.Run("ok", func(t *testing.T) {
		input := &ConvertImageInput{
			SourcePath: "testdata/overly_high.jpg",
			OutputPath: "testdata/test_output_ConvertImage.png",
		}

		err := ConvertImage(input)
		assert.NoError(t, err)

		assert.FileExists(t, input.OutputPath)

		assert.NoError(t, helper.DeleteFile(input.OutputPath))
	})
}

func TestMultimedia_MergeImagesToVideos(t *testing.T) {
	ctx := context.Background()
	t.Run("err file not found", func(t *testing.T) {
		input := &MergeImagesToVideosInput{
			ImageDurations: map[string]float64{
				"testdata/notfound.jpg": 3.0,
			},
			OutputPath: "testdata/test_output_MergeImagesToVideos.mp4",
		}

		err := MergeImagesToVideos(ctx, input)
		assert.Error(t, err)
	})

	t.Run("err unknown format", func(t *testing.T) {
		input := &MergeImagesToVideosInput{
			ImageDurations: map[string]float64{
				"testdata/overly_high.jpg": 3.0,
			},
			OutputPath: "testdata/test_output_MergeImagesToVideos.unknown_format",
		}

		err := MergeImagesToVideos(ctx, input)
		assert.Error(t, err)
	})

	t.Run("ok", func(t *testing.T) {
		input := &MergeImagesToVideosInput{
			ImageDurations: map[string]float64{
				"testdata/overly_high.jpg": 3.0,
			},
			OutputPath: "testdata/test_output_MergeImagesToVideos.mp4",
		}

		err := MergeImagesToVideos(ctx, input)
		assert.NoError(t, err)

		assert.FileExists(t, input.OutputPath)

		data, err := GetVideoData(ctx, input.OutputPath)
		assert.NoError(t, err)

		assert.Equal(t, "3.000000", data.FirstVideoStream().Duration)

		assert.NoError(t, helper.DeleteFile(input.OutputPath))
	})
}

func TestMultimedia_ScaleUpAndFillImage(t *testing.T) {
	ctx := context.Background()
	t.Run("file not found", func(t *testing.T) {
		input := &ScaleUpAndFillImageInput{
			SourcePath: "testdata/notfound.jpg",
			OutputPath: "testdata/unknown_format.uknown_format_must_begone",
			Width:      100,
			Height:     100,
			Color:      color.Transparent,
		}

		err := ScaleUpAndFillImage(ctx, input)
		assert.Error(t, err)
	})

	t.Run("not error even if the supplied format OutputPath is unknown", func(t *testing.T) {
		input := &ScaleUpAndFillImageInput{
			SourcePath: "testdata/overly_high.jpg",
			OutputPath: "testdata/unknown_format.uknown_format_must_begone",
			Width:      100,
			Height:     100,
			Color:      color.Transparent,
		}

		err := ScaleUpAndFillImage(ctx, input)
		assert.NoError(t, err)
		assert.FileExists(t, input.OutputPath)

		assert.NoError(t, helper.DeleteFile(input.OutputPath))
	})

	t.Run("ok", func(t *testing.T) {
		input := &ScaleUpAndFillImageInput{
			SourcePath: "testdata/overly_high.jpg",
			OutputPath: "testdata/test_output_ScaleUpAndFillImage.jpg",
			Width:      100,
			Height:     100,
			Color:      color.Transparent,
		}

		err := ScaleUpAndFillImage(ctx, input)
		assert.NoError(t, err)

		assert.FileExists(t, input.OutputPath)

		assert.NoError(t, helper.DeleteFile(input.OutputPath))
	})
}

func TestMultimedia_ScaleUpImageByResolution(t *testing.T) {
	ctx := context.Background()
	t.Run("file not found", func(t *testing.T) {
		input := &ScaleUpImageByResolutionInput{
			SourcePath: "testdata/notfound.jpg",
			OutputPath: "testdata/unknown_format.uknown_format_must_begone",
			MaxWidth:   100,
			MaxHeight:  100,
			Filter:     imaging.NearestNeighbor,
		}

		err := ScaleUpImageByResolution(ctx, input)
		assert.Error(t, err)
	})

	t.Run("err format output", func(t *testing.T) {
		input := &ScaleUpImageByResolutionInput{
			SourcePath: "testdata/overly_high.jpg",
			OutputPath: "testdata/unknown_format.uknown_format_must_begone",
			MaxWidth:   100,
			MaxHeight:  100,
			Filter:     imaging.NearestNeighbor,
		}

		err := ScaleUpImageByResolution(ctx, input)
		assert.Error(t, err)
	})

	t.Run("ok", func(t *testing.T) {
		input := &ScaleUpImageByResolutionInput{
			SourcePath: "testdata/overly_high.jpg",
			OutputPath: "testdata/test_output_ScaleUpImageByResolution.jpg",
			MaxWidth:   100,
			MaxHeight:  100,
			Filter:     imaging.NearestNeighbor,
		}

		err := ScaleUpImageByResolution(ctx, input)
		assert.NoError(t, err)

		assert.FileExists(t, input.OutputPath)

		assert.NoError(t, helper.DeleteFile(input.OutputPath))
	})
}
