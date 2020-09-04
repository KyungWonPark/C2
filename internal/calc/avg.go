package calc

import (
	"runtime"
	"sync"
)

func average(matBuffer [][13362]float32, divisor float32, order <-chan int, wg *sync.WaitGroup) {
	for {
		index, ok := <-order
		if ok {
			for i := range matBuffer[index] {
				matBuffer[index][i] = matBuffer[index][i] / divisor
			}

			wg.Done()
		} else {
			break
		}
	}
	return
}

// DoAverage does averaging
func DoAverage(matBuffer [][13362]float32, divisor float32) {
	order := make(chan int, runtime.NumCPU())
	var wg sync.WaitGroup

	for i := 0; i < runtime.NumCPU(); i++ {
		go average(matBuffer, divisor, order, &wg)
	}

	wg.Add(13362)
	for i := 0; i < 13362; i++ {
		order <- i
	}
	wg.Wait()

	close(order)
	return
}
