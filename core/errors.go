// Copyright 2024 The go-solana Authors
// This file is part of the go-solana library.


package core

import (
	"errors"
	"fmt"
)

var (
	ErrEmptySlice = errors.New("empty slice found")
	ErrEmptyString = errors.New("empty string found")
	ErrEmptyAccount = errors.New("empty account found")
)

// StdErr return standard Err
func StdErr(reason string, err error) error {
	return fmt.Errorf("%s Failed. Err: %w", reason, err)
}
