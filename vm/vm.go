package vm

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
)

func RegisterVm(token, server, dcName, dsName, imageName string) string {
	var statusCode string
	requestPath := "https://" + server + "/api/vcenter/vm?action=register"

	type Placement struct {
		Folder       string `json:"folder"`
		ResourcePool string `json:"resource_pool"`
	}

	type Payload struct {
		DatastorePath string    `json:"datastore_path"`
		Name          string    `json:"name"`
		Placement     Placement `json:"placement"`
	}

	data := Payload{							         // Update
		DatastorePath: "["+ dsName + "] "+ imageName + "/" + imageName + ".vmx",
		Name: imageName,
		Placement: Placement{
			Folder: "group-v4", 
			ResourcePool: "resgroup-9",
		},
	}
	payloadBytes, err := json.Marshal(data)

	if err != nil {
		fmt.Println("Error1")							// Update
	}
	body := bytes.NewReader(payloadBytes)

	req, err := http.NewRequest(http.MethodPost, requestPath, body)
	if err != nil {
		fmt.Println("Error2")							// Update
	}
	req.Header.Set("Content-Type", "application/json")
	newToken := common.trimQuotes(token)
	req.Header.Set("vmware-api-session-id", newToken)
	
	v1 := req.Header.Get("vmware-api-session-id")
	fmt.Println(v1)


	defaultTransport := http.DefaultTransport.(*http.Transport)
	customTransport := &http.Transport{
		Proxy:					defaultTransport.Proxy,
		DialContext:			defaultTransport.DialContext,
		MaxIdleConns:   		defaultTransport.MaxIdleConns,
		IdleConnTimeout: 		defaultTransport.IdleConnTimeout,
		ExpectContinueTimeout: 	defaultTransport.ExpectContinueTimeout,
		TLSHandshakeTimeout: 	defaultTransport.TLSHandshakeTimeout,
		TLSClientConfig: 		&tls.Config{InsecureSkipVerify: true},
	}

	client := &http.Client{Transport: customTransport}
	resp, err := client.Do(req)
	fmt.Println(resp)

	if err != nil {
		fmt.Println("Error3")
	}
	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		statusCode = "200"
	} else {
		statusCode = "400"
	}
	fmt.Println(statusCode)
	return statusCode
}