package info

import (
	"fmt"
	"regexp"
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

// VideoExt represent the file extension of the video based on the mime type
type VideoExt string

// Stream holds the video stream sources and content details
type Stream struct {
	URL           string
	MimeType      string
	Extension     VideoExt
	ContentLength int
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
	var cl int
	var err error
	for i := 0; i < len(formats); i++ {
		if formats[i].ContentLength != "" {
			cl, err = strconv.Atoi(formats[i].ContentLength)
			if err != nil {
				return fmt.Errorf("info.streamFromFormat.ContentLength: %w", err)
			}
		} else {
			cl = 0
		}

		streams[i] = Stream{
			URL:           formats[i].URL,
			ContentLength: cl,
			MimeType:      formats[i].MimeType,
			Quality:       VideoQuality(formats[i].Quality),
			Extension:     extensionFromType(formats[i].MimeType),
		}
	}
	return nil
}

const (
	MP4  VideoExt = ".mp4"
	TGP           = ".3gp"
	FLV           = ".flv"
	WEBM          = ".webm"
	AVI           = ".avi"
)

const videoTypePattern = "video\\/(\\w+)"

var videoTypeRegexp = regexp.MustCompile(videoTypePattern)

// mime.ExtensionByType is very heavy compared with regex, benchmark on info_test
func extensionFromType(mimeType string) VideoExt {
	match := videoTypeRegexp.FindStringSubmatch(mimeType)
	if match == nil || len(match) == 1 {
		return AVI
	}

	switch match[1] {
	case "mp4":
		return MP4
	case "3gp":
		return TGP
	case "flv":
		return FLV
	case "webm":
		return WEBM
	default:
		return AVI
	}
}
