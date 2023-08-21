package multimedia

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/sweet-go/stdlib/helper"
)

func TestAddAudioToVideo(t *testing.T) {
	audioFile := "./testdata/audio_sample.mp3"
	videoWithoutSound := "./testdata/video_without_audio.mp4"

	t.Run("ok", func(t *testing.T) {
		input := &AddAudioToVideoInput{
			VideoSourcePath: videoWithoutSound,
			AudioSourcePath: audioFile,
			OutputPath:      "./testdata/test_output.mp4",
		}

		err := AddAudioToVideo(input)
		assert.NoError(t, err)

		defer func() {
			err := helper.DeleteFile(input.OutputPath)
			assert.NoError(t, err)
		}()

		audio, err := IsVideoHasAudio(context.Background(), input.OutputPath)
		assert.NoError(t, err)
		assert.True(t, audio)
	})

	t.Run("video not found", func(t *testing.T) {
		input := &AddAudioToVideoInput{
			VideoSourcePath: "not_found.mp4",
			AudioSourcePath: audioFile,
			OutputPath:      "./testdata/test_output.mp4",
		}

		err := AddAudioToVideo(input)
		assert.Error(t, err)
	})
}
