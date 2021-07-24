package chiapos_test

import (
	"encoding/hex"
	"fmt"
	"github.com/gnc-project/poc/chiapos"
	"strings"
	"testing"
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

//113b9319a3342be618b94d54043f946db3dc61376be48223175dda7c96e798157396b314e90ec2142a5cb0db5756507da40b875b98b8b22efb2107256b43e0fddf463ef937bcbd5f88c274d7d5cdfd75d70076997e5996ad5afcc0d9db827ebdd61a685a66c0cf78b673af374b3e5ff49e21aac008553f81e3fef4080513b9911f3502f01bf089c788238e6da92a4b5b8bd014b6baa6295664176a79d21779c28b07b052957c2b4a3c73ddc2047e5f2e6d3ac33cff941dfa9b24094a8a24c1648068804f8d108c586a6b03b3a5ac237ac0cd1ba0682d6991332cd1e82b78e0ae6df608647275e3900beccf649ffb435a1ac60a3330784eb0c207488de6620c5088813699445533762371444020516372527512045505248809176954376175144524794675275
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

func TestGetInfo(t *testing.T) {
	dp, err := chiapos.NewDiskProver("./testplots/plot.dat", false)
	if err != nil {
		t.Fatalf("NewDiskProver: %s", err)
	}
	defer dp.Close()

	pi := dp.PlotInfo()
	fmt.Println(pi.FarmerPublicKey.Bytes(), pi.PoolPublicKey.Bytes())
}
