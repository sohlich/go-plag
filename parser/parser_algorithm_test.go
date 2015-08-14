package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWinnowing(t *testing.T) {
	input := []uint32{77, 72, 42, 17, 98, 50, 17, 98, 8, 88, 67, 39, 77, 72, 42, 17, 98}
	expected := []uint32{17, 17, 8, 39, 17}

	winnowing := Winnowing{4}
	fp, err := winnowing.processTokensToFingerPrint(input)

	if err != nil || !assert.ObjectsAreEqualValues(expected, fp.FingerPrint) {
		t.Errorf("TestWinnowing - Arrays not equal: %s %s", expected, fp.FingerPrint)
	}

}
