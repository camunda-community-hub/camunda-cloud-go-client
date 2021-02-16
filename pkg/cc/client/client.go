package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

var authResponsePayload AuthResponsePayload

var clusterParams ClusterParams

var clusterCreatedResponse ClusterCreatedResponse

var clusterStatusResponse ClusterStatusResponse

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

func GetClusterParams() {

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

func GetClusterDetails(clusterId string) ClusterStatus {
	req, _ := http.NewRequest("GET", "https://api.cloud.camunda.io/clusters/"+clusterId, nil)
	req.Header.Set("Authorization", "Bearer "+authResponsePayload.AccessToken)

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		log.Fatalf("failed to create client cluster details, %v", err)
	}

	//fmt.Println("response Status cluster params:", resp.Status)
	//fmt.Println("response Headers cluster params:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	//fmt.Println("response Body cluster params :", string(body))
	err2 := json.Unmarshal(body, &clusterStatusResponse)
	if err2 != nil {
		log.Fatalf("failed to parse body cluster details, %v", err2)
	}

	return clusterStatusResponse.ClusterStatus

}

func CreateCluster(clusterName string) string {
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

func Login(clientId string, clientSecret string) bool {

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

func DeleteCluster(clusterId string) bool {
	req, _ := http.NewRequest("DELETE", "https://api.cloud.camunda.io/clusters/"+clusterId, nil)
	req.Header.Set("Authorization", "Bearer "+authResponsePayload.AccessToken)

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		log.Fatalf("failed to create client cluster params, %v", err)
	}

	if resp.StatusCode == 200 {
		return true
	} else {
		return false
	}
	//fmt.Println("response Status delete cluster:", resp.Status)
	//fmt.Println("response Headers delete cluster:", resp.Header)
	//body, _ := ioutil.ReadAll(resp.Body)

	//fmt.Println("response Body delete cluster :", string(body))

}

func GetClusters() {
	req, _ := http.NewRequest("GET", "https://api.cloud.camunda.io/clusters", nil)
	req.Header.Set("Authorization", "Bearer "+authResponsePayload.AccessToken)
	client := &http.Client{}

	resp, err := client.Do(req)

	defer resp.Request.GetBody()

	if err != nil {
		log.Fatalf("failed to create client for get clusters, %v", err)
	}

	body, _ := ioutil.ReadAll(resp.Body)

	fmt.Println("response body clusters:", string(body))
}
