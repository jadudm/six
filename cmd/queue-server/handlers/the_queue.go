package handlers

import (
	. "com.jadud.search.six/pkg/backingqueue"
	. "com.jadud.search.six/pkg/types"
)

var Q *BQ[JSON]

func Init(filepath string, cron string) {
	Q = NewBackingQueue[JSON](&FileSaveConfig[JSON]{
		FilePath: filepath,
		Cron:     cron,
	})
}
