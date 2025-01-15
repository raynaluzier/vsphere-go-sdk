package vm

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"runtime"

	"github.com/raynaluzier/vsphere-go-sdk/common"
)

func RegisterVm(token, server, dcName, dsName, imageName, folderId string) string {
	var statusCode string
	requestPath := "https://" + server + "/api/vcenter/vm?action=register"

	type Placement struct {
		Folder       string `json:"folder"`         // required
		ResourcePool string `json:"resource_pool"`	// required
	}

	type Payload struct {
		DatastorePath string    `json:"datastore_path"`
		Name          string    `json:"name"`
		Placement     Placement `json:"placement"`
	}

	// if trying to use vmtx, get error "A specified parameter was not correct: path"
	// if trying to use vmdk, import is successful, but there's no network, etc.
	data := Payload{                                 // Update
		DatastorePath: "[Work2] ub20pkrt-10031746/ub20pkrt-10031746_2-flat.vmx",  //orig, vmx
		Name: "ub20pkrt-10031746",
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
	newToken := common.TrimQuotes(token)
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

	if resp.StatusCode == 200 {   // print actual status code
		statusCode = "200"
	} else {
		statusCode = "400"
	}
	fmt.Println(statusCode)
	return statusCode
}

func ConvertOvfToVmx() string {
	var cmd *exec.Cmd

	switch runtime.GOOS{
	case "windows":  // if you don't put the cmd and /c, it will fail
		cmd = exec.Command("cmd", "/c", "ovftool --help")
	default: // mac & linux
		cmd = exec.Command("ifconfig -a")   // mac "bash"
	}

	cmd.Stdout = os.Stdout

	if err := cmd.Run(); err != nil {
		fmt.Println("could not run command")
		return "Failed"
	} else {
		return "Success"
	}

}