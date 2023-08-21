// Package multimedia contains functions related to multimedia processing
package multimedia

import (
	"io"
	"os/exec"
)

// AddAudioToVideoInput is the input for AddAudioToVideo function
type AddAudioToVideoInput struct {
	VideoSourcePath string
	AudioSourcePath string
	OutputPath      string
	ErrStream       io.Writer
	OutStream       io.Writer
}

// AddAudioToVideo adds audio to video and will repeat the audio if the video is longer than the audio.
// If the name of the output matches with existing files, it will be overwritten.
func AddAudioToVideo(input *AddAudioToVideoInput) error {
	cmd := exec.Command(
		"ffmpeg",
		"-i", input.VideoSourcePath,
		"-stream_loop", "-1",
		"-i", input.AudioSourcePath,
		"-shortest",
		"-map", "0:v",
		"-map", "1:a",
		"-c:v", "copy",
		input.OutputPath,
		"-y",
	)

	cmd.Stderr = input.ErrStream
	cmd.Stdout = input.OutStream

	err := cmd.Run()
	if err != nil {
		return err
	}

	return nil
}
