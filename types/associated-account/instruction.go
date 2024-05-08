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

package associatedaccount

import (
	"github.com/cielu/go-solana/common"
	"github.com/cielu/go-solana/pkg/encodbin"
	"github.com/cielu/go-solana/types/base"
)

type Instruction struct {
	encodbin.BaseVariant
}

func (inst *Instruction) ProgramID() common.Address {
	return base.SPLAssociatedTokenAccountProgramID
}

func (inst *Instruction) Accounts() (out []*base.AccountMeta) {
	return base.GetAccounts(inst.Impl.([]*base.AccountMeta))
}

func (inst *Instruction) Data() ([]byte, error) {
	return []byte{}, nil
}

func (inst *Instruction) MarshalWithEncoder(encoder *encodbin.Encoder) error {
	return encoder.Encode(inst.Impl)
}
