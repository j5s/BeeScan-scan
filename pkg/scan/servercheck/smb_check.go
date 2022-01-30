package servercheck

import (
	"github.com/stacktitan/smb/smb"
	"log"
	"strconv"
)

/*
创建人员：云深不知处
创建时间：2022/1/5
程序功能：smb检测单元
*/



func SMBCheck(ip string,port string,user string,pwd string) bool {
	flag := false
	p,err1 := strconv.Atoi(port)
	if err1 != nil{
		log.Println(err1)
	}
	options := smb.Options{
		Host:        ip,
		Port:        p,
		User:        user,
		Password:    pwd,
		Domain:      "",
		Workstation: "",
	}

	session, err := smb.NewSession(options, false)
	if err == nil {
		session.Close()
		if session.IsAuthenticated {
			flag = true
		}
	}
	return flag
}