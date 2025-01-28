# Archive - Authentication Functions

## GetAuthClient
Takes in the authentication method ("user_pass" or "private_key"), server name or address to SSH to, port, username, and either passwor or private key file depending on the authentication method defined and calls the corresponding authentication function (`AuthUserPass` or `AuthPrivateKey`) to get an SSH authentication client.

The information is passed via a map[string]string called `authInput` with sample usage provided above the function.

If port number is left blank (""), then the default port number of "22" will be used.

#### Inputs --> Username/Password Authentication
| Name    | Description                                                                       | Type     | Required |
|---------|-----------------------------------------------------------------------------------|----------|:--------:|
| method  | "user_pass"; the auth method to use. Valid options: "user_pass" or "private_key"  | string   | TRUE     |
| user    | Username that has SSH access to the server                                        | string   | TRUE     |
| pass    | Password for the associated user/service account                                  | string   | TRUE     |
| server  | FQDN or IP address of the server to SSH into                                      | string   | TRUE     |
| port    | SSH port number                                                                   | string   | TRUE     |

#### Inputs --> Username/Private Key Authentication
| Name             | Description                                                                         | Type     | Required |
|------------------|-------------------------------------------------------------------------------------|----------|:--------:|
| method           | "private_key"; the auth method to use. Valid options: "user_pass" or "private_key"  | string   | TRUE     |
| user             | Username that has SSH access to the server                                          | string   | TRUE     |
| private_key_file | Path to the private key file that should be used to authenticate                    | string   | TRUE     |
| server           | FQDN or IP address of the server to SSH into                                        | string   | TRUE     |
| port             | SSH port number                                                                     | string   | TRUE     |

#### Outputs
| Name    | Description                      | Type               |
|---------|----------------------------------|--------------------|
| client  | Authenticated SSH client session | *sshclient.Client  |


## AuthUserPass
Performs username/password based authentication to SSH into a target server. 

The intended use was to pair with either the Windows or Linux-based File Copy Script function to copy image files from perhaps a local directory to a datastore before importing the image into vCenter and marking it as a template.

#### Inputs
| Name       | Description                                                             | Type     | Required |
|------------|-------------------------------------------------------------------------|----------|:--------:|
| user       | Username that has SSH access to the server                              | string   | TRUE     |
| pass       | Password for the associated user/service account                        | string   | TRUE     |
| server     | FQDN or IP address of the server to SSH into                            | string   | TRUE     |
| port       | SSH port number; if left blank, the standard SSH port "22" will be used | string   | TRUE     |

#### Outputs
| Name    | Description                      | Type               |
|---------|----------------------------------|--------------------|
| client  | Authenticated SSH client session | *sshclient.Client  |


## AuthPrivateKey
Performs username/private key based authentication to SSH into a target server.

The intended use was to pair with either the Windows or Linux-based File Copy Script function to copy image files from perhaps a local directory to a datastore before importing the image into vCenter and marking it as a template.

#### Inputs
| Name             | Description                                                             | Type     | Required |
|------------------|-------------------------------------------------------------------------|----------|:--------:|
| user             | Username that has SSH access to the server                              | string   | TRUE     |
| private_key_file | Path to the private key file that should be used to authenticate        | string   | TRUE     |
| server           | FQDN or IP address of the server to SSH into                            | string   | TRUE     |
| port             | SSH port number; if left blank, the standard SSH port "22" will be used | string   | TRUE     |

#### Outputs
| Name    | Description                      | Type               |
|---------|----------------------------------|--------------------|
| client  | Authenticated SSH client session | *sshclient.Client  |


## AuthPrivateKeyPhrase
**NOT FUNCTIONING --> Verified the key has a password, but getting error: "ssh: key is not password protected"**

Performs username/private key with passphrase based authentication to SSH into a target server.

The intended use was to pair with either the Windows or Linux-based File Copy Script function to copy image files from perhaps a local directory to a datastore before importing the image into vCenter and marking it as a template.

#### Inputs
| Name             | Description                                                             | Type     | Required |
|------------------|-------------------------------------------------------------------------|----------|:--------:|
| user             | Username that has SSH access to the server                              | string   | TRUE     |
| private_key_file | Path to the private key file that should be used to authenticate        | string   | TRUE     |
| server           | FQDN or IP address of the server to SSH into                            | string   | TRUE     |
| port             | SSH port number; if left blank, the standard SSH port "22" will be used | string   | TRUE     |
| passphrase       | Passphrase that was set when generating the pub/private keys            | string   | TRUE     |

#### Outputs
| Name    | Description                      | Type               |
|---------|----------------------------------|--------------------|
| client  | Authenticated SSH client session | *sshclient.Client  |