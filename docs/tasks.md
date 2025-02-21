# Tasks Functions
As described previously, these functions are intended to be used with a custom Packer plugin to perform larger, specific operations.

## GetResourceIds
Takes in a variety of required vCenter-related inputs to get and return the resource ID of the target vCenter folder and resource pool.

#### Inputs
| Name        | Description                                                                                  | Type     | Required |
|-------------|----------------------------------------------------------------------------------------------|----------|:--------:|
| vcUser      | Username of the vCenter account that will be performing the operations (ex: jdoe@domain.com)`| string   | TRUE     |
| vcPass      | Password for the provided vCenter account                                                    | string   | TRUE     |
| vcServer    | FQDN or IP address of the target vCenter server                                              | string   | TRUE     |
| dcName      | Name of the target datacenter in vCenter; will use default datacenter if blank               | string   | FALSE    |
| folderName  | Name of the target folder in vCenter; will use default root folder if blank                  | string   | FALSE    |
| resPoolName | Name of the target resource pool in vCenter; will use default pool if blank                  | string   | FALSE    |
| clusterName | Name of the target cluster in vCenter; will use default DC and pool if left blank            | string   | FALSE    |

#### Outputs
| Name      | Description                                       | Type     |
|-----------|---------------------------------------------------|----------|
| folderId  | Resource ID of the target vCenter folder          | string   |
| resPoolId | Resource ID of the target vCenter resource pool   | string   |
| err       | If applicable, resulting error for the operation  | string   |


## ConvertImportFromDownload
When the `import_no_download` flag is set to FALSE (which is default), the desired image (OVA, OVF, or VMTX) will first be downloaded to a designated output directory in a previous step, within a sub-folder named after the image name. If successul, this function will take in a number of vCenter-related inputs, the originally provided output directory, and the Artifactory download URI of the image. **The output directory should be an accessible datastore directory as this is where vCenter will import the resulting VMX from.**

Using this information, we'll get the required authentication token from vCenter, the correct source and target paths will be set, the image will be converted into VMX format, the `vmPathName` will be set as appropriate based on the target path (formatting depends on the image type), the VMX will be imported into vCenter, and then marked as a VM template (converted back to VMTX).

The result of this process is returned as a string status of either "Success" or "Failed".

#### Inputs
| Name        | Description                                                                                             | Type     | Required |
|-------------|---------------------------------------------------------------------------------------------------------|----------|:--------:|
| vcUser      | Username of the vCenter account that will be performing the operations (ex: jdoe@domain.com)`           | string   | TRUE     |
| vcPass      | Password for the provided vCenter account                                                               | string   | TRUE     |
| vcServer    | FQDN or IP address of the target vCenter server                                                         | string   | TRUE     |
| outputDir   | Properly escaped directory to where image files were be downloaded (without the image named sub-folder) | string   | TRUE     |
| downloadUri | Artifactory download URI address for the image (OVA, OVF, or VMTX)                                      | string   | TRUE     |
| dcName      | Name of the target datacenter in vCenter                                                                | string   | TRUE     |
| dsName      | Name of the target datastore in vCenter                                                                 | string   | TRUE     | 
| imageName   | Name of the image; i.e. the image file without the extension                                            | string   | TRUE     |
| folderId    | Resource ID of the target vCenter folder                                                                | string   | TRUE     |
| resPoolId   | Resource ID of the target vCenter resource pool                                                         | string   | TRUE     |

#### Outputs
| Name      | Description                       | Type     |
|-----------|-----------------------------------|----------|
| (result)  | Returns "Success" or "Failed"     | string   |


## ConvertImportNoDownload
When the `import_no_download` flag is set to TRUE, the desired image (OVA, OVF, or VMTX) is considered to already be downloaded to a designated output directory, in this case considered `sourcePath`. **Source path should be an accessible datastore directory as this is where vCenter will import the resulting VMX from.** This function will take in a number of vCenter-related inputs and the source path where the downloaded images reside.

Using this information, we'll get the required authentication token from vCenter and then the correct target paths will be set (the format of which is determined by the type of image). If the file type is OVA or OVF, when the OVFTool runs, the converted image files will be placed in the target path within an image named-based sub-folder. This is a function of the OVFTool itself and not within our control. This resulting path is considered the post-conversion target path, which is used to formulate the `vmPathName` used by vCenter during the import process.

If the image is either VMTX or already in VMX format, the target path is the same as the post-conversion target path. The OVFTool isn't involved in this case and taking a VMTX to VMX is simply a matter of renaming the file.

Additionally, if the image type is OVF specifically, a sub-directory will be created under the source path called `ovf_files` and the OVF files will be moved their first, this path will be used as the conversion source, and the resulting image files will be placed in the post-conversion target path, as described above. When an OVF package is unpacked, it results in one or more disk files that are named the same as the disk file(s) that are included in the OVF package. If we left them in the same directory, when the unpacking occurred, there would be a file conflict and the conversion would stall out.

If the conversion process succeeds, the `vmPathName` is set, the VMX will be imported into vCenter, and then marked as a VM template (converted back to VMTX). 

The result of this process is returned as a string status of either "Success" or "Failed".

#### Inputs
| Name        | Description                                                                                  | Type     | Required |
|-------------|----------------------------------------------------------------------------------------------|----------|:--------:|
| vcUser      | Username of the vCenter account that will be performing the operations (ex: jdoe@domain.com)`| string   | TRUE     |
| vcPass      | Password for the provided vCenter account                                                    | string   | TRUE     |
| vcServer    | FQDN or IP address of the target vCenter server                                              | string   | TRUE     |
| dcName      | Name of the target datacenter in vCenter                                                     | string   | TRUE     |
| dsName      | Name of the target datastore in vCenter                                                      | string   | TRUE     |
| sourcePath  | Properly escaped datastore directory where the image files reside                            | string   | TRUE     |
| folderId    | Resource ID of the target vCenter folder                                                     | string   | TRUE     |
| resPoolId   | Resource ID of the target vCenter resource pool                                              | string   | TRUE     |

#### Outputs
| Name      | Description                       | Type     |
|-----------|-----------------------------------|----------|
| (result)  | Returns "Success" or "Failed"     | string   |