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

func SetPathsFromDownloadUri(outputDir, downloadUri string) (string, string, string) {
	// Takes output directory and download URI, parses the image name from the download URI and determines 
	// the image file type, source path, and target path that will be used with the image conversion process, if needed

	// 'downloadUri' is the Artifactory path where the image files were downloaded from; this is used to determine the fileName, imageName, and source file type
	// without having the user provide it.
	// 'outputDir' is the location where the image files were originally downloaded to; the download process will put them in their own image named-based folder,
	//				so this is the path we'll use to generate the source directory of the image files for conversion
	// 'targetPath' is where the converted VMX image files will output TO; typically we'd want to make this the same place where the files were downloaded
	//              as this is the directory where vCenter will import the machine from and where those template files will exist going forward

	//     NOTE:  The conversion tool automatically creates a directory based on the image name that it places the files into, so we don't need to include
	//            an image name-based folder in our target path.
	//            
	// 'fileType' is pulled from file name (ext without ".")
	var sourcePath, targetPath string
	outputDir  = common.CheckAddSlashToPath(outputDir)
	fileName  := common.ParseUriForFilename(downloadUri)
	imageName := common.ParseFilenameForImageName(fileName)     // includes ending slash
	fileType  := common.GetFileType(fileName)				 
	sourcePath = outputDir + imageName + fileName				// E:\\Lab\\win2022\\win2022.ova

	if fileType == "vmtx" {
		fileNoExt := strings.TrimSuffix(fileName, "vmtx")		// win2022.vmtx returns: win2022.
		// targetPath used in rename
		targetPath = outputDir + imageName + fileNoExt + "vmx"  // returns: E:\\Lab\\win2022\\win2022.vmx
	} else {  // ova or ovf...
		targetPath = outputDir									// E:\\Lab --> ovftool will dump files to:  E:\\Lab\\win2022\\[win2022 VMX files]
	}
	return fileType, sourcePath, targetPath
}

func SetPathNoDownload(sourcePath string) string {
	// Target path handling for instances where we are just converting and importing, but not downloading first
	// As converted images should reside with source image files, target path is formed from source path
	var targetPath, fileName, fileType, imageName, checkPath string
	isWinPath := common.CheckPathType(sourcePath)

	if isWinPath == true {
		fileName, _ = common.FileNamePathFromWin(sourcePath)		// Ex: E:\\Lab\\win22\\win22.ova, returns: win22.ova
		fileType  = common.GetFileType(fileName)					// Ex: win22.ova, returns: ova
		imageName = common.ParseFilenameForImageName(fileName)		// returns: win22
		checkPath = imageName + "\\" + fileName						// returns: win22\\win22.ova

		if strings.Contains(sourcePath, checkPath) {		              //if 'E:\\Lab\\win22\\win22.ova' contains 'win22\\win22.ova'....
			if fileType == "vmtx" {
				trimmedPath := strings.TrimSuffix(sourcePath, "vmtx")     // Ex: E:\\Lab\\win22\\win22.vmtx, returns: E:\\Lab\\win22\\win22.
				targetPath = trimmedPath + "vmx"					      // returns: E:\\Lab\\win22\\win22.vmx
				return targetPath
			} else {  // ova or ovf
				targetPath = strings.TrimSuffix(sourcePath, checkPath)    // Ex: 'E:\\Lab\\win22\\win22.ova', returns: 'E:\\Lab\\'
				return targetPath
			}
		} else {
			if fileType == "vmtx" {
				trimmedPath := strings.TrimSuffix(sourcePath, "vmtx")		// Ex: G:\\this\\path\\somefolder\\somefile.vmtx, returns: G:\\this\\path\\somefolder\\somefile.
				targetPath = trimmedPath + "vmx"							// returns: G:\\this\\path\\somefolder\\somefile.vmx
				return targetPath
			} else {  // ova or ovf
				file, parentDir := common.GetBaseImagePathWin(sourcePath)	// Ex: G:\\this\\path\\somefolder\\somefile.ovf, returns: somefile.ovf, somefolder
				altPath := parentDir + "\\" + file							// returns: somefolder\\somefile.ovf
				targetPath = strings.TrimSuffix(sourcePath, altPath)		// returns: G:\\this\\path\\
				return targetPath
			}
		}
	} else { // linux path
		fileName, _ = common.FileNamePathFromLnx(sourcePath)		// Ex: /lab/rhel9/rhel9.ova, returns: rhel9.ova
		fileType = common.GetFileType(fileName)						// Ex: rhel9.ova, returns: ova
		imageName = common.ParseFilenameForImageName(fileName)		// returns: rhel9
		checkPath = imageName + "/" + fileName						// returns: rhel/rhel9.ova

		if strings.Contains(sourcePath, checkPath) {		                //if '/lab/rhel9/rhel9.ova' contains 'rhel9/rhel9.ova'....
			if fileType == "vmtx" {
				trimmedPath := strings.TrimSuffix(sourcePath, "vmtx")		// Ex: /lab/rhel9/rhel9.vmtx, returns: /lab/rhel9/rhel9.
				targetPath = trimmedPath + "vmx"					        // returns: /lab/rhel9/rhel9.vmx
				return targetPath
			} else {	// ova or ovf
				targetPath = strings.TrimSuffix(sourcePath, checkPath)	    // Ex: '/lab/rhel9/rhel9.ova', returns: '/lab/'
				return targetPath
			}
		} else {  // if some other path was used
			if fileType == "vmtx" {
				trimmedPath := strings.TrimSuffix(sourcePath, "vmtx")		// Ex: /this/path/somefolder/somefile.vmtx, returns: /this/path/somefolder/somefile.
				targetPath = trimmedPath + "vmx"							// returns: /this/path/somefolder/somefile.vmx
				return targetPath
			} else {  // ova or ovf
				file, parentDir := common.GetBaseImagePathLnx(sourcePath)	// Ex: /this/path/somefolder/somefile.ovf, returns: somefile.ovf, somefolder
				altPath := parentDir + "/" + file							// returns: somefolder/somefile.ovf
				targetPath = strings.TrimSuffix(sourcePath, altPath)		// returns: /this/path/
				return targetPath
			}
		}
	}
}

func ConvertImageByType(fileType, sourcePath, targetPath string) string {
	var result string

	switch fileType {
	case "ova":
		common.LogTxtHandler().Info("File type found: " + fileType + "; converting to vmx...")
		result = ConvertOvfaToVmx(sourcePath, targetPath)
	case "ovf":
		common.LogTxtHandler().Info("File type found: " + fileType + "; converting to vmx...")
		result = ConvertOvfaToVmx(sourcePath, targetPath)
	case "vmtx":
		common.LogTxtHandler().Info("File type found: " + fileType + "; converting to vmx...")
		result = common.RenameFile(sourcePath, targetPath)
	case "vmx":
		common.LogTxtHandler().Info("File is already in format: vmx")
		result = "Success"
	default:
		common.LogTxtHandler().Error("Found unsupported file type: " + fileType)
		common.LogTxtHandler().Error("Supported file types are: ova, ovf, vmtx, and vmx")
		result = "Failed"
	}
	return result
}

// we need to account for no download option...
func SetVmPathName(sourcePath, dsName string) string {
	// 'sourcePath' comes from result of SetPathsFromDownloadUri; we're using this to form the vmPathName, which should match the
	// resulting target path during the image conversion process --> Example: E:\\labimage\\labimage.ova
	isWinPath := common.CheckPathType(sourcePath)			// E:\\labimage\\labimage.ova --> true

	if isWinPath == true {
		noLetterPath := common.TrimDriveLetter(sourcePath)	// returns: labimage\\labimage.ova	
		sourcePath = common.SwapSlashes(noLetterPath)		// returns: labimage/labimage.ova
	}

	ext := common.GetFileType(sourcePath)             // extension without leading '.', example 'ova'
	path := strings.TrimSuffix(sourcePath, ext)  	  // returns:  labimage/labimage.
	vmPathName := "[" + dsName + "] " + path + "vmx"  // returns:  [datastore] labimage/labimage.vmx
	return vmPathName
}

func RegisterVm(token, vcServer, dcName, vmPathName, imageName, folderId, resPoolId string) string {
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
	// 'vmPathName' needs to match the resulting target path of the image file after the image conversion
	data := Payload{ 
		DatastorePath: vmPathName,
		Name: imageName,
		Placement: Placement{
			Folder: folderId,
			ResourcePool: resPoolId,
		},
	}

	payloadBytes, err := json.Marshal(data)

	if err != nil {
		common.LogTxtHandler().Error("Error: Unable to marshal json data - ", err)
	}
	body := bytes.NewReader(payloadBytes)

	req, err := http.NewRequest(http.MethodPost, requestPath, body)
	if err != nil {
		common.LogTxtHandler().Error("Error: Error making HTTP POST request - ", err)
	}
	req.Header.Set("Content-Type", "application/json")
	newToken := common.TrimQuotes(token) // Required or auth will fail
	req.Header.Set("vmware-api-session-id", newToken)
	
	//v1 := req.Header.Get("vmware-api-session-id")
	//fmt.Println(v1)

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
	strResp := fmt.Sprintf("%v\n", resp)
	common.LogTxtHandler().Debug(strResp)

	if err != nil {
		common.LogTxtHandler().Error("Error registering VMX with vCenter - ", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		statusCode = "200"
	} else {
		statusCode = (fmt.Sprintf("%v", resp.StatusCode))
		common.LogTxtHandler().Error("Error registering VMX with vCenter. Validate inputs and ensure image is not already in the target inventory.")
	}
	common.LogTxtHandler().Debug("Status Code: " + statusCode)
	return statusCode
}

// These OVF/OVA conversion functions require the OVFTool be installed
	// Converts OVF/OVA to VMX
	// For Windows: 'cmd' and '/c' are not included in the commands, this will fail
	// Input can be local or URL
	// Ensure local paths are escaped properly (ex: C:\\lab\\file.vmx)
func ConvertOvfaToVmx(inputPath, outputPath string) string {
	var cmd *exec.Cmd
	var ovaMatched, ovfMatched, vmxMatched bool

	inputPath = strings.ToLower(inputPath)                   // Ex: e:\\lab-servs\\image1234.ova
	outputPath = strings.ToLower(outputPath)				 // Ex: e:\\lab-servs\\image1234.vmx
	
	// Ensure input is either an OVA or OVF file
	ovfMatched, err := regexp.MatchString("ovf", inputPath)
	if ovfMatched == false {
		ovaMatched, err = regexp.MatchString("ova", inputPath)
		if ovaMatched == false {
			common.LogTxtHandler().Error("Error: Input is neither an OVA or OVF file.")
			common.LogTxtHandler().Error("Please provide the full file path to the OVA or OVF that's to be converted.") 
		}
	}
	if err != nil {
		common.LogTxtHandler().Error("Unable to search for OVA/OVF string.")
	}
	
	// Ensure output is a VMX file
	vmxMatched, err = regexp.MatchString("vmx", outputPath)
	if vmxMatched == false {
		common.LogTxtHandler().Error("Error: Output does not include a VMX file path.")
		common.LogTxtHandler().Error("Please provide the full destination file path to the resulting VMX file.") 	
	}
	if err != nil {
		common.LogTxtHandler().Error("Unable to search for VMX string.")
	}
	
	ovfCmd := "ovftool " + inputPath + " " + outputPath								// Ex:  ovftool e:\\lab-servs\\image1234.ova e:\\lab-servs\\image1234.vmx
	fmt.Println("OVF CMD: " + ovfCmd)

	common.LogTxtHandler().Info("Beginning conversion process... This could take a while.")
	switch runtime.GOOS{
	case "windows":
		fmt.Println("Running Windows shell...")
		cmd = exec.Command("cmd", "/c", ovfCmd)
	default: // mac & linux
		fmt.Println("Running Linux shell...")
		cmd = exec.Command(ovfCmd)   // mac "bash"
	}

	cmd.Stdout = os.Stdout

	if err := cmd.Run(); err != nil {
		common.LogTxtHandler().Error("Could not run the exec shell command.")
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