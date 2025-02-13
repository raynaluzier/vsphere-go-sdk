package vm

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"strconv"
	"strings"

	"github.com/raynaluzier/vsphere-go-sdk/common"
)

func SetOvfFileList(sourcePath string) ([]string, error) {
	// 'sourcePath' is full path to image file including filename
	// Verifies and sets the list of OVF files to move
	var fileName, sourcePathOnly, imageFileName string
	var fileList []string
	var fileTypes []string
	var strI string

	isWinPath := common.CheckPathType(sourcePath)
	if isWinPath == true {
		imageFileName, sourcePathOnly = common.FileNamePathFromWin(sourcePath)
	} else {
		imageFileName, sourcePathOnly = common.FileNamePathFromLnx(sourcePath)
	}

	imageName := common.ParseFilenameForImageName(imageFileName)

	items, _ := os.ReadDir(sourcePathOnly)

	fileTypes = []string{".ovf", ".mf"}
	for _, ft := range fileTypes {
		// Construct filename to check for
		fileName = imageName + ft

		// If file exists, add to copy list
		for _, item := range items {
			if item.Name() == fileName {
				fileList = append(fileList, item.Name())
			}
		}
	}

	// Check for and add associated disk files
	for i := 1; i < 15; i++ {
		strI = strconv.Itoa(i)
		fileName = imageName + "-disk" + strI + ".vmdk"

		for _, item := range items {
			if item.Name() == fileName {
				fileList = append(fileList, item.Name())
			}
		}
	}

	if len(fileList) == 0 {
		err := errors.New("No file found called: " + fileName + " in source directory: " + sourcePathOnly)
		return fileList, err 
	}
    return fileList, nil
}

func SetPathsFromDownloadUri(outputDir, downloadUri string) (string, string, string) {
	common.LogTxtHandler().Info(">>>>>>> Setting Source and Target Paths for Conversion.....")
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

	common.LogTxtHandler().Debug("Parsed Filename: " + fileName)
	common.LogTxtHandler().Debug("Parsed Image Name: " + imageName)
	common.LogTxtHandler().Debug("Parsed Image File Type: " + fileType)
	common.LogTxtHandler().Debug("Source Path for Conversion: " + sourcePath)

	if fileType == "vmtx" {
		fileNoExt := strings.TrimSuffix(fileName, "vmtx")		// win2022.vmtx returns: win2022.
		// targetPath used in rename
		targetPath = outputDir + imageName + fileNoExt + "vmx"  // returns: E:\\Lab\\win2022\\win2022.vmx
	} else {  // ova or ovf...
		targetPath = outputDir									// E:\\Lab --> ovftool will dump files to:  E:\\Lab\\win2022\\[win2022 VMX files]
	}
	common.LogTxtHandler().Debug("Target Path for Conversion: " + targetPath)

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

		common.LogTxtHandler().Debug("Parsed Filename: " + fileName)
		common.LogTxtHandler().Debug("Parsed Image Name: " + imageName)
		common.LogTxtHandler().Debug("Parsed Image File Type: " + fileType)
		common.LogTxtHandler().Debug("Image Path Being Checked For: " + checkPath)

		if strings.Contains(sourcePath, checkPath) {		              //if 'E:\\Lab\\win22\\win22.ova' contains 'win22\\win22.ova'....
			if fileType == "vmtx" {
				trimmedPath := strings.TrimSuffix(sourcePath, "vmtx")     // Ex: E:\\Lab\\win22\\win22.vmtx, returns: E:\\Lab\\win22\\win22.
				targetPath = trimmedPath + "vmx"					      // returns: E:\\Lab\\win22\\win22.vmx
				common.LogTxtHandler().Debug("Target Path: " + targetPath)
				return targetPath
			} else {  // ova or ovf
				targetPath = strings.TrimSuffix(sourcePath, checkPath)    // Ex: 'E:\\Lab\\win22\\win22.ova', returns: 'E:\\Lab\\'
				common.LogTxtHandler().Debug("Target Path: " + targetPath)
				return targetPath
			}
		} else {
			if fileType == "vmtx" {
				trimmedPath := strings.TrimSuffix(sourcePath, "vmtx")		// Ex: G:\\this\\path\\somefolder\\somefile.vmtx, returns: G:\\this\\path\\somefolder\\somefile.
				targetPath = trimmedPath + "vmx"							// returns: G:\\this\\path\\somefolder\\somefile.vmx
				common.LogTxtHandler().Debug("Target Path: " + targetPath)
				return targetPath
			} else {  // ova or ovf
				file, _ := common.GetBaseImagePathWin(sourcePath)	    // Ex: G:\\this\\path\\somefolder\\somefile.ovf, returns: somefile.ovf
				targetPath = strings.TrimSuffix(sourcePath, file)		// returns: G:\\this\\path\\somefolder\\
				common.LogTxtHandler().Debug("Target Path: " + targetPath)
				return targetPath
			}
		}
	} else { // linux path
		fileName, _ = common.FileNamePathFromLnx(sourcePath)		// Ex: /lab/rhel9/rhel9.ova, returns: rhel9.ova
		fileType = common.GetFileType(fileName)						// Ex: rhel9.ova, returns: ova
		imageName = common.ParseFilenameForImageName(fileName)		// returns: rhel9
		checkPath = imageName + "/" + fileName						// returns: rhel/rhel9.ova

		common.LogTxtHandler().Debug("Parsed Filename: " + fileName)
		common.LogTxtHandler().Debug("Parsed Image Name: " + imageName)
		common.LogTxtHandler().Debug("Parsed Image File Type: " + fileType)
		common.LogTxtHandler().Debug("Image Path Being Checked For: " + checkPath)

		if strings.Contains(sourcePath, checkPath) {		                //if '/lab/rhel9/rhel9.ova' contains 'rhel9/rhel9.ova'....
			if fileType == "vmtx" {
				trimmedPath := strings.TrimSuffix(sourcePath, "vmtx")		// Ex: /lab/rhel9/rhel9.vmtx, returns: /lab/rhel9/rhel9.
				targetPath = trimmedPath + "vmx"					        // returns: /lab/rhel9/rhel9.vmx
				common.LogTxtHandler().Debug("Target Path: " + targetPath)
				return targetPath
			} else {	// ova or ovf
				targetPath = strings.TrimSuffix(sourcePath, checkPath)	    // Ex: '/lab/rhel9/rhel9.ova', returns: '/lab/'
				return targetPath
			}
		} else {  // if some other path was used
			if fileType == "vmtx" {
				trimmedPath := strings.TrimSuffix(sourcePath, "vmtx")		// Ex: /this/path/somefolder/somefile.vmtx, returns: /this/path/somefolder/somefile.
				targetPath = trimmedPath + "vmx"							// returns: /this/path/somefolder/somefile.vmx
				common.LogTxtHandler().Debug("Target Path: " + targetPath)
				return targetPath
			} else {  // ova or ovf
				file, _ := common.GetBaseImagePathLnx(sourcePath)	    // Ex: /this/path/somefolder/somefile.ovf, returns: somefile.ovf
				targetPath = strings.TrimSuffix(sourcePath, file)		// returns: /this/path/somefolder/
				common.LogTxtHandler().Debug("Target Path: " + targetPath)
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
func SetVmPathName(path, dsName string) string {
	// 'path' is the full file path to the image file being imported to vCenter
	var vmPathName string
	isWinPath := common.CheckPathType(path)			     // E:\\labimage\\labimage.ova --> true
	ext := common.GetFileType(path)               // extension without leading '.', example 'ova'

	if isWinPath == true {
		noLetterPath := common.TrimDriveLetter(path)	// returns: labimage\\labimage.ova	
		path = common.SwapSlashes(noLetterPath)		    // returns: labimage/labimage.ova
	} else {
		path = strings.TrimPrefix(path, "/")			// trim leading slash for linux
	}

	if ext != "vmx" {
		path = strings.TrimSuffix(path, ext)  	  // returns:  labimage/labimage.
		vmPathName = "[" + dsName + "] " + path + "vmx"  // returns:  [datastore] labimage/labimage.vmx
	} else {
		vmPathName = "[" + dsName + "] " + path          // returns:  [datastore] labimage/labimage.vmx
	}
	common.LogTxtHandler().Info("vmPathName is set to: " + vmPathName)
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
		strErr := fmt.Sprintf("%v\n", err)
		common.LogTxtHandler().Error("Error: Unable to marshal json data - " + strErr)
	}
	body := bytes.NewReader(payloadBytes)

	req, err := http.NewRequest(http.MethodPost, requestPath, body)
	if err != nil {
		strErr := fmt.Sprintf("%v\n", err)
		common.LogTxtHandler().Error("Error: Error making HTTP POST request - " + strErr)
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
		strErr := fmt.Sprintf("%v\n", err)
		common.LogTxtHandler().Error("Error registering VMX with vCenter - " + strErr)
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
	var ovaMatched, ovfMatched bool

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