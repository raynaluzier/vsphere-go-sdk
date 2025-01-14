package main

import (
	"fmt"
	"os"

	_ "github.com/raynaluzier/vsphere-go-sdk/common"
	"github.com/raynaluzier/vsphere-go-sdk/util"
	_ "github.com/raynaluzier/vsphere-go-sdk/vm"
)

func main() {
	vCenter := os.Getenv("VCENTER_SERVER")
	vcUser  := os.Getenv("VCENTER_USER") 
	vcPass  := os.Getenv("VCENTER_PASSWORD")
	dcName := os.Getenv("DATACENTER_NAME")
	dsName := os.Getenv("DATASTORE_NAME")

	util.VcenterServer = vCenter
	util.VcUser = vcUser
	util.VcPassword = vcPass
	//===================================================

	fmt.Println(dcName)
	fmt.Println(dsName)

}
