package system

import (
	"context"
	"fmt"
	"testing"

	"github.com/cielu/go-solana"
	computebudget "github.com/cielu/go-solana/core/compute-budget"
)

func newClient() *solana.Client {
	rpcUrl := "https://api.mainnet-beta.solana.com"
	c, err := solana.Dial(rpcUrl)
	if err != nil {
		panic("Dial rpc endpoint failed")
	}
	return c
}

func TestSolTransfer(t *testing.T) {
	var (
		c        = newClient()
		ctx      = context.Background()
		execInst []solana.Instruction
	)
	//
	setLimitInst := computebudget.NewSetComputeUnitLimitInstruction(1000000)
	execInst = append(execInst, setLimitInst.Build())

	setPriceInst := computebudget.NewSetComputeUnitPriceInstruction(10000)
	execInst = append(execInst, setPriceInst.Build())

	transferInst := NewTransferInstruction(
		solana.StrToPublicKey("EfgnVEwyeeFLZyZ4nnnzZtqV6B3DhdtXFNsGSzdti9ZN"),
		solana.StrToPublicKey("6XViKPqw7t47tZz8UJR1bJFVzxjnQbuKtN2TBgnfZmo4"),
		1e1,
	)
	execInst = append(execInst, transferInst.Build())

	latestHash, err := c.GetLatestBlockhash(ctx)
	if err != nil {
		println("Get latest blockHash err:", err)
	}

	transaction, err := solana.NewTransaction(execInst, latestHash.LastBlock.Blockhash, solana.StrToPublicKey("EfgnVEwyeeFLZyZ4nnnzZtqV6B3DhdtXFNsGSzdti9ZN"))

	key, _ := solana.AccountFromBase58Key("")
	// /
	signTx, err := transaction.Sign([]solana.Account{key})
	if err != nil {
		fmt.Println("signErr:", err)
	}
	res, err := c.SendTransaction(ctx, signTx)
	if err != nil {
		println(err.Error())
	}
	println(res.String())

}
