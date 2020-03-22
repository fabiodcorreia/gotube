package info

import (
	"reflect"
	"testing"

	"github.com/fabiodcorreia/gotube/internal/serial"
)

func TestGetVideoInfo(t *testing.T) {
	type args struct {
		player serial.PlayerResponse
	}
	tests := []struct {
		name          string
		args          args
		wantVideoInfo VideoInfo
		wantErr       bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVideoInfo, err := GetVideoInfo(tt.args.player)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetVideoInfo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotVideoInfo, tt.wantVideoInfo) {
				t.Errorf("GetVideoInfo() = %v, want %v", gotVideoInfo, tt.wantVideoInfo)
			}
		})
	}
}

func Test_streamFromFormat(t *testing.T) {
	type args struct {
		ft []serial.Format
		st []Stream
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{

		{
			name: "Invalid Input Streams",
			args: args{
				ft: []serial.Format{{Itag: 1, ContentLength: "123", MimeType: "video/mp4;+codecs=\"avc1.64001F,+mp4a.40.2\"", Quality: "360p", URL: "http://..."}},
				st: nil,
			},
			wantErr: true,
		},
		{
			name: "Invalid Input Formats",
			args: args{
				ft: nil,
				st: make([]Stream, 0),
			},
			wantErr: true,
		},
		{
			name: "Valid Streams and Formats",
			args: args{
				ft: []serial.Format{{Itag: 1, ContentLength: "123", MimeType: "video/mp4;+codecs=\"avc1.64001F,+mp4a.40.2\"", Quality: "360p", URL: "http://..."}},
				st: make([]Stream, 1),
			},
			wantErr: false,
		},
		{
			name: "Invalid MimeType",
			args: args{
				ft: []serial.Format{{Itag: 1, ContentLength: "123", MimeType: "video", Quality: "360p", URL: "http://..."}},
				st: make([]Stream, 1),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if err := streamFromFormat(tt.args.ft, tt.args.st); (err != nil) != tt.wantErr {
				t.Errorf("streamFromFormat() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_extensionFromType(t *testing.T) {
	type args struct {
		mimeType string
	}
	tests := []struct {
		name    string
		args    args
		wantExt VideoExt
	}{

		{
			name: "Mp4-360p",
			args: args{
				mimeType: "video/mp4;+codecs=\"avc1.42001E,+mp4a.40.2\"",
			},
			wantExt: MP4,
		},

		{
			name: "Mp4-720p",
			args: args{
				mimeType: "video/mp4;+codecs=\"avc1.640028\"",
			},
			wantExt: MP4,
		},

		{
			name: "Mp4-1080p",
			args: args{
				mimeType: "video/mp4;+codecs=\"avc1.64002a\"",
			},
			wantExt: MP4,
		},

		{
			name: "3gp",
			args: args{
				mimeType: "video/3gp;+codecs=\"vp9\"",
			},
			wantExt: TGP,
		},

		{
			name: "avi",
			args: args{
				mimeType: "video/avi;+codecs=\"vp9\"",
			},
			wantExt: AVI,
		},

		{
			name: "flv",
			args: args{
				mimeType: "video/flv;+codecs=\"vp9\"",
			},
			wantExt: FLV,
		},

		{
			name: "Webm-1080p",
			args: args{
				mimeType: "video/webm;+codecs=\"vp9\"",
			},
			wantExt: WEBM,
		},

		{
			name: "Empty Type",
			args: args{
				mimeType: "",
			},
			wantExt: AVI,
		},

		{
			name: "Unknown Type",
			args: args{
				mimeType: "something",
			},
			wantExt: AVI,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotExt := extensionFromType(tt.args.mimeType)
			if gotExt != tt.wantExt {
				t.Errorf("extensionFromType() = %v, want %v", gotExt, tt.wantExt)
			}
		})
	}
}

// Benchmark with mime.ExtensionsByType is more costly
//
// BenchmarkExtensionFromType-12    	17404563	       688 ns/op	   2.91 MB/s	     456 B/op	       6 allocs/op
//
// against this implementation with
//
// BenchmarkExtensionFromType-12    	50939827	       232 ns/op	   8.62 MB/s	      32 B/op	       1 allocs/op
func BenchmarkExtensionFromType(b *testing.B) {
	b.SetBytes(2)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		extensionFromType("video/mp4;+codecs=\"avc1.64001F,+mp4a.40.2\"")
	}
}

func BenchmarkStreamFromFormat(b *testing.B) {
	b.SetBytes(2)
	b.ResetTimer()
	ft := []serial.Format{
		{Itag: 1, ContentLength: "123", MimeType: "video/mp4;+codecs=\"avc1.64001F,+mp4a.40.2\"", Quality: "360p", URL: "http://..."},
		{Itag: 2, ContentLength: "123", MimeType: "video/mp4;+codecs=\"avc1.64001F,+mp4a.40.2\"", Quality: "360p", URL: "http://..."},
		{Itag: 3, ContentLength: "123", MimeType: "video/mp4;+codecs=\"avc1.64001F,+mp4a.40.2\"", Quality: "360p", URL: "http://..."},
		{Itag: 4, ContentLength: "123", MimeType: "video/mp4;+codecs=\"avc1.64001F,+mp4a.40.2\"", Quality: "360p", URL: "http://..."},
		{Itag: 5, ContentLength: "123", MimeType: "video/mp4;+codecs=\"avc1.64001F,+mp4a.40.2\"", Quality: "360p", URL: "http://..."},
		{Itag: 6, ContentLength: "123", MimeType: "video/mp4;+codecs=\"avc1.64001F,+mp4a.40.2\"", Quality: "360p", URL: "http://..."},
		{Itag: 7, ContentLength: "123", MimeType: "video/mp4;+codecs=\"avc1.64001F,+mp4a.40.2\"", Quality: "360p", URL: "http://..."},
		{Itag: 8, ContentLength: "123", MimeType: "video/mp4;+codecs=\"avc1.64001F,+mp4a.40.2\"", Quality: "360p", URL: "http://..."},
		{Itag: 9, ContentLength: "123", MimeType: "video/mp4;+codecs=\"avc1.64001F,+mp4a.40.2\"", Quality: "360p", URL: "http://..."},
		{Itag: 10, ContentLength: "123", MimeType: "video/mp4;+codecs=\"avc1.64001F,+mp4a.40.2\"", Quality: "360p", URL: "http://..."},
		{Itag: 11, ContentLength: "123", MimeType: "video/mp4;+codecs=\"avc1.64001F,+mp4a.40.2\"", Quality: "360p", URL: "http://..."},
		{Itag: 12, ContentLength: "123", MimeType: "video/mp4;+codecs=\"avc1.64001F,+mp4a.40.2\"", Quality: "360p", URL: "http://..."},
		{Itag: 13, ContentLength: "123", MimeType: "video/mp4;+codecs=\"avc1.64001F,+mp4a.40.2\"", Quality: "360p", URL: "http://..."},
		{Itag: 14, ContentLength: "123", MimeType: "video/mp4;+codecs=\"avc1.64001F,+mp4a.40.2\"", Quality: "360p", URL: "http://..."},
		{Itag: 15, ContentLength: "123", MimeType: "video/mp4;+codecs=\"avc1.64001F,+mp4a.40.2\"", Quality: "360p", URL: "http://..."},
		{Itag: 16, ContentLength: "123", MimeType: "video/mp4;+codecs=\"avc1.64001F,+mp4a.40.2\"", Quality: "360p", URL: "http://..."},
		{Itag: 17, ContentLength: "123", MimeType: "video/mp4;+codecs=\"avc1.64001F,+mp4a.40.2\"", Quality: "360p", URL: "http://..."},
		{Itag: 18, ContentLength: "123", MimeType: "video/mp4;+codecs=\"avc1.64001F,+mp4a.40.2\"", Quality: "360p", URL: "http://..."},
		{Itag: 19, ContentLength: "123", MimeType: "video/mp4;+codecs=\"avc1.64001F,+mp4a.40.2\"", Quality: "360p", URL: "http://..."},
		{Itag: 20, ContentLength: "123", MimeType: "video/mp4;+codecs=\"avc1.64001F,+mp4a.40.2\"", Quality: "360p", URL: "http://..."},
	}

	for i := 0; i < b.N; i++ {
		st := make([]Stream, len(ft))
		streamFromFormat(ft, st)
	}
}

func BenchmarkGetVideoInfo(b *testing.B) {
	b.SetBytes(2)
	b.ResetTimer()
	p := serial.PlayerResponse{
		Playability: serial.PlayabilityStatus{
			Status: "OK",
		},
		Data: serial.StreamingData{
			Formats: []serial.Format{
				{
					Itag:          1,
					ContentLength: "1233",
					MimeType:      "video/mp4;+codecs=\"avc1.64001F,+mp4a.40.2\"",
					Quality:       "720p",
					URL:           "https://......",
				},
				{
					Itag:          2,
					ContentLength: "1233",
					MimeType:      "video/mp4;+codecs=\"avc1.64001F,+mp4a.40.2\"",
					Quality:       "720p",
					URL:           "https://......",
				},
				{
					Itag:          3,
					ContentLength: "1233",
					MimeType:      "video/mp4;+codecs=\"avc1.64001F,+mp4a.40.2\"",
					Quality:       "720p",
					URL:           "https://......",
				},
				{
					Itag:          4,
					ContentLength: "1233",
					MimeType:      "video/mp4;+codecs=\"avc1.64001F,+mp4a.40.2\"",
					Quality:       "720p",
					URL:           "https://......",
				},
				{
					Itag:          5,
					ContentLength: "1233",
					MimeType:      "video/mp4;+codecs=\"avc1.64001F,+mp4a.40.2\"",
					Quality:       "720p",
					URL:           "https://......",
				},
			},
		},
		Details: serial.VideoDetails{
			Duration: "12331",
			Title:    "Building+Hexagonal+Microservices+with+Go+-+Part+Three",
			VideoID:  "hrroijewoij",
		},
	}

	for i := 0; i < b.N; i++ {
		GetVideoInfo(p)
	}
}
