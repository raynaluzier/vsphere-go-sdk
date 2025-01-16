package common

import (
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	artcommon "github.com/raynaluzier/artifactory-go-sdk/common"
	"github.com/raynaluzier/vsphere-go-sdk/vm"
)

func VcenterAuth(user, pass, server string) string {
	server = AddUrlProtocol(server)  // checks for https:// and adds if missing
	requestPath := server + "/api/session"
	request, err := http.NewRequest("POST", requestPath, nil)
	request.SetBasicAuth(user, pass)
	
	if err != nil {
		strErr := fmt.Sprintf("%v\n", err)
		fmt.Println(strErr)
	}
	
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
	response, err := client.Do(request)
	if err != nil {
		strErr := fmt.Sprintf("%v\n", err)
		fmt.Println(strErr)
	}
	
	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		strErr := fmt.Sprintf("%v\n", err)
		fmt.Println(strErr)
	}
	
	token := string(body)
	return token
}

func AddUrlProtocol(server string) string {
	var serverUrl string
	server = strings.ToLower(server)

	matched, err := regexp.MatchString("https", server)
	if err != nil {
		fmt.Println("Error searching for string - ", err)
	}

	if matched == false {
		serverUrl = "https://" + server
		return serverUrl
	} else {
		return server
	}
}

func TrimUrlProtocol(serverUrl string) string {
	// Sets server URL to lowercase first then searches
	var server string
	serverUrl = strings.ToLower(serverUrl)

	sMatched, err := regexp.MatchString("https", serverUrl)
	if err != nil {
		fmt.Println("Error searching for string - ", err)
	}

	if sMatched == true {
		server = strings.TrimPrefix(serverUrl, "https://")
		return server
	} else {
		matched, err := regexp.MatchString("http", serverUrl)
		if err != nil {
			fmt.Println("Error searching for string - ", err)
			return ""
		}

		if matched == true {
			server = strings.TrimPrefix(serverUrl, "http://")
			return server
		} else {
			fmt.Println("URL protocol http/https not found. Continuing...")
			return serverUrl // leaves as is
		}
	}
}

func TrimQuotes(s string) string {
    if len(s) >= 2 {
        if s[0] == '"' && s[len(s)-1] == '"' {
            return s[1 : len(s)-1]
        }
    }
    return s
}

func RenameFile(oldFilePath, newFilePath string) string {
	// Full path to files
	err := os.Rename(oldFilePath, newFilePath)
    if err != nil {
        log.Fatal(err)
		return "Failed"
    } else {
		return "Success"
	}
}

func GetFileType(filePath string) string {
	filePath = strings.ToLower(filePath)

	ext := filepath.Ext(filePath)
	return ext
}

func CheckFileConvert(outputDir, downloadUri string) string {
	// Takes output directory and download URI, parses the image name from the download URI and determines 
	// the source file path
	// File type is checked; if OVA/OVF, it's converted to VMX. If VMTX, it's converted to VMX
	var result string
	var sourcePath, newPath string
	fileName := artcommon.ParseUriForFilename(downloadUri)
	imageName := artcommon.ParseFilenameForImageName(fileName)
	sourcePath = outputDir + "/" + fileName // CHECK FOR SLASH IN OUTPUT DIR
	newPath    = outputDir + "/" + imageName + ".vmx"

	fileType := GetFileType(fileName)

	switch fileType {
	case "ova":
		fmt.Println("File type found: " + fileType + "; converting to vmx...")
		result = vm.ConvertOvfaToVmx(sourcePath, newPath)
	case "ovf":
		fmt.Println("File type found: " + fileType + "; converting to vmx...")
		result = vm.ConvertOvfaToVmx(sourcePath, newPath)
	case "vmtx":
		fmt.Println("File type found: " + fileType + "; converting to vmx...")
		result = RenameFile(sourcePath, newPath)
	case "vmx":
		fmt.Println("File is already in needed format: vmx.")
		result = "Success"
	default:
		log.Fatal("Found unsupported file type: " + fileType)
		log.Fatal("Supported file types are: ova, ovf, vmtx, and vmx")
		result = "Failed"
	}
	return result
}