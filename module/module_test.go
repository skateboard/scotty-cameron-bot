package module

import (
	"Scotty/internal/account"
	"Scotty/internal/harvester"
	"Scotty/internal/task"
	"testing"
)

func TestModule(t *testing.T) {
	harv := harvester.New("harvester", "https://www.scottycameron.com/")
	go func() {
		err := harv.Initialize()
		if err != nil {
			t.Error(err)
		}
	}()
	defer harv.Destroy()

	module := ScottyTask{
		TaskBase:   task.New("test"),
		parameters: Parameters{},
		options: Options{
			Sku:        "33700",
			CategoryID: "51",
			harvester:  harv,
		},
		account: account.Account{
			Email:    "test@hoku.app",
			Password: "Cool!2345",
		},
	}

	task.RunTask(&module)
}
