package archive

import (
	"fmt"

	"github.com/helloyi/go-sshclient"
	"github.com/raynaluzier/vsphere-go-sdk/common"
)

// verified key has pass, but get error: "ssh: key is not password protected"
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
