package worker

import (
	"BeeScan-scan/pkg/job"
	"BeeScan-scan/pkg/runner"
	"log"
)

/*
创建人员：云深不知处
创建时间：2022/1/4
程序功能：任务池
*/

// 单个任务
func Woker(jobs <-chan *runner.Runner, results chan<- *runner.Output, f func(target *runner.Runner) *runner.Output) {

	for i := 1; i <= len(jobs); i++ {
		j := <-jobs
		result := f(j)
		results <- result
	}

}

// 任务池
func WokerPoll(WokerNumber int, jobs <-chan *runner.Runner, results chan *runner.Output, f func(target *runner.Runner) *runner.Output, nodestate *job.NodeState) {
	if WokerNumber < 1 {
		log.Fatalln("The wokernumber is error!")
	}

	for w := 1; w <= WokerNumber; w++ {
		nodestate.Running++
		go Woker(jobs, results, f)
		nodestate.Running--
		nodestate.Finished++
	}

}
