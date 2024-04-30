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

// Revoke the delegate's authority.
type Revoke struct {

	// [0] = [WRITE] source
	// ··········· The source account.
	//
	// [1] = [] owner
	// ··········· The source account's owner.
	//
	// [2...] = [SIGNER] signers
	// ··········· M signer accounts.
	Accounts []*base.AccountMeta `bin:"-" borsh_skip:"true"`
	Signers  []*base.AccountMeta `bin:"-" borsh_skip:"true"`
}

func (rvk *Revoke) SetAccounts(accounts []*base.AccountMeta) error {
	rvk.Accounts, rvk.Signers = core.SliceSplitFrom(accounts, 2)
	return nil
}

func (rvk Revoke) GetAccounts() (accounts []*base.AccountMeta) {
	accounts = append(accounts, rvk.Accounts...)
	accounts = append(accounts, rvk.Signers...)
	return
}

// NewRevokeInstructionBuilder creates a new `Revoke` instruction builder.
func NewRevokeInstructionBuilder() *Revoke {
	nd := &Revoke{
		Accounts: make([]*base.AccountMeta, 2),
		Signers:  make([]*base.AccountMeta, 0),
	}
	return nd
}

// SetSourceAccount sets the "source" account.
// The source account.
func (rvk *Revoke) SetSourceAccount(source common.Address) *Revoke {
	rvk.Accounts[0] = base.Meta(source).WRITE()
	return rvk
}

// GetSourceAccount gets the "source" account.
// The source account.
func (rvk *Revoke) GetSourceAccount() *base.AccountMeta {
	return rvk.Accounts[0]
}

// SetOwnerAccount sets the "owner" account.
// The source account's owner.
func (rvk *Revoke) SetOwnerAccount(owner common.Address, multisigSigners ...common.Address) *Revoke {
	rvk.Accounts[1] = base.Meta(owner)
	if len(multisigSigners) == 0 {
		rvk.Accounts[1].SIGNER()
	}
	for _, signer := range multisigSigners {
		rvk.Signers = append(rvk.Signers, base.Meta(signer).SIGNER())
	}
	return rvk
}

// GetOwnerAccount gets the "owner" account.
// The source account's owner.
func (rvk *Revoke) GetOwnerAccount() *base.AccountMeta {
	return rvk.Accounts[1]
}

func (rvk Revoke) Build() *Instruction {
	return &Instruction{BaseVariant: encodbin.BaseVariant{
		Impl:   rvk,
		TypeID: encodbin.TypeIDFromUint8(Instruction_Revoke),
	}}
}

// ValidateAndBuild validates the instruction parameters and accounts;
// if there is a validation error, it returns the error.
// Otherwise, it builds and returns the instruction.
func (rvk Revoke) ValidateAndBuild() (*Instruction, error) {
	if err := rvk.Validate(); err != nil {
		return nil, err
	}
	return rvk.Build(), nil
}

func (rvk *Revoke) Validate() error {
	// Check whether all (required) accounts are set:
	{
		if rvk.Accounts[0] == nil {
			return errors.New("accounts.Source is not set")
		}
		if rvk.Accounts[1] == nil {
			return errors.New("accounts.Owner is not set")
		}
		if !rvk.Accounts[1].IsSigner && len(rvk.Signers) == 0 {
			return fmt.Errorf("accounts.Signers is not set")
		}
		if len(rvk.Signers) > MAX_SIGNERS {
			return fmt.Errorf("too many signers; got %v, but max is 11", len(rvk.Signers))
		}
	}
	return nil
}

func (rvk Revoke) MarshalWithEncoder(encoder *encodbin.Encoder) (err error) {
	return nil
}

// NewRevokeInstruction declares a new Revoke instruction with the provided parameters and accounts.
func NewRevokeInstruction(
	// Accounts:
	source common.Address,
	owner common.Address,
	multisigSigners []common.Address,
) *Revoke {
	return NewRevokeInstructionBuilder().
		SetSourceAccount(source).
		SetOwnerAccount(owner, multisigSigners...)
}
