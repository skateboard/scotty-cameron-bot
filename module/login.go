package module

import (
	"Scotty/internal/task"
	"github.com/drizzleaio/http"
	"log"
	"net/url"
	"strings"
)

func (t *ScottyTask) login() task.State {
	solv, err := t.options.harvester.Harvest("51829642-2cda-4b09-896c-594f89d700cc")
	if err != nil {
		return FailedLogin
	}

	form := url.Values{}
	form.Add("__RequestVerificationToken", t.parameters.RequestVerificationToken)
	form.Add("Username", "test@hoku.app")
	form.Add("Password", "Cool!2345")
	form.Add("g-recaptcha-response", solv)
	form.Add("h-recaptcha-response", solv)
	form.Add("RememberMe", "true")
	form.Add("RememberMe", "false")

	request, err := http.NewRequest("POST", "https://www.scottycameron.com/store/user/login/", strings.NewReader(form.Encode()))
	if err != nil {
		return FailedLogin
	}

	request.SetHeaders([]map[string]string{
		{"content-length": ""},
		{"content-type": "application/x-www-form-urlencoded"},
		{"user-agent": t.parameters.UserAgent},
		{"accept": "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9"},
		{"sec-fetch-site": "same-origin"},
		{"sec-fetch-mode": "navigate"},
		{"sec-fetch-dest": "document"},
		{"sec-ch-ua": `" Not A;Brand";v="99", "Chromium";v="102", "Google Chrome";v="102"`},
		{"sec-ch-ua-mobile": "?0"},
		{"sec-ch-ua-platform": "windows"},
		{"referer": "https://www.scottycameron.com/store/user/login/"},
		{"accept-encoding": "gzip, deflate, br"},
		{"accept-language": "en-US,en;q=0.9"},
	})

	response, err := t.Client.Do(request)
	if err != nil {
		return FailedLogin
	}

	if response.StatusCode != 302 {
		return FailedLogin
	}

	log.Println("Logged in")

	t.Stop()
	return 0
}

func (t *ScottyTask) getLoginParameters() task.State {
	request, err := http.NewRequest("GET", "https://www.scottycameron.com/store/user/login/", nil)
	if err != nil {
		return FailedLoginParameters
	}

	request.SetHeaders([]map[string]string{
		{"user-agent": t.parameters.UserAgent},
		{"accept": "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9"},
		{"sec-fetch-site": "same-origin"},
		{"sec-fetch-mode": "navigate"},
		{"sec-fetch-dest": "document"},
		{"sec-ch-ua": `" Not A;Brand";v="99", "Chromium";v="102", "Google Chrome";v="102"`},
		{"sec-ch-ua-mobile": "?0"},
		{"sec-ch-ua-platform": "windows"},
		{"referer": "https://www.scottycameron.com/store/user/signup"},
		{"accept-encoding": "gzip, deflate, br"},
		{"accept-language": "en-US,en;q=0.9"},
	})

	response, err := t.Client.Do(request)
	if err != nil {
		return FailedLoginParameters
	}

	if response.StatusCode != 200 {
		return FailedLoginParameters
	}

	b, err := response.GetBodyBytes()
	if err != nil {
		return FailedLoginParameters
	}

	matches := requestVerifyTokenRegex.FindAllStringSubmatch(string(b), -1)
	if len(matches) == 0 {
		log.Println("No matches")
		return FailedLoginParameters
	}
	t.parameters.RequestVerificationToken = matches[0][1]

	return t.login()
}
