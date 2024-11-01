package main

import "fmt"

func calculateGrid(allPositions *[][][3]float64, ww int32, hh int32) {
	w := float64(ww)
	h := float64(hh)
	var enemySize int32 = 15
	flotES := float64(enemySize)
	var spacing int32 = 10
	var cols int32 = ww / (enemySize + spacing)
	var rows int32 = hh / (enemySize + spacing)
	paddingX := (w - (float64(cols * enemySize))) / (float64(cols) + 1) // Spaziatura orizzontale
	paddingY := (h - (float64(rows * enemySize))) / (float64(rows) + 1) // Spaziatura verticale

	//var allPositions [][][2]float64 //[[row]]
	for r := 1; r <= int(rows); r++ {
		var tempRow [][3]float64
		for c := 1; c <= int(cols); c++ {
			var tempCoord [3]float64
			tempCoord[0] = float64(c) * (flotES + paddingX)
			tempCoord[1] = float64(r) * (flotES + paddingY)
			tempCoord[2] = float64(0) ///0 dead - 1 alive
			tempRow = append(tempRow, tempCoord)
		}
		*allPositions = append(*allPositions, tempRow)
	}

}
func handleGameOfLife() {
	//TODO later as param
	//speed := 5                      // secondi/10 between cycles = 0.5s
	//density := 20                   // %
	var allPositions [][][3]float64 //[[row]]

	calculateGrid(&allPositions, 2150, 950)

	fmt.Println("asd??", allPositions)

}

func main() {
	handleGameOfLife()
}
