package db

import (
	"BeeScan-scan/pkg/config"
	"BeeScan-scan/pkg/file"
	log2 "BeeScan-scan/pkg/log"
	"BeeScan-scan/pkg/runner"
	"context"
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
	Log string `json:"log"`
}

// ElasticSearchInit es数据库初始化连接
func ElasticSearchInit() *elastic.Client {
	host := "http://" + config.GlobalConfig.DBConfig.Elasticsearch.Host + ":" + config.GlobalConfig.DBConfig.Elasticsearch.Port
	client, err := elastic.NewClient(
		elastic.SetURL(host),
		elastic.SetBasicAuth(config.GlobalConfig.DBConfig.Elasticsearch.Username, config.GlobalConfig.DBConfig.Elasticsearch.Password),
	)
	if err != nil {
		log2.Error("[ElasticSearchInit]:", err)
		fmt.Fprintln(color.Output, color.HiRedString("[ERRO]"), "["+time.Now().Format("2006-01-02 15:04:05")+"]", "[ElasticSearchInit]:", err)
		os.Exit(1)
	}
	return client
}

// EsAdd 添加结果到es数据库
func EsAdd(client *elastic.Client, res *runner.Output) {

	// 文档件存在则更新，否则插入
	_, err := client.Update().Index(config.GlobalConfig.DBConfig.Elasticsearch.Index).Id(res.Ip + "-" + res.Port + "-" + res.Domain).Doc(res).Upsert(res).Refresh("true").Do(context.Background())
	if err != nil {
		log2.Error(err)
		fmt.Fprintln(color.Output, color.HiRedString("[ERRO]"), "[DBEsUpInsert]:", err)
	}

}

func ESLogAdd(client *elastic.Client, filename string) {
	var TheNodeLog NodeLog
	var logs []byte
	var err error
	if file.Exists(filename) {
		logs, err = ioutil.ReadFile(filename)
		if err != nil {
			log2.Error(err)
			fmt.Fprintln(color.Output, color.HiRedString("[ERRO]"), "[ESLogAdd]:", err)
		}
		TheNodeLog.Log = string(logs)
		_, err = client.Update().Index(config.GlobalConfig.DBConfig.Elasticsearch.Index).Id(config.GlobalConfig.NodeConfig.NodeName + "_log").Doc(TheNodeLog).Upsert(TheNodeLog).Refresh("true").Do(context.Background())
		if err != nil {
			log2.Error(err)
			fmt.Fprintln(color.Output, color.HiRedString("[ERRO]"), "[ESLogAdd]:", err)
		}
	}
}
