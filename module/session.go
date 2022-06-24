package module

import (
	"Scotty/internal/task"
	"github.com/drizzleaio/http"
	"log"
	"regexp"
)

var (
	sensorReg = regexp.MustCompile(`<script type="text/javascript"\s\ssrc=\"(.*?)\"`)
	relicId   = regexp.MustCompile(`(?m)xpid:"(.*)",li`)
)

func (t *ScottyTask) getSession2() task.State {
	//https://www.scottycameron.com/store/accessories/
	request, err := http.NewRequest("GET", "https://www.scottycameron.com/store/accessories/", nil)
	if err != nil {
		log.Println(err)
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
		log.Println(err)
		return FailedSession
	}

	if response.StatusCode != 200 {
		log.Println("Failed to get session")
		return FailedSession
	}

	return GetProduct
}

func (t *ScottyTask) getSession() task.State {
	request, err := http.NewRequest("GET", "https://www.scottycameron.com/", nil)
	if err != nil {
		log.Println(err)
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
		log.Println(err)
		return FailedSession
	}

	if response.StatusCode != 200 {
		log.Println("Failed to get session")
		return FailedSession
	}

	b, err := response.GetBodyBytes()
	if err != nil {
		log.Println(err)
		return FailedSession
	}

	matches := sensorReg.FindAllStringSubmatch(string(b), -1)
	if len(matches) == 0 {
		log.Println("No matches")
		return FailedSession
	}

	matches1 := relicId.FindAllStringSubmatch(string(b), -1)
	if len(matches1) == 0 {
		log.Println("No matches")
		return FailedSession
	}

	t.parameters.RelicID = matches1[0][1]
	t.parameters.AkamaiSensor = matches[0][1]

	//{stamp:%27-1%27%2Cnecessary:true%2Cpreferences:true%2Cstatistics:true%2Cmarketing:true%2Cver:1%2Cutc:1655875523833%2Cregion:%27US-06%27}
	//t.Client.SetCookies("https://www.scottycameron.com/", []*http.Cookie{
	//	{
	//		Name:  "_WebStorePublishState",
	//		Value: "PRODUCTION",
	//	},
	//	{
	//		Name:  "_WebStoreculture",
	//		Value: "1",
	//	},
	//	{
	//		Name:  "culture",
	//		Value: "en",
	//	},
	//})

	return GetAkamai
}
