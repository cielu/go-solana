package solclient

import (
	"context"
	"github.com/cielu/go-solana/common"
	"github.com/cielu/go-solana/types"
	"github.com/cielu/go-solana/types/system"
	"github.com/cielu/go-solana/types/token"
	"testing"
)

func Test_TransferTokenChecked(t *testing.T) {
	var (
		c   = newClient()
		ctx = context.Background()
	)
	latestHash, err := c.GetLatestBlockhash(ctx)
	if err != nil {
		println("get latest blockHash err:", err)
	}
	var feePayer, _ = types.AccountFromBase58("3HE29Pg2c2tjbCkVxJpDKhLZuqPLEfoeF3gwjE8MTP3WzvQmLFCxHtKHkGnqNMBPPgFwTWP4vmb9b9a7hGybgtDb")
	var owner, _ = types.AccountFromBase58("3HE29Pg2c2tjbCkVxJpDKhLZuqPLEfoeF3gwjE8MTP3WzvQmLFCxHtKHkGnqNMBPPgFwTWP4vmb9b9a7hGybgtDb")
	var from = common.StrToAddress("BZYExy8yxFZF6jTp4h7X98dPLBcbQDFhvHXPdTjDb2ag")
	var to = common.StrToAddress("7dCNBJ3qtkU4z2vCNk2MsNJqJ7CSWdqCxn7Zogpp9Sr")
	var mint = common.StrToAddress("6vG61wtqP7aRgabnECQ2pYBHToJEmPtafvQrxYwmqsAL")

	transaction, err := types.NewTransaction(
		types.NewTransactionParam{
			Signatures: []types.Account{feePayer, owner},
			Message: types.NewMessage(types.NewMessageParam{
				FeePayer:        feePayer.PublicKey,
				RecentBlockhash: latestHash.LastBlock.Blockhash,
				Instructions: []types.Instruction{
					token.TransferChecked(token.TransferCheckedParam{
						From:     from,
						To:       to,
						Mint:     mint,
						Auth:     feePayer.PublicKey,
						Signers:  []common.Address{},
						Amount:   1000000e9,
						Decimals: 9,
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

func Test_TransferToken(t *testing.T) {
	var (
		c   = newClient()
		ctx = context.Background()
	)
	type args struct {
		param token.TransferParam
	}

	latestHash, err := c.GetLatestBlockhash(ctx)
	if err != nil {
		println("get latest blockHash err:", err)
	}
	var feePayer, _ = types.AccountFromBase58("2KLTas5hUpiFNyT3tdDGjuFuJbYcm1bMsqgiCFJeU6JdmgAbAykqNf2jqSFfTEP9ATz5wg3JckgH5H19L8V9r6Sb")
	var owner, _ = types.AccountFromBase58("2KLTas5hUpiFNyT3tdDGjuFuJbYcm1bMsqgiCFJeU6JdmgAbAykqNf2jqSFfTEP9ATz5wg3JckgH5H19L8V9r6Sb")
	var from = common.StrToAddress("BuAeYkSfGyTpHszavYX4VwtwvkiAiB8S3gwar6sxymrL")
	var to = common.StrToAddress("4X9nZF4gNqATtcS26CoSm9oZVoZwookgHuSCZT2dPYUa")

	transaction, err := types.NewTransaction(
		types.NewTransactionParam{
			Signatures: []types.Account{feePayer, owner},
			Message: types.NewMessage(types.NewMessageParam{
				FeePayer:        feePayer.PublicKey,
				RecentBlockhash: latestHash.LastBlock.Blockhash,
				Instructions: []types.Instruction{
					token.Transfer(token.TransferParam{
						From:    from,
						To:      to,
						Auth:    owner.PublicKey,
						Signers: []common.Address{},
						Amount:  1,
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

func TestClient_SendTransaction(t *testing.T) {
	var (
		c   = newClient()
		ctx = context.Background()
	)

	blockhash, err2 := c.GetLatestBlockhash(ctx)
	if err2 != nil {
		println("error")
	}
	var feePayer, _ = types.AccountFromBase58("2KLTas5hUpiFNyT3tdDGjuFuJbYcm1bMsqgiCFJeU6JdmgAbAykqNf2jqSFfTEP9ATz5wg3JckgH5H19L8V9r6Sb")
	var alice, _ = types.AccountFromBase58("2KLTas5hUpiFNyT3tdDGjuFuJbYcm1bMsqgiCFJeU6JdmgAbAykqNf2jqSFfTEP9ATz5wg3JckgH5H19L8V9r6Sb")

	transaction, err := types.NewTransaction(
		types.NewTransactionParam{
			Signatures: []types.Account{feePayer, feePayer},
			Message: types.NewMessage(types.NewMessageParam{
				FeePayer:        feePayer.PublicKey,
				RecentBlockhash: blockhash.LastBlock.Blockhash,
				Instructions: []types.Instruction{
					system.Transfer(system.TransferParam{
						From:   alice.PublicKey,
						To:     common.StrToAddress("CtiyYm2pRNwNKGcUPC7h9zHpdoqQnBEt9vPWQbZU9RCD"),
						Amount: 1e6, // 0.01 SOL
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
