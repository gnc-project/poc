package poc

import (
	"errors"
	"github.com/gnc-project/poc/chiapos"
	"math/big"
	"sync"
)


func VerifiedQuality(proof []byte,pid,challenge [32]byte,slot,height,k uint64) (*big.Int,error) {
	quality,err := GetVerifiedQuality(pid,int(k),proof,challenge)
	if err != nil {
		return nil,err
	}

	hashVal := HashValChia(quality, slot, height)
	q1 := Q1FactorChia(uint8(k))

	return GetQuality(q1, hashVal),nil
}

var lock sync.Mutex
func GetVerifiedQuality(pid [32]byte,k int,proof []byte,challenge [32]byte) ([]byte, error)  {
	lock.Lock()
	defer lock.Unlock()
	if k < chiapos.MinPlotSize || k > chiapos.MaxPlotSize {
		return nil, errors.New("invalid plot k size")
	}

	if pass := chiapos.PassPlotFilter(pid, challenge); !pass {
		return nil, errors.New("not passing plot filter")
	}

	verifier := chiapos.NewProofVerifier()
	defer verifier.Free()

	posChallenge := chiapos.CalculatePosChallenge(pid, challenge)

	quality,err := verifier.GetVerifiedQuality(pid[:],proof,posChallenge,k)
	if err != nil {
		return nil, err
	}

	if len(quality) == 0 {
		return nil, errors.New("empty chia pos quality")
	}

	return quality, nil
}

func GetGNCProof(challenge [32]byte, index uint32,plot *chiapos.DiskProver) ([]byte, error)  {

	if plot.Size() < chiapos.MinPlotSize || plot.Size() > chiapos.MaxPlotSize {
		return nil, errors.New("invalid plot k size")
	}

	if pass := chiapos.PassPlotFilter(plot.ID(), challenge); !pass {
		return nil, errors.New("not passing plot filter")
	}

	posChallenge := chiapos.CalculatePosChallenge(plot.ID(), challenge)
	proof, err := plot.GetFullProof(posChallenge, index)
	if err != nil {
		return nil, err
	}
	return proof,nil
}