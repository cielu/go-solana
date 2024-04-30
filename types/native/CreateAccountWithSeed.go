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

package system

import (
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/cielu/go-solana/common"
	"github.com/cielu/go-solana/pkg/encodbin"
	"github.com/cielu/go-solana/types/base"
)

// CreateAccountWithSeed Create a new account at an address derived from a base pubkey and a seed
type CreateAccountWithSeed struct {
	// Base public key
	Base *common.Address

	// String of ASCII chars, no longer than Pubkey::MAX_SEED_LEN
	Seed *string

	// Number of lamports to transfer to the new account
	Lamports *uint64

	// Number of bytes of memory to allocate
	Space *uint64

	// Owner program account address
	Owner *common.Address

	// [0] = [WRITE, SIGNER] FundingAccount
	// ··········· Funding account
	//
	// [1] = [WRITE] CreatedAccount
	// ··········· Created account
	//
	// [2] = [SIGNER] BaseAccount
	// ··········· Base account
	AccountMeta []*base.AccountMeta `bin:"-" borsh_skip:"true"`
}

// NewCreateAccountWithSeedInstructionBuilder creates a new `CreateAccountWithSeed` instruction builder.
func NewCreateAccountWithSeedInstructionBuilder() *CreateAccountWithSeed {
	nd := &CreateAccountWithSeed{
		AccountMeta: make([]*base.AccountMeta, 3),
	}
	return nd
}

// Base public key
func (cAcc *CreateAccountWithSeed) SetBase(base common.Address) *CreateAccountWithSeed {
	cAcc.Base = &base
	return cAcc
}

// String of ASCII chars, no longer than Pubkey::MAX_SEED_LEN
func (cAcc *CreateAccountWithSeed) SetSeed(seed string) *CreateAccountWithSeed {
	cAcc.Seed = &seed
	return cAcc
}

// Number of lamports to transfer to the new account
func (cAcc *CreateAccountWithSeed) SetLamports(lamports uint64) *CreateAccountWithSeed {
	cAcc.Lamports = &lamports
	return cAcc
}

// Number of bytes of memory to allocate
func (cAcc *CreateAccountWithSeed) SetSpace(space uint64) *CreateAccountWithSeed {
	cAcc.Space = &space
	return cAcc
}

// Owner program account address
func (cAcc *CreateAccountWithSeed) SetOwner(owner common.Address) *CreateAccountWithSeed {
	cAcc.Owner = &owner
	return cAcc
}

// Funding account
func (cAcc *CreateAccountWithSeed) SetFundingAccount(fundingAccount common.Address) *CreateAccountWithSeed {
	cAcc.AccountMeta[0] = base.Meta(fundingAccount).WRITE().SIGNER()
	return cAcc
}

func (cAcc *CreateAccountWithSeed) GetFundingAccount() *base.AccountMeta {
	return cAcc.AccountMeta[0]
}

// Created account
func (cAcc *CreateAccountWithSeed) SetCreatedAccount(createdAccount common.Address) *CreateAccountWithSeed {
	cAcc.AccountMeta[1] = base.Meta(createdAccount).WRITE()
	return cAcc
}

func (cAcc *CreateAccountWithSeed) GetCreatedAccount() *base.AccountMeta {
	return cAcc.AccountMeta[1]
}

// Base account
func (cAcc *CreateAccountWithSeed) SetBaseAccount(baseAccount common.Address) *CreateAccountWithSeed {
	cAcc.AccountMeta[2] = base.Meta(baseAccount).SIGNER()
	return cAcc
}

func (cAcc *CreateAccountWithSeed) GetBaseAccount() *base.AccountMeta {
	return cAcc.AccountMeta[2]
}

func (cAcc CreateAccountWithSeed) Build() *Instruction {
	{
		if *cAcc.Base != cAcc.GetFundingAccount().PublicKey {
			cAcc.AccountMeta[2] = base.Meta(*cAcc.Base).SIGNER()
		}
	}
	return &Instruction{BaseVariant: encodbin.BaseVariant{
		Impl:   cAcc,
		TypeID: encodbin.TypeIDFromUint32(Instruction_CreateAccountWithSeed, binary.LittleEndian),
	}}
}

// ValidateAndBuild validates the instruction parameters and accounts;
// if there is a validation error, it returns the error.
// Otherwise, it builds and returns the instruction.
func (cAcc CreateAccountWithSeed) ValidateAndBuild() (*Instruction, error) {
	if err := cAcc.Validate(); err != nil {
		return nil, err
	}
	return cAcc.Build(), nil
}

func (cAcc *CreateAccountWithSeed) Validate() error {
	// Check whether all (required) parameters are set:
	{
		if cAcc.Base == nil {
			return errors.New("Base parameter is not set")
		}
		if cAcc.Seed == nil {
			return errors.New("Seed parameter is not set")
		}
		if cAcc.Lamports == nil {
			return errors.New("Lamports parameter is not set")
		}
		if cAcc.Space == nil {
			return errors.New("Space parameter is not set")
		}
		if cAcc.Owner == nil {
			return errors.New("Owner parameter is not set")
		}
	}

	// Check whether all accounts are set:
	{
		if cAcc.AccountMeta[0] == nil {
			return fmt.Errorf("FundingAccount is not set")
		}
		if cAcc.AccountMeta[1] == nil {
			return fmt.Errorf("CreatedAccount is not set")
		}
	}
	return nil
}

func (cAcc CreateAccountWithSeed) MarshalWithEncoder(encoder *encodbin.Encoder) error {
	// Serialize `Base` param:
	{
		err := encoder.Encode(*cAcc.Base)
		if err != nil {
			return err
		}
	}
	// Serialize `Seed` param:
	{
		err := encoder.WriteRustString(*cAcc.Seed)
		if err != nil {
			return err
		}
	}
	// Serialize `Lamports` param:
	{
		err := encoder.Encode(*cAcc.Lamports)
		if err != nil {
			return err
		}
	}
	// Serialize `Space` param:
	{
		err := encoder.Encode(*cAcc.Space)
		if err != nil {
			return err
		}
	}
	// Serialize `Owner` param:
	{
		err := encoder.Encode(*cAcc.Owner)
		if err != nil {
			return err
		}
	}
	return nil
}

// NewCreateAccountWithSeedInstruction declares a new CreateAccountWithSeed instruction with the provided parameters and accounts.
func NewCreateAccountWithSeedInstruction(
	// Parameters:
	base common.Address,
	seed string,
	lamports uint64,
	space uint64,
	owner common.Address,
	// Accounts:
	fundingAccount common.Address,
	createdAccount common.Address,
	baseAccount common.Address) *CreateAccountWithSeed {
	return NewCreateAccountWithSeedInstructionBuilder().
		SetBase(base).
		SetSeed(seed).
		SetLamports(lamports).
		SetSpace(space).
		SetOwner(owner).
		SetFundingAccount(fundingAccount).
		SetCreatedAccount(createdAccount).
		SetBaseAccount(baseAccount)
}
