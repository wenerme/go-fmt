// Copyright © 2018 wener <wenermail@gmail.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"
	"os"

	"github.com/Sirupsen/logrus"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string
var verbose = 2

var rootCfg = &_rootCfg{}

type _rootCfg struct {
	Log struct {
		Level string
	}
}

var rootCmd = &cobra.Command{
	Use:   "goform",
	Short: "Various file form handler for go",
	Long:  ``,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.goform.yaml)")
	rootCmd.PersistentFlags().CountVarP(&verbose, "verbose", "v", "Show more verboses")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".goform" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".goform")
	}

	viper.AutomaticEnv() // read in environment variables that match

	if err := viper.ReadInConfig(); err == nil {
		logrus.Info("Using config file: {}", viper.ConfigFileUsed())
	} else {
		// Ignore not found
	}

	if err := viper.Unmarshal(rootCfg); err != nil {
		logrus.Fatal(err)
	}

	if verbose > len(logrus.AllLevels)-1 {
		verbose = len(logrus.AllLevels) - 1
	}
	rootCfg.Log.Level = logrus.AllLevels[verbose].String()

	if l, err := logrus.ParseLevel(rootCfg.Log.Level); err == nil {
		logrus.SetLevel(l)
	} else {
		logrus.Fatal(err)
	}

	logrus.Debug("set log level to ", rootCfg.Log.Level)
}
