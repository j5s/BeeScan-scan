package util

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
