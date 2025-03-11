package govmomi

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"regexp"
	"strings"

	"github.com/raynaluzier/vsphere-go-sdk/common"

	"github.com/vmware/govmomi" //go get github.com/vmware/govmomi
	"github.com/vmware/govmomi/find"
	"github.com/vmware/govmomi/object"
	_ "github.com/vmware/govmomi/view"
	_ "github.com/vmware/govmomi/vim25/mo"
)

/*  // Govmomi auth needs to be added to each function; context not available when passed
	// Creating a connection context
    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()

	// Parsing URL
	trimServer := common.TrimUrlProtocol(server) // trims off http:// or https:// off of server name
	sdkUrl := "https://" + user + ":" + pass + "@" + trimServer + ":443/sdk"

	url, err := url.Parse(sdkUrl)  // https://username:password@hostname:443/sdk  @ = %40
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error parsing URL: %s\n", err)
        os.Exit(1)
    }

    // Connecting to vCenter
    client, err := govmomi.NewClient(ctx, url, true)   // shared context, parsed URL, whether client will tolerate an insecure cert
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error connecting to vCenter: %s\n", err)
        os.Exit(1)
    }
*/

func GetResPoolId(user, pass, server, resPoolName, datacenterName, clusterName string) (string, error) {
	// If resPoolName and clusterName are left blank, the default resource pool will be used, which it will also default to if there's
	// not enough info to find the named resource pool
	//--------------------------------------------------------
	ctx, cancel := context.WithCancel(context.Background())
    defer cancel()

	// Parsing URL
	trimServer := common.TrimUrlProtocol(server) // trims off http:// or https:// off of server name
	sdkUrl := "https://" + user + ":" + pass + "@" + trimServer + ":443/sdk"

	url, err := url.Parse(sdkUrl)  // https://username:password@hostname:443/sdk  @ = %40
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error parsing URL: %s\n", err)
        os.Exit(1)
    }
 
    // Connecting to vCenter
    client, err := govmomi.NewClient(ctx, url, true)   // shared context, parsed URL, whether client will tolerate an insecure cert
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error connecting to vCenter: %s\n", err)
        os.Exit(1)
    }
	//--------------------------------------------------------
	var resPoolPath string
	var resPool *object.ResourcePool

	finder   := find.NewFinder(client.Client, true)

	if resPoolName == "" && clusterName == "" {
		dc, err := finder.DefaultDatacenter(ctx)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting default datacenter: %s\n", err)
			os.Exit(1)
		}
		finder.SetDatacenter(dc)
		resPool, err = finder.DefaultResourcePool(ctx)
	} else if resPoolName != "" && clusterName != "" {
		resPoolPath = "/" + datacenterName + "/host/" + clusterName + "/Resources/" + resPoolName
		resPool, err = finder.ResourcePool(ctx, resPoolPath)
	} else if resPoolName == "" && clusterName != "" {
		resPoolPath = "/" + datacenterName + "/host/" + clusterName + "/Resources"
		resPool, err = finder.ResourcePool(ctx, resPoolPath)
	} else {
		common.LogTxtHandler().Info("Not enough information provided to find a specific resource pool.")
		common.LogTxtHandler().Info("Using the default cluster and default resource pool...")
		dc, err := finder.DefaultDatacenter(ctx)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting default datacenter: %s\n", err)
			os.Exit(1)
		}
		finder.SetDatacenter(dc)
		resPool, err = finder.DefaultResourcePool(ctx)
	}
	if err != nil {
        fmt.Fprintf(os.Stderr, "Error getting resource pool ID: %s\n", err)
        os.Exit(1)
	}

	strResPool := resPool.String()
	_, after, _ := strings.Cut(strResPool, ":")					// Returns: 'group-v10 @ /Lab/vm'
	pathAfter := after
	before, _, _ := strings.Cut(pathAfter, "@")					// Returns: 'group-v10 '
	resPoolId := before
	resPoolId = strings.TrimSuffix(resPoolId, " ")				// Returns: 'group-v10'
	fmt.Println("Resource Pool ID: ", resPoolId)
	return resPoolId, err
}

func GetFolderId(user, pass, server, folderName, datacenterName string) (string, error) {
	//--------------------------------------------------------
	ctx, cancel := context.WithCancel(context.Background())
    defer cancel()

	// Parsing URL
	trimServer := common.TrimUrlProtocol(server) // trims off http:// or https:// off of server name
	sdkUrl := "https://" + user + ":" + pass + "@" + trimServer + ":443/sdk"

	url, err := url.Parse(sdkUrl)  // https://username:password@hostname:443/sdk  @ = %40
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error parsing URL: %s\n", err)
        os.Exit(1)
    }
 
    // Connecting to vCenter
    client, err := govmomi.NewClient(ctx, url, true)   // shared context, parsed URL, whether client will tolerate an insecure cert
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error connecting to vCenter: %s\n", err)
        os.Exit(1)
    }
	//--------------------------------------------------------
	// This is the folder where the VM resides
	// If root folder, leave folderName blank
	var folderPath string
	finder   := find.NewFinder(client.Client, true)

	if folderName == "" && datacenterName != "" {
		folderPath = "/" + datacenterName + "/vm"                // Ex: Folder:group-v10 @ /Lab/vm
	} else if folderName != "" && datacenterName != "" {
		folderPath = "/" + datacenterName + "/vm/" + folderName	 // Ex: Folder:group-v1141 @ /Lab/vm/test-vms
	} else if folderName == "" && datacenterName == "" {
		dc, err := finder.DefaultDatacenter(ctx)     			 // Ex: Datacenter:datacenter-3 @ /Lab
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting default datacenter: %s\n", err)
			os.Exit(1)
		}

		strDc := dc.String()
		_, after, _ := strings.Cut(strDc, "@ ")      			// Returns: /Lab
		datacenterName := after
		folderPath = "/" + datacenterName + "/vm"

	} else if folderName != "" && datacenterName == "" {
		common.LogTxtHandler().Info("Not enough information provided to find a specific folder.")
		common.LogTxtHandler().Info("Using the default datacenter and root folder...")
		dc, err := finder.DefaultDatacenter(ctx)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting default datacenter: %s\n", err)
			os.Exit(1)
		}
		strDc := dc.String()
		_, after, _ := strings.Cut(strDc, "@ ")      			// Returns: /Lab
		datacenterName := after
		folderPath = "/" + datacenterName + "/vm"
	}

	folder, err := finder.Folder(ctx, folderPath)
	if err != nil {
		common.LogTxtHandler().Error("Unable to get folder by path: " + folderPath)
	}

	strFolder := folder.String()
	_, after, _ := strings.Cut(strFolder, ":")					// Returns: 'group-v10 @ /Lab/vm'
	pathAfter := after
	before, _, _ := strings.Cut(pathAfter, "@")					// Returns: 'group-v10 '
	folderId := before
	folderId = strings.TrimSuffix(folderId, " ")				// Returns: 'group-v10'
	fmt.Println("Folder ID: ", folderId)
	return folderId, err
}

func MarkAsTemplate(user, pass, server, imageName, datacenterName string) (string) {
	//-----------------------------------------------------
	// Creating a connection context
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	
	// Parsing URL
	trimServer := common.TrimUrlProtocol(server) // trims off http:// or https:// off of server name
	sdkUrl := "https://" + user + ":" + pass + "@" + trimServer + ":443/sdk"
	
	url, err := url.Parse(sdkUrl)  // https://username:password@hostname:443/sdk  @ = %40
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing URL: %s\n", err)
		os.Exit(1)
	}
	
	// Connecting to vCenter
	client, err := govmomi.NewClient(ctx, url, true)   // shared context, parsed URL, whether client will tolerate an insecure cert
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error connecting to vCenter: %s\n", err)
		os.Exit(1)
	}
	//-----------------------------------------------------
	var newVmList []*object.VirtualMachine
	var dc *object.Datacenter

	finder   := find.NewFinder(client.Client, true)

	if datacenterName == "" {
		dc, err = finder.DefaultDatacenter(ctx)
		common.LogTxtHandler().Info("No datacenter specified. Using default: " + dc.String())
	} else {
		dcName := "/" + datacenterName
		dc, err = finder.Datacenter(ctx, dcName)
	}

    if err != nil {
        fmt.Fprintf(os.Stderr, "Error: %s\n", err)
        os.Exit(1)
    }
	finder.SetDatacenter(dc)

	vms, err := finder.VirtualMachineList(ctx, "*")
	if err != nil {
        fmt.Fprintf(os.Stderr, "Error: %s\n", err)
        os.Exit(1)
    }

	for x := range vms {
		matched, err := regexp.MatchString(imageName, vms[x].Name())
		if err != nil {
			break
		}

		if matched == true {
			newVmList = append(newVmList, vms[x])
		}
	}

	if len(newVmList) == 1 {
		vm := newVmList[0]
		vm.MarkAsTemplate(ctx)
		isTemplate, err := vm.IsTemplate(ctx)
		if isTemplate == true {
			common.LogTxtHandler().Info("Successfully converted " + imageName + " to a template.")
			return "Success"
		} else {
			strErr := fmt.Sprintf("%v", err)
			common.LogTxtHandler().Error("Error converting " + imageName + " to a template - " + strErr)
			return "Failure"
		}
	} else if len(newVmList) == 0 {
		common.LogTxtHandler().Error("Unable to find image in vCenter or Datacenter: " + dc.String())
		return "Failure"
	} else {
		common.LogTxtHandler().Error("More than one image named: " + imageName + " was returned.")
		return "Failure"
	}
}
