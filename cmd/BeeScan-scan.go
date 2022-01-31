package main

import (
	"BeeScan-scan/internal/runner"
	"BeeScan-scan/pkg/banner"
	"BeeScan-scan/pkg/config"
	"BeeScan-scan/pkg/db"
	"BeeScan-scan/pkg/job"
	log2 "BeeScan-scan/pkg/log"
	"BeeScan-scan/pkg/node"
	"BeeScan-scan/pkg/result"
	"BeeScan-scan/pkg/scan/fringerprint"
	"BeeScan-scan/pkg/scan/gonmap"
	"BeeScan-scan/pkg/scan/gowapp"
	"BeeScan-scan/pkg/scan/ipinfo"
	"BeeScan-scan/pkg/util"
	"BeeScan-scan/pkg/worker"
	"embed"
	"fmt"
	"github.com/fatih/color"
	redis2 "github.com/go-redis/redis"
	gowap "github.com/jiaocoll/GoWapp/pkg/core"
	"github.com/olivere/elastic/v7"
	"github.com/panjf2000/ants/v2"
	"go.uber.org/ratelimit"
	"sync"
	"time"
)

/*
创建人员：云深不知处
创建时间：2022/1/3
程序功能：主函数
*/

//go:embed nmap-probes wapp.json goby.json ip2region.db
var f embed.FS

var (
	queue      *job.Queue
	jobs       chan *runner.Runner
	results    chan *result.Output
	tmpresults chan *result.Output
	conn       *redis2.Client
	nodestate  *job.NodeState
	taskstate  *job.TaskState
	region     *ipinfo.Ip2Region
	wg         sync.WaitGroup
	rl         ratelimit.Limiter
	wapp       *gowap.Wappalyzer
	esclient   *elastic.Client
	fofaPrints *fringerprint.FofaPrints
	GoNmap     *gonmap.VScan
	p          *ants.PoolWithFunc
)

func init() {
	banner.Banner()
	fmt.Fprintln(color.Output, color.HiMagentaString("Initializing......"))
	config.Setup()
	log2.Setup()
	jobs = make(chan *runner.Runner, 1000)
	tmpresults = make(chan *result.Output, 1000)
	results = make(chan *result.Output, 1000)
	GoNmap = gonmap.GoNmapInit(f)
	wapp, _ = gowapp.GowappInit(f)
	fofaPrints = fringerprint.FOFAInit(f)
	region = ipinfo.IpInfoInit(f)
	conn = db.InitRedis()
	esclient = db.ElasticSearchInit()                                //初始化redis连接
	node.NodeRegister(conn, config.GlobalConfig.NodeConfig.NodeName) //节点注册
	queue = job.NewQueue()
	nodestate = &job.NodeState{
		Tasks:     0,
		Running:   0,
		Finished:  0,
		State:     "Free",
		StartTime: time.Now().Format("2006-01-02 15:04:05"),
	}
	taskstate = &job.TaskState{
		Name:      "",
		TargetNum: 0,
		Tasks:     0,
		Running:   0,
		Finished:  0,
		LastTime:  time.Time{},
	}
	//wg = sizedwaitgroup.New(config.GlobalConfig.WorkerConfig.WorkerNumber)
	//rl = ratelimit.New(config.GlobalConfig.WorkerConfig.Thread)
	//p, _ = ants.NewPoolWithFunc(config.GlobalConfig.WorkerConfig.WorkerNumber, func(j interface{}) {
	//	if j.(*runner.Runner) != nil {
	//		if j.(*runner.Runner).Ip != "" || j.(*runner.Runner).Domain != "" {
	//			rl.Take()
	//			nodestate.Running++
	//			taskstate.Running++
	//			if j.(*runner.Runner).Ip != "" {
	//				log2.Info("[Scanning]:", j.(*runner.Runner).Ip)
	//				fmt.Fprintln(color.Output, color.HiCyanString("[INFO]"), "["+time.Now().Format("2006-01-02 15:04:05")+"]", "[Scanning]:", j.(*runner.Runner).Ip)
	//				log2.Info("[Tasks]:", nodestate.Tasks, "[Running]:", nodestate.Running, "[Finished]:", nodestate.Finished)
	//				fmt.Fprintln(color.Output, color.HiCyanString("[INFO]"), "["+time.Now().Format("2006-01-02 15:04:05")+"]", "[Tasks]:", nodestate.Tasks, "[Running]:", nodestate.Running, "[Finished]:", nodestate.Finished)
	//			} else if j.(*runner.Runner).Domain != "" {
	//				log2.Info("[Scanning]:", j.(*runner.Runner).Domain)
	//				fmt.Fprintln(color.Output, color.HiCyanString("[INFO]"), "["+time.Now().Format("2006-01-02 15:04:05")+"]", "[Scanning]:", j.(*runner.Runner).Domain)
	//				log2.Info("[Tasks]:", nodestate.Tasks, "[Running]:", nodestate.Running, "[Finished]:", nodestate.Finished)
	//				fmt.Fprintln(color.Output, color.HiCyanString("[INFO]"), "["+time.Now().Format("2006-01-02 15:04:05")+"]", "[Tasks]:", nodestate.Tasks, "[Running]:", nodestate.Running, "[Finished]:", nodestate.Finished)
	//			}
	//			node.NodeUpdate(conn, config.GlobalConfig.NodeConfig.NodeName, nodestate)
	//			node.TaskUpdate(conn, taskstate)
	//			result := Scan(j.(*runner.Runner)) // 执行扫描
	//			nodestate.Running--
	//			taskstate.Running--
	//			nodestate.Finished++
	//			taskstate.Finished++
	//			if j.(*runner.Runner).Ip != "" {
	//				log2.Info("[Scanned]:", j.(*runner.Runner).Ip)
	//				fmt.Fprintln(color.Output, color.HiCyanString("[INFO]"), "["+time.Now().Format("2006-01-02 15:04:05")+"]", "[Scanned]:", j.(*runner.Runner).Ip)
	//				fmt.Fprintln(color.Output, color.HiCyanString("[INFO]"), "["+time.Now().Format("2006-01-02 15:04:05")+"]", "[Tasks]:", nodestate.Tasks, "[Running]:", nodestate.Running, "[Finished]:", nodestate.Finished)
	//			} else if j.(*runner.Runner).Domain != "" {
	//				log2.Info("[Scanning]:", j.(*runner.Runner).Domain)
	//				fmt.Fprintln(color.Output, color.HiCyanString("[INFO]"), "["+time.Now().Format("2006-01-02 15:04:05")+"]", "[Scanning]:", j.(*runner.Runner).Domain)
	//				log2.Info("[Tasks]:", nodestate.Tasks, "[Running]:", nodestate.Running, "[Finished]:", nodestate.Finished)
	//				fmt.Fprintln(color.Output, color.HiCyanString("[INFO]"), "["+time.Now().Format("2006-01-02 15:04:05")+"]", "[Tasks]:", nodestate.Tasks, "[Running]:", nodestate.Running, "[Finished]:", nodestate.Finished)
	//			}
	//			node.NodeUpdate(conn, config.GlobalConfig.NodeConfig.NodeName, nodestate)
	//			node.TaskUpdate(conn, taskstate)
	//			tmpresults <- result
	//			defer wg.Done()
	//		}
	//	}
	//})

	p = worker.WorkerInit(nodestate, taskstate, &wg, conn, GoNmap, region, tmpresults)
	fmt.Fprintln(color.Output, color.HiMagentaString("Initialized!"))
}

func main() {

	for true {
		// 扫描节点状态更新
		node.NodeUpdate(conn, config.GlobalConfig.NodeConfig.NodeName, nodestate)

		// 处理消息队列任务
		runner.Handlejob(conn, queue, taskstate)
		var runners []*runner.Runner
		if queue.Length > 0 { //判断队列长度
			// 任务实例集合
			runners = runner.HandleTargets(queue, fofaPrints)
		}

		if len(runners) > 0 {
			// 添加任务到任务队列
			for _, v := range runners {
				jobs <- v
				nodestate.Tasks++
			}
			runners = nil
		}

		if len(jobs) > 0 {
			nodestate.State = "Running"
			for i := 0; i < len(jobs); i++ {
				j := <-jobs
				wg.Add(1)
				_ = p.Invoke(j)
			}
			wg.Wait()
		}

		// Gowapp扫描模块，单独进行扫描，防止rod模块多线程出现panic
		if len(tmpresults) > 0 {
			for i := 0; i < len(tmpresults); i++ {
				tmpres := <-tmpresults
				if tmpres != nil {
					if tmpres.Banner != "" {
						tmpres.Wappalyzer = gowapp.GoWapp(tmpres, wapp)
						time.Sleep(500 * time.Millisecond)
						results <- tmpres
					}
				}
			}
		}

		if len(results) > 0 {
			// 遍历结果，写入数据库
			for i := 0; i < len(results); i++ {
				res := <-results
				if res != nil {
					db.EsAdd(esclient, res)
				}
			}
		}

		// 写入日志到es数据库
		if nodestate.State == "Free" && util.MinSub(db.QueryLogByID(esclient, config.GlobalConfig.NodeConfig.NodeName))%30 == 0 && util.MinSub(db.QueryLogByID(esclient, config.GlobalConfig.NodeConfig.NodeName)) > 0 {
			db.ESLogAdd(esclient, "BeeScanLogs.log")
		}

		// 扫描节点状态更新
		node.NodeUpdate(conn, config.GlobalConfig.NodeConfig.NodeName, nodestate)
		if len(jobs) == 0 {
			fmt.Fprintln(color.Output, color.HiCyanString("[INFO]"), "["+time.Now().Format("2006-01-02 15:04:05")+"]", "[ConnectCheck]")
			nodestate.State = "Free"
			taskstate.Name = ""
			taskstate.Running = 0
			taskstate.Finished = 0
			// 节点每运行15天会定时重新扫描es数据库中一定时长的目标
			if util.DaySub(nodestate.StartTime) > 15 && util.DaySub(nodestate.StartTime)%15 == 0 {
				RegularTargets := db.EsScanRegular(esclient)
				RegularTargets = util.Removesamesip(RegularTargets)
				if RegularTargets != nil {
					log2.Info("[targets]:", RegularTargets)
					fmt.Fprintln(color.Output, color.HiCyanString("[INFO]"), "["+time.Now().Format("2006-01-02 15:04:05")+"]", "[targets]:", RegularTargets)
					for _, v := range RegularTargets {
						job.Push(queue, v)
					}
				}
			}
			time.Sleep(10 * time.Second)
		}
	}
}

//// Handlejob 任务处理
//func Handlejob(c *redis2.Client, queue *job.Queue) {
//	var targets []string
//	// 查看消息队列，取出任务
//	lenval := c.LLen(config.GlobalConfig.NodeConfig.NodeQueue)
//	qlen := lenval.Val()
//	if qlen > 0 { // 若队列不空
//		for i := 1; i <= int(qlen); i++ {
//			tmpjob := db.RecvJob(c)
//			st := strings.Replace(tmpjob[1], "\"", "", -1)
//			tmptargets := strings.Split(st, ",")
//			taskstate.Tasks = len(tmptargets) - 1
//			for k, v := range tmptargets {
//				if k == 0 {
//					taskstate.Name = v
//				}
//				if k != 0 && v != "" {
//					targets = append(targets, v)
//				}
//			}
//			log2.Info("[targets]:", targets)
//			fmt.Fprintln(color.Output, color.HiCyanString("[INFO]"), "["+time.Now().Format("2006-01-02 15:04:05")+"]", "[targets]:", targets)
//		}
//		for _, t := range targets {
//			job.Push(queue, t) //将任务目标加入到任务队列中
//		}
//	}
//}
//
//// HandleTargets 生成扫描实例
//func HandleTargets(queue *job.Queue) []*runner.Runner {
//	var targets []string
//	var runners []*runner.Runner
//	for i := 0; i <= queue.Length; i++ {
//		targets = append(targets, job.Pop(queue))
//	}
//	if len(targets) > 0 {
//		for _, v := range targets {
//			target := strings.Split(v, ":")
//			if len(target) > 0 {
//				tmptarget := util.TargetsHandle(target[0]) //目标处理，若是c段地址，则返回一个ip段，若是单个ip，则直接返回单个ip切片，若是域名或url地址，则返回域名
//				for _, t := range tmptarget {
//					var runner2 *runner.Runner
//					var err1 error
//					if strings.Contains(t, "com") || strings.Contains(t, "cn") {
//						ip := getipbydomain.GetIPbyDomain(t)
//						if strings.Contains(target[1], "U:") {
//							tmp := strings.Split(target[1], ":")
//							port := tmp[1]
//							runner2, err1 = runner.NewRunner(ip, port, t, "udp", fofaPrints)
//						}
//						runner2, err1 = runner.NewRunner(ip, target[1], t, "tcp", fofaPrints)
//					} else {
//						if strings.Contains(target[1], "U") {
//							tmp := strings.Split(target[1], ":")
//							port := tmp[1]
//							runner2, err1 = runner.NewRunner(t, port, "", "udp", fofaPrints)
//						}
//						runner2, err1 = runner.NewRunner(t, target[1], "", "tcp", fofaPrints)
//					}
//					if err1 != nil {
//						log2.Error("[HandleTargets]:", err1)
//						fmt.Fprintln(color.Output, color.HiRedString("[ERROR]"), "["+time.Now().Format("2006-01-02 15:04:05")+"]", "[HandleTargets]:", err1)
//					}
//					runners = append(runners, runner2)
//				}
//			}
//		}
//	}
//	if len(runners) > 0 {
//		return runners
//	}
//	return nil
//}
//
//// Scan 扫描函数
//func Scan(target *runner.Runner) *runner.Output {
//	result := &runner.Output{}
//	// 域名存在与否
//	if target.Domain != "" {
//
//		// 主机存活探测
//		if icmp.IcmpCheckAlive(target.Domain, target.Ip) || ping.PingCheckAlive(target.Domain) || httpcheck.HttpCheck(target.Domain, target.Port, target.Ip) || tcp.TcpCheckAlive(target.Ip, target.Port) {
//
//			if tcp.TcpCheckAlive(target.Ip, target.Port) {
//				// 普通端口探测
//				nmapbanner, err := gonmap.GoNmapScan(GoNmap, target.Ip, target.Port, target.Protocol)
//
//				result.Servers = nmapbanner
//				if strings.Contains(result.Servers.Banner, "HTTP") {
//					result.Servers.Name = "http"
//					result.Servername = "http"
//				} else {
//					result.Servername = nmapbanner.Name
//				}
//				// web端口探测
//				webresult := runner.FingerResult{}
//				if result.Servername == "http" {
//					webresult = runner.Request(target)
//				}
//				result.Webbanner = webresult
//				result.Ip = target.Ip
//				result.Port = target.Port
//				result.Protocol = strings.ToUpper(target.Protocol)
//				result.Domain = target.Domain
//
//				if webresult.Header != "" {
//					result.Banner = result.Webbanner.Header
//				} else {
//					result.Banner = nmapbanner.Banner
//				}
//				// ip信息查询
//				info, err := ipinfo.GetIpinfo(region, target.Ip)
//				if err != nil {
//					log2.Warn("[GetIPInfo]:", err)
//					fmt.Fprintln(color.Output, color.HiYellowString("[WARN]"), "["+time.Now().Format("2006-01-02 15:04:05")+"]", "[GetIPInfo]:", err)
//				}
//				result.City = info.City
//				result.Region = info.Region
//				result.ISP = info.ISP
//				result.CityId = info.CityId
//				result.Province = info.Province
//				result.Country = info.Country
//				result.TargetId = target.Ip + "-" + target.Port + "-" + target.Domain
//				if result.Port == "80" {
//					result.Target = "http://" + target.Domain
//				} else {
//					result.Target = "http://" + target.Domain + ":" + result.Port
//				}
//				result.LastTime = time.Now().Format("2006-01-02 15:04:05")
//				return result
//			}
//		}
//	} else {
//		if cdncheck.IPCDNCheck(target.Ip) != true { //判断IP是否存在CDN
//
//			// 主机存活探测
//			if icmp.IcmpCheckAlive("", target.Ip) || ping.PingCheckAlive(target.Domain) || httpcheck.HttpCheck(target.Ip, target.Port, target.Ip) || tcp.TcpCheckAlive(target.Ip, target.Port) {
//
//				if tcp.TcpCheckAlive(target.Ip, target.Port) {
//					// 普通端口探测
//					nmapbanner, err := gonmap.GoNmapScan(GoNmap, target.Ip, target.Port, target.Protocol)
//					result.Servers = nmapbanner
//					if strings.Contains(result.Servers.Banner, "HTTP") {
//						result.Servers.Name = "http"
//						result.Servername = "http"
//					} else {
//						result.Servername = nmapbanner.Name
//					}
//					// web端口探测
//					webresult := runner.FingerResult{}
//					if result.Servername == "http" {
//						webresult = runner.Request(target)
//					}
//					result.Webbanner = webresult
//					result.Ip = target.Ip
//					result.Port = target.Port
//					result.Protocol = strings.ToUpper(target.Protocol)
//					result.Domain = target.Domain
//
//					if webresult.Header != "" {
//						result.Banner = result.Webbanner.Header
//					} else {
//						result.Banner = nmapbanner.Banner
//					}
//					// ip信息查询
//					info, err := ipinfo.GetIpinfo(region, target.Ip)
//					if err != nil {
//						log2.Warn("[GetIPInfo]:", err)
//						fmt.Fprintln(color.Output, color.HiYellowString("[WARNING]"), "["+time.Now().Format("2006-01-02 15:04:05")+"]", "[GetIPInfo]:", err)
//					}
//					result.City = info.City
//					result.Region = info.Region
//					result.ISP = info.ISP
//					result.CityId = info.CityId
//					result.Province = info.Province
//					result.Country = info.Country
//					result.TargetId = target.Ip + "-" + target.Port + "-" + target.Domain
//					if result.Port == "80" {
//						result.Target = "http://www." + target.Domain
//					} else {
//						result.TargetId = "http://www." + target.Domain + ":" + target.Port
//					}
//					result.LastTime = time.Now().Format("2006-01-02 15:04:05")
//					return result
//				}
//			}
//		}
//	}
//	return nil
//}
