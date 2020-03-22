package gotube

import (
	"bytes"
	"testing"

	"github.com/fabiodcorreia/gotube/internal/info"
)

func TestVideo_Download(t *testing.T) {
	tests := []struct {
		name    string
		v       *Video
		want    int64
		wantW   string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			got, err := tt.v.Download(w)
			if (err != nil) != tt.wantErr {
				t.Errorf("Video.Download() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Video.Download() = %v, want %v", got, tt.want)
			}
			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf("Video.Download() = %v, want %v", gotW, tt.wantW)
			}
		})
	}
}

func TestGetVideoDetails(t *testing.T) {
	type args struct {
		url string
	}
	tests := []struct {
		name      string
		args      args
		wantVideo Video
		wantErr   bool
	}{
		{
			name: "Get Video Details",
			args: args{
				url: "https://www.youtube.com/watch?v=urarTyKn9cg2",
			},
			wantErr: false,
			wantVideo: Video{
				VideoInfo: info.VideoInfo{
					ID:       "urarTyKn9cg",
					Duration: 1533,
					Title:    "GopherCon UK 2019: Julie Qiu - Finding Dependable Go Packages",
					Streams: []info.Stream{
						{
							ContentLength: 1,
							Extension:     "",
							MimeType:      "",
							Quality:       "",
						},
						{
							ContentLength: 1,
							Extension:     "",
							MimeType:      "",
							Quality:       "",
						},
					},
				},
			},
		},
		{
			name: "Invalid Youtube get_video_info url",
			args: args{
				url: "https://google.com",
			},
			wantErr: true,
		},
		{
			name: "Invalid url",
			args: args{
				url: "https://127.0.0.1",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVideo, err := GetVideoDetails(tt.args.url)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetVideoDetails() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotVideo.Duration != tt.wantVideo.Duration {
				t.Errorf("GetVideoDetails() Duration = %v, want %v", gotVideo, tt.wantVideo)
				return
			}
			if gotVideo.ID != tt.wantVideo.ID {
				t.Errorf("GetVideoDetails() ID = %v, want %v", gotVideo, tt.wantVideo)
				return
			}
			if gotVideo.Title != tt.wantVideo.Title {
				t.Errorf("GetVideoDetails() Title = %v, want %v", gotVideo, tt.wantVideo)
				return
			}
			if len(gotVideo.Streams) != len(tt.wantVideo.Streams) {
				t.Errorf("GetVideoDetails() Streams = %v, want %v", len(gotVideo.Streams), len(tt.wantVideo.Streams))
				return
			}
		})
	}
}
