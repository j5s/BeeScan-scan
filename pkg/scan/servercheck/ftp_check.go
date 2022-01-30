package servercheck

import (
	"fmt"
	"github.com/jlaffaye/ftp"
	"time"
)

/*
创建人员：云深不知处
创建时间：2022/1/5
程序功能：ftp检测单元
*/


func FTPCheck(ip string,port string,user string,pwd string) bool {
	client, err := ftp.Dial(fmt.Sprintf(`%s:%s`, ip, port), ftp.DialWithTimeout(6*time.Second))
	if err != nil {
		fmt.Println(err)
		return false
	}

	err = client.Login(user, pwd)
	if err != nil {
		fmt.Println(err)
		return false
	}

	defer client.Quit()

	return true
}
