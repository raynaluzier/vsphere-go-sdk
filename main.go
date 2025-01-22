package main

import (
	"fmt"
	"os"

	"github.com/raynaluzier/vsphere-go-sdk/common"
	"github.com/raynaluzier/vsphere-go-sdk/util"
)

func main() {
	vCenter := os.Getenv("VCENTER_SERVER")
	vcUser  := os.Getenv("VCENTER_USER") 
	vcPass  := os.Getenv("VCENTER_PASSWORD")
	//dcName := os.Getenv("VCENTER_DATACENTER")
	//dsName := os.Getenv("VCENTER_DATASTORE")
	//folderName := os.Getenv("VCENTER_FOLDER")
	//resPoolName := os.Getenv("VCENTER_RESOURCE_POOL")
	//clusterName := os.Getenv("VCENTER_CLUSTER")
	imageName := os.Getenv("IMAGE_NAME")

	//logLevel 	:= os.Getenv("LOGGING")
	outputDir 	:= os.Getenv("OUTPUTDIR")

	util.VcenterServer = vCenter
	util.VcUser = vcUser
	util.VcPassword = vcPass
	//util.Logging   = logLevel
	util.OutputDir = outputDir
	//===================================================
	outputDir = "E:\\ub20pkrt-10031746\\"
	sourcePath := outputDir + imageName + ".vmtx"
	newPath    := outputDir + imageName + ".vmx"

	//result := vm.CheckFileConvert(outputDir, downloadUri)
	fmt.Println(common.RenameFile(sourcePath, newPath))

	//if result == "Success" {
	//	fmt.Println(tasks.ImportVm(vcUser, vcPass, vCenter, dcName, dsName, imageName, folderName, resPoolName, clusterName))
	//}
}
