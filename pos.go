package poc

import (
	"errors"
	"fmt"
	"github.com/gnc-project/poc/chiapos"
	"math/big"
)

var proofVerifier *chiapos.ProofVerifier

func init()  {
	proofVerifier = chiapos.NewProofVerifier()
}

func getProofVerifier() *chiapos.ProofVerifier {
	return proofVerifier
}

func ValidateDeadline(pid [32]byte,k int,proof []byte,challenge [32]byte,difficulty,elapsedTime *big.Int) (bool, error) {

	quality,err := GetVerifiedQuality(pid,k,proof,challenge)
	if err != nil {
		return false,err
	}

	deadline := CalculateDeadline(challenge,quality,difficulty)
	if elapsedTime.Cmp(deadline) < 0{
		return false,fmt.Errorf("invalid deadline (elapsedTime: %v,deadline: %v)", elapsedTime, deadline)
	}

	return true, nil
}

func GetVerifiedQuality(pid [32]byte,k int,proof []byte,challenge [32]byte) ([]byte, error)  {

	if k < chiapos.MinPlotSize || k > chiapos.MaxPlotSize {
		return nil, errors.New("invalid plot k size")
	}

	if pass := chiapos.PassPlotFilter(pid,challenge); !pass {
		return nil, errors.New("not passing plot filter")
	}

	pv := getProofVerifier()
	quality,err := pv.GetVerifiedQuality(pid[:],proof,challenge,k)
	if err != nil {
		return nil, err
	}

	if len(quality) == 0 {
		return nil, errors.New("empty chia pos quality")
	}

	return quality, nil
}

func GetQuality(prover *chiapos.DiskProver,challenge [32]byte)([][]byte, error) {

	if !chiapos.PassPlotFilter(prover.ID(),challenge){
		return nil,errors.New("not passing plot filter")
	}

	return prover.GetQualitiesForChallenge(challenge)
}