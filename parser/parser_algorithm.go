package parser

import (
	"errors"
	"hash/fnv"
)

func hash(s string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(s))
	return h.Sum32()
}

type FingerPrint struct {
	FingerPrint []int
}

type FingerPrintAlgorithm interface {
	processTokensToFingerPrint([]string) *FingerPrint
}

type Winnowing struct{}

func (w *Winnowing) processTokensToFingerPrint(tokens []string) (*FingerPrint, error) {

	return nil, errors.New("Not implemeted")
}
