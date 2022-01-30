package runner

import (
	"BeeScan-scan/pkg/httpx"
	log2 "BeeScan-scan/pkg/log"
	"BeeScan-scan/pkg/scan/fringerprint"
	"BeeScan-scan/pkg/scan/httpcheck"
	"fmt"
	"github.com/fatih/color"
	"github.com/projectdiscovery/hmap/store/hybrid"
	"log"
	"net"
	"net/url"
	"strings"
	"time"
)

/*
创建人员：云深不知处
创建时间：2022/1/1
程序功能：
*/

type Runner struct {
	Ip       string
	Port     string
	Domain   string
	Protocol string
	Ht       *httpx.HTTPX
	Hm       *hybrid.HybridMap
	Fofa     *fringerprint.FofaPrints
}

// NewRunner 创建runner实例
func NewRunner(ip string, port string, domain string, protocol string, FofaPrints fringerprint.FofaPrints) (*Runner, error) {
	runner := &Runner{
		Ip:       ip,
		Port:     port,
		Domain:   domain,
		Protocol: protocol,
	}
	if hm, err := hybrid.New(hybrid.DefaultDiskOptions); err != nil {
		log.Fatalf("Could not create temporary input file: %s\n", err)
	} else {
		runner.Hm = hm
	}

	// http
	HttpOptions := &httpx.HTTPOptions{
		Timeout:          3 * time.Second,
		RetryMax:         3,
		FollowRedirects:  true,
		Unsafe:           false,
		DefaultUserAgent: httpx.GetRadnomUserAgent(),
	}
	ht, err := httpx.NewHttpx(HttpOptions)
	if err != nil {
		return nil, err
	}
	runner.Ht = ht
	runner.Fofa = &FofaPrints
	return runner, nil

}

// http请求
func do(r *Runner, fullUrl string) (*httpx.Response, error) {
	req, err := r.Ht.NewRequest("GET", fullUrl)

	if err != nil {
		return &httpx.Response{}, err
	}
	resp, err2 := r.Ht.Do(req)
	return resp, err2
}

// Request FOFA指纹识别
func Request(r *Runner) FingerResult {
	var resp *httpx.Response
	var fullUrl string
	if r.Ht.Dialer != nil {
		defer r.Ht.Dialer.Close()
	}
	defer func(Hm *hybrid.HybridMap) {
		err := Hm.Close()
		if err != nil {
			log2.Warn("[HttpRequest]:", r.Ip)
			fmt.Fprintln(color.Output, color.HiYellowString("[WARN]"), "["+time.Now().Format("2006-01-02 15:04:05")+"]", "[HttpRequest]:", r.Ip)
		}
	}(r.Hm)
	if httpcheck.HttpCheck(r.Domain, r.Port, r.Ip) {

		log2.Info("[HttpRequest]:", r.Ip)
		fmt.Fprintln(color.Output, color.HiCyanString("[INFO]"), "["+time.Now().Format("2006-01-02 15:04:05")+"]", "[HttpRequest]:", r.Ip)

		retried := false
		protocol := httpx.HTTP
	retry:
		if r.Domain != "" {
			fullUrl = fmt.Sprintf("%s://%s:%s", protocol, r.Domain, r.Port)
		} else {
			fullUrl = fmt.Sprintf("%s://%s:%s", protocol, r.Ip, r.Port)
		}

		timeStart := time.Now()

		resp = &httpx.Response{}
		var err error
		resp, err = do(r, fullUrl)
		if err != nil {
			if !retried {
				if protocol == httpx.HTTPS {
					protocol = httpx.HTTP
				} else {
					protocol = httpx.HTTPS
				}
				retried = true
				goto retry
			}
		}
		builder := &strings.Builder{}
		builder.WriteString(fullUrl)

		var title string
		if resp != nil {
			title = resp.Title
		}

		p, err := url.Parse(fullUrl)
		var ip string
		var ipArray []string
		if err != nil {
			ip = ""
		} else {
			hostname := p.Hostname()
			ip = r.Ht.Dialer.GetDialedIP(hostname)
			// ip为空，看看p.host是不是ip
			if ip == "" {
				address := net.ParseIP(hostname)
				if address != nil {
					ip = address.String()
				}
			}
		}
		dnsData, err := r.Ht.Dialer.GetDNSData(p.Host)
		if dnsData != nil && err == nil {
			ipArray = append(ipArray, dnsData.CNAME...)
			ipArray = append(ipArray, dnsData.A...)
			ipArray = append(ipArray, dnsData.AAAA...)
		}
		cname := strings.Join(ipArray, ",")

		// CDN检测
		cdn, err := r.Ht.CDNCheck(resp, r.Ip, cname)
		if err != nil {
			log2.Warn("[CDNCheck]:", err)
			fmt.Fprintln(color.Output, color.HiYellowString("[WARN]"), "["+time.Now().Format("2006-01-02 15:04:05")+"]", "[CDNCheck]:", err)
		}

		// 指纹处理
		fofaResults, err := r.Fofa.Matcher(resp)
		if err != nil {
			log2.Warn("[FOFAFRINGER]:", err)
			fmt.Fprintln(color.Output, color.HiYellowString("[WARN]"), "["+time.Now().Format("2006-01-02 15:04:05")+"]", "[FOFAFRINGER]:", err)
		}
		var webbanner FingerResult
		if resp != nil {
			if resp.TLSData != nil {
				webbanner = FingerResult{
					Title:         title,
					TLSData:       resp.TLSData,
					ContentLength: resp.ContentLength,
					StatusCode:    resp.StatusCode,
					ResponseTime:  time.Since(timeStart).String(),
					Str:           builder.String(),
					Header:        resp.HeaderStr,
					FirstLine:     resp.FirstLine,
					Headers:       resp.Headers,
					DataStr:       resp.DataStr,
					Fingers:       fofaResults,
					CDN:           cdn,
				}
			} else {
				tlsdata := &httpx.TLSData{
					DNSNames:         nil,
					Emails:           nil,
					CommonName:       nil,
					Organization:     nil,
					IssuerCommonName: nil,
					IssuerOrg:        nil,
				}
				webbanner = FingerResult{
					Title:         title,
					TLSData:       tlsdata,
					ContentLength: resp.ContentLength,
					StatusCode:    resp.StatusCode,
					ResponseTime:  time.Since(timeStart).String(),
					Str:           builder.String(),
					Header:        resp.HeaderStr,
					FirstLine:     resp.FirstLine,
					Headers:       resp.Headers,
					DataStr:       resp.DataStr,
					Fingers:       fofaResults,
					CDN:           cdn,
				}
			}
		} else {
			webbanner = FingerResult{}
		}
		return webbanner
	}
	return FingerResult{}
}

func Close(r *Runner) {
	r.Ht.Dialer.Close()
	_ = r.Hm.Close()
}
