package task

import (
	"fmt"
	"github.com/gookit/color"
	"log"
)

var (
	Red    = color.FgRed.Render
	Green  = color.FgGreen.Render
	Yellow = color.FgYellow.Render
	Cyan   = color.FgCyan.Render
)

func (b *Base) Info(message interface{}) {
	log.Println(Cyan(fmt.Sprintf("[Task %s] %s", b.ID, message)))
}

func (b *Base) Warning(message interface{}) {
	log.Println(Yellow(fmt.Sprintf("[Task %s] %s", b.ID, message)))
}

func (b *Base) Success(message interface{}) {
	log.Println(Green(fmt.Sprintf("[Task %s] %s", b.ID, message)))
}

func (b *Base) Error(message interface{}) {
	log.Println(Red(fmt.Sprintf("[Task %s] %s", b.ID, message)))
}
