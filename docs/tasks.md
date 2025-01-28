# Tasks Functions
As described previously, these functions are intended to be used with a custom Packer plugin, but can be called independently if desired.

## ImportVm
Takes in the information about the target vCenter environment and placement, and builds the JSON payload for the API call. The credential information is used to first authenticate to vCenter via the `VcenterAuth` function and gets the needed token (which is stripped of its surrounding quotes). Next, the vSphere ID of the target folder is found using the `GetFolderId` function, and then the vSphere ID of the target resource pool is found using the `GetResPoolId` function.

If successful, the vCenter information along with the authentication token is passed to the `RegisterVm` function to convert the image file to VMX, register it with vCenter, and then mark it as a VM Template.

#### Inputs
| Name        | Description                                                                              | Type     | Required |
|-------------|------------------------------------------------------------------------------------------|----------|:--------:|
| vcUser      | Username of the vCenter account that will be performing the operations (ex: jdoe@domain.com)  | string   | TRUE     |
| vcPass      | Password for the provided vCenter account                                                     | string   | TRUE     |
| vcServer    | FQDN or IP address of the target vCenter server                                               | string   | TRUE     |
| dcName      | Target datacenter name where the template will be imported                                    | string   | TRUE     |
| dsName      | Datastore name where the image files are located                                              | string   | TRUE     |
| imageName   | Name of the image (no extension or path info)                                                 | string   | TRUE     |
| folderName  | Target folder name where the template will be imported                                        | string   | TRUE     |
| resPoolName | Target resource pool name where the template will be imported                                 | string   | TRUE     |
| clusterName | Target cluster name where the template will be imported                                       | string   | TRUE     |

#### Outputs
| Name     | Description                                                                                  | Type     |
|----------|----------------------------------------------------------------------------------------------|----------|
| result   | Returns string of either "Success" or "Failure" of the resulting template conversion action  | string   |

