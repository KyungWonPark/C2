package main

import (
	"sync"

	"github.com/KyungWonPark/C2/internal/calc"
	"github.com/KyungWonPark/nifti"
)

// Voxel represents fMRI voxel coordinates
type Voxel struct {
	x int
	y int
	z int
}

var fileList []string
var convKernel [3][3][3]float32
var greyVoxels [13362]Voxel

func convolution(img *nifti.Nifti1Image, timePoint int, seed Voxel) float32 {
	var value float32
	value = 0

	for k := -1; k < 2; k++ {
		for j := -1; j < 2; j++ {
			for i := -1; i < 2; i++ {
				value += img.GetAt(uint32(seed.x+i), uint32(seed.y+j), uint32(seed.z+k), uint32(timePoint)) * convKernel[k+1][j+1][i+1]
			}
		}
	}

	return value / 8
}

func sampling(img *nifti.Nifti1Image, order <-chan int, wg *sync.WaitGroup, timeSeries [][600]float32) {
	for {
		timePoint, ok := <-order
		if ok {
			for i, vox := range greyVoxels {
				seed := Voxel{1 + 2*vox.x, 1 + 2*vox.y, 2 + 2*vox.z}
				timeSeries[i][timePoint-300] = convolution(img, timePoint, seed)
			}
			wg.Done()
		} else {
			break
		}
	}

	return
}

func doSampling(path string, timeSeries [][600]float32, workerConfig *calc.Config) {
	var img nifti.Nifti1Image
	img.LoadImage(path, true)

	order := make(chan int, workerConfig.NumLoader)
	var wg sync.WaitGroup
	for i := 0; i < workerConfig.NumLoader; i++ {
		go sampling(&img, order, &wg, timeSeries)
	}

	wg.Add(600)
	for timePoint := 300; timePoint < 900; timePoint++ {
		order <- timePoint
	}
	wg.Wait()
	close(order)
	return
}
