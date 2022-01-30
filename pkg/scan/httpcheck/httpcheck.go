package httpcheck

import (
	log2 "BeeScan-scan/pkg/log"
	"fmt"
	"github.com/fatih/color"
	"github.com/valyala/fasthttp"
	"time"
)

/*
创建人员：云深不知处
创建时间：2022/1/7
程序功能：http检测主机是否存在
*/

// HttpCheck HTTP检测主机存活
func HttpCheck(domain string, port string, ip string) bool {
	var url string
	log2.Info("[HttpCheck]:", ip)
	fmt.Fprintln(color.Output, color.HiCyanString("[INFO]"), "["+time.Now().Format("2006-01-02 15:04:05")+"]", "[HttpCheck]:", ip)
	if port == "80" {
		url = "http://" + domain
	} else {
		url = "http://" + domain + ":" + port
	}
	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)
	req.Header.SetMethod("GET")
	req.SetRequestURI(url)
	err := fasthttp.DoTimeout(req, nil, 3*time.Second)
	if err != nil {
		log2.Warn("[HttpCheck]:", err)
		fmt.Fprintln(color.Output, color.HiYellowString("[WARN]"), "["+time.Now().Format("2006-01-02 15:04:05")+"]", "[HttpCheck]:", err)
		return false
	} else {
		return true
	}

}
