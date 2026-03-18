package backend

import (
	"testing"
	"time"
)

func TestQueueAsyncDownloadExecutesTask(t *testing.T) {
	done := make(chan struct{})

	if err := QueueAsyncDownload(func() {
		close(done)
	}); err != nil {
		t.Fatalf("QueueAsyncDownload returned error: %v", err)
	}

	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("queued async task was not executed")
	}
}

func TestQueueAsyncDownloadProcessesTasksInOrder(t *testing.T) {
	result := make(chan int, 2)

	if err := QueueAsyncDownload(func() {
		result <- 1
	}); err != nil {
		t.Fatalf("QueueAsyncDownload returned error: %v", err)
	}

	if err := QueueAsyncDownload(func() {
		result <- 2
	}); err != nil {
		t.Fatalf("QueueAsyncDownload returned error: %v", err)
	}

	select {
	case v := <-result:
		if v != 1 {
			t.Fatalf("first task executed out of order, got %d", v)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("first queued task did not run")
	}

	select {
	case v := <-result:
		if v != 2 {
			t.Fatalf("second task executed out of order, got %d", v)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("second queued task did not run")
	}
}
