package cli

import "github.com/skateboard/scotty-cameron-bot/cmd/console"

type Cli struct {
	User    string
	Version string
}

func New(user, version string) *Cli {
	return &Cli{
		User:    user,
		Version: version,
	}
}

func (c *Cli) Start() {
	c.Success("Welcome to ScottyIO!")
	c.Info("Version: " + c.Version)
	console.SetBase(c.Version)

	c.mainMenu()
}
