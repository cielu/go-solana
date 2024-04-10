## go-solana

#### Example

```go
package main

import (
	"context"
	"fmt"
	"github.com/cielu/go-solana/common"
	"github.com/cielu/go-solana/core"
	"github.com/cielu/go-solana/rpc"
	"github.com/cielu/go-solana/solclient"
	"math/big"
)

func main() {
	var (
		ctx = context.Background()
	)
	c, err := solclient.Dial(rpc.DevnetRPCEndpoint)
	// err
	if err != nil {
		panic("Failed Dial Solana RPC")
	}
	// account
	account := common.Base58ToAddress("So11111111111111111111111111111111111111112")
	// get AccountInfo
	res, err := c.GetAccountInfo(ctx, account)
	// has err
	if err != nil {
		fmt.Errorf("GetAccountInfo Failed: %w", err)
		return
	}
	
	core.BeautifyConsole("AccountInfo:", res)

	// request Airdrop
	signature, err := c.RequestAirdrop(ctx, account, big.NewInt(1000000000))
	// has err
	if err != nil {
		fmt.Errorf("RequestAirdrop Failed: %w", err)
		return
	}

	core.BeautifyConsole("signature:", signature)
}



```


