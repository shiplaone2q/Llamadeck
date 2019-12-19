package models

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"sync"
	"time"
)

type JobStatus string

const (
	JobStatusQueued      JobStatus = "queued"
	JobStatusDownloading JobStatus = "downloading"
	JobStatusCompleted   JobStatus = "completed"
	JobStatusFailed      JobStatus = "failed"
	JobStatusCancelled   JobStatus = "cancelled"
)

type Progress struct {
	BytesDownloaded int64  `json:"bytes_downloaded"`
	TotalBytes      int64  `json:"total_bytes"`
	CurrentFile     string `json:"current_file"`
}

type Job struct {
	ID          string             `json:"id"`
	Repo        string             `json:"repo"`
	Tag         string             `json:"tag"`
	Status      JobStatus          `json:"status"`
	Progress    Progress           `json:"progress"`
	Error       string             `json:"error,omitempty"`
	CreatedAt   time.Time          `json:"created_at"`
	CompletedAt *time.Time         `json:"completed_at,omitempty"`
	CancelFunc  context.CancelFunc `json:"-"`
	ModelPath   string             `json:"model_path,omitempty"`
}

type JobStore struct {
	jobs   map[string]*Job
	mutex  sync.RWMutex
	ctx    context.Context
	cancel context.CancelFunc
}

func NewJobStore() *JobStore {
	ctx, cancel := context.WithCancel(context.Background())

	s := &JobStore{
		jobs:   make(map[string]*Job),
		ctx:    ctx,
		cancel: cancel,
	}

	// Start background job cleanup
	go s.cleanupLoop()

	return s
}

func (s *JobStore) Create(repo, tag string) (*Job, error) {
	jobID, err := s.generateJobID()
	if err != nil {
		return nil, fmt.Errorf("failed to generate job ID: %w", err)
	}

