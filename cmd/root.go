package cmd

import (
	"fmt"
	"os"
	"path"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use:   "rkenum",
	Short: "Need an enum? Rye Knot",
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

var cfgFile string

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(
		&cfgFile, "config", "",
		"config file (default is $HOME/.config/rk/rkenum.toml)")
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)
		configDir := path.Join(home, ".config", "rk")
		viper.AddConfigPath(configDir)
		viper.SetConfigType("toml")
		viper.SetConfigName("rkenum")
	}
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		_, err = fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
		if err != nil {
			panic(err)
		}
	}
}
