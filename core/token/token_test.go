package token

import (
	"context"
	"fmt"
	"testing"

	"github.com/cielu/go-solana"
)

func newClient() *solana.Client {
	rpcUrl := "https://mainnet-beta.solana.com"
	c, err := solana.Dial(rpcUrl)
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

	payer := solana.StrToPublicKey("F8HCC3DyoR6KN9SSK9NL1V6weRgsEvp8hjL26EnTxNTF")

	instruction := NewTransferCheckedInstruction(
		1e9,
		9,
		solana.StrToPublicKey("BZYExy8yxFZF6jTp4h7X98dPLBcbQDFhvHXPdTjDb2ag"),
		solana.StrToPublicKey("6vG61wtqP7aRgabnECQ2pYBHToJEmPtafvQrxYwmqsAL"),
		solana.StrToPublicKey("EXC6EAnN7HMXbTWomY6j7tQZY1cfZ52LRJpwZ6i3CY66"),
		payer,
		[]solana.PublicKey{payer},
	).Build()

	// core.BeautifyConsole("instruction", instruction)

	latestHash, err := c.GetLatestBlockhash(ctx)
	if err != nil {
		println("get latest blockHash err:", err)
	}

	transaction, err := solana.NewTransaction([]solana.Instruction{instruction}, latestHash.LastBlock.Blockhash, payer)
	if err != nil {
		fmt.Println("create transaction err:", err)
		return
	}

	key, _ := solana.AccountFromBase58Key("3HE29Pg2c2tjbCkVxJpDKhLZuqPLEfoeF3gwjE8MTP3WzvQmLFCxHtKHkGnqNMBPPgFwTWP4vmb9b9a7hGybgtDb")

	_, signErr := transaction.Sign([]solana.Account{key})

	if signErr != nil {
		fmt.Println("sign error:", signErr)
		return
	}

	binary, err := transaction.MarshalBinary()

	res, err := c.SendTransaction(ctx, binary)
	if err != nil {
		println(err.Error())
	}
	println(res.String())
}
