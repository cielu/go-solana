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

// CloseAccount Close an account by transferring all its SOL to the destination account.
// Non-native accounts may only be closed if its token amount is zero.
type CloseAccount struct {

	// [0] = [WRITE] account
	// ··········· The account to close.
	//
	// [1] = [WRITE] destination
	// ··········· The destination account.
	//
	// [2] = [] owner
	// ··········· The account's owner.
	//
	// [3...] = [SIGNER] signers
	// ··········· M signer accounts.
	Accounts []*base.AccountMeta `bin:"-" borsh_skip:"true"`
	Signers  []*base.AccountMeta `bin:"-" borsh_skip:"true"`
}

func (cloAcc *CloseAccount) SetAccounts(accounts []*base.AccountMeta) error {
	cloAcc.Accounts, cloAcc.Signers = core.SliceSplitFrom(accounts, 3)
	return nil
}

func (cloAcc CloseAccount) GetAccounts() (accounts []*base.AccountMeta) {
	accounts = append(accounts, cloAcc.Accounts...)
	accounts = append(accounts, cloAcc.Signers...)
	return
}

// NewCloseAccountInstructionBuilder creates a new `CloseAccount` instruction builder.
func NewCloseAccountInstructionBuilder() *CloseAccount {
	nd := &CloseAccount{
		Accounts: make([]*base.AccountMeta, 3),
		Signers:  make([]*base.AccountMeta, 0),
	}
	return nd
}

// SetAccount sets the "account" account.
// The account to close.
func (cloAcc *CloseAccount) SetAccount(account common.Address) *CloseAccount {
	cloAcc.Accounts[0] = base.Meta(account).WRITE()
	return cloAcc
}

// GetAccount gets the "account" account.
// The account to close.
func (cloAcc *CloseAccount) GetAccount() *base.AccountMeta {
	return cloAcc.Accounts[0]
}

// SetDestinationAccount sets the "destination" account.
// The destination account.
func (cloAcc *CloseAccount) SetDestinationAccount(destination common.Address) *CloseAccount {
	cloAcc.Accounts[1] = base.Meta(destination).WRITE()
	return cloAcc
}

// GetDestinationAccount gets the "destination" account.
// The destination account.
func (cloAcc *CloseAccount) GetDestinationAccount() *base.AccountMeta {
	return cloAcc.Accounts[1]
}

// SetOwnerAccount sets the "owner" account.
// The account's owner.
func (cloAcc *CloseAccount) SetOwnerAccount(owner common.Address, multisigSigners ...common.Address) *CloseAccount {
	cloAcc.Accounts[2] = base.Meta(owner)
	if len(multisigSigners) == 0 {
		cloAcc.Accounts[2].SIGNER()
	}
	for _, signer := range multisigSigners {
		cloAcc.Signers = append(cloAcc.Signers, base.Meta(signer).SIGNER())
	}
	return cloAcc
}

// GetOwnerAccount gets the "owner" account.
// The account's owner.
func (cloAcc *CloseAccount) GetOwnerAccount() *base.AccountMeta {
	return cloAcc.Accounts[2]
}

func (cloAcc CloseAccount) Build() *Instruction {
	return &Instruction{BaseVariant: encodbin.BaseVariant{
		Impl:   cloAcc,
		TypeID: encodbin.TypeIDFromUint8(Instruction_CloseAccount),
	}}
}

// ValidateAndBuild validates the instruction parameters and accounts;
// if there is a validation error, it returns the error.
// Otherwise, it builds and returns the instruction.
func (cloAcc CloseAccount) ValidateAndBuild() (*Instruction, error) {
	if err := cloAcc.Validate(); err != nil {
		return nil, err
	}
	return cloAcc.Build(), nil
}

func (cloAcc *CloseAccount) Validate() error {
	// Check whether all (required) accounts are set:
	{
		if cloAcc.Accounts[0] == nil {
			return errors.New("accounts.Account is not set")
		}
		if cloAcc.Accounts[1] == nil {
			return errors.New("accounts.Destination is not set")
		}
		if cloAcc.Accounts[2] == nil {
			return errors.New("accounts.Owner is not set")
		}
		if !cloAcc.Accounts[2].IsSigner && len(cloAcc.Signers) == 0 {
			return fmt.Errorf("accounts.Signers is not set")
		}
		if len(cloAcc.Signers) > MAX_SIGNERS {
			return fmt.Errorf("too many signers; got %v, but max is 11", len(cloAcc.Signers))
		}
	}
	return nil
}

func (cloAcc CloseAccount) MarshalWithEncoder(encoder *encodbin.Encoder) (err error) {
	return nil
}

// NewCloseAccountInstruction declares a new CloseAccount instruction with the provided parameters and accounts.
func NewCloseAccountInstruction(
	// Accounts:
	account common.Address,
	destination common.Address,
	owner common.Address,
	multisigSigners []common.Address,
) *CloseAccount {
	return NewCloseAccountInstructionBuilder().
		SetAccount(account).
		SetDestinationAccount(destination).
		SetOwnerAccount(owner, multisigSigners...)
}
