/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/lppgo/nova/utils/log"

	"github.com/flylog/colorstyle"

	"github.com/spf13/cobra"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

var (
	cfgFile         string
	AgentHost       string
	AgentServerPort string
	Tag             string // eg: redis:latest
	NameSpace       string
	DownloadDir     string
	Rename          string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "pull-cli",
	Short: "pull-cli for download file...",
	Long: `pull-cli是一个下载通过P2P方式下载文件的CLI. 
For example:
pull-cli --agent-host=127.0.0.1 --agent-server-port=16002 --tag=redis:latest --download-dir=/tmp/nuwa/ --rename=redis2`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		log.Infof("start pull :%s", colorstyle.Yellow(Tag))

		if AgentServerPort == "" {
			log.Fatalf("agent-server-port is not empty")
		}
		names := strings.Split(Tag, ":")
		if len(names) != 2 {
			log.Fatalf("tag is invalid")
		}

		addr := AgentHost + ":" + AgentServerPort
		filename := Rename
		if Rename == "" {
			filename = names[0]
		}

		err := PullFile(addr, Tag, NameSpace, filename, DownloadDir)
		if err != nil {
			log.Fatalf("pull file is err:%s", err.Error())
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.pull-cli.yaml)")

	rootCmd.PersistentFlags().StringVar(&AgentHost, "agent-host", "localhost", "the agent server host/IP")
	rootCmd.PersistentFlags().StringVar(&AgentServerPort, "agent-server-port", "", "the agent server port")
	rootCmd.PersistentFlags().StringVar(&Tag, "tag", "", "the tag is identify for download file")
	rootCmd.PersistentFlags().StringVar(&NameSpace, "namespace", ".*", "represent a namespace")
	rootCmd.PersistentFlags().StringVar(&DownloadDir, "download-dir", "", "represent download directory")
	rootCmd.PersistentFlags().StringVar(&DownloadDir, "rename", "", "the file rename,default origin name")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".pull-cli" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".pull-cli")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}
