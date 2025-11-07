// Copyright 2024 The go-solana Authors
// This file is part of the go-solana library.

package solana

import (
	"bytes"
	"crypto/ed25519"
	"database/sql/driver"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math/big"

	"github.com/cielu/go-solana/core"
	"github.com/mr-tron/base58"
)

// Lengths of signatures and PublicKeys in bytes.
const (
	// HashLength is the expected length of the hash
	HashLength = 32
	// PublicKeyLength is the expected length of the PublicKey
	PublicKeyLength = 32
	// SignatureLength is the expected length of the signature
	SignatureLength = 64
)

// ///// -------------------------------------------------///////
// ///// -------------------------------------------------///////
// ///// -------------------- PublicKey --------------------///////
// ///// -------------------- PublicKey --------------------///////
// ///// -------------------------------------------------///////
// ///// -------------------------------------------------///////

// PublicKey The PublicKey
type PublicKey [PublicKeyLength]byte

// BytesToPublicKey returns PublicKey with value b.
func BytesToPublicKey(b []byte) (a PublicKey) {
	a.SetBytes(b)
	return
}

// BigIntToPublicKey returns PublicKey with byte values of b.
func BigIntToPublicKey(b *big.Int) PublicKey { return BytesToPublicKey(b.Bytes()) }

// StrToPublicKey returns PublicKey with byte values of b.
// Notice: only support base58/base64 str
func StrToPublicKey(b string) PublicKey {
	// decode base58 str
	if d, err := base58.Decode(b); err == nil {
		return BytesToPublicKey(d)
	}
	// decode base64 str
	if d, err := base64.StdEncoding.DecodeString(b); err == nil {
		return BytesToPublicKey(d)
	}
	// empty
	return PublicKey{}
}

// Base58ToPublicKey returns PublicKey with byte values of b.
func Base58ToPublicKey(b string) PublicKey {
	// decode base58
	d, _ := base58.Decode(b)
	// bytes to PublicKey
	return BytesToPublicKey(d)
}

// Base64ToPublicKey returns PublicKey with byte values of b.
func Base64ToPublicKey(b string) PublicKey {
	// decode base64
	d, _ := base64.StdEncoding.DecodeString(b)
	// bytes to PublicKey
	return BytesToPublicKey(d)
}

// IsEmpty PublicKey is empty
func (a PublicKey) IsEmpty() bool {
	return a == PublicKey{}
}

// Equals compares PublicKey a eq b
func (a PublicKey) Equals(b PublicKey) bool  {
	return a==b
}

// Cmp compares two PublicKeyes.
func (a PublicKey) Cmp(other PublicKey) int {
	return bytes.Compare(a[:], other[:])
}

// Bytes return PublicKey bytes
func (a PublicKey) Bytes() []byte { return a[:] }

// Big return PublicKey to *big.Int
func (a PublicKey) Big() *big.Int { return new(big.Int).SetBytes(a[:]) }

// Base58 return base58 account
func (a PublicKey) Base58() string {
	return base58.Encode(a[:])
}

// String return base58 account
func (a PublicKey) String() string {
	return a.Base58()
}

// SetBytes sets the PublicKey to the value of b.
func (a *PublicKey) SetBytes(b []byte) {
	if len(b) > len(a) {
		b = b[len(b)-PublicKeyLength:]
	}
	copy(a[PublicKeyLength-len(b):], b)
}

// MarshalText returns base58 str account
func (a PublicKey) MarshalText() ([]byte, error) {
	input, err := json.Marshal(a.Base58())
	return input[1 : len(input)-1], err
}

// UnmarshalText parses an account in base58 syntax.
func (a *PublicKey) UnmarshalText(input []byte) error {
	a.SetBytes(input)
	return nil
}

// UnmarshalJSON parses an account in base58 syntax.
func (a *PublicKey) UnmarshalJSON(input []byte) error {
	// Unmarshal data to []byte
	data, _, err := core.UnmarshalDataByEncoding(input)
	// set string to Hash
	a.SetBytes(data)
	return err
}

// Scan implements Scanner for database/sql.
func (a *PublicKey) Scan(src interface{}) error {
	srcB, ok := src.([]byte)
	if !ok {
		return fmt.Errorf("can't scan %T into PublicKey", src)
	}
	if len(srcB) != PublicKeyLength {
		return fmt.Errorf("can't scan []byte of len %d into PublicKey, want %d", len(srcB), PublicKeyLength)
	}
	copy(a[:], srcB)
	return nil
}

// Value implements valuer for database/sql.
func (a PublicKey) Value() (driver.Value, error) {
	return a[:], nil
}

// ImplementsGraphQLType returns true if Hash implements the specified GraphQL type.
func (a PublicKey) ImplementsGraphQLType(name string) bool { return name == "PublicKey" }

// UnmarshalGraphQL unmarshals the provided GraphQL query data.
func (a *PublicKey) UnmarshalGraphQL(input interface{}) error {
	var err error
	switch input := input.(type) {
	case string:
		err = a.UnmarshalText([]byte(input))
	default:
		err = fmt.Errorf("unexpected type %T for PublicKey", input)
	}
	return err
}

// ///// ----------------------------------------------///////
// ///// ----------------------------------------------///////
// ///// -------------------- Hash --------------------///////
// ///// -------------------- Hash --------------------///////
// ///// ----------------------------------------------///////
// ///// ----------------------------------------------///////

// Hash The Hash
type Hash [HashLength]byte

// BytesToHash returns Hash with value b.
func BytesToHash(b []byte) (h Hash) {
	h.SetBytes(b)
	return
}

// BigToHash returns Hash with byte values of b.
func BigToHash(b *big.Int) Hash { return BytesToHash(b.Bytes()) }

// StrToHash returns Hash with byte values of b.
// Notice: only support base58/base64 str
func StrToHash(b string) Hash {
	// decode base58 str
	if d, err := base58.Decode(b); err == nil {
		return BytesToHash(d)
	}
	// decode base64 str
	if d, err := base64.StdEncoding.DecodeString(b); err == nil {
		return BytesToHash(d)
	}
	// base 64
	return Hash{}
}

// Base58ToHash returns Hash with byte values of b.
func Base58ToHash(b string) Hash {
	// decode base58
	d, _ := base58.Decode(b)
	// bytes to Hash
	return BytesToHash(d)
}

// Base64ToHash returns Hash with byte values of b.
func Base64ToHash(b string) Hash {
	// decode base64
	d, _ := base64.StdEncoding.DecodeString(b)
	// bytes to PublicKey
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
	return input[1 : len(input)-1], err
}

// UnmarshalText parses a hash in base58 syntax.
func (h *Hash) UnmarshalText(input []byte) error {
	h.SetBytes(input)
	return nil
}

// UnmarshalJSON parses a hash in base58 syntax.
func (h *Hash) UnmarshalJSON(input []byte) error {
	// Unmarshal data to []byte
	data, _, err := core.UnmarshalDataByEncoding(input)
	// set string to Hash
	h.SetBytes(data)
	// return err
	return err
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

// ///// -------------------------------------------------///////
// ///// -------------------------------------------------///////
// ///// -------------------- Base58 ---------------------///////
// ///// -------------------- Base58 ---------------------///////
// ///// -------------------------------------------------///////
// ///// -------------------------------------------------///////

// type Base58 []byte
//
// func (t Base58) MarshalJSON() ([]byte, error) {
// 	return json.Marshal(base58.Encode(t))
// }
//
// func (t *Base58) UnmarshalJSON(data []byte) (err error) {
// 	var s string
// 	err = json.Unmarshal(data, &s)
// 	if err != nil {
// 		return
// 	}
// 	if s == "" {
// 		*t = []byte{}
// 		return nil
// 	}
// 	*t, err = base58.Decode(s)
// 	return
// }
//
// func (t Base58) String() string {
// 	return base58.Encode(t)
// }


// ///// -------------------------------------------------///////
// ///// -------------------------------------------------///////
// ///// -------------------- SolData --------------------///////
// ///// -------------------- SolData --------------------///////
// ///// -------------------------------------------------///////
// ///// -------------------------------------------------///////

// SolData base58, base64 data
type SolData struct {
	RawData  []byte
	Encoding string
}

// BytesToSolData default base58
func BytesToSolData(data []byte) (sd SolData) {
	sd.SetBytes(data)
	return
}

// Base58 return base58 str
func (sd SolData) Base58() string {
	return base58.Encode(sd.RawData)
}

func (sd SolData) Base64() string {
	return base64.StdEncoding.EncodeToString(sd.RawData)
}

// String return base58 str
func (sd SolData) String() string {
	// base64
	if sd.Encoding == "base64" {
		return sd.Base64()
	}
	return sd.Base58()
}

// SetBytes sets the SolData to the value of sd. (default base58)
func (sd *SolData) SetBytes(input []byte) {
	sd.RawData = input
}

// SetSolData sets the SolData
func (sd *SolData) SetSolData(data []byte, encoding string) {
	sd.RawData = data
	sd.Encoding = encoding
}

// MarshalText returns base58/base64 str
func (sd SolData) MarshalText() ([]byte, error) {
	input, err := json.Marshal(sd.String())
	return input[1 : len(input)-1], err
}

// UnmarshalText parses data in base58 syntax.
func (sd *SolData) UnmarshalText(input []byte) error {
	sd.SetBytes(input)
	return nil
}

// UnmarshalJSON parses data in base58 syntax.
func (sd *SolData) UnmarshalJSON(input []byte) error {
	// Unmarshal data to []byte
	data, encoding, err := core.UnmarshalDataByEncoding(input)
	// has err
	if err != nil {
		return err
	}
	// set SolData by encoding
	if encoding == "" {
		sd.SetBytes(data)
	} else {
		sd.SetSolData(data, encoding)
	}
	return nil
}

// Value implements valuer for database/sql.
func (sd SolData) Value() (driver.Value, error) {
	return sd.RawData[:], nil
}

// ///// ---------------------------------------------------///////
// ///// ---------------------------------------------------///////
// ///// -------------------- Signature --------------------///////
// ///// -------------------- Signature --------------------///////
// ///// ---------------------------------------------------///////
// ///// ---------------------------------------------------///////

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
	// bytes to PublicKey
	return BytesToSignature(d)
}

// StrToSignature returns Signature with byte values of b.
// Notice: only support base58/base64 str
func StrToSignature(b string) Signature {
	// decode base58 str
	if d, err := base58.Decode(b); err == nil {
		return BytesToSignature(d)
	}
	// decode base64 str
	if d, err := base64.StdEncoding.DecodeString(b); err == nil {
		return BytesToSignature(d)
	}
	// empty
	return Signature{}
}

// Cmp compares two PublicKeyes.
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

// SetBytes sets the PublicKey to the value of b.
func (s *Signature) SetBytes(b []byte) {
	if len(b) > len(s) {
		b = b[len(b)-SignatureLength:]
	}
	copy(s[SignatureLength-len(b):], b)
}

// MarshalText returns base58 str account
func (s Signature) MarshalText() ([]byte, error) {
	input, err := json.Marshal(s.Base58())
	return input[1 : len(input)-1], err
}

// UnmarshalText parses an account in base58 syntax.
func (s *Signature) UnmarshalText(input []byte) error {
	s.SetBytes(input)
	return nil
}

// UnmarshalJSON parses an account in base58 syntax.
func (s *Signature) UnmarshalJSON(input []byte) error {
	// Unmarshal data to []byte
	data, _, err := core.UnmarshalDataByEncoding(input)
	// set string to Hash
	s.SetBytes(data)
	// return err
	return err
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

func (s Signature) Sign(message []byte) []byte {
	return ed25519.Sign(s.Bytes(), message)
}
