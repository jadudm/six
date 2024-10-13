package backingqueue

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"

	"github.com/edwingeng/deque/v2"
	"github.com/robfig/cron/v3"
)

var IsSafeQueueName = regexp.MustCompile(`^[a-zA-Z-_]+$`).MatchString

type qMap[T any] map[string]*deque.Deque[T]

type StrOrBytes interface {
	string | []byte
}

type BQ[T StrOrBytes] struct {
	SC     SaveConfig[T]
	Queues qMap[T]
	cron   *cron.Cron
}
type SaveConfig[T StrOrBytes] interface {
	Save(qc BQ[T])
	Load()
	GetFrequency() string
}

type FileSaveConfig[T StrOrBytes] struct {
	FilePath string
	Cron     string
}

func getQueueValues[T StrOrBytes](q *deque.Deque[T]) [][]byte {
	raw := q.Dump()
	rawb := make([][]byte, 0)
	for _, o := range raw {
		rawb = append(rawb, []byte(o))
	}
	return rawb
}

func (fsc *FileSaveConfig[T]) Save(qc BQ[T]) {

	os.MkdirAll(fsc.FilePath, 0755)
	for name, q := range qc.Queues {
		// log.Println("Saving queue:", name)
		dest := filepath.Join(fsc.FilePath, fmt.Sprintf("%s.json", name))
		f, err := os.Create(dest)
		if err != nil {
			log.Panic("cannot create queue log")
		}
		defer f.Close()
		raw := getQueueValues(q)
		f.WriteString("[\n")
		for _, o := range raw {
			f.Write([]byte(o))
			f.WriteString("\n")
		}

		f.WriteString("]\n")

	}
}

func (fsc *FileSaveConfig[T]) Load() {

}

func (fsc *FileSaveConfig[T]) GetFrequency() string {
	return fsc.Cron
}

func NewBackingQueue[T StrOrBytes](sc SaveConfig[T]) *BQ[T] {
	bq := BQ[T]{
		SC:     sc,
		Queues: make(qMap[T], 0),
		cron:   cron.New(),
	}
	return &bq
}

func (qc BQ[T]) NewQueue(name string) {

	// Make sure we're only ASCII for queue names
	//if utf8string.NewString(name).IsASCII() {
	if IsSafeQueueName(name) && (len(name) < 256) {
		// Only create the queue once
		if _, ok := qc.Queues[name]; !ok {
			qc.Queues[name] = deque.NewDeque[T]()
			qc.cron.AddFunc(qc.SC.GetFrequency(), func() { qc.SC.Save(qc) })
			qc.cron.Start()
			// Fast, for debugging.
			// go func() {
			// 	for {
			// 		time.Sleep(10 * time.Second)
			// 		qc.SC.Save(qc)
			// 	}
			// }()
		}
	}
}

func (qc BQ[T]) check_and_init(name string) {
	if _, ok := qc.Queues[name]; !ok {
		qc.NewQueue(name)
	}
}

// We only push to the back of the queue; FIFO.
func (qc BQ[T]) Enqueue(name string, obj T) {
	qc.check_and_init(name)
	//log.Println("EQ", obj)
	qc.Queues[name].Enqueue(obj)
}

func (qc BQ[T]) Dequeue(name string) T {
	return qc.Queues[name].Dequeue()
}

func (qc BQ[T]) Length(name string) int {
	if _, ok := qc.Queues[name]; !ok {
		return -1
	}
	return qc.Queues[name].Len()
}

func (qc BQ[T]) IsEmpty(name string) bool {
	qc.check_and_init(name)
	return qc.Queues[name].IsEmpty()
}

func (qc BQ[T]) Save(name string) {
	qc.check_and_init(name)
	qc.SC.Save(qc)
}

func (qc BQ[T]) SaveAll() {
	for name := range qc.Queues {
		qc.Save(name)
	}
}
