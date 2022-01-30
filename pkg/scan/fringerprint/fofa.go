package fringerprint

import (
	"BeeScan-scan/pkg/httpx"
	log2 "BeeScan-scan/pkg/log"
	"embed"
	"encoding/json"
	"fmt"
	"github.com/boy-hack/govaluate"
	"github.com/fatih/color"
	"strings"
	"time"
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

func FOFAInit(f embed.FS) FofaPrints {

	FofaJson, err := f.ReadFile("goby.json")
	if err != nil {
		log2.Error("[FOFAInit]:", err)
		fmt.Fprintln(color.Output, color.HiRedString("[ERRO]"), "["+time.Now().Format("2006-01-02 15:04:05")+"]", "[FOFAInit]:", err)
	}
	var fofas FofaPrints
	err1 := json.Unmarshal(FofaJson, &fofas)
	if err1 != nil {
		log2.Error("[FOFAInit]:", err1)
		fmt.Fprintln(color.Output, color.HiRedString("[ERRO]"), "["+time.Now().Format("2006-01-02 15:04:05")+"]", "[FOFAInit]:", err1)
	}
	return fofas
}

func (f *Fofa) Matcher(response *httpx.Response) (bool, error) {
	expString := f.Condition
	expression, err := govaluate.NewEvaluableExpressionWithFunctions(expString, HelperFunctions(response))
	if err != nil {
		return false, err
	}
	paramters := make(map[string]interface{})
	if response != nil {
		if response.Title != "" {
			paramters["title"] = response.Title
		} else {
			paramters["title"] = ""
		}
		if response.GetHeader("server") != "" {
			paramters["server"] = response.GetHeader("server")
		} else {
			paramters["server"] = ""
		}
	}
	paramters["protocol"] = "http"
	result, err := expression.Evaluate(paramters)
	if err != nil {
		return false, err
	}
	t := result.(bool)
	return t, err
}
func (f *FofaPrints) Matcher(response *httpx.Response) (FofaPrints, error) {
	var ret FofaPrints
	for _, item := range *f {
		v, err := item.Matcher(response)
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
func HelperFunctions(resp *httpx.Response) (functions map[string]govaluate.ExpressionFunction) {
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
		return false, nil
	}

	functions["banner_contains"] = func(args ...interface{}) (interface{}, error) {
		return false, nil
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
		return false, nil
	}

	functions["port_contains"] = func(args ...interface{}) (interface{}, error) {
		return false, nil
	}

	return functions
}

func toString(v interface{}) string {
	return fmt.Sprint(v)
}

func toInt(v interface{}) int {
	return int(v.(float64))
}
