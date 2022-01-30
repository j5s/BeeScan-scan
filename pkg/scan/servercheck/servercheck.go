package servercheck

import (
	"io/ioutil"
	"log"
	"os"
	"strings"
)

/*
创建人员：云深不知处
创建时间：2022/1/6
程序功能：检测流程
*/


func ServerCheck(ip string, port string,filename1 string,filename2 string) string {

	var User []string
	var Pwd []string
	file1, err1 := os.Open(filename1)
	if err1 != nil{
		log.Println(err1)
	}
	file2, err2 := os.Open(filename2)
	if err2 != nil{
		log.Println(err2)
	}
	content1,err3 := ioutil.ReadAll(file1)
	if err3 != nil{
		log.Println(err3)
	}
	content2,err4 := ioutil.ReadAll(file2)
	if err4 != nil{
		log.Println(err4)
	}
	User = strings.Split(string(content1),"\n")
	Pwd = strings.Split(string(content2),"\n")
	tmp := []string{"mysql","postgres","mssql"}

	for _,user := range User{
		for _,pwd := range Pwd{
			if SSHCheck(ip,port,user,pwd){
				return "SSH"
			}else if SNMPCheck(ip, port, pwd) {
				return "SNMP"
			}else if SMBCheck(ip, port, user, pwd) {
				return "SMB"
			}else if RedisCheck(ip, port, user, pwd) {
				return "Redis"
			}else if MongoDBCheck(ip, port, user, pwd) {
				return "MongoDB"
			}else if FTPCheck(ip, port, user, pwd) {
				return "FTP"
			}else if ElasticSearchCheck(ip, port, user, pwd) {
				return "ElasticSearch"
			}
		}
	}
	for _,hostType := range tmp{
		for _,user := range User{
			for _,pwd := range Pwd{
				if SQLCheck(hostType, ip, port, user, pwd) {
					return hostType
				}
			}
		}
	}
	return ""
}