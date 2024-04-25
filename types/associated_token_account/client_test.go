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
	rpcUrl := "https://small-sly-pool.solana-mainnet.quiknode.pro/83bb6bc144ded159a51a1c28b45ae82ae95053af/"
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
		feePayer, _ = types.AccountFromBase58("2KLTas5hUpiFNyT3tdDGjuFuJbYcm1bMsqgiCFJeU6JdmgAbAykqNf2jqSFfTEP9ATz5wg3JckgH5H19L8V9r6Sb")
		auth, _     = types.AccountFromBase58("3DjxzhaPzhq4oVM7u3n92UgNa5BtGBoYzDLboZCoYZDbgvuGycAFm4UeqhZqzTJM61zU4nRoBdxekRxDj3duk3W9")
		mint        = common.StrToAddress("DezXAZ8z7PnrnRJjz3wXBoRgixCa6xjnB7YaB1pPB263")
	)
	ata, _, err := account.FindAssociatedTokenAddress(auth.PublicKey, mint)
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
					Owner:                  auth.PublicKey,
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
