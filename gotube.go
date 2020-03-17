package gotube

import (
	"fmt"
	"io"
	"net/http"

	"github.com/fabiodcorreia/gotube/internal/client"
	"github.com/fabiodcorreia/gotube/internal/info"
	"github.com/fabiodcorreia/gotube/internal/parse"
	"github.com/fabiodcorreia/gotube/internal/serial"
)

type Video struct {
	client *client.Client
	info.VideoInfo
}

func (v Video) DownloadDefault(w io.Writer) (int64, error) {
	resp, err := http.DefaultClient.Get(v.Streams[0].URL)
	if err != nil {
		return 0, fmt.Errorf("DownloadVideo: %w", err)
	}

	defer resp.Body.Close()
	return io.Copy(w, resp.Body)
}

func (v Video) Download(w io.Writer) (int64, error) {
	resp, err := v.client.Get(v.Streams[1].URL)
	if err != nil {
		return 0, fmt.Errorf("DownloadVideo: %w", err)
	}
	defer resp.Close()
	return io.Copy(w, resp)
}

var c = client.NewClient()

func GetVideoDetails(url string) (video Video, err error) {
	videoID, err := parse.FindVideoID(url)

	resp, err := http.DefaultClient.Get(fmt.Sprintf("https://youtube.com/get_video_info?video_id=%s", videoID))

	if err != nil {
		return video, fmt.Errorf("gotube.GetVideoDetails.Get: %w", err)
	}

	json, err := parse.ExtractPlayerResp(resp.Body)
	resp.Body.Close()
	if err != nil {
		return video, fmt.Errorf("gotube.GetVideoDetails.ExtractPlayerResponse: %w", err)
	}

	//var player *serial.PlayerResponse
	player := serial.PlayerResponse{}
	err = serial.GetPlayerResponse(json, &player)
	if err != nil {
		return video, fmt.Errorf("gotube.GetVideoDetails.GetPlayerResponse: %w", err)
	}

	inf, err := info.GetVideoInfo(player)
	if err != nil {
		return video, fmt.Errorf("gotube.GetVideoDetails.GetVideoInfo: %w", err)
	}

	return Video{
		VideoInfo: inf,
		client:    &c,
	}, nil
}
