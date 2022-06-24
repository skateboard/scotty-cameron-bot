package main

import "Scotty/cmd/cli"

func main() {
	c := cli.New("Scotty", "0.0.1")
	c.Start()
}
