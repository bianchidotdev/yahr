package core

import (
	"os"
	"testing"
)

func TestReadConfig(t *testing.T) {
	t.Run("with a non-existent file", func(t *testing.T) {
		file := "../fixtures/super-fake.yaml"
		config, err := ReadConfig(file)

		if err == nil {
			t.Errorf("expected file read error, but got config %s", config)
		}
	})

	t.Run("with a valid yaml file", func(t *testing.T) {
		os.Setenv("ENV_MATEY", "Wha' be an env var")
		file := "../fixtures/pirate.yaml"
		config, err := ReadConfig(file)

		if err != nil {
			t.Error(err)
			t.FailNow()
		}

		expected := `requests:
  monkeyness:
    host: pirate.monkeyness.com
    headers:
      who-are-ye: Wha' be an env var
    requests:
      insult:
        path: /api/insult
      translate:
        path: /api/translate
        queryparams:
          english: where is the nearest restroom
`
		if expected != string(config) {
			t.Errorf("expected '%s' but got '%s'", expected, config)
		}

	})
}

func TestParseConfig(t *testing.T) {
	// TODO: implement tests
	t.Run("with an empty config", func(t *testing.T) {})
	t.Run("with an config missing requests key", func(t *testing.T) {})
	t.Run("with an config containing requests", func(t *testing.T) {})
}
