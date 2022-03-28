package cmd

import (
	"fmt"
	"os"
	"path"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var logger *zap.Logger
var cfgFile string

var rootCmd = &cobra.Command{
	Use:     "flusso",
	Short:   "A distributed commit log service",
	Version: "0.1.0",
}

func init() {
	cobra.OnInitialize(initLogger, initConfig)
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "config file (default is $HOME/.flusso/config.yaml)")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		flussodir := path.Join(home, ".flusso")
		viper.AddConfigPath(flussodir)
		viper.SetConfigType("yaml")
		viper.SetConfigName("config")
	}

	// read in environment variables that match
	viper.AutomaticEnv()

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		// logger.Debug(fmt.Sprintf("using config file: %s", viper.ConfigFileUsed()))
	}
}

func initLogger() {
	l, err := zap.NewDevelopment()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	logger = l.Named("cli")
}
