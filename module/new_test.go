package module

import (
	"github.com/skateboard/scotty-cameron-bot/internal/account"
	"github.com/skateboard/scotty-cameron-bot/internal/harvester"
	"github.com/skateboard/scotty-cameron-bot/internal/task"
	"testing"
)

func TestTest(t *testing.T) {
	harv := harvester.New("harvester", "https://www.scottycameron.com/")
	go func() {
		err := harv.Initialize()
		if err != nil {
			t.Error(err)
		}
	}()
	defer harv.Destroy()

	module := ScottyTask{
		Base:       task.New("test"),
		parameters: Parameters{},
		options: Options{
			Sku:        "33700",
			CategoryID: "51",
		},
		account: account.Account{
			Email:    "test@hoku.app",
			Password: "Cool!2345",
		},
		harvester: harv,
	}

	task.RunTask(&module)
}
