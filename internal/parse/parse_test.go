package parse

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"testing"
)

type readerWithError struct{}

func (r readerWithError) Close() error {
	return nil
}

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

func TestExtractPlayerResp(t *testing.T) {
	type args struct {
		gvInfo io.Reader
	}
	tests := []struct {
		name    string
		args    args
		want    int
		wantErr bool
	}{
		struct {
			name    string
			args    args
			want    int
			wantErr bool
		}{
			name: "Parse with success",
			args: args{
				gvInfo: getReader("../../testdata/get_video_info"),
			},
			want:    72769,
			wantErr: false,
		},
		{
			name: "Parse fail with nil input",
			args: args{
				gvInfo: nil,
			},
			want:    0,
			wantErr: true,
		},
		{
			name: "Parse fail with invalid input reader",
			args: args{
				gvInfo: readerWithError{},
			},
			want:    0,
			wantErr: true,
		},
		{
			name: "Parse fail with invalid input content",
			args: args{
				gvInfo: strings.NewReader("invalid%input%content"),
			},
			want:    0,
			wantErr: true,
		},
		{
			name: "Parse fail with no status",
			args: args{
				gvInfo: strings.NewReader("reason=Invalid+parameters.&errorcode=2"),
			},
			want:    0,
			wantErr: true,
		},
		{
			name: "Parse fail with status not ok",
			args: args{
				gvInfo: getReader("../../testdata/get_video_info_nok"),
			},
			want:    0,
			wantErr: true,
		},
		{
			name: "Parse fail with status not ok without reason",
			args: args{
				gvInfo: getReader("../../testdata/get_video_info_nok no_reason"),
			},
			want:    0,
			wantErr: true,
		},
		{
			name: "Parse fail with no player_response",
			args: args{
				gvInfo: getReader("../../testdata/get_video_info_no_player"),
			},
			want:    0,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := ExtractPlayerResp(tt.args.gvInfo)
			if (err != nil) != tt.wantErr {
				t.Errorf("ExtPlayerResp() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got == nil && tt.wantErr {
				return
			}
			content, _ := ioutil.ReadAll(got)

			if len(content) != tt.want {
				t.Errorf("ExtPlayerResp() = %v, want %v", len(content), tt.want)
			}
		})
	}
}

func TestFindVideoID(t *testing.T) {
	type args struct {
		url string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		struct {
			name    string
			args    args
			want    string
			wantErr bool
		}{
			name: "Empty URL",
			args: args{
				url: "",
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "Invalid URL",
			args: args{
				url: "invalid url",
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "Valid URL",
			args: args{
				url: "https://www.youtube.com/watch?v=zzAdEt3xZ1M",
			},
			want:    "zzAdEt3xZ1M",
			wantErr: false,
		},
		{
			name: "Valid URL from playlist",
			args: args{
				url: "https://www.youtube.com/watch?v=EFJfdWzBHwE&list=PL2ntRZ1ySWBdDyspRTNBIKES1Y-P__59_&index=4",
			},
			want:    "EFJfdWzBHwE",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := FindVideoID(tt.args.url)
			if (err != nil) != tt.wantErr {
				t.Errorf("FindVideoID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("FindVideoID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func BenchmarkFindVideoID(b *testing.B) {
	b.SetBytes(2)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		FindVideoID("https://www.youtube.com/watch?v=EFJfdWzBHwE&list=PL2ntRZ1ySWBdDyspRTNBIKES1Y-P__59_&index=4")
	}
}

func BenchmarkTestExtPlayerResp(b *testing.B) {
	r := getReader("../../testdata/get_video_info")
	b.SetBytes(2)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		ExtractPlayerResp(r)
	}
}
