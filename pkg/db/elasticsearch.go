package db

import (
	"BeeScan-scan/pkg/config"
	log2 "BeeScan-scan/pkg/log"
	"BeeScan-scan/pkg/result"
	"BeeScan-scan/pkg/util"
	"context"
	"encoding/json"
	"fmt"
	"github.com/fatih/color"
	"github.com/olivere/elastic/v7"
	"io/ioutil"
	"os"
	"time"
)

/*
创建人员：云深不知处
创建时间：2022/1/14
程序功能：elasticsearch数据库操作
*/

type NodeLog struct {
	Log      string `json:"log"`
	LastTime string `json:"lastTime"`
}

// ElasticSearchInit es数据库初始化连接
func ElasticSearchInit() *elastic.Client {
	host := "http://" + config.GlobalConfig.DBConfig.Elasticsearch.Host + ":" + config.GlobalConfig.DBConfig.Elasticsearch.Port
	client, err := elastic.NewClient(
		elastic.SetURL(host),
		elastic.SetBasicAuth(config.GlobalConfig.DBConfig.Elasticsearch.Username, config.GlobalConfig.DBConfig.Elasticsearch.Password),
		elastic.SetSniff(false),
	)
	if err != nil {
		log2.Error("[ElasticSearchInit]:", err)
		fmt.Fprintln(color.Output, color.HiRedString("[ERRO]"), "["+time.Now().Format("2006-01-02 15:04:05")+"]", "[ElasticSearchInit]:", err)
		os.Exit(1)
	}
	return client
}

// EsAdd 添加结果到es数据库
func EsAdd(client *elastic.Client, res *result.Output) {

	// 文档件存在则更新，否则插入
	_, err := client.Update().Index(config.GlobalConfig.DBConfig.Elasticsearch.Index).Id(res.Ip + "-" + res.Port + "-" + res.Domain).Doc(res).Upsert(res).Refresh("true").Do(context.Background())
	if err != nil {
		log2.Error(err)
		fmt.Fprintln(color.Output, color.HiRedString("[ERRO]"), "[DBEsUpInsert]:", err)
	}

}

// ESLogAdd 日志写入es数据库
func ESLogAdd(client *elastic.Client, filename string) {
	var TheNodeLog NodeLog
	var logs []byte
	var err error
	if util.FileExist(filename) {
		logs, err = ioutil.ReadFile(filename)
		if err != nil {
			log2.Error(err)
			fmt.Fprintln(color.Output, color.HiRedString("[ERRO]"), "[ESLogAdd]:", err)
		}
		TheNodeLog.Log = string(logs)
		TheNodeLog.LastTime = time.Now().Format("2006-01-02 15:04:05")
		_, err = client.Update().Index(config.GlobalConfig.DBConfig.Elasticsearch.Index).Id(config.GlobalConfig.NodeConfig.NodeName + "_log").Doc(TheNodeLog).Upsert(TheNodeLog).Refresh("true").Do(context.Background())
		if err != nil {
			log2.Error(err)
			fmt.Fprintln(color.Output, color.HiRedString("[ERRO]"), "[ESLogAdd]:", err)
		}
	}
}

// EsScanRegular 定期重新扫描es数据库中时间超过30天的目标
func EsScanRegular(client *elastic.Client) []string {
	var res *elastic.SearchResult
	var count int64
	var err error
	var targets []string
	count, err = client.Count().Index(config.GlobalConfig.DBConfig.Elasticsearch.Index).Do(context.Background())
	res, err = client.Search(config.GlobalConfig.DBConfig.Elasticsearch.Index).Size(int(count)).From(0).Do(context.Background())
	if err != nil {
		log2.Error(err)
		fmt.Fprintln(color.Output, color.HiRedString("[ERRO]"), "[EsScanRegular]:", err)
	}
	if res != nil {
		if res.Hits != nil {
			if res.Hits.Hits != nil {
				for _, item := range res.Hits.Hits {
					out := result.Output{}
					err = json.Unmarshal(item.Source, &out)
					if err != nil {
						fmt.Fprintln(color.Output, color.HiRedString("[ERRO]"), "[EsScanRegular]:", err)
					}
					if util.DaySub(out.LastTime) >= 30 {
						if out.Domain != "" && out.Protocol == "UDP" {
							target := out.Domain + ":" + "U:" + out.Port
							targets = append(targets, target)
						} else if out.Domain != "" && out.Protocol == "TCP" {
							target := out.Domain + ":" + out.Port
							targets = append(targets, target)
						} else if out.Ip != "" && out.Protocol == "UDP" {
							target := out.Ip + ":" + "U:" + out.Port
							targets = append(targets, target)
						} else if out.Ip != "" && out.Protocol == "TCP" {
							target := out.TargetId + ":" + out.Port
							targets = append(targets, target)
						}
					}
				}
			}
		}
	}
	return targets
}

func QueryLogByID(client *elastic.Client, nodename string) string {
	var res *elastic.GetResult
	var err error
	var TheNodeLog NodeLog
	res, err = client.Get().Index(config.GlobalConfig.DBConfig.Elasticsearch.Index).Id(nodename + "_log").Do(context.Background())
	if err != nil {
		fmt.Fprintln(color.Output, color.HiRedString("[ERRO]"), "[ESQueryLog]:", err)
	}

	if res != nil {
		if res.Found {
			if res.Source != nil {
				err = json.Unmarshal(res.Source, &TheNodeLog)
				if err != nil {
					fmt.Fprintln(color.Output, color.HiRedString("[ERRO]"), "[ESQueryLog]:", err)
				}
			}
		}
	}
	return TheNodeLog.LastTime
}
