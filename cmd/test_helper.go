package cmd

import (
	"github.com/urfave/cli/v2"

	"github.com/bianchidotdev/yahr/core"
)

func MockData() string {
	return `requests:
  monkey_island:
    host: monkeyisland.example.com
    requests:
      escape:
        method: "post"
        path: /escape
      characters:
        path: /characters
  davie_jones:
    host: daviejones.example.com
    requests:
      locker:
        path: /locker`
}

func MockApp(mockData string) *cli.App {
	if mockData == "" {
		mockData = MockData()
	}
	return &cli.App{
		Name: "yahr",
		Commands: []*cli.Command{
			RequestCmd,
			RunCmd,
		},
		Before: func(cCtx *cli.Context) error {
			requestData := mockData
			config, err := core.ParseConfig([]byte(requestData))
			if err != nil {
				return err
			}
			err = core.SetConfig(config)
			if err != nil {
				return err
			}
			return nil
		},
	}
}
