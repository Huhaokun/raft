package raft

import (
	"fmt"
	"log"
	"net"
	"net/rpc"
	"strconv"
)

type RaftNode struct {
	id          int
	state       State
	leaderState LeaderState
	timer       Timer
	winElection chan bool
	timeup      chan bool
	machines    map[int]string
}

func (node *RaftNode) Init(id string) {
	fmt.Println("Initing...")
	state := State{CurrentTerm: 0}
	node.timeup = make(chan bool, 1)
	node.winElection = make(chan bool, 1)
	node.id, _ = strconv.Atoi(id)
	node.state = state
	node.machines = map[int]string{0: "20000", 1: "20001", 2: "20002", 3: "20003", 4: "20004"}
}
func (node *RaftNode) RunFollower() string {
	fmt.Println(node.id, "am follower")
	go node.StartTimer()
	for {
		select {
		case <-node.timeup:
			// if time is up transfer to candidate
			return node.RunCandidate()
		}
	}
	return "I am follower"
}
func (node *RaftNode) RunCandidate() string {
	fmt.Println(node.id, "am the Candidate")

	//ask for vote
	StartElection := func() {
		node.state.CurrentTerm++
		voteNum := 1
		node.state.VotedFor = node.id // vote for itself
		for id, port := range node.machines {
			if id != node.id {
				args := &RequestVoteArgs{
					Term:        node.state.CurrentTerm,
					CandidateId: node.id,
				}
				var reply RequestVoteReply
				call(port, "RaftNode.RequestVoteRPC", &args, &reply)
				if reply.VoteGranted == true {
					voteNum++
				}
			}
		}
		if voteNum > 3 { //hardcode
			//<-channel
			fmt.Println(node.id, "win the election")
			node.winElection <- true
		}
	}

	for {
		go node.StartTimer()
		go StartElection()
		select {
		//		case <-node.rpcfromleader:
		//			return node.RunFollower()
		case <-node.timeup:
			fmt.Println("time is up,let's run another election")
			continue
		case <-node.winElection:
			return node.RunLeader()
			fmt.Println("this can't happen")
		}
	}

}

func (node *RaftNode) RunLeader() string {
	fmt.Println(node.id, "am leader")
	return "I am leader"
}
func (node *RaftNode) RunRPCServer(port string) {
	rpcs := rpc.NewServer()
	rpcs.Register(node)
	l, e := net.Listen("tcp", ":"+port)
	if e != nil {
		log.Fatal("listen error:", e)
	}

	for true {
		conn, err := l.Accept()
		if err == nil {
			go rpcs.ServeConn(conn)
		} else if err != nil {
			conn.Close()
		}
	}
}

func (node *RaftNode) RequestVoteRPC(args *RequestVoteArgs, reply *RequestVoteReply) error {
	if node.state.CurrentTerm <= args.Term &&
		(node.state.VotedFor == 0 ||
			node.state.VotedFor == args.CandidateId) { //TODO
		node.state.VotedFor = args.CandidateId
		reply.VoteGranted = true
		//make sure it's in follower state
		node.timer.Reset()
	} else {
		reply.VoteGranted = false
	}
	reply.Term = node.state.CurrentTerm
	return nil
}
