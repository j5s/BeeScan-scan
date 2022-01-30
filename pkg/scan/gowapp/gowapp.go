package gowapp

import (
	log2 "BeeScan-scan/pkg/log"
	"BeeScan-scan/pkg/runner"
	"BeeScan-scan/pkg/scan/httpcheck"
	"embed"
	"fmt"
	"github.com/fatih/color"
	gowap "github.com/jiaocoll/GoWapp/pkg/core"
	"os"
	"time"
)

/*
创建人员：云深不知处
创建时间：2022/1/14
程序功能：Wappalyzer模块
*/

type TargetInfo struct {
	Urls         []Urls         `json:"urls"`
	Technologies []Technologies `json:"technologies"`
}
type Urls struct {
	URL    string `json:"url"`
	Status int    `json:"status"`
}
type Categories struct {
	ID   int    `json:"id"`
	Slug string `json:"slug"`
	Name string `json:"name"`
}
type Technologies struct {
	Slug       string       `json:"slug"`
	Name       string       `json:"name"`
	Confidence int          `json:"confidence"`
	Version    string       `json:"version"`
	Icon       string       `json:"icon"`
	Website    string       `json:"website"`
	Cpe        string       `json:"cpe"`
	Categories []Categories `json:"categories"`
}

func GowappConfig() *gowap.Config {
	//Create a Config object and customize it
	wapconfig := gowap.NewConfig()
	//Timeout in seconds for fetching the url
	wapconfig.TimeoutSeconds = 5
	//Timeout in seconds for loading the page
	wapconfig.LoadingTimeoutSeconds = 5
	//Don't analyze page when depth superior to this number. Default (0) means no recursivity (only first page will be analyzed)
	wapconfig.MaxDepth = 0
	//Max number of pages to visit. Exit when reached
	wapconfig.MaxVisitedLinks = 5
	//Delay in ms between requests
	wapconfig.MsDelayBetweenRequests = 200
	//Choose scraper between rod (default) and colly
	wapconfig.Scraper = "rod"
	//Override the user-agent string
	wapconfig.UserAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/97.0.4692.71 Safari/537.36"
	//Output as a JSON string
	wapconfig.JSON = false
	return wapconfig
}

func GowappInit(f embed.FS) (*gowap.Wappalyzer, error) {
	wapconfig := GowappConfig()
	wapp, err := gowap.Init(wapconfig, f)
	if err != nil {
		log2.Error("[GoWappInit]:", err)
		fmt.Fprintln(color.Output, color.HiRedString("[ERRO]"), "["+time.Now().Format("2006-01-02 15:04:05")+"]", "[GoWappInit]:", err)
		os.Exit(1)
	}
	return wapp, nil
}

// GoWapp Wappalyzer识别模块
func GoWapp(r *runner.Output, wapp *gowap.Wappalyzer) *gowap.Output {
	if httpcheck.HttpCheck(r.Domain, r.Port, r.Ip) {
		fmt.Fprintln(color.Output, color.HiCyanString("[INFO]"), "["+time.Now().Format("2006-01-02 15:04:05")+"]", "[GoWapp]:", r.Ip)
		var fullUrl string
		targetinfo := &gowap.Output{}
		protocol := "http"
		if r.Domain != "" {
			fullUrl = fmt.Sprintf("%s://%s:%s/", protocol, r.Domain, r.Port)
		} else {
			fullUrl = fmt.Sprintf("%s://%s:%s/", protocol, r.Ip, r.Port)
		}
		res, _ := wapp.Analyze(fullUrl)
		if res != nil {
			targetinfo = res.(*gowap.Output)
		}
		return targetinfo
	}
	return nil
}
