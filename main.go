package main

import (
	"sync"
	"time"
	"vitornilson1998/gostock/service"
	"vitornilson1998/gostock/utils"
)

func main() {
	init := time.Now()
	utils.CreateCsv()
	// Breaking the nasdaq.csv in chunks of 50 tickers.
	arr := utils.ChunkSlice(utils.ReadCsv(), 50)

	pipe := make(chan []string)

	workGroup := new(sync.WaitGroup)

	// Creating a list of goroutines based on size of chunks. Each chunk have one goroutine.
	for i := 0; i < len(arr); i++ {
		workGroup.Add(1)
		go service.ProcessChunk(pipe, workGroup)
	}

	// Appending chunks to channel and starting process.
	for _, list := range arr {
		pipe <- list
	}

	close(pipe)
	workGroup.Wait()

	println("Total time: ", int64(time.Since(init).Seconds()))

	// Running python code that refine search based on market analisys.
	utils.RefineResearch()
}
