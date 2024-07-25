package asset

import (
	"github.com/cielu/go-solana/common"
	"github.com/cielu/go-solana/pkg/encodbin"
	"github.com/cielu/go-solana/types/base"
)

type TransferLamportsToAta struct {

	//资产账户
	Asset common.Address
	//Ata账户(强制owner是Asset)
	WarpSolAccount common.Address
	//Mint地址
	Mint common.Address
	//资产账户管理员地址
	Administrator common.Address
	//工作者
	Worker common.Address
	//转移数量
	Amount                 uint64
	AssociatedTokenProgram common.Address
	TokenProgram           common.Address
	SystemProgram          common.Address

	base.AccountMetaSlice `bin:"-" borsh_skip:"true"`
}

func (init *TransferLamportsToAta) SetAsset(val common.Address) *TransferLamportsToAta {
	init.Asset = val
	return init
}

func (init *TransferLamportsToAta) SetAdministrator(val common.Address) *TransferLamportsToAta {
	init.Administrator = val
	return init
}

func (init *TransferLamportsToAta) SetMint(val common.Address) *TransferLamportsToAta {
	init.Mint = val
	return init
}

func (init *TransferLamportsToAta) SetWarpSolAccount(val common.Address) *TransferLamportsToAta {
	init.WarpSolAccount = val
	return init
}

func (init *TransferLamportsToAta) SetWorker(val common.Address) *TransferLamportsToAta {
	init.Worker = val
	return init
}

func (init *TransferLamportsToAta) SetAmount(val uint64) *TransferLamportsToAta {
	init.Amount = val
	return init
}

func NewTransferLamportsToAtaInstructionBuilder() *TransferLamportsToAta {
	init := &TransferLamportsToAta{}
	return init
}

func (init TransferLamportsToAta) MarshalWithEncoder(encoder *encodbin.Encoder) (err error) {
	// Serialize `Amount` param:
	err = encoder.Encode(init.Amount)
	if err != nil {
		return err
	}
	return nil
}

func (init TransferLamportsToAta) Build() *Instruction {

	keys := []*base.AccountMeta{
		{
			PublicKey:  init.Asset,
			IsSigner:   false,
			IsWritable: true,
		},
		{
			PublicKey:  init.WarpSolAccount,
			IsSigner:   false,
			IsWritable: true,
		},
		{
			PublicKey:  init.Mint,
			IsSigner:   false,
			IsWritable: false,
		},
		{
			PublicKey:  init.Administrator,
			IsSigner:   false,
			IsWritable: false,
		},
		{
			PublicKey:  init.Worker,
			IsSigner:   true,
			IsWritable: true,
		},
		{
			PublicKey:  base.SPLAssociatedTokenAccountProgramID,
			IsSigner:   false,
			IsWritable: false,
		},
		{
			PublicKey:  base.TokenProgramID,
			IsSigner:   false,
			IsWritable: false,
		},
		{
			PublicKey:  base.SystemProgramID,
			IsSigner:   false,
			IsWritable: false,
		},
	}

	init.AccountMetaSlice = keys

	return &Instruction{BaseVariant: encodbin.BaseVariant{
		Impl:   init,
		TypeID: encodbin.TypeIDFromUint64(Discriminator_transferLamportsToAta, encodbin.LE),
	}}
}

func NewTransferLamportsToAtaInstruction(
	Asset common.Address,
	WarpSolAccount common.Address,
	Mint common.Address,
	Administrator common.Address,
	Worker common.Address,
	Amount uint64,
) *TransferLamportsToAta {
	return NewTransferLamportsToAtaInstructionBuilder().
		SetAsset(Asset).
		SetWarpSolAccount(WarpSolAccount).
		SetMint(Mint).
		SetAdministrator(Administrator).
		SetWorker(Worker).
		SetAmount(Amount)
}
