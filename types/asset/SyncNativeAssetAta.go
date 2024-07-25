package asset

import (
	"github.com/cielu/go-solana/common"
	"github.com/cielu/go-solana/pkg/encodbin"
	"github.com/cielu/go-solana/types/base"
)

type SyncNativeAssetAta struct {

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

	AssociatedTokenProgram common.Address
	TokenProgram           common.Address
	SystemProgram          common.Address

	base.AccountMetaSlice `bin:"-" borsh_skip:"true"`
}

func (init *SyncNativeAssetAta) SetAsset(val common.Address) *SyncNativeAssetAta {
	init.Asset = val
	return init
}

func (init *SyncNativeAssetAta) SetAdministrator(val common.Address) *SyncNativeAssetAta {
	init.Administrator = val
	return init
}

func (init *SyncNativeAssetAta) SetMint(val common.Address) *SyncNativeAssetAta {
	init.Mint = val
	return init
}

func (init *SyncNativeAssetAta) SetWarpSolAccount(val common.Address) *SyncNativeAssetAta {
	init.WarpSolAccount = val
	return init
}

func (init *SyncNativeAssetAta) SetWorker(val common.Address) *SyncNativeAssetAta {
	init.Worker = val
	return init
}

func NewSyncNativeAssetAtaInstructionBuilder() *SyncNativeAssetAta {
	init := &SyncNativeAssetAta{}
	return init
}

func (init SyncNativeAssetAta) MarshalWithEncoder(encoder *encodbin.Encoder) (err error) {
	return nil
}

func (init SyncNativeAssetAta) Build() *Instruction {

	keys := []*base.AccountMeta{
		{
			PublicKey:  init.Asset,
			IsSigner:   false,
			IsWritable: false,
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
			PublicKey:  base.TokenProgramID,
			IsSigner:   false,
			IsWritable: false,
		},
	}

	init.AccountMetaSlice = keys

	return &Instruction{BaseVariant: encodbin.BaseVariant{
		Impl:   init,
		TypeID: encodbin.TypeIDFromUint64(Discriminator_syncNativeAta, encodbin.LE),
	}}
}

func NewSyncNativeAssetAtaInstruction(
	Asset common.Address,
	WarpSolAccount common.Address,
	Mint common.Address,
	Administrator common.Address,
	Worker common.Address,
) *SyncNativeAssetAta {
	return NewSyncNativeAssetAtaInstructionBuilder().
		SetAsset(Asset).
		SetWarpSolAccount(WarpSolAccount).
		SetMint(Mint).
		SetAdministrator(Administrator).
		SetWorker(Worker)
}
