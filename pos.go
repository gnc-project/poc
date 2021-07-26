package poc

import (
	"crypto/sha256"
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


func NextChallenge(challenge [32]byte,quality []byte) [32]byte {
	return sha256.Sum256(append(challenge[:],quality...))
}

func ValidateDeadline(pid [32]byte,k int,proof []byte,challenge,parentCh [32]byte,difficulty,elapsedTime *big.Int) (bool, error) {

	quality,err := GetVerifiedQuality(pid,k,proof,parentCh)
	if err != nil {
		return false,err
	}

	if challenge != NextChallenge(parentCh,quality) {
		return false,fmt.Errorf("invalid challenge")
	}

	deadline := CalculateDeadline(pid,parentCh[:],difficulty)
	if elapsedTime.Cmp(deadline) < 0{
		return false,fmt.Errorf("invalid deadline (elapsedTime: %v,deadline: %v)", elapsedTime, deadline)
	}

	return true, nil
}

func GetVerifiedQuality(pid [32]byte,k int,proof []byte,challenge [32]byte) ([]byte, error)  {

	if k < chiapos.MinPlotSize || k > chiapos.MaxPlotSize {
		return nil, errors.New("invalid plot k size")
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

	if prover.Size() < chiapos.MinPlotSize || prover.Size() > chiapos.MaxPlotSize {
		return nil, errors.New("invalid plot k size")
	}

	return prover.GetQualitiesForChallenge(challenge)
}
