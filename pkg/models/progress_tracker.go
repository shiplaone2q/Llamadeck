package models

type ProgressTracker struct {
	jobStore *JobStore
}

func NewProgressTracker(jobStore *JobStore) *ProgressTracker {
	return &ProgressTracker{
		jobStore: jobStore,
	}
