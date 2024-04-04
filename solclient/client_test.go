package solclient

import (
	"context"
	"github.com/cielu/go-solana/common"
	"github.com/cielu/go-solana/core"
	"github.com/cielu/go-solana/rpc"
	"github.com/cielu/go-solana/types"
	"os"
	"testing"
)

// remove this, if not
func init() {
	os.Setenv("http_proxy", "http://127.0.0.1:7890")
	os.Setenv("https_proxy", "http://127.0.0.1:7890")
	os.Setenv("all_proxy", "socks5://127.0.0.1:7890")
}

func newClient() *Client {
	//
	rpcUrl := rpc.DevnetRPCEndpoint
	//rpcUrl := rpc.MainnetRPCEndpoint
	// dial rpc
	c, err := Dial(rpcUrl)
	if err != nil {
		panic("Dial rpc endpoint failed")
	}
	return c
}

func getDefaultRpcCfg() types.RpcCommitmentWithMinSlotCfg {
	return types.RpcCommitmentWithMinSlotCfg{
		Commitment: types.RpcCommitmentConfirmed,
	}
}

func TestClient_GetAccountInfo(t *testing.T) {
	//
	var (
		c   = newClient()
		ctx = context.Background()
	)
	account := common.Base58ToAddress("So11111111111111111111111111111111111111112")

	//res, err := c.GetAccountInfo(ctx, account)

	accounts := []common.Address{account}
	res, err := c.GetMultipleAccounts(ctx, accounts)

	if err != nil {
		t.Error("GetAccountInfo Failed: %w", err)
	}
	core.BeautifyConsole("AccountInfo:", res)
}

func TestClient_GetBalance(t *testing.T) {
	var (
		c   = newClient()
		ctx = context.Background()
	)
	account := common.Base58ToAddress("So11111111111111111111111111111111111111112")
	res, err := c.GetBalance(ctx, account)
	if err != nil {
		t.Error("GetBalance Failed: %w", err)
	}
	core.BeautifyConsole("AccountBalance:", res)
}

func TestClient_GetBlock(t *testing.T) {
	var (
		c   = newClient()
		ctx = context.Background()
	)
	res, err := c.GetBlock(ctx, 256731099)
	if err != nil {
		t.Error("GetBlock Failed: %w", err)
	}
	core.BeautifyConsole("BlockInfo:", res)
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
	core.BeautifyConsole("BlockCommitment:", res)
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
	core.BeautifyConsole("GetBlockHeight:", res)
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
	core.BeautifyConsole("GetBlockProduction:", res)
}

func TestClient_GetBlocks(t *testing.T) {
	var (
		c   = newClient()
		ctx = context.Background()
	)
	res, err := c.GetBlocks(ctx, 256731099)
	if err != nil {
		t.Error("GetBlocks Failed: %w", err)
	}
	core.BeautifyConsole("GetBlocks:", res)
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
	core.BeautifyConsole("GetEpochInfo:", res)
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
	core.BeautifyConsole("GetEpochSchedule:", res)
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
	core.BeautifyConsole("GetFeeForMessage:", res)
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
	//
	cfg := types.RpcCombinedCfg{
		Filter: []map[string]interface{}{
			{
				"dataSize": 17,
			},
			{
				"memcmp": memcmp,
			},
		},
	}
	account := common.Base58ToAddress("metaqbxxUerdq28cj1RbAWkYQm3ybzjb6a8bt518x1s")
	//
	res, err := c.GetProgramAccounts(ctx, account, cfg)
	if err != nil {
		t.Error("GetProgramAccounts Failed: %w", err)
	}
	core.BeautifyConsole("GetProgramAccounts:", res)
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
	core.BeautifyConsole("Res:", res)
}

func TestClient_GetTokenAccountsByDelegate(t *testing.T) {
	var (
		c   = newClient()
		ctx = context.Background()
	)
	account := common.Base58ToAddress("2gKXoChNdN3LpMijfLdx4AVs62JBseXTHtrqxCSYreWs")
	prog := common.Base58ToAddress("TokenkegQfeZyiNwAJbNbGKPFXCWuBvf9Ss623VQ5DA")
	//
	res, err := c.GetTokenAccountsByDelegate(ctx, account, types.RpcMintWithProgramID{ProgramId: &prog})
	if err != nil {
		t.Error("Res Failed: %w", err)
	}
	core.BeautifyConsole("Res:", res)
}

func TestClient_GetTokenLargestAccounts(t *testing.T) {
	var (
		c   = newClient()
		ctx = context.Background()
	)

	token := common.Base58ToAddress("3p7U58GR11MnfRuWCBufj9AW3Y7P1x848CWgtECpNQpt")

	res, err := c.GetTokenLargestAccounts(ctx, token)
	if err != nil {
		t.Error("Res Failed: %w", err)
	}
	core.BeautifyConsole("Res:", res)
}


func TestClient_GetTokenSupply(t *testing.T) {
	var (
		c   = newClient()
		ctx = context.Background()
	)

	account := common.Base58ToAddress("3p7U58GR11MnfRuWCBufj9AW3Y7P1x848CWgtECpNQpt")

	res, err := c.GetTokenSupply(ctx, account)
	if err != nil {
		t.Error("Res Failed: %w", err)
	}
	core.BeautifyConsole("Res:", res)
}

func TestClient_GetTransaction(t *testing.T) {
	var (
		c   = newClient()
		ctx = context.Background()
	)

	signature := common.Base58ToSignature("4mQwRJD8yjNf7wMzTMsKX6fkzK2uLMk9UZRCKnjLRaePkd4LrYsGJPagQ4pqYLeT4MLV1FCRMbApyqT2b5mRUhPn")

	res, err := c.GetTransaction(ctx, signature)
	if err != nil {
		t.Error("Res Failed: %w", err)
	}
	core.BeautifyConsole("Res:", res)
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
	core.BeautifyConsole("Res:", res)
}







