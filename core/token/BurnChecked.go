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

// BurnChecked Burns tokens by removing them from an account.  `BurnChecked` does not
// support accounts associated with the native mint, use `CloseAccount`
// instead.
//
// This instruction differs from Burn in that the decimals value is checked
// by the caller. This may be useful when creating transactions offline or
// within a hardware wallet.
type BurnChecked struct {
	// The amount of tokens to burn.
	Amount *uint64

	// Expected number of base 10 digits to the right of the decimal place.
	Decimals *uint8

	// [0] = [WRITE] source
	// ··········· The account to burn from.
	//
	// [1] = [WRITE] mint
	// ··········· The token mint.
	//
	// [2] = [] owner
	// ··········· The account's owner/delegate.
	//
	// [3...] = [SIGNER] signers
	// ··········· M signer accounts.
	Accounts []*solana.AccountMeta `bin:"-" borsh_skip:"true"`
	Signers  []*solana.AccountMeta `bin:"-" borsh_skip:"true"`
}

func (brCkd *BurnChecked) SetAccounts(accounts []*solana.AccountMeta) error {
	brCkd.Accounts, brCkd.Signers = library.SliceSplitFrom(accounts, 3)
	return nil
}

func (brCkd BurnChecked) GetAccounts() (accounts []*solana.AccountMeta) {
	accounts = append(accounts, brCkd.Accounts...)
	accounts = append(accounts, brCkd.Signers...)
	return
}

// NewBurnCheckedInstructionBuilder creates a new `BurnChecked` instruction builder.
func NewBurnCheckedInstructionBuilder() *BurnChecked {
	nd := &BurnChecked{
		Accounts: make([]*solana.AccountMeta, 3),
		Signers:  make([]*solana.AccountMeta, 0),
	}
	return nd
}

// SetAmount sets the "amount" parameter.
// The amount of tokens to burn.
func (brCkd *BurnChecked) SetAmount(amount uint64) *BurnChecked {
	brCkd.Amount = &amount
	return brCkd
}

// SetDecimals sets the "decimals" parameter.
// Expected number of base 10 digits to the right of the decimal place.
func (brCkd *BurnChecked) SetDecimals(decimals uint8) *BurnChecked {
	brCkd.Decimals = &decimals
	return brCkd
}

// SetSourceAccount sets the "source" account.
// The account to burn from.
func (brCkd *BurnChecked) SetSourceAccount(source solana.PublicKey) *BurnChecked {
	brCkd.Accounts[0] = solana.Meta(source).WRITE()
	return brCkd
}

// GetSourceAccount gets the "source" account.
// The account to burn from.
func (brCkd *BurnChecked) GetSourceAccount() *solana.AccountMeta {
	return brCkd.Accounts[0]
}

// SetMintAccount sets the "mint" account.
// The token mint.
func (brCkd *BurnChecked) SetMintAccount(mint solana.PublicKey) *BurnChecked {
	brCkd.Accounts[1] = solana.Meta(mint).WRITE()
	return brCkd
}

// GetMintAccount gets the "mint" account.
// The token mint.
func (brCkd *BurnChecked) GetMintAccount() *solana.AccountMeta {
	return brCkd.Accounts[1]
}

// SetOwnerAccount sets the "owner" account.
// The account's owner/delegate.
func (brCkd *BurnChecked) SetOwnerAccount(owner solana.PublicKey, multisigSigners ...solana.PublicKey) *BurnChecked {
	brCkd.Accounts[2] = solana.Meta(owner)
	if len(multisigSigners) == 0 {
		brCkd.Accounts[2].SIGNER()
	}
	for _, signer := range multisigSigners {
		brCkd.Signers = append(brCkd.Signers, solana.Meta(signer).SIGNER())
	}
	return brCkd
}

// GetOwnerAccount gets the "owner" account.
// The account's owner/delegate.
func (brCkd *BurnChecked) GetOwnerAccount() *solana.AccountMeta {
	return brCkd.Accounts[2]
}

func (brCkd BurnChecked) Build() *Instruction {
	return &Instruction{BaseVariant: encodbin.BaseVariant{
		Impl:   brCkd,
		TypeID: encodbin.TypeIDFromUint8(Instruction_BurnChecked),
	}}
}

// ValidateAndBuild validates the instruction parameters and accounts;
// if there is a validation error, it returns the error.
// Otherwise, it builds and returns the instruction.
func (brCkd BurnChecked) ValidateAndBuild() (*Instruction, error) {
	if err := brCkd.Validate(); err != nil {
		return nil, err
	}
	return brCkd.Build(), nil
}

func (brCkd *BurnChecked) Validate() error {
	// Check whether all (required) parameters are set:
	{
		if brCkd.Amount == nil {
			return errors.New("Amount parameter is not set")
		}
		if brCkd.Decimals == nil {
			return errors.New("Decimals parameter is not set")
		}
	}

	// Check whether all (required) accounts are set:
	{
		if brCkd.Accounts[0] == nil {
			return errors.New("accounts.Source is not set")
		}
		if brCkd.Accounts[1] == nil {
			return errors.New("accounts.Mint is not set")
		}
		if brCkd.Accounts[2] == nil {
			return errors.New("accounts.Owner is not set")
		}
		if !brCkd.Accounts[2].IsSigner && len(brCkd.Signers) == 0 {
			return fmt.Errorf("accounts.Signers is not set")
		}
		if len(brCkd.Signers) > MAX_SIGNERS {
			return fmt.Errorf("too many signers; got %v, but max is 11", len(brCkd.Signers))
		}
	}
	return nil
}

func (brCkd BurnChecked) MarshalWithEncoder(encoder *encodbin.Encoder) (err error) {
	// Serialize `Amount` param:
	err = encoder.Encode(brCkd.Amount)
	if err != nil {
		return err
	}
	// Serialize `Decimals` param:
	err = encoder.Encode(brCkd.Decimals)
	if err != nil {
		return err
	}
	return nil
}

// NewBurnCheckedInstruction declares a new BurnChecked instruction with the provided parameters and accounts.
func NewBurnCheckedInstruction(
	// Parameters:
	amount uint64,
	decimals uint8,
	// Accounts:
	source solana.PublicKey,
	mint solana.PublicKey,
	owner solana.PublicKey,
	multisigSigners []solana.PublicKey,
) *BurnChecked {
	return NewBurnCheckedInstructionBuilder().
		SetAmount(amount).
		SetDecimals(decimals).
		SetSourceAccount(source).
		SetMintAccount(mint).
		SetOwnerAccount(owner, multisigSigners...)
}
