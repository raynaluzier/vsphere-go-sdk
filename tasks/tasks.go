package tasks

import (
	"fmt"

	"github.com/raynaluzier/vsphere-go-sdk/common"
	"github.com/raynaluzier/vsphere-go-sdk/govmomi"
	"github.com/raynaluzier/vsphere-go-sdk/vm"
)

// download happens from artifactory sdk, called from plugin
// Check file type and conversion/rename happens in vsphere sdk, called from plugin
// true/false + copy to datastore happens in vsphere sdk, called from plugin

func ImportVm(vcUser, vcPass, vcServer, dcName, dsName, imageName, folderName string) {
	token := common.VcenterAuth(vcUser, vcPass, vcServer)

	folderId := govmomi.GetFolderId(vcUser, vcPass, vcServer, folderName, dcName)

	statusCode := vm.RegisterVm(token, vcServer, dcName, dsName, imageName, folderId)
	fmt.Println(statusCode)

	result := govmomi.MarkAsTemplate(vcUser, vcPass, vcServer, imageName, dsName)
	fmt.Println(result)
}