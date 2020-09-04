package calc

import (
	"math"
	"sync"
)

func zScoring(timeSeries [][600]float32, order <-chan int, wg *sync.WaitGroup) {
	for {
		index, ok := <-order
		if ok {
			var valAcc float32
			var sqrAcc float32

			for _, value := range timeSeries[index] {
				valAcc += value
				sqrAcc += value * value
			}

			avg := valAcc / 600
			sqrMean := sqrAcc / 600
			stddev := float32(math.Sqrt(float64(sqrMean) - float64(avg*avg)))

			for i, value := range timeSeries[index] {
				timeSeries[index][i] = (value - avg) / stddev
			}

			wg.Done()
		} else {
			break
		}
	}

	return
}

func doZScoring(timeSeries [][600]float32, workerConfig *Config) {
	numWorkers := (*workerConfig).NumComputer

	order := make(chan int, numWorkers)
	var wg sync.WaitGroup
	for i := 0; i < numWorkers; i++ {
		go zScoring(timeSeries, order, &wg)
	}

	wg.Add(13362)
	for i := 0; i < 13362; i++ {
		order <- i
	}
	wg.Wait()
	close(order)
	return
}
