package tlp

import (
	"log"

	gtst "com.jadud.search.six/pkg/types"
)

func ShowMsg(ch_in <-chan gtst.JSON, ch_out chan<- gtst.JSON) {
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
		log.Println(msg)
	}
}
