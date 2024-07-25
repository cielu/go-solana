package asset

import (
	"github.com/cielu/go-solana/common"
	"github.com/cielu/go-solana/pkg/encodbin"
	"github.com/cielu/go-solana/types/base"
)

type ProxyRaydiumSwapBaseIn struct {
	AmmProgram        common.Address
	Amm               common.Address
	AmmAuthority      common.Address
	AmmOpenOrders     common.Address
	AmmTargetOrders   common.Address
	AmmCoinVault      common.Address
	AmmPcVault        common.Address
	MarketProgram     common.Address
	Market            common.Address
	MarketBids        common.Address
	MarketAsks        common.Address
	MarketEventQueue  common.Address
	MarketCoinVault   common.Address
	MarketPcVault     common.Address
	MarketVaultSigner common.Address
	//out by only owner asset ata.
	UserTokenSource common.Address
	//in by only owner asset ata.
	UserTokenDestination common.Address
	Asset                common.Address
	Administrator        common.Address
	Worker               common.Address
	//交易得到的代币Mint地址,填写错误会导致模拟交易失败!
	//only destination token mint.
	Mint             common.Address
	AmountIn         uint64
	MinimumAmountOut uint64
	//当source余额为0时,是否自动关闭ata账户并且返回SOL到Worker.
	CanClose bool

	base.AccountMetaSlice `bin:"-" borsh_skip:"true"`
}

func (init ProxyRaydiumSwapBaseIn) MarshalWithEncoder(encoder *encodbin.Encoder) (err error) {
	// Serialize `Amount` param:
	err = encoder.Encode(init.AmountIn)
	if err != nil {
		return err
	}
	// Serialize `MinimumAmountOut` param:
	err = encoder.Encode(init.MinimumAmountOut)
	if err != nil {
		return err
	}
	// Serialize `CanClose` param:
	err = encoder.Encode(init.CanClose)
	if err != nil {
		return err
	}
	return nil
}

func (init ProxyRaydiumSwapBaseIn) Build() *Instruction {

	keys := []*base.AccountMeta{
		{
			PublicKey:  init.AmmProgram,
			IsSigner:   false,
			IsWritable: false,
		},
		{
			PublicKey:  init.Amm,
			IsSigner:   false,
			IsWritable: true,
		},
		{
			PublicKey:  init.AmmAuthority,
			IsSigner:   false,
			IsWritable: false,
		},
		{
			PublicKey:  init.AmmOpenOrders,
			IsSigner:   false,
			IsWritable: true,
		},
		{
			PublicKey:  init.AmmTargetOrders,
			IsSigner:   false,
			IsWritable: true,
		},
		{
			PublicKey:  init.AmmCoinVault,
			IsSigner:   false,
			IsWritable: true,
		},
		{
			PublicKey:  init.AmmPcVault,
			IsSigner:   false,
			IsWritable: true,
		},
		{
			PublicKey:  init.MarketProgram,
			IsSigner:   false,
			IsWritable: false,
		},
		{
			PublicKey:  init.Market,
			IsSigner:   false,
			IsWritable: true,
		},
		{
			PublicKey:  init.MarketBids,
			IsSigner:   false,
			IsWritable: true,
		},
		{
			PublicKey:  init.MarketAsks,
			IsSigner:   false,
			IsWritable: true,
		},
		{
			PublicKey:  init.MarketEventQueue,
			IsSigner:   false,
			IsWritable: true,
		},
		{
			PublicKey:  init.MarketCoinVault,
			IsSigner:   false,
			IsWritable: true,
		},
		{
			PublicKey:  init.MarketPcVault,
			IsSigner:   false,
			IsWritable: true,
		},
		{
			PublicKey:  init.MarketVaultSigner,
			IsSigner:   false,
			IsWritable: true,
		},
		{
			PublicKey:  init.UserTokenSource,
			IsSigner:   false,
			IsWritable: true,
		},
		{
			PublicKey:  init.UserTokenDestination,
			IsSigner:   false,
			IsWritable: true,
		},
		{
			PublicKey:  init.Asset,
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
			PublicKey:  init.Mint,
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
		{
			PublicKey:  base.SPLAssociatedTokenAccountProgramID,
			IsSigner:   false,
			IsWritable: false,
		},
	}

	init.AccountMetaSlice = keys

	return &Instruction{BaseVariant: encodbin.BaseVariant{
		Impl:   init,
		TypeID: encodbin.TypeIDFromUint64(Discriminator_proxyRaydiumSwapBaseIn, encodbin.LE),
	}}
}

func NewProxyRaydiumSwapBaseIn(
	AmmProgram common.Address,
	Amm common.Address,
	AmmAuthority common.Address,
	AmmOpenOrders common.Address,
	AmmTargetOrders common.Address,
	AmmCoinVault common.Address,
	AmmPcVault common.Address,
	MarketProgram common.Address,
	Market common.Address,
	MarketBids common.Address,
	MarketAsks common.Address,
	MarketEventQueue common.Address,
	MarketCoinVault common.Address,
	MarketPcVault common.Address,
	MarketVaultSigner common.Address,
	UserTokenSource common.Address,
	UserTokenDestination common.Address,
	Asset common.Address,
	Administrator common.Address,
	Worker common.Address,
	Mint common.Address,
	AmountIn uint64,
	MinimumAmountOut uint64,
	CanClose bool,
) *ProxyRaydiumSwapBaseIn {
	return &ProxyRaydiumSwapBaseIn{
		AmmProgram:           AmmProgram,
		Amm:                  Amm,
		AmmAuthority:         AmmAuthority,
		AmmOpenOrders:        AmmOpenOrders,
		AmmTargetOrders:      AmmTargetOrders,
		AmmCoinVault:         AmmCoinVault,
		AmmPcVault:           AmmPcVault,
		MarketProgram:        MarketProgram,
		Market:               Market,
		MarketBids:           MarketBids,
		MarketAsks:           MarketAsks,
		MarketEventQueue:     MarketEventQueue,
		MarketCoinVault:      MarketCoinVault,
		MarketPcVault:        MarketPcVault,
		MarketVaultSigner:    MarketVaultSigner,
		UserTokenSource:      UserTokenSource,
		UserTokenDestination: UserTokenDestination,
		Asset:                Asset,
		Administrator:        Administrator,
		Worker:               Worker,
		Mint:                 Mint,
		AmountIn:             AmountIn,
		MinimumAmountOut:     MinimumAmountOut,
		CanClose:             CanClose,
	}
}
