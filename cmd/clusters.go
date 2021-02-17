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
	"encoding/json"
	"fmt"

	"github.com/salaboy/camunda-cloud-go-client/pkg/cc/client"
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

var clusterCmd = &cobra.Command{
	Use:   "clusters [options]",
	Short: "Manage your cluster's resources on Camunda Cloud",
	Long: `Used together [OPTIONS] like get, create, delete for manage your resources on Camunda Cloud. For example:

  # List all clusters
  camunda-cloud-go-cli clusters get --all
   
  # Get cluster by name
  camunda-cloud-go-cli clusters get --name=<cluster_name> (If your cluster have a composite name, use: --name='<cluster name>')
	
  # Create cluster with default configuration
  camunda-cloud-go-cli clusters create --default --name=<cluster_name> (If your cluster have a composite name, use: --name='<cluster name>')
 
  # Crate cluster with custom configuration
  camunda-cloud-go-cli clusters create 
    --name=<cluster_name> (If your cluster have a composite name, use: --name='<cluster name>'
	--channel=<channel_id>
	--generation=<generation_id>
	--region=<region_id>
	--plan=<plan_type_id>
	
  # Delete cluster by id
  camunda-cloud-go-cli clusters delete --id=<cluster_id>`,
}

var getClusterCmd = &cobra.Command{
	Use:   "get",
	Short: "Get your clusters resources",
	Long: `Used together with clusters command, to get your clusters on Camunda Cloud. For example:

  # List all clusters
  camunda-cloud-go-cli clusters get --all
 
  # Get cluster by name
  camunda-cloud-go-cli clusters get --name=<cluster_name> (If your cluster have a composite name, use: --name='<cluster name>')`,
	Run: func(cmd *cobra.Command, args []string) {

		if name != "" {

			cluster, _ := client.GetClusterByName(name)
			showCluster(cluster)

		} else {

			all, _ := cmd.Flags().GetBool("all")

			if all {
				clusters, _ := client.GetClusters()
				showClusters(clusters)
			}
		}
	},
}

var createClusterCmd = &cobra.Command{
	Use:   "create",
	Short: "Create cluster",
	Long: `Used together clusters command, to create your clusters on Camunda Cloud. For example:
	Used together with clusters command, to get your clusters on Camunda Cloud. For example:

  # Create cluster with default configuration
  camunda-cloud-go-cli clusters create --default --name=<cluster_name> (If your cluster have a composite name, use: --name='<cluster name>')
 
  # Crate cluster with own configuration
  camunda-cloud-go-cli clusters create 
    --name=<cluster_name> (If your cluster have a composite name, use: --name='<cluster name>'
	--channel=<channel_id>
	--generation=<generation_id>
	--region=<region_id>
	--plan=<plan_type_id>
  `,
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

		if name != "" {

			client.GetClusterParams()
			clusterID, err := client.CreateCluster(name)

			if err != nil {
				fmt.Println(err.Error())
			} else {
				fmt.Println("Cluster created successfully. Cluster id: ", clusterID)
			}
		}
	},
}

var deleteClusterCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete cluster",
	Long: `Used together with clusters command, to delete your clusters on Camunda Cloud. For example:

  # Delete cluster by id
  camunda-cloud-go-cli clusters delete --id=<cluster_id>
`,
	Run: func(cmd *cobra.Command, args []string) {

		if id != "" {
			client.DeleteCluster(id)
			fmt.Println("Cluster deleted successfully")
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
	getClusterCmd.Flags().StringVarP(&name, "name", "n", "", "camunda-cloud-go-cli clusters get --name='<cluster_name>'")

	// delete cmd
	deleteClusterCmd.PersistentFlags().StringVarP(&id, "id", "i", "", "camunda-cloud-go-cli clusters delete --id=<cluster_id>")
	deleteClusterCmd.MarkFlagRequired("id")

	// create cmd
	createClusterCmd.Flags().BoolP("default", "d", true, "camunda-cloud-go-cli clusters create --default=(true|false)")
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