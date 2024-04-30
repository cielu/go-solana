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
	"github.com/cielu/go-solana/common"
	"github.com/cielu/go-solana/pkg/encodbin"
	"github.com/cielu/go-solana/types/base"
)

// Initializes a new account to hold tokens.  If this account is associated
// with the native mint then the token balance of the initialized account
// will be equal to the amount of SOL in the account. If this account is
// associated with another mint, that mint must be initialized before this
// command can succeed.
//
// The `InitializeAccount` instruction requires no signers and MUST be
// included within the same Transaction as the system program's
// `CreateAccount` instruction that creates the account being initialized.
// Otherwise another party can acquire ownership of the uninitialized
// account.
type InitializeAccount struct {

	// [0] = [WRITE] account
	// ··········· The account to initialize.
	//
	// [1] = [] mint
	// ··········· The mint this account will be associated with.
	//
	// [2] = [] owner
	// ··········· The new account's owner/multisignature.
	//
	// [3] = [] $(SysVarRentPubkey)
	// ··········· Rent sysvar.
	AccountMeta []*base.AccountMeta `bin:"-" borsh_skip:"true"`
}

// NewInitializeAccountInstructionBuilder creates a new `InitializeAccount` instruction builder.
func NewInitializeAccountInstructionBuilder() *InitializeAccount {
	nd := &InitializeAccount{
		AccountMeta: make([]*base.AccountMeta, 4),
	}
	nd.AccountMeta[3] = base.Meta(base.SysVarRentPubkey)
	return nd
}

// SetAccount sets the "account" account.
// The account to initialize.
func (inst *InitializeAccount) SetAccount(account common.Address) *InitializeAccount {
	inst.AccountMeta[0] = base.Meta(account).WRITE()
	return inst
}

// GetAccount gets the "account" account.
// The account to initialize.
func (inst *InitializeAccount) GetAccount() *base.AccountMeta {
	return inst.AccountMeta[0]
}

// SetMintAccount sets the "mint" account.
// The mint this account will be associated with.
func (inst *InitializeAccount) SetMintAccount(mint common.Address) *InitializeAccount {
	inst.AccountMeta[1] = base.Meta(mint)
	return inst
}

// GetMintAccount gets the "mint" account.
// The mint this account will be associated with.
func (inst *InitializeAccount) GetMintAccount() *base.AccountMeta {
	return inst.AccountMeta[1]
}

// SetOwnerAccount sets the "owner" account.
// The new account's owner/multisignature.
func (inst *InitializeAccount) SetOwnerAccount(owner common.Address) *InitializeAccount {
	inst.AccountMeta[2] = base.Meta(owner)
	return inst
}

// GetOwnerAccount gets the "owner" account.
// The new account's owner/multisignature.
func (inst *InitializeAccount) GetOwnerAccount() *base.AccountMeta {
	return inst.AccountMeta[2]
}

// SetSysVarRentPubkeyAccount sets the "$(SysVarRentPubkey)" account.
// Rent sysvar.
func (inst *InitializeAccount) SetSysVarRentPubkeyAccount(SysVarRentPubkey common.Address) *InitializeAccount {
	inst.AccountMeta[3] = base.Meta(SysVarRentPubkey)
	return inst
}

// GetSysVarRentPubkeyAccount gets the "$(SysVarRentPubkey)" account.
// Rent sysvar.
func (inst *InitializeAccount) GetSysVarRentPubkeyAccount() *base.AccountMeta {
	return inst.AccountMeta[3]
}

func (inst InitializeAccount) Build() *Instruction {
	return &Instruction{BaseVariant: encodbin.BaseVariant{
		Impl:   inst,
		TypeID: encodbin.TypeIDFromUint8(Instruction_InitializeAccount),
	}}
}

// ValidateAndBuild validates the instruction parameters and accounts;
// if there is a validation error, it returns the error.
// Otherwise, it builds and returns the instruction.
func (inst InitializeAccount) ValidateAndBuild() (*Instruction, error) {
	if err := inst.Validate(); err != nil {
		return nil, err
	}
	return inst.Build(), nil
}

func (inst *InitializeAccount) Validate() error {
	// Check whether all (required) accounts are set:
	{
		if inst.AccountMeta[0] == nil {
			return errors.New("accounts.Account is not set")
		}
		if inst.AccountMeta[1] == nil {
			return errors.New("accounts.Mint is not set")
		}
		if inst.AccountMeta[2] == nil {
			return errors.New("accounts.Owner is not set")
		}
		if inst.AccountMeta[3] == nil {
			return errors.New("accounts.SysVarRentPubkey is not set")
		}
	}
	return nil
}

func (obj InitializeAccount) MarshalWithEncoder(encoder *encodbin.Encoder) (err error) {
	return nil
}

// NewInitializeAccountInstruction declares a new InitializeAccount instruction with the provided parameters and accounts.
func NewInitializeAccountInstruction(
	// Accounts:
	account common.Address,
	mint common.Address,
	owner common.Address) *InitializeAccount {
	return NewInitializeAccountInstructionBuilder().
		SetAccount(account).
		SetMintAccount(mint).
		SetOwnerAccount(owner).
		SetSysVarRentPubkeyAccount(base.SysVarRentPubkey)
}
