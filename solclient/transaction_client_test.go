package solclient

import (
	"context"
	"github.com/cielu/go-solana/common"
	"github.com/cielu/go-solana/types"
	"github.com/cielu/go-solana/types/system"
	"testing"
)

func TestClient_SendTransaction(t *testing.T) {
	var (
		c   = newClient()
		ctx = context.Background()
	)

	blockhash, err2 := c.GetLatestBlockhash(ctx)
	if err2 != nil {
		println("error")
	}

	var feePayer, _ = types.AccountFromBase58("3DjxzhaPzhq4oVM7u3n92UgNa5BtGBoYzDLboZCoYZDbgvuGycAFm4UeqhZqzTJM61zU4nRoBdxekRxDj3duk3W9")
	var alice, _ = types.AccountFromBase58("3DjxzhaPzhq4oVM7u3n92UgNa5BtGBoYzDLboZCoYZDbgvuGycAFm4UeqhZqzTJM61zU4nRoBdxekRxDj3duk3W9")

	transaction, err := types.NewTransaction(
		types.NewTransactionParam{
			Signatures: []types.Account{feePayer, feePayer},
			Message: types.NewMessage(types.NewMessageParam{
				FeePayer:        feePayer.PublicKey,
				RecentBlockhash: blockhash.LastBlock.Blockhash,
				Instructions: []types.Instruction{
					system.Transfer(system.TransferParam{
						From:   alice.PublicKey,
						To:     common.StrToAddress("5qV6Xh7pjHNTsXFRRckUJChVxpNmEPL6X58QXPc6qht9"),
						Amount: 1e5, // 0.01 SOL
					}),
				},
			}),
		},
	)
	serialize, err := transaction.Serialize()
	if err != nil {
		println("err")
	}
	res, err := c.SendTransaction(ctx, common.SolData{serialize, "base58"})
	if err != nil {
		println(err.Error())
	}
	println(res.String())

}
