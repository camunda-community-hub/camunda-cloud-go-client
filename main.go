package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

var clientId = os.Getenv("CC_CLIENT_ID")
var clientSecret = os.Getenv("CC_CLIENT_SECRET")

func handler(w http.ResponseWriter, r *http.Request) {

	fmt.Fprintf(w, "Client ID:  "+clientId+"\n")
	fmt.Fprintf(w, "Client Secret:  "+clientSecret+"\n")

	var clusterName = "hello there"

	fmt.Println("Attempting to Login ...")
	var loginOk = login(clientId, clientSecret)
	if loginOk {
		fmt.Println("Login Successful!")
		fmt.Println("Fetching Cluster Creation Params ...")
		getClusterParams()

		fmt.Println("Creating Cluster", clusterName, " ... ")

		var clusterId = createCluster(clusterName)

		fmt.Println("Cluster",clusterName, " created with Id: ", clusterId)

	}else{
		fmt.Println("Login Failed.")
	}
	//fmt.Fprintf(w, "Hello from:  "+title+"\n")
}

type ClusterCreationParams struct {
	ClusterName  string `json:"name"`
	ChannelId    string `json:"channelId"`
	GenerationId string `json:"generationId"`
	RegionId     string `json:"regionId"`
	PlanTypeId   string `json:"planTypeId"`
}

func NewClusterCreationParams(clusterName string, channelId string,
	generationId string, regionId string,
	planTypeId string) ClusterCreationParams {

	clusterCreationParams := ClusterCreationParams{}
	clusterCreationParams.ClusterName = clusterName
	clusterCreationParams.ChannelId = channelId
	clusterCreationParams.GenerationId = generationId
	clusterCreationParams.RegionId = regionId
	clusterCreationParams.PlanTypeId = planTypeId

	return clusterCreationParams
}

type AuthRequestPayload struct {
	GrantType    string `json:"grant_type"`
	Audience     string `json:"audience"`
	ClientId     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
}

func NewAuthRequestPayload(clientId string, clientSecret string) AuthRequestPayload {
	authRequestPayload := AuthRequestPayload{}
	authRequestPayload.GrantType = "client_credentials"
	authRequestPayload.Audience = "api.cloud.camunda.io"
	authRequestPayload.ClientId = clientId
	authRequestPayload.ClientSecret = clientSecret
	return authRequestPayload
}

type AuthResponsePayload struct {
	AccessToken string `json:"access_token"`
}

type ClusterParams struct {
	Channels         []Channel          `json:"channels"`
	ClusterPlanTypes []ClusterPlantType `json:"clusterPlanTypes"`
	Regions          []Region           `json:"regions"`
}

type Channel struct {
	Id                string       `json:"uuid"`
	Name              string       `json:"name"`
	AllowedGeneration []Generation `json:"allowedGenerations"`
	IsDefault         bool         `json:"isDefault"`
	DefaultGeneration Generation   `json:"defaultGeneration"`
}

type Generation struct {
	Id   string `json:"uuid"`
	Name string `json:"name"`
}

type ClusterPlantType struct {
	Id   string `json:"uuid"`
	Name string `json:"name"`
}

type Region struct {
	Id     string `json:"uuid"`
	Name   string `json:"name"`
	Region string `json:"region"`
	Zone   string `json:"zone"`
}

type ClusterCreatedResponse struct {
	ClusterId string `json:"clusterId"`
}

var authResponsePayload AuthResponsePayload

var clusterParams ClusterParams

var clusterCreatedResponse ClusterCreatedResponse

func getDefaultClusterChannel() Channel {
	var selectedChannel = Channel{}

	for _, c := range clusterParams.Channels {
		if c.IsDefault {
			selectedChannel = c
		}
	}
	return selectedChannel
}

func getDevelopmentClusterPlan() ClusterPlantType {
	var developmentClusterPlanType = ClusterPlantType{}
	for _, cp := range clusterParams.ClusterPlanTypes {
		if cp.Name == "Development" {
			developmentClusterPlanType = cp
		}
	}

	return developmentClusterPlanType
}

func getDefaultRegion() Region {
	var defaultRegion = Region{}
	for _, r := range clusterParams.Regions {
		if r.Name == "Europe West 1D" {
			defaultRegion = r
		}

	}
	return defaultRegion
}

func getClusterParams() {

	req, _ := http.NewRequest("GET", "https://api.cloud.camunda.io/clusters/parameters", nil)
	req.Header.Set("Authorization", "Bearer "+authResponsePayload.AccessToken)

	//fmt.Println("Request cluster params:", req)

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		log.Fatalf("failed to create client cluster params, %v", err)
	}
	//fmt.Println("response Status cluster params:", resp.Status)
	//fmt.Println("response Headers cluster params:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	//fmt.Println("response Body cluster params :", string(body))
	err2 := json.Unmarshal(body, &clusterParams)
	if err2 != nil {
		log.Fatalf("failed to parse body cluster params, %v", err2)
	}
	//marshal, _ := json.Marshal(clusterParams)
	//
	//fmt.Println("parsed: ", string(marshal))
}

func createCluster(clusterName string) string {
	fmt.Println("Creating Cluster Creation Params")
	var channel = getDefaultClusterChannel()
	var clusterPlan = getDevelopmentClusterPlan()
	var region = getDefaultRegion()
	var jsonStr, _ = json.Marshal(NewClusterCreationParams(clusterName,
		channel.Id,
		channel.DefaultGeneration.Id,
		region.Id,
		clusterPlan.Id))

	req, err := http.NewRequest("POST", "https://api.cloud.camunda.io/clusters/", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+authResponsePayload.AccessToken)

	//fmt.Println("Request create cluster :", req)

	client := &http.Client{}
	resp, err := client.Do(req)

	defer resp.Body.Close()
	//fmt.Println("\n\n\nCreate Cluster Response Status:", resp.Status)
	//fmt.Println("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	//fmt.Println("response Body:", string(body))
	if err != nil {
		log.Fatalf("failed to create client, %v", err)
	}

	err2 := json.Unmarshal(body, &clusterCreatedResponse)

	if err2 != nil {
		log.Fatalf("failed to parse body for login, %v", err2)
	}

	return clusterCreatedResponse.ClusterId
}

func login(clientId string, clientSecret string) bool {

	jsonStr, _ := json.Marshal(NewAuthRequestPayload(clientId, clientSecret))

	req, err := http.NewRequest("POST", "https://login.cloud.camunda.io/oauth/token", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	//fmt.Println("Request :", req)
	client := &http.Client{}
	resp, err := client.Do(req)

	defer resp.Body.Close()

	if err != nil {
		log.Fatalf("failed to create client for login, %v", err)
	}
	body, _ := ioutil.ReadAll(resp.Body)
	//fmt.Println("response Status:", resp.Status)
	//fmt.Println("response Body:", string(body))
	if resp.StatusCode == 200 {
		err2 := json.Unmarshal(body, &authResponsePayload)
		log.Printf("json from login parsed!")
		if err2 != nil {
			log.Fatalf("failed to parse body for login, %v", err2)
		}
		return true
	} else {
		log.Fatalf("HTTP Error trying to login, %v", resp.StatusCode)
		return false
	}
	return false
}

func main() {
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}
