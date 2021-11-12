package difficulty

import (
	"fmt"
	"math/big"
	"testing"
	"time"
)

func TestCalcNextRequiredDifficulty(t *testing.T) {
	last := int64(1636687068)
	lastTime := time.Unix(last,0)
	blockTime := time.Unix(last+30,0)
	fmt.Println("elapsed",blockTime.Sub(lastTime))
	diff := CalcNextRequiredDifficulty(lastTime,big.NewInt(100000000),blockTime)
	fmt.Println("diff",diff)
}
