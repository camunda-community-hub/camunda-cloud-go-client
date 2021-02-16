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

	"github.com/salaboy/camunda-cloud-go-client/pkg/cc/client"
	"github.com/spf13/cobra"
)

func handlerGet(value string, cmd cobra.Command) {

	if value != "get" {
		return
	}

	all, _ := cmd.Flags().GetBool("all")

	if all {
		ok := client.Login("bq~NI~MzG5KacsZn", "qZ1EwAdFF2F5MgyfnjnSO2TAvitMBn1s")

		if ok {
			fmt.Println("all:", all)
			client.GetClusters()
		}

		return
	}

	params, _ := cmd.Flags().GetBool("params")

	if params {
		ok := client.Login("bq~NI~MzG5KacsZn", "qZ1EwAdFF2F5MgyfnjnSO2TAvitMBn1s")

		if ok {
			fmt.Println("params:", params)
			client.GetClusterParams()
		}
	}

}

func handlerCreate(value string, cmd cobra.Command) {
	if value != "create" {
		return
	}

	name, _ := cmd.Flags().GetString("name")

	if name != "" {

		ok := client.Login("bq~NI~MzG5KacsZn", "qZ1EwAdFF2F5MgyfnjnSO2TAvitMBn1s")

		if ok {
			client.GetClusterParams()
			client.CreateCluster(name)
		}
	}

}

func handlerDelete(value string, cmd cobra.Command) {
	if value == "delete" {

		fmt.Println(value)
		id, _ := cmd.Flags().GetString("id")

		fmt.Println("id:", id)

		if id != "" {

			ok := client.Login("iQfNDE7Yrupv5tXl", "QWMZz6Se6hlGAEgcNpMTepA4~1B1.NMQ")

			if ok {
				client.DeleteCluster(id)
			}
		}

	}
}

// clusterCmd represents the cluster command
var clusterCmd = &cobra.Command{
	Use:   "clusters",
	Short: "Manage your cluster's resources on Camunda Cloud",
	Long:  `Add conditional long usage here, that dependents of parent`,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		handlers := []func(string, cobra.Command){
			handlerDelete,
			handlerCreate,
			handlerGet,
		}

		if len(args) > 0 {
			for _, handler := range handlers {
				handler(args[0], *cmd)
			}
		}

	},
}

func init() {
	rootCmd.AddCommand(clusterCmd)
	clusterCmd.Flags().BoolP("all", "a", false, "Get all clusters")
	clusterCmd.Flags().BoolP("params", "p", false, "Get cluster creation parameters")
	clusterCmd.Flags().String("name", "", "Create cluster")
	clusterCmd.Flags().String("id", "", "Cluster's id")
}
