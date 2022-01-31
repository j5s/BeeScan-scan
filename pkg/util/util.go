package util

import (
	"fmt"
	"os"
	"time"
)

/*
创建人员：云深不知处
创建时间：2022/1/3
程序功能：工具包
*/

func StrInSlice(i string, array []string) bool {
	ret := false
	for _, a := range array {
		if i == a {
			ret = true
			break
		}
	}
	return ret
}

// Removesamesip 去重函数
func Removesamesip(ips []string) (result []string) {
	result = make([]string, 0)
	tempMap := make(map[string]bool, len(ips))
	for _, e := range ips {
		if tempMap[e] == false {
			tempMap[e] = true
			result = append(result, e)
		}
	}
	return result
}

// DaySub 天数差
func DaySub(BeforeData string) int {
	current := time.Now().Unix()

	loc, _ := time.LoadLocation("Local") //获取时区
	tmp, _ := time.ParseInLocation("2006-01-02 15:04:05", BeforeData, loc)
	timestamp := tmp.Unix() //转化为时间戳 类型是int64

	res := (current - timestamp) / 86400 //相差值
	return int(res)

}

// HourSub 小时差
func HourSub(BeforeData string) int {
	current := time.Now().Unix()

	loc, _ := time.LoadLocation("Local") //获取时区
	tmp, _ := time.ParseInLocation("2006-01-02 15:04:05", BeforeData, loc)
	timestamp := tmp.Unix() //转化为时间戳 类型是int64

	res := (current - timestamp) / 3600 //相差值
	return int(res)
}

func MinSub(BeforeData string) int {
	current := time.Now().Unix()

	loc, _ := time.LoadLocation("Local") //获取时区
	tmp, _ := time.ParseInLocation("2006-01-02 15:04:05", BeforeData, loc)
	timestamp := tmp.Unix() //转化为时间戳 类型是int64

	res := (current - timestamp) / 60 //相差值
	return int(res)
}

// FileExist 判断所给路径文件/文件夹是否存在
func FileExist(path string) bool {
	_, err := os.Stat(path) //os.Stat获取文件信息
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}

func StrToSlince(SourceData []string) string {
	var res string
	for k, v := range SourceData {
		if k == 0 {
			res += fmt.Sprintf("%s ", v)
		} else {
			res += fmt.Sprintf("%s ", v)
		}
	}
	return res
}
