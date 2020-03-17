package parse

import (
	"bytes"
	"fmt"
	"io"
	"net/url"
	"regexp"
	"strings"
)

const (
	statusOK   = "ok"
	statusFail = "fail"

	paramStatus = "status"
	paramReason = "reason"
	paramPlayer = "player_response"
)

// ExtractPlayerResp receives the content of get_video_info and extracts the player_response json
func ExtractPlayerResp(gvInfo io.Reader) (player io.Reader, err error) {
	if gvInfo == nil {
		return player, fmt.Errorf("parse.ExtPlayerResp.query: value can't be empty")
	}

	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(gvInfo)

	if err != nil {
		return player, fmt.Errorf("parse.ExtPlayerResp.Read: %w", err)
	}

	params, err := url.ParseQuery(buf.String())
	if err != nil {
		return player, fmt.Errorf("parse.ExtPlayerResp.ParseQuery: %w", err)
	}

	status, fStatus := params[paramStatus]
	if !fStatus {
		return player, fmt.Errorf("parse.ExtPlayerResp.Status: status param not found")
	}

	if status[0] != statusOK {
		if status[0] == statusFail {
			reason, fReason := params[paramReason]
			if fReason {
				return player, fmt.Errorf("parse.ExtPlayerResp.Status: status not ok - '%s' with reason '%s'", status[0], reason[0])
			}
		}
		return player, fmt.Errorf("parse.ExtPlayerResp.Status: status not ok - '%s'", status[0])
	}

	p, pStatus := params[paramPlayer]
	if !pStatus {
		return player, fmt.Errorf("parse.ExtPlayerResp.PlayerResponse: player_response param not found")
	}

	return strings.NewReader(p[0]), nil
}

// videoIDPattern is the regex pattern to get the VideoID
const (
	videoIDPattern = "v=(\\w{11})"
	maxLenURL      = 150
)

// videoRegexp is the compiled regex to find videoID on the URL
var videoRegexp = regexp.MustCompile(videoIDPattern)

// FindVideoID receives an youtube video url and return the videoId of that video
func FindVideoID(url string) (videoID string, err error) {
	if url == "" {
		return videoID, fmt.Errorf("parse.FindVideoID.url: url can't be empty")
	}

	if len(url) > maxLenURL {
		return videoID, fmt.Errorf("parse.FindVideoID.url: url can't have more then %d chars", maxLenURL)
	}

	match := videoRegexp.FindStringSubmatch(url)

	if match == nil || len(match) == 1 {
		return videoID, fmt.Errorf("parse.FindVideoID: VideoID not found on URL %s", url)
	}

	return match[1], err
}
