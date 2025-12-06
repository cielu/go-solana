// Copyright 2024 The go-solana Authors
// This file is part of the go-solana library.

package solana

import (
	"bytes"
	"crypto/ed25519"
	"database/sql/driver"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"

	bin "github.com/cielu/go-solana/pkg/encodbin"
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

// StrToPublicKey returns PublicKey with byte values of b.
// Notice: only support base58/base64 str
func StrToPublicKey(b string) PublicKey {
	// decode base58 str
	if d, err := base58.Decode(b); err == nil {
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

// IsEmpty PublicKey is empty
func (p PublicKey) IsEmpty() bool {
	return p == PublicKey{}
}

// Equals compares PublicKey a eq b
func (p PublicKey) Equals(b PublicKey) bool {
	return p == b
}

// Cmp compares two PublicKeyes.
func (p PublicKey) Cmp(other PublicKey) int {
	return bytes.Compare(p[:], other[:])
}

// Bytes return PublicKey bytes
func (p PublicKey) Bytes() []byte { return p[:] }

// Base58 return base58 account
func (p PublicKey) Base58() string {
	return base58.Encode(p[:])
}

// String return base58 account
func (p PublicKey) String() string {
	return p.Base58()
}

// SetBytes sets the PublicKey to the value of b.
func (p *PublicKey) SetBytes(b []byte) {
	if len(b) > len(p) {
		b = b[len(b)-PublicKeyLength:]
	}
	copy(p[PublicKeyLength-len(b):], b)
}

// MarshalText returns base58 str account
func (p PublicKey) MarshalText() ([]byte, error) {
	input, err := json.Marshal(p.Base58())
	return input[1 : len(input)-1], err
}

// UnmarshalText parses an account in base58 syntax.
func (p *PublicKey) UnmarshalText(input []byte) error {
	p.SetBytes(input)
	return nil
}

func (p PublicKey) MarshalJSON() ([]byte, error) {
	return json.Marshal(p.Base58())
}

// UnmarshalJSON parses an account in base58 syntax.
func (p *PublicKey) UnmarshalJSON(input []byte) error {
	// Unmarshal data to []byte
	var s string
	if err := json.Unmarshal(input, &s); err != nil {
		return err
	}
	// Decode
	if val, err := base58.Decode(s); err != nil {
		return err
	} else {
		p.SetBytes(val)
	}
	return nil
}

// Scan implements Scanner for database/sql.
func (p *PublicKey) Scan(src interface{}) error {
	srcB, ok := src.([]byte)
	if !ok {
		return fmt.Errorf("can't scan %T into PublicKey", src)
	}
	if len(srcB) != PublicKeyLength {
		return fmt.Errorf("can't scan []byte of len %d into PublicKey, want %d", len(srcB), PublicKeyLength)
	}
	p.SetBytes(srcB)
	return nil
}

// Value implements valuer for database/sql.
func (p PublicKey) Value() (driver.Value, error) {
	return p.String(), nil
}

// ImplementsGraphQLType returns true if Hash implements the specified GraphQL type.
func (p PublicKey) ImplementsGraphQLType(name string) bool { return name == "PublicKey" }

// UnmarshalGraphQL unmarshals the provided GraphQL query data.
func (p *PublicKey) UnmarshalGraphQL(input interface{}) error {
	var err error
	switch v := input.(type) {
	case string:
		err = p.UnmarshalText([]byte(v))
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

// StrToHash returns Hash with byte values of b.
// Notice: only support base58/base64 str
func StrToHash(b string) Hash {
	// decode base58 str
	if d, err := base58.Decode(b); err == nil {
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

// Cmp compares two Hashes.
func (h Hash) Cmp(other Hash) int {
	return bytes.Compare(h[:], other[:])
}

// Bytes return Hash bytes
func (h Hash) Bytes() []byte { return h[:] }

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

func (h Hash) MarshalJSON() ([]byte, error) {
	return json.Marshal(h.Base58())
}

// UnmarshalJSON parses a hash in base58 syntax.
func (h *Hash) UnmarshalJSON(input []byte) error {
	// Unmarshal data to []byte
	var s string
	if err := json.Unmarshal(input, &s); err != nil {
		return err
	}
	// Decode
	if val, err := base58.Decode(s); err != nil {
		return err
	} else {
		h.SetBytes(val)
	}
	// return err
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
	h.SetBytes(srcB)
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

type Base58Data []byte

func (t Base58Data) MarshalJSON() ([]byte, error) {
	return json.Marshal(base58.Encode(t))
}

func (t *Base58Data) UnmarshalJSON(data []byte) (err error) {
	var s string
	err = json.Unmarshal(data, &s)
	if err != nil {
		return
	}
	if s == "" {
		*t = []byte{}
		return nil
	}
	*t, err = base58.Decode(s)
	return
}

func (t Base58Data) String() string {
	return t.Base58()
}

func (t Base58Data) Hex() string {
	return hex.EncodeToString(t)
}

func (t Base58Data) Base58() string {
	return base58.Encode(t)
}

// ///// -------------------------------------------------///////
// ///// -------------------------------------------------///////
// ///// -------------------- SolData --------------------///////
// ///// -------------------- SolData --------------------///////
// ///// -------------------------------------------------///////
// ///// -------------------------------------------------///////

// SolData base58, base64 data
type SolData struct {
	rawData  []byte
	encoding EncodingEnum
}

// BytesToSolData default base58
func BytesToSolData(data []byte) (sd SolData) {
	sd.SetBytes(data)
	return
}

func (sd SolData) RawData() []byte {
	return sd.rawData
}

func (sd SolData) Encoding() EncodingEnum {
	return sd.encoding
}

// Base58 return base58 str
func (sd SolData) Base58() string {
	return base58.Encode(sd.rawData)
}

func (sd SolData) Base64() string {
	return base64.StdEncoding.EncodeToString(sd.rawData)
}

// String return base58 str
func (sd SolData) String() string {
	// base64
	if sd.encoding == EncodingBase64 {
		return sd.Base64()
	}
	return sd.Base58()
}

// SetBytes sets the SolData to the value of sd. (default base58)
func (sd *SolData) SetBytes(input []byte) {
	sd.rawData = input
}

// SetSolData sets the SolData
func (sd *SolData) SetSolData(data []byte, encoding EncodingEnum) {
	sd.rawData = data
	sd.encoding = encoding
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

func (sd SolData) MarshalJSON() ([]byte, error) {
	return json.Marshal([]interface{}{sd.String(), sd.encoding})
}

// UnmarshalJSON parses data in base58 syntax.
func (sd *SolData) UnmarshalJSON(input []byte) error {
	// Unmarshal data to []byte
	var (
		err  error
		data []string
	)
	if err = json.Unmarshal(input, &data); err != nil {
		return err
	}

	if len(data) != 2 {
		return fmt.Errorf("invalid length for SolData, expected 2, found %d", len(data))
	}
	// SetEncoding
	sd.encoding = EncodingEnum(data[1])

	content := data[0]
	// empty string
	if content == "" {
		return nil
	}
	// Decode by encoding
	switch sd.encoding {
	case EncodingBase58:
		sd.rawData, err = base58.Decode(content)
	case EncodingBase64:
		sd.rawData, err = base64.StdEncoding.DecodeString(content)
	default:
		return fmt.Errorf("unknown encoding: %s", sd.encoding)
	}
	return err
}

func (sd SolData) MarshalWithEncoder(encoder *bin.Encoder) (err error) {
	err = encoder.WriteBytes(sd.rawData, true)
	if err != nil {
		return err
	}
	err = encoder.WriteString(string(sd.encoding))
	if err != nil {
		return err
	}
	return nil
}

func (sd *SolData) UnmarshalWithDecoder(decoder *bin.Decoder) (err error) {
	sd.rawData, err = decoder.ReadByteSlice()
	if err != nil {
		return err
	}
	{
		enc, err := decoder.ReadString()
		if err != nil {
			return err
		}
		sd.encoding = EncodingEnum(enc)
	}
	return nil
}

// Value implements valuer for database/sql.
func (sd SolData) Value() (driver.Value, error) {
	return sd.rawData[:], nil
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
func BytesToSignature(b []byte) (s Signature) {
	s.SetBytes(b)
	return
}

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
	// empty
	return Signature{}
}

// Bytes return Signature bytes
func (s Signature) Bytes() []byte { return s[:] }

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

func (s Signature) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.Base58())
}

// UnmarshalJSON parses an account in base58 syntax.
func (s *Signature) UnmarshalJSON(input []byte) error {
	// Unmarshal data to []byte
	var data string
	if err := json.Unmarshal(input, &data); err != nil {
		return err
	}
	// Decode
	if val, err := base58.Decode(data); err != nil {
		return err
	} else {
		s.SetBytes(val)
	}
	// return err
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
	return s.String(), nil
}

// ImplementsGraphQLType returns true if Hash implements the specified GraphQL type.
func (s Signature) ImplementsGraphQLType(name string) bool {
	return name == "Signature"
}

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
