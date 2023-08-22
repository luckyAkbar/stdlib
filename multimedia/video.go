package multimedia

import (
	"context"
	"fmt"
	"io"
	"os/exec"

	"github.com/sweet-go/stdlib/helper"
	ffmpeg "github.com/u2takey/ffmpeg-go"
	"gopkg.in/vansante/go-ffprobe.v2"
)

// GetVideoData returns the video data utilizing ffprobe
func GetVideoData(ctx context.Context, path string) (*ffprobe.ProbeData, error) {
	data, err := ffprobe.ProbeURL(ctx, path)
	if err != nil {
		return nil, err
	}

	return data, nil
}

// IsVideoHasAudio returns true if the video has audio stream
func IsVideoHasAudio(ctx context.Context, path string) (bool, error) {
	data, err := GetVideoData(ctx, path)
	if err != nil {
		return false, err
	}

	if audio := data.FirstAudioStream(); audio == nil {
		return false, nil
	}

	return true, nil
}

// ScaleVideoInput is the input for ScaleVideo function
type ScaleVideoInput struct {
	SourcePath string
	OutputPath string

	// ScaleRatio is the ratio of the video to be scaled. Example valid value: 1280:720
	ScaleRatio string
	ErrStream  io.Writer
	OutStream  io.Writer
}

// ScaleVideo scales the video using ffmpeg
func ScaleVideo(ctx context.Context, input *ScaleVideoInput) error {
	cmd := exec.CommandContext(
		ctx,
		"ffmpeg",
		"-y",
		"-i", input.SourcePath,
		"-vf", fmt.Sprintf("scale=%s", input.ScaleRatio),
		"-acodec", "aac",
		input.OutputPath,
	)

	cmd.Stderr = input.ErrStream
	cmd.Stdout = input.OutStream

	println(cmd.String())

	err := cmd.Run()
	if err != nil {
		return err
	}

	return nil
}

// WebmToMP4Input is the input for ConvertWebmToMP4 function
type WebmToMP4Input struct {
	SourcePath   string
	OutputPath   string
	InputKwargs  []ffmpeg.KwArgs
	OutputKwargs []ffmpeg.KwArgs
}

// ConvertWebmToMP4 converts webm to mp4 using ffmpeg
func ConvertWebmToMP4(input *WebmToMP4Input) error {
	return ffmpeg.Input(input.SourcePath, input.InputKwargs...).Output(input.OutputPath, input.OutputKwargs...).ErrorToStdOut().OverWriteOutput().Run()
}

// ConcatMP4VideosInput is the input for ConcatMP4Videos function
type ConcatMP4VideosInput struct {
	SourcePaths []string
	OutputPath  string

	// ListFile is the path of the txt file that contains the list of the source paths. The file generated here will be deleted once the process
	// has finished
	ListFile  string
	ErrStream io.Writer
	OutStream io.Writer
}

// ConcatMP4Videos concatenates mp4 videos using ffmpeg.
// Will overwrite the output path if it already exists.
func ConcatMP4Videos(ctx context.Context, input *ConcatMP4VideosInput) error {
	err := generateConcatFileList(input.SourcePaths, input.ListFile)
	if err != nil {
		return err
	}

	// avoid having too many unused txt file from this process
	defer func() {
		helper.LogIfError(helper.DeleteFile(input.ListFile))
	}()

	cmd := exec.CommandContext(
		ctx,
		"ffmpeg",
		"-f", "concat",
		"-safe", "0",
		"-i", input.ListFile,
		"-c", "copy",
		input.OutputPath,
		"-y",
	)

	cmd.Stderr = input.ErrStream
	cmd.Stdout = input.OutStream

	return cmd.Run()
}

// GetVideoAspectRatio returns the width and height of the video
func GetVideoAspectRatio(ctx context.Context, path string) (width, height int, err error) {
	data, err := ffprobe.ProbeURL(ctx, path)
	if err != nil {
		return 0, 0, err
	}

	return data.FirstVideoStream().Width, data.FirstVideoStream().Height, nil
}

// VideoOrientation is the orientation of the video
type VideoOrientation string

// list of video orientations
const (
	VideoOrientationLandscape VideoOrientation = "landscape"
	VideoOrientationPortrait  VideoOrientation = "portrait"
	VideoOrientationSquare    VideoOrientation = "square"
)

// DetermineVideoOrientation determines the orientation of the video
func DetermineVideoOrientation(ctx context.Context, path string) (VideoOrientation, error) {
	w, h, err := GetVideoAspectRatio(ctx, path)
	if err != nil {
		return "", err
	}

	if w == h {
		return VideoOrientationSquare, nil
	}

	if w > h {
		return VideoOrientationLandscape, nil
	}

	return VideoOrientationPortrait, nil
}

// TransformLandscapeVideoToPortraitInput is the input for TransformLandscapeVideoToPortrait function
type TransformLandscapeVideoToPortraitInput struct {
	SourcePath string
	OutputPath string
	ErrStream  io.Writer
	OutStream  io.Writer
	Width      int
	Height     int
}

// TransformLandscapeVideoToPortrait transforms landscape video to portrait video.
// If the video is already in portrait orientation, it will assign the output path with the supplied input path.
// Will overwrite the output path if it already exists.
func TransformLandscapeVideoToPortrait(ctx context.Context, input *TransformLandscapeVideoToPortraitInput) error {
	orientation, err := DetermineVideoOrientation(ctx, input.SourcePath)
	if err != nil {
		return err
	}

	// if the video is already in portrait orientation, just return the source path
	if orientation == VideoOrientationPortrait {
		input.OutputPath = input.SourcePath
		return nil
	}

	cmd := exec.CommandContext(
		ctx,
		"ffmpeg",
		"-i", input.SourcePath,
		"-vf", fmt.Sprintf("scale=%d:-1,pad=%d:%d:(ow-iw)/2:(oh-ih)/2,setsar=1", input.Width, input.Width, input.Height),
		"-c:a", "copy",
		input.OutputPath, "-y",
	)

	cmd.Stderr = input.ErrStream
	cmd.Stdout = input.OutStream

	return cmd.Run()
}

// ResizeAndEncodeVideoInput is the input for ResizeAndEncodeVideo function
type ResizeAndEncodeVideoInput struct {
	SourcePath string
	OutputPath string
	ErrStream  io.Writer
	OutStream  io.Writer
	Width      int

	// Tune is the tune parameter for ffmpeg. See here https://trac.ffmpeg.org/wiki/Encode/H.264#Tune
	Tune string

	// Preset is the preset parameter for ffmpeg. See here https://trac.ffmpeg.org/wiki/Encode/H.264#Preset
	Preset string
}

// ResizeAndEncodeVideo resizes and encodes the video using ffmpeg.
// Highly hardcoded for our use case. If in the future needs more flexibility, must add more field on the ResizeAndEncodeVideoInput.
// Will overwrite the output path if it already exists.
func ResizeAndEncodeVideo(ctx context.Context, input *ResizeAndEncodeVideoInput) error {
	cmd := exec.CommandContext(
		ctx,
		"ffmpeg",
		"-i", input.SourcePath,
		"-vf", fmt.Sprintf("scale=%d:-2,setsar=1", input.Width),
		"-c:v", "libx264", "-preset", input.Preset, "-tune", input.Tune,
		"-c:a", "aac", "-strict", "experimental",
		input.OutputPath, "-y",
	)

	cmd.Stderr = input.ErrStream
	cmd.Stdout = input.OutStream

	return cmd.Run()
}
