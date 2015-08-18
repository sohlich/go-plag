package main

import (
	log "github.com/Sirupsen/logrus"
	"github.com/sohlich/go-plag/parser"
)

//TODO implement assignment check
func checkAssignment(assignment *Assignment) {
	//Obtain all assignment files
	// submissionFiles, err := mongo.FindAllSubmissionsByAssignment(assignment.ID.Hex())
	// if err != nil {
	// 	log.Error(err)
	// 	return
	// }
}

//Generates non-repeating
//tuples from give array
func generateTuples(files []SubmissionFile) <-chan OutputComparisonResult {
	output := make(chan OutputComparisonResult, 10)
	go func(chan OutputComparisonResult) {
		for i := 0; i < len(files); i++ {
			for j := i + 1; j < len(files); j++ {
				tuple := OutputComparisonResult{
					Files: []string{files[i].ID.Hex(), files[j].ID.Hex()},
				}
				output <- tuple
			}
		}
		close(output)
	}(output)
	return output
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
