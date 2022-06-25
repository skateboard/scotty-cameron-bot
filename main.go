package main

import "github.com/skateboard/scotty-cameron-bot/cmd/cli"

func main() {
	c := cli.New("Scotty", "0.0.1")
	c.Start()
}
