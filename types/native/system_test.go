package native

import (
	"context"
	"fmt"
	"github.com/cielu/go-solana/common"
	"github.com/cielu/go-solana/crypto"
	"github.com/cielu/go-solana/solclient"
	"github.com/cielu/go-solana/types"
	computebudget "github.com/cielu/go-solana/types/compute-budget"
	"testing"
)

func newClient() *solclient.Client {
	rpcUrl := "https://api.mainnet-beta.solana.com"
	c, err := solclient.Dial(rpcUrl)
	if err != nil {
		panic("Dial rpc endpoint failed")
	}
	return c
}

func TestSolTransfer(t *testing.T) {
	var (
		c        = newClient()
		ctx      = context.Background()
		execInst []types.Instruction
	)
	//
	setLimitInst := computebudget.NewSetComputeUnitLimitInstruction(1000000)
	execInst = append(execInst, setLimitInst.Build())

	setPriceInst := computebudget.NewSetComputeUnitPriceInstruction(10000)
	execInst = append(execInst, setPriceInst.Build())

	transferInst := NewTransferInstruction(
		common.StrToAddress("EfgnVEwyeeFLZyZ4nnnzZtqV6B3DhdtXFNsGSzdti9ZN"),
		common.StrToAddress("6XViKPqw7t47tZz8UJR1bJFVzxjnQbuKtN2TBgnfZmo4"),
		1e1,
	)
	execInst = append(execInst, transferInst.Build())

	latestHash, err := c.GetLatestBlockhash(ctx)
	if err != nil {
		println("Get latest blockHash err:", err)
	}

	transaction, err := types.NewTransaction(execInst, latestHash.LastBlock.Blockhash, common.StrToAddress("EfgnVEwyeeFLZyZ4nnnzZtqV6B3DhdtXFNsGSzdti9ZN"))

	key, _ := crypto.AccountFromBase58Key("")
	// /
	signTx, err := transaction.Sign([]crypto.Account{key})
	if err != nil {
		fmt.Println("signErr:", err)
	}
	res, err := c.SendTransaction(ctx, signTx)
	if err != nil {
		println(err.Error())
	}
	println(res.String())

}
