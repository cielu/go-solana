# Claims

> **Every performance claim this project makes in public — in the
> README, in a blog post, in a tweet — must be recorded here and
> backed by a checked-in benchmark. No speculative numbers.**

## Philosophy

1. **No claim without a benchmark.** Every row in the tables below
   links to the `Benchmark…` function that produced it.
2. **Numbers come from a reference machine.** See §"Reference
   hardware" below. Different hardware gets a different table.
3. **We list our weaknesses honestly.** If `gagliardetto/solana-go`
   beats us on a benchmark, it goes in the table with the real ratio.
4. **Benchmarks live in version control.** `benchmarks/testdata/`
   holds the raw RPC fixtures so runs are reproducible.

## Reference hardware

_TBD. Record the exact CPU, RAM, OS, Go version, and date of the
run in this section whenever the numbers table is refreshed._

```
CPU:   <model, cores>
RAM:   <amount, type>
OS:    <name + version>
Go:    <go version go1.x.y>
Date:  YYYY-MM-DD
SHA:   <commit SHA of this repo at the time of the run>
```

## Comparison table

> Empty table. Populated by the Phase A milestone. Units: ns/op,
> allocations/op, B/op. Lower is better.

| # | Benchmark | solana-go (ours) | gagliardetto/solana-go | Ratio | Source |
|---|---|---|---|---|---|
| 1 | `BenchmarkPublicKey_Base58` | — | — | — | — |
| 2 | `BenchmarkPublicKey_UnmarshalJSON` | — | — | — | — |
| 3 | `BenchmarkMessage_MarshalLegacy` | — | — | — | — |
| 4 | `BenchmarkMessage_UnmarshalLegacy` | — | — | — | — |
| 5 | `BenchmarkMessage_MarshalV0` | — | — | — | — |
| 6 | `BenchmarkMessage_UnmarshalV0` | — | — | — | — |
| 7 | `BenchmarkTransaction_Sign_1signer` | — | — | — | — |
| 8 | `BenchmarkTransaction_Sign_4signers` | — | — | — | — |
| 9 | `BenchmarkGetAccountInfo_Decode` | — | — | — | — |
| 10 | `BenchmarkGetBlock_Decode` | — | — | — | — |

## Non-performance claims

These are qualitative claims we are willing to make in public once
each line item is honestly true.

- [ ] **Zero reflection in hot paths.** Verified by `grep -r reflect.
      rpc/ internal/binary` returning nothing.
- [ ] **Zero library-code `panic`s.** Verified by `grep -rn 'panic(' .
      | grep -v _test.go | grep -v cmd/` returning nothing.
- [ ] **Zero `encoding/json` on hot paths.** Verified by
      `grep -rn 'encoding/json' rpc/ internal/` returning nothing.
- [ ] **Apache-2.0 throughout, no LGPL vendored code.** Verified by
      `grep -rn LGPL` returning nothing.
- [ ] **`go vet ./...` clean.** Verified in CI on every push.
- [ ] **`staticcheck ./...` clean.** Verified in CI.
- [ ] **`go test -race ./...` clean.** Verified in CI.
- [ ] **>80% coverage on core types**
      (`solana.go`, `message.go`, `transaction.go`).
- [ ] **Every exported identifier has a doc comment.**

## Known weaknesses

_This section exists on purpose. Being honest here is a feature._

- _No numbers yet; populate after Phase A._
