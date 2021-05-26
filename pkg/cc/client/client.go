package client

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/trace/jaeger"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/semconv"
	"go.opentelemetry.io/otel/trace"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

type CCClient struct {
	AuthResponsePayload AuthResponsePayload

	ClusterParams ClusterParams

	ClusterCreatedResponse ClusterCreatedResponse

	ClusterStatusResponse ClusterStatusResponse

	ZeebeClientCreate ZeebeClientCreatedResponse

	tracer trace.Tracer

	tracingEnabled bool

	tracerURL string

	ccApiURL string
}

func (c *CCClient) TracingEnabled(tracingEnabled bool ){
	c.tracingEnabled = tracingEnabled
}

func (c *CCClient) SetCCApiURL(ccApiURL string ){
	c.ccApiURL = ccApiURL
}

func (c *CCClient) SetTracerURL(tracerURL string ){
	c.tracerURL = tracerURL
}

func (c *CCClient) InitTracer() func() {

	// Create and install Jaeger export pipeline.
	flush, err := jaeger.InstallNewPipeline(
		jaeger.WithCollectorEndpoint("http://" + c.tracerURL + "/api/traces"),
		jaeger.WithSDKOptions(
			sdktrace.WithSampler(sdktrace.AlwaysSample()),
			sdktrace.WithResource(resource.NewWithAttributes(
				semconv.ServiceNameKey.String("cc-ctl"),
				attribute.String("exporter", "jaeger"),
				attribute.Float64("float", 312.23),
			)),
		),
	)
	if err != nil {
		log.Fatal(err)
	}
	c.tracer = otel.Tracer("cc-ctl")
	return flush
}

func (c *CCClient) getDefaultStableClusterChannel() Channel {
	var selectedChannel = Channel{}

	for _, channel := range c.ClusterParams.Channels {
		if channel.Name == "Stable" {
			selectedChannel = channel
		}
	}
	return selectedChannel
}

func (c *CCClient) getClusterChannelByName(channelName string) Channel {
	var selectedChannel = Channel{}

	for _, channel := range c.ClusterParams.Channels {
		if strings.Contains(channel.Name, channelName) {
			selectedChannel = channel
		}
	}
	return selectedChannel
}

func (c *CCClient) getClusterPlanByName(clusterPlanName string) ClusterPlantType {
	var developmentClusterPlanType = ClusterPlantType{}
	for _, cp := range c.ClusterParams.ClusterPlanTypes {
		if cp.Name == clusterPlanName {
			developmentClusterPlanType = cp
		}
	}

	return developmentClusterPlanType
}

func (c *CCClient) getDevelopmentClusterPlan() ClusterPlantType {
	var developmentClusterPlanType = ClusterPlantType{}
	for _, cp := range c.ClusterParams.ClusterPlanTypes {
		if cp.Name == "Development" {
			developmentClusterPlanType = cp
		}
	}

	return developmentClusterPlanType
}

func (c *CCClient) getDefaultRegion() Region {
	//chose the first one as default
	return c.ClusterParams.Regions[0]
}

func (c *CCClient) GetClusterParams() (*ClusterParams, error) {
	ctx := context.Background()
	return c.GetClusterParamsWithContext(ctx)
}

func (c *CCClient) GetClusterParamsWithContext(ctx context.Context) (*ClusterParams, error) {
	if c.tracingEnabled {
		_, span := c.tracer.Start(ctx, "getClusterParams")
		defer span.End()
	}

	req, _ := http.NewRequest("GET", "https://api."+c.ccApiURL+"/clusters/parameters", nil)
	req.Header.Set("Authorization", "Bearer "+c.AuthResponsePayload.AccessToken)

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		log.Printf("failed to create client cluster params, %v", err)
		return &c.ClusterParams, err
	}

	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response Body Get cluster params :", string(body))
	err2 := json.Unmarshal(body, &c.ClusterParams)

	if err2 != nil {
		log.Printf("failed to parse body cluster params, %v, %s", err2, string(body))
		return &c.ClusterParams, err2
	}

	return &c.ClusterParams, nil
}
func (c *CCClient) GetClusterDetails(clusterId string) (ClusterStatus, error) {
	ctx := context.Background()
	return c.GetClusterDetailsWithContext(ctx, clusterId)
}

func (c *CCClient) GetClusterDetailsWithContext(ctx context.Context, clusterId string) (ClusterStatus, error) {
	if c.tracingEnabled {
		_, span := c.tracer.Start(ctx, "getClusterDetails")
		defer span.End()
	}
	req, _ := http.NewRequest("GET", "https://api."+c.ccApiURL+"/clusters/"+clusterId, nil)
	req.Header.Set("Authorization", "Bearer "+c.AuthResponsePayload.AccessToken)

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
	err2 := json.Unmarshal(body, &c.ClusterStatusResponse)
	if err2 != nil {
		//log.Printf("failed to parse body cluster details, %v,  %s", err2, string(body))
		clusterStatus.Ready = "Not Found"
		return clusterStatus, nil
	}
	clusterStatus = c.ClusterStatusResponse.ClusterStatus
	return clusterStatus, nil

}


func (c *CCClient) CreateClusterCustomConfig(clusterParams ClusterCreationParams) (string, error) {
	ctx := context.Background()
	return c.CreateClusterCustomConfigWithContext(ctx, clusterParams)
}

func (c *CCClient) CreateClusterCustomConfigWithContext(ctx context.Context, clusterParams ClusterCreationParams) (string, error) {

	if c.tracingEnabled {
		_, span := c.tracer.Start(ctx, "createClusterCustomConfig")
		defer span.End()
	}

	_, existsErr := c.clusterExistsValidator(clusterParams.ClusterName)

	if existsErr != nil {
		return "", existsErr
	}

	jsonStr, _ := json.Marshal(clusterParams)

	req, _ := http.NewRequest("POST", "https://api."+c.ccApiURL+"/clusters", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.AuthResponsePayload.AccessToken)

	fmt.Println("Request create cluster :", req)

	client := &http.Client{}
	resp, err := client.Do(req)

	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Printf("failed to create client, %v", err)
		return "", err
	}

	err2 := json.Unmarshal(body, &c.ClusterCreatedResponse)

	if err2 != nil {
		log.Printf("Body to unmarshal: %s", string(body))
		log.Printf("failed to parse body for create cluster, %v", err2)
		return "", err2
	}

	return c.ClusterCreatedResponse.ClusterId, nil
}

func (c *CCClient) CreateClusterWithParams(clusterName string, clusterPlanName string, channelName string, generationName string, clusterRegion string) (string, error) {
	ctx := context.Background()
	return c.CreateClusterWithParamsAndContext(ctx, clusterName, clusterPlanName, channelName, generationName, clusterRegion)
}

func (c *CCClient) CreateClusterWithParamsAndContext(ctx context.Context, clusterName string, clusterPlanName string, channelName string, generationName string, clusterRegion string) (string, error) {
	if c.tracingEnabled {
		_, span := c.tracer.Start(ctx, "createClusterWithParams")
		defer span.End()
	}
	_, existsErr := c.clusterExistsValidator(clusterName)

	if existsErr != nil {
		return "", existsErr
	}

	var channel = Channel{}
	var clusterPlan = ClusterPlantType{}
	var region = Region{}
	var generation = Generation{}

	if clusterRegion != "" {
		region, _ = c.getClusterRegionByName(clusterRegion)
	} else {
		region = c.getDefaultRegion()
	}

	if channelName != "" {
		channel = c.getClusterChannelByName(channelName)
	} else {
		channel = c.getDefaultStableClusterChannel()
	}

	if generationName != "" {
		generation = c.getGenerationByNameForSelectedChannel(channel, clusterPlanName)
	} else {
		generation = channel.DefaultGeneration
	}

	if clusterPlanName != "" {
		clusterPlan = c.getClusterPlanByName(clusterPlanName)
	} else {
		clusterPlan = c.getDevelopmentClusterPlan()
	}

	var jsonStr, _ = json.Marshal(NewClusterCreationParams(clusterName,
		channel.Id,
		generation.Id,
		region.Id,
		clusterPlan.Id))

	req, err := http.NewRequest("POST", "https://api."+c.ccApiURL+"/clusters/", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.AuthResponsePayload.AccessToken)

	fmt.Println("Request create cluster :", req)

	client := &http.Client{}
	resp, err := client.Do(req)

	defer resp.Body.Close()
	fmt.Println("\n\n\nCreate Cluster Response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Printf("failed to create client, %v", err)
		return "", err
	}

	err2 := json.Unmarshal(body, &c.ClusterCreatedResponse)

	if err2 != nil {
		log.Printf("Body to unmarshal: %s ", string(body))
		log.Printf("failed to parse body for create cluster, %v", err2)
		return "", err2
	}

	return c.ClusterCreatedResponse.ClusterId, nil
}

func (c *CCClient) getGenerationByNameForSelectedChannel(channel Channel, generationName string) Generation {
	var generataion = Generation{}
	for _, ag := range channel.AllowedGeneration {
		if ag.Name == generationName {
			generataion = ag
		}
	}
	return generataion
}

func (c *CCClient) getClusterRegionByName(regionName string) (Region, error) {
	var selectedRegion = Region{}
	for _, r := range c.ClusterParams.Regions {
		if r.Name == regionName {
			selectedRegion = r
		}

	}
	if selectedRegion.Name == "" {
		return Region{}, errors.New("No Region Found with name: " + regionName)
	}
	return selectedRegion, nil

}

func (c *CCClient) CreateClusterDefault(clusterName string) (string, error) {
	ctx := context.Background()
	return c.CreateClusterDefaultWithContext(ctx, clusterName)
}

func (c *CCClient) CreateClusterDefaultWithContext(ctx context.Context, clusterName string) (string, error) {
	if c.tracingEnabled {
		_, span := c.tracer.Start(ctx, "createClusterDefault")
		defer span.End()
	}
	_, existsErr := c.clusterExistsValidator(clusterName)

	if existsErr != nil {
		return "", existsErr
	}

	var channel = c.getDefaultStableClusterChannel()
	var clusterPlan = c.getDevelopmentClusterPlan()
	var region = c.getDefaultRegion()
	var jsonStr, _ = json.Marshal(NewClusterCreationParams(clusterName,
		channel.Id,
		channel.DefaultGeneration.Id,
		region.Id,
		clusterPlan.Id))

	req, err := http.NewRequest("POST", "https://api."+c.ccApiURL+"/clusters/", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.AuthResponsePayload.AccessToken)

	fmt.Println("Request create cluster :", req)

	client := &http.Client{}
	resp, err := client.Do(req)

	defer resp.Body.Close()
	fmt.Println("\n\n\nCreate Cluster Response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Printf("failed to create client, %v", err)
		return "", err
	}

	err2 := json.Unmarshal(body, &c.ClusterCreatedResponse)

	if err2 != nil {
		log.Printf("Body to unmarshal: %s ", string(body))
		log.Printf("failed to parse body for create cluster, %v", err2)
		return "", err2
	}

	return c.ClusterCreatedResponse.ClusterId, nil
}

func (c *CCClient) Login(clientId string, clientSecret string) (bool, error) {
	ctx := context.Background()
	return c.LoginWithContext(ctx, clientId, clientSecret)
}

func (c *CCClient) LoginWithContext(ctx context.Context, clientId string, clientSecret string) (bool, error) {
	if c.tracingEnabled {
		_, span := c.tracer.Start(ctx, "login")
		defer span.End()
	}
	jsonStr, _ := json.Marshal(NewAuthRequestPayload(clientId, clientSecret))

	req, err := http.NewRequest("POST", "https://login."+c.ccApiURL+"/oauth/token", bytes.NewBuffer(jsonStr))
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
		err2 := json.Unmarshal(body, &c.AuthResponsePayload)
		//		log.Printf("json from login parsed!")
		if err2 != nil {
			log.Printf("failed to parse body for login, %v, %s", err2, string(body))
			return false, err2
		}
		return true, nil
	} else {
		log.Printf("HTTP Error trying to login, %v", resp.StatusCode)
		return false, errors.New(fmt.Sprintf("HTTP Error trying to login: %d", resp.StatusCode))
	}
}

func (c *CCClient) DeleteCluster(clusterId string) (bool, error) {
	ctx := context.Background()
	return c.DeleteClusterWithContext(ctx, clusterId)
}

func (c *CCClient) DeleteClusterWithContext(ctx context.Context, clusterId string) (bool, error) {

	if c.tracingEnabled {
		_, span := c.tracer.Start(ctx, "deleteCluster")
		defer span.End()
	}
	req, _ := http.NewRequest("DELETE", "https://api."+c.ccApiURL+"/clusters/"+clusterId, nil)
	req.Header.Set("Authorization", "Bearer "+c.AuthResponsePayload.AccessToken)

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		log.Printf("failed to create client cluster params, %v", err)
		return false, errors.New(fmt.Sprintf("HTTP Error trying to login: %d", resp.StatusCode))
	}

	if resp.StatusCode == 200 {
		return true, nil
	}
	return false, errors.New(fmt.Sprintf("HTTP Error trying to login: %d", resp.StatusCode))
	//fmt.Println("response Status delete cluster:", resp.Status)
	//fmt.Println("response Headers delete cluster:", resp.Header)
	//body, _ := ioutil.ReadAll(resp.Body)

	//fmt.Println("response Body delete cluster :", string(body))

}


func (c *CCClient) GetClusters() ([]Cluster, error) {
	ctx := context.Background()
	return c.GetClustersWithContext(ctx)
}

// GetClusters from Camunda Cloud
func (c *CCClient) GetClustersWithContext(ctx context.Context) ( []Cluster, error) {

	if c.tracingEnabled {
		_, span := c.tracer.Start(ctx, "getClusters")
		defer span.End()
	}
	data := []Cluster{}

	req, _ := http.NewRequest("GET", "https://api."+c.ccApiURL+"/clusters", nil)

	req.Header.Set("Authorization", "Bearer "+c.AuthResponsePayload.AccessToken)

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


func (c *CCClient) GetClusterByName(name string) (Cluster, error) {
	ctx := context.Background()
	return c.GetClusterByNameWithContext(ctx, name)
}

func (c *CCClient) GetClusterByNameWithContext(ctx context.Context, name string) (Cluster, error) {

	if c.tracingEnabled {
		_, span := c.tracer.Start(ctx, "getClusterByName")
		defer span.End()
	}
	data := Cluster{}

	clusters, err := c.GetClusters()

	if err != nil {
		return data, err
	}

	for _, cluster := range clusters {

		if cluster.Name == name {
			return cluster, nil
		}
	}

	return data, nil
}

func (c *CCClient) clusterExistsValidator(clusterName string) (string, error) {
	ctx := context.Background()
	return c.clusterExistsValidatorWithContext(ctx, clusterName)
}

func (c *CCClient) clusterExistsValidatorWithContext(ctx context.Context, clusterName string) (string, error) {

	if c.tracingEnabled {
		_, span := c.tracer.Start(ctx, "clusterExistsValidator")
		defer span.End()
	}

	cluster, err := c.GetClusterByName(clusterName)

	if err != nil {
		return "", err
	}

	if cluster.ID != "" {
		return "", errors.New("Cluster name already exists on Camunda Cloud")
	}

	return "", nil
}

func (c *CCClient) GetZeebeClients(clusterID string) ([]ZeebeClientResponse, error) {
	ctx := context.Background()
	return c.GetZeebeClientsWithContext(ctx, clusterID)
}
// GetZeebeClients - List all Zeebe clients
func (c *CCClient) GetZeebeClientsWithContext(ctx context.Context, clusterID string) ([]ZeebeClientResponse, error) {

	if c.tracingEnabled {
		_, span := c.tracer.Start(ctx, "getZeebeClients")
		defer span.End()
	}

	data := []ZeebeClientResponse{}

	if len(clusterID) == 0 {
		return data, NewError("Cluster id should not be empty")
	}

	req, _ := http.NewRequest("GET", "https://api."+c.ccApiURL+"/clusters/"+clusterID+"/clients", nil)

	req.Header.Set("Authorization", "Bearer "+c.AuthResponsePayload.AccessToken)

	client := &http.Client{}

	resp, err := client.Do(req)

	if err != nil {
		log.Printf("Failed to get zeebe clients")
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

func (c *CCClient) GetZeebeClientDetails(clusterID string, clientID string) (ZeebeClientDetailsResponse, error) {
	ctx := context.Background()
	return c.GetZeebeClientDetailsWithContext(ctx, clusterID, clientID)
}

func (c *CCClient) GetZeebeClientDetailsWithContext(ctx context.Context, clusterID string, clientID string) (ZeebeClientDetailsResponse, error) {

	if c.tracingEnabled {
		_, span := c.tracer.Start(ctx, "getZeebeClientDetails")
		defer span.End()
	}

	data := ZeebeClientDetailsResponse{}

	if len(clusterID) == 0 {
		return data, NewError("Cluster id should not be empty")
	}

	if len(clientID) == 0 {
		return data, NewError("Client id should not be empty")
	}

	req, _ := http.NewRequest("GET", "https://api."+c.ccApiURL+"/clusters/"+clusterID+"/clients/"+clientID, nil)

	req.Header.Set("Authorization", "Bearer "+c.AuthResponsePayload.AccessToken)

	client := &http.Client{}

	resp, err := client.Do(req)

	if err != nil {
		log.Printf("Failed to get zeebe details")
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

func (c *CCClient) CreateZeebeClient(clusterID string, clientName string) (ZeebeClientCreatedResponse, error) {
	ctx := context.Background()
	return c.CreateZeebeClientWithContext(ctx, clusterID, clientName)
}

func (c *CCClient) CreateZeebeClientWithContext(ctx context.Context, clusterID string, clientName string) (ZeebeClientCreatedResponse, error) {

	if c.tracingEnabled {
		_, span := c.tracer.Start(ctx, "createZeebeClient")
		defer span.End()
	}

	zeebeClient := ZeebeClientCreatePayload{
		ClientName: clientName,
	}

	if len(clusterID) == 0 {
		return ZeebeClientCreatedResponse{}, NewError("Cluster id should not be empty")
	}

	if len(clientName) == 0 {
		return ZeebeClientCreatedResponse{}, NewError("Client name should not be empty")
	}

	jsonStr, _ := json.Marshal(zeebeClient)

	req, _ := http.NewRequest("POST", "https://api."+c.ccApiURL+"/clusters/"+clusterID+"/clients", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.AuthResponsePayload.AccessToken)

	client := &http.Client{}
	resp, err := client.Do(req)

	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Printf("failed to create zeebe client, %v", err)
		return ZeebeClientCreatedResponse{}, err
	}

	err2 := json.Unmarshal(body, &c.ZeebeClientCreate)

	if err2 != nil {
		log.Printf("Body to unmarshal: %s", string(body))
		log.Printf("failed to parse body for zeebe client, %v", err2)
		return ZeebeClientCreatedResponse{}, err2
	}

	return c.ZeebeClientCreate, nil
}

func (c *CCClient) DeleteZeebeClient(clusterID string, clientID string) (bool, error) {
	ctx := context.Background()
	return c.DeleteZeebeClientWithContext(ctx, clusterID, clientID)
}

func (c *CCClient) DeleteZeebeClientWithContext(ctx context.Context, clusterID string, clientID string) (bool, error) {

	if c.tracingEnabled {
		_, span := c.tracer.Start(ctx, "deleteZeebeClient")
		defer span.End()
	}

	if len(clusterID) == 0 {
		return false, NewError("Cluster id should not be empty")
	}

	req, _ := http.NewRequest("DELETE", "https://api."+c.ccApiURL+"/clusters/"+clusterID+"/clients/"+clientID, nil)
	req.Header.Set("Authorization", "Bearer "+c.AuthResponsePayload.AccessToken)

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		log.Printf("Failed to delete zeebe client, %v", err)
		return false, errors.New(fmt.Sprintf("HTTP Error trying to delete zeebe client: %d", resp.StatusCode))
	}

	if resp.StatusCode == 200 {
		return true, nil
	}
	return false, errors.New(fmt.Sprintf("HTTP Error trying to delete zeebe client: %d", resp.StatusCode))

}
