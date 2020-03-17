package serial

import (
	"encoding/json"
	"fmt"
	"io"
)

// playabilityStatusOK represents the success status
const playabilityStatusOK = "OK"

// GetPlayerResponse receives an io.Reader with the player_response JSON and return an object representation
func GetPlayerResponse(jsonStr io.Reader, rp *PlayerResponse) error {
	if jsonStr == nil {
		return fmt.Errorf("serial.GetPlayerResponse.jsonStr: value can't be nil")
	}

	err := json.NewDecoder(jsonStr).Decode(rp)

	if err != nil {
		return fmt.Errorf("serial.GetPlayerResponse.Unmarshal: %w", err)
	}

	if rp.Playability.Status != playabilityStatusOK {
		(*rp) = PlayerResponse{} //reset the rp value to zero value
		return fmt.Errorf("serial.GetPlayerResponse.PlayStatus not Ok: %s", rp.Playability.Status)
	}

	return nil
}

// The structs for JSON unmarshal don't represent the full object but just the required data

// PlayerResponse is the struct representation of player_response
//
// {
//	"playabilityStatus": {...},
//	"streamingData": {...},
//	"videoDetails": {...},
// }
type PlayerResponse struct {
	Playability PlayabilityStatus `json:"playabilityStatus"`
	Details     VideoDetails      `json:"videoDetails"`
	Data        StreamingData     `json:"streamingData"`
}

// PlayabilityStatus is the struct representation of playabilityStatus
//
// {
// 	"status": "OK"
// }
type PlayabilityStatus struct {
	Status string `json:"status"`
}

// VideoDetails is the struct representation of videoDetails
//
// {
//  "videoId": "QybXz4SpPxE",
//  "title": "Video+Title+With+Spaces",
//	"lengthSeconds": "944",
// }
type VideoDetails struct {
	VideoID  string `json:"videoId"`
	Title    string `json:"title"`
	Duration string `json:"lengthSeconds"`
}

// StreamingData is the struct representation of streamingData
//
// "formats": [
//	{...},...
// ],
type StreamingData struct {
	Formats []Format `json:"formats"`
}

// Format is the struct representation of streaming formats
//
// {
// 	"itag": 18,
//	"url": "https:....",
//	"mimeType": "video/mp4;+codecs=\"avc1.42001E,+mp4a.40.2\"",
//	"contentLength": "32870018",
//	"qualityLabel": "360p",
// }
type Format struct {
	Itag          int    `json:"itag"`
	URL           string `json:"url"`
	MimeType      string `json:"mimeType"`
	ContentLength string `json:"contentLength"`
	Quality       string `json:"qualityLabel"`
}
