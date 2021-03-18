package client

type ClusterStatusResponse struct {
	ClusterId     string        `json:"uuid"`
	ClusterStatus ClusterStatus `json:"status"`
}

type ClusterStatus struct {
	OperateStatus  string `json:"operateStatus"`
	OperateURL     string `json:"operateUrl"`
	Ready          string `json:"ready"`
	ZeebeStatus    string `json:"zeebeStatus"`
	ZeebeURL       string `json:"zeebeUrl"`
	TaskListStatus string `json:"tasklistStatus"`
	TaskListURL    string `json:"tasklistUrl"`
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

type Cluster struct {
	ID              string          `json:"uuid"`
	Name            string          `json:"name"`
	Channel         Channel         `json:"channel"`
	Generation      Generation      `json:"generation"`
	Created         string          `json:"created"`
	K8sContext      K8sContext      `json:"k8sContext"`
	ClusterMetadata ClusterMetadata `json:"metadata"`
}

type K8sContext struct {
	ID     string `json:"uuid"`
	Name   string `json:"name"`
	Region string `json:"region"`
	Zone   string `json:"zone"`
}

type ClusterMetadata struct {
	ID                string `json:"uid"`
	CreationTimestamp string `json:"creationTimestamp"`
	Generation        int    `json:"generation"`
	Name              string `json:"name"`
	ResourceVersion   string `json:"resourceVersion"`
	SelfLink          string `json:"selfLink"`
}

type ZeebeClientResponse struct {
	ClientID    string   `json:"clientId"`
	Created     string   `json:"created"`
	CreatedBy   string   `json:"createdBy"`
	UUID        string   `json:"uuid"`
	Name        string   `json:"name"`
	Internal    bool     `json:"internal"`
	Permissions []string `json:"permissions"`
}

type ZeebeClientDetailsResponse struct {
	Name                        string `json:"name"`
	ZEEBEADDRESS                string `json:"ZEEBE_ADDRESS"`
	ZEEBECLIENTID               string `json:"ZEEBE_CLIENT_ID"`
	ZEEBEAUTHORIZATIONSERVERURL string `json:"ZEEBE_AUTHORIZATION_SERVER_URL"`
}

type ZeebeClientCreatedResponse struct {
	Name         string `json:"name"`
	ClientID     string `json:"clientId"`
	ClientSecret string `json:"clientSecret"`
}

type ZeebeClientCreatePayload struct {
	ClientName string `json:"clientName"`
}
