<img src="https://img.shields.io/badge/Community%20Extension-An%20open%20source%20community%20maintained%20project-FF4700"> <img src="https://img.shields.io/badge/Lifecycle-Incubating-blue">

# Camunda Cloud CLI and Go Client Library
Camunda Cloud CLI to interact with Camunda Cloud Resources.

This repository contains a CLI tool to interact with your Camunda Cloud account via the command-line. 

This CLI interacts with the Camudna Cloud Management REST API.

This repository also contains a Go Library client to consume the Camunda Cloud Management REST APIs from other Go programs. 

## Usage

You can create a Camunda Cloud Account here: https://accounts.cloud.camunda.io/signup
  
Then you need to go to `<YOUR USER>(Top right corner) -> Organization Settings -> Cloud Management API` 
and Create a new client. 
You need to copy and save the Client Id and the Secret Id from that client. 

You need to export the following variables for this command to interact with a Camunda Cloud Account:
  - export CC_CLIENT_ID=`<YOUR CLIENT ID>`
  - export CC_CLIENT_SECRET=`<YOUR CLIENT SECRET>`
  
  Available Commands:  
  
  **List all clusters**
  `cc-ctl clusters get --all`

  **Get cluster from id**
  `cc-ctl clusters get --id <cluster_id>`

  **Get cluster from name**
  `cc-ctl clusters get --name <cluster_name>`

  **Delete cluster from id**
  `cc-ctl clusters delete --id <cluster_id>`

  **Delete cluster from name**
  `cc-ctl clusters delete --name <cluster_name>`

  **Create cluster from default configuration**
  `cc-ctl clusters create --default --name <cluster_name>`

# Feedback / Contribute back

This is a super simple project for you to contribute. Feel free to create issues or send PRs with improvements. 

Consumers: 
- [Zeebe Kubernetes Operator V3](https://github.com/salaboy/zeebe-operator-cc)

