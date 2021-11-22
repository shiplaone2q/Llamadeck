package models

import (
	"context"
	"testing"
	"time"
)

func TestJobStore_JobStateTransitions(t *testing.T) {
	store := NewJobStore()
	defer store.Close()

	// Create job
	job, err := store.Create("org/model", "Q4_K_M")
	if err != nil {
		t.Fatalf("failed to create job: %v", err)
	}

	if job.Status != JobStatusQueued {
		t.Errorf("initial status = %v, want %v", job.Status, JobStatusQueued)
	}

	// Update to downloading
	store.UpdateStatus(job.ID, JobStatusDownloading)
	retrieved, _ := store.Get(job.ID)
	if retrieved.Status != JobStatusDownloading {
		t.Errorf("status after update = %v, want %v", retrieved.Status, JobStatusDownloading)
	}

	// Complete job
	store.Complete(job.ID)
	retrieved, _ = store.Get(job.ID)
	if retrieved.Status != JobStatusCompleted {
		t.Errorf("status after complete = %v, want %v", retrieved.Status, JobStatusCompleted)
	}

	if retrieved.CompletedAt == nil {
		t.Error("CompletedAt should be set after completion")
	}
}

func TestJobStore_FailJob(t *testing.T) {
	store := NewJobStore()
	defer store.Close()

	job, _ := store.Create("org/model", "Q4_K_M")

	errorMsg := "download failed: network timeout"
	store.Fail(job.ID, errorMsg)

	retrieved, _ := store.Get(job.ID)
	if retrieved.Status != JobStatusFailed {
		t.Errorf("status = %v, want %v", retrieved.Status, JobStatusFailed)
	}

	if retrieved.Error != errorMsg {
		t.Errorf("error = %q, want %q", retrieved.Error, errorMsg)
