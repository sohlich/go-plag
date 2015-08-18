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

func TestJaccard(t *testing.T) {

	mapA := map[uint32]int{
		uint32(10): 2,
		uint32(12): 1,
		uint32(13): 1,
	}

	mapB := map[uint32]int{
		uint32(13): 5,
		uint32(10): 1,
	}

	result := Jaccard.Compare(mapA, mapB)

	assert.EqualValues(t, float32(1)/float32(4), result, float64(0.0001))
}
