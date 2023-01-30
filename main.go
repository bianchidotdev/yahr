package main

import (
	"log"
	"os"
	"time"

	"github.com/michaeldbianchi/yahr/cmd"
	"github.com/michaeldbianchi/yahr/common"
	"github.com/urfave/cli/v2"
)

// used by goreleaser
var (
	shortened = false
	version   = "dev"
	commit    = "none"
	date      = "unknown"
	output    = "json"
)

func NewApp() *cli.App {
	app := &cli.App{
		Name:     "yahr",
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
		Flags: []cli.Flag{
			&cli.StringFlag{Name: "cfgFile", Aliases: []string{"c"}, Value: "./yahr.yaml"},
		},
		Commands: []*cli.Command{
			cmd.RequestCmd,
			cmd.RunCmd,
		},
		Before: func(cCtx *cli.Context) error {
			configBytes, err := common.ReadConfig(cCtx.String("cfgFile"))
			if err != nil {
				log.Println("Error reading config: ", err)
				return err
			}

			appConfig, err := common.ParseConfig(configBytes)

			err = common.SetConfig(appConfig)
			if err != nil {
				log.Println("Error reading config: ", err)
				return err
			}
			return nil
		},
	}
	return app
}

func main() {
	app := NewApp()

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
