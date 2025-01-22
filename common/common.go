package common

import (
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"log/slog"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/raynaluzier/vsphere-go-sdk/util"
)

var logLevel slog.Level

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
			//fmt.Println("URL protocol http/https not found. Continuing...")
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

func ParseUriForFilename(artifactUri string) string {
	// This can be the download or artifact URI
	fileName := path.Base(artifactUri)
	return fileName
}

func ParseFilenameForImageName(fileName string) string {
	ext := filepath.Ext(fileName)
	imageName := strings.TrimSuffix(fileName, ext)
	return imageName
}

func CheckPathType(path string) bool {
	// Checks path to see if path is Unix-based (has '/') or Windows-based (has '\')
	isWinPath := strings.Contains(path, "\\")
	return isWinPath
}

func CheckAddSlashToPath(path string) string {
	lastChar := path[len(path)-1:]
	winPath := CheckPathType(path)

	if winPath == true {
		if lastChar == "\\" {
			LogTxtHandler().Debug("Path: '" + path + "' is formatted properly")
			return path
		} else {
			// Add backslash to path
			path = path + "\\"
			return path
		}
	} else {  // Unix Path
		if lastChar == "/" {
			LogTxtHandler().Debug("Path: '" + path + "' is formatted properly")
			return path
		} else {
			// Add forwardslash to path
			path = path + "/"
			return path
		}
	}
}

func SetLoggingLevel() slog.Level {
	level := util.Logging

	switch level {
	case "INFO":
		logLevel = slog.LevelInfo
	case "WARN":
		logLevel = slog.LevelWarn
	case "ERROR":
		logLevel = slog.LevelError
	case "DEBUG":
		logLevel = slog.LevelDebug
	default:
		logLevel = slog.LevelInfo
	}
	return logLevel
}

func LogTxtHandler() *slog.Logger {
	loggingLevel := SetLoggingLevel()
	opts := &slog.HandlerOptions{
		Level: slog.Level(loggingLevel),
	}
	handler   := slog.NewTextHandler(os.Stdout, opts)
	txtLogger := slog.New(handler)
	return txtLogger
}

func LogJsonHandler() *slog.Logger {
	loggingLevel := SetLoggingLevel()
	opts := &slog.HandlerOptions{
		Level: slog.Level(loggingLevel),
	}
	handler := slog.NewJSONHandler(os.Stdout, opts)
	jsonLogger := slog.New(handler)
	return jsonLogger
}

