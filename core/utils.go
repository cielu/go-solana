// Copyright 2024 The go-solana Authors
// This file is part of the go-solana library.


package core

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/mr-tron/base58"
)

// Has0xPrefix input has 0x prefix
func Has0xPrefix(input string) bool {
	return len(input) >= 2 && input[0] == '0' && (input[1] == 'x' || input[1] == 'X')
}

// UnmarshalDataByEncoding Unmarshal data to string by encoding
func UnmarshalDataByEncoding(input []byte) ([]byte, string, error) {
	var (
		err      error
		data     interface{}
		encoding string
	)
	// Unmarshal
	if err = json.Unmarshal(input, &data); err != nil {
		return input, "", err
	}
	// get active type
	switch v := data.(type) {
	case string:
		encoding = "base58"
		input = DecodeBase58Str(v)
	// slice
	case []interface{}:
		// none data
		if len(v) == 0 {
			return nil, "", err
		}
		// decode to string
		switch v[1] {
		case "base58":
			encoding = "base58"
			fmt.Println(v[0].(string))
			input = DecodeBase58Str(v[0].(string))
		case "base64":
			encoding = "base64"
			input, _ = base64.StdEncoding.DecodeString(v[0].(string))
		default:
			return nil, "", fmt.Errorf("UnmarshalDataByEncoding Err: %s", v[1])
		}
	}
	return input, encoding, err
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

// DecodeBase58Str input string
func DecodeBase58Str(input string) []byte {
	data, _ := base58.Decode(input)
	return data
}

// BeautifyConsole console the content with json format
func BeautifyConsole(title string, content any) {
	// MarshalIndent
	jsonData, _ := json.MarshalIndent(content, "", "    ")
	// print data
	fmt.Println(title, string(jsonData))
}


