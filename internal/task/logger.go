package task

import (
	"fmt"
	"github.com/gookit/color"
	"log"
)

//func TaskLog(message string, clr string, ID string) {
//	go func() {
//		switch clr {
//		case "cyan":
//			log.Println(Cyan(fmt.Sprintf("[Task %s] %s", ID, message)))
//		case "blue":
//			log.Println(Blue(fmt.Sprintf("[Task %s] %s", ID, message)))
//		case "red":
//			log.Println(Red(fmt.Sprintf("[Task %s] %s", ID, message)))
//		case "yellow":
//			log.Println(Yellow(fmt.Sprintf("[Task %s] %s", ID, message)))
//		case "green":
//			log.Println(Green(fmt.Sprintf("[Task %s] %s", ID, message)))
//		}
//	}()
//
//}

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
