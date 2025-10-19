package game

import (
	"bytes"
	"fmt"
	"strconv"
)

// Grid parses the server response and returns a grid of integers
func Grid(response []byte) ([]int, error) {
	payload := response[4 : len(response)-1]

	rawValues := bytes.Split(payload, []byte(","))

	grid := make([]int, len(rawValues))
	for i, value := range rawValues {
		trimmed := bytes.Trim(value, "\"")

		num, err := strconv.Atoi(string(trimmed))
		if err != nil {
			return nil, fmt.Errorf("error converting value to int: %v", err)
		}
		grid[i] = num
	}
	return grid, nil
}
