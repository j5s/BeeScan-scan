package servercheck

import (
	"github.com/gosnmp/gosnmp"
	"log"
	"strconv"
	"time"
)

/*
创建人员：云深不知处
创建时间：2022/1/5
程序功能：snmp检测单元
*/



func SNMPCheck(ip string, port string, pwd string) bool {
	flag := false
	p,err1 := strconv.Atoi(port)
	if err1 != nil{
		log.Println(err1)
	}
	gosnmp.Default.Target = ip
	gosnmp.Default.Port = uint16(p)
	gosnmp.Default.Community = pwd
	gosnmp.Default.Timeout = 4 * time.Second

	err := gosnmp.Default.Connect()
	if err == nil {
		oidList := []string{"1.3.6.1.2.1.1.4.0", "1.3.6.1.2.1.1.7.0"}
		_, err := gosnmp.Default.Get(oidList)
		if err == nil {
			flag = true
		}
	}

	return flag
}