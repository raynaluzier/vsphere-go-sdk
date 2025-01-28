# Govmomi Functions
Govmomi-based authentication needs to be added to each Govmomi-based function because the required context is not available when passed from function to function.

A commented-out copy of the authentication piece is available at the top of the Govmomi package for reference.


## GetResPoolId
A vSphere resource pool ID is required when importing the VMX into vCenter before marking it as a template, even if just the default resource pool is used. If `resPoolName` and `clusterName` are left blank (""), the default resource pool will be used. If not enough information is provided to properly derive the resource pool ID, then the default resource pool will be used instead.

This function is used in support of the `RegisterVm` function to gather the necessary information to pass into payload.

#### Inputs
| Name           | Description                                                                                   | Type     | Required |
|----------------|-----------------------------------------------------------------------------------------------|----------|:--------:|
| user           | Username of the vCenter account that will be performing the operations (ex: jdoe@domain.com)  | string   | TRUE     |
| pass           | Password for the provided vCenter account                                                     | string   | TRUE     |
| server         | FQDN or IP address of the target vCenter server                                               | string   | TRUE     |
| resPoolName    | Target resource pool name where the template will be imported                                 | string   | TRUE     |
| datacenterName | Target datacenter name where the template will be imported                                    | string   | TRUE     |
| clusterName    | Target cluster name where the template will be imported                                       | string   | TRUE     |

#### Outputs
| Name       | Description                                              | Type    |
|------------|----------------------------------------------------------|---------|
| resPoolId  | vSphere ID of the target resource pool                   | string  |
| err        | Error if there's any issue getting the resource pool ID  | error   |


## GetFolderId
A vSphere folder ID is required when importing the VMX into vCenter before marking it as a template, even if just the root folder is used. If `folderName` is left blank (""), the root folder will be used instead. 

This function is used in support of the `RegisterVm` function to gather the necessary information to pass into payload.

#### Inputs
| Name           | Description                                                                                   | Type     | Required |
|----------------|-----------------------------------------------------------------------------------------------|----------|:--------:|
| user           | Username of the vCenter account that will be performing the operations (ex: jdoe@domain.com)  | string   | TRUE     |
| pass           | Password for the provided vCenter account                                                     | string   | TRUE     |
| server         | FQDN or IP address of the target vCenter server                                               | string   | TRUE     |
| folderName     | Target resource pool name where the template will be imported                                 | string   | TRUE     |
| datacenterName | Target datacenter name where the template will be imported                                    | string   | TRUE     |

#### Outputs
| Name      | Description                                       | Type    |
|-----------|---------------------------------------------------|---------|
| folderId  | vSphere ID of the target vCenter folder           | string  |
| err       | Error if there's any issue getting the folder ID  | error   |


## MarkAsTemplate
Takes in the vCenter credential information, image name (which is used to find the virtual machine - i.e. the VMX that we converted the image to), and the datacenter where the virtual machine (VMX) resides and marks the VM as a VM template (VMTX). 

For instances where the image that was downloaded originated as a VMTX, and then was renamed to VMX before the RegisterVm function, this process is required because we cannot import a VMTX file. It has to be a VMX file and THEN we can convert it back to a VM template (VMTX).

#### Inputs
| Name           | Description                                                                                   | Type     | Required |
|----------------|-----------------------------------------------------------------------------------------------|----------|:--------:|
| user           | Username of the vCenter account that will be performing the operations (ex: jdoe@domain.com)  | string   | TRUE     |
| pass           | Password for the provided vCenter account                                                     | string   | TRUE     |
| server         | FQDN or IP address of the target vCenter server                                               | string   | TRUE     |
| folderName     | Target resource pool name where the template will be imported                                 | string   | TRUE     |
| datacenterName | Target datacenter name where the template will be imported                                    | string   | TRUE     |

#### Outputs
| Name       | Description                                                                                  | Type     |
|------------|----------------------------------------------------------------------------------------------|----------|
| (result)   | Returns string of either "Success" or "Failure" of the resulting template conversion action  | string   |