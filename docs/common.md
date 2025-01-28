# Common Functions

## VcenterAuth
Takes in the credential and server information for a vCenter account with the adequate level of permissions to authenticate with vCenter, register new VMs, convert VMs into templates, query for resources such as datacenters, datastores, resource pools, and folders to get the object's ID. The result is a token that will be used to authorize subsequent vSphere API calls.

#### Inputs
| Name    | Description                                                                                   | Type    | Required |
|---------|-----------------------------------------------------------------------------------------------|---------|:--------:|
| user    | Username of the vCenter account that will be performing the operations (ex: jdoe@domain.com)  | string  | TRUE     |
| pass    | Password for the provided vCenter account                                                     | string  | TRUE     |
| server  | FQDN or IP address of the target vCenter server                                               | string  | TRUE     |

#### Outputs
| Name   | Description                                                                                    | Type     |
|--------|------------------------------------------------------------------------------------------------|----------|
| token  | Return result from the vCenter authentication request; used to authorize subsequent API calls  | string   |


## AddUrlProtocol
Looks at the server address and if it's missing 'https://', the function adds this. If the server address already contains the http/https protocol, the existing name is returned. This supports the `VcenterAuth` function to ensure the provided server address is formatted for the API call.

#### Inputs
| Name       | Description                                        | Type     | Required |
|------------|----------------------------------------------------|----------|:--------:|
| server     | FQDN or IP address of the target vCenter server    | string   | TRUE     |

#### Outputs
| Name                | Description                        | Type     |
|---------------------|------------------------------------|----------|
| serverUrl / server* | URL of the target vCenter server   | string   |
* If server address already contains http/https protocol


## TrimUrlProtocol
Checks for and if exists, trims off 'http://' or 'https://' from the URL to get the server address. This is primarily used to support the Govmomi package functions that require the username and password of the vCenter account passed in as part of the server URL to authenticate.

#### Inputs
| Name       | Description                       | Type    | Required |
|------------|-----------------------------------|---------|:--------:|
| serverUrl  | URL of the target vCenter server  | string  | TRUE     |

#### Outputs
| Name                | Description                                      | Type     |
|---------------------|--------------------------------------------------|----------|
| server / serverUrl* | FQDN or IP address of the target vCenter server  | string   |
* If the http/https protocol isn't found in the server address, the server address is left as-is and returned


## TrimQuotes
This function removes the inherent quotes from the string and returns the result. This required when passing the vCenter `vmware-api-session-id` token in subsequent vSphere API calls otherwise the call will fail.

#### Inputs
| Name | Description                                 | Type    | Required |
|------|---------------------------------------------|---------|:--------:|
| s    | String that should have its quotes removed  | string  | TRUE     |

#### Outputs
| Name | Description                           | Type    |
|------|---------------------------------------|---------|
| s    | Resulting string with quotes removed  | string  |


## Rename File
Takes in the full path to the file and renames it to the new file path. This function supports renaming a VMTX file (VM Template) to VMX file (virtual machine) so that it can be imported into vCenter before we convert it back to a template.

#### Inputs
| Name         | Description                                   | Type     | Required |
|--------------|-----------------------------------------------|----------|:--------:|
| oldFilePath  | Full path to the file that will be renamed    | string   | TRUE     |
| newFilePath  | Full path to the resulting file after rename  | string   | TRUE     |

#### Outputs
| Name     | Description                          | Type    |
|----------|--------------------------------------|---------|
| (result) | Returns either "Failed" or "Success" | string  |


## GetFileType
Takes in the path to a file and extracts the file extension. The file extension is then returned. 

This supports the `CheckFileConvert` function which uses an image's primary Artifactory download URI to determine the image type and then converts it to a VMX as appropriate in prep to import it into vCenter.

#### Inputs
| Name      | Description                                                                          | Type    | Required |
|-----------|--------------------------------------------------------------------------------------|---------|:--------:|
| filePath  | Path to the image file (extracted from the image's primary Artifactory download URI) | string  | TRUE     |

#### Outputs
| Name | Description                | Type     |
|------|----------------------------|----------|
| ext  | Extracted file extension   | string   |


## ParseUriForFilename
Takes in either the image's Artifactory artifact URI or download URI address for the primary image file (OVA, OVF, or VMTX) and parses it for the file name.

#### Inputs
| Name         | Description                                                 | Type     | Required |
|--------------|-------------------------------------------------------------|----------|:--------:|
| artifactUri  | An image's Artifactory artifact or download URI address     | string   | TRUE     |

#### Outputs
| Name        | Description                                                 | Type    |
|-------------|-------------------------------------------------------------|---------|
| fileName    | Resulting file name parsed from the artifact's URI address  | string  |


## ParseFilenameForImageName
Takes in the image's primary filename (OVA, OVF, or VMTX file) and strips off the extension to get the resulting image name.

#### Inputs
| Name       | Description                                                | Type     | Required |
|------------|------------------------------------------------------------|----------|:--------:|
| fileName   | Primary filename of a given image (OVA, OVF, or VMTX file) | string   | TRUE     |

#### Outputs
| Name        | Description                                     | Type     |
|-------------|-------------------------------------------------|----------|
| imageName   | Name of the image extracted from the file name  | string   |


## CheckPathType
Checks the provided path to see if it's Unix-based (has '/') or Windows-based (has '\'). This is often used in combination with `CheckAddSlashToPath` to add the appropriate ending slash type to given path if needed.

#### Inputs
| Name  | Description                                          | Type   | Required |
|-------|------------------------------------------------------|--------|:--------:|
| path  | Path to provided directory; such as Output Directory | string | TRUE     |

#### Outputs
| Name       | Description                                        | Type |
|------------|----------------------------------------------------|------|
| isWinPath | Returns true of the provided path is Windows-based  | bool |


## CheckAddSlashToPath
Used with `CheckPathType`; based on path type (Windows vs. Unix), checks the provided path to see if it ends with the appropriate back or forward slashes. If not present, the function will add a slash as appropriate to the platform type. This ensures the output directory path provided is formatted as required.

#### Inputs
| Name       | Description                                              | Type   | Required |
|------------|----------------------------------------------------------|--------|:--------:|
| inputStr   | String that was provided through some kind of user input | string | TRUE     |
| actualStr  | String pulled from actual object name                    | string | TRUE     |

#### Outputs
| Name       | Description                                               | Type |
|------------|-----------------------------------------------------------|------|
| ture/false | Returns true if compared string match, regardless of case | bool |


## SetLoggingLevel
Uses the Global Variable `util.Logging` to set the desired logging level and returns the slog.Level equivalent value to be used by the desired logging handlers (LogTxtHandler or LogJsonHandler). If not specified, logging level defaults to INFO.

#### Inputs
Takes no inputs

#### Outputs
| Name      | Description                                                                                | Type       |
|-----------|--------------------------------------------------------------------------------------------|------------|
| logLevel  | Will be slog.LevelInfo, slog.LevelWarn, slog.LevelError, or slog.LevelDebug based on input | slog.Level |


## LogTxtHandler
Takes in the appropriate logging level type from `SetLoggingLevel()` and sets the level in the handler options. Then a new Text handler interface is created with the specified logging format and defines where they are written to, in this case, Stdout.

Output example:  `time=2024-12-02T10:35:41.267-07:00 level=INFO msg="This is your info message."`

#### Inputs
| Name        | Description                                             | Type     | Required |
|-------------|---------------------------------------------------------|----------|:--------:|
| LOGGING     | Desired log level; Accepts: INFO, WARN, ERROR, DEBUG    | string   | FALSE    |

#### Outputs
| Name      | Description                                                                                | Type       |
|-----------|--------------------------------------------------------------------------------------------|------------|
| logLevel  | Will be slog.LevelInfo, slog.LevelWarn, slog.LevelError, or slog.LevelDebug based on input | slog.Level |

#### Usage
Example: `someLogLevel := common.SetLoggingLevel()`
         `common.LogTxtHandler(someLogLevel).Info("Info stuff. All is well!")`
         `common.LogTxthandler(someLogLevel).Debug("Found object: test-artifact.txt")`

If 'INFO' is set in .env, then only the .Info, .Warn, and .Error logs will be output.
If 'WARN' is set, then only .Warn, and .Error logs will be output.
If 'ERROR' is set, then only .Error logs will be output.
If 'DEBUG' is set, then all logs - .Info, .Warn, .Error, and .Debug - will be output.


## LogJsonHandler
Takes in the appropriate logging level type from `SetLoggingLevel()` and sets the level in the handler options. Then a new JSON handler interface is created with the specified logging format and defines where they are written to, in this case, Stdout. The JSON handler is useful for parsing and performing other actions based on the output, or writing to an external logging system.

Output example:  `{"time":"2024-12-02T10:13:31.252815-07:00","level":"INFO","msg":"Some JSON Info message."}`

#### Inputs
| Name        | Description                                             | Type     | Required |
|-------------|---------------------------------------------------------|----------|:--------:|
| LOGGING     | Desired log level; Accepts: INFO, WARN, ERROR, DEBUG    | string   | FALSE    |

#### Outputs
| Name      | Description                                                                                | Type       |
|-----------|--------------------------------------------------------------------------------------------|------------|
| logLevel  | Will be slog.LevelInfo, slog.LevelWarn, slog.LevelError, or slog.LevelDebug based on input | slog.Level |

#### Usage
Example: `someLogLevel := common.SetLoggingLevel()`
         `common.LogJsonHandler(someLogLevel).Info("Info stuff. All is well!")`
         `common.LogJsonhandler(someLogLevel).Debug("Found object: test-artifact.txt")`

If 'INFO' is set in .env, then only the .Info, .Warn, and .Error logs will be output.
If 'WARN' is set, then only .Warn, and .Error logs will be output.
If 'ERROR' is set, then only .Error logs will be output.
If 'DEBUG' is set, then all logs - .Info, .Warn, .Error, and .Debug - will be output.