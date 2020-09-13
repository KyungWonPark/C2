package calc

import (
	"runtime"
	"sync"

	"gonum.org/v1/gonum/mat"
)

func thresholding(matBuffer [][13362]float32, dest *mat.SymDense, threshold float64, order <-chan int, wg *sync.WaitGroup) {
	for {
		work, ok := <-order
		if ok {
			for i := range matBuffer[work] {
				val := float64(matBuffer[work][i])
				if val < threshold {
					val = 0
				}
				dest.SetSym(work, i, val)
			}
			wg.Done()
		} else {
			break
		}
	}

	return
}

// DoThresholding does Thresholding
func DoThresholding(matBuffer [][13362]float32, dest *mat.SymDense, threshold float64) {
	order := make(chan int, runtime.NumCPU())

	var wg sync.WaitGroup

	for i := 0; i < runtime.NumCPU(); i++ {
		go thresholding(matBuffer, dest, threshold, order, &wg)
	}

	wg.Add(13362)
	for i := 0; i < 13362; i++ {
		order <- i
	}
	wg.Wait()
	close(order)
	return
}
