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

package computebudget

import (
	"errors"
	"github.com/cielu/go-solana/pkg/encodbin"
	"github.com/cielu/go-solana/types/base"
)

type SetComputeUnitPrice struct {
	MicroLamports uint64
}

func (setPrice *SetComputeUnitPrice) SetAccounts(accounts []*base.AccountMeta) error {
	return nil
}

func (setPrice SetComputeUnitPrice) GetAccounts() (accounts []*base.AccountMeta) {
	return
}

// NewSetComputeUnitPriceInstructionBuilder creates a new `SetComputeUnitPrice` instruction builder.
func NewSetComputeUnitPriceInstructionBuilder() *SetComputeUnitPrice {
	nd := &SetComputeUnitPrice{}
	return nd
}

func (setPrice *SetComputeUnitPrice) SetMicroLamports(microLamports uint64) *SetComputeUnitPrice {
	setPrice.MicroLamports = microLamports
	return setPrice
}

func (setPrice SetComputeUnitPrice) Build() *Instruction {
	return &Instruction{
		BaseVariant: encodbin.BaseVariant{
			Impl:   setPrice,
			TypeID: encodbin.TypeIDFromUint8(Instruction_SetComputeUnitPrice),
		},
	}
}

// ValidateAndBuild validates the instruction parameters and accounts;
// if there is a validation error, it returns the error.
// Otherwise, it builds and returns the instruction.
func (setPrice SetComputeUnitPrice) ValidateAndBuild() (*Instruction, error) {
	if err := setPrice.Validate(); err != nil {
		return nil, err
	}
	return setPrice.Build(), nil
}

func (setPrice *SetComputeUnitPrice) Validate() error {
	// Check whether all (required) parameters are set:
	{
		if setPrice.MicroLamports == 0 {
			return errors.New("MicroLamports parameter is not set")
		}
	}
	return nil
}

func (setPrice SetComputeUnitPrice) MarshalWithEncoder(encoder *encodbin.Encoder) (err error) {
	// Serialize `MicroLamports` param:
	err = encoder.Encode(setPrice.MicroLamports)
	if err != nil {
		return err
	}
	return nil
}

// NewSetComputeUnitPriceInstruction declares a new SetComputeUnitPrice instruction with the provided parameters and accounts.
func NewSetComputeUnitPriceInstruction(
	// Parameters:
	microLamports uint64,
) *SetComputeUnitPrice {
	return NewSetComputeUnitPriceInstructionBuilder().SetMicroLamports(microLamports)
}
