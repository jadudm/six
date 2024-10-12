package main

import (
	"log"
	"sync"

	tlp "com.jadud.search.six/cmd/update-orchestrator/tlp"
	GSTS "com.jadud.search.six/pkg/types"
	"github.com/joho/godotenv"
)

func load_dotenv() {
	// https://github.com/joho/godotenv
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

// I'm the update-orchestrator.
// That means I do the following:
// 1. Check the HEAD queue for incoming jobs (these could be domains or pages... does not matter)
// 2. Depending on whether it is HTML/PDF/etc., put it on the right indexer queue.
// 3. Indexers pull work from their queues. They both walk for links (back to HEAD), or
//    they pull content and send it to the DB.

func main() {
	load_dotenv()

	ch_q_out := make(chan GSTS.JSON)
	// ch_sm_bh := make(chan GSTS.JSON)
	//ch_bh := make(chan GSTS.JSON)

	var wg sync.WaitGroup
	wg.Add(1)

	go tlp.CheckQueue("@every 1m", ch_q_out)
	//go tlp.ShowMsg(ch_q_out, ch_sm_bh)
	go tlp.NoisyBlackHole[GSTS.JSON](ch_q_out)

	wg.Wait()
	log.Println("we will never see this")
}
