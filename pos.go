package poc

import (
	"errors"
	"github.com/gnc-project/poc/chiapos"
)

var proofVerifier *chiapos.ProofVerifier

func init()  {
	proofVerifier = chiapos.NewProofVerifier()
}

func getProofVerifier() *chiapos.ProofVerifier {
	return proofVerifier
}

func ValidateProof(pid []byte,proof []byte,challenge []byte,k int) ([]byte, error) {
	if len(challenge) != 32 {
		return nil,errors.New("invalid challenge")
	}

	if k < chiapos.MinPlotSize || k > chiapos.MaxPlotSize {
		return nil, errors.New("invalid plot k size")
	}

	var ch [32]byte
	copy(ch[:],challenge[:32])
	pv := getProofVerifier()

	return pv.GetVerifiedQuality(pid,proof,ch,k)
}
