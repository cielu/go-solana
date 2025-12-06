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
	"github.com/cielu/go-solana/pkg/encodbin"
)

// InitializeAccount3 Like InitializeAccount2, but does not require the Rent sysvar to be provided.
type InitializeAccount3 struct {
	// The new account's owner/multisignature.
	Owner *solana.PublicKey

	// [0] = [WRITE] account
	// ··········· The account to initialize.
	//
	// [1] = [] mint
	// ··········· The mint this account will be associated with.
	solana.AccountMetaSlice `bin:"-" borsh_skip:"true"`
}

// NewInitializeAccount3InstructionBuilder creates a new `InitializeAccount3` instruction builder.
func NewInitializeAccount3InstructionBuilder() *InitializeAccount3 {
	nd := &InitializeAccount3{
		AccountMetaSlice: make([]*solana.AccountMeta, 2),
	}
	return nd
}

// SetOwner sets the "owner" parameter.
// The new account's owner/multisignature.
func (initAcc3 *InitializeAccount3) SetOwner(owner solana.PublicKey) *InitializeAccount3 {
	initAcc3.Owner = &owner
	return initAcc3
}

// SetAccount sets the "account" account.
// The account to initialize.
func (initAcc3 *InitializeAccount3) SetAccount(account solana.PublicKey) *InitializeAccount3 {
	initAcc3.AccountMetaSlice[0] = solana.Meta(account).WRITE()
	return initAcc3
}

// GetAccount gets the "account" account.
// The account to initialize.
func (initAcc3 *InitializeAccount3) GetAccount() *solana.AccountMeta {
	return initAcc3.AccountMetaSlice[0]
}

// SetMintAccount sets the "mint" account.
// The mint this account will be associated with.
func (initAcc3 *InitializeAccount3) SetMintAccount(mint solana.PublicKey) *InitializeAccount3 {
	initAcc3.AccountMetaSlice[1] = solana.Meta(mint)
	return initAcc3
}

// GetMintAccount gets the "mint" account.
// The mint this account will be associated with.
func (initAcc3 *InitializeAccount3) GetMintAccount() *solana.AccountMeta {
	return initAcc3.AccountMetaSlice[1]
}

func (initAcc3 InitializeAccount3) Build() *Instruction {
	return &Instruction{BaseVariant: encodbin.BaseVariant{
		Impl:   initAcc3,
		TypeID: encodbin.TypeIDFromUint8(Instruction_InitializeAccount3),
	}}
}

// ValidateAndBuild validates the instruction parameters and accounts;
// if there is a validation error, it returns the error.
// Otherwise, it builds and returns the instruction.
func (initAcc3 InitializeAccount3) ValidateAndBuild() (*Instruction, error) {
	if err := initAcc3.Validate(); err != nil {
		return nil, err
	}
	return initAcc3.Build(), nil
}

func (initAcc3 *InitializeAccount3) Validate() error {
	// Check whether all (required) parameters are set:
	{
		if initAcc3.Owner == nil {
			return errors.New("Owner parameter is not set")
		}
	}

	// Check whether all (required) accounts are set:
	{
		if initAcc3.AccountMetaSlice[0] == nil {
			return errors.New("accounts.Account is not set")
		}
		if initAcc3.AccountMetaSlice[1] == nil {
			return errors.New("accounts.Mint is not set")
		}
	}
	return nil
}

func (initAcc3 InitializeAccount3) MarshalWithEncoder(encoder *encodbin.Encoder) (err error) {
	// Serialize `Owner` param:
	err = encoder.Encode(initAcc3.Owner)
	if err != nil {
		return err
	}
	return nil
}

// NewInitializeAccount3Instruction declares a new InitializeAccount3 instruction with the provided parameters and accounts.
func NewInitializeAccount3Instruction(
	// Parameters:
	owner solana.PublicKey,
	// Accounts:
	account solana.PublicKey,
	mint solana.PublicKey) *InitializeAccount3 {
	return NewInitializeAccount3InstructionBuilder().
		SetOwner(owner).
		SetAccount(account).
		SetMintAccount(mint)
}
