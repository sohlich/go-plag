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
	processTokensToFingerPrint([]string) *FingerPrint
}

type Winnowing struct {
	windowLengt int
}

func (alg *Winnowing) processTokensToFingerPrint(tokens []uint32) (*FingerPrint, error) {
	fp := &FingerPrint{make([]uint32, 0)}
	n := len(tokens)
	lastMin := -1
	for i := 0; i < n-alg.windowLengt+1; i++ {
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
