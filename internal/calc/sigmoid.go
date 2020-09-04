package calc

import (
	"math"
	"sync"
)

func sigmoid(timeSeries [][600]float32, stats []LinStatEle, order <-chan int, wg *sync.WaitGroup) {
	for {
		index, ok := <-order
		if ok {
			var valAcc float32
			var sqrAcc float32

			for i, value := range timeSeries[index] {
				newVal := float32(2/(1+math.Exp(-float64(value))) - 1)
				valAcc += newVal
				sqrAcc += newVal * newVal
				timeSeries[index][i] = newVal
			}

			avgVal := valAcc / 600
			avgSqr := sqrAcc / 600
			stdDev := float32(math.Sqrt(float64(avgSqr) - float64(avgVal*avgVal)))

			stats[index] = LinStatEle{
				avg:    avgVal,
				stddev: stdDev,
			}

			wg.Done()
		} else {
			break
		}
	}

	return
}

func doSigmoid(timeSeries [][600]float32, stats []LinStatEle, workerConfig *Config) {
	numWorkers := (*workerConfig).NumComputer

	order := make(chan int, numWorkers)

	var wg sync.WaitGroup

	for i := 0; i < numWorkers; i++ {
		go sigmoid(timeSeries, stats, order, &wg)
	}

	wg.Add(13362)
	for i := 0; i < 13362; i++ {
		order <- i
	}
	wg.Wait()
	close(order)
	return
}
