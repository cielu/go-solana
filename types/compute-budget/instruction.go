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
	"bytes"
	"fmt"

	"github.com/cielu/go-solana"
	"github.com/cielu/go-solana/pkg/encodbin"
	"github.com/cielu/go-solana/types/base"
)

const (
	// Instruction_RequestUnitsDeprecated Deprecated
	// after feature remove_deprecated_request_unit_ix::id() is activated
	Instruction_RequestUnitsDeprecated uint8 = iota

	// Instruction_RequestHeapFrame Request a specific transaction-wide program heap region size in bytes.
	// The value requested must be a multiple of 1024. This new heap region
	// size applies to each program executed in the transaction,
	// including all calls to CPIs.
	Instruction_RequestHeapFrame

	// Instruction_SetComputeUnitLimit Set a specific compute unit limit that the transaction is allowed to consume.
	Instruction_SetComputeUnitLimit

	// Instruction_SetComputeUnitPrice Set a compute unit price in "micro-lamports" to pay a higher transaction
	// fee for higher transaction prioritization.
	Instruction_SetComputeUnitPrice
)

type Instruction struct {
	encodbin.BaseVariant
}

func (inst *Instruction) ProgramID() solana.PublicKey {
	return base.ComputeBudget
}

func (inst *Instruction) Accounts() (out []*solana.AccountMeta) {
	return inst.Impl.(solana.AccountsGettable).GetAccounts()
}

func (inst *Instruction) Data() ([]byte, error) {
	buf := new(bytes.Buffer)
	if err := encodbin.NewBinEncoder(buf).Encode(inst); err != nil {
		return nil, fmt.Errorf("unable to encode instruction: %w", err)
	}
	return buf.Bytes(), nil
}

func (inst *Instruction) MarshalWithEncoder(encoder *encodbin.Encoder) error {
	err := encoder.WriteUint8(inst.TypeID.Uint8())
	if err != nil {
		return fmt.Errorf("unable to write variant type: %w", err)
	}
	return encoder.Encode(inst.Impl)
}
