package new

import (
	"Scotty/internal/task"
	"bytes"
	"encoding/base64"
	"encoding/json"
	"github.com/drizzleaio/http"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func (t *ScottyTask) getPixel() task.State {
	request, err := http.NewRequest("GET", t.parameters.PixelUrl, nil)
	if err != nil {
		t.Error(err)
		return FailedGetPixel
	}

	request.SetHeaders([]map[string]string{
		{"user-agent": t.parameters.UserAgent},
		{"accept": "*/*"},
		{"accept-language": "en-US,en;q=0.5"},
		{"accept-encoding": "zip, deflate, br"},
		{"referer": "https://www.scottycameron.com/"},
		{"x-newrelic-id": t.parameters.RelicID},
		{"content-type": "text/plain;charset=UTF-8"},
		{"content-length": ""},
		{"origin": "https://www.scottycameron.com"},
		{"dnt": "1"},
		{"sec-fetch-dest": "document"},
		{"sec-fetch-mode": "navigate"},
		{"sec-fetch-site": "none"},
		{"te": "trailers"},
	})

	response, err := t.Client.Do(request)
	if err != nil {
		t.Error(err)
		return FailedGetPixel
	}

	if response.StatusCode != 200 {
		t.Error(err)
		return FailedGetPixel
	}

	b, err := response.GetBodyBytes()
	if err != nil {
		t.Error(err)
		return FailedGetPixel
	}

	bodyText := string(b)

	s := strings.Index(bodyText, `g=_[`) + 4

	e := strings.Index(bodyText[s:], `]`) + s
	test := bodyText[s:e]

	i, err := strconv.Atoi(test)
	if err != nil {
		log.Println("Error Parsing Pixel")
		time.Sleep(5000 * time.Millisecond)
		return FailedGetPixel
	}

	if i == 0 {
		return FailedGetPixel
	}

	x := kms2(bodyText, `,"`, i) + 1
	y := strings.Index(bodyText[x+1:], `"`) + x + 2

	scripVal, err := strconv.Unquote(bodyText[x:y])
	if err != nil {
		t.Error(err)
		return FailedGetPixel
	}
	t.parameters.PixelScriptVal = scripVal
	if len(t.parameters.PixelScriptVal) < 32 {
		log.Println(err)
		return FailedGetPixel
	}

	str2 := strings.Split(t.parameters.PixelUrl2, "?")[1]
	substring := str2[strings.Index(str2, "=")+1:]
	if strings.Contains(substring, "&") {
		substring = substring[:strings.Index(substring, "&")]
	}
	decoded, _ := base64.StdEncoding.DecodeString(substring)
	re := regexp.MustCompile("t=([0-z]*)")
	matcher := re.FindStringSubmatch(string(decoded))
	if len(matcher) == 0 {
		t.Error("Error Parsing Pixel")
		return FailedGetPixel
	}
	t.parameters.PixelT = matcher[1]

	return LoginParameters
}

func (t *ScottyTask) submitPixel() task.State {
	payload := t.generatePixelPayload()
	if payload == nil {
		t.Error("Error Generating Pixel Payload")
		return FailedSubmitPixel
	}

	request, err := http.NewRequest("POST", t.parameters.PixelUrl, bytes.NewBuffer([]byte(*payload)))
	if err != nil {
		t.Error(err)
		return FailedGetPixel
	}

	request.SetHeaders([]map[string]string{
		{"content-length": ""},
		{"sec-ch-ua": t.parameters.SecChUA},
		{"x-newrelic-id": t.parameters.RelicID},
		{"content-type": "application/x-www-form-urlencoded"},
		{"sec-ch-ua-mobile": "?0"},
		{"user-agent": t.parameters.UserAgent},
		{"sec-ch-ua-platform": t.parameters.SecChUaPlatform},
		{"accept": "*/*"},
		{"origin": "https://www.scottycameron.com"},
		{"sec-fetch-site": "same-origin"},
		{"sec-fetch-mode": "cors"},
		{"sec-fetch-dest": "empty"},
		{"referer": "https://www.scottycameron.com/"},
		{"accept-encoding": "gzip, deflate, br"},
		{"accept-language": "en-US,en;q=0.9"},
	})

	response, err := t.Client.Do(request)
	if err != nil {
		t.Error(err)
		return FailedSubmitPixel
	}

	if response.StatusCode != 200 {
		t.Error(err)
		return FailedSubmitPixel
	}

	return LoginParameters
}

type PixelGeneratePayload struct {
	PixelID   string `json:"pixel_id"`
	T         string `json:"t"`
	ScriptVal string `json:"script_val"`
	UserAgent string `json:"user_agent"`
}

func (t *ScottyTask) generatePixelPayload() *string {
	payload := PixelGeneratePayload{
		PixelID:   t.parameters.PixelID,
		T:         t.parameters.PixelT,
		ScriptVal: t.parameters.PixelScriptVal,
		UserAgent: t.parameters.UserAgent,
	}
	jsonBytes, _ := json.Marshal(payload)

	request, err := http.NewRequest("POST", "https://api.scrimcloud.xyz/v2/akamai/pixel", bytes.NewBuffer(jsonBytes))
	if err != nil {
		return nil
	}

	request.Header.Set("Content-Type", "application/json")

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil
	}

	if response.StatusCode != 200 {
		return nil
	}

	var generateResponse GenerateResponse
	err = response.GetBodyJSON(&generateResponse)
	if err != nil {
		return nil
	}

	return &generateResponse.Data
}

func kms2(s, find string, n int) int {
	i := 0
	for m := 1; m <= n; m++ {
		x := strings.Index(s[i:], find)
		if x < 0 {
			break
		}
		i += x
		if m == n {
			return i
		}
		i += len(find)
	}
	return i
}
