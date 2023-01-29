package main

import (
	"log"
	"os"
	"time"
	"github.com/spf13/viper"
	"github.com/urfave/cli/v2"

	"github.com/michaeldbianchi/yahr/cmd"
)

// used by goreleaser
var (
	shortened  = false
	version    = "dev"
	commit     = "none"
	date       = "unknown"
	output     = "json"
)

func main() {
	viper.SetConfigFile("./yahr.yaml")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatal("Can't read config:", err)
	}

	app := &cli.App{
        Name:  "yahr",
		Version:  version,
		Compiled: time.Now(),
        Authors: []*cli.Author{
            &cli.Author{
                Name:  "Michael Bianchi",
                Email: "michael@bianchi.dev",
            },
        },
        Usage: `A yaml-driven http client for being able to easily define
and run http requests and easily share them with your team.`,
        Commands: []*cli.Command{
			cmd.RequestCmd,
			cmd.RunCmd,
		},
    }

    if err := app.Run(os.Args); err != nil {
        log.Fatal(err)
    }
}
