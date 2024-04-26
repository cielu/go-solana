package associated_token_account

import (
	"context"
	"fmt"
	"github.com/cielu/go-solana/common"
	"github.com/cielu/go-solana/solclient"
	"github.com/cielu/go-solana/types"
	"github.com/cielu/go-solana/types/account"
	"log"
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

func Test_CreateAssociationAccount(t *testing.T) {
	var (
		c           = newClient()
		ctx         = context.Background()
		feePayer, _ = types.AccountFromBase58("3HE29Pg2c2tjbCkVxJpDKhLZuqPLEfoeF3gwjE8MTP3WzvQmLFCxHtKHkGnqNMBPPgFwTWP4vmb9b9a7hGybgtDb")
		auth        = common.StrToAddress("Fjw2S4TzcxCkJdfho6mhseDCPKQ9QrxUf2aZiGpi7ar3")
		mint        = common.StrToAddress("6vG61wtqP7aRgabnECQ2pYBHToJEmPtafvQrxYwmqsAL")
	)
	ata, _, err := account.FindAssociatedTokenAddress(auth, mint)
	if err != nil {
		log.Fatalf("find ata error, err: %v", err)
	}
	fmt.Println("ata:", ata.String())

	latestHash, err := c.GetLatestBlockhash(ctx)
	if err != nil {
		println("get latest blockHash err:", err)
	}
	tx, err := types.NewTransaction(types.NewTransactionParam{
		Message: types.NewMessage(types.NewMessageParam{
			FeePayer:        feePayer.PublicKey,
			RecentBlockhash: latestHash.LastBlock.Blockhash,
			Instructions: []types.Instruction{
				Create(CreateParam{
					Funder:                 feePayer.PublicKey,
					Owner:                  auth,
					Mint:                   mint,
					AssociatedTokenAccount: ata,
				}),
			},
		}),
		Signatures: []types.Account{feePayer},
	})
	serialize, err := tx.Serialize()
	if err != nil {
		println("err")
	}
	res, err := c.SendTransaction(ctx, common.SolData{serialize, "base58"})
	if err != nil {
		println(err.Error())
	}
	println(res.String())
}
