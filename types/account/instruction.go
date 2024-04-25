package account

import (
	"crypto/sha256"
	"errors"
	"filippo.io/edwards25519"
	"github.com/cielu/go-solana/common"
)

const (
	PublicKeyLength = 32
	MaxSeedLength   = 32
	MaxSeed         = 16
)

func FindAssociatedTokenAddress(walletAddress, tokenMintAddress common.Address) (common.Address, uint8, error) {
	seeds := [][]byte{}
	seeds = append(seeds, walletAddress.Bytes())
	seeds = append(seeds, common.TokenProgramID.Bytes())
	seeds = append(seeds, tokenMintAddress.Bytes())

	return FindProgramAddress(seeds, common.SPLAssociatedTokenAccountProgramID)
}

func FindProgramAddress(seed [][]byte, programID common.Address) (common.Address, uint8, error) {
	var pubKey common.Address
	var err error
	var nonce uint8 = 0xff
	for nonce != 0x0 {
		pubKey, err = CreateProgramAddress(append(seed, []byte{byte(nonce)}), programID)
		if err == nil {
			return pubKey, nonce, nil
		}
		nonce--
	}
	return common.Address{}, nonce, errors.New("unable to find a viable program address")
}

func IsOnCurve(p common.Address) bool {
	_, err := new(edwards25519.Point).SetBytes(p.Bytes())
	return err == nil
}

func CreateProgramAddress(seeds [][]byte, programId common.Address) (common.Address, error) {
	if len(seeds) > MaxSeed {
		return common.Address{}, errors.New("length of the seed is too long for address generation")
	}

	buf := []byte{}
	for _, seed := range seeds {
		if len(seed) > MaxSeedLength {
			return common.Address{}, errors.New("length of the seed is too long for address generation")
		}
		buf = append(buf, seed...)
	}
	buf = append(buf, programId[:]...)
	buf = append(buf, []byte("ProgramDerivedAddress")...)
	h := sha256.Sum256(buf)

	pubkey := PublicKeyFromBytes(h[:])
	if IsOnCurve(pubkey) {
		return common.Address{}, errors.New("invalid seeds, address must fall off the curve")
	}
	return pubkey, nil
}

func PublicKeyFromBytes(b []byte) common.Address {
	var pubkey common.Address
	if len(b) > PublicKeyLength {
		b = b[:PublicKeyLength]
	}
	copy(pubkey[PublicKeyLength-len(b):], b)
	return pubkey
}
