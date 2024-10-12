package tlp

import (
	"log"

	GTST "com.jadud.search.six/pkg/types"
)

func ShowMsg(ch_in <-chan GTST.JSON, ch_out chan<- GTST.JSON) {
	for {
		msg := <-ch_in
		log.Println(string(msg))
		ch_out <- msg
	}
}

func BlackHole[T any](ch_in <-chan T) {
	for {
		<-ch_in
	}
}

func NoisyBlackHole[T any](ch_in <-chan T) {
	for {
		msg := <-ch_in
		log.Printf("%s", msg)
	}
}
