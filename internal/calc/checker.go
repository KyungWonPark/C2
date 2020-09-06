package calc

import (
	"fmt"
	"sync"
)

func check(timeSeries [][13362]float32, order <-chan int, wg *sync.WaitGroup) {
	for {
		index, ok := <-order
		if ok {
			for i, value := range timeSeries[index] {
				if value > 1 {
					fmt.Printf("ValueError: mat[%d][%d] is over 1!\n", index, i)
				} else if value < -1 {
					fmt.Printf("ValueError: mat[%d][%d] is under -1!\n", index, i)
				}
			}

			wg.Done()
		} else {
			break
		}
	}

	return
}

// DoCheck checks value range
func DoCheck(timeSeries [][13362]float32, workerConfig *Config) {
	numWorkers := (*workerConfig).NumComputer

	order := make(chan int, numWorkers)
	var wg sync.WaitGroup
	for i := 0; i < numWorkers; i++ {
		go check(timeSeries, order, &wg)
	}

	wg.Add(13362)
	for i := 0; i < 13362; i++ {
		order <- i
	}
	wg.Wait()
	close(order)
	return
}
