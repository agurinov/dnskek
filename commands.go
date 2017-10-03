package main

import (
	"fmt"
	"net"
	"os"
	"os/user"
	"runtime"
	"text/template"

	"github.com/urfave/cli"

	"github.com/agurinov/dnskek/log"
	"github.com/agurinov/dnskek/server"
)

func serveCmd(c *cli.Context) {
	log.SetDebug(c.GlobalBool("debug"))
	// Phase 1. Create and start server
	server.New(
		net.ParseIP("127.0.0.1"),
		c.GlobalInt("port"),
		nil,
		c.GlobalInt("ttl"),
	).Serve()
}

func installCmd(c *cli.Context) {
	log.SetDebug(c.GlobalBool("debug"))
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
}

func uninstallCmd(c *cli.Context) {
	log.SetDebug(c.GlobalBool("debug"))
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
}
