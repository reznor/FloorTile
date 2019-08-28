package main

import "fmt"
import "math/rand"
import "time"
import "github.com/fatih/color"

func init() {
	rand.Seed(time.Now().UnixNano())
}

type tile int

const (
	Z tile = iota // Unlaid
	A             // 12x12
	B             // 12x24 horizontal
	C             // 12x24 vertical
	D             // 24x24
)

const (
	ROWS    = 15
	COLUMNS = 60
)

var floor [ROWS][COLUMNS]tile
var tileCounts map[tile]int32

func isTileLaid(row, column int) bool {
	return floor[row][column] != Z
}

func isWithinFloor(row, column int) bool {
	return row >= 0 && row < ROWS && column >= 0 && column < COLUMNS
}

func canLayTile(row, column int) bool {
	return isWithinFloor(row, column) && !isTileLaid(row, column)
}

func getCandidateTiles(row, column int) []tile {
	var candidates []tile

	if canLayTile(row, column) {
		candidates = append(candidates, A)
	}

	if canLayTile(row, column) && canLayTile(row, column+1) {
		candidates = append(candidates, B)
	}

	if canLayTile(row, column) && canLayTile(row+1, column) {
		candidates = append(candidates, C)
	}

	if canLayTile(row, column) && canLayTile(row+1, column) && canLayTile(row, column+1) && canLayTile(row+1, column+1) {
		candidates = append(candidates, D)
	}

	return candidates
}

func getLeftTile(row, column int) tile {
	if !isWithinFloor(row, column-1) {
		return Z
	}

	return floor[row][column-1]
}

func getAboveTile(row, column int) tile {
	if !isWithinFloor(row-1, column) {
		return Z
	}

	return floor[row-1][column]
}

func getProblematicTiles(candidateTiles []tile, row, column int) []tile {
	var problematicTiles []tile

	for _, candidateTile := range candidateTiles {
		// Try to avoid picking:
		//
		//	1) the same tile as what's immediately to the left.
		// 	2) The same tile as what's immediately above.
		if candidateTile == getLeftTile(row, column) || candidateTile == getAboveTile(row, column) {
			problematicTiles = append(problematicTiles, candidateTile)
		}
	}

	return problematicTiles
}

func removeProblematicTiles(candidateTiles, problematicTiles []tile) []tile {
	var filteredTiles []tile

NextCandidateTile:
	for _, candidateTile := range candidateTiles {
		for _, problematicTile := range problematicTiles {
			if problematicTile == candidateTile {
				continue NextCandidateTile
			}
		}

		filteredTiles = append(filteredTiles, candidateTile)
	}

	return filteredTiles
}

func layTile(row, column int, t tile) {
	defer func() {
		if tileCounts == nil {
			tileCounts = make(map[tile]int32)
		}

		tileCounts[t]++
	}()

	switch t {
	case A:
		floor[row][column] = t
		return
	case B:
		floor[row][column] = t
		floor[row][column+1] = t
		return
	case C:
		floor[row][column] = t
		floor[row+1][column] = t
		return
	case D:
		floor[row][column] = t
		floor[row+1][column] = t
		floor[row][column+1] = t
		floor[row+1][column+1] = t
		return
	}
}

func makePattern() {
	for currentRow := 0; currentRow < ROWS; currentRow++ {
		for currentColumn := 0; currentColumn < COLUMNS; currentColumn++ {
			candidateTiles := getCandidateTiles(currentRow, currentColumn)
			if len(candidateTiles) == 0 {
				continue
			}

			problematicTiles := getProblematicTiles(candidateTiles, currentRow, currentColumn)

			// If at least one of the candidates is not problematic, filter out the problematic tiles; otherwise, we have no option but to pick a problematic one.
			if len(candidateTiles) > len(problematicTiles) {
				candidateTiles = removeProblematicTiles(candidateTiles, problematicTiles)
			}

			layTile(currentRow, currentColumn, candidateTiles[rand.Intn(len(candidateTiles))])
		}
	}
}

func getTileColor(t tile) *color.Color {
	switch t {
	case A:
		return color.New(color.FgRed)
	case B:
		return color.New(color.FgBlue)
	case C:
		return color.New(color.FgGreen)
	case D:
		return color.New(color.FgYellow)
	case Z:
		return color.New(color.FgWhite)
	}

	return nil
}

func printTile(t tile) {
	getTileColor(t).Add(color.Bold).Printf("%d", t)
}

func printFloor() {
	for t,n := range tileCounts {
		fmt.Printf("%v:%d\n", t, n)
	}

	for i := 0; i < ROWS; i++ {
		for j := 0; j < COLUMNS; j++ {
			printTile(floor[i][j])
		}
		fmt.Printf("\n")
	}
}

func main() {
	makePattern()
	printFloor()
}
