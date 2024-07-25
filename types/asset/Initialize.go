package asset

import (
	"github.com/cielu/go-solana/common"
	"github.com/cielu/go-solana/pkg/encodbin"
	"github.com/cielu/go-solana/types/base"
)

type Initialize struct {

	//资产账户
	Asset common.Address
	//FeePayer
	User common.Address

	//资产合约超级管理员
	Administrator common.Address
	//多签管理员一
	SingerA common.Address
	//多签管理员二
	SingerB common.Address
	//多签管理员三
	SingerC common.Address
	//工作者(主)
	PrimaryWorker common.Address
	//工作者(次)
	SecondaryWorker common.Address
	// [0] = [WRITE] asset(资产管理合约)
	// [1] = [WRITE、Singer] feePay(支付者)
	// [2] = [SystemProgram]

	base.AccountMetaSlice `bin:"-" borsh_skip:"true"`
}

// 合约超级管理员.
func (init *Initialize) SetAdministrator(asset common.Address) *Initialize {
	init.Administrator = asset
	return init
}

// 多签A
func (init *Initialize) SetSingerA(asset common.Address) *Initialize {
	init.SingerA = asset
	return init
}

// 多签B
func (init *Initialize) SetSingerB(asset common.Address) *Initialize {
	init.SingerB = asset
	return init
}

// 多签C
func (init *Initialize) SetSingerC(asset common.Address) *Initialize {
	init.SingerC = asset
	return init
}

// 工作者(主)
func (init *Initialize) SetPrimaryWorker(asset common.Address) *Initialize {
	init.PrimaryWorker = asset
	return init
}

// 工作者(副)
func (init *Initialize) SetSecondaryWorker(asset common.Address) *Initialize {
	init.SecondaryWorker = asset
	return init
}

// 资产账户
func (init *Initialize) SetAssetAccount(asset common.Address) *Initialize {
	init.Asset = asset
	return init
}

// FeePayer
func (init *Initialize) SetUserAccount(asset common.Address) *Initialize {
	init.User = asset
	return init
}

func NewInitInstructionBuilder() *Initialize {
	init := &Initialize{}
	return init
}

func (init Initialize) MarshalWithEncoder(encoder *encodbin.Encoder) (err error) {
	// Serialize `Amount` param:
	err = encoder.Encode(init.Administrator)
	if err != nil {
		return err
	}
	// Serialize `Decimals` param:
	err = encoder.Encode(init.SingerA)
	if err != nil {
		return err
	}
	err = encoder.Encode(init.SingerB)
	if err != nil {
		return err
	}
	// Serialize `Decimals` param:
	err = encoder.Encode(init.SingerC)
	if err != nil {
		return err
	}
	err = encoder.Encode(init.PrimaryWorker)
	if err != nil {
		return err
	}
	// Serialize `Decimals` param:
	err = encoder.Encode(init.SecondaryWorker)
	if err != nil {
		return err
	}
	return nil
}

func (init Initialize) Build() *Instruction {

	keys := []*base.AccountMeta{
		{
			PublicKey:  init.Asset,
			IsSigner:   false,
			IsWritable: true,
		},
		{
			PublicKey:  init.User,
			IsSigner:   true,
			IsWritable: true,
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
		TypeID: encodbin.TypeIDFromUint64(Discriminator_initialize, encodbin.LE),
	}}
}

func NewInitInstruction(
	Administrator common.Address,
	SingerA common.Address,
	SingerB common.Address,
	SingerC common.Address,
	PrimaryWorker common.Address,
	SecondaryWorker common.Address,
	AssetAccount common.Address,
	UserAccount common.Address,
) *Initialize {
	return NewInitInstructionBuilder().
		SetAdministrator(Administrator).
		SetSingerA(SingerA).
		SetSingerB(SingerB).
		SetSingerC(SingerC).
		SetPrimaryWorker(PrimaryWorker).
		SetSecondaryWorker(SecondaryWorker).
		SetAssetAccount(AssetAccount).
		SetUserAccount(UserAccount)
}
