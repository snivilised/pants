<!-- MD029/ol-prefix - Ordered list item prefix -->
<!-- MarkDownLint-disable MD029 -->

# Go Testing: Integration and Advanced Patterns (Ginkgo/Gomega)

Detailed reference for suite-level setup, acceptance testing, and real transport
testing using Ginkgo and Gomega.
Sources: Google Go Style Guide (best-practices), adapted to Ginkgo/Gomega.

---

## Suite-Level Setup

> **Source**: Google Go Style Guide (best-practices), adapted

Use `BeforeSuite`/`AfterSuite` when **all specs in the suite** require common
setup that needs teardown (e.g., a shared database). This should **not be your
first choice** — prefer `BeforeEach`/`DeferCleanup` scoped to the specs that
need them.

```go
var db *sql.DB

var _ = BeforeSuite(func() {
    ctx := context.Background()
    d, err := setupDatabase(ctx)
    Expect(err).NotTo(HaveOccurred())
    db = d
})

var _ = AfterSuite(func() {
    Expect(db.Close()).To(Succeed())
})

var _ = Describe("Insert", func() {
    It("inserts a row", func() {
        // uses db
    })
})

var _ = Describe("Select", func() {
    It("selects a row", func() {
        // uses db
    })
})
```

Key points:

- Run the suite via a `TestXxx(t *testing.T)` entrypoint that calls
  `RunSpecs(t, "Xxx Suite")` — this is the only place `*testing.T` appears
- `BeforeSuite`/`AfterSuite` run once per process, not once per spec
- Ensure individual specs remain hermetic — reset any global state they modify
  in `BeforeEach`/`AfterEach`, not by relying on suite ordering

---

## Acceptance Testing

> **Source**: Google Go Style Guide (best-practices), adapted

Acceptance testing validates that an implementation upholds a contract, treating
it as a black box. This pattern is useful when users implement your interfaces
and you want to provide a reusable validation suite.

### Structure

1. Create a test helper package (e.g., `chesstest` for package `chess`)
2. Export a validation function that accepts the implementation under test and
   returns an error (never a `*testing.T`, so it stays framework-agnostic):

```go
// Package chesstest provides acceptance tests for chess.Player implementations.
package chesstest

// ExercisePlayer tests a Player implementation in a single turn.
// Returns nil if the player makes a correct move, or an error describing
// the violation.
func ExercisePlayer(b *chess.Board, p chess.Player) error {
    move := p.Move()
    if putsOwnKingIntoCheck(b, move) {
        return &IllegalMoveError{Move: move, Reason: "puts own king in check"}
    }
    return nil
}
```

3. End users write a Ginkgo spec against the validation function:

```go
var _ = Describe("DeepBlue acceptance", func() {
    It("never puts its own king in check", func() {
        player := deepblue.New()
        err := chesstest.ExerciseGame(chesstest.SimpleGame, player)
        Expect(err).NotTo(HaveOccurred())
    })
})
```

Reserve hard suite aborts (`Fail` via `Expect` in a `BeforeEach`) for setup
failures only — validation errors from the acceptance function should be
asserted with `Expect(err).NotTo(HaveOccurred())`, not panic-style aborts.

---

## Use Real Transports

> **Source**: Google Go Style Guide (best-practices), adapted

When testing component integrations over HTTP or RPC, prefer real transport
round-trips over hand-implemented client mocks:

```go
var _ = Describe("API integration", func() {
    var srv *httptest.Server

    BeforeEach(func() {
        srv = httptest.NewServer(newFakeHandler())
        DeferCleanup(srv.Close)
    })

    It("fetches a user", func() {
        client := api.NewClient(srv.URL)
        result, err := client.GetUser(context.Background(), "user-123")
        Expect(err).NotTo(HaveOccurred())
        Expect(result.Name).To(Equal("Test User"))
    })
})
```

Using the production client with a test server ensures your spec exercises as
much real code as possible, avoiding the complexity of imitating client
behaviour.

---

## Common Mistakes

### Relying on spec execution order for setup/teardown

Ginkgo randomises spec order by default. Never rely on one spec's side effects
being visible to another. Use `BeforeEach`/`AfterEach` (or `BeforeSuite` only
for truly suite-wide, read-only resources) instead of mutable package-level
state shared across specs:

```go
// Bad: assumes specs run in declared order and share mutated state
var shared *Resource

var _ = Describe("First", func() {
    It("sets up shared", func() {
        shared = NewResource()
    })
})

var _ = Describe("Second", func() {
    It("uses shared", func() {
        Expect(shared).NotTo(BeNil()) // flaky: order is not guaranteed
    })
})

// Good: each spec is self-contained
var _ = Describe("Second", func() {
    var shared *Resource

    BeforeEach(func() {
        shared = NewResource()
    })

    It("uses shared", func() {
        Expect(shared).NotTo(BeNil())
    })
})
```

### Forgetting DeferCleanup

`DeferCleanup` is the Ginkgo equivalent of `t.Cleanup` — it runs after each
spec (or after the suite, when called in `BeforeSuite`) even if the spec fails:

```go
// Bad: defer in BeforeEach does not run between specs as expected
BeforeEach(func() {
    conn := mustConnect()
    defer conn.Close() // closes immediately, not after the spec
})

// Good
BeforeEach(func() {
    conn := mustConnect()
    DeferCleanup(conn.Close)
})
```
