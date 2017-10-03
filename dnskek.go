package main

import (
	"fmt"
	"net"
	"os"
	"os/user"
	"runtime"
	"sort"
	"text/template"
	"time"

	"github.com/urfave/cli"

	"bitbucket.org/agurinov/dnskek/log"
	"bitbucket.org/agurinov/dnskek/server"
)

func main() {
	// Phase 1. Get cli options, some validation checks and configure working env
	// errors from this phase must be paniced with traceback and os.exit(1)
	app := cli.NewApp()
	app.Name = "dnskek"
	app.Compiled = time.Now()
	app.Authors = []cli.Author{
		cli.Author{
			Name:  "Alexander Gurinov",
			Email: "alexander.gurinov@gmail.com",
		},
	}
	app.Usage = "DNS server across docker-machines"
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
			Name:  "serve",
			Usage: "Run DNS server",
			Action: func(c *cli.Context) error {
				// set debug mode
				log.SetDebug(c.GlobalBool("debug"))
				// let's start server
				server.New(
					net.ParseIP("127.0.0.1"),
					c.GlobalInt("port"),
					nil,
					c.GlobalInt("ttl"),
				).Serve()
				return nil
			},
		},
		{
			Name:  "install",
			Usage: "Prints terminal commands to INSTALL dnskek",
			Action: func(c *cli.Context) error {

				// TODO create some common wrapper for enable debug for all actions

				// log.Info("DMKDMKL", c.App.Flags)

				// Phase 1. get templates and help templates functions
				funcMap := template.FuncMap{
					"executable": os.Executable,
					"user":       user.Current,
				}
				tplChain := template.Must(
					template.New("").Funcs(funcMap).ParseGlob(
						fmt.Sprintf("./tpls/%s-%s/*.tpl", runtime.GOOS, runtime.GOARCH),
					),
				)

				// Phase 2. Output with actual dynamic data
				data := map[string]interface{}{
					"appName": c.App.Name,
					"debug":   c.GlobalBool("debug"),
					"ip":      "127.0.0.1",
					"port":    c.GlobalInt("port"),
					"zone":    c.GlobalString("zone"),
				}
				tplChain.ExecuteTemplate(os.Stdout, "install.tpl", data)

				return nil
			},
		},
		{
			Name:  "uninstall",
			Usage: "Prints terminal commands to UNINSTALL dnskek",
			Action: func(c *cli.Context) error {
				// Phase 1. get templates and help templates functions
				funcMap := template.FuncMap{
					"executable": os.Executable,
					"user":       user.Current,
				}
				tplChain := template.Must(
					template.New("").Funcs(funcMap).ParseGlob(
						fmt.Sprintf("./tpls/%s-%s/*.tpl", runtime.GOOS, runtime.GOARCH),
					),
				)

				// Phase 2. Output with actual dynamic data
				data := map[string]interface{}{
					"appName": c.App.Name,
					"zone":    c.GlobalString("zone"),
				}
				tplChain.ExecuteTemplate(os.Stdout, "uninstall.tpl", data)

				return nil
			},
		},
	}
	// configure sorting for help
	sort.Sort(cli.FlagsByName(app.Flags))
	sort.Sort(cli.CommandsByName(app.Commands))
	// run
	app.Run(os.Args)
}
