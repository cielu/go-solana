package asset

import (
	"bytes"
	"fmt"
	"github.com/cielu/go-solana/common"
	"github.com/cielu/go-solana/pkg/encodbin"
	"github.com/cielu/go-solana/types/base"
)

const (
	Discriminator_initialize             uint64 = 0xafaf6d1f0d989bed
	Discriminator_transferLamportsToAta  uint64 = 0x19abfdc53787513b
	Discriminator_syncNativeAta          uint64 = 0xf22a38bcd6c68795
	Discriminator_proxyOrcaSwap          uint64 = 0xddbc843aea104171
	Discriminator_proxyRaydiumSwapBaseIn uint64 = 0x1c49426140d5541a
)

type Instruction struct {
	encodbin.BaseVariant
	AssertExecutorProgramID common.Address
}

func (inst *Instruction) ProgramID() common.Address {
	if inst.AssertExecutorProgramID.IsEmpty() {
		return base.AssetExecutorProgramID
	}
	return base.AssetExecutorProgramID
}

func (inst *Instruction) SetProgramID(tokenProgramID common.Address) {
	inst.AssertExecutorProgramID = tokenProgramID
}

func (inst *Instruction) Accounts() (out []*base.AccountMeta) {
	return inst.Impl.(base.AccountsGettable).GetAccounts()
}

func (inst *Instruction) Data() ([]byte, error) {
	buf := new(bytes.Buffer)
	if err := encodbin.NewBinEncoder(buf).Encode(inst); err != nil {
		return nil, fmt.Errorf("unable to encode instruction: %w", err)
	}
	return buf.Bytes(), nil
}

func (inst *Instruction) MarshalWithEncoder(encoder *encodbin.Encoder) error {
	err := encoder.WriteUint64(inst.TypeID.Uint64(), encodbin.BE)
	if err != nil {
		return fmt.Errorf("unable to write variant type: %w", err)
	}
	return encoder.Encode(inst.Impl)
}
