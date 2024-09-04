package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/aleroxac/docker-healthcheck/config"
	"github.com/aleroxac/docker-healthcheck/internal/infra/client_http"
)

func main() {
	config, err := config.Init()
	if err != nil {
		log.Fatalf("Fail to initialize config: %v", err)
	}

	client := client_http.NewClient(5)
	request := client_http.ClientRequest{
		Method:   "GET",
		Protocol: config.Protocol,
		Host:     config.Host,
		Port:     config.Port,
		Path:     config.Path,
	}
	response, err := client.Request(&request)
	if err != nil {
		log.Fatalf("Fail to make the request: %v", err)
	}

	res_json, err := json.Marshal(response)
	if err != nil {
		log.Fatalf("Fail to marshall response: %v", err)
	}

	var client_response client_http.ClientResponse
	err = json.Unmarshal(res_json, &client_response)
	if err != nil {
		log.Fatalf("Fail to unmarshall reponse: %v", err)
	}
	fmt.Println(string(res_json))
}
