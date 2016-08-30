/* ****************************************************************************
A Constraint Propagation Search algorithm to solve general sudoku puzzles of
any dimension based on the work by Peter Norvig: http://norvig.com/sudoku.html

Copyright (c) 2016 Everett Robinson

Permission is hereby granted, free of charge, to any person obtaining a copy of
this software and associated documentation files (the "Software"), to deal in
Software without restriction, including without limitation the rights to use,
copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the
Software, and to permit persons to whom the Software is furnished to do so,
subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS
FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER
IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
* ****************************************************************************/

package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"
)

// Modified from https://stackoverflow.com/questions/9862443/golang-is-there-a-better-way-read-a-file-of-integers-into-an-array
// Read in the start state of the sudoku puzzle (of arbitrary dimension) in a single line presentation.
func readInOneLine(r io.Reader, line int, delimiter string, emptyValue string, blockXDim int, blockYDim int) (puzzle map[string]string, e error) {

	scanner := bufio.NewScanner(r)
	scanner.Split(bufio.ScanLines)

	// Start puzzle at line 1 (more user friendly)
	lineCounter := 1

	// For each line
	for scanner.Scan() {

		// Check if it's the line we selected
		if line == lineCounter {

			// Read the puzzle text in and split it into it's components
			puzzleText := scanner.Text()
			puzzleElements := strings.Split(puzzleText, delimiter)
			puzzleDim := blockXDim * blockYDim
			puzzle = make(map[string]string)

			for i := 0; i < puzzleDim; i++ {
				row := numberToAlpha(i + 1)
				for j := 0; j < puzzleDim; j++ {
					column := strconv.Itoa(j + 1)
					position := row + column
					element := puzzleElements[(i*puzzleDim)+j]
					if element != emptyValue {
						puzzle[position] = element
					} else {
						puzzle[position] = "0"
					}
				}
			}
		}

		lineCounter++
	}

	return puzzle, scanner.Err()
}

// return the number of digits in an int up to 4. Sudoku puzzles of greater than
// 9999*9999 are probably not practical.
func numDigits(n int) int {
	if n < 0 {
		n = -1 * n
	}
	if n < 10 {
		return 1
	}
	if n < 100 {
		return 2
	}
	if n < 1000 {
		return 3
	}

	return 4
}

// Make the puzzle look like a sudoku
func printPuzzle(puzzle map[string]string, blockXDim int, blockYDim int) {

	puzzleDim := blockXDim * blockYDim

	width := numDigits(blockXDim * blockYDim)

	for r := 0; r < puzzleDim; r++ {
		if r > 0 && r%blockYDim == 0 {
			for i := 0; i < (blockXDim-1)+(width+2)*blockXDim*blockYDim; i++ {
				fmt.Printf("%s", "-")
			}
			fmt.Println()
		}
		for c := 0; c < puzzleDim; c++ {
			if c > 0 && c%blockXDim == 0 {
				fmt.Printf("|")
			}
			row := numberToAlpha(r + 1)
			column := strconv.Itoa(c + 1)
			position := row + column
			if puzzle[position] != "0" {
				fmt.Printf("%-*s%s ", width-len(puzzle[position])+1, " ", puzzle[position])
			} else {
				fmt.Printf("%-*s ", width+1, " ")
			}
		}
		fmt.Printf("\n")
	}

}

// Convert a solved puzzle from the values representation to the solved representation
func convertPuzzle(originalPuzzle map[string][]string) (convertedPuzzle map[string]string) {
	convertedPuzzle = make(map[string]string)

	// copy each map element
	for key, value := range originalPuzzle {
		convertedPuzzle[key] = value[0]
	}

	return convertedPuzzle
}

// Copy the map representation of the initial or solved puzzles
func copyPuzzle(originalPuzzle map[string]string) (copiedPuzzle map[string]string) {

	copiedPuzzle = make(map[string]string)

	// copy each map element
	for key, value := range originalPuzzle {
		copiedPuzzle[key] = value
	}

	return copiedPuzzle
}

// Copy the values map representation of the unsolved puzzle.
func copyValues(values map[string][]string) (copiedValues map[string][]string) {

	copiedValues = make(map[string][]string, len(values))

	// copy each map element
	for key := range values {
		copiedValues[key] = make([]string, len(values[key]))
		copy(copiedValues[key], values[key])
	}
	return copiedValues
}

// Return the integer quotient and remainder from long division
func divmod(number int, base int) (quotient int, remainder int) {
	quotient = number / base
	remainder = number % base

	return quotient, remainder
}

// Convert a base 10 number to a corresponding alpha character where 1 = A, ..., 26 = Z.
func numberToAlpha(number int) (alpha string) {

	alphabet := []string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N",
		"O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z"}
	alphaString := ""

	base := 26

	var remainder int

	for number > 0 {

		number, remainder = divmod(number-1, base)
		alphaString = alphabet[remainder] + alphaString
	}

	return alphaString
}

// Create a slice for the column headers and digits in the puzzle
func makeDigits(puzzleDim int) []string {

	digits := make([]string, puzzleDim)

	for i := 0; i < puzzleDim; i++ {
		digits[i] = strconv.Itoa(i + 1)
	}

	return digits
}

// Create a slice of the alphabetic row labels
func makeRows(puzzleDim int) []string {

	rows := make([]string, puzzleDim)

	for i := 0; i < puzzleDim; i++ {
		rows[i] = numberToAlpha(i + 1)
	}

	return rows
}

// Create a map of each square to the units in which it is a member
func makeUnits(squares []string, unitList [][]string) map[string][][]string {

	units := make(map[string][][]string)
	unitCount := make(map[string]int)

	for i := 0; i < len(unitList); i++ {
		for j := 0; j < len(unitList[i]); j++ {

			square := unitList[i][j]

			_, ok := units[square]

			if !ok {
				units[square] = make([][]string, 3)
				units[square][0] = unitList[i]
				unitCount[square] = 1
			} else {
				units[square][unitCount[square]] = unitList[i]
				unitCount[square]++
			}
		}
	}
	return units
}

// Creae a 2D slice contaning every unit in the puzzle and the corresponding squares
// unitlist = ([cross(rows, c) for c in cols] + [cross(r, cols) for r in rows] +
//             [cross(rs, cs) for rs in ('ABC','DEF','GHI') for cs in ('123','456','789')])
func makeUnitList(rows []string, columns []string, blockXDim int, blockYDim int) (unitList [][]string) {

	unitCount := 3 * blockXDim * blockYDim

	unitList = make([][]string, unitCount)

	i := 0

	// [cross(rows, c) for c in cols]
	for c := range columns {
		unitList[i] = cross(rows, columns[c:c+1])
		i++
	}

	// cross(r, cols) for r in rows]
	for r := range rows {
		unitList[i] = cross(rows[r:r+1], columns)
		i++
	}

	// Dimensions are swapped because the number of blocks in each dimension is equal to the other dimension's
	// blockDim. For rows; block count = puzzleDim/blockYDim; (blockXDim*blockYDim)/blockYDim = blockXDim.
	// The generalized form of: [cross(rs, cs) for rs in ('ABC','DEF','GHI') for cs in ('123','456','789')]
	for rs := 0; rs < blockXDim; rs++ {

		subRow := rows[rs*blockXDim : (rs+1)*blockXDim]

		for cs := 0; cs < blockYDim; cs++ {

			subColumn := columns[cs*blockYDim : (cs+1)*blockYDim]

			unitList[i] = cross(subRow, subColumn)

			i++
		}
	}

	return unitList
}

// Create the peers map for each square of the puzzle
func makePeers(units map[string][][]string, blockXDim int, blockYDim int) map[string]map[string]bool {

	puzzleDim := blockXDim * blockYDim

	// The number of peers each square has
	peerCount := 3*(puzzleDim-1) - (blockXDim - 1) - (blockYDim - 1)

	peers := make(map[string]map[string]bool)

	for square, allUnits := range units {
		subPeers := make(map[string]bool, peerCount)
		for _, unit := range allUnits {
			for _, peer := range unit {
				if square != peer {
					subPeers[peer] = true
				}
			}
		}
		peers[square] = subPeers
	}

	return peers
}

// Create the combination of each element of A with each element of B
func cross(A []string, B []string) []string {

	aDim := len(B)
	bDim := len(A)

	crossed := make([]string, aDim*bDim)

	for i, a := range A {
		for j, b := range B {
			crossed[(i*aDim)+j] = a + b
		}
	}

	return crossed
}

// Borrowed from: http://stackoverflow.com/questions/10485743/contains-method-for-a-slice
func contains(s []string, e string) (int, bool) {
	for i, a := range s {
		if a == e {
			return i, true
		}
	}
	return -1, false
}

// Eliminate all the other values (except d) from values[s] and propagate.
// Return values, except return False if a contradiction is detected.
func assign(inValues map[string][]string, s string, d string, peers map[string]map[string]bool, units map[string][][]string) (values map[string][]string, contradiction bool) {

	otherValues := copyValues(inValues)
	digitIndex, _ := contains(otherValues[s], d)
	otherValues[s] = append(otherValues[s][:digitIndex], otherValues[s][digitIndex+1:]...)

	for i := 0; i < len(otherValues[s]); i++ {
		inValues, contradiction = eliminate(inValues, s, otherValues[s][i], peers, units)
		if contradiction {
			return nil, contradiction
		}
	}

	return inValues, false
}

// Eliminate d from values[s]; propagate when values or places <= 2.
// Return values, except return False if a contradiction is detected.
func eliminate(values map[string][]string, s string, d string, peers map[string]map[string]bool, units map[string][][]string) (outValues map[string][]string, contradiction bool) {

	digitIndex, digitFound := contains(values[s], d)

	if !digitFound {
		return values, false
	}

	outValues = copyValues(values)
	outValues[s] = append(outValues[s][:digitIndex], outValues[s][digitIndex+1:]...)

	// (1) If a square s is reduced to one value d2, then eliminate d2 from the peers.
	if len(outValues[s]) == 0 {
		return nil, true //Contradiction: removed last value
	} else if len(outValues[s]) == 1 {
		d2 := outValues[s][0]
		for s2 := range peers[s] {
			outValues, contradiction = eliminate(outValues, s2, d2, peers, units)
			if contradiction {
				return nil, contradiction
			}
		}
	}
	// (2) If a unit u is reduced to only one place for a value d, then put it there.
	for _, unit := range units[s] {
		dplaces := []string{}
		for _, s3 := range unit {
			_, digitFound := contains(outValues[s3], d)
			if digitFound {
				dplaces = append(dplaces, s3)
			}
		}
		if len(dplaces) == 0 {
			return nil, true // Contradiction: no place for this value
		} else if len(dplaces) == 1 {
			// d can only be in one place in unit; assign it there
			outValues, contradiction = assign(outValues, dplaces[0], d, peers, units)
			if contradiction {
				return nil, contradiction
			}
		}
	}
	return outValues, false
}

// Convert puzzle to a dict of possible values, {square: digits}, or
// return False if a contradiction is detected.
// To start, every square can be any digit; then assign values from the grid.
func parseGrid(puzzle map[string]string, digits []string, peers map[string]map[string]bool, units map[string][][]string) (values map[string][]string, contradiction bool) {

	values = make(map[string][]string)

	for square := range puzzle {
		values[square] = digits
	}

	for square, digit := range puzzle {
		_, digitFound := contains(digits, digit)
		if digitFound {
			values, contradiction = assign(values, square, digit, peers, units)
			if contradiction {
				return nil, contradiction
			}
		}
	}

	return values, false
}

// Use backtracking to search through the remaining possibilities after the constraints have been propagated.
func search(values map[string][]string, digits []string, peers map[string]map[string]bool, units map[string][][]string) (map[string][]string, bool) {

	// Check if the puzzle has been solved already
	solved := true

	for s := range values {
		if len(values[s]) != 1 {
			solved = false
		}
	}

	if solved {
		return values, false
	}

	// Figure out which unsolved square has minimum number of possibilities
	minLength := len(digits) + 1
	minSquare := ""

	for s, sValues := range values {
		length := len(sValues)
		if length > 1 && length < minLength {
			minSquare = s
			minLength = length
		}
	}

	// Assign one of the possibilities to that square and recursively search
	for _, searchDigit := range values[minSquare] {
		searchValues := copyValues(values)
		searchValues, contradiction := assign(searchValues, minSquare, searchDigit, peers, units)
		if !contradiction {
			resultValues, contradiction := search(searchValues, digits, peers, units)
			if !contradiction {
				// We found a solution: return it!
				return resultValues, false
			}
		}
	}

	// No solution found after searching: contradiction must exist
	return nil, true
}

// Use constraint propagation on the original puzzle via parseGrid(), then backtack/search through the remaining possibilities
func solve(puzzle map[string]string, digits []string, peers map[string]map[string]bool, units map[string][][]string) (map[string]string, bool) {

	values, contradiction := parseGrid(puzzle, digits, peers, units)

	if contradiction {
		return nil, contradiction
	}

	searchedPuzzle, contradiction := search(values, digits, peers, units)

	if contradiction {
		return nil, contradiction
	}

	solvedPuzzle := convertPuzzle(searchedPuzzle)

	return solvedPuzzle, false
}

func main() {
	start := time.Now()

	inputModePtr := flag.String("m", "one-line", "An input mode used to interpret the input file")
	delimiterPtr := flag.String("del", "", "The delimeter used to separate the puzzle squares in the input")
	emptyValuePtr := flag.String("e", ".", "The character used to indicate an empty square in the puzzle")
	dimPtr := flag.String("d", "3x3", "The dimensions of one of the puzzle blocks (eg. standard sudoku is 3x3)")
	filePtr := flag.String("f", "puzzles.txt", "The filename to be checked")
	linePtr := flag.String("l", "1", "The line of the puzzle to be solved")

	flag.Parse()

	puzzleLine, _ := strconv.Atoi(*linePtr)
	puzzleDim := strings.Split(*dimPtr, "x")
	blockXDim, _ := strconv.Atoi(puzzleDim[0])
	blockYDim, _ := strconv.Atoi(puzzleDim[1])

	inFile, err := os.Open(*filePtr)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Use string arrays to avoid ambiguity when puzzle dimesnions exceed the 9
	// column digits in a standard sudoku puzzle
	var digits []string
	var rows []string
	var columns []string
	var squares []string
	var unitList [][]string
	var units map[string][][]string
	var peers map[string]map[string]bool

	var originalPuzzle map[string]string

	if *inputModePtr == "one-line" {
		// Read the file into an array
		originalPuzzle, err = readInOneLine(inFile, puzzleLine, *delimiterPtr, *emptyValuePtr, blockXDim, blockYDim)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	} else {
		fmt.Println("No appropriate input mode for the puzzle was entered.")
		os.Exit(1)
	}

	digits = makeDigits(blockXDim * blockYDim)
	rows = makeRows(blockXDim * blockYDim)
	columns = digits
	squares = cross(rows, columns)
	unitList = makeUnitList(rows, columns, blockXDim, blockYDim)
	units = makeUnits(squares, unitList)
	peers = makePeers(units, blockXDim, blockYDim)

	fmt.Println()
	fmt.Println("Original Puzzle:")
	printPuzzle(originalPuzzle, blockXDim, blockYDim)

	solvedPuzzle, contradiction := solve(originalPuzzle, digits, peers, units)

	if !contradiction {
		fmt.Println()
		fmt.Println("Solved Puzzle:")
		printPuzzle(solvedPuzzle, blockXDim, blockYDim)
	} else {
		fmt.Println()
		fmt.Println("No viable solution to the puzzle was found.\n")
		fmt.Println()
	}

	elapsed := time.Since(start)

	fmt.Printf("Execution completed in %s \n", elapsed)
}
