package ping

import (
	log2 "BeeScan-scan/pkg/log"
	"bytes"
	"os/exec"
	"runtime"
	"strings"
)

/*
创建人员：云深不知处
创建时间：2022/1/1
程序功能：ping检测主机是否存活
*/

// PingCheckAlive PING检测主机存活
func PingCheckAlive(host string) bool {
	log2.Info("[PingCheck]:", host)
	sysType := runtime.GOOS
	if sysType == "windows" {
		cmd := exec.Command("ping", "-n", "2", host)
		var output bytes.Buffer
		cmd.Stdout = &output
		cmd.Run()
		if strings.Contains(output.String(), "TTL=") && strings.Contains(output.String(), host) {
			return true
		}
	} else if sysType == "linux" || sysType == "darwin" {
		cmd := exec.Command("ping", "-c", "2", host)
		var output bytes.Buffer
		cmd.Stdout = &output
		cmd.Run()
		if strings.Contains(output.String(), "ttl=") && strings.Contains(output.String(), host) {
			return true
		}
	}
	return false
}
