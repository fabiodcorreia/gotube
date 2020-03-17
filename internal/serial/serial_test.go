package serial

import (
	"fmt"
	"io"
	"log"
	"os"
	"reflect"
	"strings"
	"testing"
)

type readerWithError struct{}

func (r readerWithError) Read(p []byte) (n int, err error) {
	return 0, fmt.Errorf("Fake io.Reader Error")
}

func getReader(path string) io.Reader {
	r, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	return r
}

func expUnmarshalSuccess() PlayerResponse {
	return PlayerResponse{
		Playability: PlayabilityStatus{
			Status: playabilityStatusOK,
		},
		Details: VideoDetails{
			Duration: "944",
			Title:    "Building+Hexagonal+Microservices+with+Go+-+Part+Three",
			VideoID:  "QyBXz9SpPqE",
		},
		Data: StreamingData{
			Formats: []Format{
				{
					Itag:          18,
					ContentLength: "32870018",
					MimeType:      "video/mp4;+codecs=\"avc1.42001E,+mp4a.40.2\"",
					Quality:       "360p",
					URL:           "https://r2---sn-8vq54vox2u-apns.googlevideo.com/videoplayback?expire=1584034585&ei=uR5qXqa8A_H4sAKX0rmgDQ&ip=77.54.222.154&id=o-ADs5ZADu6k96QVrzg8zyxlePsZRFp-9VuVG64kYTFu3b&itag=18&source=youtube&requiressl=yes&mh=2r&mm=31%2C26&mn=sn-8vq54vox2u-apns%2Csn-h5q7knez&ms=au%2Conr&mv=m&mvi=1&pl=19&initcwndbps=1225000&vprv=1&mime=video%2Fmp4&gir=yes&clen=32870018&ratebypass=yes&dur=943.542&lmt=1572303210674293&mt=1584012907&fvip=5&fexp=23842630%2C23882514&c=WEB&txp=2311222&sparams=expire%2Cei%2Cip%2Cid%2Citag%2Csource%2Crequiressl%2Cvprv%2Cmime%2Cgir%2Cclen%2Cratebypass%2Cdur%2Clmt&sig=ADKhkGMwRQIgYU25-xKMXgcY1v7rXvFbNruDmipyUmTb8Z_eiJA99lwCIQCm9LFQvHcORjsXpyp_WlCsggh7fsIF_wYYowy8FaOTcw%3D%3D&lsparams=mh%2Cmm%2Cmn%2Cms%2Cmv%2Cmvi%2Cpl%2Cinitcwndbps&lsig=ABSNjpQwRQIgc_ZGZn819YWW6Tgkx4NAe3PdNgMHEvXgKxeIq8OgtOECIQCBlygi3Z9ROfJtCKwlrbHRNevkv7hwB3bSwEc_lvvYVg%3D%3D",
				},
				{
					Itag:          22,
					ContentLength: "",
					MimeType:      "video/mp4;+codecs=\"avc1.64001F,+mp4a.40.2\"",
					Quality:       "720p",
					URL:           "https://r2---sn-8vq54vox2u-apns.googlevideo.com/videoplayback?expire=1584034585&ei=uR5qXqa8A_H4sAKX0rmgDQ&ip=77.54.222.154&id=o-ADs5ZADu6k96QVrzg8zyxlePsZRFp-9VuVG64kYTFu3b&itag=22&source=youtube&requiressl=yes&mh=2r&mm=31%2C26&mn=sn-8vq54vox2u-apns%2Csn-h5q7knez&ms=au%2Conr&mv=m&mvi=1&pl=19&initcwndbps=1225000&vprv=1&mime=video%2Fmp4&ratebypass=yes&dur=943.542&lmt=1572303329020022&mt=1584012907&fvip=5&fexp=23842630%2C23882514&c=WEB&txp=2316222&sparams=expire%2Cei%2Cip%2Cid%2Citag%2Csource%2Crequiressl%2Cvprv%2Cmime%2Cratebypass%2Cdur%2Clmt&sig=ADKhkGMwRQIgMR-Ifi6f9n5OcnI9ZPYyxRpLjjecuCAkq6z7NQAa-okCIQCA0kLJGMrQ3ZeYdj4GQwyI0PMc67br_8nxXHu_6ugvrg%3D%3D&lsparams=mh%2Cmm%2Cmn%2Cms%2Cmv%2Cmvi%2Cpl%2Cinitcwndbps&lsig=ABSNjpQwRQIgc_ZGZn819YWW6Tgkx4NAe3PdNgMHEvXgKxeIq8OgtOECIQCBlygi3Z9ROfJtCKwlrbHRNevkv7hwB3bSwEc_lvvYVg%3D%3D",
				},
			},
		},
	}
}

func TestGetPlayerResponse(t *testing.T) {
	type args struct {
		jsonStr io.Reader
	}
	tests := []struct {
		name    string
		args    args
		want    PlayerResponse
		wantErr bool
	}{
		struct {
			name    string
			args    args
			want    PlayerResponse
			wantErr bool
		}{
			name: "Unmarshal with success",
			args: args{
				jsonStr: getReader("../../testdata/player_response_ok.json"),
			},
			want:    expUnmarshalSuccess(),
			wantErr: false,
		},
		{
			name: "Unmarshal with status not ok",
			args: args{
				jsonStr: getReader("../../testdata/player_response_status_nok.json"),
			},
			wantErr: true,
		},
		{
			name: "Unmarshal with nil input",
			args: args{
				jsonStr: nil,
			},
			wantErr: true,
		},
		{
			name: "Unmarshal with error on Read",
			args: args{
				jsonStr: readerWithError{},
			},
			wantErr: true,
		},
		{
			name: "Unmarshal with invalid json",
			args: args{
				jsonStr: strings.NewReader("invalid json"),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := PlayerResponse{}
			err := GetPlayerResponse(tt.args.jsonStr, &got)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetPlayerResponse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetPlayerResponse() = %v, want %v", got, tt.want)
			}
		})
	}
}

// Benchmarks

func BenchmarkGetPlayerResponse(b *testing.B) {
	r := getReader("../../testdata/player_response_ok.json")
	b.SetBytes(2)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		rp := PlayerResponse{}
		GetPlayerResponse(r, &rp)
	}
}
