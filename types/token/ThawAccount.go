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

// ThawAccount Thaw a Frozen account using the Mint's freeze_authority (if set).
type ThawAccount struct {

	// [0] = [WRITE] account
	// ··········· The account to thaw.
	//
	// [1] = [] mint
	// ··········· The token mint.
	//
	// [2] = [] authority
	// ··········· The mint freeze authority.
	//
	// [3...] = [SIGNER] signers
	// ··········· M signer accounts.
	Accounts []*base.AccountMeta `bin:"-" borsh_skip:"true"`
	Signers  []*base.AccountMeta `bin:"-" borsh_skip:"true"`
}

func (tAcc *ThawAccount) SetAccounts(accounts []*base.AccountMeta) error {
	tAcc.Accounts, tAcc.Signers = core.SliceSplitFrom(accounts, 3)
	return nil
}

func (tAcc ThawAccount) GetAccounts() (accounts []*base.AccountMeta) {
	accounts = append(accounts, tAcc.Accounts...)
	accounts = append(accounts, tAcc.Signers...)
	return
}

// NewThawAccountInstructionBuilder creates a new `ThawAccount` instruction builder.
func NewThawAccountInstructionBuilder() *ThawAccount {
	nd := &ThawAccount{
		Accounts: make([]*base.AccountMeta, 3),
		Signers:  make([]*base.AccountMeta, 0),
	}
	return nd
}

// SetAccount sets the "account" account.
// The account to thaw.
func (tAcc *ThawAccount) SetAccount(account common.Address) *ThawAccount {
	tAcc.Accounts[0] = base.Meta(account).WRITE()
	return tAcc
}

// GetAccount gets the "account" account.
// The account to thaw.
func (tAcc *ThawAccount) GetAccount() *base.AccountMeta {
	return tAcc.Accounts[0]
}

// SetMintAccount sets the "mint" account.
// The token mint.
func (tAcc *ThawAccount) SetMintAccount(mint common.Address) *ThawAccount {
	tAcc.Accounts[1] = base.Meta(mint)
	return tAcc
}

// GetMintAccount gets the "mint" account.
// The token mint.
func (tAcc *ThawAccount) GetMintAccount() *base.AccountMeta {
	return tAcc.Accounts[1]
}

// SetAuthorityAccount sets the "authority" account.
// The mint freeze authority.
func (tAcc *ThawAccount) SetAuthorityAccount(authority common.Address, multisigSigners ...common.Address) *ThawAccount {
	tAcc.Accounts[2] = base.Meta(authority)
	if len(multisigSigners) == 0 {
		tAcc.Accounts[2].SIGNER()
	}
	for _, signer := range multisigSigners {
		tAcc.Signers = append(tAcc.Signers, base.Meta(signer).SIGNER())
	}
	return tAcc
}

// GetAuthorityAccount gets the "authority" account.
// The mint freeze authority.
func (tAcc *ThawAccount) GetAuthorityAccount() *base.AccountMeta {
	return tAcc.Accounts[2]
}

func (tAcc ThawAccount) Build() *Instruction {
	return &Instruction{BaseVariant: encodbin.BaseVariant{
		Impl:   tAcc,
		TypeID: encodbin.TypeIDFromUint8(Instruction_ThawAccount),
	}}
}

// ValidateAndBuild validates the instruction parameters and accounts;
// if there is a validation error, it returns the error.
// Otherwise, it builds and returns the instruction.
func (tAcc ThawAccount) ValidateAndBuild() (*Instruction, error) {
	if err := tAcc.Validate(); err != nil {
		return nil, err
	}
	return tAcc.Build(), nil
}

func (tAcc *ThawAccount) Validate() error {
	// Check whether all (required) accounts are set:
	{
		if tAcc.Accounts[0] == nil {
			return errors.New("accounts.Account is not set")
		}
		if tAcc.Accounts[1] == nil {
			return errors.New("accounts.Mint is not set")
		}
		if tAcc.Accounts[2] == nil {
			return errors.New("accounts.Authority is not set")
		}
		if !tAcc.Accounts[2].IsSigner && len(tAcc.Signers) == 0 {
			return fmt.Errorf("accounts.Signers is not set")
		}
		if len(tAcc.Signers) > MAX_SIGNERS {
			return fmt.Errorf("too many signers; got %v, but max is 11", len(tAcc.Signers))
		}
	}
	return nil
}

func (tAcc ThawAccount) MarshalWithEncoder(encoder *encodbin.Encoder) (err error) {
	return nil
}

// NewThawAccountInstruction declares a new ThawAccount instruction with the provided parameters and accounts.
func NewThawAccountInstruction(
	// Accounts:
	account common.Address,
	mint common.Address,
	authority common.Address,
	multisigSigners []common.Address,
) *ThawAccount {
	return NewThawAccountInstructionBuilder().
		SetAccount(account).
		SetMintAccount(mint).
		SetAuthorityAccount(authority, multisigSigners...)
}
