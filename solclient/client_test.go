package solclient

import (
	"context"
	"go-solana/common"
	"go-solana/core"
	"go-solana/rpc"
	"go-solana/types"
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
	c, err := Dial(rpc.DevnetRPCEndpoint)
	if err != nil {
		panic("Dial rpc endpoint failed")
	}
	return c
}

func getDefaultRpcCfg() types.RpcCommitmentAndMinSlotCfg {
	return types.RpcCommitmentAndMinSlotCfg{
		Commitment: types.RpcCommitmentConfirmed,
	}
}

func TestClient_GetAccountInfo(t *testing.T) {
	//
	var (
		c = newClient()
		ctx = context.Background()
	)
	account := common.Base58ToAddress("So11111111111111111111111111111111111111112")
	res, err := c.GetAccountInfo(ctx, account)
	if err != nil {
		t.Error("GetAccountInfo Failed: %w", err)
	}
	core.BeautifyConsole("AccountInfo:", res)
}

func TestClient_GetBalance(t *testing.T) {
	var (
		c = newClient()
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
		c = newClient()
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
		c = newClient()
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
		c = newClient()
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
		c = newClient()
		ctx = context.Background()
	)
	res, err := c.GetBlockProduction(ctx)
	if err != nil {
		t.Error("GetBlockProduction Failed: %w", err)
	}
	core.BeautifyConsole("GetBlockProduction:", res)
}
