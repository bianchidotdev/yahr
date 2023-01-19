package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
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
var Verbose bool
var Silent bool

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "./yahr.yaml", "config file (default is ./config.yaml)")
	rootCmd.PersistentFlags().BoolVarP(&Verbose, "verbose", "v", false, "verbose output")
	rootCmd.PersistentFlags().BoolVarP(&Silent, "silent", "s", false, "silence all output other than the response body")
	viper.BindPFlag("verbose", rootCmd.PersistentFlags().Lookup("verbose"))
	viper.BindPFlag("silent", rootCmd.PersistentFlags().Lookup("silent"))

	if viper.GetBool("verbose") {
		log.SetOutput(os.Stderr)
	}
}

func initConfig() {
	// TODO: don't blow up on missing config

	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	}
	// else {
	// 	viper.AddConfigPath("./")
	// 	viper.SetConfigName("yahr")

	// 	// // Find home directory.
	// 	// home, err := homedir.Dir()
	// 	// if err != nil {
	// 	//   log.Fatal(err)
	// 	// }

	// 	// // Search config in home directory with name ".yahr" (without extension).
	// 	// viper.AddConfigPath(home)
	// 	// viper.SetConfigName(".yahr")
	// }
	if err := viper.ReadInConfig(); err != nil {
		log.Fatal("Can't read config:", err)
	}
}

func Execute(versionCmd *cobra.Command) {
	rootCmd.AddCommand(versionCmd)
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
