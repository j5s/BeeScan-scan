package node

import (
	"BeeScan-scan/pkg/job"
	redis2 "github.com/go-redis/redis"
	"time"
)

/*
创建人员：云深不知处
创建时间：2022/1/5
程序功能：节点
*/

var state map[string]interface{}

// NodeRegister 节点注册
func NodeRegister(c *redis2.Client, nodename string) {
	nodeflag := c.Exists(nodename).Val()
	if nodeflag == 0 {
		state = make(map[string]interface{})
		state["tasks"] = 0
		state["running"] = 0
		state["finished"] = 0
		state["lasttime"] = time.Now().Format("2006-01-02 15:04:05")
		c.HMSet(nodename, state)
	} else {
		c.HSet(nodename, "lasttime", time.Now().Format("2006-01-02 15:04:05"))
	}
}

// NodeUpdate 节点当前任务执行状况
func NodeUpdate(c *redis2.Client, nodename string, nodestate *job.NodeState) {
	c.HSet(nodename, "tasks", nodestate.Tasks)
	c.HSet(nodename, "running", nodestate.Running)
	c.HSet(nodename, "finished", nodestate.Finished)
	c.HSet(nodename, "lasttime", time.Now().Format("2006-01-02 15:04:05"))
}

// TaskUpdate 任务状态更新
func TaskUpdate(c *redis2.Client, task job.TaskState) {
	c.HSet(task.Name, "running", task.Running)
	c.HSet(task.Name, "finished", task.Finished)
	c.HSet(task.Name, "lasttime", time.Now().Format("2006-01-02 15:04:05"))
}
