package parser

import (
	"math"
	//"errors"
	"hash/fnv"
)

func hash(s string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(s))
	return h.Sum32()
}

type FingerPrint struct {
	FingerPrint []uint32
}

type FingerPrintAlgorithm interface {
	processTokensToFingerPrint([]string) (*FingerPrint, error)
}

type FingerPrintComparator interface {
	compare([]uint32, []uint32) (float32, error)
}

//Winnowing algorithm to fingerprint documents
//selects minimum hash value from chosen window as
//fingerprint values
type Winnowing struct {
	windowLengt int
}

//Creates fingerprint using winnowing algorithm
func (alg *Winnowing) processTokensToFingerPrint(tokens []uint32) (*FingerPrint, error) {
	fp := &FingerPrint{make([]uint32, 0)}
	n := len(tokens)
	lastMin := -1
	for i := 0; i <= n-alg.windowLengt; i++ {
		min := uint32(math.MaxUint32)
		index := 0
		for j := 0; j < alg.windowLengt; j++ {
			if tokens[i+j] < uint32(min) {
				min = uint32(tokens[i+j])
				index = i + j
			}
		}

		if lastMin != index {
			fp.FingerPrint = append(fp.FingerPrint, min)
			lastMin = index
		}
	}
	return fp, nil
}

//JaccardIndex computes the similarity of two sets
//based on their content. This special case
//compares not only the the "type" of
//the item but also the count in set.
var Jaccard *jaccardIndexComparator = &jaccardIndexComparator{}

type jaccardIndexComparator struct{}

func (j *jaccardIndexComparator) compare(mapA map[uint32]int, mapB map[uint32]int) float32 {

	//Compute intersection
	intrSec := 0
	tCntA := 0
	tCntB := 0
	for item, val := range mapA {
		tCntA += val
		valB := mapB[item]
		tCntB += valB
		intrSec += int(math.Min(float64(val), float64(valB)))
	}

	totalBlcks := tCntA + tCntB - intrSec
	jaccardIndx := float32(intrSec) / float32(totalBlcks)

	return jaccardIndx
}
