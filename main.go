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
	"gonum.org/v1/gonum/blas/blas64"
	"gonum.org/v1/gonum/mat"
	blas_netlib "gonum.org/v1/netlib/blas/netlib"
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
	// OpenBLAS
	blas64.Use(blas_netlib.Implementation{})

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

	fmt.Printf("Averaging: %d samples.\n", len(fileList))

	calc.DoAverage(matBuffer, float32(len(fileList)))

	fmt.Println("Finished Calculation.")

	backend := make([]float64, 13362*13362)
	afterThreshold := mat.NewSymDense(13362, backend)

	thsArr := []float64{0}

	for _, threshold := range thsArr {
		calc.DoThresholding(matBuffer, afterThreshold, threshold)

		var eigSym mat.EigenSym
		ok := eigSym.Factorize(afterThreshold, true)
		if !ok {
			fmt.Printf("Failed to do eigen decomposition for threshold: %f!\n", threshold)
		}

		eigVals := eigSym.Values(nil)
		var eigVecs mat.Dense
		eigSym.VectorsTo(&eigVecs)

		postFix := "-threshold-" + fmt.Sprintf("%f", threshold)
		writeEntireEigenValue(eigVals, "eigenValue"+postFix)
		writeEntireEigenVector(&eigVecs, "eigenVector"+postFix)

		fmt.Printf("Processed thr: %f\n", threshold)
	}

	return
}

func writeEntireEigenValue(eigVals []float64, fileName string) {
	path := os.Getenv("RESULT") + "/" + fileName + ".csv"
	f, err := os.Create(path)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	for i := range eigVals {
		var line string
		line += fmt.Sprintf("%.*e", 6, eigVals[i])
		fmt.Fprintf(f, "%s\n", line)
	}

	return
}

func writeEntireEigenVector(eigVecs *mat.Dense, fileName string) {
	path := os.Getenv("RESULT") + "/" + fileName + ".csv"
	f, err := os.Create(path)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	rows, cols := eigVecs.Dims()
	for i := 0; i < rows; i++ {
		var line string
		for j := 0; j < cols; j++ {
			line += fmt.Sprintf("%.*e", 6, eigVecs.At(i, j))
			if j < (cols - 1) {
				line += ", "
			}
		}
		fmt.Fprintf(f, "%s\n", line)
	}

	return
}
