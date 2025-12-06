// Copyright 2021 github.com/gagliardetto
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package token

import (
	"errors"

	"github.com/cielu/go-solana"
	"github.com/cielu/go-solana/core/base"
	"github.com/cielu/go-solana/pkg/encodbin"
)

// InitializeMint Initializes a new mint and optionally deposits all the newly minted
// tokens in an account.
//
// The `InitializeMint` instruction requires no signers and MUST be
// included within the same Transaction as the system program's
// `CreateAccount` instruction that creates the account being initialized.
// Otherwise another party can acquire ownership of the uninitialized
// account.
type InitializeMint struct {
	// Number of base 10 digits to the right of the decimal place.
	Decimals *uint8

	// The authority/multisignature to mint tokens.
	MintAuthority *solana.PublicKey

	// The freeze authority/multisignature of the mint.
	FreezeAuthority *solana.PublicKey `bin:"optional"`

	// [0] = [WRITE] mint
	// ··········· The mint to initialize.
	//
	// [1] = [] $(SysVarRentPubkey)
	// ··········· Rent sysvar.
	solana.AccountMetaSlice `bin:"-" borsh_skip:"true"`
}

// NewInitializeMintInstructionBuilder creates a new `InitializeMint` instruction builder.
func NewInitializeMintInstructionBuilder() *InitializeMint {
	nd := &InitializeMint{
		AccountMetaSlice: make([]*solana.AccountMeta, 2),
	}
	nd.AccountMetaSlice[1] = solana.Meta(base.SysVarRentPubkey)
	return nd
}

// SetDecimals sets the "decimals" parameter.
// Number of base 10 digits to the right of the decimal place.
func (initMint *InitializeMint) SetDecimals(decimals uint8) *InitializeMint {
	initMint.Decimals = &decimals
	return initMint
}

// SetMintAuthority sets the "mint_authority" parameter.
// The authority/multisignature to mint tokens.
func (initMint *InitializeMint) SetMintAuthority(mint_authority solana.PublicKey) *InitializeMint {
	initMint.MintAuthority = &mint_authority
	return initMint
}

// SetFreezeAuthority sets the "freeze_authority" parameter.
// The freeze authority/multisignature of the mint.
func (initMint *InitializeMint) SetFreezeAuthority(freeze_authority solana.PublicKey) *InitializeMint {
	initMint.FreezeAuthority = &freeze_authority
	return initMint
}

// SetMintAccount sets the "mint" account.
// The mint to initialize.
func (initMint *InitializeMint) SetMintAccount(mint solana.PublicKey) *InitializeMint {
	initMint.AccountMetaSlice[0] = solana.Meta(mint).WRITE()
	return initMint
}

// GetMintAccount gets the "mint" account.
// The mint to initialize.
func (initMint *InitializeMint) GetMintAccount() *solana.AccountMeta {
	return initMint.AccountMetaSlice[0]
}

// SetSysVarRentPubkeyAccount sets the "$(SysVarRentPubkey)" account.
// Rent sysvar.
func (initMint *InitializeMint) SetSysVarRentPubkeyAccount(SysVarRentPubkey solana.PublicKey) *InitializeMint {
	initMint.AccountMetaSlice[1] = solana.Meta(SysVarRentPubkey)
	return initMint
}

// GetSysVarRentPubkeyAccount gets the "$(SysVarRentPubkey)" account.
// Rent sysvar.
func (initMint *InitializeMint) GetSysVarRentPubkeyAccount() *solana.AccountMeta {
	return initMint.AccountMetaSlice[1]
}

func (initMint InitializeMint) Build() *Instruction {
	return &Instruction{BaseVariant: encodbin.BaseVariant{
		Impl:   initMint,
		TypeID: encodbin.TypeIDFromUint8(Instruction_InitializeMint),
	}}
}

// ValidateAndBuild validates the instruction parameters and accounts;
// if there is a validation error, it returns the error.
// Otherwise, it builds and returns the instruction.
func (initMint InitializeMint) ValidateAndBuild() (*Instruction, error) {
	if err := initMint.Validate(); err != nil {
		return nil, err
	}
	return initMint.Build(), nil
}

func (initMint *InitializeMint) Validate() error {
	// Check whether all (required) parameters are set:
	{
		if initMint.Decimals == nil {
			return errors.New("Decimals parameter is not set")
		}
		if initMint.MintAuthority == nil {
			return errors.New("MintAuthority parameter is not set")
		}
	}

	// Check whether all (required) accounts are set:
	{
		if initMint.AccountMetaSlice[0] == nil {
			return errors.New("accounts.Mint is not set")
		}
		if initMint.AccountMetaSlice[1] == nil {
			return errors.New("accounts.SysVarRentPubkey is not set")
		}
	}
	return nil
}

func (initMint InitializeMint) MarshalWithEncoder(encoder *encodbin.Encoder) (err error) {
	// Serialize `Decimals` param:
	err = encoder.Encode(initMint.Decimals)
	if err != nil {
		return err
	}
	// Serialize `MintAuthority` param:
	err = encoder.Encode(initMint.MintAuthority)
	if err != nil {
		return err
	}
	// Serialize `FreezeAuthority` param (optional):
	{
		if initMint.FreezeAuthority == nil {
			err = encoder.WriteBool(false)
			if err != nil {
				return err
			}
		} else {
			err = encoder.WriteBool(true)
			if err != nil {
				return err
			}
			err = encoder.Encode(initMint.FreezeAuthority)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// NewInitializeMintInstruction declares a new InitializeMint instruction with the provided parameters and accounts.
func NewInitializeMintInstruction(
	// Parameters:
	decimals uint8,
	mint_authority solana.PublicKey,
	freeze_authority solana.PublicKey,
	// Accounts:
	mint solana.PublicKey) *InitializeMint {
	return NewInitializeMintInstructionBuilder().
		SetDecimals(decimals).
		SetMintAuthority(mint_authority).
		SetFreezeAuthority(freeze_authority).
		SetMintAccount(mint).
		SetSysVarRentPubkeyAccount(base.SysVarRentPubkey)
}
