package poc

import (
	"crypto/sha256"
	"fmt"
	"math/big"
	"math/rand"
	"testing"
)

func TestGenerator(t *testing.T) {

	diff := big.NewInt(0).SetUint64(130310957598748952)
	for i:=0;i<100000;i++  {
		pid := sha256.Sum256([]byte(fmt.Sprintf("sdaf%d",rand.Intn(1000000000))))
		ch := sha256.Sum256([]byte(fmt.Sprintf("sdaf%d",rand.Intn(1000000000))))
		d := CalculateDeadline(pid,ch[:],diff)
		if d.Cmp(big.NewInt(18)) < 0 {
			aj := big.NewInt(0).Div(diff,big.NewInt(2))
			diff.Add(diff,aj)
			fmt.Println("diff--------->",diff)
		}else if d.Cmp(big.NewInt(18)) > 0 {
			aj := big.NewInt(0).Div(diff,big.NewInt(2))
			diff.Sub(diff,aj)
			fmt.Println("diff--------->",diff)
		}
		fmt.Println("deadline------>",d)
	}
	//115792089237316195423570985008687907853269984665640564039457584007913129639936
}
