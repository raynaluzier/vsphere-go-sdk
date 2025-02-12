package main

import (
	_ "log"
	"os"
	_ "strings"

	_ "github.com/raynaluzier/vsphere-go-sdk/common"
	_ "github.com/raynaluzier/vsphere-go-sdk/govmomi"
	"github.com/raynaluzier/vsphere-go-sdk/util"
	_ "github.com/raynaluzier/vsphere-go-sdk/vm"
	//"golang.org/x/crypto/ssh" // go get golang.org/x/crypto/ssh
)

func main() {
	vcServer := os.Getenv("VCENTER_SERVER")
	vcUser  := os.Getenv("VCENTER_USER") 
	vcPass  := os.Getenv("VCENTER_PASSWORD")
	//dcName := os.Getenv("VCENTER_DATACENTER")
	//dsName := os.Getenv("VCENTER_DATASTORE")
	//folderName := os.Getenv("VCENTER_FOLDER")
	//resPoolName := os.Getenv("VCENTER_RESOURCE_POOL")
	//clusterName := os.Getenv("VCENTER_CLUSTER")
	//imageName := os.Getenv("IMAGE_NAME")

	logLevel 	:= os.Getenv("LOGGING")
	outputDir 	:= os.Getenv("OUTPUTDIR")

	util.VcServer = vcServer
	util.VcUser = vcUser
	util.VcPassword = vcPass
	util.Logging   = logLevel
	util.OutputDir = outputDir
	//===================================================

}
