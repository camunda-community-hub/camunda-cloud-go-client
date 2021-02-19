/*
Copyright © 2021 Matheus Cruz matheuscruz.dev@gmail.com

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
	name       string
	id         string
	channel    string
	generation string
	region     string
	plan       string
)

var (
	deleteExample = `

  # Delete cluster by id
  cc clusters delete --id=<cluster_id>`

	getExample = `

  # List all clusters
  cc clusters get --all
   
  # Get cluster by name
  cc clusters get --name=<cluster_name> (If your cluster have a composite name, use: --name='<cluster name>')

  # Get params to create a cluster
  cc clusters get --params`
	createExample = `

  # Create cluster with default configuration
  cc clusters create --default --name=<cluster_name> (If your cluster have a composite name, use: --name='<cluster name>')
 
  # Crate cluster with custom configuration
  cc clusters create 
    --name=<cluster_name> (If your cluster have a composite name, use: --name='<cluster name>'
	--channel=<channel_id>
	--generation=<generation_id>
	--region=<region_id>
	--plan=<plan_type_id>`
)

var clusterCmd = &cobra.Command{
	Use:   "clusters [options]",
	Short: "Manage your cluster's resources on Camunda Cloud",
	Long:  "Used together [OPTIONS] like get, create, delete for manage your resources on Camunda Cloud. For example:" + getExample + createExample + deleteExample,
}

var getClusterCmd = &cobra.Command{
	Use:   "get",
	Short: "Get clusters",
	Long:  "Used together with clusters command, to get your clusters on Camunda Cloud. For example:" + getExample,
	Run: func(cmd *cobra.Command, args []string) {
		all, _ := cmd.Flags().GetBool("all")
		params, _ := cmd.Flags().GetBool("params")

		if name != "" && params && all {
			fmt.Println("Error: --all and --name and --params cannot be specified together")
			return
		}

		if name != "" && params {
			fmt.Println("Error: --name and --params cannot be specified together")
			return
		}

		if name != "" && all {
			fmt.Println("Error: --all and --name cannot be specified together")
			return
		}

		if all && params {
			fmt.Println("Error: --all and --params cannot be specified together")
			return
		}

		if name != "" {
			cluster, _ := client.GetClusterByName(name)
			showCluster(cluster)
			return
		}

		if all {
			clusters, _ := client.GetClusters()
			showClusters(clusters)
			return
		}

		if params {
			params, _ := client.GetClusterParams()
			showParams(*params)
		}

	},
}

var createClusterCmd = &cobra.Command{
	Use:   "create",
	Short: "Create cluster",
	Long:  "Used together clusters command, to create your clusters on Camunda Cloud. For example:" + createExample,
	PreRun: func(cmd *cobra.Command, args []string) {
		def, _ := cmd.Flags().GetBool("default")

		if !def {
			cmd.MarkFlagRequired("channel")
			cmd.MarkFlagRequired("generation")
			cmd.MarkFlagRequired("region")
			cmd.MarkFlagRequired("plan")
		}
	},
	Run: func(cmd *cobra.Command, args []string) {

		def, _ := cmd.Flags().GetBool("default")

		if !def {
			clusterID, err := client.CreateClusterCustomConfig(client.NewClusterCreationParams(
				name, channel, generation, region, plan,
			))

			if err != nil {
				fmt.Println(err.Error())
			} else {
				fmt.Println("Cluster create successfully. Cluster id:", clusterID)
			}
		} else {

			if name != "" {

				client.GetClusterParams()
				clusterID, err := client.CreateClusterDefault(name)

				if err != nil {
					fmt.Println(err.Error())
				} else {
					fmt.Println("Cluster created successfully. Cluster id:", clusterID)
				}
			}
		}
	},
}

var deleteClusterCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete cluster",
	Long:  "Used together with clusters command, to delete your clusters on Camunda Cloud. For example:" + deleteExample,
	Run: func(cmd *cobra.Command, args []string) {

		if id != "" {
			success, _ := client.DeleteCluster(id)

			if success {
				fmt.Println("Cluster deleted successfully")
			} else {
				fmt.Println("Error: We can't delete your cluster")
			}
		}
	},
}

func init() {

	clusterCmd.AddCommand(getClusterCmd)
	clusterCmd.AddCommand(createClusterCmd)
	clusterCmd.AddCommand(deleteClusterCmd)
	rootCmd.AddCommand(clusterCmd)

	// get cmd
	getClusterCmd.Flags().BoolP("all", "a", false, "Get all clusters: camunda-cloud-go-cli get --all")
	getClusterCmd.Flags().BoolP("params", "p", false, "Get params to create a cluster: camunda-cloud-go-cli get --params")
	getClusterCmd.Flags().StringVarP(&name, "name", "n", "", "camunda-cloud-go-cli clusters get --name='<cluster_name>'")

	// delete cmd
	deleteClusterCmd.PersistentFlags().StringVarP(&id, "id", "i", "", "camunda-cloud-go-cli clusters delete --id=<cluster_id>")
	deleteClusterCmd.MarkFlagRequired("id")

	// create cmd
	createClusterCmd.Flags().BoolP("default", "d", false, "camunda-cloud-go-cli clusters create --default=(true|false)")
	createClusterCmd.Flags().StringVarP(&name, "name", "n", "", "Cluster's name")
	createClusterCmd.Flags().StringVarP(&channel, "channel", "c", "", "Cluster's channel id")
	createClusterCmd.Flags().StringVarP(&generation, "generation", "g", "", "Cluster's generation id")
	createClusterCmd.Flags().StringVarP(&region, "region", "r", "", "Cluster's region id")
	createClusterCmd.Flags().StringVarP(&plan, "plan", "p", "", "Cluster's plan type id")
	createClusterCmd.MarkFlagRequired("name")
}

func showCluster(cluster client.Cluster) {
	data, _ := json.MarshalIndent(cluster, "", "  ")
	fmt.Println(string(data))
}

func showClusters(clusters []client.Cluster) {
	data, _ := json.MarshalIndent(clusters, "", "  ")
	fmt.Println(string(data))
}

func showParams(params client.ClusterParams) {
	data, _ := json.MarshalIndent(params, "", "  ")
	fmt.Println(string(data))
}
