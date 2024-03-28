// Copyright 2024 The go-solana Authors
// This file is part of the go-solana library.


package core

import (
	"encoding/json"
	"fmt"
)

// Has0xPrefix input has 0x prefix
func Has0xPrefix(input string) bool {
	return len(input) >= 2 && input[0] == '0' && (input[1] == 'x' || input[1] == 'X')
}

// BeautifyConsole console the content with json format
func BeautifyConsole(title, content any)  {
	// MarshalIndent
	jsonData, _ := json.MarshalIndent(content, "", "    ")
	// print data
	fmt.Println(title, string(jsonData))
}
