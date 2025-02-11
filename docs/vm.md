# VM Functions

## SetPathsFromDownloadUri
This function is called when the `import_no_download` flag is set to FALSE; meaning the image artifacts are downloaded prior to conversion and import into vCenter.

The function takes in the provided output directory and download URI. The download URI is parsed for the image name and image file type (OVA, OVF, or VMTX), and the output directory is used to set the source and target path locations. The source path matches the directory path where the image artifacts were downloaded. The target path will ultimately be this same directory. How it's passed within the process is determined by the image file type.  

--> For OVA and OVF files, the OVFTOOL automatically places converted files into a folder based on the image's name (ex: win2022.ova will be placed into /win2022/win2022.ova). The download process also places the image files in a folder based on the image's name. Therefore, we will set the target path to match the `outputDir` so when the conversion process runs, it will use the same image folder that was used during the download process so the image files are all in the same directory for administrative neatness.

For example:  outputDir = "E:\lab-servs", and we're downloading the image 'win2022.ova' from Artifactory...
* The process will download 'win2022.ova' into "E:\lab-servs\win2022\win2022.ova". <-- This is the `sourcePath` that will be passed to the conversion process in function `ConvertImageByType`.
* The `targetPath` we will use is "E:\lab-servs". When the conversion runs, it will automatically place the converted image files in "E:\lab-servs\win2022\" so all of the image files are in the same directory.

This is the path used to set the `vmPathName` used within the payload that gets sent in the `RegisterVm` function. 

--> For VMTX files, the `sourcePath` process is the same. However, the `targetPath` used is the same as the sourcepath except the file extension is 'VMX' here. VMTX files use the `RenameFile` function instead during the conversion process.

#### Inputs
| Name        | Description                                                                 | Type     | Required |
|-------------|-----------------------------------------------------------------------------|----------|:--------:|
| outputDir   | Accessible datastore path where the image files were downloaded             | string   | TRUE     |
| downloadUri | Artifactory download URI of the image; determines image name and image type | string   | TRUE     |

#### Outputs
| Name       | Description                                                                                                        | Type     |
|------------|--------------------------------------------------------------------------------------------------------------------|----------|
| fileType   | Type of image file (OVA, OVF, or VMTX)                                                                             | string   |
| sourcePath | Full file path to the source image files being converted (set from `outputDir`+ imageName + fileName)              | string   |
| targetPath | Depending on image type, target path of conversion matching either the `outputDir` or full path to target VMX file | string   |


## SetPathNoDownload
This function is called when the `import_no_download` flag is set to TRUE; meaning the image artifacts WILL NOT be downloaded prior to conversion and import into vCenter.

The function takes in the source path which includes the filename and extension of the image, and it's checked for the type of OS platform it belongs to. Based on Windows or Linux path type, the path is parsed for the image's file name, image file type (OVA, OVF, or VMTX), and base image name (ex: win2022). We assume the target path will ultimately be the same as the source path post-conversion for neatness, but how it's formed and passed is dependent on the image file type. We also assume that the source path to the image will include a folder structure such as 'image_name/image_name.ext' (because that's the structure we use for downloads and also the typical organization structure used in datastores to host VMs). 

--> For OVA and OVF files, the OVFTOOL automatically places converted files into a folder based on the image's name (ex: win2022.ova will be placed into /win2022/win2022.ova). The download process also places the image files in a folder based on the image's name. Therefore, we will strip off the 'image_name/image_name.ext' from the source path and set this value as the target path.

For example:  sourcePath = "E:\lab-servs\win2022\win2022.ova"
* This is the same `sourcePath` that will be passed to the conversion process in function `ConvertImageByType`.
* We'll strip off 'win2022\win2022.ova' so the `targetPath` we will use is "E:\lab-servs". When the conversion runs, it will automatically place the converted image files in "E:\lab-servs\win2022\" so all of the image files are in the same directory.
* If for some reason you are using a non-typical pathing scheme (for example: E:\path\somefolder\somefile.ova), we will detect this and pass 'E:\lab-servs' as the target path. **However, the OVFTOOL will place the converted files into 'E:\lab-servs\somefile\' instead automatically. So you will have the source files in 'E:\path\somefolder\' and the converted files in 'E:\lab-servs\somefile\'. This is a function of the OVFTOOL outside of our control that cannot be changed.** 

--> For VMTX files, the `sourcePath` is the same. However, the `targetPath` used is the same as the sourcepath except the file extension is 'VMX' here. VMTX files use the `RenameFile` function instead during the conversion process.

#### Inputs
| Name       | Description                                                       | Type     | Required |
|------------|-------------------------------------------------------------------|----------|:--------:|
| sourcePath | Full file path to the source images file being converted          | string   | TRUE     |

#### Outputs
| Name       | Description                                                                                                                      | Type     |
|------------|----------------------------------------------------------------------------------------------------------------------------------|----------|
| targetPath | Depending on image type, full path to target VMX file or base folder path of source path without the image's parent folder/file  | string   |


## ConvertImageByType
Takes in the image's file type, source directory where the current image file(s) reside, the target directory where the converted image files will be placed, and the `fileType` (ova, ovf, or vmtx) which will determine the type image conversion process that needs to take place to get to type VMX. For OVA and OVF files, the target directory will be the same as the `outputDir` provided. For example:

    File Type:  ova
    Source Path:  E:\lab-servs\win22\win22.ova
    Target Path:  E:\lab-servs\ --> OVFTOOL automatically places the files in E:\lab-servs\win22\

    File Type:  ovf
    Source Path:  E:\lab-servs\win22\win22.ovf
    Target Path:  E:\lab-servs\ --> OVFTOOL automatically places the files in E:\lab-servs\win22\

    File Type:  vmtx
    Source Path:  E:\lab-servs\win22\win22.vmtx
    Target Path:  E:\lab-servs\win22\win22.vmx  

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


## SetVmPathName
Takes in the 'sourcePath' that's a result of `SetPathsFromDownloadUri` or as a direct input if the 'import_no_download' flag is set to TRUE. The 'sourcePath' would be something like "E:\\labimage\\labimage.ova" or "/labimage/labimage.ova". This function also takes in the datastore name (dsName).

Since the VMX files are assumed to be in the same directory as the originally downloaded images files, the sourcePath is reformatted to create the `vmPathName` that will be used by the `RegisterVm` function as the datastore path. The OS-platform, taken by the source path formatting, is taken into account. 

For example: sourcePath = "E:\\labimage\\labimage.ova" and dsName = "lab-servs". The result of this function would be:  vmPathName = "[lab-servs] labimage/labimage.vmx".

Likewise, sourcePath = "E:\\windows-servers\\labimage22\\labimage22.ova" and dsName = "lab-servs". The result of this function would be:  vmPathName = "[lab-servs] windows-servers/labimage/labimage.vmx".

#### Inputs
| Name       | Description                                                                                        | Type     | Required |
|------------|----------------------------------------------------------------------------------------------------|----------|:--------:|
| sourcepath | Full directory path with filename of the source image file(s) (before conversion; OVA, OVF, VMTX)  | string   | TRUE     |
| dsName     | Name of the datastore where the image files reside                                                 | string   | TRUE     |

#### Outputs
| Name        | Description                                                                               | Type     |
|-------------|-------------------------------------------------------------------------------------------|----------|
| vmPathName  | Properly formatted datastore name and file path to VMX file used for import into vCenter  | string   |


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
