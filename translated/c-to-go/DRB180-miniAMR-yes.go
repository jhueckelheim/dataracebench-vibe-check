/*
Simplified miniAMR proxy app to reproduce data race behavior.
Data Race Pair, in@60:43:R vs. in@52:43:W
                work@65:19@W vs. work@65:19@W
                bp->array@65:36@R vs. bp->array@75:19@W
                bp->array@66:36@R vs. bp->array@75:19@W
                bp->array@67:36@R vs. bp->array@75:19@W
                bp->array@68:36@R vs. bp->array@75:19@W
                bp->array@69:36@R vs. bp->array@75:19@W
                bp->array@70:36@R vs. bp->array@75:19@W
                bp->array@71:36@R vs. bp->array@75:19@W
*/

package main

import (
	"math/rand"
	"sync"
)

type numSz int64

var (
	maxNumBlocks int
	numRefine    int
	numVars      int
	xBlockSize   int
	yBlockSize   int
	zBlockSize   int
	errorTol     int
	tol          float64
)

type block struct {
	number  numSz
	level   int
	refine  int
	newProc int
	parent  numSz // if original block -1, else if on node, number in structure, else (-2 - parent->number)
	array   [][][][]float64
}

var blocks []block

func stencilCalc(variable int, stencilIn int) {
	// Shared work array that creates race condition
	// In the original C code, this was declared as private but used unsafely
	work := make([][][]float64, xBlockSize+2)
	for i := range work {
		work[i] = make([][]float64, yBlockSize+2)
		for j := range work[i] {
			work[i][j] = make([]float64, zBlockSize+2)
		}
	}

	var wg sync.WaitGroup

	// Parallel region - race condition occurs here
	for in := 0; in < maxNumBlocks; in++ {
		wg.Add(1)
		go func(blockIdx int) {
			defer wg.Done()

			bp := &blocks[blockIdx]

			// First phase: compute work array
			// Race condition: multiple goroutines access shared work array
			for i := 1; i <= xBlockSize; i++ {
				for j := 1; j <= yBlockSize; j++ {
					for k := 1; k <= zBlockSize; k++ {
						// Race: work array is shared among all goroutines
						work[i][j][k] = (bp.array[variable][i-1][j][k] +
							bp.array[variable][i][j-1][k] +
							bp.array[variable][i][j][k-1] +
							bp.array[variable][i][j][k] +
							bp.array[variable][i][j][k+1] +
							bp.array[variable][i][j+1][k] +
							bp.array[variable][i+1][j][k]) / 7.0
					}
				}
			}

			// Second phase: copy back to array
			// Race condition: reading from shared work array while others might be writing
			for i := 1; i <= xBlockSize; i++ {
				for j := 1; j <= yBlockSize; j++ {
					for k := 1; k <= zBlockSize; k++ {
						// Race: bp.array write vs bp.array read from other goroutines
						bp.array[variable][i][j][k] = work[i][j][k]
					}
				}
			}
		}(in)
	}

	wg.Wait()
}

func allocate() {
	blocks = make([]block, maxNumBlocks)

	for n := 0; n < maxNumBlocks; n++ {
		blocks[n].number = -1
		blocks[n].array = make([][][][]float64, numVars)
		for m := 0; m < numVars; m++ {
			blocks[n].array[m] = make([][][]float64, xBlockSize+2)
			for i := 0; i < xBlockSize+2; i++ {
				blocks[n].array[m][i] = make([][]float64, yBlockSize+2)
				for j := 0; j < yBlockSize+2; j++ {
					blocks[n].array[m][i][j] = make([]float64, zBlockSize+2)
				}
			}
		}
	}
}

func initialize() {
	// Initialize blocks
	for o := 0; o < maxNumBlocks && o < 1; o++ {
		bp := &blocks[o]
		bp.level = 0
		bp.number = numSz(o)

		for variable := 0; variable < numVars; variable++ {
			for ib := 1; ib <= xBlockSize; ib++ {
				for jb := 1; jb <= yBlockSize; jb++ {
					for kb := 1; kb <= zBlockSize; kb++ {
						bp.array[variable][ib][jb][kb] = rand.Float64()
					}
				}
			}
		}
	}
}

func driver() {
	initialize()

	for variable := 0; variable < numVars; variable++ {
		stencilCalc(variable, 7)
	}
}

func main() {
	maxNumBlocks = 500
	numRefine = 5
	numVars = 40
	xBlockSize = 10
	yBlockSize = 10
	zBlockSize = 10

	allocate()
	driver()
}
