package vm

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"strings"

	"github.com/raynaluzier/vsphere-go-sdk/common"
)

func RegisterVm(token, vcServer, dcName, dsName, imageName, folderId string) string {
	var statusCode string
	requestPath := "https://" + vcServer + "/api/vcenter/vm?action=register"

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
		DatastorePath: "["+ dsName + "] "+ imageName + "/" + imageName + ".vmx",
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


// These OVF/OVA conversion functions require the OVFTool be installed
	// For Windows: 'cmd' and '/c' are not included in the commands, this will fail
	// Input can be local or URL
	// Ensure local paths are escaped properly (ex: C:\\lab\\file.vmx)
func ConvertOvfaToVmx(inputPath, outputPath string) string {
	var cmd *exec.Cmd
	var ovaMatched, ovfMatched, vmxMatched bool

	inputPath = strings.ToLower(inputPath)
	outputPath = strings.ToLower(outputPath)
	
	// Ensure input is either an OVA or OVF file
	ovfMatched, err := regexp.MatchString("ovf", inputPath)
	if ovfMatched == false {
		ovaMatched, err = regexp.MatchString("ova", inputPath)
		if ovaMatched == false {
			fmt.Println("Error: Input is neither an OVA or OVF file.")
			fmt.Println("Please provide the full file path to the OVA or OVF that's to be converted.") 
		}
	}
	if err != nil {
		fmt.Println("Unable to search for OVA/OVF string.")
	}
	
	// Ensure output is a VMX file
	vmxMatched, err = regexp.MatchString("vmx", outputPath)
	if vmxMatched == false {
		fmt.Println("Error: Output does not include a VMX file path.")
		fmt.Println("Please provide the full destination file path to the resulting VMX file.") 	
	}
	if err != nil {
		fmt.Println("Unable to search for VMX string.")
	}
	
	ovfCmd := "ovftool " + inputPath + " " + outputPath

	fmt.Println("Beginning conversion process... This could take a while.")
	switch runtime.GOOS{
	case "windows":
		cmd = exec.Command("cmd", "/c", ovfCmd)
	default: // mac & linux
		cmd = exec.Command(ovfCmd)   // mac "bash"
	}

	cmd.Stdout = os.Stdout

	if err := cmd.Run(); err != nil {
		fmt.Println("Could not run the exec shell command.")
		// if mac or linux, do we need to prefice cmd with "bash"?
		return "Failed"
	} else {
		return "Success"
	}
}

func ConvertVmxToOvfa(inputPath, outputPath string) string {
	var cmd *exec.Cmd
	var ovaMatched, ovfMatched, vmxMatched bool

	inputPath = strings.ToLower(inputPath)
	outputPath = strings.ToLower(inputPath)
	
	// Ensure input is a VMX file
	vmxMatched, err := regexp.MatchString("vmx", inputPath)
	if vmxMatched == false {
		fmt.Println("Error: Input does not include a VMX file path.")
		fmt.Println("Please provide the full file path to the VMX that's to be converted.") 	
	}
	if err != nil {
		fmt.Println("Unable to search for VMX string.")
	}
	
	// Ensure output is either an OVA or OVF file
	ovfMatched, err = regexp.MatchString("ovf", outputPath)
	if ovfMatched == false {
		ovaMatched, err = regexp.MatchString("ova", outputPath)
		if ovaMatched == false {
			fmt.Println("Error: Output is neither an OVA or OVF file.")
			fmt.Println("Please provide the full destination file path to the resulting OVA or OVF file.")
		}
	}
	if err != nil {
		fmt.Println("Unable to search for OVA/OVF string.")
	}
	
	ovfCmd := "ovftool " + inputPath + " " + outputPath

	fmt.Println("Beginning conversion process... This could take a while.")
	switch runtime.GOOS{
	case "windows":
		cmd = exec.Command("cmd", "/c", ovfCmd)
	default: // mac & linux
		cmd = exec.Command(ovfCmd)   // mac "bash"
	}

	cmd.Stdout = os.Stdout

	if err := cmd.Run(); err != nil {
		fmt.Println("Could not run the exec shell command.") 
		// if mac or linux, do we need to prefice cmd with "bash"?
		return "Failed"
	} else {
		return "Success"
	}
}