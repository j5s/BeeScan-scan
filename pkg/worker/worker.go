package worker

import (
	"BeeScan-scan/internal/runner"
	"BeeScan-scan/pkg/config"
	"BeeScan-scan/pkg/job"
	log2 "BeeScan-scan/pkg/log"
	"BeeScan-scan/pkg/node"
	"BeeScan-scan/pkg/result"
	"BeeScan-scan/pkg/scan/gonmap"
	"BeeScan-scan/pkg/scan/ipinfo"
	"fmt"
	"github.com/fatih/color"
	redis2 "github.com/go-redis/redis"
	"github.com/panjf2000/ants/v2"
	"go.uber.org/ratelimit"
	"sync"
	"time"
)

/*
创建人员：云深不知处
创建时间：2022/1/4
程序功能：任务池
*/

func WorkerInit(nodestate *job.NodeState, taskstate *job.TaskState, wg *sync.WaitGroup, conn *redis2.Client, GoNmap *gonmap.VScan, region *ipinfo.Ip2Region, tmpresults chan *result.Output) *ants.PoolWithFunc {
	rl := ratelimit.New(config.GlobalConfig.WorkerConfig.Thread)
	p, _ := ants.NewPoolWithFunc(config.GlobalConfig.WorkerConfig.WorkerNumber, func(j interface{}) {
		if j.(*runner.Runner) != nil {
			if j.(*runner.Runner).Ip != "" || j.(*runner.Runner).Domain != "" {
				nodestate.Running++
				taskstate.Running++
				rl.Take()
				if j.(*runner.Runner).Ip != "" {
					log2.Info("[Scanning]:", j.(*runner.Runner).Ip)
					fmt.Fprintln(color.Output, color.HiCyanString("[INFO]"), "["+time.Now().Format("2006-01-02 15:04:05")+"]", "[Scanning]:", j.(*runner.Runner).Ip)
					log2.Info("[Tasks]:", nodestate.Tasks, "[Running]:", nodestate.Running, "[Finished]:", nodestate.Finished)
					fmt.Fprintln(color.Output, color.HiCyanString("[INFO]"), "["+time.Now().Format("2006-01-02 15:04:05")+"]", "[Tasks]:", nodestate.Tasks, "[Running]:", nodestate.Running, "[Finished]:", nodestate.Finished)
				} else if j.(*runner.Runner).Domain != "" {
					log2.Info("[Scanning]:", j.(*runner.Runner).Domain)
					fmt.Fprintln(color.Output, color.HiCyanString("[INFO]"), "["+time.Now().Format("2006-01-02 15:04:05")+"]", "[Scanning]:", j.(*runner.Runner).Domain)
					log2.Info("[Tasks]:", nodestate.Tasks, "[Running]:", nodestate.Running, "[Finished]:", nodestate.Finished)
					fmt.Fprintln(color.Output, color.HiCyanString("[INFO]"), "["+time.Now().Format("2006-01-02 15:04:05")+"]", "[Tasks]:", nodestate.Tasks, "[Running]:", nodestate.Running, "[Finished]:", nodestate.Finished)
				}
				node.NodeUpdate(conn, config.GlobalConfig.NodeConfig.NodeName, nodestate)
				node.TaskUpdate(conn, taskstate)
				result := runner.Scan(j.(*runner.Runner), GoNmap, region) // 执行扫描
				nodestate.Running--
				taskstate.Running--
				nodestate.Finished++
				taskstate.Finished++
				if j.(*runner.Runner).Ip != "" {
					log2.Info("[Scanned]:", j.(*runner.Runner).Ip)
					fmt.Fprintln(color.Output, color.HiCyanString("[INFO]"), "["+time.Now().Format("2006-01-02 15:04:05")+"]", "[Scanned]:", j.(*runner.Runner).Ip)
					fmt.Fprintln(color.Output, color.HiCyanString("[INFO]"), "["+time.Now().Format("2006-01-02 15:04:05")+"]", "[Tasks]:", nodestate.Tasks, "[Running]:", nodestate.Running, "[Finished]:", nodestate.Finished)
				} else if j.(*runner.Runner).Domain != "" {
					log2.Info("[Scanning]:", j.(*runner.Runner).Domain)
					fmt.Fprintln(color.Output, color.HiCyanString("[INFO]"), "["+time.Now().Format("2006-01-02 15:04:05")+"]", "[Scanning]:", j.(*runner.Runner).Domain)
					log2.Info("[Tasks]:", nodestate.Tasks, "[Running]:", nodestate.Running, "[Finished]:", nodestate.Finished)
					fmt.Fprintln(color.Output, color.HiCyanString("[INFO]"), "["+time.Now().Format("2006-01-02 15:04:05")+"]", "[Tasks]:", nodestate.Tasks, "[Running]:", nodestate.Running, "[Finished]:", nodestate.Finished)
				}
				node.NodeUpdate(conn, config.GlobalConfig.NodeConfig.NodeName, nodestate)
				node.TaskUpdate(conn, taskstate)
				tmpresults <- result
				defer wg.Done()
			}
		}
	})
	return p
}
