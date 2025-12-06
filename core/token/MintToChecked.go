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
	"github.com/cielu/go-solana/library"
	"github.com/cielu/go-solana/pkg/encodbin"
)

// MintToChecked Mints new tokens to an account.  The native mint does not support minting.
//
// This instruction differs from MintTo in that the decimals value is
// checked by the caller.  This may be useful when creating transactions
// offline or within a hardware wallet.
type MintToChecked struct {
	// The amount of new tokens to mint.
	Amount *uint64

	// Expected number of base 10 digits to the right of the decimal place.
	Decimals *uint8

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

func (mCkd *MintToChecked) SetAccounts(accounts []*solana.AccountMeta) error {
	mCkd.Accounts, mCkd.Signers = library.SliceSplitFrom(accounts, 3)
	return nil
}

func (mCkd MintToChecked) GetAccounts() (accounts []*solana.AccountMeta) {
	accounts = append(accounts, mCkd.Accounts...)
	accounts = append(accounts, mCkd.Signers...)
	return
}

// NewMintToCheckedInstructionBuilder creates a new `MintToChecked` instruction builder.
func NewMintToCheckedInstructionBuilder() *MintToChecked {
	nd := &MintToChecked{
		Accounts: make([]*solana.AccountMeta, 3),
		Signers:  make([]*solana.AccountMeta, 0),
	}
	return nd
}

// SetAmount sets the "amount" parameter.
// The amount of new tokens to mint.
func (mCkd *MintToChecked) SetAmount(amount uint64) *MintToChecked {
	mCkd.Amount = &amount
	return mCkd
}

// SetDecimals sets the "decimals" parameter.
// Expected number of base 10 digits to the right of the decimal place.
func (mCkd *MintToChecked) SetDecimals(decimals uint8) *MintToChecked {
	mCkd.Decimals = &decimals
	return mCkd
}

// SetMintAccount sets the "mint" account.
// The mint.
func (mCkd *MintToChecked) SetMintAccount(mint solana.PublicKey) *MintToChecked {
	mCkd.Accounts[0] = solana.Meta(mint).WRITE()
	return mCkd
}

// GetMintAccount gets the "mint" account.
// The mint.
func (mCkd *MintToChecked) GetMintAccount() *solana.AccountMeta {
	return mCkd.Accounts[0]
}

// SetDestinationAccount sets the "destination" account.
// The account to mint tokens to.
func (mCkd *MintToChecked) SetDestinationAccount(destination solana.PublicKey) *MintToChecked {
	mCkd.Accounts[1] = solana.Meta(destination).WRITE()
	return mCkd
}

// GetDestinationAccount gets the "destination" account.
// The account to mint tokens to.
func (mCkd *MintToChecked) GetDestinationAccount() *solana.AccountMeta {
	return mCkd.Accounts[1]
}

// SetAuthorityAccount sets the "authority" account.
// The mint's minting authority.
func (mCkd *MintToChecked) SetAuthorityAccount(authority solana.PublicKey, multisigSigners ...solana.PublicKey) *MintToChecked {
	mCkd.Accounts[2] = solana.Meta(authority)
	if len(multisigSigners) == 0 {
		mCkd.Accounts[2].SIGNER()
	}
	for _, signer := range multisigSigners {
		mCkd.Signers = append(mCkd.Signers, solana.Meta(signer).SIGNER())
	}
	return mCkd
}

// GetAuthorityAccount gets the "authority" account.
// The mint's minting authority.
func (mCkd *MintToChecked) GetAuthorityAccount() *solana.AccountMeta {
	return mCkd.Accounts[2]
}

func (mCkd MintToChecked) Build() *Instruction {
	return &Instruction{BaseVariant: encodbin.BaseVariant{
		Impl:   mCkd,
		TypeID: encodbin.TypeIDFromUint8(Instruction_MintToChecked),
	}}
}

// ValidateAndBuild validates the instruction parameters and accounts;
// if there is a validation error, it returns the error.
// Otherwise, it builds and returns the instruction.
func (mCkd MintToChecked) ValidateAndBuild() (*Instruction, error) {
	if err := mCkd.Validate(); err != nil {
		return nil, err
	}
	return mCkd.Build(), nil
}

func (mCkd *MintToChecked) Validate() error {
	// Check whether all (required) parameters are set:
	{
		if mCkd.Amount == nil {
			return errors.New("Amount parameter is not set")
		}
		if mCkd.Decimals == nil {
			return errors.New("Decimals parameter is not set")
		}
	}

	// Check whether all (required) accounts are set:
	{
		if mCkd.Accounts[0] == nil {
			return errors.New("accounts.Mint is not set")
		}
		if mCkd.Accounts[1] == nil {
			return errors.New("accounts.Destination is not set")
		}
		if mCkd.Accounts[2] == nil {
			return errors.New("accounts.Authority is not set")
		}
		if !mCkd.Accounts[2].IsSigner && len(mCkd.Signers) == 0 {
			return fmt.Errorf("accounts.Signers is not set")
		}
		if len(mCkd.Signers) > MAX_SIGNERS {
			return fmt.Errorf("too many signers; got %v, but max is 11", len(mCkd.Signers))
		}
	}
	return nil
}

func (mCkd MintToChecked) MarshalWithEncoder(encoder *encodbin.Encoder) (err error) {
	// Serialize `Amount` param:
	err = encoder.Encode(mCkd.Amount)
	if err != nil {
		return err
	}
	// Serialize `Decimals` param:
	err = encoder.Encode(mCkd.Decimals)
	if err != nil {
		return err
	}
	return nil
}

// NewMintToCheckedInstruction declares a new MintToChecked instruction with the provided parameters and accounts.
func NewMintToCheckedInstruction(
	// Parameters:
	amount uint64,
	decimals uint8,
	// Accounts:
	mint solana.PublicKey,
	destination solana.PublicKey,
	authority solana.PublicKey,
	multisigSigners []solana.PublicKey,
) *MintToChecked {
	return NewMintToCheckedInstructionBuilder().
		SetAmount(amount).
		SetDecimals(decimals).
		SetMintAccount(mint).
		SetDestinationAccount(destination).
		SetAuthorityAccount(authority, multisigSigners...)
}
