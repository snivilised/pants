# Extensible Validation APIs (Ginkgo/Gomega)

Detailed reference for designing reusable validation functions that callers can
use for acceptance testing from within Ginkgo specs.
Sources: Google Go Style Guide (best-practices), adapted to Ginkgo/Gomega.

---

## The `*test` Package Export Pattern

When you own an interface that others implement, export a validation function
in a companion `*test` package. This lets implementers verify correctness
without duplicating your spec logic, and without that package depending on
Ginkgo or Gomega at all.

```go
// Package storagetest provides acceptance tests for storage.Backend.
package storagetest

// Verify runs a validation suite against any storage.Backend.
// Returns an error describing the first violation, or nil on success.
func Verify(b storage.Backend) error {
    if err := verifyRoundTrip(b); err != nil {
        return fmt.Errorf("round-trip: %w", err)
    }
    if err := verifyNotFound(b); err != nil {
        return fmt.Errorf("not-found: %w", err)
    }
    return nil
}
```

Callers write a thin Ginkgo spec that plugs in their implementation:

```go
var _ = Describe("MyBackend acceptance", func() {
    It("satisfies the storage.Backend contract", func() {
        b := mybackend.New()
        Expect(storagetest.Verify(b)).To(Succeed())
    })
})
```

---

## Designing Extensible Validation Functions

**Return errors, not framework assertions.** This keeps validation functions
usable as plain Go functions, independent of Ginkgo/Gomega — callers decide
whether a violation becomes an `Expect(...).To(Succeed())` failure or
something else entirely.

```go
// Good: Returns error — caller controls how the failure surfaces
func ExercisePlayer(b *chess.Board, p chess.Player) error {
    move := p.Move()
    if putsOwnKingIntoCheck(b, move) {
        return &IllegalMoveError{Move: move, Reason: "puts own king in check"}
    }
    return nil
}

// Bad: Takes a Gomega instance and asserts internally — caller loses control
func ExercisePlayer(g Gomega, b *chess.Board, p chess.Player) {
    move := p.Move()
    g.Expect(putsOwnKingIntoCheck(b, move)).To(BeFalse())
}
```

**Use custom error types** for rich diagnostics when needed — Gomega's
`MatchError` and `%v`/`%+v` formatting in failure messages will surface these
naturally:

```go
type IllegalMoveError struct {
    Move   chess.Move
    Reason string
}

func (e *IllegalMoveError) Error() string {
    return fmt.Sprintf("illegal move %v: %s", e.Move, e.Reason)
}
```

Asserting on the specific error type from a spec:

```go
It("rejects a move into check", func() {
    err := chesstest.ExercisePlayer(board, recklessPlayer)
    var illegal *chesstest.IllegalMoveError
    Expect(errors.As(err, &illegal)).To(BeTrue())
    Expect(illegal.Reason).To(Equal("puts own king in check"))
})
```

---

## When to Use Validation APIs vs Simple Helpers

| Situation | Use |
| --------- | --- |
| Interface you own, others implement | Validation API in `*test` package |
| Shared setup across specs in one package | `BeforeEach` helper |
| Complex assertion reused in 2-3 specs | Helper returning `error` or `bool` |
| One-off setup or comparison | Inline spec code |

**Validation APIs** are worth the extra package when:

- Multiple external packages will implement your interface
- The contract has non-obvious invariants that are easy to get wrong
- You want a single source of truth for "correct behavior"

**Simple helpers** are better when:

- The helper is a straightforward setup or comparison function
- The reuse is incidental, not part of a published contract

---

## Naming Conventions

Name the function with a verb that signals scope: `Verify`, `Exercise`,
`RunConformance`. Accept the interface under test as a parameter — never
construct the implementation inside the validation package, and never import
Ginkgo or Gomega into the validation package itself.

| Package | Function | Purpose |
| ------- | -------- | ------- |
| `storagetest` | `Verify` | Validates a `storage.Backend` |
| `chesstest` | `ExercisePlayer` | Validates a `chess.Player` |
| `cachetest` | `RunConformance` | Full conformance suite for `cache.Cache` |
