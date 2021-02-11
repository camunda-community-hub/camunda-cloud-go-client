package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
	login()

	getClusterParams()

	createCluster()

	//fmt.Fprintf(w, "Hello from:  "+title+"\n")
}

type AuthPayload struct {
	AccessToken string `json:"access_token"`
}

type ClusterParams struct{
	Channels[] Channel `json:"channels"`
	ClusterPlanTypes[] ClusterPlantType `json:"clusterPlanTypes"`
	Regions[] Region `json:"regions"`
}

type Channel struct {
	Id   string `json:"uuid"`
	Name string `json:"name"`
	AllowedGeneration[] AllowedGeneration `json:"allowedGenerations"`
}

type AllowedGeneration struct {
	Id   string `json:"uuid"`
	Name string `json:"name"`
}

type ClusterPlantType struct {
	Id   string `json:"uuid"`
	Name string `json:"name"`

}

type Region struct {
	Id   string `json:"uuid"`
	Name string `json:"name"`
	Region string `json:"region"`
	Zone string `json:"zone"`

}
var authPayload AuthPayload

var clusterParams ClusterParams

func getClusterParams() {
	req, _ := http.NewRequest("GET", "https://api.cloud.camunda.io/clusters/parameters", nil)
	req.Header.Set("Authorization", "Bearer "+authPayload.AccessToken)

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		log.Fatalf("failed to create client, %v", err)
	}
	//fmt.Println("response Status:", resp.Status)
	//fmt.Println("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	//fmt.Println("response Body:", string(body))
	err2 := json.Unmarshal(body, &clusterParams)
	if err2 != nil {
		log.Fatalf("failed to parse body, %v", err2)
	}
	marshal, _ := json.Marshal(clusterParams)

	fmt.Println("parsed: ",string(marshal) )
}
func createCluster() {

	var jsonStr = []byte(`{
  "name": "my cool cluster",
  
}`)
	req, err := http.NewRequest("POST", "https://api.cloud.camunda.io/clusters/", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer " + authPayload.AccessToken)

	client := &http.Client{}
	resp, err := client.Do(req)

	defer resp.Body.Close()
	fmt.Println("\n\n\nCreate Cluster Response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response Body:", string(body))
	if err != nil {
		log.Fatalf("failed to create client, %v", err)
	}
}
func login() {

	var jsonStr = []byte(`{"grant_type":"client_credentials", "audience":"api.cloud.camunda.io", "client_id":"7v7k.rQj199QUD-Y", "client_secret":"2~XrPs.QyM7PGFWluRalkiSzBZT1INZL"}`)

	req, err := http.NewRequest("POST", "https://login.cloud.camunda.io/oauth/token", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)

	defer resp.Body.Close()

	if err != nil {
		log.Fatalf("failed to create client, %v", err)
	}

	//fmt.Println("response Status:", resp.Status)
	//fmt.Println("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	//fmt.Println("response Body:", string(body))

	err2 := json.Unmarshal(body, &authPayload)
	if err2 != nil {
		log.Fatalf("failed to parse body, %v", err2)
	}

	//fmt.Println("Response Status:", responseBody.AccessToken)
}

func main() {
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}
