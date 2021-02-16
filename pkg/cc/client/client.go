package client

import (
	"bytes"
	"encoding/json"
	"errors"
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

func GetClusterParams() error {

	req, _ := http.NewRequest("GET", "https://api.cloud.camunda.io/clusters/parameters", nil)
	req.Header.Set("Authorization", "Bearer "+authResponsePayload.AccessToken)

	//fmt.Println("Request cluster params:", req)

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		log.Printf("failed to create client cluster params, %v", err)
		return err
	}
	//fmt.Println("response Status cluster params:", resp.Status)
	//fmt.Println("response Headers cluster params:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	//fmt.Println("response Body cluster params :", string(body))
	err2 := json.Unmarshal(body, &clusterParams)
	if err2 != nil {
		log.Printf("failed to parse body cluster params, %v, %s", err2, string(body))
		return err2
	}
	//marshal, _ := json.Marshal(clusterParams)
	//
	//fmt.Println("parsed: ", string(marshal))
	return nil
}

func GetClusterDetails(clusterId string) (ClusterStatus, error) {
	req, _ := http.NewRequest("GET", "https://api.cloud.camunda.io/clusters/"+clusterId, nil)
	req.Header.Set("Authorization", "Bearer "+authResponsePayload.AccessToken)

	var clusterStatus = ClusterStatus{}
	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		log.Printf("failed to create client cluster details, %v", err)
		return clusterStatus, err
	}

	//fmt.Println("response Status cluster params:", resp.Status)
	//fmt.Println("response Headers cluster params:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	//fmt.Println("response Body cluster params :", string(body))
	err2 := json.Unmarshal(body, &clusterStatusResponse)
	if err2 != nil {
		log.Printf("failed to parse body cluster details, %v,  %s", err2, string(body))
		clusterStatus.Ready = "Not Found"
		return clusterStatus, nil
	}
	clusterStatus = clusterStatusResponse.ClusterStatus
	return clusterStatus, nil

}

func CreateCluster(clusterName string) (string, error) {
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

	if err != nil {
		log.Printf("failed to create client, %v", err)
		return "", err
	}

	err2 := json.Unmarshal(body, &clusterCreatedResponse)

	if err2 != nil {
		log.Printf("Body to unmarshal: ", string(body))
		log.Printf("failed to parse body for create cluster, %v", err2)
		return "", err2
	}

	return clusterCreatedResponse.ClusterId, nil
}

func Login(clientId string, clientSecret string) (bool, error) {

	jsonStr, _ := json.Marshal(NewAuthRequestPayload(clientId, clientSecret))

	req, err := http.NewRequest("POST", "https://login.cloud.camunda.io/oauth/token", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	//fmt.Println("Request :", req)
	client := &http.Client{}
	resp, err := client.Do(req)

	defer resp.Body.Close()

	if err != nil {
		log.Printf("failed to create client for login, %v", err)
		return false, err
	}
	body, _ := ioutil.ReadAll(resp.Body)
	//fmt.Println("response Status:", resp.Status)
	//fmt.Println("response Body:", string(body))
	if resp.StatusCode == 200 {
		err2 := json.Unmarshal(body, &authResponsePayload)
//		log.Printf("json from login parsed!")
		if err2 != nil {
			log.Printf("failed to parse body for login, %v, %s", err2, string(body))
			return false, err2
		}
		return true, nil
	} else {
		log.Printf("HTTP Error trying to login, %v", resp.StatusCode)
		return false, errors.New(fmt.Sprintf("HTTP Error trying to login: %i", resp.StatusCode))
	}
}

func DeleteCluster(clusterId string) (bool, error) {
	req, _ := http.NewRequest("DELETE", "https://api.cloud.camunda.io/clusters/"+clusterId, nil)
	req.Header.Set("Authorization", "Bearer "+authResponsePayload.AccessToken)

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		log.Printf("failed to create client cluster params, %v", err)
		return false, errors.New(fmt.Sprintf("HTTP Error trying to login: %i", resp.StatusCode))
	}

	if resp.StatusCode == 200 {
		return true, nil
	}
	return false, errors.New(fmt.Sprintf("HTTP Error trying to login: %i", resp.StatusCode))
	//fmt.Println("response Status delete cluster:", resp.Status)
	//fmt.Println("response Headers delete cluster:", resp.Header)
	//body, _ := ioutil.ReadAll(resp.Body)

	//fmt.Println("response Body delete cluster :", string(body))

}

// GetClusters from Camunda Cloud
func GetClusters() ([]Cluster, error) {

	data := []Cluster{}

	req, _ := http.NewRequest("GET", "https://api.cloud.camunda.io/clusters", nil)

	req.Header.Set("Authorization", "Bearer "+authResponsePayload.AccessToken)

	client := &http.Client{}

	resp, err := client.Do(req)

	if err != nil {
		log.Printf("Failed to get all clusters")
		return data, err
	}

	body, _ := ioutil.ReadAll(resp.Body)

	err2 := json.Unmarshal(body, &data)

	if err2 != nil {
		log.Printf("Failed to unmarshal response body ->  %s", string(body))
		return data, err2
	}

	return data, nil
}
