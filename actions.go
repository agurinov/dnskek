package main

import (
	"fmt"
	"net"
	"os"
	"os/exec"
	"os/user"
	"runtime"
	"text/template"

	"github.com/urfave/cli"

	"github.com/agurinov/dnskek/log"
	"github.com/agurinov/dnskek/server"
)

var (
	// Actions
	serveActionUsage     = "Run DNS server"
	installActionUsage   = "Prints terminal commands to INSTALL dnskek"
	uninstallActionUsage = "Prints terminal commands to UNINSTALL dnskek"
	// Flags
	zoneFlagUsage  = "DNS zone for local server"
	portFlagUsage  = "Port for running DNS server on"
	debugFlagUsage = "Debugging mode"
	ttlFlagUsage   = "TTL for DNS records. Also, this flag used as ttl of Docker machine registry"
)

func serveAction(c *cli.Context) {
	log.SetDebug(c.GlobalBool("debug"))
	// Phase 1. Create and start server
	server.New(
		net.ParseIP("127.0.0.1"),
		c.GlobalInt("port"),
		nil,
		c.GlobalInt("ttl"),
	).Serve()
}

func installAction(c *cli.Context) {
	log.SetDebug(c.GlobalBool("debug"))
	// Phase 1. get templates and help templates functions
	funcMap := template.FuncMap{
		"executable": os.Executable,
		"user":       user.Current,
		"lookup":     exec.LookPath,
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
}

func uninstallAction(c *cli.Context) {
	log.SetDebug(c.GlobalBool("debug"))
	// Phase 1. get templates and help templates functions
	funcMap := template.FuncMap{
		"executable": os.Executable,
		"user":       user.Current,
		"lookup":     exec.LookPath,
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
}
