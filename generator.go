package poc

import (
	"crypto/sha256"
	"encoding/hex"
	"math/big"
	"sync"
)

type Generator struct {
	sync.Map
}

var gen  *Generator

func init()  {
	gen = &Generator{}
}

func GetGenerator() *Generator {
	return gen
}

func (ge *Generator)SetCurrent(state *State)  {
	ge.Store(state.BlockNumber,state)
}

func (ge *Generator)GetCurrent(number *big.Int) *State  {
	if v, ok := ge.Load(number);ok{
		return v.(*State)
	}
	return nil
}

func (ge *Generator)AddState(state *State)  {
	ge.Store(state.GeyKey(),state)
}

func (ge *Generator)GetState(pn string)*State  {
	if st,ok := ge.Load(pn);ok{
		return st.(*State)
	}else {
		return nil
	}
}

type State struct {
	Pid []byte 					`json:"pid"`
	BlockNumber *big.Int 		`json:"block_number"`
	Challenge	[]byte 			`json:"challenge"`
	Deadline	*big.Int 		`json:"deadline"`
	Proof  		[]byte	  		`json:"proof"`
	Reward 		string  `json:"reward"`
}

func NewState(pid []byte,number *big.Int,challenge []byte,difficulty *big.Int,proof []byte,reward string)*State  {
	st := &State{
		Pid: pid,
		BlockNumber: number,
		Challenge: challenge,
		Proof: proof,
		Reward: reward,
	}
	st.Deadline = CalcuteDeadline(pid,challenge,difficulty)
	return st
}

func (s *State)GeyKey() string {
	key := sha256.Sum256(append(s.Pid,s.BlockNumber.Bytes()...))
	return hex.EncodeToString(key[:])
}

func (s *State)NextChallenge() []byte {
	challenge := sha256.Sum256(append(s.Challenge,s.Pid...))
	return challenge[:]
}

