package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/notfoundy/reamcli/internal/app"
	"github.com/notfoundy/reamcli/internal/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func main() {
	Execute()
}

var cfgFile, cfgDir, cfgPath string

var rootCmd = &cobra.Command{
	Run: func(cmd *cobra.Command, args []string) {
		app, err := app.NewApp()
		if err == nil {
			app.Run()
		}
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $XDG_CONFIG_HOME/reamcli/config.yaml)")

	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		cfgDir, err := utils.GetConfigDir()
		cobra.CheckErr(err)
		cfgPath = filepath.Join(cfgDir, "config.yaml")
		viper.SetConfigFile(cfgPath)
		viper.SetConfigType("yaml")
	}

	viper.AutomaticEnv() // read in environment variables that match

	if _, err := os.Stat(viper.ConfigFileUsed()); os.IsNotExist(err) {
		err := os.WriteFile(viper.ConfigFileUsed(), []byte("# reamcli config\n"), 0600)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Failed to create config file:", err)
		}
	}

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}
