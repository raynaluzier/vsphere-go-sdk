package archive

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/helloyi/go-sshclient"
	"github.com/raynaluzier/vsphere-go-sdk/common"
)

/*
authInput := make(map[string]string)
authInput["method"] = "private_key"
authInput["user"] = "rt-user"
authInput["server"] = "192.168.1.xxx"
authInput["port"] = "22"
authInput["private_key_file"] = "C:\\Users\\me\\.ssh\\id_ecdsa"
*/

/*
authInput["method"] = "user_pass"
authInput["user"] = "domain.local\\someuser"
authInput["pass"] = "xyz123abc"
authInput["port"] = "22"
authInput["server"] = "192.168.1.xxx"

// Create copy list
imageName := "image1111"
imageType := "ova"
fileSuffix := ""
sourceDir := "\\\\192.168.1.xxx\\c$\\test_dir\\image1234"
targetDir := "\\\\192.168.1.xxx\\c$\\lab\\" + imageName

// Get auth client
client := common.GetAuthClient(authInput)

copyList, err := common.MakeFileCopyList(sourceDir, imageName, imageType, fileSuffix)
if err != nil {
	fmt.Println(err)
}
// Check source dir for win or lin to determine copy method

// Copy files
result := common.WinCopyFiles(sourceDir, targetDir, copyList, client)
fmt.Println(result)
*/

func MakeFileCopyList(sourceDir, imageName, imageType, fileSuffix string) ([]string, error) {
	// sourceDir ex: c:\\lab or /lab
	// imageType ex: ova, ovf, vmtx
	var fileName string
	var copyList []string
	var fileTypes []string
	var strI string
	imageType = strings.ToLower(imageType)

	// Make sure directories have ending slash, based on platform type
	sourceDir = common.CheckAddSlashToPath(sourceDir)

	items, _ := os.ReadDir(sourceDir)

	if imageType == "ova" {
		// Construct filename to check for
		if fileSuffix != "" {
			fileName = imageName + "-" + fileSuffix + ".ova"
		} else {
			fileName = imageName + ".ova"
		}

		// If file exists, add to copy list
		for _, item := range items {
			if item.Name() == fileName {
				copyList = append(copyList, item.Name())
			}
		}

		if len(copyList) == 0 {
			err := errors.New("No file found called: " + fileName + " in source directory: " + sourceDir)
			return copyList, err 
		}

	} else if imageType == "ovf" {
		fileTypes = []string{".ovf", ".mf"}
		for _, ft := range fileTypes {
			// Construct filename to check for
			if fileSuffix != "" {
				fileName = imageName + "-" + fileSuffix + ft
			} else {
				fileName = imageName + ft
			}

			// If file exists, add to copy list
			for _, item := range items {
				if item.Name() == fileName {
					copyList = append(copyList, item.Name())
				}
			}
		}

		// Check for and add associated disk files
		for i := 1; i < 15; i++ {
			strI = strconv.Itoa(i)
			if fileSuffix != "" {
				fileName = imageName + "-" + fileSuffix + "-disk-" + strI + ".vmdk"
			} else {
				fileName = imageName + "-disk-" + strI + ".vmdk"
			}

			for _, item := range items {
				if item.Name() == fileName {
					copyList = append(copyList, item.Name())
				}
			}
		}

		if len(copyList) == 0 {
			err := errors.New("No file found called: " + fileName + " in source directory: " + sourceDir)
			return copyList, err 
		}

	} else if imageType == "vmtx" {
		fileTypes = []string{".nvram", ".vmsd", ".vmtx", ".vmxf"}
		for _, ft := range fileTypes {
			// Construct file names from static file types
			if fileSuffix != "" {
				fileName = imageName + "-" + fileSuffix + ft
			} else {
				fileName = imageName + ft
			}

			// If file exists, add to copy list
			for _, item := range items {
				if item.Name() == fileName {
					copyList = append(copyList, item.Name())
				}
			}
		}

		// Construct and search for non-numbered virtual disk file -------->
		if fileSuffix != "" {
			fileName = imageName + "-" + fileSuffix + ".vmdk"
		} else {
			fileName = imageName + ".vmdk"
		}

		for _, item := range items {
			if item.Name() == fileName {
				copyList = append(copyList, item.Name())
			}
		}

		// Construct and search for numbered virtual disk files ------------>
		for i := 1; i < 15; i++ {
			strI = strconv.Itoa(i)
			if fileSuffix != "" {
				fileName = imageName + "-" + fileSuffix + "_" + strI + ".vmdk"
			} else {
				fileName = imageName + "_" + strI + ".vmdk"
			}

			for _, item := range items {
				if item.Name() == fileName {
					copyList = append(copyList, item.Name())
				}
			}
		}

		// Construct and search for -ctk disk files ------------------------>
		// Non-numbered -ctk files
		if fileSuffix != "" {
			fileName = imageName + "-" + fileSuffix + "-ctk.vmdk"
		} else {
			fileName = imageName + "-ctk.vmdk"
		}

		for _, item := range items {
			if item.Name() == fileName {
				copyList = append(copyList, item.Name())
			}
		}

		// Numbered -ctk files
		for i := 1; i < 15; i++ {
			strI = strconv.Itoa(i)
			if fileSuffix != "" {
				fileName = imageName + "-" + fileSuffix + "_" + strI + "-ctk.vmdk" 
			} else {
				fileName = imageName + "_" + strI + "-ctk.vmdk"
			}

			for _, item := range items {
				if item.Name() == fileName {
					copyList = append(copyList, item.Name())
				}
			}
		}

		// Construct and search for -flat disk files ------------------------>
		// Non-numbered virtual disk file
		if fileSuffix != "" {
			fileName = imageName + "-" + fileSuffix + "-flat.vmdk"
		} else {
			fileName = imageName + "-flat.vmdk"
		}

		for _, item := range items {
			if item.Name() == fileName {
				copyList = append(copyList, item.Name())
			}
		}

		// Numbered virtual disk files
		// Numbered -ctk files
		for i := 1; i < 15; i++ {
			strI = strconv.Itoa(i)
			if fileSuffix != "" {
				fileName = imageName + "-" + fileSuffix + "_" + strI + "-flat.vmdk" 
			} else {
				fileName = imageName + "_" + strI + "-flat.vmdk"
			}

			for _, item := range items {
				if item.Name() == fileName {
					copyList = append(copyList, item.Name())
				}
			}
		}

		// vmware.log file
		fileName = "vmware.log"
		for _, item := range items {
			if item.Name() == fileName {
				copyList = append(copyList, item.Name())
			}
		}

		if len(copyList) == 0 {
			err := errors.New("No matching files found in source directory: " + sourceDir)
			common.LogTxtHandler().Error("Check that the provided source director and/or image name exists.")
			return copyList, err 
		}
	}
	return copyList, nil
}

func RunScriptTestSSH(client *sshclient.Client) (string) {
	script := `
		cmd
		/c
		ipconfig /all
	`
	out, err := client.Script(script).Output()
	if err != nil {
		fmt.Println(err)
	}
	defer client.Close()
	fmt.Println(string(out))

	return "Complete"
}

func WinCopyFiles(sourceDir, targetDir string, copyList []string, client *sshclient.Client) (string) {
	// Check path type, then add end slash if needed
	sourceDir = common.CheckAddSlashToPath(sourceDir)
	targetDir = common.CheckAddSlashToPath(targetDir)

	// Check if target directory exists, if not, create...
	_, err := os.Stat(targetDir)
	if os.IsNotExist(err) {
		err = os.Mkdir(targetDir, 0755)
		if err != nil {
			strErr := fmt.Sprintf("%v\n", err)
			common.LogTxtHandler().Error("Error creating directory: " + targetDir + " - " + strErr)
		} else {
			common.LogTxtHandler().Info("Successfully created directory: " + targetDir)
		}
	}
	
	for _, file := range copyList {
		scriptB := fmt.Sprintf(`robocopy %s %s %s`, sourceDir, targetDir, file) 

		out, err := client.Script(scriptB).Output()
		if err != nil {
			strErr := fmt.Sprintf("%v\n", err)
			common.LogTxtHandler().Error("Error copying file: " + file + " - " + strErr)
			return "Copy Process Failed"
		} else {
			common.LogTxtHandler().Info(string(out))
		}
	}

	defer client.Close()
	return "End of Copy Process"
}
