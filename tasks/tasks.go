package tasks

import (
	"fmt"
	"strings"

	"github.com/raynaluzier/vsphere-go-sdk/common"
	"github.com/raynaluzier/vsphere-go-sdk/govmomi"
	"github.com/raynaluzier/vsphere-go-sdk/util"
	"github.com/raynaluzier/vsphere-go-sdk/vm"
)

func GetResourceIds(vcUser, vcPass, vcServer, dcName, folderName, resPoolName, clusterName string) (string, string, error) {
	util.VcUser     = vcUser
	util.VcPassword = vcPass
	util.VcServer   = vcServer

	common.LogTxtHandler().Info("Getting resource IDs for target vSphere folder and resource pool...")

	common.LogTxtHandler().Debug("Getting folder ID...")
	folderId, err := govmomi.GetFolderId(vcUser, vcPass, vcServer, folderName, dcName)
	common.LogTxtHandler().Info("Folder ID: " + folderId)
	if err != nil {
		strErr := fmt.Sprintf("%v", err)
		common.LogTxtHandler().Error("Error getting folder ID: %s - " + strErr)
		return "", "", err
	}

	common.LogTxtHandler().Debug("Getting resource pool ID...")
	resPoolId, err := govmomi.GetResPoolId(vcUser, vcPass, vcServer, resPoolName, dcName, clusterName)
	common.LogTxtHandler().Info("Resource Pool ID: " + resPoolId)
	if err != nil {
		strErr := fmt.Sprintf("%v", err)
		common.LogTxtHandler().Error("Error getting resource pool ID: %s - " + strErr)
		return folderId, "", err
	}
	return folderId, resPoolId, nil
}

func ConvertImportFromDownload(vcUser, vcPass, vcServer, outputDir, downloadUri, dcName, dsName, dsImagePath, imageName, folderId, resPoolId string) string {
	// download process automatically puts image into own image-based dir; so don't need to include in the output dir path...
	util.VcUser     = vcUser
	util.VcPassword = vcPass
	util.VcServer   = vcServer
	util.OutputDir  = outputDir
	var vmPathName, fileType, sourcePath, targetPath string
	var imageFileName, sourceFolderPath, postConvTargetPath, convertResult, postConvTargetFilePath string

	// outputDir:   /mnt/servers										E:\\Lab\\win22
	// downloadUri: https://art.server.com/repo/folder/ub20/ub20.ovf	https://art.server.com/repo/folder/win22/win22.ova

	if outputDir != "" && downloadUri != "" {
		common.LogTxtHandler().Info("--- Setting source and target paths...")
		fileType, sourcePath, targetPath = vm.SetPathsFromDownloadUri(outputDir, downloadUri)
		// Example return:
			// fileType:    ovf							ova								vmtx
			// sourcePath:	/mnt/servers/ub20/ub20.ovf  E:\\Lab\\win22\\win22.ova		/mnt/server/ub24/ub24.vmtx
			// targetPath:	/mnt/servers/				/mnt/servers/

		common.LogTxtHandler().Info("File Type: " + fileType)			
		common.LogTxtHandler().Info("Source Path: " + sourcePath)		
		common.LogTxtHandler().Info("Target Path: " + targetPath)		
		
		common.LogTxtHandler().Info("Converting image by image type... Please wait...")

		// Before we convert, if ovf, we need to create an ovf_files folder and move the files into it to avoid file conflicts - and this will be new source
		isWinPath  := common.CheckPathType(sourcePath)
		if isWinPath == true {
			common.LogTxtHandler().Debug("---> Path type is WINDOWS-based.")
			imageFileName, sourceFolderPath = common.FileNamePathFromWin(sourcePath)
			imageName = common.ParseFilenameForImageName(imageFileName)		            
			// imageFileName:     win22.ova
			// sourceFolderPath:  E:\\Lab\\win22
			// imageName:         win22
				
			if fileType == "ova" || fileType == "ovf" {		// vmtx files have a target path that includes full path to VMX file, the other types just have a folder target
				common.LogTxtHandler().Debug("File type is either OVA or OVF...")
				common.LogTxtHandler().Debug("Converted files will be placed in the target path within an image name-based sub folder.")
				common.LogTxtHandler().Debug("This is the post-conversion target path (used for vmPathName to vCenter).")

				postConvTargetPath = targetPath + imageName
				postConvTargetPath = common.CheckAddSlashToPath(postConvTargetPath)
				// postConvTargetPath = E:\\Lab\\win22\\
			} else {
				common.LogTxtHandler().Debug("File type is either VMTX or VMX...")
				common.LogTxtHandler().Debug("Target path is the same as the post-conversion target path.")

				postConvTargetPath = targetPath
				// postConvTargetPath = E:\\Lab\\win19\\win19.vmx
			}
		} else {  // Linux path
			common.LogTxtHandler().Debug("---> Path type is LINUX-based.")
			imageFileName, sourceFolderPath = common.FileNamePathFromLnx(sourcePath)
			imageName = common.ParseFilenameForImageName(imageFileName)
			// imageFileName:     ub20.ovf
			// sourceFolderPath:  /mnt/servers/ub20
			// imageName:         ub20

			if fileType == "ova" || fileType == "ovf" {
				common.LogTxtHandler().Debug("File type is either OVA or OVF...")
				common.LogTxtHandler().Debug("Converted files will be placed in the target path within an image name-based sub folder.")
				common.LogTxtHandler().Debug("This is the post-conversion target path (used for vmPathName to vCenter).")

				postConvTargetPath = targetPath + imageName
				postConvTargetPath = common.CheckAddSlashToPath(postConvTargetPath)
				// postConvTargetPath = /mnt/servers/ub20/
			} else {
				common.LogTxtHandler().Debug("File type is either VMTX or VMX...")
				common.LogTxtHandler().Debug("Target path is the same as the post-conversion target path.")

				postConvTargetPath = targetPath
				// postConvTargetPath = /mnt/server/ub24/ub24.vmx
			}
		}

		common.LogTxtHandler().Info("Image Filename: " + imageFileName)
		common.LogTxtHandler().Info("Image Name: " + imageName)
		common.LogTxtHandler().Info("File Type: " + fileType)
		common.LogTxtHandler().Info("Source Path: " + sourcePath)
		common.LogTxtHandler().Info("Target Path (before conversion): " + targetPath)
		common.LogTxtHandler().Info("Source Folder Path: " + sourceFolderPath)
		common.LogTxtHandler().Info("Post Conversion Target Path: " + postConvTargetPath)

		// If this is an OVF image, we need to first move the image files into a sub dir called "ovf_files" and update the conversion source path to here
		// If not, we'll get a file conflict with the disk file(s)
		if fileType == "ovf" {
			common.LogTxtHandler().Info("OVF file detected...")
			common.LogTxtHandler().Info("Moving OVF files into subdirectory of source path called 'ovf_files'...")
			common.LogTxtHandler().Debug("OVF files contain disk files named the same as the resulting disk files. If we don't move them, there will be a file conflict during conversion.")
			
			destDir := sourceFolderPath + "ovf_files"					  // /mnt/servers/ub20/ovf_files         E:\\Lab\\win22\\ovf_files
			destDir = common.CheckAddSlashToPath(destDir)       		  // add ending slash by os type - /mnt/servers/ub20/ovf_files/     E:\\Lab\\win22\\ovf_files\\
			moveList, err := vm.SetOvfFileList(sourcePath)                // Get list of OVF files to move
			err = common.MoveFiles(moveList, sourceFolderPath, destDir)   // [file list], 'E:\\path\\to\\win2022\\', 'E:\\path\\to\\win2022\\ovf_files\\'
			if err != nil {
				strErr := fmt.Sprintf("%v", err)
				common.LogTxtHandler().Error("Error moving files: " + strErr)
			} else {
				common.LogTxtHandler().Info("Files moved successfully!")
			}

			// Setting the new conversion source path to the ovf_file dir for the conversion process only
			newSourcePath := destDir + imageFileName				        // ex: E:\\path\\to\\win2022\\ovf_files\\win2022.ovf
			common.LogTxtHandler().Debug("New Source Path for Image Files (after OVF files moved): " + newSourcePath)
			common.LogTxtHandler().Info("Checking image type and converting if necessary. This may time some time...")
			convertResult = vm.ConvertImageByType(fileType, newSourcePath, targetPath)

		} else {  // ova and vmtx don't need to be moved first
			common.LogTxtHandler().Info("OVA, VMTX, or VMX files detected...")
			common.LogTxtHandler().Info("Checking image type and converting if necessary. This may time some time...")
			convertResult = vm.ConvertImageByType(fileType, sourcePath, targetPath)
		}

		common.LogTxtHandler().Info("---> Conversion Result: " + convertResult)
		// Set the vmPathName value in prep for vCenter import
		if convertResult != "Failed" {
			// dsImagePath - represents datastore path to image without any mount point info from Packer server
			common.LogTxtHandler().Info("Setting vmPathName....")
			if dsImagePath == "" {
				if fileType == "ova" || fileType == "ovf" { // don't do for vmx or vmtx files as full path to file is already being passed in those cases
					postConvTargetFilePath = postConvTargetPath + imageName + ".vmx"   
					vmPathName = vm.SetVmPathName(postConvTargetFilePath, dsName)
					// postConvTargetFilePath = /mnt/servers/ub20/ub20.vmx (this would break...)          
					//						  = E:\\Lab\\win22\\win22.vmx 
					// vmPathName = [dsName] mnt/server/ub24/ub24.vmx (this would break...)
					//				[dsName] Lab/win22/win22.vmx					

				} else {
					vmPathName = vm.SetVmPathName(postConvTargetPath, dsName)
					// postConvTargetFilePath = /mnt/server/ub24/ub24.vmx (this would break...)
					//  						E:\\Lab\\win22\\win22.vmx
					// vmPathName = [dsName] mnt/server/ub24/ub24.vmx (this would break...)
					//				[dsName] Lab/win22/win22.vmx
				}
			} else {
				common.LogTxtHandler().Info("Datastore image pathing set; using this to set vmPathName...")
				dsImagePath = common.CheckAddSlashToPath(dsImagePath)   // Ex: /dev-servers/
				dsImagePath = dsImagePath + imageName					// Ex: /dev-servers/ub24
				dsImagePath = common.CheckAddSlashToPath(dsImagePath)   // Ex: /dev-servers/ub24/
				dsImagePath = dsImagePath + imageFileName				// Ex: /dev-servers/ub24/ub24.ovf
				vmPathName = vm.SetVmPathName(dsImagePath, dsName)		// Will rename file in path to be VMX
				// dsImagePath = /dev-servers/ub24/ub24.vmx
				// vmPathName = [dsName] dev-servers/ub24/ub24.vmx
			}

			common.LogTxtHandler().Info("Beginning import into vCenter....")
			vcToken := common.VcenterAuth(vcUser, vcPass, vcServer)
			statusCode := vm.RegisterVm(vcToken, vcServer, dcName, vmPathName, imageName, folderId, resPoolId)
			common.LogTxtHandler().Info("Status Code of Register VM task: " + statusCode)
		
			if statusCode == "200" {
				common.LogTxtHandler().Info("Import successful. Marking image as a VM Template...")
				tempResult := govmomi.MarkAsTemplate(vcUser, vcPass, vcServer, imageName, dcName)
				common.LogTxtHandler().Info(tempResult)
		
					if strings.Contains(tempResult, "Success") {
						common.LogTxtHandler().Info("The image import and template conversion completed successfully.")
						return "Success"
					} else {
						common.LogTxtHandler().Error("Error: Unable to import and/or convert the image into a VM Template.")
						return "Failed"
					}
			} else {
				common.LogTxtHandler().Error("Error registering VMX file with vCenter.")
				return "Failed"
			}
		} else {
			common.LogTxtHandler().Error("Error during image type check and file conversion process.")
			return "Failed"
		}
	} else {
		common.LogTxtHandler().Error("Missing output directory and/or download URI.")
		return "Failed"
	}
}

func ConvertImportNoDownload(vcUser, vcPass, vcServer, dcName, dsName, sourcePath, dsImagePath, folderId, resPoolId string) string {
	util.VcUser     = vcUser
	util.VcPassword = vcPass
	util.VcServer   = vcServer

	var imageFileName, sourceFolderPath, vmPathName, fileType, convertResult string
	var postConvTargetPath, postConvTargetFilePath, imageName string
	common.LogTxtHandler().Info("---> 'import_no_download' flag is set to TRUE. Skipping artifact download....")

	if sourcePath != "" {
		common.LogTxtHandler().Info("--- Setting target path and determining path type...")
		targetPath := vm.SetPathNoDownload(sourcePath)					                // Ex: ova/ovf = E:\Lab, vmtx = E:\Lab\win22\win22.vmx
		isWinPath  := common.CheckPathType(sourcePath)
		if isWinPath == true {
			common.LogTxtHandler().Debug("---> Path type is WINDOWS-based.")
			imageFileName, sourceFolderPath = common.FileNamePathFromWin(sourcePath)	// Ex: E:\Lab\win22\win22.ova, returns: win22.ova, E:\Lab\win22\
			imageName = common.ParseFilenameForImageName(imageFileName)		            // Ex: win22.ova, returns win22
			fileType  = common.GetFileType(imageFileName)
				
			if fileType == "ova" || fileType == "ovf" {		// vmtx files have a target path that includes full path to VMX file, the other types just have a folder target
				common.LogTxtHandler().Debug("File type is either OVA or OVF...")
				common.LogTxtHandler().Debug("Converted files will be placed in the target path within an image name-based sub folder.")
				common.LogTxtHandler().Debug("This is the post-conversion target path (used for vmPathName to vCenter).")

				postConvTargetPath = targetPath + imageName
				postConvTargetPath = common.CheckAddSlashToPath(postConvTargetPath)
			} else {
				common.LogTxtHandler().Debug("File type is either VMTX or VMX...")
				common.LogTxtHandler().Debug("Target path is the same as the post-conversion target path.")

				postConvTargetPath = targetPath
				// since we're grabbing the sourceFolderPath regardless of type, we can use this VMTX postConvert value as it will be the same
			}
		} else {
			common.LogTxtHandler().Debug("---> Path type is LINUX-based.")
			imageFileName, sourceFolderPath = common.FileNamePathFromLnx(sourcePath)		// Ex: /lab/rhel9/rhel9.ova, returns: rhel9.ova, /lab/rhel9/
			imageName = common.ParseFilenameForImageName(imageFileName)		                // Ex: rhel9.ova, returns rhel9
			fileType  = common.GetFileType(imageFileName)

			if fileType == "ova" || fileType == "ovf" {
				common.LogTxtHandler().Debug("File type is either OVA or OVF...")
				common.LogTxtHandler().Debug("Converted files will be placed in the target path within an image name-based sub folder.")
				common.LogTxtHandler().Debug("This is the post-conversion target path (used for vmPathName to vCenter).")

				postConvTargetPath = targetPath + imageName
				postConvTargetPath = common.CheckAddSlashToPath(postConvTargetPath)
			} else {
				common.LogTxtHandler().Debug("File type is either VMTX or VMX...")
				common.LogTxtHandler().Debug("Target path is the same as the post-conversion target path.")

				postConvTargetPath = targetPath
			}
		}

		common.LogTxtHandler().Info("Image Filename: " + imageFileName)
		common.LogTxtHandler().Info("Image Name: " + imageName)
		common.LogTxtHandler().Info("File Type: " + fileType)
		common.LogTxtHandler().Info("Source Path: " + sourcePath)
		common.LogTxtHandler().Info("Target Path (before conversion): " + targetPath)
		common.LogTxtHandler().Info("Source Folder Path: " + sourceFolderPath)
		common.LogTxtHandler().Info("Post Conversion Target Path: " + postConvTargetPath)

		// If this is an OVF image, we need to first move the image files into a sub dir called "ovf_files" and update the conversion source path to here
		// If not, we'll get a file conflict with the disk file(s)
		if fileType == "ovf" {
			common.LogTxtHandler().Info("OVF file detected...")
			common.LogTxtHandler().Info("Moving OVF files into subdirectory of source path called 'ovf_files'...")
			common.LogTxtHandler().Debug("OVF files contain disk files named the same as the resulting disk files. If we don't move them, there will be a file conflict during conversion.")
			destDir := sourceFolderPath + "ovf_files"		              // ex: 'E:\\path\\to\\win2022\\ovf_files'
			destDir = common.CheckAddSlashToPath(destDir)                 // add ending slash by os type; 'E:\\path\\to\\win2022\\ovf_files\\'
			moveList, err := vm.SetOvfFileList(sourcePath)                // Get list of OVF files to move
			err = common.MoveFiles(moveList, sourceFolderPath, destDir)   // [file list], 'E:\\path\\to\\win2022\\', 'E:\\path\\to\\win2022\\ovf_files\\'
			if err != nil {
				strErr := fmt.Sprintf("%v", err)
				common.LogTxtHandler().Error("Error moving files: " + strErr)
			} else {
				common.LogTxtHandler().Info("Files moved successfully!")
			}

			// Setting the new conversion source path to the ovf_file dir for the conversion process only
			newSourcePath := destDir + imageFileName				        // ex: E:\\path\\to\\win2022\\ovf_files\\win2022.ovf
			common.LogTxtHandler().Debug("New Source Path for Image Files (after OVF files moved): " + newSourcePath)
			common.LogTxtHandler().Info("Checking image type and converting if necessary. This may time some time...")
			convertResult = vm.ConvertImageByType(fileType, newSourcePath, targetPath)
				
		} else {  // ova and vmtx don't need to be moved first
			common.LogTxtHandler().Info("OVA, VMTX, or VMX files detected...")
			common.LogTxtHandler().Info("Checking image type and converting if necessary. This may time some time...")
			convertResult = vm.ConvertImageByType(fileType, sourcePath, targetPath)
		}

		common.LogTxtHandler().Info("--> Conversion Result: " + convertResult)

		if convertResult != "Failed" {
			common.LogTxtHandler().Info("Setting vmPathName....")
			// dsImagePath - represents datastore path to image without any mount point info from Packer server
			if dsImagePath == "" {
				if fileType == "ova" || fileType == "ovf" { // don't do for vmx or vmtx files as full path to file is already being passed in those cases
				postConvTargetFilePath = postConvTargetPath + imageName + ".vmx"
				vmPathName = vm.SetVmPathName(postConvTargetFilePath, dsName)
				} else {
					vmPathName = vm.SetVmPathName(postConvTargetPath, dsName)
				}
			} else {
				common.LogTxtHandler().Info("Datastore image pathing set; using this to set vmPathName...")
				dsImagePath = common.CheckAddSlashToPath(dsImagePath)   // Ex: /dev-servers/
				dsImagePath = dsImagePath + imageName					// Ex: /dev-servers/ub24
				dsImagePath = common.CheckAddSlashToPath(dsImagePath)   // Ex: /dev-servers/ub24/
				dsImagePath = dsImagePath + imageFileName				// Ex: /dev-servers/ub24/ub24.ovf
				vmPathName = vm.SetVmPathName(dsImagePath, dsName)		// Will rename file in path to be VMX
				// dsImagePath = /dev-servers/ub24/ub24.vmx
				// vmPathName = [dsName] dev-servers/ub24/ub24.vmx
			}

			common.LogTxtHandler().Info("vmPathName: " + vmPathName)
			common.LogTxtHandler().Info("Beginning import into vCenter....")

			vcToken := common.VcenterAuth(vcUser, vcPass, vcServer)
			statusCode := vm.RegisterVm(vcToken, vcServer, dcName, vmPathName, imageName, folderId, resPoolId)
			common.LogTxtHandler().Info("Status Code of Register VM task: " + statusCode)
	
			if statusCode == "200" {
				common.LogTxtHandler().Info("Import successful. Marking image as a VM Template...")
				tempResult := govmomi.MarkAsTemplate(vcUser, vcPass, vcServer, imageName, dcName)
				common.LogTxtHandler().Info(tempResult)
	
				if strings.Contains(tempResult, "Success") {
					common.LogTxtHandler().Info("The image import and template conversion completed successfully.")
					return "Success"
				} else {
					common.LogTxtHandler().Error("Error: Unable to import and/or convert the image into a VM Template.")
					return "Failed"
				}
			} else {
				common.LogTxtHandler().Error("Error registering VMX file with vCenter.")
				return "Failed"
			}
		} else {
			common.LogTxtHandler().Error("Error during image type check and file conversion process.")
			return "Failed"
		}
	}
	return "Success"
}