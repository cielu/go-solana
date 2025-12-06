package solana

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
)

func newClient() *Client {
	//
	rpcUrl := "https://api.mainnet-beta.solana.com"
	// dial rpc
	c, err := Dial(rpcUrl)
	if err != nil {
		panic("Dial rpc endpoint failed")
	}
	c.SetDebug(true)
	return c
}

func getDefaultRpcCfg() RpcCommitmentWithMinSlotCfg {
	return RpcCommitmentWithMinSlotCfg{
		Commitment: RpcCommitmentConfirmed,
	}
}

func beautifulConsole(title string, content any) {
	// MarshalIndent
	jsonData, _ := json.MarshalIndent(content, "", "    ")
	// print data
	fmt.Println(title, string(jsonData))
}

func TestClient_GetAccountInfo(t *testing.T) {
	//
	var (
		c   = newClient()
		ctx = context.Background()
	)
	account := Base58ToPublicKey("89jFVQvaVeLs3h35BC3P592UKF5xB8fK2eyjdBvSFeWw")

	// res, err := c.GetAccountInfo(ctx, account)

	accounts := []PublicKey{account}
	res, err := c.GetMultipleAccounts(ctx, accounts)

	if err != nil {
		t.Error("GetAccountInfo Failed: %w", err)
	}
	// Define Data detail

	// var dataDetail AccountData

	// encodbin.NewBinDecoder(res.Accounts[0].Data).Decode(&dataDetail)
	//
	// core.BeautifyConsole("dataDetail", dataDetail)

	beautifulConsole("AccountInfo:", res)
}

func TestClient_GetBalance(t *testing.T) {
	var (
		c   = newClient()
		ctx = context.Background()
	)
	account := Base58ToPublicKey("So11111111111111111111111111111111111111112")
	res, err := c.GetBalance(ctx, account)
	if err != nil {
		t.Error("GetBalance Failed: %w", err)
	}
	beautifulConsole("AccountBalance:", res)
}

func TestClient_GetBlock(t *testing.T) {
	var (
		c   = newClient()
		ctx = context.Background()
	)
	res, err := c.GetBlock(ctx, 378542246)
	if err != nil {
		t.Error("GetBlock Failed: %w", err)
	}
	// range block transactions
	for i, tx := range res.Transactions {
		fmt.Println("Tx index:", i+1)
		fmt.Println("Signature:", tx.Transaction.Signatures[0])
		// foreach Instruction
		for i2, instruction := range tx.Transaction.Message.Instructions {
			fmt.Println("	Instruction Index:", i2)
			fmt.Println("	Program ID:", instruction.ProgramIDIndex)
			fmt.Println("	Data:", instruction.Data)
		}
		fmt.Println("===================================================")
	}
	beautifulConsole("BlockInfo:", res)
}

func TestClient_GetBlockCommitment(t *testing.T) {
	var (
		c   = newClient()
		ctx = context.Background()
	)
	res, err := c.GetBlockCommitment(ctx, 256778078)
	if err != nil {
		t.Error("GetBlockCommitment Failed: %w", err)
	}
	beautifulConsole("BlockCommitment:", res)
}

func TestClient_GetBlockHeight(t *testing.T) {
	var (
		c   = newClient()
		ctx = context.Background()
	)
	res, err := c.GetBlockHeight(ctx, getDefaultRpcCfg())
	if err != nil {
		t.Error("GetBlockHeight Failed: %w", err)
	}
	beautifulConsole("GetBlockHeight:", res)
}

func TestClient_GetBlockProduction(t *testing.T) {
	var (
		c   = newClient()
		ctx = context.Background()
	)
	res, err := c.GetBlockProduction(ctx)
	if err != nil {
		t.Error("GetBlockProduction Failed: %w", err)
	}
	beautifulConsole("GetBlockProduction:", res)
}

func TestClient_GetBlocks(t *testing.T) {
	var (
		c   = newClient()
		ctx = context.Background()
	)
	res, err := c.GetBlocks(ctx, 378542246)
	if err != nil {
		t.Error("GetBlocks Failed: %w", err)
	}
	beautifulConsole("GetBlocks:", res)
}

func TestClient_GetEpochInfo(t *testing.T) {
	var (
		c   = newClient()
		ctx = context.Background()
	)
	res, err := c.GetEpochInfo(ctx)
	if err != nil {
		t.Error("GetEpochInfo Failed: %w", err)
	}
	beautifulConsole("GetEpochInfo:", res)
}

func TestClient_GetEpochSchedule(t *testing.T) {
	// GetEpochSchedule
	var (
		c   = newClient()
		ctx = context.Background()
	)
	res, err := c.GetEpochSchedule(ctx)
	if err != nil {
		t.Error("GetEpochSchedule Failed: %w", err)
	}
	beautifulConsole("GetEpochSchedule:", res)
}

func TestClient_GetFeeForMessage(t *testing.T) {
	// GetFeeForMessage
	var (
		c   = newClient()
		ctx = context.Background()
		msg = "AQABAgIAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAEBAQAA"
	)
	res, err := c.GetFeeForMessage(ctx, msg)
	if err != nil {
		t.Error("GetFeeForMessage Failed: %w", err)
	}
	beautifulConsole("GetFeeForMessage:", res)
}

func TestClient_GetProgramAccounts(t *testing.T) {
	// GetProgramAccounts
	var (
		c   = newClient()
		ctx = context.Background()
	)
	// memcmp
	memcmp := map[string]interface{}{
		"offset": 4,
		"bytes":  "3Mc6vR",
	}
	minSlot := uint64(378543712)
	//
	cfg := RpcCombinedCfg{
		MinContextSlot: &minSlot,
		Filter: []map[string]interface{}{
			{
				"dataSize": 17,
			},
			{
				"memcmp": memcmp,
			},
		},
	}
	account := Base58ToPublicKey("CAMMCzo5YL8w4VFF8KVHrK22GGUsp5VTaW7grrKgrWqK")
	//
	res, err := c.GetProgramAccounts(ctx, account, cfg)
	if err != nil {
		t.Error("GetProgramAccounts Failed: %w", err)
	}
	beautifulConsole("GetProgramAccounts:", res)
}

func TestClient_GetRecentPerformanceSamples(t *testing.T) {
	var (
		c   = newClient()
		ctx = context.Background()
	)
	res, err := c.GetRecentPrioritizationFees(ctx)
	if err != nil {
		t.Error("Res Failed: %w", err)
	}
	beautifulConsole("Res:", res)
}

func TestClient_GetSignatureStatuses(t *testing.T) {
	var (
		c          = newClient()
		ctx        = context.Background()
		signatures = []Signature{
			Base58ToSignature("4BC9UMSQrLqEbTvcya6Ukt3Lvq1ZzWpCYRZG6ygGCu26mJYQ1uAiQNVVbu3jPz5SBS2oWKmbhNTR3h6x6wyBELS5"),
		}
		// account = Base58ToPublicKey("Vote111111111111111111111111111111111111111")
	)
	res, err := c.GetSignatureStatuses(ctx, signatures, RpcSearchTxHistoryCfg{true})
	// res, err := c.GetSignaturesForPublicKey(ctx, account)
	if err != nil {
		t.Error("Res Failed: %w", err)
	}
	beautifulConsole("Res:", res)
}

func TestClient_GetTokenAccountsByDelegate(t *testing.T) {
	var (
		c   = newClient()
		ctx = context.Background()
	)
	account := Base58ToPublicKey("2gKXoChNdN3LpMijfLdx4AVs62JBseXTHtrqxCSYreWs")

	prog := Base58ToPublicKey("TokenkegQfeZyiNwAJbNbGKPFXCWuBvf9Ss623VQ5DA")
	//
	res, err := c.GetTokenAccountsByDelegate(ctx, account, RpcMintWithProgramID{ProgramId: &prog})
	if err != nil {
		t.Error("Res Failed: %w", err)
	}
	beautifulConsole("Res:", res)
}

func TestClient_GetTokenLargestAccounts(t *testing.T) {
	var (
		c   = newClient()
		ctx = context.Background()
	)

	token := Base58ToPublicKey("3p7U58GR11MnfRuWCBufj9AW3Y7P1x848CWgtECpNQpt")

	res, err := c.GetTokenLargestAccounts(ctx, token)
	if err != nil {
		t.Error("Res Failed: %w", err)
	}
	beautifulConsole("Res:", res)
}

func TestClient_GetTokenSupply(t *testing.T) {
	var (
		c   = newClient()
		ctx = context.Background()
	)

	account := Base58ToPublicKey("J1toso1uCk3RLmjorhTtrVwY9HJ7X8V9yYac6Y7kGCPn")

	res, err := c.GetTokenSupply(ctx, account)
	if err != nil {
		t.Error("Res Failed: %w", err)
	}
	beautifulConsole("Res:", res)
}

func TestClient_GetTransaction(t *testing.T) {
	var (
		c   = newClient()
		ctx = context.Background()
	)

	signature := Base58ToSignature("3rGHV82oPnL2gZphq78bzXPpc4PHRPbQGfKLng4VrS6HYkyCn8CmgMzhmTuchbZZQ1RiUPzrqqgfwGygeeJqRy5c")
	// GetTransaction
	res, err := c.GetTransaction(ctx, signature, RpcGetTransactionCfg{
		Encoding:              EncodingBase64,
		Commitment:            RpcCommitmentConfirmed,
		MaxSupportedTxVersion: 1,
	})
	if err != nil {
		t.Error("Res Failed: %w", err)
	}
	addrTableLookups := res.Transaction.Message.AddressTableLookups
	// utils
	beautifulConsole("table IDs", addrTableLookups.GetTableIDs())
	// // has Address TableLookups
	// if len(addrTableLookups) > 0 {
	// 	for _, address := range addrTableLookups {
	// 		accRes, _ := c.GetAccountInfo(ctx, address.AccountKey, RpcAccountInfoCfg{
	// 			Encoding:   EncodingBase64,
	// 			Commitment: RpcCommitmentConfirmed,
	// 		})
	// 		beautifulConsole("AccInfo:", accRes.AccountInfo)
	// 		// Decode AccRes
	// 		lookupState, _ := DecodeAddressLookupTableState(accRes.AccountInfo.Data.RawData())
	// 		// Decode Address lookups table
	// 		beautifulConsole("Res:", lookupState)
	// 	}
	// }
}

func TestClient_GetVoteAccounts(t *testing.T) {
	//
	var (
		c   = newClient()
		ctx = context.Background()
	)

	res, err := c.GetVoteAccounts(ctx)
	if err != nil {
		t.Error("Res Failed: %w", err)
	}
	beautifulConsole("Res:", res)
}

func TestClient_SendTransaction(t *testing.T) {
	var (
		c   = newClient()
		ctx = context.Background()
	)

	var (
		instrs []Instruction
	)

	recentHash, err := c.GetLatestBlockhash(ctx)
	if err != nil {
		fmt.Println("GetLatestBlockhash Failed:", err)
		return
	}

	payer, _ := GenerateAccount()
	// payer, _ := AccountFromBase58Key("payer private key")

	// instrs = append(instrs, computebudget.NewSetComputeUnitLimitInstruction(1000000).Build())
	// instrs = append(instrs, computebudget.NewSetComputeUnitPriceInstruction(10000).Build())

	// transferInst := system.NewTransferInstruction(
	// 	StrToPublicKey("EfgnVEwyeeFLZyZ4nnnzZtqV6B3DhdtXFNsGSzdti9ZN"),
	// 	StrToPublicKey("6XViKPqw7t47tZz8UJR1bJFVzxjnQbuKtN2TBgnfZmo4"),
	// 	1e1,
	// )

	// instrs = append(instrs, transferInst.Build())
	// return
	tx, err := NewTransaction(instrs, recentHash.LastBlock.Blockhash, payer.PublicKey)

	sigTx, err := tx.Sign([]Account{payer})
	//
	if err != nil {
		fmt.Println("Sign Tx Failed:", err)
		return
	}

	res, err := c.SendTransaction(ctx, sigTx)
	if err != nil {
		t.Error("Res Failed: %w", err)
	}

	beautifulConsole("Res:", res)
}
