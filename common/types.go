// Copyright 2024 The go-solana Authors
// This file is part of the go-solana library.

package common

import (
	"bytes"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"github.com/mr-tron/base58"
	"math/big"
)

// Lengths of signatures and addresses in bytes.
const (
	// HashLength is the expected length of the hash
	HashLength = 32
	// AddressLength is the expected length of the address
	AddressLength = 32
	// SignatureLength is the expected length of the signature
	SignatureLength = 64
)

/////// -------------------------------------------------///////
/////// -------------------------------------------------///////
/////// -------------------- Address --------------------///////
/////// -------------------- Address --------------------///////
/////// -------------------------------------------------///////
/////// -------------------------------------------------///////

// Address The address
type Address [AddressLength]byte

// BytesToAddress returns Address with value b.
func BytesToAddress(b []byte) (a Address) {
	a.SetBytes(b)
	return
}

// BigToAddress returns Address with byte values of b.
func BigToAddress(b *big.Int) Address { return BytesToAddress(b.Bytes()) }

// Base58ToAddress returns Address with byte values of b.
func Base58ToAddress(b string) Address {
	// decode base58
	d, _ := base58.Decode(b)
	// bytes to address
	return BytesToAddress(d)
}

// Cmp compares two addresses.
func (a Address) Cmp(other Address) int {
	return bytes.Compare(a[:], other[:])
}

// Bytes return Address bytes
func (a Address) Bytes() []byte { return a[:] }

// Big return Address to *big.Int
func (a Address) Big() *big.Int { return new(big.Int).SetBytes(a[:]) }

// Base58 return base58 account
func (a Address) Base58() string {
	return base58.Encode(a[:])
}

// String return base58 account
func (a Address) String() string {
	return a.Base58()
}

// SetBytes sets the address to the value of b.
func (a *Address) SetBytes(b []byte) {
	if len(b) > len(a) {
		b = b[len(b)-AddressLength:]
	}
	copy(a[AddressLength-len(b):], b)
}

// MarshalText returns base58 str account
func (a Address) MarshalText() ([]byte, error) {
	input, err := json.Marshal(a.Base58())
	return input[1:len(input)-1], err
}

// UnmarshalText parses an account in base58 syntax.
func (a *Address) UnmarshalText(input []byte) error {
	a.SetBytes(input)
	return nil
}

// UnmarshalJSON parses an account in base58 syntax.
func (a *Address) UnmarshalJSON(input []byte) error {
	a.SetBytes(input)
	return nil
}

// Scan implements Scanner for database/sql.
func (a *Address) Scan(src interface{}) error {
	srcB, ok := src.([]byte)
	if !ok {
		return fmt.Errorf("can't scan %T into Address", src)
	}
	if len(srcB) != AddressLength {
		return fmt.Errorf("can't scan []byte of len %d into Address, want %d", len(srcB), AddressLength)
	}
	copy(a[:], srcB)
	return nil
}

// Value implements valuer for database/sql.
func (a Address) Value() (driver.Value, error) {
	return a[:], nil
}

// ImplementsGraphQLType returns true if Hash implements the specified GraphQL type.
func (a Address) ImplementsGraphQLType(name string) bool { return name == "Address" }

// UnmarshalGraphQL unmarshals the provided GraphQL query data.
func (a *Address) UnmarshalGraphQL(input interface{}) error {
	var err error
	switch input := input.(type) {
	case string:
		err = a.UnmarshalText([]byte(input))
	default:
		err = fmt.Errorf("unexpected type %T for Address", input)
	}
	return err
}

/////// ----------------------------------------------///////
/////// ----------------------------------------------///////
/////// -------------------- Hash --------------------///////
/////// -------------------- Hash --------------------///////
/////// ----------------------------------------------///////
/////// ----------------------------------------------///////

// Hash The Hash
type Hash [HashLength]byte

// BytesToHash returns Hash with value b.
func BytesToHash(b []byte) (h Hash) {
	h.SetBytes(b)
	return
}

// BigToHash returns Hash with byte values of b.
func BigToHash(b *big.Int) Hash { return BytesToHash(b.Bytes()) }

// Base58ToHash returns Hash with byte values of b.
func Base58ToHash(b string) Hash {
	// decode base58
	d, _ := base58.Decode(b)
	// bytes to Hash
	return BytesToHash(d)
}

// Cmp compares two Hashes.
func (h Hash) Cmp(other Hash) int {
	return bytes.Compare(h[:], other[:])
}

// Bytes return Hash bytes
func (h Hash) Bytes() []byte { return h[:] }

// Big return Hash to *big.Int
func (h Hash) Big() *big.Int { return new(big.Int).SetBytes(h[:]) }

// Base58 return base58 account
func (h Hash) Base58() string {
	return base58.Encode(h[:])
}

// String return base58 account
func (h Hash) String() string {
	return h.Base58()
}

// SetBytes sets the Hash to the value of b.
func (h *Hash) SetBytes(b []byte) {
	if len(b) > len(h) {
		b = b[len(b)-HashLength:]
	}
	copy(h[HashLength-len(b):], b)
}

// MarshalText returns base58 str hash
func (h Hash) MarshalText() ([]byte, error) {
	input, err := json.Marshal(h.Base58())
	return input[1:len(input)-1], err
}

// UnmarshalText parses a hash in base58 syntax.
func (h *Hash) UnmarshalText(input []byte) error {
	h.SetBytes(input)
	return nil
}

// UnmarshalJSON parses a hash in base58 syntax.
func (h *Hash) UnmarshalJSON(input []byte) error {
	h.SetBytes(input)
	return nil
}

// Scan implements Scanner for database/sql.
func (h *Hash) Scan(src interface{}) error {
	srcB, ok := src.([]byte)
	if !ok {
		return fmt.Errorf("can't scan %T into Hash", src)
	}
	if len(srcB) != HashLength {
		return fmt.Errorf("can't scan []byte of len %d into Hash, want %d", len(srcB), HashLength)
	}
	copy(h[:], srcB)
	return nil
}

// Value implements valuer for database/sql.
func (h Hash) Value() (driver.Value, error) {
	return h[:], nil
}

// ImplementsGraphQLType returns true if Hash implements the specified GraphQL type.
func (h Hash) ImplementsGraphQLType(name string) bool { return name == "Hash" }

// UnmarshalGraphQL unmarshals the provided GraphQL query data.
func (h *Hash) UnmarshalGraphQL(input interface{}) error {
	var err error
	switch input := input.(type) {
	case string:
		err = h.UnmarshalText([]byte(input))
	default:
		err = fmt.Errorf("unexpected type %T for Hash", input)
	}
	return err
}

/////// ----------------------------------------------------///////
/////// ----------------------------------------------------///////
/////// -------------------- Base58Data --------------------///////
/////// -------------------- Base58Data --------------------///////
/////// ----------------------------------------------------///////
/////// ----------------------------------------------------///////


// Base58Data base58 data
type Base58Data []byte

// Base58 return base58 str
func (b Base58Data) Base58() string {
	return base58.Encode(b)
}

// String return base58 str
func (b Base58Data) String() string {
	return b.Base58()
}

// SetBytes sets the Base58Data to the value of b.
func (b *Base58Data) SetBytes(input []byte) {
	*b = input
}

// MarshalText returns base58 str
func (b Base58Data) MarshalText() ([]byte, error) {
	input, err := json.Marshal(b.Base58())
	return input[1:len(input)-1], err
}

// UnmarshalText parses data in base58 syntax.
func (b *Base58Data) UnmarshalText(input []byte) error {
	b.SetBytes(input)
	return nil
}

// UnmarshalJSON parses data in base58 syntax.
func (b *Base58Data) UnmarshalJSON(input []byte) error {
	b.SetBytes(input)
	return nil
}

// Value implements valuer for database/sql.
func (b Base58Data) Value() (driver.Value, error) {
	return b[:], nil
}

/////// ---------------------------------------------------///////
/////// ---------------------------------------------------///////
/////// -------------------- Signature --------------------///////
/////// -------------------- Signature --------------------///////
/////// ---------------------------------------------------///////
/////// ---------------------------------------------------///////

// Signature The signature
type Signature [SignatureLength]byte

// BytesToSignature returns Signature with value b.
func BytesToSignature(b []byte) (a Signature) {
	a.SetBytes(b)
	return
}

// BigToSignature returns Signature with byte values of b.
func BigToSignature(b *big.Int) Signature { return BytesToSignature(b.Bytes()) }

// Base58ToSignature returns Signature with byte values of b.
func Base58ToSignature(b string) Signature {
	// decode base58
	d, _ := base58.Decode(b)
	// bytes to address
	return BytesToSignature(d)
}

// Cmp compares two addresses.
func (s Signature) Cmp(other Signature) int {
	return bytes.Compare(s[:], other[:])
}

// Bytes return Signature bytes
func (s Signature) Bytes() []byte { return s[:] }

// Big return Signature to *big.Int
func (s Signature) Big() *big.Int { return new(big.Int).SetBytes(s[:]) }

// Base58 return base58 account
func (s Signature) Base58() string {
	return base58.Encode(s[:])
}

// String return base58 account
func (s Signature) String() string {
	return s.Base58()
}

// SetBytes sets the address to the value of b.
func (s *Signature) SetBytes(b []byte) {
	if len(b) > len(s) {
		b = b[len(b)-SignatureLength:]
	}
	copy(s[SignatureLength-len(b):], b)
}

// MarshalText returns base58 str account
func (s Signature) MarshalText() ([]byte, error) {
	input, err := json.Marshal(s.Base58())
	return input[1:len(input)-1], err
}

// UnmarshalText parses an account in base58 syntax.
func (s *Signature) UnmarshalText(input []byte) error {
	s.SetBytes(input)
	return nil
}

// UnmarshalJSON parses an account in base58 syntax.
func (s *Signature) UnmarshalJSON(input []byte) error {
	s.SetBytes(input)
	return nil
}

// Scan implements Scanner for database/sql.
func (s *Signature) Scan(src interface{}) error {
	srcB, ok := src.([]byte)
	if !ok {
		return fmt.Errorf("can't scan %T into Signature", src)
	}
	if len(srcB) != SignatureLength {
		return fmt.Errorf("can't scan []byte of len %d into Signature, want %d", len(srcB), SignatureLength)
	}
	copy(s[:], srcB)
	return nil
}

// Value implements valuer for database/sql.
func (s Signature) Value() (driver.Value, error) {
	return s[:], nil
}

// ImplementsGraphQLType returns true if Hash implements the specified GraphQL type.
func (s Signature) ImplementsGraphQLType(name string) bool { return name == "Signature" }

// UnmarshalGraphQL unmarshals the provided GraphQL query dats.
func (s *Signature) UnmarshalGraphQL(input interface{}) error {
	var err error
	switch input := input.(type) {
	case string:
		err = s.UnmarshalText([]byte(input))
	default:
		err = fmt.Errorf("unexpected type %T for Signature", input)
	}
	return err
}

