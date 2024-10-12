package backingqueue

import (
	"errors"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func setupTest[T any](t *testing.T) (func(t *testing.T), *BQ[T]) {
	// log.Println("setup test")

	return func(t *testing.T) {
			// log.Println("teardown test")
		}, NewBackingQueue[T](&FileSaveConfig[T]{
			FilePath: "test_data",
			Cron:     "*/1 0 0 0 0",
		})
}

func TestCreateQueue(t *testing.T) {
	_ = NewBackingQueue[int](&FileSaveConfig[int]{
		FilePath: "test_data",
		Cron:     "*/1 0 0 0 0",
	})
}

func TestEmptyQueue(t *testing.T) {
	bq := NewBackingQueue[int](&FileSaveConfig[int]{
		FilePath: "test_data",
		Cron:     "*/1 0 0 0 0",
	})
	bq.NewQueue("alice")
	assert.Equal(t, bq.IsEmpty("alice"), true)
}

func TestInsertIntoQueue(t *testing.T) {
	bq := NewBackingQueue[int](&FileSaveConfig[int]{
		FilePath: "test_data",
		Cron:     "*/1 0 0 0 0",
	})
	bq.NewQueue("alice")
	bq.Enqueue("alice", 1)
	assert.Equal(t, bq.IsEmpty("alice"), false)
}

func TestInsertIntoTwoQueues(t *testing.T) {
	bq := NewBackingQueue[int](&FileSaveConfig[int]{
		FilePath: "test_data",
		Cron:     "*/1 0 0 0 0",
	})
	bq.NewQueue("alice")
	bq.NewQueue("bob")
	bq.Enqueue("alice", 1)
	bq.Enqueue("alice", 2)
	bq.Enqueue("bob", 1)
	assert.Equal(t, bq.Length("alice"), 2)
	assert.Equal(t, bq.Length("bob"), 1)
}

func TestInsertStruct(t *testing.T) {
	type meaning struct {
		life int
	}
	bq := NewBackingQueue[meaning](&FileSaveConfig[meaning]{
		FilePath: "test_data",
		Cron:     "*/1 0 0 0 0",
	})
	bq.NewQueue("alice")
	bq.Enqueue("alice", meaning{life: 1})
	assert.Equal(t, bq.IsEmpty("alice"), false)
}

func TestRetrieveStruct(t *testing.T) {
	type meaning struct {
		life int
	}
	bq := NewBackingQueue[meaning](&FileSaveConfig[meaning]{
		FilePath: "test_data",
		Cron:     "*/1 0 0 0 0",
	})
	bq.NewQueue("alice")
	bq.Enqueue("alice", meaning{life: 42})
	f := bq.Dequeue("alice")
	assert.Equal(t, f.life, 42)
}

func TestSave(t *testing.T) {
	bq := NewBackingQueue[int](&FileSaveConfig[int]{
		FilePath: "test_data",
		Cron:     "@every 1m",
	})
	bq.NewQueue("alice")
	bq.Enqueue("alice", 1)
	bq.Enqueue("alice", 2)
	assert.Equal(t, bq.IsEmpty("alice"), false)
	bq.Save("alice")
}

func TestSaveStruct(t *testing.T) {
	type meaning struct {
		// Lowercase will not serialize.
		Life int
	}
	bq := NewBackingQueue[meaning](&FileSaveConfig[meaning]{
		FilePath: "test_data",
		Cron:     "*/1 0 0 0 0",
	})
	bq.NewQueue("beatrice")
	bq.Enqueue("beatrice", meaning{Life: 1})
	bq.Enqueue("beatrice", meaning{Life: 2})
	assert.Equal(t, bq.IsEmpty("beatrice"), false)
	bq.Save("beatrice")
	// https://stackoverflow.com/questions/12518876/how-to-check-if-a-file-exists-in-go
	if _, err := os.Stat("test_data/beatrice.json"); errors.Is(err, os.ErrNotExist) {
		t.Error("BQ not saved", "beatrice")
	}

}
