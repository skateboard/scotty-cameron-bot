package new

import (
	"Scotty/internal/task"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/drizzleaio/http"
	"log"
	"strings"
	"time"
)

type GeneratePayload struct {
	Site      string `json:"site"`
	UserAgent string `json:"userAgent"`
	Abck      string `json:"abck"`
}

type GenerateResponse struct {
	Data string `json:"data"`
}

func (t *ScottyTask) submitSensor(amount int) {
	currentSubmitted := 0

	for amount > currentSubmitted {
		sensor := t.generatePayload()
		if sensor == nil {
			log.Println("Submit Sensor: Failed to generate sensor.")
			continue
		}

		ok := t.postSensor(*sensor)
		if !ok {
			log.Println("Submit Sensor: Failed to post sensor.")
			continue
		}

		currentSubmitted++

		time.Sleep(time.Second)
	}
}

func (t *ScottyTask) postSensor(sensorData string) bool {
	data := fmt.Sprintf(`{"sensor_data": "%v"}`, sensorData)

	request, err := http.NewRequest("POST", "https://www.scottycameron.com"+t.parameters.AkamaiSensorUrl, strings.NewReader(data))
	if err != nil {
		log.Println("Post Sensor: Failed to create request.")
		return false
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
		log.Println("Post Sensor: Failed to send request.")
		return false
	}

	if response.StatusCode != 201 {
		log.Println("Post Sensor: Failed to send request.")
		return false
	}

	return true
}

func (t *ScottyTask) generatePayload() *string {
	abck := t.Client.GetCookieFromDomain("https://www.scottycameron.com", "_abck", true)
	if abck == nil {
		return nil
	}

	payload := GeneratePayload{
		Site:      "https://www.scottycameron.com/",
		UserAgent: t.parameters.UserAgent,
		Abck:      abck.Value,
	}
	jsonBytes, _ := json.Marshal(payload)

	request, err := http.NewRequest("POST", "https://api.scrimcloud.xyz/v2/akamai/sensor", bytes.NewBuffer(jsonBytes))
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

func (t *ScottyTask) getAkamaiSensorEP() task.State {
	request, err := http.NewRequest("GET", "https://www.scottycameron.com"+t.parameters.AkamaiSensorUrl, nil)
	if err != nil {
		return FailedGetAkamai
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
		return FailedGetAkamai
	}

	if response.StatusCode != 200 {
		return FailedGetAkamai
	}

	t.submitSensor(3)
	return GetPixel
}
