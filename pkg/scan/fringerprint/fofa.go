package fringerprint

import (
	"BeeScan-scan/pkg/httpx"
	log2 "BeeScan-scan/pkg/log"
	"BeeScan-scan/pkg/scan/gonmap"
	"BeeScan-scan/pkg/util"
	"embed"
	"encoding/json"
	"fmt"
	"github.com/boy-hack/govaluate"
	"strings"
)

/*
创建人员：云深不知处
创建时间：2022/1/1
程序功能：指纹识别
*/

type Fofa struct {
	RuleId         string `json:"rule_id"`
	Level          string `json:"level"`
	SoftHard       string `json:"softhard"`
	Product        string `json:"product"`
	Company        string `json:"company"`
	Category       string `json:"category"`
	ParentCategory string `json:"parent_category"`
	Condition      string `json:"Condition"`
}
type FofaPrints []Fofa

var FofaJson []byte

func FOFAInit(f embed.FS) *FofaPrints {
	var err error
	FofaJson, err = f.ReadFile("goby.json")
	if err != nil {
		log2.Error("[FOFAInit]:", err)
	}
	fofas := &FofaPrints{}
	err1 := json.Unmarshal(FofaJson, fofas)
	if err1 != nil {
		log2.Error("[FOFAInit]:", err1)
	}
	return fofas
}

func (f *Fofa) Matcher(response *httpx.Response, gomapres *gonmap.Result, port string) (bool, error) {
	expString := f.Condition
	expression, err := govaluate.NewEvaluableExpressionWithFunctions(expString, HelperFunctions(response, gomapres, port))
	if err != nil {
		return false, err
	}
	paramters := make(map[string]interface{})
	if response != nil && gomapres != nil {
		if response.Title != "" {
			paramters["title"] = response.Title
		} else {
			paramters["title"] = ""
		}

		if response.GetHeader("server") != "" {
			paramters["server"] = response.GetHeader("server")
		} else if response.GetHeader("server") == "" && gomapres.Name != "" {
			paramters["server"] = gomapres.Service.Name
		} else {
			paramters["server"] = ""
		}

		if gomapres.Service.Protocol != "" {
			paramters["protocol"] = gomapres.Service.Protocol
		} else {
			paramters["protocol"] = "http"
		}

		if response.HeaderStr != "" {
			paramters["header"] = response.HeaderStr
		} else {
			paramters["header"] = ""
		}

		if response.DataStr != "" {
			paramters["body"] = response.DataStr
		} else {
			paramters["body"] = ""
		}

		if gomapres.Banner != "" {
			paramters["banner"] = gomapres.Banner
		} else {
			paramters["banner"] = ""
		}

		if response.TLSData != nil {
			var cert string
			cert += util.StrToSlince(response.TLSData.DNSNames) + "\n"
			cert += util.StrToSlince(response.TLSData.IssuerCommonName) + "\n"
			cert += util.StrToSlince(response.TLSData.Organization) + "\n"
			cert += util.StrToSlince(response.TLSData.CommonName) + "\n"
			cert += util.StrToSlince(response.TLSData.Emails) + "\n"
			cert += util.StrToSlince(response.TLSData.IssuerOrg)
			paramters["cert"] = cert
		} else {
			paramters["cert"] = ""
		}

		if port != "" {
			paramters["port"] = port
		} else {
			paramters["port"] = ""
		}

	}

	result, err := expression.Evaluate(paramters)
	if err != nil {
		return false, err
	}
	t := result.(bool)
	return t, err
}
func (f *FofaPrints) Matcher(response *httpx.Response, gomapres *gonmap.Result, port string) (FofaPrints, error) {
	var ret FofaPrints
	for _, item := range *f {
		v, err := item.Matcher(response, gomapres, port)
		if err != nil {
			return nil, err
		}
		if v {
			ret = append(ret, item)
		}
	}
	return ret, nil
}

// HelperFunctions contains the dsl functions
func HelperFunctions(resp *httpx.Response, gomapres *gonmap.Result, port string) (functions map[string]govaluate.ExpressionFunction) {
	functions = make(map[string]govaluate.ExpressionFunction)
	functions["title_contains"] = func(args ...interface{}) (interface{}, error) {
		pattern := strings.ToLower(toString(args[0]))
		var title string
		if resp != nil {
			if resp.Title != "" {
				title = strings.ToLower(resp.Title)
			} else {
				title = ""
			}
		}
		return strings.Index(title, pattern) != -1, nil
	}

	functions["body_contains"] = func(args ...interface{}) (interface{}, error) {
		pattern := strings.ToLower(toString(args[0]))
		var data string
		if resp != nil {
			if resp.DataStr != "" {
				data = strings.ToLower(resp.DataStr)
			} else {
				data = ""
			}
		}
		return strings.Index(data, pattern) != -1, nil
	}

	functions["protocol_contains"] = func(args ...interface{}) (interface{}, error) {
		pattern := strings.ToLower(toString(args[0]))
		var protocol string
		if gomapres != nil {
			if gomapres.Service.Protocol != "" {
				protocol = strings.ToLower(gomapres.Service.Protocol)
			}
			if resp != nil {
				if gomapres.Service.Protocol == "" && resp.HeaderStr != "" {
					protocol = "http"
				}
			}
		}

		return strings.Index(protocol, pattern) != -1, nil
	}

	functions["banner_contains"] = func(args ...interface{}) (interface{}, error) {
		pattern := strings.ToLower(toString(args[0]))
		var banner string
		if gomapres != nil {

			if gomapres.Banner != "" {
				banner = strings.ToLower(gomapres.Banner)
			}
		}
		return strings.Index(banner, pattern) != -1, nil
	}

	functions["header_contains"] = func(args ...interface{}) (interface{}, error) {
		pattern := strings.ToLower(toString(args[0]))
		var data string
		if resp != nil {
			if resp.HeaderStr != "" {
				data = strings.ToLower(resp.HeaderStr)
			} else {
				data = ""
			}
		}
		return strings.Index(data, pattern) != -1, nil
	}

	functions["server_contains"] = func(args ...interface{}) (interface{}, error) {
		pattern := strings.ToLower(toString(args[0]))
		var server string
		if resp != nil {
			server = resp.GetHeader("server")
		}
		return strings.Index(server, pattern) != -1, nil
	}

	functions["cert_contains"] = func(args ...interface{}) (interface{}, error) {
		pattern := strings.ToLower(toString(args[0]))
		var cert string
		if resp != nil {
			if resp.TLSData != nil {
				cert += util.StrToSlince(resp.TLSData.DNSNames) + "\n"
				cert += util.StrToSlince(resp.TLSData.IssuerCommonName) + "\n"
				cert += util.StrToSlince(resp.TLSData.Organization) + "\n"
				cert += util.StrToSlince(resp.TLSData.CommonName) + "\n"
				cert += util.StrToSlince(resp.TLSData.Emails) + "\n"
				cert += util.StrToSlince(resp.TLSData.IssuerOrg) + "\n"
			}
		}
		return strings.Index(cert, pattern) != -1, nil
	}

	functions["port_contains"] = func(args ...interface{}) (interface{}, error) {
		pattern := strings.ToLower(toString(args[0]))
		if port != "" {
			return strings.Index(port, pattern) != -1, nil
		}
		return false, nil
	}

	return functions
}

func toString(v interface{}) string {
	return fmt.Sprint(v)
}
