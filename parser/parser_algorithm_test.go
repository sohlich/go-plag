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

//Test basic jaccard index
func TestJaccardBasic(t *testing.T) {
	mapA := map[string]int{
		"10": 2,
		"12": 1,
		"13": 1,
	}
	mapB := map[string]int{
		"13": 5,
		"10": 1,
	}
	result := Jaccard.Compare(mapA, mapB)

	assert.EqualValues(t, float32(1)/float32(4), result, float64(0.0001))
}

//Test empty map
func TestJaccardEmptyMap(t *testing.T) {
	mapA := make(map[string]int)
	mapB := make(map[string]int)
	result := Jaccard.Compare(mapA, mapB)
	assert.EqualValues(t, 0, result, float64(0.0001))
}

//Test nil value for jaccard function
func TestJaccardNillMap(t *testing.T) {
	result := Jaccard.Compare(nil, nil)
	assert.EqualValues(t, 0, result, float64(0.0001))
}
