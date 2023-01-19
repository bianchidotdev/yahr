package main

import (
	"fmt"
	"github.com/michaeldbianchi/yahr/cmd"
	"github.com/spf13/cobra"
)

// used by goreleaser
var (
	shortened  = false
	version    = "dev"
	commit     = "none"
	date       = "unknown"
	output     = "json"
	versionCmd = &cobra.Command{
		Use:   "version",
		Short: "Version will output the current build information",
		Run: func(_ *cobra.Command, _ []string) {
			fmt.Println("Version:", version)
		},
	}
)

func main() {
	cmd.Execute(versionCmd)
}
