// Copyright 2024 The go-solana Authors
// This file is part of the go-solana library.


package core

import (
	"encoding/json"
	"fmt"
	"github.com/mr-tron/base58"
)

// Has0xPrefix input has 0x prefix
func Has0xPrefix(input string) bool {
	return len(input) >= 2 && input[0] == '0' && (input[1] == 'x' || input[1] == 'X')
}

// DecodeBase58Str input string
func DecodeBase58Str(input string) []byte {
	data, _ := base58.Decode(input)
	return data
}

// BeautifyConsole console the content with json format
func BeautifyConsole(title, content any) {
	// MarshalIndent
	jsonData, _ := json.MarshalIndent(content, "", "    ")
	// print data
	fmt.Println(title, string(jsonData))
}


