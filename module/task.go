package module

import (
	"Scotty/internal/account"
	"Scotty/internal/harvester"
	"Scotty/internal/task"
	"fmt"
	"github.com/drizzleaio/http"
	"regexp"
	"strings"
	"time"
)

type ScottyTask struct {
	*task.Base

	options Options

	parameters Parameters

	account account.Account

	harvester *harvester.Harvester
}

type Options struct {
	Sku        string
	CategoryID string
	Url        string
}

type Parameters struct {
	UserAgent       string
	SecChUaPlatform string
	SecChUA         string

	RelicID         string
	AkamaiSensorUrl string

	PixelUrl       string
	PixelUrl2      string
	PixelID        string
	PixelVersion   string
	PixelScriptVal string
	PixelT         string

	RequestVerificationID string
}

const (
	Session task.State = iota
	FailedSession

	GetAkamai
	FailedGetAkamai

	GetPixel
	FailedGetPixel

	SubmitPixel
	FailedSubmitPixel

	LoginParameters
	FailedLoginParameters

	Login
	FailedLogin
)

var (
	requestVerifyTokenRegex = regexp.MustCompile(`(?m)<input name="__RequestVerificationToken" type="hidden" value="(.*)" />`)
)

func (t *ScottyTask) Next(state task.State) (task.State, error) {
	switch state {
	case task.Initialize:
		t.Info("Initializing...")
		t.Client.UpdateServerName("www.scottycameron.com")
		t.Client.SetProxy("http://127.0.0.1:8888")
		t.Client.DoRedirect(false)

		t.parameters.UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.4896.127 Safari/537.36"

		if !t.parseUserAgent() {
			t.Stop()
			return 0, nil
		}

		_ = t.Client.SetCookies("https://www.scottycameron.com/", []*http.Cookie{
			{
				Name:  "CookieConsent",
				Value: "{stamp:%27-1%27%2Cnecessary:true%2Cpreferences:true%2Cstatistics:true%2Cmarketing:true%2Cver:1%2Cutc:1655816472136%2Cregion:%27IN%27}",
			},
			{
				Name:  "scottyCameronHasVisitedTwo",
				Value: "true",
			},
			{
				Name:  "culture",
				Value: "en",
			},
			{
				Name:  "persona",
				Value: "other",
			},
			{
				Name:  "isAdmin",
				Value: "False",
			},
		})

		t.Info("Initialized!")
		return Session, nil
	case Session:
		t.Info("Getting session...")
		return t.getSession(), nil
	case FailedSession:
		t.Error("Failed to get session!")
		time.Sleep(time.Second * 5)
		return Session, nil
	case GetAkamai:
		t.Info("Getting Akamai...")
		return t.getAkamaiSensorEP(), nil
	case FailedGetAkamai:
		t.Error("Failed to get Akamai!")
	case GetPixel:
		t.Info("Getting Pixel...")
		return t.getPixel(), nil
	case FailedGetPixel:
		t.Error("Failed to get pixel!")
	//case SubmitPixel:
	//	t.Info("Submitting pixel...")
	//	return t.submitPixel(), nil
	//case FailedSubmitPixel:
	//	t.Error("Failed to submit pixel!")
	case LoginParameters:
		t.Info("Getting Login Parameters..")
		return t.getLoginParameters(), nil
	case FailedLoginParameters:
		t.Error("Failed to get Login Parameters!")
	case Login:
		t.Info("Logging in...")
		return t.login(), nil
	case FailedLogin:
		t.Error("Failed to login!")
	}

	t.Stop()
	return 0, nil
}

func (t *ScottyTask) parseUserAgent() bool {
	re, err := regexp.Compile("Chrome/([0-9]*)")
	if err != nil {
		fmt.Println(err)
		return false
	}
	version := re.FindAllStringSubmatch(t.parameters.UserAgent, -1)
	if len(version) == 0 {
		fmt.Println("Failed to parse user agent!")
		return false
	}
	t.parameters.SecChUA = fmt.Sprintf(`"Not A;Brand";v="99", "Chromium";v="%v", "Google Chrome";v="%v"`, version[0][1], version[0][1])

	if strings.Contains(t.parameters.UserAgent, "Windows NT") {
		t.parameters.SecChUaPlatform = "\"Windows\""
	} else if strings.Contains(t.parameters.UserAgent, "Mac OS X") {
		t.parameters.SecChUaPlatform = "\"macOS\""
	} else {
		t.parameters.SecChUaPlatform = "\"Android\""
	}

	return true
}
