package calc

import (
	"runtime"
	"sync"
)

type pWork struct {
	from int
	to   int
}

func pearson(timeSeries [][600]float32, stats []LinStatEle, matBuffer [][13362]float32, order <-chan int, wg *sync.WaitGroup) {
	for {
		work, ok := <-order
		if ok {
			for i := work; i < 13362; i++ {
				var accProd float32
				for t := 0; t < 600; t++ {
					accProd += timeSeries[work][t] * timeSeries[i][t]
				}

				cov := (accProd / 600) - (stats[work].avg * stats[i].avg)

				pearson := cov / (stats[work].stddev * stats[i].stddev)

				matBuffer[work][i] += pearson
				if work != i {
					matBuffer[i][work] += pearson
				}
			}

			wg.Done()
		} else {
			break
		}
	}

	return
}

// DoPearson does Pearson's correlation calculation
func DoPearson(timeSeries [][600]float32, stats []LinStatEle, matBuffer [][13362]float32) {
	order := make(chan int, runtime.NumCPU())

	var wg sync.WaitGroup

	for i := 0; i < runtime.NumCPU(); i++ {
		go pearson(timeSeries, stats, matBuffer, order, &wg)
	}

	wg.Add(13362)
	for i := 0; i < 13362; i++ {
		order <- i
	}
	wg.Wait()
	close(order)
	return
}
