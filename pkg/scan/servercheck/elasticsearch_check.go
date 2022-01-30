package servercheck

import (
	"fmt"
	"github.com/olivere/elastic"
	"context"
)

/*
创建人员：云深不知处
创建时间：2022/1/5
程序功能：elasticsearch检测单元
*/


func ElasticSearchCheck(ip string,port string,user string,pwd string) bool {
	flag := false

	client, err := elastic.NewClient(elastic.SetURL(fmt.Sprintf("http://%s:%s", ip, port)),
		elastic.SetBasicAuth(user, pwd),
	)
	if err == nil {
		ctx := context.Background()
		_, _, err = client.Ping(fmt.Sprintf("http://%s:%s", ip, port)).Do(ctx)
		if err == nil {
			flag = true
		}
	}
	return flag
}