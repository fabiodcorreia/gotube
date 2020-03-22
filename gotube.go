package gotube

import (
	"fmt"
	"io"

	"github.com/fabiodcorreia/gotube/internal/client"
	"github.com/fabiodcorreia/gotube/internal/info"
	"github.com/fabiodcorreia/gotube/internal/parse"
	"github.com/fabiodcorreia/gotube/internal/serial"
)

const videoInfoDetailsURL = "https://youtube.com/get_video_info?video_id=%s"

// Video holds the video details
type Video struct {
	client *client.CustomClient
	info.VideoInfo
}

func (v *Video) Download(w io.Writer) (int64, error) {
	resp, err := v.client.Get(v.Streams[0].URL)
	if err != nil {
		return 0, fmt.Errorf("DownloadVideo: %w", err)
	}
	defer resp.Close()
	//TODO try https://golang.org/pkg/io/#CopyBuffer with bigger or smaller buffer
	return io.Copy(w, resp)
}

var c = client.NewClient()

func GetVideoDetails(url string) (video Video, err error) {

	videoID, err := parse.FindVideoID(url)
	if err != nil {
		return video, fmt.Errorf("gotube.GetVideoDetails.FindVideID: %w", err)
	}

	inf, err := fetchVideoDetails(fmt.Sprintf(videoInfoDetailsURL, videoID), &c)
	if err != nil {
		return video, fmt.Errorf("gotube.GetVideoDetails.GetVideoInfo: %w", err)
	}

	return Video{
		VideoInfo: inf,
		client:    &c,
	}, nil
}

func fetchVideoDetails(url string, c *client.CustomClient) (inf info.VideoInfo, err error) {

	resp, err := c.Get(url)
	if err != nil {
		return inf, fmt.Errorf("gotube.fetchVideoDetails.Get: %w", err)
	}
	defer resp.Close()

	json, err := parse.ExtractPlayerResp(resp)
	if err != nil {
		return inf, fmt.Errorf("gotube.GetVideoDetails.ExtractPlayerResponse: %w", err)
	}

	var player serial.PlayerResponse
	err = serial.GetPlayerResponse(json, &player)
	if err != nil {
		return inf, fmt.Errorf("gotube.GetVideoDetails.GetPlayerResponse: %w", err)
	}

	return info.GetVideoInfo(player)
}
