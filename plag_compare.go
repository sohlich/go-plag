package main

import (
	"errors"
	"sync"

	"github.com/sohlich/go-plag/parser"
	"golang.org/x/net/context"
)

//Errors
var ErrNoAssignment = errors.New("No assignment found")

var inProgressMap = NewJobMap()

//JobMap is struct to hold comparison job contexts.
//In case that some submission come before the
//end of comparison job. The job is cancelled and
//new job will be lounched to be sure all files are compared.
type JobMap struct {
	mutex *sync.Mutex
	//holds cancel func for assignmentID
	jobMap map[string]context.CancelFunc
}

//NewJobMap creates new instance
//of JobMap structure. Which holds the
//current running comparison jobs for
//assignmentID
func NewJobMap() *JobMap {
	return &JobMap{
		new(sync.Mutex),
		make(map[string]context.CancelFunc),
	}
}

//GetJob returns the context of jobrunning job.
//To be able to cancel the job.
//Wraps the map two value assignment, to be
// thread safe.
func (j *JobMap) GetJob(assignmentID string) (context.CancelFunc, bool) {
	j.mutex.Lock()
	cancel, ok := j.jobMap[assignmentID]
	j.mutex.Unlock()
	return cancel, ok
}

//Directly tries to cancel the comparison job
//on assignmentID key. If the key is present in map,
//cancel function in context.Context will be called.
func (j *JobMap) TryCancelJobFor(assignmentID string) bool {
	cancel, ok := j.GetJob(assignmentID)
	if ok {
		Log.Debugln("Found previous running comparison job")
		cancel()
	}
	return ok
}

//Puts the cancel function of give assignmentID
//to storage.
func (j *JobMap) PutJob(assignmentID string, cancel context.CancelFunc) {
	j.mutex.Lock()
	j.jobMap[assignmentID] = cancel
	j.mutex.Unlock()
}

//Check the assignment for plagiats by comparing all
//submissions in assignment.
func checkAssignment(assignment *Assignment) error {

	//Handle null assignment
	if assignment == nil {
		return ErrNoAssignment
	}

	assignmentID := assignment.ID.Hex()

	//try to cancell previous running processes
	inProgressMap.TryCancelJobFor(assignmentID)

	ctx, cancelFunc := context.WithCancel(context.TODO())
	inProgressMap.PutJob(assignmentID, cancelFunc)

	//Obtain all assignment files
	submissionFiles, err :=
		mongo.FindAllSubmissionsByAssignment(assignmentID)
	if err != nil {
		return err
	}

	//Pipe the tuples -> comparison
	processChanel := generateTuples(ctx, submissionFiles)
	outpuchannel := compareFiles(ctx, processChanel)

	//Wait to process all comparisons
	for {
		select {
		case <-ctx.Done():
			Log.Infoln("Assignment check cancelled")
			return nil
		case comparison, isOpen := <-outpuchannel:
			if !isOpen {
				//iIf all comparisons are done
				//and channel is closed
				Log.Infof("Comparison of %s done", assignmentID)
				apacErr := syncWithApac(assignmentID)
				if apacErr != nil {
					Log.Error(apacErr)
				}
				return nil
			}
			_, err := mongo.Save(&comparison)
			if err != nil {
				Log.Error(err)
			}

		}
	}

	return nil
}

//Generates non-repeating
//tuples from give array
func generateTuples(ctx context.Context,
	files []SubmissionFile) <-chan OutputComparisonResult {

	output := make(chan OutputComparisonResult)

	go func(chan OutputComparisonResult) {
		defer close(output)
		for i := 0; i < len(files); i++ {
			for j := i + 1; j < len(files); j++ {
				select {
				case <-ctx.Done():
					Log.Debugln("Cancelling tuple generation")
					return
				default:
					//Do not compare files form same
					//submission and owner
					noCmpr := files[i].Submission == files[j].Submission
					noCmpr = noCmpr || files[i].Owner == files[j].Owner
					if noCmpr {
						continue
					}

					outFls := []string{files[i].ID.Hex(), files[j].ID.Hex()}
					outSmsn := []string{files[i].Submission, files[j].Submission}
					outTkns := []map[string]int{files[i].TokenMap, files[j].TokenMap}

					output <- OutputComparisonResult{
						Assignment:      files[i].Assignment,
						Files:           outFls,
						Submissions:     outSmsn,
						Tokens:          outTkns,
						SimilarityIndex: -1,
					}
				}
			}
		}
	}(output)
	return output
}

//Receives OutputComparisonResult
//from channel and return an output channel
//with filled entity with comparison result
func compareFiles(ctx context.Context,
	inputChannel <-chan OutputComparisonResult) <-chan OutputComparisonResult {

	outputChannel := make(chan OutputComparisonResult)

	go func(inChan <-chan OutputComparisonResult) {
		defer close(outputChannel)
		for {
			select {
			case <-ctx.Done():
				Log.Debugln("Cancelling comparison")
				return
			case toCompare, ok := <-inChan:
				if !ok {
					return
				}
				//Metrics
				metrics.ComparisonInc()
				sbmsnOne := toCompare.Tokens[0]
				sbmsnTwo := toCompare.Tokens[1]
				toCompare.SimilarityIndex =
					parser.Jaccard.Compare(sbmsnOne, sbmsnTwo)
				Log.Debugf("Comparison of %v", toCompare)
				outputChannel <- toCompare

			}
		}
	}(inputChannel)

	return outputChannel
}
