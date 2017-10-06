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
			Usage: zoneFlagUsage,
			Value: "lo",
		},
		cli.IntFlag{
			Name:  "port",
			Usage: portFlagUsage,
			Value: 5354,
		},
		cli.BoolFlag{
			Name:  "debug",
			Usage: debugFlagUsage,
		},
		cli.IntFlag{
			Name:  "ttl",
			Usage: ttlFlagUsage,
			Value: 300,
		},
	}
	app.Commands = []cli.Command{
		{
			Name:   "serve",
			Usage:  serveActionUsage,
			Action: serveAction,
		},
		{
			Name:   "install",
			Usage:  installActionUsage,
			Action: installAction,
		},
		{
			Name:   "uninstall",
			Usage:  uninstallActionUsage,
			Action: uninstallAction,
		},
	}
	// configure sorting for help
	sort.Sort(cli.FlagsByName(app.Flags))
	sort.Sort(cli.CommandsByName(app.Commands))
	// run
	app.Run(os.Args)
}
