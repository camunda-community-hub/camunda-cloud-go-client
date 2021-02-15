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
	"github.com/spf13/cobra"
)

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Gets a most important information about the specified camunda cloud resources",
	Long: `You can get specific information about camunda cloud cluster resources. For example:

# Get all clusters
camunda-cloud-go-client get cluster --all

# Get cluster
camunda-cloud-go-client get cluster <cluster-id>

# Get cluster creation parameters
camunda-cloud-go-client get cluster --parameters

# Get all Zeebe clients
camunda-cloud-go-client zeebe --cluster <cluster_id> --all

# Get Zeebe client details
camunda-cloud-go-client zeebe --cluster <cluster_id> --client <client_id>
	`,
}

func init() {
	rootCmd.AddCommand(getCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// getCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// getCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
