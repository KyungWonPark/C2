package util

import (
	"fmt"
	"log"
	"os"
	"strings"
)

// MatWrite writes a matrix to text file
func MatWrite(mat [][13362]float32, fileName string) {
	path := os.Getenv("RESULT") + fileName + ".txt"
	f, err := os.Create(path)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	for i := range mat {
		var line string
		for j := range mat[i] {
			line += fmt.Sprintf("%.*e, ", 6, mat[i][j])
		}
		line = strings.TrimSuffix(line, ",")
		fmt.Fprintf(f, "%s\n", line)
	}
	return
}
