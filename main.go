package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"runtime"
	"strconv"
	"strings"

	"github.com/KyungWonPark/C2/internal/calc"
)

func init() {
	// Sampling setting
	for z := 0; z < 2; z++ {
		for y := 0; y < 2; y++ {
			for x := 0; x < 2; x++ {
				taxiDist := math.Abs(float64(x-1)) + math.Abs(float64(y-1)) + math.Abs(float64(z-1))
				convKernel[z][y][x] = float32(math.Pow(2, -1*taxiDist))
			}
		}
	}

	f, err := os.Open("files/greyList.txt")
	if err != nil {
		log.Fatal("Failed to open greyList.txt file!", err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	i := 0
	for scanner.Scan() {
		line := scanner.Text()
		xyz := strings.Split(line, ",")
		x, err0 := strconv.Atoi(xyz[0])
		y, err1 := strconv.Atoi(xyz[1])
		z, err2 := strconv.Atoi(xyz[2])
		if err0 != nil || err1 != nil || err2 != nil {
			log.Fatal("Failed to convert ascii to integer!", err0)
		}

		greyVoxels[i] = Voxel{x, y, z}
		i++
	}

	g, err := os.Open("files/fileList.txt")
	if err != nil {
		log.Fatal("Failed to open fileList.txt!", err)
	}
	defer g.Close()

	scanner = bufio.NewScanner(g)
	for scanner.Scan() {
		line := scanner.Text()
		fileList = append(fileList, line)
	}

	return
}

func main() {
	// Worker number setting
	var workerConfig calc.Config

	workerConfig.NumLoader = 4
	workerConfig.NumComputer = runtime.NumCPU() - workerConfig.NumLoader + 2
	workerConfig.NumBufferSize = 8

	// Prepare matBuffer
	matBuffer := make([][13362]float32, 13362)
	buffer := make(ringBuffer, workerConfig.NumBufferSize)
	for i := range buffer {
		buffer[i].isEmpty = true
		buffer[i].data = make([][600]float32, 13362)
	}
	bufferCh := make(chan int, workerConfig.NumBufferSize)

	go load(fileList, buffer, bufferCh, &workerConfig)

	compute(buffer, bufferCh, matBuffer, &workerConfig)

	calc.DoAverage(matBuffer, float32(len(fileList)))

	fmt.Println("Finished Calculation.")

	f, err := os.Create("output-matrix.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	for i := range matBuffer {
		for j := range matBuffer[i] {

			fmt.Fprintf(f, "%.*e", 6, matBuffer[i][j])
			if j != len(matBuffer[i])-1 {
				fmt.Fprintf(f, "%s", ",")
			}
		}
		fmt.Fprintf(f, "%s", "\n")
	}

	return
}
