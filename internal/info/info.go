package info

import (
	"fmt"
	"mime"
	"strconv"

	"github.com/fabiodcorreia/gotube/internal/serial"
)

// VideoInfo holds the video information details and streams to download
type VideoInfo struct {
	ID       string
	Title    string
	Duration int
	Streams  []Stream
}

// VideoQuality represents the video quality of the stream, 360p or 720p
type VideoQuality string

// Stream holds the video stream sources and content details
type Stream struct {
	URL           string
	MimeType      string
	Extension     string
	ContentLength string
	Quality       VideoQuality
}

// GetVideoInfo receives a PlayerResponse and return a VideoInfo for that
func GetVideoInfo(player serial.PlayerResponse) (videoInfo VideoInfo, err error) {
	d, err := strconv.Atoi(player.Details.Duration)
	if err != nil {
		return videoInfo, fmt.Errorf("info.GetVideoInfo.Duration: %w", err)
	}

	st := make([]Stream, len(player.Data.Formats))
	err = streamFromFormat(player.Data.Formats, st)
	if err != nil {
		return videoInfo, err
	}

	return VideoInfo{
		ID:       player.Details.VideoID,
		Title:    player.Details.Title, //strings.ReplaceAll(player.Details.Title, "+", " "),
		Duration: d,
		Streams:  st,
	}, nil
}

func streamFromFormat(formats []serial.Format, streams []Stream) error {

	if formats == nil {
		return fmt.Errorf("info.streamFromFormat.formats: slice is nil")
	}
	if streams == nil {
		return fmt.Errorf("info.streamFromFormat.streams: slice is nil")
	}
	for i := 0; i < len(formats); i++ {
		ext, err := extensionFromType(formats[i].MimeType)
		if err != nil {
			return err
		}
		streams[i] = Stream{
			URL:           formats[i].URL,
			ContentLength: formats[i].ContentLength,
			MimeType:      formats[i].MimeType,
			Quality:       VideoQuality(formats[i].Quality),
			Extension:     ext,
		}
	}
	return nil
}

//TODO this mime.ExtensionsByType is very heavy
func extensionFromType(mimeType string) (ext string, err error) {
	exts, err := mime.ExtensionsByType(mimeType)
	if err != nil {
		return ext, fmt.Errorf("info.extensionFromType.ExtensionByType: %w", err)
	}
	if len(exts) == 0 {
		return ext, fmt.Errorf("info.extensionFromType.NoMatch: Extension not found for type %s", mimeType)
	}
	return exts[0], nil
}
