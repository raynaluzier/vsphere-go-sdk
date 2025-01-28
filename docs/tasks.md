# Tasks Functions
As described previously, these functions are intended to be used with a custom Packer plugin, but can be called independently if desired.

## ImportVm


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

