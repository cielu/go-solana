# Roadmap

> **Status**: Pre-1.0. This document is the public phase plan for
> `github.com/cielu/solana-go`. Every phase has acceptance criteria,
> and the project does not advance to the next phase until they are
> met.

## Vision

The fastest, cleanest, most honestly-documented Solana client in Go.

The direct reference point is `gagliardetto/solana-go`. Our goal is
**not** to match them feature-for-feature — it is to be **measurably
better** on performance, API design, and documentation, and to build
an ecosystem moat via a first-class Anchor IDL code generator.

Every public performance claim we make must be backed by a
benchmark in this repository. No marketing numbers.

## Release tracks

- `v0.x.x`: unstable. Breaking changes allowed on every minor bump.
- `v1.0.0-rcN`: API frozen modulo fatal bugs. Benchmarks must match
  the claims in CLAIMS.md on a reference machine.
- `v1.0.0`: the shippable, stable, advertised release.

---

## Phase 0 — Bootstrap

**Status: 🟡 in progress**

- [x] Module path, license, Go version decided.
- [x] Repository scaffolded with `CLAUDE.md`, `ROADMAP.md`,
      `CLAIMS.md`, `README.md`, `LICENSE`, `.github/workflows/ci.yml`.
- [ ] First commit lands on `main`. CI green.
- [ ] JSON library decision (benchmark-driven).

**Exit criterion**: `go build ./...` and `go vet ./...` succeed on a
fresh clone, CI is green, the JSON dependency is decided and
justified in `CLAIMS.md`.

---

## Phase A — Performance-first core (weeks 1–3)

**Goal**: make the foundational code measurably faster than the
reference implementation, with benchmarks to prove it.

### Deliverables

1. `solana.go` — core types (`PublicKey`, `Hash`, `Signature`,
   `Account`), each with full JSON/text/sql round-trip tests,
   including the edge cases the older repo got wrong
   (silent truncation on oversized input, empty-string handling).
2. `internal/binary` — zero-copy binary `Encoder`/`Decoder` with a
   `sync.Pool` buffer. Must beat `encoding/binary` on a representative
   transaction marshal benchmark.
3. `message.go`, `transaction.go`, `instruction.go` — the
   message/transaction model with:
   - Legacy + v0 message encoding (explicit 0x80 high-bit convention,
     **not** ad-hoc arithmetic).
   - `NoSignedMarshalBinary` that actually emits zero-filled
     signature placeholders.
   - `Sign` with O(1) signer lookup and clean reset semantics.
4. `rpc/` — our own JSON-RPC transport. No vendored go-ethereum code.
   Pluggable JSON codec, pluggable HTTP client, pluggable
   `RetryPolicy`, `context.Context` plumbed through every call.
5. `benchmarks/testdata/` — real mainnet fixtures for `getAccountInfo`,
   `getBlock` (versioned, ~100 tx), `getTransaction`,
   `getSignaturesForAddress`.
6. A benchmark suite that covers:
   - `BenchmarkPublicKey_Base58`
   - `BenchmarkPublicKey_UnmarshalJSON`
   - `BenchmarkMessage_Marshal{Legacy,V0}`
   - `BenchmarkMessage_Unmarshal{Legacy,V0}`
   - `BenchmarkTransaction_Sign/{1,4,8}signers`
   - `BenchmarkGetAccountInfo_Decode`
   - `BenchmarkGetBlock_Decode`

### Acceptance criteria

- Every listed benchmark is in the repo and runs in CI.
- Each benchmark has its `gagliardetto/solana-go` counterpart in the
  same file (behind a build tag that pulls in the competitor as a
  dev-only dependency).
- `CLAIMS.md` has a populated comparison table with **real numbers**
  from a reference machine, and we are faster (or equal) on at least
  **8 of 10** benchmarks. If we are slower anywhere, that is explicitly
  listed as a known weakness in `CLAIMS.md`.

---

## Phase B — Ergonomics and reliability (weeks 3–5)

**Goal**: the library should be a pleasure to use and hard to misuse.

### Deliverables

1. `client.go` — typed wrappers for every stable Solana RPC method.
   - Functional options for every optional parameter. No
     `...interface{}`.
   - Every method returns wrapped errors with method name and context.
   - Limit and range validation surfaces as descriptive errors, never
     silent clipping.
2. `wsclient.go` — WebSocket subscription wrappers for
   `accountSubscribe`, `blockSubscribe`, `logsSubscribe`,
   `programSubscribe`, `signatureSubscribe`, `slotSubscribe`,
   `rootSubscribe`, `slotsUpdatesSubscribe`. Each has a typed channel.
3. `helpers/send.go` — `SendAndConfirmTransaction`:
   - Pluggable `CommitmentLevel`.
   - Automatic blockhash refresh on `BlockhashNotFound`.
   - Bounded retry with jitter.
   - Optional pre-flight simulation.
4. `helpers/priorityfee.go` — `EstimatePriorityFee`: reads
   `getRecentPrioritizationFees`, returns p50/p75/p95 with optional
   per-account scoping.
5. `helpers/simulate.go` — `SimulateTransaction` with decoded
   instruction errors and a readable diff of account changes.
6. `errors.go` — sentinel errors (`ErrBlockhashExpired`,
   `ErrInsufficientFunds`, `ErrRateLimited`, …) and an `ErrRPC`
   type that carries the RPC method name and raw response body for
   diagnostics.

### Acceptance criteria

- Unit test coverage ≥ 80% on `solana.go`, `message.go`,
  `transaction.go`, `client.go`, `wsclient.go`, `helpers/`.
- Integration tests guarded by `//go:build integration` that hit a
  local `solana-test-validator`. CI runs them against a validator
  container.
- Public API frozen for v1.0: every exported identifier has a doc
  comment, every breaking change after this point requires a version
  bump and a `CHANGELOG` entry.

---

## Phase C — IDL code generator (weeks 5–8)

**Goal**: the ecosystem moat. Any Anchor program can be bound in
minutes, with generated code that is cleaner than hand-written code.

### Deliverables

1. `cmd/anchoridl-gen/` — CLI:
   ```
   anchoridl-gen -idl path/to/program.json -out programs/foo [-package foo]
   ```
2. Generated output per program:
   - `accounts.go`: every account struct with discriminator constant,
     `Unmarshal` that does one zero-copy pass, no reflection.
   - `instructions.go`: one builder per instruction using our
     functional-options convention. Each builder validates required
     signers/accounts at build time.
   - `errors.go`: Go `error` sentinel per Anchor error code.
   - `events.go`: struct per event with `ParseEvent` helper.
   - `types.go`: defined types (enums, structs, bitflags).
   - `pda.go`: `FindXxxAddress` helpers for every PDA declared in the
     IDL, with correct seed decoding (including arg-typed seeds).
3. Golden-file tests: the generator's output for a known IDL must be
   byte-identical to a checked-in expected file. Drift is a bug.
4. At least **four** real program IDLs generated as showcase:
   - Metaplex Token Metadata
   - Orca Whirlpool
   - Pyth Oracle
   - Jupiter Aggregator v6

   These generated bindings live in a **separate** repository
   (`cielu/solana-programs`), not in this one.

### Acceptance criteria

- Generator output compiles, passes `go vet`, has no unused
  identifiers, and runs the generator's own unit tests.
- The four showcase programs build and round-trip at least one
  representative instruction through `SimulateTransaction` against a
  local test validator.
- Our generated code is **fewer lines** than equivalent `anchor-go`
  output for the same IDL, and the generated types are plain Go
  structs (no reflection helpers).

---

## Phase D — Showcase, docs, launch (weeks 8–10)

**Goal**: make the library findable, learnable, and credible.

### Deliverables

1. `examples/` — a runnable cookbook:
   - `01-hello` — connect and print an account balance.
   - `02-send-transfer` — build, sign, send a SOL transfer.
   - `03-versioned-tx` — address lookup tables end-to-end.
   - `04-send-and-confirm` — the helper, with simulated failure.
   - `05-subscribe-logs` — WebSocket logs with clean shutdown.
   - `06-idl-generate` — generate a small Anchor program binding
     from scratch.
2. `docs/` — a minimal mdBook-style docs site generated from the
   examples. No rocket science.
3. `CLAIMS.md` locked with final numbers for `v1.0.0-rc1`, dated,
   with machine specs.
4. `README.md` with:
   - One-line pitch.
   - Benchmark highlight table (3 rows, linked to `CLAIMS.md`).
   - `go get` line.
   - Five-line quickstart.
   - Link to cookbook and ROADMAP.
5. Launch post for the Solana Go community (Twitter/X thread,
   Medium/mirror article) — **only after** CI green, benchmarks
   published, four IDL bindings generated.

### Acceptance criteria

- `go get github.com/cielu/solana-go` from a fresh module works.
- Every example in `examples/` is compiled and run in CI against the
  local test validator.
- README passes `markdown-lint` and every link resolves.
- `v1.0.0-rc1` tagged, GitHub Release created, announcement posted.

---

## Phase E — Continuous hardening (post-1.0)

- Follow every new Solana RPC method within 2 weeks of availability.
- Treat every bug report that costs a user real money as a
  drop-everything priority.
- Publish `CLAIMS.md` re-runs on every minor release.
- Maintain the four showcase IDL bindings in `cielu/solana-programs`
  alongside upstream program upgrades.
- Community PRs get a first response within 48 hours, merged or
  rejected with reason within 7 days.

---

## Things we will not build

- **Wallet adapters / frontend helpers**. Out of scope. A server-side
  Go SDK has no business handling browser wallets.
- **Indexer / historical data store**. Different shape of problem,
  different library.
- **Staking pool management tooling**. Leave to specialized tools.
- **Our own JSON-RPC server** beyond what is needed for
  `solana-test-validator`-driven integration tests.
