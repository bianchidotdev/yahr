package cmd

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"regexp"
	"testing"
)

func buildMockData(host string) string {
	return fmt.Sprintf(`requests:
  monkey_island:
    host: %s
    scheme: http
    requests:
      escape:
        method: "post"
        path: /escape
      characters:
        path: /characters`, host)
}

func TestRunCmd(t *testing.T) {
	t.Run("with no args", func(t *testing.T) {
		var output bytes.Buffer // capture output

		app := MockApp("")
		app.Writer = &output

		args := []string{"yahr", "run"}

		err := app.Run(args)

		// Verify
		if err != nil {
			t.Error(err)
			t.FailNow()
		}

		regex := `No group or request specified\n\nNAME:.*`
		match, err := regexp.Match(regex, []byte(output.String()))
		if err != nil {
			t.Errorf("got error %s", err)
		}
		if !match {
			t.Errorf("expected '%s' to match but got '%s'", regex, output.String())
		}
	})
	t.Run("with a request", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintln(w, `Strike yer colors!`)
		}))
		defer ts.Close()

		mockServerURL, err := url.Parse(ts.URL)
		if err != nil {
			t.Error(err)
			t.FailNow()
		}

		var output bytes.Buffer // capture output

		mockYaml := buildMockData(mockServerURL.Host)

		app := MockApp(mockYaml)
		app.Writer = &output

		args := []string{"yahr", "run", "monkey_island", "escape"}

		err = app.Run(args)

		// Verify
		if err != nil {
			t.Error(err)
			t.FailNow()
		}

		matches := []string{
			`POST /escape HTTP/1.1`,
			`Status: 200`,
			`Response Body:\nStrike yer colors!`,
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
	t.Run("with a request and silent flag", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintln(w, `Hoist the Jolly Roger!`)
		}))
		defer ts.Close()

		mockServerURL, err := url.Parse(ts.URL)
		if err != nil {
			t.Error(err)
			t.FailNow()
		}

		var output bytes.Buffer // capture output

		mockYaml := buildMockData(mockServerURL.Host)

		app := MockApp(mockYaml)
		app.Writer = &output

		args := []string{"yahr", "run", "-s", "monkey_island", "escape"}

		err = app.Run(args)

		// Verify
		if err != nil {
			t.Error(err)
			t.FailNow()
		}

		expected := "Hoist the Jolly Roger!\n"

		if expected != output.String() {
			t.Errorf("expected '%s' but got '%s'", expected, output.String())
		}
	})
}
