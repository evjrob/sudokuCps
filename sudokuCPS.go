/* ****************************************************************************
A Constraint Propagation Search algorithm to solve general sudoku puzzles of
any dimension based on the work by Peter Norvig: http://norvig.com/sudoku.html
* ****************************************************************************/

package sudokuCPS

import (
	"strconv"
)

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

func convertIntSliceToMap(originalPuzzle [][]int) (convertedPuzzle map[string]string) {
	convertedPuzzle = make(map[string]string)

	for i := range originalPuzzle {
		row := numberToAlpha(i + 1)
		for j := range originalPuzzle[i] {
			column := strconv.Itoa(j + 1)
			position := row + column
			element := strconv.Itoa(originalPuzzle[i][j])
			convertedPuzzle[position] = element
		}
	}
	return convertedPuzzle
}

func convertMapToIntSlice(originalPuzzle map[string]string, puzzleDim int) (convertedPuzzle [][]int) {
	convertedPuzzle = make([][]int, puzzleDim, puzzleDim)
	for i := 0; i < puzzleDim; i++ {
		convertedPuzzle[i] = make([]int, puzzleDim, puzzleDim)
		row := numberToAlpha(i + 1)
		for j := 0; j < puzzleDim; j++ {
			column := strconv.Itoa(j + 1)
			position := row + column
			element, _ := strconv.Atoi(originalPuzzle[position])
			convertedPuzzle[i][j] = element
		}
	}
	return convertedPuzzle
}

// Convert a solved puzzle from the values representation to the solved representation
func convertValuesToSolved(originalPuzzle map[string][]string) (convertedPuzzle map[string]string) {
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

// Integer exponentiation. return a^b.
func pow(a, b int) int {
	p := 1
	for b > 0 {
		if b&1 != 0 {
			p *= a
		}
		b >>= 1
		a *= a
	}
	return p
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

// Convert a base 10 number to a corresponding alpha character where A = 1, ..., Z = 26.
func alphaToNumber(alphaString string) (number int) {

	alphabet := map[string]int{
	"A":1, "B":2, "C":3, "D":4, "E":5, "F":6, "G":7, "H":8, "I":9, "J":10, "K":11,
	"L":12, "M":13, "N":14, "O":15, "P":16, "Q":17, "R":18, "S":19, "T":20, "U":21,
	"V":22, "W":23, "X":24, "Y":25, "Z":26,
	}

	number = 0
	base := 26
	characters := []rune(alphaString)
	characterCount := len(characters)

	for i, c := range characters {
		number = number + alphabet[string(c)]*pow(base, characterCount-1-i)
	}

	return number
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

		subRow := rows[rs*blockYDim : (rs+1)*blockYDim]

		for cs := 0; cs < blockYDim; cs++ {

			subColumn := columns[cs*blockXDim : (cs+1)*blockXDim]

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

// Solve uses constraint propagation on the original puzzle via parseGrid(), then backtack/search through the remaining possibilities
func Solve(originalPuzzle [][]int, blockXDim, blockYDim int) ([][]int, bool) {

	puzzle := convertIntSliceToMap(originalPuzzle)

	digits := makeDigits(blockXDim * blockYDim)
	rows := makeRows(blockXDim * blockYDim)
	columns := digits
	squares := cross(rows, columns)
	unitList := makeUnitList(rows, columns, blockXDim, blockYDim)
	units := makeUnits(squares, unitList)
	peers := makePeers(units, blockXDim, blockYDim)

	values, contradiction := parseGrid(puzzle, digits, peers, units)

	if contradiction {
		return nil, false
	}

	searchedPuzzle, contradiction := search(values, digits, peers, units)

	if contradiction {
		return nil, false
	}

	solvedPuzzleMap := convertValuesToSolved(searchedPuzzle)
	solvedPuzzle := convertMapToIntSlice(solvedPuzzleMap, blockXDim * blockYDim)

	return solvedPuzzle, true
}
