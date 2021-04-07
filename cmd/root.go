/*
Copyright © 2021

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

	cc "github.com/camunda-community-hub/camunda-cloud-go-client/pkg/cc/client"

	"github.com/spf13/cobra"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

var cfgFile string
var client cc.CCClient

var ClientId = os.Getenv("CC_CLIENT_ID")
var ClientSecret = os.Getenv("CC_CLIENT_SECRET")

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:                   "cc-ctl",
	DisableFlagsInUseLine: true,
	Short:                 "Camunda Cloud CLI to manage Camunda Cloud Resources",
	Long: `Camunda Cloud CLI to interact with Camunda Cloud Resources.
  You can create a Camunda Cloud Account here: https://accounts.cloud.camunda.io/signup
  Then you need to go to <YOUR USER>(Top right corner) -> Organization Settings -> Cloud Management API 
  and Create a new client. You need to copy and save the Client Id and the Secret Id from that client. 
  You need to export the following variables for this command to interact with a Camunda Cloud Account:
  - export CC_CLIENT_ID=<YOUR CLIENT ID>
  - export CC_CLIENT_SECRET=<YOUR CLIENT SECRET>
  
  Available Commands:  
  # List all clusters
  cc-ctl clusters get --all

  # Get cluster from id
  cc-ctl clusters get --id <cluster_id>

  # Get cluster from name
  cc-ctl clusters get --name <cluster_name>

  # Delete cluster from id
  cc-ctl clusters delete --id <cluster_id>

  # Delete cluster from name
  cc-ctl clusters delete --name <cluster_name>

  # Create cluster from default configuration
  cc-ctl clusters create --default --name <cluster_name>`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if ClientId == "" || ClientSecret == "" {
		fmt.Println(rootCmd.Long)
		os.Exit(1)
	}
	login, err := client.Login(ClientId, ClientSecret)

	if err != nil || !login {

		fmt.Errorf("Error trying to Login to Camunda Cloud, "+
			"please check your CC_CLIENT_ID and CC_CLIENT_SECRET! \n %s", err)
		os.Exit(1)

	}
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	cobra.OnInitialize(initConfig)
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

		// Search config in home directory with name ".camunda-cloud-go-client" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".camunda-cloud-go-client")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}
