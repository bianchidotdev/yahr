package main

import (
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/joho/godotenv/autoload"
	"github.com/urfave/cli/v2"

	"github.com/michaeldbianchi/yahr/cmd"
	"github.com/michaeldbianchi/yahr/core"
)

// used by goreleaser
var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
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
			&cli.StringFlag{
				Name:    "cfgFile",
				Aliases: []string{"c"},
				Value:   "./yahr.yaml",
				EnvVars: []string{"YAHR_CONFIG_FILE"},
			},
		},
		Commands: []*cli.Command{
			cmd.RequestCmd,
			cmd.RunCmd,
		},
		Before: func(cCtx *cli.Context) error {
			// err := godotenv.Load()
			// if err != nil {
			// 	return fmt.Errorf("Error loading .env file - %s", err)
			// }

			configBytes, err := core.ReadConfig(cCtx.String("cfgFile"))
			if err != nil {
				return fmt.Errorf("Error reading config - %s", err)
			}

			appConfig, err := core.ParseConfig(configBytes)

			err = core.SetConfig(appConfig)
			if err != nil {
				return fmt.Errorf("Error reading config - %s", err)
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
