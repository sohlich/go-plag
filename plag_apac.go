package main

type SubmissionSimilarity struct {
	Uuid       string
	Similarity float64
}

type PlagiarismSync struct {
	Baseuuid    string
	Similarity  float64
	Submissions []SubmissionSimilarity
}
