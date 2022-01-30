package job

import (
	"BeeScan-scan/pkg/runner"
	"time"
)

/*
创建人员：云深不知处
创建时间：2022/1/4
程序功能：任务队列
*/

type Job struct {
	Targets []*runner.Runner
	State   string
}

type NodeState struct {
	Tasks    int
	Running  int
	Finished int
}

type TaskState struct {
	Name      string
	TargetNum int
	Tasks     int
	Running   int
	Finished  int
	LastTime  time.Time
}

type Queue struct {
	Head   *QueueNode //头结点
	Tail   *QueueNode //尾结点
	Length int        //长度
}

type QueueNode struct {
	data string
	next *QueueNode
}

// NewQueue 创建链表队列
func NewQueue() (q *Queue) {
	head := QueueNode{}
	return &Queue{&head, &head, 0}
}

// Push 入队操作
func Push(q *Queue, data string) {
	if q == nil || q.Tail == nil {
		return
	}
	node := &QueueNode{data: data, next: nil}
	q.Tail.next = node
	q.Tail = node
	q.Length++
}

// Pop 出队操作
func Pop(q *Queue) string {
	if q == nil || q.Head == nil || q.Head == q.Tail {
		return ""
	}

	if q.Head.next == q.Tail {
		q.Tail = q.Head
	}
	data := q.Head.next.data
	q.Head.next = q.Head.next.next
	q.Length--
	return data
}
