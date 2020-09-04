package main

import (
	"os"

	"github.com/KyungWonPark/C2/internal/calc"
)

type ringBuffer []struct {
	isEmpty bool
	data    [][600]float32
}

func load(fileList []string, buffer ringBuffer, bufferCh chan<- int, workerConfig *calc.Config) {
	bufferIndex := 0
	dataDir := os.Getenv("DATA") + "fMRI-Smoothed/"

	for i := 0; i < len(fileList); i++ {
		if !buffer[bufferIndex].isEmpty {
			i = i - 1
			bufferIndex = (bufferIndex + 1) % workerConfig.NumBufferSize
			continue
		}

		path := dataDir + fileList[i]
		doSampling(path, buffer[bufferIndex].data, workerConfig)
		buffer[bufferIndex].isEmpty = false
		bufferCh <- bufferIndex
		bufferIndex = (bufferIndex + 1) % workerConfig.NumBufferSize
	}
	close(bufferCh)

	return
}

func compute(buffer ringBuffer, bufferCh <-chan int, matBuffer [][13362]float32, workerConfig *calc.Config) {
	stats := make([]calc.LinStatEle, 13362)

	for {
		bufferIndex, ok := <-bufferCh
		if ok {
			timeSeries := buffer[bufferIndex].data
			// z-score
			calc.DoZScoring(timeSeries, workerConfig)

			// sigmoid
			calc.DoSigmoid(timeSeries, stats, workerConfig)

			// pearson & accumulation
			calc.DoPearson(timeSeries, stats, matBuffer)

			buffer[bufferIndex].isEmpty = true
		} else {
			break
		}
	}
	return
}
