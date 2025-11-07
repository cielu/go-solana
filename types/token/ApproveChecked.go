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

// ApproveChecked  A delegate is given the authority over tokens on
// behalf of the source account's owner.
//
// This instruction differs from Approve in that the token mint and
// decimals value is checked by the caller.  This may be useful when
// creating transactions offline or within a hardware wallet.
type ApproveChecked struct {
	// The amount of tokens the delegate is approved for.
	Amount *uint64

	// Expected number of base 10 digits to the right of the decimal place.
	Decimals *uint8

	// [0] = [WRITE] source
	// ··········· The source account.
	//
	// [1] = [] mint
	// ··········· The token mint.
	//
	// [2] = [] delegate
	// ··········· The delegate.
	//
	// [3] = [] owner
	// ··········· The source account owner.
	//
	// [4...] = [SIGNER] signers
	// ··········· M signer accounts.
	Accounts []*solana.AccountMeta `bin:"-" borsh_skip:"true"`
	Signers  []*solana.AccountMeta `bin:"-" borsh_skip:"true"`
}

func (apprCkd *ApproveChecked) SetAccounts(accounts []*solana.AccountMeta) error {
	apprCkd.Accounts, apprCkd.Signers = core.SliceSplitFrom(accounts, 4)
	return nil
}

func (apprCkd ApproveChecked) GetAccounts() (accounts []*solana.AccountMeta) {
	accounts = append(accounts, apprCkd.Accounts...)
	accounts = append(accounts, apprCkd.Signers...)
	return
}

// NewApproveCheckedInstructionBuilder creates a new `ApproveChecked` instruction builder.
func NewApproveCheckedInstructionBuilder() *ApproveChecked {
	nd := &ApproveChecked{
		Accounts: make([]*solana.AccountMeta, 4),
		Signers:  make([]*solana.AccountMeta, 0),
	}
	return nd
}

// SetAmount sets the "amount" parameter.
// The amount of tokens the delegate is approved for.
func (apprCkd *ApproveChecked) SetAmount(amount uint64) *ApproveChecked {
	apprCkd.Amount = &amount
	return apprCkd
}

// SetDecimals sets the "decimals" parameter.
// Expected number of base 10 digits to the right of the decimal place.
func (apprCkd *ApproveChecked) SetDecimals(decimals uint8) *ApproveChecked {
	apprCkd.Decimals = &decimals
	return apprCkd
}

// SetSourceAccount sets the "source" account.
// The source account.
func (apprCkd *ApproveChecked) SetSourceAccount(source solana.PublicKey) *ApproveChecked {
	apprCkd.Accounts[0] = solana.Meta(source).WRITE()
	return apprCkd
}

// GetSourceAccount gets the "source" account.
// The source account.
func (apprCkd *ApproveChecked) GetSourceAccount() *solana.AccountMeta {
	return apprCkd.Accounts[0]
}

// SetMintAccount sets the "mint" account.
// The token mint.
func (apprCkd *ApproveChecked) SetMintAccount(mint solana.PublicKey) *ApproveChecked {
	apprCkd.Accounts[1] = solana.Meta(mint)
	return apprCkd
}

// GetMintAccount gets the "mint" account.
// The token mint.
func (apprCkd *ApproveChecked) GetMintAccount() *solana.AccountMeta {
	return apprCkd.Accounts[1]
}

// SetDelegateAccount sets the "delegate" account.
// The delegate.
func (apprCkd *ApproveChecked) SetDelegateAccount(delegate solana.PublicKey) *ApproveChecked {
	apprCkd.Accounts[2] = solana.Meta(delegate)
	return apprCkd
}

// GetDelegateAccount gets the "delegate" account.
// The delegate.
func (apprCkd *ApproveChecked) GetDelegateAccount() *solana.AccountMeta {
	return apprCkd.Accounts[2]
}

// SetOwnerAccount sets the "owner" account.
// The source account owner.
func (apprCkd *ApproveChecked) SetOwnerAccount(owner solana.PublicKey, multisigSigners ...solana.PublicKey) *ApproveChecked {
	apprCkd.Accounts[3] = solana.Meta(owner)
	if len(multisigSigners) == 0 {
		apprCkd.Accounts[3].SIGNER()
	}
	for _, signer := range multisigSigners {
		apprCkd.Signers = append(apprCkd.Signers, solana.Meta(signer).SIGNER())
	}
	return apprCkd
}

// GetOwnerAccount gets the "owner" account.
// The source account owner.
func (apprCkd *ApproveChecked) GetOwnerAccount() *solana.AccountMeta {
	return apprCkd.Accounts[3]
}

func (apprCkd ApproveChecked) Build() *Instruction {
	return &Instruction{BaseVariant: encodbin.BaseVariant{
		Impl:   apprCkd,
		TypeID: encodbin.TypeIDFromUint8(Instruction_ApproveChecked),
	}}
}

// ValidateAndBuild validates the instruction parameters and accounts;
// if there is a validation error, it returns the error.
// Otherwise, it builds and returns the instruction.
func (apprCkd ApproveChecked) ValidateAndBuild() (*Instruction, error) {
	if err := apprCkd.Validate(); err != nil {
		return nil, err
	}
	return apprCkd.Build(), nil
}

func (apprCkd *ApproveChecked) Validate() error {
	// Check whether all (required) parameters are set:
	{
		if apprCkd.Amount == nil {
			return errors.New("Amount parameter is not set")
		}
		if apprCkd.Decimals == nil {
			return errors.New("Decimals parameter is not set")
		}
	}

	// Check whether all (required) accounts are set:
	{
		if apprCkd.Accounts[0] == nil {
			return errors.New("accounts.Source is not set")
		}
		if apprCkd.Accounts[1] == nil {
			return errors.New("accounts.Mint is not set")
		}
		if apprCkd.Accounts[2] == nil {
			return errors.New("accounts.Delegate is not set")
		}
		if apprCkd.Accounts[3] == nil {
			return errors.New("accounts.Owner is not set")
		}
		if !apprCkd.Accounts[3].IsSigner && len(apprCkd.Signers) == 0 {
			return fmt.Errorf("accounts.Signers is not set")
		}
		if len(apprCkd.Signers) > MAX_SIGNERS {
			return fmt.Errorf("too many signers; got %v, but max is 11", len(apprCkd.Signers))
		}
	}
	return nil
}

func (apprCkd ApproveChecked) MarshalWithEncoder(encoder *encodbin.Encoder) (err error) {
	// Serialize `Amount` param:
	err = encoder.Encode(apprCkd.Amount)
	if err != nil {
		return err
	}
	// Serialize `Decimals` param:
	err = encoder.Encode(apprCkd.Decimals)
	if err != nil {
		return err
	}
	return nil
}

// NewApproveCheckedInstruction declares a new ApproveChecked instruction with the provided parameters and accounts.
func NewApproveCheckedInstruction(
	// Parameters:
	amount uint64,
	decimals uint8,
	// Accounts:
	source solana.PublicKey,
	mint solana.PublicKey,
	delegate solana.PublicKey,
	owner solana.PublicKey,
	multisigSigners []solana.PublicKey,
) *ApproveChecked {
	return NewApproveCheckedInstructionBuilder().
		SetAmount(amount).
		SetDecimals(decimals).
		SetSourceAccount(source).
		SetMintAccount(mint).
		SetDelegateAccount(delegate).
		SetOwnerAccount(owner, multisigSigners...)
}
