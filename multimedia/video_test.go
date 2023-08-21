package multimedia

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/sweet-go/stdlib/helper"
)

func TestGetVideoData(t *testing.T) {
	ctx := context.Background()
	videoPath := "./testdata/common_testing_video.mp4"

	t.Run("ok", func(t *testing.T) {
		data, err := GetVideoData(ctx, videoPath)
		assert.NoError(t, err)

		assert.Equal(t, "QuickTime / MOV", data.Format.FormatLongName)
	})

	t.Run("video not found", func(t *testing.T) {
		_, err := GetVideoData(ctx, "not_found.mp4")
		assert.Error(t, err)
	})
}

func TestIsVideoHasAudio(t *testing.T) {
	ctx := context.Background()
	videoWithAudio := "./testdata/common_testing_video.mp4"
	videoWithoutAudio := "./testdata/video_without_audio.mp4"

	t.Run("yes", func(t *testing.T) {
		hasAudio, err := IsVideoHasAudio(ctx, videoWithAudio)
		assert.NoError(t, err)

		assert.True(t, hasAudio)
	})

	t.Run("no", func(t *testing.T) {
		hasAudio, err := IsVideoHasAudio(ctx, videoWithoutAudio)
		assert.NoError(t, err)

		assert.False(t, hasAudio)
	})

	t.Run("video not found", func(t *testing.T) {
		_, err := IsVideoHasAudio(ctx, "not_found.mp4")
		assert.Error(t, err)
	})
}

func TestScaleVideo(t *testing.T) {
	ctx := context.Background()
	videoPath := "./testdata/common_testing_video.mp4"
	outputPath := "./testdata/test_output_scale.mp4"

	t.Run("ok", func(t *testing.T) {
		err := ScaleVideo(ctx, &ScaleVideoInput{
			SourcePath: videoPath,
			OutputPath: outputPath,
			ScaleRatio: "100:200",
		})
		assert.NoError(t, err)

		defer func() {
			err := helper.DeleteFile(outputPath)
			assert.NoError(t, err)
		}()

		data, err := GetVideoData(ctx, outputPath)
		assert.NoError(t, err)

		assert.Equal(t, 100, data.FirstVideoStream().Width)
		assert.Equal(t, 200, data.FirstVideoStream().Height)
	})

	t.Run("video not found", func(t *testing.T) {
		err := ScaleVideo(ctx, &ScaleVideoInput{
			SourcePath: "not_found.mp4",
			OutputPath: outputPath,
			ScaleRatio: "1280:720",
		})
		assert.Error(t, err)
	})
}

func TestConvertWebmToMP4(t *testing.T) {
	input := &WebmToMP4Input{
		SourcePath: "./testdata/webm_video_input.webm",
		OutputPath: "./testdata/test_output+ConvertWebmToMP4.mp4",
	}

	t.Run("ok", func(t *testing.T) {
		err := ConvertWebmToMP4(input)
		assert.NoError(t, err)

		defer func() {
			err := helper.DeleteFile(input.OutputPath)
			assert.NoError(t, err)
		}()
	})

	t.Run("err", func(t *testing.T) {
		inputNotFound := &WebmToMP4Input{
			SourcePath: "./testdata/notfound.webm",
			OutputPath: "./testdata/test_output+ConvertWebmToMP4.mp4",
		}

		err := ConvertWebmToMP4(inputNotFound)
		assert.Error(t, err)
	})
}

func TestConcatMP4Videos(t *testing.T) {
	ctx := context.Background()

	input := &ConcatMP4VideosInput{
		SourcePaths: []string{
			"./testdata/common_testing_video.mp4",
			"./testdata/video_without_audio.mp4",
		}, OutputPath: "./testdata/test_output_ConcatMP4Videos.mp4",
		ListFile:  fmt.Sprintf("./%s.txt", helper.GenerateID()),
		ErrStream: os.Stderr,
	}

	t.Run("ok", func(t *testing.T) {
		err := ConcatMP4Videos(ctx, input)
		assert.NoError(t, err)

		defer func() {
			err := helper.DeleteFile(input.OutputPath)
			assert.NoError(t, err)
		}()
	})

	t.Run("forbidden list file path", func(t *testing.T) {
		forbiddenInput := &ConcatMP4VideosInput{
			SourcePaths: []string{
				"./testdata/common_testing_video.mp4",
				"./testdata/video_without_audio.mp4",
			}, OutputPath: "./testdata/test_output_ConcatMP4Videos.mp4",
			ListFile:  fmt.Sprintf("/root/%s.txt", helper.GenerateID()),
			ErrStream: os.Stderr,
		}

		err := ConcatMP4Videos(ctx, forbiddenInput)
		assert.Error(t, err)
	})

	t.Run("video not found", func(t *testing.T) {
		inputNotFound := &ConcatMP4VideosInput{
			SourcePaths: []string{
				"./testdata/notfound2.mp4",
				"./testdata/notfound.mp4",
			}, OutputPath: "./testdata/test_output_ConcatMP4Videos.mp4",
			ListFile: fmt.Sprintf("./%s.txt", helper.GenerateID()),
		}

		err := ConcatMP4Videos(ctx, inputNotFound)
		assert.Error(t, err)
	})
}

func TestMultimedia_GetVideoAspectRatio(t *testing.T) {
	ctx := context.Background()
	t.Run("err", func(t *testing.T) {
		_, _, err := GetVideoAspectRatio(ctx, "testdata/notfound.mp4")
		assert.Error(t, err)
	})

	t.Run("ok", func(t *testing.T) {
		w, h, err := GetVideoAspectRatio(ctx, "testdata/common_testing_video.mp4")
		assert.NoError(t, err)

		assert.Equal(t, 460, w)
		assert.Equal(t, 252, h)
	})
}

func TestMultimedia_DetermineVideoOrientation(t *testing.T) {
	ctx := context.Background()
	t.Run("landscape", func(t *testing.T) {
		orientation, err := DetermineVideoOrientation(ctx, "testdata/common_testing_video.mp4")
		assert.NoError(t, err)

		assert.Equal(t, VideoOrientationLandscape, orientation)
	})

	t.Run("portrait", func(t *testing.T) {
		orientation, err := DetermineVideoOrientation(ctx, "testdata/portrait_video.mp4")
		assert.NoError(t, err)

		assert.Equal(t, VideoOrientationPortrait, orientation)
	})

	t.Run("square", func(t *testing.T) {
		orientation, err := DetermineVideoOrientation(ctx, "testdata/square_video.mp4")
		assert.NoError(t, err)

		assert.Equal(t, VideoOrientationSquare, orientation)
	})

	t.Run("video not found", func(t *testing.T) {
		_, err := DetermineVideoOrientation(ctx, "testdata/notfound.mp4")
		assert.Error(t, err)
	})
}

func TestMultimedia_TransformLandscapeVideoToPortrait(t *testing.T) {
	t.Run("file not found", func(t *testing.T) {
		err := TransformLandscapeVideoToPortrait(context.Background(), &TransformLandscapeVideoToPortraitInput{
			SourcePath: "testdata/notfound.mp4",
			OutputPath: "testdata/test_output_TransformLandscapeVideoToPortrait.mp4",
		})
		assert.Error(t, err)
	})

	t.Run("video already potrait", func(t *testing.T) {
		err := TransformLandscapeVideoToPortrait(context.Background(), &TransformLandscapeVideoToPortraitInput{
			SourcePath: "testdata/portrait_video.mp4",
			OutputPath: "testdata/test_output_TransformLandscapeVideoToPortrait.mp4",
		})
		assert.NoError(t, err)
	})

	t.Run("err", func(t *testing.T) {
		err := TransformLandscapeVideoToPortrait(context.Background(), &TransformLandscapeVideoToPortraitInput{
			SourcePath: "testdata/common_testing_video.mp4",

			// invalid width and height
			Width:      -100,
			Height:     -100,
			OutputPath: "testdata/test_output_TransformLandscapeVideoToPortrait.mp4",
		})
		assert.Error(t, err)

		assert.NoError(t, helper.DeleteFile("testdata/test_output_TransformLandscapeVideoToPortrait.mp4"))
	})
}

func TestMultimedia_ResizeAndEncodeVideo(t *testing.T) {
	t.Run("error", func(t *testing.T) {
		err := ResizeAndEncodeVideo(context.Background(), &ResizeAndEncodeVideoInput{
			SourcePath: "testdata/notfound.mp4",
			OutputPath: "testdata/test_output_ResizeAndEncodeVideo.mp4",
		})
		assert.Error(t, err)
	})

	t.Run("ok", func(t *testing.T) {
		err := ResizeAndEncodeVideo(context.Background(), &ResizeAndEncodeVideoInput{
			SourcePath: "testdata/common_testing_video.mp4",
			OutputPath: "testdata/test_output_ResizeAndEncodeVideo.mp4",
			ErrStream:  os.Stderr,
			OutStream:  os.Stdout,
			Width:      100,
			Tune:       "film",
			Preset:     "ultrafast",
		})
		assert.NoError(t, err)

		assert.NoError(t, helper.DeleteFile("testdata/test_output_ResizeAndEncodeVideo.mp4"))
	})
}
