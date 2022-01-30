package servercheck

import (
	"fmt"
	"github.com/go-redis/redis"
	"time"

)

/*
创建人员：云深不知处
创建时间：2022/1/5
程序功能：redis检测单元
*/


func RedisCheck(ip string,port string,user string,pwd string) bool {
	client := redis.NewClient(&redis.Options{
		Addr:        fmt.Sprintf(`%s:%s`, ip, port),
		Password:    pwd,
		DB:          0,
		DialTimeout: 6 * time.Second,
	})

	defer client.Close()

	status, err := client.Ping().Result()
	if err != nil {
		return false
	}

	if status == "PONG" {
		return true
	}

	return false
}
