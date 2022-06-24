package cli

import (
	"fmt"
	"github.com/gookit/color"
)

var (
	Red    = color.FgRed.Render
	Green  = color.FgGreen.Render
	Yellow = color.FgYellow.Render
	Cyan   = color.FgCyan.Render
)

func (c *Cli) Info(message interface{}) {
	fmt.Println(Cyan(message))
}

func (c *Cli) Warning(message interface{}) {
	fmt.Println(Yellow(message))
}

func (c *Cli) Success(message interface{}) {
	fmt.Println(Green(message))
}

func (c *Cli) Error(message interface{}) {
	fmt.Println(Red(message))
}
