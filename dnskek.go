package main

import (
	"os"
	"sort"
	"time"

	"github.com/urfave/cli"
)

const (
	NAME    = "dnskek"
	VERSION = "0.0.1"
)

func main() {
	// Phase 1. Get cli options, some validation checks and configure working env
	// errors from this phase must be paniced with traceback and os.exit(1)
	app := cli.NewApp()
	app.Name = NAME
	app.Version = VERSION
	app.Compiled = time.Now()
	app.Authors = []cli.Author{
		{
			Name:  "Alexander Gurinov",
			Email: "alexander.gurinov@gmail.com",
		},
	}
	app.Usage = "DNS server across docker machines"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "zone",
			Value: "lo",
			Usage: "DNS zone for local server",
		},
		cli.IntFlag{
			Name:  "port",
			Value: 5354,
			Usage: "Port for running DNS server on",
		},
		cli.BoolFlag{
			Name:  "debug",
			Usage: "Debugging mode",
		},
		cli.IntFlag{
			Name:  "ttl",
			Value: 300,
			Usage: "TTL for DNS records",
		},
	}
	app.Commands = []cli.Command{
		{
			Name:   "serve",
			Usage:  "Run DNS server",
			Action: serveCmd,
		},
		{
			Name:   "install",
			Usage:  "Prints terminal commands to INSTALL dnskek",
			Action: installCmd,
		},
		{
			Name:   "uninstall",
			Usage:  "Prints terminal commands to UNINSTALL dnskek",
			Action: uninstallCmd,
		},
	}
	// configure sorting for help
	sort.Sort(cli.FlagsByName(app.Flags))
	sort.Sort(cli.CommandsByName(app.Commands))
	// run
	app.Run(os.Args)
}
