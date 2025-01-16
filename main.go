package main

import (
	"os"

	"github.com/raynaluzier/vsphere-go-sdk/util"
)

func main() {
	vCenter := os.Getenv("VCENTER_SERVER")
	vcUser  := os.Getenv("VCENTER_USER") 
	vcPass  := os.Getenv("VCENTER_PASSWORD")
	//dcName := os.Getenv("DATACENTER_NAME")
	//dsName := os.Getenv("DATASTORE_NAME")

	util.VcenterServer = vCenter
	util.VcUser = vcUser
	util.VcPassword = vcPass
	//===================================================


}
