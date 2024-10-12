package backingqueue

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/edwingeng/deque/v2"
	"github.com/robfig/cron/v3"
	"golang.org/x/exp/utf8string"
)

type SaveConfig[T any] interface {
	Save(qc BQ[T])
	Load()
	GetFrequency() string
}

type FileSaveConfig[T any] struct {
	FilePath string
	Cron     string
}

func (fsc *FileSaveConfig[T]) Save(qc BQ[T]) {
	os.MkdirAll(fsc.FilePath, 0755)
	for name, q := range qc.Queues {
		raw := q.Dump()
		data, _ := json.MarshalIndent(raw, "", " ")
		dest := filepath.Join(fsc.FilePath, fmt.Sprintf("%s.json", name))
		_ = os.WriteFile(dest, data, 0644)

	}
}

func (fsc *FileSaveConfig[T]) Load() {

}

func (fsc *FileSaveConfig[T]) GetFrequency() string {
	return fsc.Cron
}

type qMap[T any] map[string]*deque.Deque[T]

type BQ[T any] struct {
	SC     SaveConfig[T]
	Queues qMap[T]
	cron   *cron.Cron
}

func NewBackingQueue[T any](sc SaveConfig[T]) *BQ[T] {
	bq := BQ[T]{
		SC:     sc,
		Queues: make(qMap[T], 0),
		cron:   cron.New(),
	}
	return &bq
}

func (qc BQ[T]) NewQueue(name string) {
	if utf8string.NewString(name).IsASCII() {
		if _, ok := qc.Queues[name]; !ok {
			qc.Queues[name] = deque.NewDeque[T]()
		}
		qc.cron.AddFunc(qc.SC.GetFrequency(), func() { qc.SC.Save(qc) })
		qc.cron.Start()
	}
}

// We only push to the back of the queue; FIFO.
func (qc BQ[T]) Push(name string, obj T) {
	qc.Queues[name].Enqueue(obj)
}

func (qc BQ[T]) Pop(name string) T {
	return qc.Queues[name].Dequeue()
}

func (qc BQ[T]) Length(name string) int {
	return qc.Queues[name].Len()
}

func (qc BQ[T]) IsEmpty(name string) bool {
	return qc.Queues[name].IsEmpty()
}

func (qc BQ[T]) Save(name string) {
	qc.SC.Save(qc)
}

func (qc BQ[T]) SaveAll(name string) {
	for name, _ := range qc.Queues {
		qc.Save(name)
	}
}
