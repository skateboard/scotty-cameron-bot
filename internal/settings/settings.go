package settings

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"sync"
)

var (
	settingsOnce sync.Once
	settings     *Settings
)

type Settings struct {
	Key            string `json:"key"`
	DiscordWebhook string `json:"discordWebhook"`
}

func GetSettings() *Settings {
	settingsOnce.Do(func() {
		settings = LoadSettings()
	})

	return settings
}

func LoadSettings() *Settings {
	var settings Settings
	data, err := ioutil.ReadFile("settings/settings.json")
	if err != nil {
		fmt.Println(err)
		return nil
	}
	err = json.Unmarshal(data, &settings)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	return &settings
}
