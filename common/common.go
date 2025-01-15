package common

import (
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
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