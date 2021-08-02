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
	Pid [32]byte 					`json:"pid"`
	K 	int 						`json:"k"`
	BlockNumber *big.Int 			`json:"block_number"`
	Challenge	[32]byte 			`json:"challenge"`
	Quality     []byte				`json:"quality"`
	Deadline	*big.Int 			`json:"deadline"`
	Proof  		[]byte	  			`json:"proof"`
}

func NewState(pid [32]byte,k int,number *big.Int,challenge [32]byte,baseTarget *big.Int,proof []byte,quality []byte)*State  {
	st := &State{
		Pid: pid,
		K: k,
		BlockNumber: number,
		Challenge: challenge,
		Proof: proof,
		Quality: quality,
	}
	st.Deadline = CalculateDeadline(challenge,quality,baseTarget)
	return st
}

func (s *State)GeyKey() string {
	key := sha256.Sum256(append(s.Pid[:],s.BlockNumber.Bytes()...))
	return hex.EncodeToString(key[:])
}

func (s *State)NextChallenge() [32]byte {
	return NextChallenge(s.Challenge,s.Quality)
}
