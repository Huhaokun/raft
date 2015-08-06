package raft

import (
	"fmt"
	"net/rpc"
)

type State struct {
	CurrentTerm int   //last term server has seen (initialized to 0 on frist boot,increases monotonically)
	VotedFor    int   //candidatedId that received vote in current term(or null if none)
	Log         []int //log entries
	CommitIndex int   //index of highest log entry known to be committed (initialized to 0, increases monotonically)
	LastApplied int   //index of highest log entry applied to state machine (initialized to 0, increases monotonically)
}
type LeaderState struct {
	NextIndex  []int //for each server, index of the next log entry to send to that server (initialized to leader last log index + 1)
	MatchIndex []int //for each server, index of highest log entry known to be replicated on server (initialized to 0, increases monotonically)
}

type RequestVoteArgs struct {
	Term         int //candidate’s term
	CandidateId  int //candidate requesting vote
	LastLogIndex int //index of candidate’s last log entry
	LastLogTerm  int //term of candidate’s last log entry
}

type RequestVoteReply struct {
	Term        int  //currentTerm, for candidate to update itself
	VoteGranted bool //true means candidate received vote
}

func call(port string, rpcname string, args interface{}, reply interface{}) bool {
	serverAddress := "127.0.0.1"
	c, errx := rpc.Dial("tcp", serverAddress+":"+port)
	if errx != nil {
		return false
	}
	defer c.Close()

	err := c.Call(rpcname, args, reply)
	if err == nil {
		return true
	}
	fmt.Println(err)
	return false
}
