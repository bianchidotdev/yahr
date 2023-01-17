package cmd

import (
	"testing"

	"github.com/spf13/viper"
)

func TestVersionCmd(t *testing.T) {
	viper.AddConfigPath("../")
	err := versionCmd.Execute()
	if err != nil {
		t.Errorf("Failed to run version command, %v", err)
	}
}
