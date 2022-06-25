package module

import (
	"fmt"
	"github.com/drizzleaio/http"
	"github.com/skateboard/scotty-cameron-bot/internal/task"
	"regexp"
	"strings"
)

var (
	sensorReg = regexp.MustCompile(`<script type="text/javascript"\s\ssrc=\"(.*?)\"`)
	relicId   = regexp.MustCompile(`(?m)xpid:"(.*)",li`)
	pixelUrl  = regexp.MustCompile(`(https://www.scottycameron.com/akam.*?)\"`)
)

func (t *ScottyTask) getSession() task.State {
	request, err := http.NewRequest("GET", "https://www.scottycameron.com/", nil)
	if err != nil {
		t.Error(err)
		return FailedSession
	}

	request.SetHeaders([]map[string]string{
		{"user-agent": t.parameters.UserAgent},
		{"accept": "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,*/*;q=0.8"},
		{"accept-language": "en-US,en;q=0.5"},
		{"accept-encoding": "zip, deflate, br"},
		{"dnt": "1"},
		{"upgrade-insecure-requests": "1"},
		{"sec-fetch-dest": "document"},
		{"sec-fetch-mode": "navigate"},
		{"sec-fetch-site": "none"},
		{"sec-fetch-user": "?1"},
		{"te": "trailers"},
	})

	response, err := t.Client.Do(request)
	if err != nil {
		t.Error(err)
		return FailedSession
	}

	if response.StatusCode != 200 {
		t.Error(fmt.Sprintf("Failed to get session! (%d)", response.StatusCode))
		return FailedSession
	}

	b, err := response.GetBodyBytes()
	if err != nil {
		t.Error(err)
		return FailedSession
	}
	text := string(b)

	matches := sensorReg.FindAllStringSubmatch(text, -1)
	if len(matches) == 0 {
		t.Error("Failed to get sensor link!")
		return FailedSession
	}

	matches1 := relicId.FindAllStringSubmatch(text, -1)
	if len(matches1) == 0 {
		t.Error("Failed to get relic id!")
		return FailedSession
	}

	matches3 := pixelUrl.FindAllStringSubmatch(text, -1)
	if len(matches3) == 0 {
		t.Error("Failed to get pixel url!")
		return FailedSession
	}

	t.parameters.RelicID = matches1[0][1]
	t.parameters.AkamaiSensorUrl = matches[0][1]
	t.parameters.PixelUrl = matches3[0][1]
	t.parameters.PixelUrl2 = matches3[1][1]

	c := strings.Index(string(b), `bazadebezolkohpepadr="`) + 22
	t.parameters.PixelID = strings.Split(text[c:], `"<`)[0]

	//Check Pixel Version to use for API
	if strings.Contains(t.parameters.PixelUrl, "13") {
		t.parameters.PixelVersion = "13"
	} else {
		t.parameters.PixelVersion = "11"
	}

	return GetAkamai
}
