package handlers

import (
	. "com.jadud.search.six/pkg/backingqueue"
	. "com.jadud.search.six/pkg/types"
)

var TheMultiqueue *BQ[Job]

const TheQueue string = "jobs"

func Init(filepath string, cron string) {
	TheMultiqueue = NewBackingQueue[Job](&FileSaveConfig[Job]{
		FilePath: filepath,
		Cron:     cron,
	})
	TheMultiqueue.NewQueue(TheQueue)
}
