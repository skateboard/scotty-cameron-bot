package cli

import (
	"fmt"
)

func (c *Cli) mainMenu() func() {
	c.Warning(`
	1. Start Account Gen
	2. Start Tasks
	`)
	c.Info("Enter your choice (1-3): ")
	var input int
	fmt.Scanln(&input)

	switch input {
	case 1:
		return c.mainMenu()
	case 2:

		return c.mainMenu()

	}

	return c.mainMenu()
}
