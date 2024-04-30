package native

import (
	"context"
	"fmt"
	"github.com/cielu/go-solana/common"
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

func TestSolTransfer(t *testing.T) {
	var (
		c   = newClient()
		ctx = context.Background()
	)

	instruction := NewTransferInstruction(
		common.StrToAddress("F8HCC3DyoR6KN9SSK9NL1V6weRgsEvp8hjL26EnTxNTF"),
		common.StrToAddress("4DBkBYx6NTSg75BQrQGyDUNbb21j1H7Dt416gdmW4785"),
		1e6,
	).Build()

	latestHash, err := c.GetLatestBlockhash(ctx)
	if err != nil {
		println("get latest blockHash err:", err)
	}

	transaction, err := types.NewTransaction([]types.Instruction{instruction}, latestHash.LastBlock.Blockhash, common.StrToAddress("F8HCC3DyoR6KN9SSK9NL1V6weRgsEvp8hjL26EnTxNTF"))
	key, _ := crypto.AccountFromBase58Key("3HE29Pg2c2tjbCkVxJpDKhLZuqPLEfoeF3gwjE8MTP3WzvQmLFCxHtKHkGnqNMBPPgFwTWP4vmb9b9a7hGybgtDb")
	signErr := transaction.Sign([]crypto.Account{
		key,
	})
	if signErr != nil {
		fmt.Println("signErr:", signErr)
	}

	binary, err := transaction.MarshalBinary()
	res, err := c.SendTransaction(ctx, common.SolData{binary, "base58"})
	if err != nil {
		println(err.Error())
	}
	println(res.String())

}
