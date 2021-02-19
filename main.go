package main

import (
	"fmt"
	"github.com/camunda-community-hub/camunda-cloud-go-client/cmd"
	"github.com/camunda-community-hub/camunda-cloud-go-client/pkg/cc/client"
	"os"
)

var ClientId = os.Getenv("CC_CLIENT_ID")
var ClientSecret = os.Getenv("CC_CLIENT_SECRET")

func main() {
	login, err := client.Login(ClientId, ClientSecret)

	if err != nil || !login {

		fmt.Errorf("Error trying to Login to Camunda Cloud, " +
			"please check your CC_CLIENT_ID and CC_CLIENT_SECRET! \n %s", err)
		os.Exit(1)

	}
	cmd.Execute()
}
