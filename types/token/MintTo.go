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
	"fmt"

	"github.com/cielu/go-solana"
	"github.com/cielu/go-solana/core"
	"github.com/cielu/go-solana/pkg/encodbin"
)

// MintTo Mints new tokens to an account.  The native mint does not support
// minting.
type MintTo struct {
	// The amount of new tokens to mint.
	Amount *uint64

	// [0] = [WRITE] mint
	// ··········· The mint.
	//
	// [1] = [WRITE] destination
	// ··········· The account to mint tokens to.
	//
	// [2] = [] authority
	// ··········· The mint's minting authority.
	//
	// [3...] = [SIGNER] signers
	// ··········· M signer accounts.
	Accounts []*solana.AccountMeta `bin:"-" borsh_skip:"true"`
	Signers  []*solana.AccountMeta `bin:"-" borsh_skip:"true"`
}

func (mto *MintTo) SetAccounts(accounts []*solana.AccountMeta) error {
	mto.Accounts, mto.Signers = core.SliceSplitFrom(accounts, 3)
	return nil
}

func (mto MintTo) GetAccounts() (accounts []*solana.AccountMeta) {
	accounts = append(accounts, mto.Accounts...)
	accounts = append(accounts, mto.Signers...)
	return
}

// NewMintToInstructionBuilder creates a new `MintTo` instruction builder.
func NewMintToInstructionBuilder() *MintTo {
	nd := &MintTo{
		Accounts: make([]*solana.AccountMeta, 3),
		Signers:  make([]*solana.AccountMeta, 0),
	}
	return nd
}

// SetAmount sets the "amount" parameter.
// The amount of new tokens to mint.
func (mto *MintTo) SetAmount(amount uint64) *MintTo {
	mto.Amount = &amount
	return mto
}

// SetMintAccount sets the "mint" account.
// The mint.
func (mto *MintTo) SetMintAccount(mint solana.PublicKey) *MintTo {
	mto.Accounts[0] = solana.Meta(mint).WRITE()
	return mto
}

// GetMintAccount gets the "mint" account.
// The mint.
func (mto *MintTo) GetMintAccount() *solana.AccountMeta {
	return mto.Accounts[0]
}

// SetDestinationAccount sets the "destination" account.
// The account to mint tokens to.
func (mto *MintTo) SetDestinationAccount(destination solana.PublicKey) *MintTo {
	mto.Accounts[1] = solana.Meta(destination).WRITE()
	return mto
}

// GetDestinationAccount gets the "destination" account.
// The account to mint tokens to.
func (mto *MintTo) GetDestinationAccount() *solana.AccountMeta {
	return mto.Accounts[1]
}

// SetAuthorityAccount sets the "authority" account.
// The mint's minting authority.
func (mto *MintTo) SetAuthorityAccount(authority solana.PublicKey, multisigSigners ...solana.PublicKey) *MintTo {
	mto.Accounts[2] = solana.Meta(authority)
	if len(multisigSigners) == 0 {
		mto.Accounts[2].SIGNER()
	}
	for _, signer := range multisigSigners {
		mto.Signers = append(mto.Signers, solana.Meta(signer).SIGNER())
	}
	return mto
}

// GetAuthorityAccount gets the "authority" account.
// The mint's minting authority.
func (mto *MintTo) GetAuthorityAccount() *solana.AccountMeta {
	return mto.Accounts[2]
}

func (mto MintTo) Build() *Instruction {
	return &Instruction{BaseVariant: encodbin.BaseVariant{
		Impl:   mto,
		TypeID: encodbin.TypeIDFromUint8(Instruction_MintTo),
	}}
}

// ValidateAndBuild validates the instruction parameters and accounts;
// if there is a validation error, it returns the error.
// Otherwise, it builds and returns the instruction.
func (mto MintTo) ValidateAndBuild() (*Instruction, error) {
	if err := mto.Validate(); err != nil {
		return nil, err
	}
	return mto.Build(), nil
}

func (mto *MintTo) Validate() error {
	// Check whether all (required) parameters are set:
	{
		if mto.Amount == nil {
			return errors.New("Amount parameter is not set")
		}
	}

	// Check whether all (required) accounts are set:
	{
		if mto.Accounts[0] == nil {
			return errors.New("accounts.Mint is not set")
		}
		if mto.Accounts[1] == nil {
			return errors.New("accounts.Destination is not set")
		}
		if mto.Accounts[2] == nil {
			return errors.New("accounts.Authority is not set")
		}
		if !mto.Accounts[2].IsSigner && len(mto.Signers) == 0 {
			return fmt.Errorf("accounts.Signers is not set")
		}
		if len(mto.Signers) > MAX_SIGNERS {
			return fmt.Errorf("too many signers; got %v, but max is 11", len(mto.Signers))
		}
	}
	return nil
}

func (mto MintTo) MarshalWithEncoder(encoder *encodbin.Encoder) (err error) {
	// Serialize `Amount` param:
	err = encoder.Encode(mto.Amount)
	if err != nil {
		return err
	}
	return nil
}

// NewMintToInstruction declares a new MintTo instruction with the provided parameters and accounts.
func NewMintToInstruction(
	// Parameters:
	amount uint64,
	// Accounts:
	mint solana.PublicKey,
	destination solana.PublicKey,
	authority solana.PublicKey,
	multisigSigners []solana.PublicKey,
) *MintTo {
	return NewMintToInstructionBuilder().
		SetAmount(amount).
		SetMintAccount(mint).
		SetDestinationAccount(destination).
		SetAuthorityAccount(authority, multisigSigners...)
}
