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

	"github.com/cielu/go-solana"
	"github.com/cielu/go-solana/pkg/encodbin"
)

// CreateAccount Create a new account
type CreateAccount struct {
	// Number of lamports to transfer to the new account
	Lamports *uint64

	// Number of bytes of memory to allocate
	Space *uint64

	// PublicKey of program that will own the new account
	Owner *solana.PublicKey

	// [0] = [WRITE, SIGNER] FundingAccount
	// ··········· Funding account
	//
	// [1] = [WRITE, SIGNER] NewAccount
	// ··········· New account
	solana.AccountMetaSlice `bin:"-" borsh_skip:"true"`
}

// NewCreateAccountInstructionBuilder creates a new `CreateAccount` instruction builder.
func NewCreateAccountInstructionBuilder() *CreateAccount {
	nd := &CreateAccount{
		AccountMetaSlice: make([]*solana.AccountMeta, 2),
	}
	return nd
}

// Number of lamports to transfer to the new account
func (cAcc *CreateAccount) SetLamports(lamports uint64) *CreateAccount {
	cAcc.Lamports = &lamports
	return cAcc
}

// Number of bytes of memory to allocate
func (cAcc *CreateAccount) SetSpace(space uint64) *CreateAccount {
	cAcc.Space = &space
	return cAcc
}

// PublicKey of program that will own the new account
func (cAcc *CreateAccount) SetOwner(owner solana.PublicKey) *CreateAccount {
	cAcc.Owner = &owner
	return cAcc
}

// Funding account
func (cAcc *CreateAccount) SetFundingAccount(fundingAccount solana.PublicKey) *CreateAccount {
	cAcc.AccountMetaSlice[0] = solana.Meta(fundingAccount).WRITE().SIGNER()
	return cAcc
}

func (cAcc *CreateAccount) GetFundingAccount() *solana.AccountMeta {
	return cAcc.AccountMetaSlice[0]
}

// New account
func (cAcc *CreateAccount) SetNewAccount(newAccount solana.PublicKey) *CreateAccount {
	cAcc.AccountMetaSlice[1] = solana.Meta(newAccount).WRITE().SIGNER()
	return cAcc
}

func (cAcc *CreateAccount) GetNewAccount() *solana.AccountMeta {
	return cAcc.AccountMetaSlice[1]
}

func (cAcc CreateAccount) Build() *Instruction {
	return &Instruction{BaseVariant: encodbin.BaseVariant{
		Impl:   cAcc,
		TypeID: encodbin.TypeIDFromUint32(Instruction_CreateAccount, binary.LittleEndian),
	}}
}

func (cAcc CreateAccount) MarshalWithEncoder(encoder *encodbin.Encoder) error {
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

// NewCreateAccountInstruction declares a new CreateAccount instruction with the provided parameters and accounts.
func NewCreateAccountInstruction(
	// Parameters:
	lamports uint64,
	space uint64,
	owner solana.PublicKey,
	// Accounts:
	fundingAccount solana.PublicKey,
	newAccount solana.PublicKey) *CreateAccount {
	return NewCreateAccountInstructionBuilder().
		SetLamports(lamports).
		SetSpace(space).
		SetOwner(owner).
		SetFundingAccount(fundingAccount).
		SetNewAccount(newAccount)
}
