package runner

import (
	"BeeScan-scan/pkg/httpx"
	"BeeScan-scan/pkg/scan/fringerprint"
	"BeeScan-scan/pkg/scan/gonmap"
	"encoding/json"
	gowap "github.com/jiaocoll/GoWapp/pkg/core"
)

/*
创建人员：云深不知处
创建时间：2022/1/1
程序功能：扫描结果
*/

type Result interface {
	STR() string
	JSON() string
}

type Output struct {
	Ip         string        `json:"ip"`
	TargetId   string        `json:"target_id"`
	Port       string        `json:"port"`
	Protocol   string        `json:"protocol"`
	Domain     string        `json:"domain"`
	Webbanner  FingerResult  `json:"webbanner"`
	Servers    gonmap.Result `json:"servers"`
	CityId     int64         `json:"cityId"`
	Country    string        `json:"country"`
	Region     string        `json:"region"`
	Province   string        `json:"province"`
	City       string        `json:"city"`
	ISP        string        `json:"isp"`
	Servername string        `json:"servername"`
	Wappalyzer *gowap.Output `json:"wappalyzer"`
	Banner     string        `json:"banner"`
	Target     string        `json:"target"`
	LastTime   string        `json:"lastTime"`
}

type FingerResult struct {
	Title         string                  `json:"title"`
	ContentLength int                     `json:"content-length"`
	TLSData       *httpx.TLSData          `json:"tls,omitempty"`
	StatusCode    int                     `json:"status-code"`
	ResponseTime  string                  `json:"response-time"`
	CDN           string                  `json:"cdn"`
	Fingers       fringerprint.FofaPrints `json:"fingers"`
	Str           string                  `json:"str"`
	Header        string                  `json:"header"`
	FirstLine     string                  `json:"firstLine"`
	Headers       map[string][]string     `json:"headers"`
	DataStr       string                  `json:"datastr"`
}

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

func (r *FingerResult) JSON() string {
	if js, err := json.Marshal(r); err == nil {
		return string(js)
	}

	return ""
}
func (r *FingerResult) STR() string {
	return r.Str
}
