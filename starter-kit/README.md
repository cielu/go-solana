# solana-go

> **Pre-1.0. API unstable. Do not use in production yet.**
> Follow [ROADMAP.md](./ROADMAP.md) for progress.

A fast, clean Solana client for Go.

## Why another Solana Go SDK

The existing standard, [`gagliardetto/solana-go`](https://github.com/gagliardetto/solana-go),
has served the ecosystem well but carries years of accumulated design
debt: reflection-heavy decoding, `math/big.Int` everywhere, legacy
`...interface{}` APIs, and no first-class Anchor IDL generator. We
think the Go Solana community deserves better, so this project is a
**clean-room rewrite** built on three principles:

1. **Performance first, benchmarked, honest.** Every public performance
   claim is backed by a checked-in benchmark on real mainnet fixtures.
   See [CLAIMS.md](./CLAIMS.md).
2. **API that is a pleasure to use.** Typed functional options
   everywhere, `context.Context` plumbed through every I/O path, error
   messages that name the failing RPC method. No `...interface{}`, no
   silent argument-dropping, no library-code `panic`s.
3. **Ecosystem moat via IDL codegen.** `cmd/anchoridl-gen` turns any
   Anchor IDL into a small, idiomatic Go binding in seconds. Program
   bindings live in a sibling repository and track upstream releases
   within days, not months.

## Status

This repository is in **Phase 0 — bootstrap**. The public surface is
being built; there is no stable `import` path yet. Watch
[ROADMAP.md](./ROADMAP.md) for phase transitions.

## Design goals

- **Go 1.22+**, generics used where they earn their keep.
- **Apache-2.0**, no LGPL-tainted vendored code.
- **Zero reflection** in hot decode paths.
- **Zero `panic`s** in library code.
- **Deterministic, offline unit tests** plus `solana-test-validator`
  integration tests behind a build tag.
- **JSON decoding** uses a benchmark-chosen high-performance library
  (decision recorded in `CLAIMS.md`).

## Non-goals

- Feature parity with `gagliardetto/solana-go` on day one. We ship
  fewer things done better and catch up over time.
- Bundling DeFi protocol bindings in this repository. Those live in
  `cielu/solana-programs` (to be created), generated from their
  upstream IDLs.
- Browser wallet adapters. Not a server-side-SDK concern.
- Historical indexers. Different problem shape.

## Quickstart

> _Coming in Phase A. The snippet below is aspirational._

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/cielu/solana-go"
)

func main() {
    client, err := solana.DialContext(context.Background(), solana.MainnetBetaEndpoint)
    if err != nil {
        log.Fatal(err)
    }
    defer client.Close()

    wsol := solana.MustPublicKey("So11111111111111111111111111111111111111112")
    info, err := client.GetAccountInfo(context.Background(), wsol)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("wSOL owner: %s\n", info.Value.Owner)
}
```

## Benchmarks

> _Coming at the end of Phase A. Placeholder table in
> [CLAIMS.md](./CLAIMS.md). No numbers will appear in this README
> until they are reproducible on a reference machine from a
> checked-in benchmark._

## Installation

```bash
go get github.com/cielu/solana-go
```

## Documentation

- [`CLAUDE.md`](./CLAUDE.md) — internal memory file for AI-assisted
  development. Read this if you want to understand every design
  decision.
- [`ROADMAP.md`](./ROADMAP.md) — public phase plan.
- [`CLAIMS.md`](./CLAIMS.md) — performance comparison table.
- [`CHANGELOG.md`](./CHANGELOG.md) — release notes.
- [`examples/`](./examples) — runnable cookbook _(coming in Phase D)_.

## Contributing

Until v1.0.0-rc1, the API is unstable and external PRs are not
actively solicited. Issues and feature requests are very welcome —
they help shape what v1.0 looks like.

## License

Apache License 2.0. See [`LICENSE`](./LICENSE).
