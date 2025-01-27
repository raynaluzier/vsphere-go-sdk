# vsphere-go-sdk

## Summary
This SDK is a collection of Golang functions used to interact with VMware vSphere/vCenter. While the function calls can be called independently, the intent is to use them with a custom Packer plugin integration.

## Pre-Requisites
To run functions from this module, the following pre-requisites must be met:

1. A fully configured vCenter environment with appropriately configured hosts, network, storage, etc. 

2. An account (such as a service account) on vCenter that has access to authenticate with vCenter, register new VMs, convert VMs into templates, query for resources such as datacenters, datastores, resource pools, and folders to get the object's ID. If the account doesn't have enough permissions or permissions to the desired resources, then the functions will not be able to return the necessary information to complete the operations and the process will fail.

3. Go is installed on the system where the functions/plugin will be run: https://go.dev/doc/install.

4. Ensure the `GOPATH` (the local directory to where the Go app is installed) is configured, as well as the environment variable to the path of the Go binary (on Windows, this would typically be `C:\Program Files\Go\bin`, for example).

5. Populating the Global Variables (see: `/util/util.go`): **VcServer**, **VcUser**, **VcPassword**, **OutputDir**, and optionally, **Logging**. This can be done when calling a function in the `tasks` package (e.g. as part of a plugin operation), by configuring the `.env` file, or statically (only recommended for testing).

Using .env File: Configure the `.env` file with vCenter credentials, vCenter server, and then the Output Directory is used for downloading the image file(s) to ensure they're placed in the desired datastore location. Logging provides an option to change the level of logging to display. Logging provides an option to change the level of logging to display.

    * `VCENTER_SERVER`   - vCenter Server   --> Ex: VCENTER_SERVER=vc01.domain.com
    * `VCENTER_USER`     - vCenter User     --> Ex: VCENTER_USER=admin@domain.com
    * `VCENTER_PASSWORD` - vCenter Password --> Ex: VCENTER_PASSWORD=P@s$w0rd123!
    * `OUTPUTDIR`        - Output directory on the desired datastore for downloading artifacts --> Ex: OUTPUTDIR=/servs/path/ or H:\servs\path
    * `LOGGING`          - Logging level (INFO, WARN, ERROR, DEBUG); defaults to 'INFO'        --> Ex: LOGGING=DEBUG

Then use `os.Getenv` to set `util.VcServer`, `util.VcUser`, `util.VcPasword`, `util.OutputDir`, and `util.Logging` respectively.

## About
This SDK is broken into several packages: `common`, `govmomi`, `vm`, `tasks`, and `util` based on the underlying behavior of the functions. Some functions are specifically related to certain behaviors so they have been grouped together into packages as described below.

### Common
These functions perform small, generalized supporting tasks for the other focused modules. These functions can be found under the `common.go` file.

### Govmomi
These functions make use of the [Govmomi](https://github.com/vmware/govmomi) package to establish an authentication client and using that client to gather specific vCenter resource information that is required to perform other functions and tasks.

### Tasks
These functions are larger operations that first set the global variables, and then make a series of function calls to perform specific activities. While they can be called independently, they were created in support of a custom Packer plugin to streamline passing environment-specific variables, such as the vCenter credentials, server, logging, and output directory. Rather than passing one or more of these to every function in the SDK (in addition to the required inputs), they are passed in ONCE to the desired function, the global variables are set, and then they are used automatically when calling each sub-function without having to pass them in over and over.

These larger tasks also group the targeted functions of a desired behavior into a single operation and keep the plugin code to a minimum and simplify performing that desired behavior.

### Util
This is a list of the global variables used within this SDK. As with any Go package, they can be used by importing the `util` package path and then referencing them as `util.Token`, `util.ServerApi`, etc.

### VM
These functions are related to VM-related operations such as converting from OVA/OVF to VMX, converting VMX to OVA/OVF, registering a VM in vCenter (we may also refer to this as 'importing'), and checking the image's file type from the primary download URI (from Artifactory) provided and either converting it or renaming it in preparation for importing into vCenter.

### Archive
Archive also exists as a package, but it's really just a place to hold potentially useful functions that were created but have no immediate use. Artifacts are split into archive files that match their associated behaviors; as of now, either the `archive-ssh-auths.go` or `archive-files.go` files. 

The `archive-ssh-auths.go` file contains functions that are perform different types of authentication methods. The original intention was to include an option to authenticate between systems and copy a list of files (based on image type) over to a datastore prior to image import into vCenter. This has been tabled in favor of using a pre-created share on the target datastore as the output directory when downloading images. 

The `archive-files.go` file contains functions related to file copy operations related to images, including making a list images files to copy, a Windows-based copy function, a placeholder for a Linux-based copy function, and a test SSH function that was used for script development. These were intended to be used with the SSH authentication functions and thus tabled in favor of the pre-created shared.

## Function Reference
A reference outline of each function's behavior and any special notes can be found in the corresponding documents below.

- [Common](https://github.com/raynaluzier/vsphere-go-sdk/blob/main/docs/common.md)
- [Govmomi](https://github.com/raynaluzier/vsphere-go-sdk/blob/main/docs/govmomi.md)
- [Tasks](https://github.com/raynaluzier/vsphere-go-sdk/blob/main/tasks/tasks.go)
- [VM](https://github.com/raynaluzier/vsphere-go-sdk/blob/main/docs/vm.md)
