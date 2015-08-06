package raft

import (
	"math/rand"
	"time"
)

type Timer struct {
	originNum int
	countDown int
	//	sigBreak  bool
}

func (t *Timer) Start(countDown int) bool {
	t.originNum = countDown
	t.countDown = countDown

	for {
		time.Sleep(time.Millisecond)
		t.countDown--
		if t.countDown == 0 {
			return true
		}
	}

}

func (t *Timer) Reset() {
	t.countDown = t.originNum
}

func (node *RaftNode) StartTimer() {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	interval := 150 + r.Intn(150)
	node.timeup <- node.timer.Start(interval)
}
