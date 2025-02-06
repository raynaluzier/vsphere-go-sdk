# VM Functions

## SetPathsFromDownloadUri
This function is part of the workflow where the image is downloaded from Artifactory, image files converted to VMX, before being imported into vCenter and then marking it as a VM Template (VMTX).

It takes in the provided output directory and download URI. The download URI is parsed for the image name and image file type (OVA, OVF, or VMTX), and the output directory is used to set the full file path to the source and target locations. It is assumed that the converted image will reside in the same location as the originating image files that were downloaded from Artifactory. So for example:

    Source Path:  E:\lab-servs\win22\win22.ova
    Target Path:  E:\lab-servs\win22\win22.vmx (and associated standard VM files)

#### Inputs
| Name        | Description                                                                 | Type     | Required |
|-------------|-----------------------------------------------------------------------------|----------|:--------:|
| outputDir   | Accessible datastore path where the image files were downloaded             | string   | TRUE     |
| downloadUri | Artifactory download URI of the image; determines image name and image type | string   | TRUE     |

#### Outputs
| Name       | Description                                                                                                     | Type     |
|------------|-----------------------------------------------------------------------------------------------------------------|----------|
| fileType   | Type of image file (OVA, OVF, or VMTX)                                                                          | string   |
| sourcePath | Full file path to the source image file being converted (set from `outputDir`)                                  | string   |
| targetPath | Full file path to the target image VMX file (and associated files) that will result after the image conversion  | string   |


## ConvertImageByType
Takes in the image's file type, source directory where the current image file(s) reside, the target directory where the converted image files will be placed, and the `fileType` (ova, ovf, or vmtx) which will determine the type image conversion process that needs to take place to get to type VMX. For example:

    File Type:  ova
    Source Path:  E:\lab-servs\win22\win22.ova
    Target Path:  E:\lab-servs\win22\win22.vmx (and associated standard VM files)

#### Inputs
| Name       | Description                                                                                                    | Type     | Required |
|------------|----------------------------------------------------------------------------------------------------------------|----------|:--------:|
| fileType   | Type of image file (OVA, OVF, or VMTX)                                                                         | string   | TRUE     |
| sourcePath | Full file path to the source image file being converted (set from `outputDir`)                                 | string   | TRUE     |
| targetPath | Full file path to the target image VMX file (and associated files) that will result after the image conversion | string   | TRUE     |

#### Outputs
| Name    | Description                                                       | Type     |
|---------|-------------------------------------------------------------------|----------|
| result  | Resulting string result of the conversion ("Success" or "Failed") | string   |


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
