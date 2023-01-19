package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
)

var rootCmd = &cobra.Command{
	Use:   "yahr",
	Short: "yahr is a yaml-driven http client",
	Long: `A yaml-driven http client for being able to easily define
and run http requests and easily share them with your team.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Do Stuff Here
	},
}

var cfgFile string

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is ./config.yaml)")
}

func initConfig() {
	// TODO: don't blow up on missing config

	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		viper.AddConfigPath("./")
		viper.SetConfigName("yahr")
		// // Find home directory.
		// home, err := homedir.Dir()
		// if err != nil {
		//   log.Println(err)
		//   os.Exit(1)
		// }

		// // Search config in home directory with name ".yahr" (without extension).
		// viper.AddConfigPath(home)
		// viper.SetConfigName(".yahr")
	}
	if err := viper.ReadInConfig(); err != nil {
		log.Fatal("Can't read config:", err)
	}
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
