package base

import (
	"crypto/sha256"
	"errors"
	"math"

	"github.com/cielu/go-solana/common"
)

const (
	// MaxSeedLength Maximum length of derived pubkey seed.
	MaxSeedLength = 32
	// MaxSeeds Maximum number of seeds.
	MaxSeeds = 16
	// Number of bytes in a signature.
)

const PDA_MARKER = "ProgramDerivedAddress"

var ErrMaxSeedLengthExceeded = errors.New("Max seed length exceeded")

// CreateProgramAddress Create a program address.
// Ported from https://github.com/solana-labs/solana/blob/216983c50e0a618facc39aa07472ba6d23f1b33a/sdk/program/src/pubkey.rs#L204
func CreateProgramAddress(seeds [][]byte, programID common.Address) (common.Address, error) {
	if len(seeds) > MaxSeeds {
		return common.Address{}, ErrMaxSeedLengthExceeded
	}

	for _, seed := range seeds {
		if len(seed) > MaxSeedLength {
			return common.Address{}, ErrMaxSeedLengthExceeded
		}
	}

	var buf []byte
	for _, seed := range seeds {
		buf = append(buf, seed...)
	}

	buf = append(buf, programID[:]...)
	buf = append(buf, []byte(PDA_MARKER)...)
	hash := sha256.Sum256(buf)

	if IsOnCurve(hash[:]) {
		return common.Address{}, errors.New("invalid seeds; address must fall off the curve")
	}

	return common.BytesToAddress(hash[:]), nil
}

type incomparable [0]func()

// Point represents a point on the edwards25519 curve.
//
// This type works similarly to math/big.Int, and all arguments and receivers
// are allowed to alias.
//
// The zero value is NOT valid, and it may be used only as a receiver.
type Point struct {
	// The point is internally represented in extended coordinates (X, Y, Z, T)
	// where x = X/Z, y = Y/Z, and xy = T/Z per https://eprint.iacr.org/2008/522.
	x, y, z, t Element

	// Make the type not comparable (i.e. used with == or as a map key), as
	// equivalent points can be represented by different Go values.
	_ incomparable
}

// d is a constant in the curve equation.
var d, _ = new(Element).SetBytes([]byte{
	0xa3, 0x78, 0x59, 0x13, 0xca, 0x4d, 0xeb, 0x75,
	0xab, 0xd8, 0x41, 0x41, 0x4d, 0x0a, 0x70, 0x00,
	0x98, 0xe8, 0x79, 0x77, 0x79, 0x40, 0xc7, 0x8c,
	0x73, 0xfe, 0x6f, 0x2b, 0xee, 0x6c, 0x03, 0x52})

// SetBytes sets v = x, where x is a 32-byte encoding of v. If x does not
// represent a valid point on the curve, SetBytes returns nil and an error and
// the receiver is unchanged. Otherwise, SetBytes returns v.
//
// Note that SetBytes accepts all non-canonical encodings of valid points.
// That is, it follows decoding rules that match most implementations in
// the ecosystem rather than RFC 8032.
func (v *Point) SetBytes(x []byte) (*Point, error) {
	// Specifically, the non-canonical encodings that are accepted are
	//   1) the ones where the field element is not reduced (see the
	//      (*field.Element).SetBytes docs) and
	//   2) the ones where the x-coordinate is zero and the sign bit is set.
	//
	// This is consistent with crypto/ed25519/internal/edwards25519. Read more
	// at https://hdevalence.ca/blog/2020-10-04-its-25519am, specifically the
	// "Canonical A, R" section.

	y, err := new(Element).SetBytes(x)
	if err != nil {
		return nil, errors.New("edwards25519: invalid point encoding length")
	}

	// -x² + y² = 1 + dx²y²
	// x² + dx²y² = x²(dy² + 1) = y² - 1
	// x² = (y² - 1) / (dy² + 1)

	// u = y² - 1
	y2 := new(Element).Square(y)
	u := new(Element).Subtract(y2, feOne)

	// v = dy² + 1
	vv := new(Element).Multiply(y2, d)
	vv = vv.Add(vv, feOne)

	// x = +√(u/v)
	xx, wasSquare := new(Element).SqrtRatio(u, vv)
	if wasSquare == 0 {
		return nil, errors.New("edwards25519: invalid point encoding")
	}

	// Select the negative square root if the sign bit is set.
	xxNeg := new(Element).Negate(xx)
	xx = xx.Select(xxNeg, xx, int(x[31]>>7))

	v.x.Set(xx)
	v.y.Set(y)
	v.z.One()
	v.t.Multiply(xx, y) // xy = T / Z

	return v, nil
}

// Check if the provided `b` is on the ed25519 curve.
func IsOnCurve(b []byte) bool {
	_, err := new(Point).SetBytes(b)
	isOnCurve := err == nil
	return isOnCurve
}

// FindProgramAddress Find a valid program address and its corresponding bump seed.
func FindProgramAddress(seed [][]byte, programID common.Address) (common.Address, uint8, error) {
	var address common.Address
	var err error
	bumpSeed := uint8(math.MaxUint8)
	for bumpSeed != 0 {
		address, err = CreateProgramAddress(append(seed, []byte{bumpSeed}), programID)
		if err == nil {
			return address, bumpSeed, nil
		}
		bumpSeed--
	}
	return common.Address{}, bumpSeed, errors.New("unable to find a valid program address")
}

func FindAssociatedTokenAddress(account common.Address, mint common.Address, options ...common.Address) (common.Address, uint8, error) {
	return FindAssociatedTokenAddressAndBumpSeed(account, mint, SPLAssociatedTokenAccountProgramID, options...)
}

func FindAssociatedTokenAddressAndBumpSeed(account common.Address, splTokenMintAddress common.Address, programID common.Address, options ...common.Address) (common.Address, uint8, error) {
	tokenProgramID := TokenProgramID
	if len(options) > 0 && options[0] == Token2022ProgramID {
		tokenProgramID = Token2022ProgramID
	}
	return FindProgramAddress([][]byte{account[:], tokenProgramID[:], splTokenMintAddress[:]}, programID)
}
