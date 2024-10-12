package handlers

import (
	. "com.jadud.search.six/pkg/backingqueue"
	. "com.jadud.search.six/pkg/types"
)

var TheQueue *BQ[Job]

func Init(filepath string, cron string) {
	TheQueue = NewBackingQueue[Job](&FileSaveConfig[Job]{
		FilePath: filepath,
		Cron:     cron,
	})

}
