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
		if j, ok := pt.jobStore.jobs[jobID]; ok {
			j.Progress.BytesDownloaded += bytes
		}
		pt.jobStore.mutex.Unlock()
	}
}

func (pt *ProgressTracker) AddToTotalBytes(jobID string, bytes int64) {
	pt.jobStore.mutex.Lock()
	defer pt.jobStore.mutex.Unlock()

