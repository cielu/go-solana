package associatedaccount

import (
	"github.com/cielu/go-solana"
	"github.com/cielu/go-solana/core/base"
	"github.com/cielu/go-solana/pkg/encodbin"
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
	Payer          solana.PublicKey `bin:"-" borsh_skip:"true"`
	Wallet         solana.PublicKey `bin:"-" borsh_skip:"true"`
	Mint           solana.PublicKey `bin:"-" borsh_skip:"true"`
	TokenProgramID solana.PublicKey `bin:"-" borsh_skip:"true"`
	CreateType     uint8            `bin:"-" borsh_skip:"true"`

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
	solana.AccountMetaSlice `bin:"-" borsh_skip:"true"`
}

// NewCreateInstructionBuilder creates a new `Create` instruction builder.
func NewCreateInstructionBuilder() *Create {
	nd := &Create{}
	return nd
}

func (inst *Create) SetPayer(payer solana.PublicKey) *Create {
	inst.Payer = payer
	return inst
}

func (inst *Create) SetWallet(wallet solana.PublicKey) *Create {
	inst.Wallet = wallet
	return inst
}

func (inst *Create) SetMint(mint solana.PublicKey) *Create {
	inst.Mint = mint
	return inst
}

func (inst *Create) SetTokenProgramID(tokenProgramID solana.PublicKey) *Create {
	inst.TokenProgramID = tokenProgramID
	return inst
}

// SetCreateIdempotent 0 = Create, 1 = CreateIdempotent
func (inst *Create) SetCreateIdempotent() *Create {
	// setType --> 1
	inst.CreateType = 1
	return inst
}

func (inst Create) Build() *Instruction {

	var associatedTokenAddress solana.PublicKey

	// Find the associatedTokenAddress;
	switch inst.TokenProgramID {
	case base.Token2022ProgramID:
		associatedTokenAddress, _, _ = base.FindAssociatedTokenAddress(inst.Wallet, inst.Mint, base.Token2022ProgramID)
	default:
		inst.TokenProgramID = base.TokenProgramID
		associatedTokenAddress, _, _ = base.FindAssociatedTokenAddress(inst.Wallet, inst.Mint)
	}

	keys := []*solana.AccountMeta{
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

	typeID := encodbin.NoTypeIDDefaultID
	// custom associated account create type
	switch inst.CreateType {
	case 1: // CreateIdempotent
		typeID = encodbin.TypeIDFromUint8(inst.CreateType)
	}

	return &Instruction{BaseVariant: encodbin.BaseVariant{
		Impl:   inst,
		TypeID: typeID,
	}}
}

func (inst Create) MarshalWithEncoder(encoder *encodbin.Encoder) error {
	return encoder.WriteBytes([]byte{}, false)
}

func NewCreateInstruction(
	payer solana.PublicKey,
	walletAddress solana.PublicKey,
	splTokenMintAddress solana.PublicKey,
) *Create {
	return NewCreateInstructionBuilder().
		SetPayer(payer).
		SetWallet(walletAddress).
		SetMint(splTokenMintAddress)
}
