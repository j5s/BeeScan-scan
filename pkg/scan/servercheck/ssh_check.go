package servercheck

import (
	"fmt"
	"golang.org/x/crypto/ssh"
	"net"
	"time"
)

/*
创建人员：云深不知处
创建时间：2022/1/5
程序功能：ssh检测单元
*/


func SSHCheck(ip string,port string,user string,pwd string) bool {
	result := false
	config := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{ssh.Password(pwd)},
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
		Timeout: 6 * time.Second,
	}

	client, err := ssh.Dial("tcp", fmt.Sprintf(`%s:%s`, ip, port), config)
	if err == nil {
		defer client.Close()
		session, err := client.NewSession()
		if err == nil {
			errEcho := session.Run("echo BeeScan")
			if errEcho == nil {
				defer session.Close()
				result = true
			}
		}

	}
	return result
}