package solclient

import (
	"context"
	"github.com/cielu/go-solana/common"
	"github.com/cielu/go-solana/core"
	"github.com/cielu/go-solana/types"
	"testing"
)

func TestClient_AccountSubscribe(t *testing.T) {
	var (
		c             = newClient()
		ctx           = context.Background()
		accountNotify = make(chan *types.SubAccountInfo)
	)
	account := common.Base58ToAddress("CM78CPUeXjn8o3yroDHxUtKsZZgoy4GPkPPXfouKNH12")
	//
	res, err := c.AccountSubscribe(ctx, accountNotify, account)
	if err != nil {
		t.Error("AccountSubscribe Failed: %w", err)
	}
	core.BeautifyConsole("AccountSubscribe:", res)
}
