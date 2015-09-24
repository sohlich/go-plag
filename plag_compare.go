package main

import (
	"sync"

	"github.com/sohlich/go-plag/parser"
	"golang.org/x/net/context"
)

var inProgressMap = make(map[string]context.CancelFunc)
var mutex = new(sync.Mutex)

//TODO implement assignment check
func checkAssignment(assignment *Assignment) int {

	//try to cancell previous running processes
	mutex.Lock()
	cancel, ok := inProgressMap[assignment.ID.Hex()]
	mutex.Unlock()

	if ok {
		Log.Debugln("Found previous running comparison job")
		cancel()
	}
	ctx, cancelFunc := context.WithCancel(context.TODO())
	inProgressMap[assignment.ID.Hex()] = cancelFunc

	//Obtain all assignment files
	submissionFiles, err := mongo.FindAllSubmissionsByAssignment(assignment.ID.Hex())
	if err != nil {
		Log.Error(err)
		return 0
	}
	processChanel := generateTuples(ctx, submissionFiles)
	outpuchannel := compareFiles(ctx, processChanel)

	compCount := 0

	for {
		select {
		case <-ctx.Done():
			Log.Infoln("Assignment check cancelled")
			return compCount
		case comparison, ok := <-outpuchannel:
			if !ok {
				Log.Infof("Comparison of %s done", assignment.ID.Hex())
				mongo.FindMaxSimilarityBySubmission(assignment.ID.Hex())
				return compCount
			}
			compCount++
			// if !math.IsNaN(float64(comparison.SimilarityIndex)) {
			_, err := mongo.Save(&comparison)
			if err != nil {
				Log.Error(err)
			}
			// }

		}
	}
}

//Generates non-repeating
//tuples from give array
func generateTuples(ctx context.Context, files []SubmissionFile) <-chan OutputComparisonResult {
	output := make(chan OutputComparisonResult)
	go func(chan OutputComparisonResult) {
		defer close(output)
		for i := 0; i < len(files); i++ {
			for j := i + 1; j < len(files); j++ {

				select {
				case <-ctx.Done():
					Log.Debugln("Generating of tuples canceled")
					return
				default:
					//Do not compare files form same submission
					if files[i].Submission == files[j].Submission {
						continue
					}

					tuple := OutputComparisonResult{
						Assignment:      files[i].Assignment,
						Files:           []string{files[i].ID.Hex(), files[j].ID.Hex()},
						Submissions:     []string{files[i].Submission, files[j].Submission},
						Tokens:          []map[string]int{files[i].TokenMap, files[j].TokenMap},
						SimilarityIndex: -1,
					}
					output <- tuple
				}
			}
		}
	}(output)
	return output
}

//Receives OutputComparisonResult
//from channel and return an output channel
//with filled entity with comparison result
func compareFiles(ctx context.Context, inputChannel <-chan OutputComparisonResult) <-chan OutputComparisonResult {
	outputChannel := make(chan OutputComparisonResult)
	go func(inChan <-chan OutputComparisonResult) {
		defer close(outputChannel)
		for {
			select {
			case <-ctx.Done():
				Log.Infoln("Cancelling comparison")
				return
			case toCompare, ok := <-inChan:
				if !ok {
					return
				}
				//Metrics
				comparison_count.Add(1)
				Log.Debugf("Starting to compare {}", toCompare.Tokens[0])
				sbmsnOne := toCompare.Tokens[0]
				sbmsnTwo := toCompare.Tokens[1]
				toCompare.SimilarityIndex = parser.Jaccard.Compare(sbmsnOne, sbmsnTwo)
				outputChannel <- toCompare

			}
		}
	}(inputChannel)

	return outputChannel
}
