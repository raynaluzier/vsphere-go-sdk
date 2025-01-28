# Archive - Files Functions

## MakeFileCopyList
Takes in the source directory where the image file(s) are located (ensure path is properly escaped 'c:\\\lab\\\'), image name, image type (OVA, OVF, or VMTX), and file suffix (ex: "1.1.0" or "20250128") if used; otherwise leave blank (""). The function will read the source directory, and based on the image type and image name, compile a list of expected files and add them to the `copyList` to be used for optionally copying to a target datastore. 

** The list of expected files assumes typical vSphere file naming conventions using the machine or image name as a base and a series of specific disk files depending on image type, and the current vmware.log file. The list will cycle through and check for up to 15 disk numbers.

The `fileSuffix` option is a unique identifier that should be tacked on at then end of an image name in instances where a standard base name (ex: 'win2022') is used as a base and therefore may not be unique. If not used, pass an empty string ("") in its place. 

#### Inputs
| Name       | Description                                           | Type      | Required |
|------------|-------------------------------------------------------|-----------|:--------:|
| sourceDir  | Source directory where the image file(s) are located  | string    | TRUE     |
| imageName  | Name of the image (no extension or path info)         | string    | TRUE     |
| imageType  | OVA, OVF, or VMTX; case insensitive                   | string    | TRUE     |
| fileSuffix | Optional unique identifier added to a base image name | string    | FALSE    |

#### Outputs
| Name      | Description                                               | Type     |
|-----------|-----------------------------------------------------------|----------|
| copyList  | List of files to copy based on given image type and name  | []string |


## RunScriptTestSSH
Takes in a pre-authenticated SSH client (*sshclient.Client) and builds a simple command-line script to be used on a Windows-based OS. The script executes and outputs the results. This is just used for some simple development/testing.

This is used with the `GetAuthClient` function for authentication either by user/pass or user/private key.

#### Inputs
| Name    | Description                      | Type              | Required |
|---------|----------------------------------|-------------------|:--------:|
| client  | Authenticated SSH client session | *sshclient.Client | TRUE     |

#### Outputs
| Name        | Description                    | Type     |
|-------------|--------------------------------|----------|
| "Complete"  | Marks the end of the function  | string   |


## WinCopyFiles
Takes in a source directory path, target directory path (this should be a datastore accessible to the system), the list of image files to copy, and the SSH client session to copy image files to/from the directory paths between Windows-based systems.

Directory paths should be properly escaped (ex: 'e:\\\lab\\\').

#### Inputs
| Name      | Description                                                                                | Type              | Required |
|-----------|--------------------------------------------------------------------------------------------|-------------------|:--------:|
| sourceDir | Source directory where the image file(s) are to be copied from                             | string            | TRUE     |
| targetDir | Destination directory where the image file(s) should be copied; this should be a datastore | string            | TRUE     |
| copyList  | List of files to copy based image type                                                     | []string          | TRUE     |
| client    | Authenticated SSH client session                                                           | *sshclient.Client | TRUE     |

#### Outputs
| Name     | Description                                                                | Type     |
|----------|----------------------------------------------------------------------------|----------|
| (result) | Resulting string of either "Copy Process Failed" or "End of Copy Process"  | string   |


## LinuxCopyFiles (PLACEHOLDER)

#### Inputs
| Name      | Description                                                               | Type     | Required |
|-----------|---------------------------------------------------------------------------|----------|:--------:|

#### Outputs
| Name     | Description                                                                | Type     |
|----------|----------------------------------------------------------------------------|----------|
