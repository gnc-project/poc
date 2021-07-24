package poc

import (
	"errors"
	"fmt"
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

	q,err := pv.GetVerifiedQuality(pid,proof,ch,k)
	if err != nil {
		return nil, err
	}
	if q == nil {
		return nil, fmt.Errorf("q is nil")
	}
	if len(q) == 0 {
		return nil, fmt.Errorf("q is zero")
	}

	return q ,nil
}
