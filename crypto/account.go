// Copyright 2024 The go-solana Authors
// This file is part of the go-solana library.

package crypto

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/mr-tron/base58"
	"github.com/tyler-smith/go-bip39"
	"go-solana/common"
	"go-solana/core"
	"go-solana/core/hdwallet"
	"os"
	"regexp"
)

type Account struct {
	Address    common.Address
	PrivateKey ed25519.PrivateKey
}

// GenerateAccount Random a new account from ed25519
func GenerateAccount() (Account, error) {
	var account Account
	// Random generateKey
	pub, prv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return account, err
	}
	copy(account.Address[:], pub)
	account.PrivateKey = prv
	// return account
	return account, err
}

// GenerateBase58PrvKey return base58 private key
func GenerateBase58PrvKey(a Account) (string, error) {
	// empty account
	if len(a.PrivateKey) == 0 {
		return "", core.ErrEmptyAccount
	}
	return base58.Encode(a.PrivateKey), nil
}

// GenerateHexPrvKey return hex private key
func GenerateHexPrvKey(a Account) (string, error) {
	// empty account
	if len(a.PrivateKey) == 0 {
		return "", core.ErrEmptyAccount
	}
	enHexKey := hex.EncodeToString(a.PrivateKey)
	// return 0x + enHexKey
	return "0x" + enHexKey, nil
}

// AccountFromBytes generate an account by bytes
func AccountFromBytes(b []byte) (Account, error) {
	// match privateKeySize
	if len(b) != ed25519.PrivateKeySize {
		return Account{}, fmt.Errorf("PrivateKey size mismatch, expected: %v, got: %v", ed25519.PrivateKeySize, len(b))
	}
	account := Account{PrivateKey: ed25519.PrivateKey(b)}
	// bytes to address
	account.Address = common.BytesToAddress(account.PrivateKey.Public().(ed25519.PublicKey))
	// return account
	return account, nil
}

// AccountFromBase58Key generate an account by base58 private key
func AccountFromBase58Key(key string) (Account, error) {
	// empty string
	if len(key) == 0 {
		return Account{}, core.ErrEmptyString
	}
	b, err := base58.Decode(key)
	// if err
	if err != nil {
		return Account{}, core.StdErr("AccountFromBase58", err)
	}
	return AccountFromBytes(b)
}

// AccountFromHexKey generate an account by hex private key
func AccountFromHexKey(key string) (Account, error) {
	// empty string
	if len(key) == 0 {
		return Account{}, core.ErrEmptyString
	}
	// has 0x prefix
	if core.Has0xPrefix(key) {
		key = key[2:]
	}
	// DecodeString
	b, err := hex.DecodeString(key)
	// if err
	if err != nil {
		return Account{}, core.StdErr("AccountFromHex", err)
	}
	return AccountFromBytes(b)
}

// AccountFromSeed generate an account by seed
func AccountFromSeed(seed []byte) (Account, error) {
	pk := ed25519.NewKeyFromSeed(seed)
	return AccountFromBytes(pk)
}

// AccountFromMnemonic generate an account by mnemonic
// @params [args]: password--> string
// @params [args]: path --> bool[true] (default: m/44'/501'/0'/0')
// @params [args]: path --> string (format: m/44'/501'/0'/0')
func AccountFromMnemonic(mnemonic string, args ...interface{}) (Account, error) {
	var (
		err      error
		seed     []byte
		path     string
		password string
	)
	// has password
	if len(args) > 0 {
		// get
		for _, arg := range args {
			switch v := arg.(type) {
			case string:
				// regexp path
				if ok, _ := regexp.MatchString(`^m(/\d+')*$`, v); ok {
					path = v
				} else {
					password = v
				}
			case bool:
				if v {
					path = "m/44'/501'/0'/0'"
				}
			}
		}
	}
	// Check mnemonic
	seed, err = bip39.NewSeedWithErrorChecking(mnemonic, password)
	// New Seed Failed
	if err != nil {
		return Account{}, core.StdErr("NewSeedWithErrorChecking", err)
	}
	// has path
	if path != "" {
		derivedKey, _ := hdwallet.Derived(path, seed)
		// override seed
		seed = derivedKey.PrivateKey
	}
	// return AccountFromSeed
	return AccountFromSeed(seed[:ed25519.SeedSize])
}

// AccountFromKeygenFile generate an account by keygen file
func AccountFromKeygenFile(file string) (Account, error) {
	// read file
	content, err := os.ReadFile(file)
	if err != nil {
		return Account{}, core.StdErr("read keygen file", err)
	}

	var values []byte
	// Unmarshal content
	if err = json.Unmarshal(content, &values); err != nil {
		return Account{}, core.StdErr("decode keygen file", err)
	}
	return AccountFromBytes(values)
}

// Sign the message with account
func (a Account) Sign(message []byte) []byte {
	return ed25519.Sign(a.PrivateKey, message)
}
