package models

type ProgressTracker struct {
	jobStore *JobStore
}

func NewProgressTracker(jobStore *JobStore) *ProgressTracker {
	return &ProgressTracker{
		jobStore: jobStore,
	}
}

func (pt *ProgressTracker) Track(jobID string, progress <-chan int64) {
	for bytes := range progress {
		pt.jobStore.mutex.Lock()
