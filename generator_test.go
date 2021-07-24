package poc

import (
	"fmt"
	"testing"
)

func TestGenerator(t *testing.T) {

	ge := GetGenerator()
	if ge == nil {
		fmt.Println("-------------------")
	}else {
		fmt.Println("*********************")
	}

}
