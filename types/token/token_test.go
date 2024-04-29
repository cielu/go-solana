package token

import (
	"context"
	"fmt"
	"github.com/cielu/go-solana/common"
	"github.com/cielu/go-solana/core"
	"github.com/cielu/go-solana/crypto"
	"github.com/cielu/go-solana/solclient"
	"github.com/cielu/go-solana/types"
	"testing"
)

func newClient() *solclient.Client {
	rpcUrl := "https://delicate-capable-wish.solana-devnet.quiknode.pro/e48425abfdab96e8263779f8e3334e4a5da10696/"
	c, err := solclient.Dial(rpcUrl)
	if err != nil {
		panic("Dial rpc endpoint failed")
	}
	return c
}

func TestTokenTransfer(t *testing.T) {
	var (
		c   = newClient()
		ctx = context.Background()
	)

	payer := common.StrToAddress("F8HCC3DyoR6KN9SSK9NL1V6weRgsEvp8hjL26EnTxNTF")

	instruction := NewTransferCheckedInstruction(
		1e9,
		9,
		common.StrToAddress("BZYExy8yxFZF6jTp4h7X98dPLBcbQDFhvHXPdTjDb2ag"),
		common.StrToAddress("6vG61wtqP7aRgabnECQ2pYBHToJEmPtafvQrxYwmqsAL"),
		common.StrToAddress("EXC6EAnN7HMXbTWomY6j7tQZY1cfZ52LRJpwZ6i3CY66"),
		payer,
		[]common.Address{payer},
	).Build()

	core.BeautifyConsole("instruction", instruction)

	latestHash, err := c.GetLatestBlockhash(ctx)
	if err != nil {
		println("get latest blockHash err:", err)
	}

	transaction, err := types.NewTransaction([]types.Instruction{instruction}, latestHash.LastBlock.Blockhash, payer)
	if err != nil {
		fmt.Println("create transaction err:", err)
		return
	}

	key, _ := crypto.AccountFromBase58Key("3HE29Pg2c2tjbCkVxJpDKhLZuqPLEfoeF3gwjE8MTP3WzvQmLFCxHtKHkGnqNMBPPgFwTWP4vmb9b9a7hGybgtDb")

	signErr := transaction.Sign([]crypto.Account{key})

	if signErr != nil {
		fmt.Println("sign error:", signErr)
		return
	}

	binary, err := transaction.MarshalBinary()

	res, err := c.SendTransaction(ctx, common.SolData{RawData: binary})
	if err != nil {
		println(err.Error())
	}
	println(res.String())
}
