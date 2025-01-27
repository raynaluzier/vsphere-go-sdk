package archive

import (
	"fmt"

	"github.com/helloyi/go-sshclient"
	"github.com/raynaluzier/vsphere-go-sdk/common"
)

/*
var authInput map[string]string  //mapstructure "auth"
auth = {
	method = "user_pass"    // "private_key"
	user = "jsmith"
	pass = "pass12345!"
	server = "server123.domain.com"
	port = "22"  // must be string
	//private_key_file
}*/

/*
authInput := make(map[string]string)
authInput["method"] = "private_key"
authInput["user"] = "rt-user"
authInput["server"] = "192.168.1.xxx"
authInput["port"] = "22"
authInput["private_key_file"] = "C:\\Users\\me\\.ssh\\id_ecdsa"
*/

/*
authInput["method"] = "user_pass"
authInput["user"] = "domain.local\\someuser"
authInput["pass"] = "xyz123abc"
authInput["port"] = "22"
authInput["server"] = "192.168.1.xxx"
*/

var client *sshclient.Client

func GetAuthClient(authInput map[string]string) (*sshclient.Client) {
	var user, pass, server, port, key string
	var client *sshclient.Client
	if len(authInput) > 0 {
		if authInput["method"] == "user_pass" {
			user   = authInput["user"]
			pass   = authInput["pass"]
			server = authInput["server"]
			port   = authInput["port"]				
			
			common.LogTxtHandler().Debug("User: " + user)
			common.LogTxtHandler().Debug("Server: " + server)
			common.LogTxtHandler().Debug("Port: " + port)
			client = AuthUserPass(user, pass, server, port)
			
		} else if authInput["method"] == "private_key" {
			user   = authInput["user"]
			server = authInput["server"]
			port   = authInput["port"]
			key    = authInput["private_key_file"]
			
			common.LogTxtHandler().Debug("User: " + user)
			common.LogTxtHandler().Debug("Server: " + server)
			common.LogTxtHandler().Debug("Port: " + port)
			client = AuthPrivateKey(server, port, user, key)
			
		} else {
			// error: unrecognized method
		}
	}
	return client
}


func AuthUserPass(user, pass, server, port string) (*sshclient.Client) {
	if port == "" {
		port = "22"
	}
	sshServer := server + ":" + port
	client, err := sshclient.DialWithPasswd(sshServer, user, pass)
	if err != nil {
		strErr := fmt.Sprintf("%v\n", err)
		common.LogTxtHandler().Error("Error authenticating with username/password - " + strErr)
	} else {
		strClient := fmt.Sprintf("%v\n", err)
		common.LogTxtHandler().Debug(strClient)
	}
	
	return client
}

func AuthPrivateKey(server, port, user, privKeyFile string) (*sshclient.Client) {
	if port == "" {
		port = "22"
	}
	sshServer := server + ":" + port
	client, err := sshclient.DialWithKey(sshServer, user, privKeyFile)
	if err != nil {
		strErr := fmt.Sprintf("%v\n", err)
		common.LogTxtHandler().Error("Error authenticating with private key file - " + strErr)
	} else {
		strClient := fmt.Sprintf("%v\n", err)
		common.LogTxtHandler().Debug(strClient)
	}
	
	return client
}

// NOTE ---> verified key has pass, but get error: "ssh: key is not password protected"
func AuthPrivateKeyPhrase(server, port, user, privKeyFile, passphrase string) (*sshclient.Client) {
	sshServer := server + ":" + port
	client, err := sshclient.DialWithKeyWithPassphrase(sshServer, user, privKeyFile, passphrase)
	if err != nil {
		strErr := fmt.Sprintf("%v\n", err)
		common.LogTxtHandler().Error("Error authenticating with private key file and passphrase - " + strErr)
	} else {
		strClient := fmt.Sprintf("%v\n", err)
		common.LogTxtHandler().Debug(strClient)
	}
	
	return client
}
