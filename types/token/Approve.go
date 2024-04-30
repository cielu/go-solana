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
	"github.com/cielu/go-solana/common"
	"github.com/cielu/go-solana/core"
	"github.com/cielu/go-solana/pkg/encodbin"
	"github.com/cielu/go-solana/types/base"
)

// Approve A delegate is given the authority over tokens on
// behalf of the source account's owner.
type Approve struct {
	// The amount of tokens the delegate is approved for.
	Amount *uint64

	// [0] = [WRITE] source
	// ··········· The source account.
	//
	// [1] = [] delegate
	// ··········· The delegate.
	//
	// [2] = [] owner
	// ··········· The source account owner.
	//
	// [3...] = [SIGNER] signers
	// ··········· M signer accounts.
	Accounts []*base.AccountMeta `bin:"-" borsh_skip:"true"`
	Signers  []*base.AccountMeta `bin:"-" borsh_skip:"true"`
}

func (appr *Approve) SetAccounts(accounts []*base.AccountMeta) error {
	appr.Accounts, appr.Signers = core.SliceSplitFrom(accounts, 3)
	return nil
}

func (appr Approve) GetAccounts() (accounts []*base.AccountMeta) {
	accounts = append(accounts, appr.Accounts...)
	accounts = append(accounts, appr.Signers...)
	return
}

// NewApproveInstructionBuilder creates a new `Approve` instruction builder.
func NewApproveInstructionBuilder() *Approve {
	nd := &Approve{
		Accounts: make([]*base.AccountMeta, 3),
		Signers:  make([]*base.AccountMeta, 0),
	}
	return nd
}

// SetAmount sets the "amount" parameter.
// The amount of tokens the delegate is approved for.
func (appr *Approve) SetAmount(amount uint64) *Approve {
	appr.Amount = &amount
	return appr
}

// SetSourceAccount sets the "source" account.
// The source account.
func (appr *Approve) SetSourceAccount(source common.Address) *Approve {
	appr.Accounts[0] = base.Meta(source).WRITE()
	return appr
}

// GetSourceAccount gets the "source" account.
// The source account.
func (appr *Approve) GetSourceAccount() *base.AccountMeta {
	return appr.Accounts[0]
}

// SetDelegateAccount sets the "delegate" account.
// The delegate.
func (appr *Approve) SetDelegateAccount(delegate common.Address) *Approve {
	appr.Accounts[1] = base.Meta(delegate)
	return appr
}

// GetDelegateAccount gets the "delegate" account.
// The delegate.
func (appr *Approve) GetDelegateAccount() *base.AccountMeta {
	return appr.Accounts[1]
}

// SetOwnerAccount sets the "owner" account.
// The source account owner.
func (appr *Approve) SetOwnerAccount(owner common.Address, multisigSigners ...common.Address) *Approve {
	appr.Accounts[2] = base.Meta(owner)
	if len(multisigSigners) == 0 {
		appr.Accounts[2].SIGNER()
	}
	for _, signer := range multisigSigners {
		appr.Signers = append(appr.Signers, base.Meta(signer).SIGNER())
	}
	return appr
}

// GetOwnerAccount gets the "owner" account.
// The source account owner.
func (appr *Approve) GetOwnerAccount() *base.AccountMeta {
	return appr.Accounts[2]
}

func (appr Approve) Build() *Instruction {
	return &Instruction{BaseVariant: encodbin.BaseVariant{
		Impl:   appr,
		TypeID: encodbin.TypeIDFromUint8(Instruction_Approve),
	}}
}

// ValidateAndBuild validates the instruction parameters and accounts;
// if there is a validation error, it returns the error.
// Otherwise, it builds and returns the instruction.
func (appr Approve) ValidateAndBuild() (*Instruction, error) {
	if err := appr.Validate(); err != nil {
		return nil, err
	}
	return appr.Build(), nil
}

func (appr *Approve) Validate() error {
	// Check whether all (required) parameters are set:
	{
		if appr.Amount == nil {
			return errors.New("Amount parameter is not set")
		}
	}

	// Check whether all (required) accounts are set:
	{
		if appr.Accounts[0] == nil {
			return errors.New("accounts.Source is not set")
		}
		if appr.Accounts[1] == nil {
			return errors.New("accounts.Delegate is not set")
		}
		if appr.Accounts[2] == nil {
			return errors.New("accounts.Owner is not set")
		}
		if !appr.Accounts[2].IsSigner && len(appr.Signers) == 0 {
			return fmt.Errorf("accounts.Signers is not set")
		}
		if len(appr.Signers) > MAX_SIGNERS {
			return fmt.Errorf("too many signers; got %v, but max is 11", len(appr.Signers))
		}
	}
	return nil
}

func (appr Approve) MarshalWithEncoder(encoder *encodbin.Encoder) (err error) {
	// Serialize `Amount` param:
	err = encoder.Encode(appr.Amount)
	if err != nil {
		return err
	}
	return nil
}

// NewApproveInstruction declares a new Approve instruction with the provided parameters and accounts.
func NewApproveInstruction(
	// Parameters:
	amount uint64,
	// Accounts:
	source common.Address,
	delegate common.Address,
	owner common.Address,
	multisigSigners []common.Address,
) *Approve {
	return NewApproveInstructionBuilder().
		SetAmount(amount).
		SetSourceAccount(source).
		SetDelegateAccount(delegate).
		SetOwnerAccount(owner, multisigSigners...)
}
