package db

import (
	"BeeScan-scan/pkg/config"
	log2 "BeeScan-scan/pkg/log"
	"fmt"
	"github.com/fatih/color"
	"github.com/go-redis/redis"
	"os"
	"time"
)

/*
创建人员：云深不知处
创建时间：2022/1/4
程序功能：redis连接与操作
*/

// InitRedis 初始化连接
func InitRedis() *redis.Client {
	conn := redis.NewClient(&redis.Options{
		Addr:     config.GlobalConfig.DBConfig.Redis.Host + ":" + config.GlobalConfig.DBConfig.Redis.Port,
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	Pong, err := conn.Ping().Result()
	if err != nil {
		log2.Error("[RedisInit]:", err)
		fmt.Fprintln(color.Output, color.HiRedString("[ERROR]"), "["+time.Now().Format("2006-01-02 15:04:05")+"]", "[RedisInit]:", err)
		os.Exit(1)
	} else if Pong == "PONG" {
		return conn
	}
	return conn
}

// RecvJob 从消息队列接收任务
func RecvJob(c *redis.Client) []string {
	val := c.BLPop(3*time.Second, config.GlobalConfig.NodeConfig.NodeQueue)
	return val.Val()
}
