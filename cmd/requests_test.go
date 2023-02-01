package cmd

import (
	"bytes"
	"regexp"
	"testing"

	"github.com/urfave/cli/v2"

	"github.com/michaeldbianchi/yahr/core"
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
		Flags: []cli.Flag{
			&cli.StringFlag{Name: "cfgFile", Value: "../examples/pirate.yaml"},
		},
		Commands: []*cli.Command{
			RequestCmd,
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

func TestRequestsListCmd(t *testing.T) {
	t.Run("with no args", func(t *testing.T) {
		var output bytes.Buffer // capture output

		app := MockApp("")
		app.Writer = &output

		args := []string{"yahr", "requests", "list"}

		err := app.Run(args)

		// Verify
		if err != nil {
			t.Error(err)
			t.FailNow()
		}

		matches := []string{
			`| Group         | Name       | Method | Endpoint                                    |`,
			`| monkey_island | escape     | post   | https://monkeyisland.example.com/escape     |`,
			`| monkey_island | characters | get    | https://monkeyisland.example.com/characters |`,
			`| davie_jones   | locker     | get    | https://davie_jones.example.com/locker      |`,
		}

		for _, regex := range matches {
			match, err := regexp.Match(regex, []byte(output.String()))
			if !match {
				t.Errorf("expected '%s' to match but got '%s'", regex, output.String())
			}
			if err != nil {
				t.Errorf("got error %s", err)
			}
		}
	})

	t.Run("with a group specified", func(t *testing.T) {
		var output bytes.Buffer // capture output

		app := MockApp("")
		app.Writer = &output

		args := []string{"yahr", "requests", "list", "davie_jones"}

		err := app.Run(args)

		// Verify
		if err != nil {
			t.Error(err)
			t.FailNow()
		}

		matches := []string{
			`| Group         | Name       | Method | Endpoint                                    |`,
			`| davie_jones   | locker     | get    | https://davie_jones.example.com/locker      |`,
		}

		for _, regex := range matches {
			match, err := regexp.Match(regex, []byte(output.String()))
			if !match {
				t.Errorf("expected '%s' to match but got '%s'", regex, output.String())
			}
			if err != nil {
				t.Errorf("got error %s", err)
			}
		}
	})

	t.Run("with no requests found", func(t *testing.T) {
		var output bytes.Buffer // capture output

		app := MockApp("requests:\n  group:")
		app.Writer = &output

		args := []string{"yahr", "requests", "list"}

		err := app.Run(args)

		// Verify
		if err != nil {
			t.Error(err)
			t.FailNow()
		}

		expected := "No requests found\n"
		if output.String() != expected {
			t.Errorf("expected '%s', but got '%s'", expected, output.String())
		}
	})
}
