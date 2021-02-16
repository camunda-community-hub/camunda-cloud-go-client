package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/salaboy/camunda-cloud-go-client/pkg/cc/client"
)

var clientId = os.Getenv("CC_CLIENT_ID")
var clientSecret = os.Getenv("CC_CLIENT_SECRET")

func handler(w http.ResponseWriter, r *http.Request) {

	fmt.Fprintf(w, "Client ID:  "+clientId+"\n")
	fmt.Fprintf(w, "Client Secret:  "+clientSecret+"\n")

	var clusterName = "hello there"

	fmt.Println("Attempting to Login ...")
	var loginOk = client.Login(clientId, clientSecret)
	if loginOk {
		fmt.Println("Login Successful!")
		fmt.Println("Fetching Cluster Creation Params ...")
		client.GetClusterParams()

		fmt.Println("Creating Cluster", clusterName, " ... ")

		var clusterId = client.CreateCluster(clusterName)

		fmt.Println("Cluster", clusterName, " created with Id: ", clusterId)

		for true {
			time.Sleep(5 * time.Second)
			var status = client.GetClusterDetails(clusterId)
			if status.Ready == "Healthy" {
				marshal, _ := json.Marshal(status)
				fmt.Println("> Cluster Status and details: ", string(marshal))
			} else {
				fmt.Println("> Waiting For the Cluster: ", clusterId, " to be ready ... ")
			}
		}

		//fmt.Println("Deleting Cluster with Id: ", clusterId)
		//
		//var deleted = client.DeleteCluster(clusterId)
		//
		//fmt.Println("Cluster with Id: ", clusterId, "deleted: ", deleted)

	} else {
		fmt.Println("Login Failed.")
	}
	//fmt.Fprintf(w, "Hello from:  "+title+"\n")
}

func main() {
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
	// cmd.Execute()
}
