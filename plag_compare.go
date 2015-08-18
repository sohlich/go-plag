package main

import (
	log "github.com/Sirupsen/logrus"
	"github.com/sohlich/go-plag/parser"
)

func checkAssignment(assignment *Assignment) {
	//Obtain all assignment files
	submissionFiles, err := mongo.FindAllSubmissionsByAssignment(assignment.ID.Hex())
	if err != nil {
		log.Error(err)
		return
	}

	output := make([]OutputComparisonResult, 0)

	for i := 0; i < len(submissionFiles); i++ {
		for j := i + 1; j < len(submissionFiles); j++ {
			tuple := OutputComparisonResult{
				Files: []string{submissionFiles[i].ID.Hex(), submissionFiles[j].ID.Hex()},
			}
			output = append(output, tuple)
		}
	}

	for _, out := range output {
		log.Debugln(out)
	}
}

//Receives OutputComparisonResult
//from channel and return an output channel
//with filled entity with comparison result
func compareFiles(inputChannel <-chan OutputComparisonResult) <-chan OutputComparisonResult {
	outputChannel := make(chan OutputComparisonResult, 10)
	go func(inChan <-chan OutputComparisonResult) {
		for toCompare := range inChan {
			log.Debugf("Starting to compare {}", toCompare.Files[0])
			sbmsnOne, err := mongo.FindSubmissionFileById(toCompare.Files[0])
			if err != nil {
				log.Error(err)
				continue
			}
			sbmsnTwo, err := mongo.FindSubmissionFileById(toCompare.Files[1])
			if err != nil {
				log.Error(err)
				continue
			}
			toCompare.SimilarityIndex = parser.Jaccard.Compare(sbmsnOne.TokenMap, sbmsnTwo.TokenMap)
			outputChannel <- toCompare
		}
	}(inputChannel)

	return outputChannel
}
