package gomet

import (
	"testing"
)

func Test_NewWorkerInChan(t *testing.T) {
	b := NewBroadcaster()

	id, workerInChan := b.NewWorkerInChan()
	if id <= int64(0) {
		t.Fatalf("Failed to generate a unique id")
	}
	if workerInChan == nil {
		t.Fatalf("Failed to generate a worker input channel")
	}
}

func Test_DeleteWorkerInChan(t *testing.T) {
	b := NewBroadcaster()

	id, _ := b.NewWorkerInChan()
	if id <= int64(0) {
		t.Fatalf("Failed to generate a unique id")
	}

	mapLength := 0
	b.WorkerInChans.Range(func(_, _ interface{}) bool {
		mapLength++
		return true
	})

	if mapLength != 1 {
		t.Fatalf("Failed to store the generated worker input channel")
	}

	b.DeleteWorkerInChan(id)

	mapLength = 0
	b.WorkerInChans.Range(func(_, _ interface{}) bool {
		mapLength++
		return true
	})

	if mapLength != 0 {
		t.Fatalf("Failed to store the generated worker input channel")
	}
}
