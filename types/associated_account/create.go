package associated_account

import (
	"github.com/cielu/go-solana/common"
	"github.com/cielu/go-solana/pkg/encodbin"
	"github.com/cielu/go-solana/types"
	"github.com/cielu/go-solana/types/base"
)

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

type Create struct {
	Payer          common.Address `bin:"-" borsh_skip:"true"`
	Wallet         common.Address `bin:"-" borsh_skip:"true"`
	Mint           common.Address `bin:"-" borsh_skip:"true"`
	TokenProgramID common.Address `bin:"-" borsh_skip:"true"`

	// [0] = [WRITE, SIGNER] Payer
	// ··········· Funding account
	//
	// [1] = [WRITE] AssociatedTokenAccount
	// ··········· Associated token account address to be created
	//
	// [2] = [] Wallet
	// ··········· Wallet address for the new associated token account
	//
	// [3] = [] TokenMint
	// ··········· The token mint for the new associated token account
	//
	// [4] = [] SystemProgram
	// ··········· System program ID
	//
	// [5] = [] TokenProgram
	// ··········· SPL token program ID
	//
	// [6] = [] SysVarRent
	// ··········· SysVarRentPubkey
	AccountMetaSlice []*types.AccountMeta `bin:"-" borsh_skip:"true"`
}

// NewCreateInstructionBuilder creates a new `Create` instruction builder.
func NewCreateInstructionBuilder() *Create {
	nd := &Create{}
	return nd
}

func (inst *Create) SetPayer(payer common.Address) *Create {
	inst.Payer = payer
	return inst
}

func (inst *Create) SetWallet(wallet common.Address) *Create {
	inst.Wallet = wallet
	return inst
}

func (inst *Create) SetMint(mint common.Address) *Create {
	inst.Mint = mint
	return inst
}

func (inst *Create) SetTokenProgramID(tokenProgramID common.Address) *Create {
	inst.TokenProgramID = tokenProgramID
	return inst
}

func (inst Create) Build() *Instruction {

	// Find the associatedTokenAddress;
	associatedTokenAddress, _, _ := FindAssociatedTokenAddress(
		inst.Wallet,
		inst.Mint,
	)
	if inst.TokenProgramID == base.Token2022ProgramID {
		associatedTokenAddress, _, _ = FindAssociatedTokenAddress(inst.Wallet, inst.Mint, base.Token2022ProgramID)
	}

	keys := []*types.AccountMeta{
		{
			PublicKey:  inst.Payer,
			IsSigner:   true,
			IsWritable: true,
		},
		{
			PublicKey:  associatedTokenAddress,
			IsSigner:   false,
			IsWritable: true,
		},
		{
			PublicKey:  inst.Wallet,
			IsSigner:   false,
			IsWritable: false,
		},
		{
			PublicKey:  inst.Mint,
			IsSigner:   false,
			IsWritable: false,
		},
		{
			PublicKey:  base.SystemProgramID,
			IsSigner:   false,
			IsWritable: false,
		},
		{
			PublicKey:  inst.TokenProgramID,
			IsSigner:   false,
			IsWritable: false,
		},
		{
			PublicKey:  base.SysVarRentPubkey,
			IsSigner:   false,
			IsWritable: false,
		},
	}

	inst.AccountMetaSlice = keys

	return &Instruction{BaseVariant: encodbin.BaseVariant{
		Impl:   inst,
		TypeID: encodbin.NoTypeIDDefaultID,
	}}
}

func (inst Create) MarshalWithEncoder(encoder *encodbin.Encoder) error {
	return encoder.WriteBytes([]byte{}, false)
}

func NewCreateInstruction(
	payer common.Address,
	walletAddress common.Address,
	splTokenMintAddress common.Address,
) *Create {
	return NewCreateInstructionBuilder().
		SetPayer(payer).
		SetWallet(walletAddress).
		SetMint(splTokenMintAddress).
		SetTokenProgramID(base.TokenProgramID)
}
