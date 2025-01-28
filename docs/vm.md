# VM Functions

## CheckFileConvert
Takes the provided output directory (datastore location) and download URI for the primary image file (OVA, OVF, or VMTX), parses the image name from the download URI, and determines the source file path. The file type is checked and depending on whether it's OVA, OVF, or VMTX, it's converted to a VMX as appropriate. This is in preparation for importing the template (as a virtual machine) into vCenter and marking it as a VM Template.

#### Inputs
| Name         | Description                                                                  | Type     | Required |
|--------------|------------------------------------------------------------------------------|----------|:--------:|
| outputDir    | Connected datastore location where the downloaded image files reside         | string   | TRUE     |
| downloadUri  | Download URI for the primary image file (OVA, OVF, or VMTX) from Artifactory | string   | TRUE     |

#### Outputs
| Name        | Description                                                | Type     |
|-------------|------------------------------------------------------------|----------|
| (result)    | Result of conversion process; either "Success" or "Failed" | string   |


## RegisterVm
Takes in the information about the target vCenter environment and placement, and builds the JSON payload for the API call. The authentication token is stripped of its surrounding quotes before the API call imports the target VMX file into vCenter.

This function is used in conjunction with the `MarkAsTemplate` function.

#### Inputs
| Name      | Description                                                                       | Type     | Required |
|-----------|-----------------------------------------------------------------------------------|----------|:--------:|
| token     | vCenter authentication token; quotes will be removed from string before API call  | string   | TRUE     |
| vcServer  | FQDN or IP address of the target vCenter server                                   | string   | TRUE     |
| dcName    | Target datacenter name where the template will be imported                        | string   | TRUE     |
| dsName    | Datastore name where the image files are located                                  | string   | TRUE     |
| imageName | Name of the image (no extension or path info)                                     | string   | TRUE     |
| folderId  | vSphere ID of the target folder where the template will be placed                 | string   | TRUE     |
| resPoolId | vSphere ID of the target resource pool where the template will be placed          | string   | TRUE     |

#### Outputs
| Name        | Description                                                                                           | Type     |
|-------------|-------------------------------------------------------------------------------------------------------|----------|
| statusCode  | Resulting status code returned from the import process; 200 if successful or specific code otherwise  | string   |


## ConvertOvfaToVmx
**Requires the OVFTool be installed on the machine that's executing the conversion commands.**

Takes in the path (inputPath) to the OVA or OVF file to be converted and the output path where the resulting VMX and associated files should be placed. The input can be a local path or URL.

Ensure local paths are escaped properly (ex: 'C:\\\lab\\\file.vmx').

#### Inputs
| Name       | Description                                                                     | Type     | Required |
|------------|---------------------------------------------------------------------------------|----------|:--------:|
| inputPath  | Source path of the OVA/OVF files; this can be a local path or URL to the image  | string   | TRUE     |
| outputPath | Destination path of the converted image files                                   | string   | TRUE     |

#### Outputs
| Name        | Description                                             | Type     |
|-------------|---------------------------------------------------------|----------|
| (result)    | Returns "Failed" or "Success" after conversion process  | string   |


## ConvertVmxToOvfa
**Requires the OVFTool be installed on the machine that's executing the conversion commands.**

Takes in the path (inputPath) to the VMX and associated files to be converted and the output path where the resulting OVA or OVF files should be placed. 

Ensure local paths are escaped properly (ex: C:\\lab\\file.vmx).

#### Inputs
| Name       | Description                                                                                | Type     | Required |
|------------|--------------------------------------------------------------------------------------------|----------|:--------:|
| inputPath  | Source path of the VMX and associated files; this can be a local path or URL to the image  | string   | TRUE     |
| outputPath | Destination path of the converted OVA or OVF image files                                   | string   | TRUE     |

#### Outputs
| Name        | Description                                             | Type     |
|-------------|---------------------------------------------------------|----------|
| (result)    | Returns "Failed" or "Success" after conversion process  | string   |
