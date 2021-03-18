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
	"encoding/json"
	"fmt"

	"github.com/camunda-community-hub/camunda-cloud-go-client/pkg/cc/client"
	"github.com/spf13/cobra"
)

var (
	cluster string
)

var (
	zeebeClientDeleteExample = ""
	zeebeClientGetExample    = `
  # List all Zeebe clients
  cc-ctl zb-client get --cluster=<cluster_id> --all`
	zeebeClientCreateExample = ""
)

// zbClientCmd represents the zb-client command
func CreateZbClientCmd() *cobra.Command {

	zbClientCmd := &cobra.Command{
		Use:   "zb-client [options]",
		Short: "Manage your zeebe clients resources on Camunda Cloud",
		Long: `Used together [OPTIONS] like get, create, delete for manage your zeebe clients resources on Camunda Cloud. For example:` +
			zeebeClientGetExample + zeebeClientCreateExample + zeebeClientDeleteExample,
	}

	return zbClientCmd
}

func init() {

	zbClientCmd := CreateZbClientCmd()
	zbClientGetCmd := CreateZbClientGetCmd()

	zbClientGetCmd.Flags().StringVarP(&cluster, "cluster", "n", "", "cc-ctl zb-client get --cluster=<cluster_id>")
	zbClientCmd.MarkFlagRequired("cluster")

	zbClientCmd.AddCommand(zbClientGetCmd)
	rootCmd.AddCommand(zbClientCmd)
}

func CreateZbClientGetCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get",
		Short: "Get Zeebe clients",
		Long:  "Used together with zb-client command, to get your zeebe clients on Camunda Cloud. For example:" + zeebeClientGetExample,
		RunE:  ZbClientGetRunE,
	}

	return cmd
}

func ZbClientGetRunE(cmd *cobra.Command, args []string) error {

	clients, err := client.GetZeebeClients(cluster)

	if err != nil {
		fmt.Println("err 1")
		return err
	}

	data, err2 := json.MarshalIndent(clients, "", "  ")

	if err2 != nil {
		return err2
	}

	fmt.Println(string(data))

	return nil
}
