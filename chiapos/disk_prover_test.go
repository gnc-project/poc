package chiapos_test

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/gnc-project/poc/chiapos"
	"math/big"
	"math/rand"
	"strings"
	"testing"
	"time"
)

func TestDiskProver_Close(t *testing.T) {

	prover,err := chiapos.NewDiskProver("/nvme/plots/plot-k32-2021-07-12-14-52-fa216e51dafd2a1bb964bd4184a4168ced78b0ac51de2eec114890807b8df5ce.plot",true)
	if err != nil {
		panic(err)
	}
	id := prover.ID()
	mem := prover.Memo()
	info := prover.PlotInfo()
	fname := prover.Filename()
	fmt.Printf("id=%s mem=%s fname=%s info=%v\n",hex.EncodeToString(id[:]),hex.EncodeToString(mem),fname,info)

	challenge := "66687aadf862bd776c8fc18b8e9f8e20089714856ee233b3902a591d0d5f2925"

	ch,err:= hex.DecodeString(challenge)

	var c [32]byte
	copy(c[:],ch[:32])

	v := chiapos.NewProofVerifier()
	defer v.Free()
	d := big.NewInt(0).SetUint64(18446744073709551615)
	for  {
		hash := sha256.Sum256([]byte(fmt.Sprintf("%d",rand.Intn(1000000))))

		quas,err := prover.GetQualitiesForChallenge(hash)
		if err != nil {
			time.Sleep(1*time.Second)
			continue
		}

		for _,v := range quas {
			q := big.NewInt(0).SetBytes(v)
			r := q.Div(q,d).Div(q,d).Div(q,d)
			fmt.Println(d.Div(d,r))
		}

		//proof,err := prover.GetFullProof(hash,0)
		//if err != nil {
		//	fmt.Printf("not found proof hash=%s\n",hex.EncodeToString(hash[:]))
		//	//time.Sleep(1*time.Second)
		//	continue
		//}
		//fmt.Printf("len--->%d hash=%s proof--->%s\n", len(proof),hex.EncodeToString(hash[:]),hex.EncodeToString(proof))
		//b,err := v.GetVerifiedQuality(id[:],proof,c,32)
		//if err != nil {
		//	panic(err)
		//}
		////fmt.Println(big.NewInt(0).SetBytes(b))
		//fmt.Println("b------>",hex.EncodeToString(b))
		//time.Sleep(1*time.Second)
	}


	for i:=0;i>=0;i++{
		qualities, err := prover.GetQualitiesForChallenge(c)
		if err !=nil {
			panic(err)
		}
		fmt.Println("qualities----------->",len(qualities))

		v := chiapos.NewProofVerifier()
		defer v.Free()
		for i,_ := range qualities {
			//k := prover.GetSize()
			proof,err := prover.GetFullProof(c,uint32(i))
			if err != nil {
				panic(err)
			}
			fmt.Printf("len--->%d proof--->%s\n", len(proof),hex.EncodeToString(proof))
			b,err := v.GetVerifiedQuality(id[:],proof,c,32)
			if err != nil {
				panic(err)
			}
			//fmt.Println(big.NewInt(0).SetBytes(b))
			fmt.Println("b------>",hex.EncodeToString(b))
		}
	}



}

func TestDiskProver(t *testing.T) {
	dp, err := chiapos.NewDiskProver("./testplots/plot.dat", false)
	if err != nil {
		t.Fatalf("NewDiskProver: %s", err)
	}
	defer dp.Close()

	id := dp.ID()
	size := dp.Size()
	fmt.Println(dp.Filename(), id, size)

	challenge := [...]byte{2, 47, 180, 44, 8, 193, 45, 227, 166, 175, 5, 56, 128, 25, 152, 6, 83, 46, 121, 81, 95, 148, 232, 52, 97, 97, 33, 1, 249, 65, 47, 158}
	qs, err := dp.GetQualitiesForChallenge(challenge)
	if err != nil {
		t.Fatalf("GetQualitiesForChallenge: %s", err)
	}
	fmt.Println(qs)

	fp, err := dp.GetFullProof(challenge, 300)
	if err != nil {
		if strings.Contains(err.Error(), "No proof of space for this challenge") {
			t.Logf("GetFullProof: %s", err)
		} else {
			t.Fatalf("GetFullProof: %s", err)
		}
	}
	fmt.Println(fp)
}

// func TestGetInfo(t *testing.T) {
// 	dp, err := chiapos.NewDiskProver("./testplots/plot.dat", false)
// 	if err != nil {
// 		t.Fatalf("NewDiskProver: %s", err)
// 	}
// 	defer dp.Close()

// 	pi := dp.PlotInfo()
// 	fmt.Println(pi.FarmerPublicKey.Bytes(), pi.PoolPublicKey.Bytes())
// }
