# Test Organization Reference (Ginkgo/Gomega)

Sources: Google Go Style Guide (best-practices, decisions), adapted to
Ginkgo/Gomega.

---

## Test Double Types

| Double | Purpose | State? | Verifies calls? |
| ------ | ------- | ------ | --------------- |
| Stub | Returns canned data | No | No |
| Fake | Working but simplified implementation | Yes | No |
| Spy | Records calls for later inspection | Yes | Yes |

**Prefer fakes over mocks.** Fakes are more readable and don't require mock
frameworks. Reserve spies for verifying side effects (e.g., an analytics
event). These types are framework-agnostic — the same fakes and spies work
whether asserted on via Gomega matchers in Ginkgo specs.

```go
// Fake: Working in-memory implementation
type FakeUserStore struct {
    users map[string]*User
}

func (f *FakeUserStore) GetUser(id string) (*User, error) {
    u, ok := f.users[id]
    if !ok {
        return nil, ErrNotFound
    }
    return u, nil
}

// Spy: Records calls for later assertion
type SpyEmailSender struct{ Sent []string }

func (s *SpyEmailSender) Send(to, body string) error {
    s.Sent = append(s.Sent, to)
    return nil
}
```

Asserting on a spy in a Ginkgo spec:

```go
var _ = Describe("Notifier", func() {
    It("sends an email to the recipient", func() {
        spy := &SpyEmailSender{}
        notifier := NewNotifier(spy)

        Expect(notifier.Notify("user@example.com", "hello")).To(Succeed())
        Expect(spy.Sent).To(ConsistOf("user@example.com"))
    })
})
```

---

## Test Double Naming Conventions

> **Advisory**: Follow consistent naming for test doubles (stubs, fakes,
> spies).

**Package naming**: Create a `*test` package alongside production code (e.g.,
`creditcardtest` for package `creditcard`, `fakeauthservice` for a standalone
fake service).

```go
// Good: In package creditcardtest

// Single double — use simple name
type Stub struct{}
func (Stub) Charge(*creditcard.Card, money.Money) error { return nil }

// Multiple behaviors — name by behavior
type AlwaysCharges struct{}
type AlwaysDeclines struct{}

// Multiple types — include type name
type StubService struct{}
type StubStoredValue struct{}
```

**Local variables**: Prefix test double variables with the double type for
clarity at the call site:

```go
// Good: Double type is immediately visible
spyCC := &creditcardtest.Spy{}
stubDB := &dbtest.Stub{Balance: 100}

// Bad: Ambiguous — is this real or a double?
cc := &creditcardtest.Spy{}
db := &dbtest.Stub{Balance: 100}
```

---

## Standalone Test Helper Packages

Create a standalone test helper package when multiple packages need the same
double, the helper has enough logic to warrant its own specs, or you want to
provide an acceptance test suite for interface implementers.

| Pattern | When to use | Example |
| ------- | ----------- | ------- |
| `footest` | General test helpers for package `foo` | `creditcardtest`, `usertest` |
| `fakeX` | Standalone fake service package | `fakeauthservice`, `fakestorage` |

```go
package usertest

// NewFakeStore builds a FakeUserStore pre-populated with the given users.
// Callers wire teardown themselves via DeferCleanup if the store holds any
// external resources; this in-memory store needs none.
func NewFakeStore(users ...*user.User) *FakeUserStore {
    store := &FakeUserStore{users: make(map[string]*user.User)}
    for _, u := range users {
        store.users[u.ID] = u
    }
    return store
}
```

Export plain constructors (no `*testing.T` parameter needed) so they remain
usable from both Ginkgo specs and any other Go code; use `DeferCleanup` at the
call site inside a spec when teardown is required.

---

## Suite Files and Bootstrap

Every package with Ginkgo specs needs exactly one suite bootstrap file that
calls `RunSpecs` — this is the only place `*testing.T` appears in a
Ginkgo-based package.

```go
package parser_test

import (
    "testing"

    . "github.com/onsi/ginkgo/v2"
    . "github.com/onsi/gomega"
)

func TestParser(t *testing.T) {
    RegisterFailHandler(Fail)
    RunSpecs(t, "Parser Suite")
}
```

| File Suffix | Use Case |
| ----------- | -------- |
| `_suite_test.go` | The single `RunSpecs` bootstrap per package |
| `_test.go` | `Describe`/`Context`/`It` spec definitions |

**Package declaration choice still applies the same way as with `testing`:**

| Package Declaration | Use Case |
| ------------------- | -------- |
| `package foo` | White-box specs, can access unexported identifiers |
| `package foo_test` | Black-box specs, avoids circular dependencies |

**Use `package foo` (white-box)** when you need to exercise unexported
functions or internal state from within specs.

**Use `package foo_test` (black-box)** when testing only the public API,
breaking import cycles, or verifying external usability.

```go
package parser_test // Black-box: only tests exported API

import (
    "mymodule/parser"

    . "github.com/onsi/ginkgo/v2"
    . "github.com/onsi/gomega"
)

var _ = Describe("Parse", func() {
    It("parses simple input", func() {
        got, err := parser.Parse("input")
        Expect(err).NotTo(HaveOccurred())
        // ...
    })
})
```

If a black-box suite needs an unexported symbol, create `export_test.go` in
`package foo` (not `foo_test`) that exposes it. Use this sparingly.

---

## Setup Scoping

> **Advisory**: Keep setup scoped to the specs that need it.

Explicit setup in `BeforeEach` within the relevant `Describe`/`Context` is
clearer and avoids penalizing unrelated specs elsewhere in the suite:

```go
// Good: Explicit setup scoped to the Describe that needs it
var _ = Describe("ParseData", func() {
    var data []byte

    BeforeEach(func() {
        data = mustLoadDataset()
    })

    It("parses the dataset", func() {
        // uses data
    })
})

var _ = Describe("Unrelated", func() {
    // Doesn't pay for dataset loading
})
```

**Avoid `BeforeSuite` for setup that only some specs need** — it runs once for
the entire process regardless of which specs end up running (including with
`--focus` filters), even specs that never touch the resource it sets up.

**Nested setup**: Use a parent `Describe`/`Context` with a shared `BeforeEach`
when a group of specs shares setup:

```go
var _ = Describe("Database", func() {
    var db *sql.DB

    BeforeEach(func() {
        db = setupTestDB()
    })

    Describe("Insert", func() {
        It("inserts a row", func() {
            // uses db
        })
    })

    Describe("Select", func() {
        It("selects a row", func() {
            // uses db
        })
    })
})
```

This scopes the database lifecycle to the specs that need it, re-running setup
fresh for each spec rather than sharing mutated state across them. Use
`BeforeSuite`/`AfterSuite` only as a last resort, and only for genuinely
suite-wide, ideally read-only, resources (see
[INTEGRATION.md](INTEGRATION.md)).
