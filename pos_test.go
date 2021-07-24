package poc

import (
	"encoding/hex"
	"fmt"
	"testing"
)

func TestValidateProof(t *testing.T) {
	for i:=0;i < 100000000;i++ {
		cal()
	}
}

func cal()  {
	id, err := hex.DecodeString("fa216e51dafd2a1bb964bd4184a4168ced78b0ac51de2eec114890807b8df5ce")
	if err != nil {
		panic(err)
	}
	proof, err := hex.DecodeString("30b6c4eb073e53a07bdd729b00eee7c478a250bf1d372015ee591224eb3bdd3978454185ac8a0367f3b5330bbf53b80df85d8a1c659ec0edfbd9f106279b8e9f4bd323d5fc791a9839d0431f890a41879b4e15d365f11eaca660541c25c840af2b152ad08d8ecdcc544430c7a2c9f2bda6a0de54f55d2c9be25e4f09ac5402bf4e17512b7a1c0b3dbe65b87dafd98463f9ae758fb80f119aa219dcb73635e165fe528cd9de4a851161b772d8a82585cbce4c325dc45e721466967fd1c10cc0e932ed27e5f0889751a24ad2cc4ded842d1fbf1f04b5c31b9209dfbe05294302b82d9d7153a48cc5ce3ed30e8e4a1561e11d1580432ced5f35ac7661c18aaa21f3")
	if err != nil {
		panic(err)
	}
	ch, err := hex.DecodeString("66687aadf862bd776c8fc18b8e9f8e20089714856ee233b3902a591d0d5f2925")
	if err != nil {
		panic(err)
	}
	q,err := ValidateProof(id,proof,ch,32)
	if err != nil {
		panic(err)
	}

	fmt.Println(hex.EncodeToString(q))
}