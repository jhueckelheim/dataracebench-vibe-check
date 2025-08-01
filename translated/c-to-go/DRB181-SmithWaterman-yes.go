/*
Smith-Waterman sequence alignment algorithm with data race
Data Race Pair, *maxPos@179:9:W vs. *maxPos@177:17:R
                H[index]@173:5:W vs. H[*maxPos]@177:15:W
*/

package main

import (
	"fmt"
	"math"
	"math/rand"
	"sync"
	"time"
)

// Constants
const (
	PATH     = -1
	NONE     = 0
	UP       = 1
	LEFT     = 2
	DIAGONAL = 3
)

// Global Variables
var (
	m              int64 // Columns - Size of string a
	n              int64 // Lines - Size of string b
	matchScore     = 5
	missmatchScore = -3
	gapScore       = -4
	a, b           []byte
)

func min(x, y int64) int64 {
	if x < y {
		return x
	}
	return y
}

func max(x, y int64) int64 {
	if x > y {
		return x
	}
	return y
}

func generate() {
	rand.Seed(time.Now().UnixNano())

	// Generate string a
	for i := int64(0); i < m; i++ {
		aux := rand.Intn(4)
		switch aux {
		case 0:
			a[i] = 'A'
		case 2:
			a[i] = 'C'
		case 3:
			a[i] = 'G'
		default:
			a[i] = 'T'
		}
	}

	// Generate string b
	for i := int64(0); i < n; i++ {
		aux := rand.Intn(4)
		switch aux {
		case 0:
			b[i] = 'A'
		case 2:
			b[i] = 'C'
		case 3:
			b[i] = 'G'
		default:
			b[i] = 'T'
		}
	}
}

func nElement(i int64) int64 {
	if i < m && i < n {
		return i
	} else if i < max(m, n) {
		minVal := min(m, n)
		return minVal - 1
	} else {
		minVal := min(m, n)
		return 2*minVal - i + int64(math.Abs(float64(m-n))) - 2
	}
}

func matchMissmatchScore(i, j int64) int {
	if a[j-1] == b[i-1] {
		return matchScore
	}
	return missmatchScore
}

func similarityScore(i, j int64, H, P []int, maxPos *int64, mutex *sync.Mutex) {
	// Stores index of element
	index := m*i + j

	// Get element above
	up := H[index-m] + gapScore

	// Get element on the left
	left := H[index-1] + gapScore

	// Get element on the diagonal
	diag := H[index-m-1] + matchMissmatchScore(i, j)

	// Calculate the maximum
	maxVal := NONE
	pred := NONE

	if diag > maxVal {
		maxVal = diag
		pred = DIAGONAL
	}

	if up > maxVal {
		maxVal = up
		pred = UP
	}

	if left > maxVal {
		maxVal = left
		pred = LEFT
	}

	// Insert the value in the similarity and predecessor matrices
	H[index] = maxVal
	P[index] = pred

	// Race condition: Reading H[*maxPos] without proper synchronization
	// The critical section only protects the write to *maxPos but not the read of H[*maxPos]
	if maxVal > H[*maxPos] {
		mutex.Lock()
		// Race: *maxPos might be modified by another goroutine between the check and this assignment
		*maxPos = index
		mutex.Unlock()
	}
}

func calcFirstDiagElement(i int64) (si, sj int64) {
	if i < n {
		si = i
		sj = 1
	} else {
		si = n - 1
		sj = i - n + 2
	}
	return
}

func main() {
	m = 2048
	n = 2048

	fmt.Printf("\nMatrix[%d][%d]\n", n, m)

	// Allocate a and b
	a = make([]byte, m)
	b = make([]byte, n)

	// Generate random arrays a and b before incrementing m and n
	generate()

	// Because now we have zeros
	m++
	n++

	// Allocate similarity matrix H
	H := make([]int, m*n)

	// Allocate predecessor matrix P
	P := make([]int, m*n)

	// Start position for backtrack
	var maxPos int64 = 0

	// Calculate the similarity matrix
	nDiag := m + n - 3

	var mutex sync.Mutex
	var wg sync.WaitGroup

	for i := int64(1); i <= nDiag; i++ {
		nEle := nElement(i)
		si, sj := calcFirstDiagElement(i)

		for j := int64(1); j <= nEle; j++ {
			wg.Add(1)
			go func(jVal int64) {
				defer wg.Done()
				ai := si - jVal + 1
				aj := sj + jVal - 1
				similarityScore(ai, aj, H, P, &maxPos, &mutex)
			}(j)
		}
	}

	wg.Wait()
	fmt.Printf("Completed with maxPos: %d\n", maxPos)
}
