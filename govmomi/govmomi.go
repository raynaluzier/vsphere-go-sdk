package govmomi

import (
	"context"
	"fmt"
	"log"
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


func GetFolderId(user, pass, server, folderName, datacenterName string) string {
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

	if folderName == "" {
		folderPath = "/" + datacenterName + "/vm"                // Ex: Folder:group-v10 @ /Lab/vm
	} else {
		folderPath = "/" + datacenterName + "/vm/" + folderName	 // Ex: Folder:group-v1141 @ /Lab/vm/test-vms
	}

	folder, err := finder.Folder(ctx, folderPath)
	if err != nil {
		fmt.Println("Unable to get folder by path: " + folderPath)
	}

	strFolder := folder.String()
	_, after, _ := strings.Cut(strFolder, ":")					// Returns: 'group-v10 @ /Lab/vm'
	pathAfter := after
	before, _, _ := strings.Cut(pathAfter, "@")					// Returns: 'group-v10 '
	folderId := before
	folderId = strings.TrimSuffix(folderId, " ")				// Returns: 'group-v10'
	fmt.Println(folderId)
	return folderId
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
		fmt.Println("No datacenter specified. Using default: " + dc.String())
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
			fmt.Println("Successfully converted " + imageName + " to a template.")
			return "Success"
		} else {
			fmt.Println("Error converting " + imageName + " to a template - ", err)
			return "Failure"
		}
	} else if len(newVmList) == 0 {
		log.Fatal("Unable to find image in vCenter.")
		return "Failure"
	} else {
		log.Fatal("More than one image named: " + imageName + " was returned.")
		return "Failure"
	}
}
