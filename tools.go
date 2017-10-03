package main

import (
	"fmt"

	"github.com/urfave/cli"
	// "dnskek/log"
)

func loggingCmdWrapper(a func(c *cli.Context)) func(c *cli.Context) {
	fmt.Println("DMJKDNDJKNDKJDNJKDNJKDNDKJNDJKDNJKDNJDKNDJKDNJKDNJDKNDJKDNJKDNDJKDNJKDNJDKNDJK")
	// set debug mode
	// log.SetDebug(c.GlobalBool("debug"))

	return a
}
