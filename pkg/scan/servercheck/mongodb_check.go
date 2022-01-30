package servercheck

import (
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"context"
)

/*
创建人员：云深不知处
创建时间：2022/1/5
程序功能：mongodb检测单元
*/


func MongoDBCheck(ip string,port string,user string,pwd string) bool {
	// 设置客户端连接配置
	clientOptions := options.Client().ApplyURI(
		fmt.Sprintf(`mongodb://%s:%s@%s:%s`, user, pwd, ip, port),
	)

	// 连接到MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		return false
	}

	// 检查连接
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		return false
	}

	return true
}