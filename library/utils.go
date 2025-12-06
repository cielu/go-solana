// Copyright 2024 The go-solana Authors
// This file is part of the go-solana library.


package library

// Has0xPrefix input has 0x prefix
func Has0xPrefix(input string) bool {
	return len(input) >= 2 && input[0] == '0' && (input[1] == 'x' || input[1] == 'X')
}

// UniqueAppend judge and append key
func UniqueAppend[T comparable](slice []T, lookup T) []T {
	// append unique key
	hasKey := false
	for _, key := range slice {
		// found the key
		if key == lookup {
			hasKey = true
			break
		}
	}
	// not found
	if !hasKey {
		slice = append(slice, lookup)
	}
	return slice
}

// SliceSplitFrom Split Slice
func SliceSplitFrom[T comparable](slice []T, index int) (first []T, second []T) {

	if index < 0 {
		panic("negative index")
	}
	if index == 0 {
		second = slice
		return
	}
	if index > len(slice)-1 {
		first = slice
		return
	}
	first = slice[:index]
	second = slice[index:]
	return
}

