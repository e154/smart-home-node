package driver

import (
	"testing"
	"fmt"
)

const expect = 250

func TestLRC(t *testing.T) {
	var i byte
	for i=1;i<255;i++ {
		fmt.Println(i + 4, "lrc:", LRC([]byte{i,3,0,0,0,1}))
	}
	lrc := LRC([]byte{1,3,0,0,0,1})
	if lrc != expect {
		t.Fatalf("lrc expected %v, actual %v", expect, lrc)
	}
}