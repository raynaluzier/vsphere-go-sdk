package tasks

import (
	"fmt"
	"os"

	"github.com/raynaluzier/vsphere-go-sdk/common"
	"github.com/raynaluzier/vsphere-go-sdk/govmomi"
	"github.com/raynaluzier/vsphere-go-sdk/vm"
)

// Download to datastore happens from artifactory sdk, called from plugin
// Check file type and conversion/rename happens in vsphere sdk, called from plugin
// RegisterVm and convert to template happens in vsphere sdk, called from plugin

func ImportVm(vcUser, vcPass, vcServer, dcName, dsName, imageName, folderName, resPoolName, clusterName string) string {
	token := common.VcenterAuth(vcUser, vcPass, vcServer)

	folderId, err := govmomi.GetFolderId(vcUser, vcPass, vcServer, folderName, dcName)
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error getting folder ID: %s\n", err)
        os.Exit(1)
	}

	resPoolId, err := govmomi.GetResPoolId(vcUser, vcPass, vcServer, resPoolName, dcName, clusterName)
	if err != nil {
        fmt.Fprintf(os.Stderr, "Error getting resource pool ID: %s\n", err)
        os.Exit(1)
	}

	// missing check file and convert process...
	statusCode := vm.RegisterVm(token, vcServer, dcName, dsName, imageName, folderId, resPoolId)
	fmt.Println("Status Code of Register VM task: ", statusCode)

	result := govmomi.MarkAsTemplate(vcUser, vcPass, vcServer, imageName, dcName)

	return result
}