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
	// Don't forget to read config either from cfgFile or from home directory!
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		viper.AddConfigPath("./")
		viper.SetConfigName("config")
		// // Find home directory.
		// home, err := homedir.Dir()
		// if err != nil {
		//   fmt.Println(err)
		//   os.Exit(1)
		// }

		// // Search config in home directory with name ".yahr" (without extension).
		// viper.AddConfigPath(home)
		// viper.SetConfigName(".yahr")
	}
	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("Can't read config:", err)
		os.Exit(1)
	}
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
