package common

import (
	"crypto/tls"
	"fmt"
	"io"
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
	var token string
	server = AddHttpsProtocol(server)  // checks for http and adds https:// if missing
	requestPath := server + "/api/session"
	request, err := http.NewRequest("POST", requestPath, nil)
	request.SetBasicAuth(user, pass)
	
	if err != nil {
		strErr := fmt.Sprintf("%v\n", err)
		LogTxtHandler().Error("Error with POST request - " + strErr)
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
		LogTxtHandler().Error("Response error: " + strErr)
		LogTxtHandler().Error("Check the server name/IP provided.")
	}
	
	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		strErr := fmt.Sprintf("%v\n", err)
		LogTxtHandler().Error("Error reading response body - " + strErr)
	}
	
	// if server name contains https, try http
	respBody := string(body)
	if strings.Contains(respBody, "Authentication required.") && strings.Contains(server, "https://") {
		server = AddHttpProtocol(server)  // Replaces https with http first, retries auth
		token = VcenterAuth(user, pass, server)
		return token
	} else {
		token = string(body)
		return token
	}
}

func AddHttpsProtocol(server string) string {
	var serverUrl string
	server = strings.ToLower(server)

	matched, err := regexp.MatchString("http", server)
	if err != nil {
		strErr := fmt.Sprintf("%v\n", err)
		LogTxtHandler().Error("Error searching for string - " + strErr)
	}

	if matched == false {
		serverUrl = "https://" + server
		LogTxtHandler().Info("Web protocol missing. Adding 'https://' to server...")
		return serverUrl
	} else {
		return server
	}
}

func AddHttpProtocol(server string) string {
	var serverUrl string
	server = strings.ToLower(server)

	matched, err := regexp.MatchString("https", server)
	if err != nil {
		strErr := fmt.Sprintf("%v\n", err)
		LogTxtHandler().Error("Error searching for string - " + strErr)
	}

	if matched == true {
		server = strings.TrimPrefix(server, "https://")
		serverUrl = "http://" + server
		LogTxtHandler().Info("Adding 'http://' to server...")
		return serverUrl
	} else {
		if strings.Contains(server, "http://") {
			return server
		} else { // server doesn't contain either https or http
			serverUrl = "http://" + server
			return serverUrl
		}
	}
}

func TrimUrlProtocol(serverUrl string) string {
	// Sets server URL to lowercase first then searches
	var server string
	serverUrl = strings.ToLower(serverUrl)

	sMatched, err := regexp.MatchString("https", serverUrl)
	if err != nil {
		strErr := fmt.Sprintf("%v\n", err)
		LogTxtHandler().Error("Error searching for string - " + strErr)
	}

	if sMatched == true {
		server = strings.TrimPrefix(serverUrl, "https://")
		return server
	} else {
		matched, err := regexp.MatchString("http", serverUrl)
		if err != nil {
			strErr := fmt.Sprintf("%v\n", err)
			LogTxtHandler().Error("Error searching for string - " + strErr)
			return ""
		}

		if matched == true {
			server = strings.TrimPrefix(serverUrl, "http://")
			return server
		} else {
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
		strErr := fmt.Sprintf("%v\n", err)
        LogTxtHandler().Error("Error renaming file - " + strErr)
		return "Failed"
    } else {
		return "Success"
	}
}

func GetFileType(filePath string) string {
	filePath = strings.ToLower(filePath)

	ext := filepath.Ext(filePath)
	if ext == "" {
		LogTxtHandler().Error("No file was included in file path; file type is empty string.")
		return ext
	} else {
		if strings.Contains(ext, ".") {
			ext = strings.TrimPrefix(ext, ".")
			return ext
		}
		return ext
	}
}

func ParseUriForFilename(artifactUri string) string {
	// This can be the download or artifact URI
	fileName := path.Base(artifactUri)
	if strings.Contains(fileName, ".") {
		return fileName
	} else {
		LogTxtHandler().Error("URI doesn't contain a complete filename.")
		LogTxtHandler().Error("Artifact URI provided: " + artifactUri)
		LogTxtHandler().Error("Filename that was parsed from URI: " + fileName)
		return fileName
	}
}

func ParseFilenameForImageName(fileName string) string {
	//Returns image name even if no ext found
	ext := filepath.Ext(fileName)
	imageName := strings.TrimSuffix(fileName, ext)
	return imageName
}

func CheckPathType(path string) bool {
	// Checks path to see if path is Unix-based (has '/') or Windows-based (has '\')
	isWinPath := strings.Contains(path, "\\")
	return isWinPath
}

func FileNamePathFromWin(path string) (string, string) {
	segments := strings.Split(path, "\\")	    // Split file path into segments
	fileName := segments[len(segments)-1]	     // Determine filename from path
	filePath := path[:len(path)-len(fileName)]   // Determine just path without filename
	return fileName, filePath
}

func FileNamePathFromLnx(path string) (string, string) {
	segments := strings.Split(path, "/")	    // Split file path into segments
	fileName := segments[len(segments)-1]	     // Determine filename from path
	filePath := path[:len(path)-len(fileName)]   // Determine just path without filename
	return fileName, filePath
}

func GetBaseImagePathWin(sourcePath string) (string, string) {
	segments := strings.Split(sourcePath, "\\")	    // Split file path into segments
	fileName := segments[len(segments)-1]	    // Determine filename from path (no begin or end slash)
	parentDir := segments[len(segments)-2]		// Determine parent directory of file (no begin or end slash)
	return fileName, parentDir
}

func GetBaseImagePathLnx(sourcePath string) (string, string) {
	segments := strings.Split(sourcePath, "/")	    // Split file path into segments
	fileName := segments[len(segments)-1]	    // Determine filename from path (no begin or end slash)
	parentDir := segments[len(segments)-2]		// Determine parent directory of file (no begin or end slash)
	return fileName, parentDir
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

func TrimDriveLetter(path string) string {
	// For path 'c:\lab\file.txt', returns 'lab\file.txt'
	i := strings.Index(path, ":")
	if i > -1 {
		remainingPath := path[i+2:]
		return remainingPath
	}
	return path
}

func SwapSlashes(path string) string {
	// Changes win path to unix path
	if strings.Contains(path, "\\") {
		newPath := strings.ReplaceAll(path, "\\", "/")
		return newPath
	}
	return path
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
