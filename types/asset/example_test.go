package asset

import (
	"context"
	"fmt"
	"github.com/cielu/go-solana/common"
	"github.com/cielu/go-solana/crypto"
	"github.com/cielu/go-solana/solclient"
	"github.com/cielu/go-solana/types"
	"github.com/cielu/go-solana/types/base"
	"testing"
)

func newClient() *solclient.Client {
	rpcUrl := "https://prettiest-shy-haze.solana-devnet.quiknode.pro/ae7f6e085169a5d0401ffb5c3f0071a362b72bbb/"
	c, err := solclient.Dial(rpcUrl)
	if err != nil {
		panic("Dial rpc endpoint failed")
	}
	return c
}

func TestProxyraydiumSwapBaseIn(t *testing.T) {
	var (
		c   = newClient()
		ctx = context.Background()
	)

	administrator := common.StrToAddress("AMPd3qHJsJxjh8T9YDZXPDCQf7bjd13Mqj7urhe3xo5N")
	seeds := [][]byte{[]byte("asset-account"), administrator.Bytes()}
	assetAccount, _, _ := base.FindProgramAddress(seeds, base.AssetExecutorProgramID)
	fmt.Printf(assetAccount.String())
	workerWaller, _ := crypto.AccountFromBase58Key("worker")
	fromMint := base.SolMint
	toMint := common.StrToAddress("4mDwDPrTwsu91qArZ3M3tSTtHVY8VwPJcsGmq9koE22U")
	fromAta, _, _ := base.FindAssociatedTokenAddress(assetAccount, fromMint)
	toAta, _, _ := base.FindAssociatedTokenAddress(assetAccount, toMint)
	inst := NewProxyRaydiumSwapBaseIn(
		common.StrToAddress("HWy1jotHpo6UqeQxx49dpYYdQB8wj9Qk9MdxwjLvDHB8"),
		common.StrToAddress("8Trh3X5RLzAD1jx9bpG7rqZRX6V4NM6uT6mToKcVpnD3"),
		common.StrToAddress("DbQqP6ehDYmeYjcBaMRuA8tAJY1EjDUz9DpwSLjaQqfC"),
		common.StrToAddress("4iTD5j5yDnQrSMaSdQTxQPuFyKW233ivCNeSgJEjVzBs"),
		common.StrToAddress("4as8thCQhz4fxyQKNyBFJKepM7EXB1URoC8Jy1EhupF4"),
		common.StrToAddress("Cy8oppH5VrabsKixfKPuhnMaMCq4ZmXocFoMRJmwJ2gm"),
		common.StrToAddress("3429rpbuU4SbkuABk3FGNsRnuqytHWuAL6F7xxpkRLzi"),
		common.StrToAddress("EoTcMgcDRTJVZDMZWBoU6rhYHZfkNTVEAfz3uUJRcYGj"),
		common.StrToAddress("9VWaSwHo6ivmwTjUff27xfgCUNZ1e1B49P8F2EySho1H"),
		common.StrToAddress("6155Bw4kZubxRDbdA93atVybY57ExeeAtRMbDDmvfTR4"),
		common.StrToAddress("2EyJxMUWbbpCFTv1uQ7kzmmbXgjFmxz7SAvodgiDby27"),
		common.StrToAddress("5R8tGzZotqyXhQ4FhqPS63s6kPZPo9RHiqjY747kgTPF"),
		common.StrToAddress("CNLHQJ1n3H8CJyNywBDPDb4FWW1qXAd4nFagyQgSeSY6"),
		common.StrToAddress("DKtZzRuKom6sLm2wh93xp4eJAhiuSeHsd98x2CUYcWhR"),
		common.StrToAddress("9VWaSwHo6ivmwTjUff27xfgCUNZ1e1B49P8F2EySho1H"),
		fromAta,
		toAta,
		assetAccount,
		administrator,
		workerWaller.Address,
		toMint,
		1e8,
		0,
		true,
	).Build()

	latestHash, err := c.GetLatestBlockhash(ctx)
	if err != nil {
		println("get latest blockHash err:", err.Error())
	}

	transaction, err := types.NewTransaction([]types.Instruction{inst}, latestHash.LastBlock.Blockhash, workerWaller.Address)
	if err != nil {
		fmt.Println("create transaction err:", err)
		return
	}

	signErr := transaction.Sign([]crypto.Account{workerWaller})

	if signErr != nil {
		fmt.Println("sign error:", signErr)
		return
	}

	binary, err := transaction.MarshalBinary()

	fmt.Printf("十六进制格式: %X\n", binary)

	res, err := c.SendTransaction(ctx, binary)
	if err != nil {
		println(err.Error())
	}
	println(res.String())

}

func TestSyncNativeAssetAta(t *testing.T) {
	var (
		c   = newClient()
		ctx = context.Background()
	)
	administrator := common.StrToAddress("6tcU7PRfuSectsCjpJNua2Dc1ygkY41sMEbXdb3R16Lp")
	workerWaller, _ := crypto.AccountFromBase58Key("worker")
	seeds := [][]byte{[]byte("asset-account"), administrator.Bytes()}
	assetAccount, _, _ := base.FindProgramAddress(seeds, base.AssetExecutorProgramID)
	mint := base.SolMint
	warpSolAta, _, _ := base.FindAssociatedTokenAddress(assetAccount, base.SolMint)

	instruction := NewSyncNativeAssetAtaInstruction(assetAccount, warpSolAta, mint, administrator, workerWaller.Address).Build()

	latestHash, err := c.GetLatestBlockhash(ctx)
	if err != nil {
		println("get latest blockHash err:", err.Error())
	}

	transaction, err := types.NewTransaction([]types.Instruction{instruction}, latestHash.LastBlock.Blockhash, workerWaller.Address)
	if err != nil {
		fmt.Println("create transaction err:", err)
		return
	}

	signErr := transaction.Sign([]crypto.Account{workerWaller})

	if signErr != nil {
		fmt.Println("sign error:", signErr)
		return
	}

	binary, err := transaction.MarshalBinary()

	fmt.Printf("十六进制格式: %X\n", binary)

	res, err := c.SendTransaction(ctx, binary)
	if err != nil {
		println(err.Error())
	}
	println(res.String())

}

func TestTransferLamportsToAta(t *testing.T) {
	var (
		c   = newClient()
		ctx = context.Background()
	)
	administrator := common.StrToAddress("6tcU7PRfuSectsCjpJNua2Dc1ygkY41sMEbXdb3R16Lp")
	workerWaller, _ := crypto.AccountFromBase58Key("worker")
	seeds := [][]byte{[]byte("asset-account"), administrator.Bytes()}
	assetAccount, _, _ := base.FindProgramAddress(seeds, base.AssetExecutorProgramID)

	warpSolAta, _, _ := base.FindAssociatedTokenAddress(assetAccount, base.SolMint)
	fmt.Println(warpSolAta)
	mint := base.SolMint
	amount := 2352480
	instruction := NewTransferLamportsToAtaInstruction(assetAccount, warpSolAta, mint, administrator, workerWaller.Address, uint64(amount)).Build()

	latestHash, err := c.GetLatestBlockhash(ctx)
	if err != nil {
		println("get latest blockHash err:", err.Error())
	}

	transaction, err := types.NewTransaction([]types.Instruction{instruction}, latestHash.LastBlock.Blockhash, workerWaller.Address)
	if err != nil {
		fmt.Println("create transaction err:", err)
		return
	}

	signErr := transaction.Sign([]crypto.Account{workerWaller})

	if signErr != nil {
		fmt.Println("sign error:", signErr)
		return
	}

	binary, err := transaction.MarshalBinary()

	fmt.Printf("十六进制格式: %X\n", binary)

	res, err := c.SendTransaction(ctx, binary)
	if err != nil {
		println(err.Error())
	}
	println(res.String())
}

/*
*

	警告:
	1.同一个资产账户不能初始化两次.
	2.多签A、B、C请勿设置为同一个账户，不然无法通过投票！！！！
	3.超级管理员和多签账户一旦设置就无法更改！！！
	4.资产账户一旦开启则无法进行关闭处理
	5.同一个FeePay无法初始两个资产账户，资产账户根据FeePay的公钥 + 特定算法计算出固定地址.
*/
func TestInit(t *testing.T) {
	var (
		c   = newClient()
		ctx = context.Background()
	)

	wallet, _ := crypto.AccountFromBase58Key("admin")
	seeds := [][]byte{[]byte("asset-account"), wallet.Address.Bytes()}
	assetAccount, _, _ := base.FindProgramAddress(seeds, base.AssetExecutorProgramID)

	var administrator = wallet.Address
	var signerA = wallet.Address
	var signerB = common.StrToAddress("7tEoehCWALh4CaVaF6d1RyF6HjVj1Z4x33qCpGChoUw")
	var signerC = common.StrToAddress("HgHHmVsCk3ddS1RdS8V7svy7fkzbr9eeU8RUJhKroq1J")
	var pw = common.StrToAddress("7tEoehCWALh4CaVaF6d1RyF6HjVj1Z4x33qCpGChoUw")
	var sw = common.StrToAddress("HgHHmVsCk3ddS1RdS8V7svy7fkzbr9eeU8RUJhKroq1J")

	fmt.Printf(assetAccount.String())

	instruction := NewInitInstruction(
		administrator,
		signerA,
		signerB,
		signerC,
		pw,
		sw,
		assetAccount,
		wallet.Address,
	).Build()

	latestHash, err := c.GetLatestBlockhash(ctx)
	if err != nil {
		println("get latest blockHash err:", err.Error())
	}

	transaction, err := types.NewTransaction([]types.Instruction{instruction}, latestHash.LastBlock.Blockhash, wallet.Address)
	if err != nil {
		fmt.Println("create transaction err:", err)
		return
	}

	signErr := transaction.Sign([]crypto.Account{wallet})

	if signErr != nil {
		fmt.Println("sign error:", signErr)
		return
	}

	binary, err := transaction.MarshalBinary()

	fmt.Printf("十六进制格式: %X\n", binary)

	res, err := c.SendTransaction(ctx, binary)
	if err != nil {
		println(err.Error())
	}
	println(res.String())

}
