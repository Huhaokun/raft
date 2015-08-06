package main

import (
	"fmt"
	"raft"
	"strconv"
)

func main() {
	// myNode := &raft.RaftNode{}
	// myNode.Init(os.Args[0])
	// myNode.RunRPCServer(os.Args[1])
	// fmt.Println(myNode.RunFollower())

	for i := 0; i < 5; i++ {
		fmt.Println(i)
		go func(i int) {
			myNode := &raft.RaftNode{}
			myNode.Init(strconv.Itoa(i))
			go myNode.RunRPCServer("2000" + strconv.Itoa(i))
			myNode.RunFollower()
		}(i)
	}
	ch := make(chan bool)
	<-ch

}
