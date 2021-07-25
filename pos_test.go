package poc

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/gnc-project/poc/chiapos"
	"math/big"
	"math/rand"
	"testing"
)

func TestValidateProof(t *testing.T) {
	for i := 0; i<10; i++ {
		hash := sha256.Sum256([]byte(fmt.Sprintf("%d",rand.Intn(100000000))))
		he := hex.EncodeToString(hash[:])
		fmt.Println(he)
	}

}

func cal()  {
	ch, err := hex.DecodeString("66687aadf862bd776c8fc18b8e9f8e20089714856ee233b3902a591d0d5f2925")
	if err != nil {
		panic(err)
	}
	challenge := BytesTo32(ch)
	prover,err := chiapos.NewDiskProver("/nvme/plots/plot-k32-2021-07-12-14-52-fa216e51dafd2a1bb964bd4184a4168ced78b0ac51de2eec114890807b8df5ce.plot",true)
	if err != nil {
		panic(err)
	}

	qualities,err := prover.GetQualitiesForChallenge(challenge)
	if err != nil {
		panic(err)
	}

	for i := 0;i< len(qualities);i++ {
		proof,err := prover.GetFullProof(challenge,uint32(i))
		if err != nil {
			panic(err)
		}
		b,err := ValidateDeadline(prover.ID(),32,proof,challenge,challenge,big.NewInt(17179869184),big.NewInt(20))
		if err != nil {
			panic(err)
		}
		fmt.Println(b)
	}

}

