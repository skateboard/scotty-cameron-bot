package new

import (
	"Scotty/internal/task"
	"fmt"
	"github.com/drizzleaio/http"
	"github.com/justhyped/OrderedForm"
	"strings"
)

func (t *ScottyTask) getLoginParameters() task.State {
	request, err := http.NewRequest("GET", "https://www.scottycameron.com/store/user/login/", nil)
	if err != nil {
		t.Error(err)
		return FailedLoginParameters
	}

	request.SetHeaders([]map[string]string{
		{"upgrade-insecure-requests": "1"},
		{"user-agent": t.parameters.UserAgent},
		{"accept": "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9"},
		{"sec-fetch-site": "same-origin"},
		{"sec-fetch-mode": "navigate"},
		{"sec-fetch-user": "?1"},
		{"sec-fetch-dest": "document"},
		{"sec-ch-ua": t.parameters.SecChUA},
		{"sec-ch-ua-mobile": "?0"},
		{"sec-ch-ua-platform": t.parameters.SecChUaPlatform},
		{"referer": "https://www.scottycameron.com/store/user/signup"},
		{"accept-encoding": "gzip, deflate, br"},
		{"accept-language": "en-US,en;q=0.9"},
	})

	response, err := t.Client.Do(request)
	if err != nil {
		t.Error(err)
		return FailedLoginParameters
	}

	if response.StatusCode != 200 {
		t.Error(fmt.Sprintf("Failed to get login parameters! (%d)", response.StatusCode))
		return FailedLoginParameters
	}

	b, err := response.GetBodyBytes()
	if err != nil {
		t.Error(err)
		return FailedLoginParameters
	}

	matches := requestVerifyTokenRegex.FindAllStringSubmatch(string(b), -1)
	if len(matches) == 0 {
		t.Error("Failed to get request verification token")
		return FailedLoginParameters
	}
	t.parameters.RequestVerificationID = matches[0][1]

	sessID := t.Client.GetCookieFromDomain("https://www.scottycameron.com/", "ASP.NET_SessionId", true)
	if sessID == nil {
		t.Error("Failed to get session ID")
		return FailedLoginParameters
	}
	_ = t.Client.SetCookie("https://www.scottycameron.com/", &http.Cookie{
		Name:  "SCSessionId",
		Value: sessID.Value,
	})

	return Login
}

func (t *ScottyTask) login() task.State {
	t.Warning("Waiting for captcha...")
	solution, err := t.harvester.Harvest("51829642-2cda-4b09-896c-594f89d700cc")
	if err != nil {
		t.Error(err)
		return FailedLogin
	}
	t.Success("Captcha solved!")

	form := new(OrderedForm.OrderedForm)
	form.Set("__RequestVerificationToken", t.parameters.RequestVerificationID)
	form.Set("Username", t.account.Email)
	form.Set("Password", t.account.Password)
	form.Set("g-recaptcha-response", solution)
	form.Set("h-recaptcha-response", solution)
	form.Set("RememberMe", "true")
	form.Set("RememberMe", "false")

	request, err := http.NewRequest("POST", "https://www.scottycameron.com/store/user/login/", strings.NewReader(form.URLEncode()))
	if err != nil {
		t.Error(err)
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
		t.Error(err)
		return FailedLogin
	}

	if response.StatusCode != 302 {
		t.Error(fmt.Sprintf("Failed to login! (%d)", response.StatusCode))
		return FailedLogin
	}

	if loc := response.Header.Get("location"); loc != "/store/" {
		t.Error(fmt.Sprintf("Failed to login! (%s)", loc))
		return FailedLogin
	}

	t.Info("Logged in!")

	t.Stop()
	return 0
}
