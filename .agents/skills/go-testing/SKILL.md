---
name: go-testing
description: Use when writing, reviewing, or improving Go test code using Ginkgo/Gomega — including DescribeTable specs, nested Describe/Context containers, spec helpers, test doubles, and assertions with Gomega matchers and cmp.Diff. Also use when a user asks to write a test for a Go function, even if they don't mention specific patterns like DescribeTable or containers. Only Ginkgo/Gomega specs are used in this project — no stdlib testing-package tests or testify. Does not cover benchmark performance testing (see go-performance).
license: Apache-2.0
compatibility: Uses github.com/onsi/ginkgo/v2, github.com/onsi/gomega, and github.com/google/go-cmp for cmp.Diff comparisons
metadata:
  sources: "Google Style Guide, Uber Style Guide, adapted to Ginkgo/Gomega"
allowed-tools: Bash(bash:*)
---

# Go Testing (Ginkgo/Gomega)

> **Project rule**: All Go tests in this project use Ginkgo/Gomega. Never
> generate `testing`-package-only tests (`func TestXxx(t *testing.T)` with
> `t.Error`/`t.Fatal` assertions) or testify (`assert`/`require`). The only
> permitted appearance of `*testing.T` is in the single per-package suite
> bootstrap file that calls `RunSpecs`.

## Quick Reference

| Pattern | Use When |
|---------|----------|
| `Expect(...).To(...)` | Default — Gomega matcher assertion |
| Multiple `Expect` in one `It` | Later checks depend on earlier ones succeeding |
| Separate `It` blocks | Independent properties that should each report on their own |
| `cmp.Diff` / `BeComparableTo` | Comparing structs, slices, maps, protos |
| `DescribeTable` | Many entries share identical logic |
| `Describe` / `Context` / `When` | Grouping specs, signalling subject/condition |
| `DeferCleanup()` | Teardown in helpers or `BeforeEach`, instead of defer |
| `BeforeEach` / `AfterEach` | Setup/teardown scoped to specs in a container |
| `BeforeSuite` / `AfterSuite` | Rare — only genuinely suite-wide resources |

---

## Useful Spec Failures

> **Normative**: Spec failures must be diagnosable without reading the spec
> source.

Every failure must make clear: what was tested, the actual (got), and the
expected (want). Gomega matchers do this automatically for most cases
(`Expect(got).To(Equal(want))` produces a clear diff-style failure); add a
custom message only when the matcher alone wouldn't make the inputs obvious:

```go
// Good — Gomega's own message is enough:
Expect(got).To(Equal(want))

// Good — custom message adds the input that produced got:
Expect(got).To(Equal(want), "Add(2, 3) produced an unexpected result")

// Bad: vague custom message that adds no information
Expect(got).To(Equal(want), "mismatch")
```

---

## No Hand-Rolled Assertion Helpers, No Testify

> **Normative**: Do not use testify (`assert`/`require`) or build custom
> assertion libraries. Use Gomega matchers directly, and `cmp.Diff`
> (surfaced via `Expect(...).To(BeEmpty())` or `BeComparableTo`) for complex
> comparisons.

```go
Expect(cmp.Diff(want, got)).To(BeEmpty())
// or
Expect(got).To(BeComparableTo(want))
```

For protocol buffers, add `protocmp.Transform()` as a `cmp` option. Avoid
comparing JSON/serialized output — compare semantically instead.

> Read [references/TEST-HELPERS.md](references/TEST-HELPERS.md) when writing
> custom comparison helpers or domain-specific spec utilities.

---

## Expect Semantics: When to Split Into Separate `It` Blocks

> **Normative**: Prefer one `It` per independent property so an early failure
> doesn't hide a later one. Use multiple sequential `Expect` calls within a
> single `It` only when each check depends on the previous one succeeding.

**Keep checks in the same `It` when:**
- Setup fails and continuing is meaningless (e.g., decode after encode)
- The checks are really one logical assertion expressed in steps

**Split into separate `It` blocks when:**
- The properties are independent and you want each to report on its own

**Never call `Expect` from a goroutine** other than the spec's own goroutine
— pass a `Gomega` instance explicitly via `Eventually(func(g Gomega) {...})`
instead.

> Read [references/TEST-HELPERS.md](references/TEST-HELPERS.md) when writing
> helpers that need to choose between single and multi-`Expect` specs, or for
> detailed examples of both.

---

## DescribeTable Specs

> See `assets/table-test-template.go` when scaffolding a new `DescribeTable`
> spec and need the canonical `Describe`/`DescribeTable`/`Entry` layout.

> **Advisory**: Use `DescribeTable` when many entries share identical logic.

**Use `DescribeTable` when:** all entries run the same code path with no
conditional setup, mocking, or assertions. A single `wantErr bool` parameter
is acceptable.

**Don't use `DescribeTable` when:** entries need complex setup, conditional
mocking, or multiple branches — write separate `It` blocks or nested
`Describe`/`Context` containers instead.

**Key rules:**
- Give every `Entry` a clear description as its first argument — never rely
  on entry position/index for identification
- Include the inputs in custom failure messages when the matcher alone
  wouldn't make them obvious

> Read [references/TABLE-DRIVEN-TESTS.md](references/TABLE-DRIVEN-TESTS.md)
> when writing `DescribeTable` specs, nested containers, or reasoning about
> parallel (`ginkgo -p`) execution.

> **Validation**: After generating or modifying specs, run
> `ginkgo -focus="<SpecName>" -v` (or `go test ./... -run TestSuite -v` for
> the package's suite entrypoint) to verify the specs compile and pass. Fix
> any compilation errors before proceeding.

---

## Spec Helpers

> **Normative**: Spec helpers use `Expect`/`DeferCleanup` directly — no
> `t.Helper()` equivalent is needed since Gomega already reports the
> matcher's own call site.

```go
func setupTestDB() *sql.DB {
    db, err := sql.Open("sqlite3", ":memory:")
    Expect(err).NotTo(HaveOccurred(), "Could not open database")
    DeferCleanup(func() { db.Close() })
    return db
}
```

> Read [references/TEST-HELPERS.md](references/TEST-HELPERS.md) when writing
> spec helpers, cleanup functions, or custom comparison utilities.

---

## Spec Error Semantics

> **Advisory**: Test error semantics, not error message strings.

```go
// Bad: Brittle string comparison
Expect(err.Error()).To(Equal("invalid input"))

// Good: Semantic check
Expect(errors.Is(err, ErrInvalidInput)).To(BeTrue())
// or, more idiomatically in Gomega:
Expect(err).To(MatchError(ErrInvalidInput))
```

For simple presence checks when specific semantics don't matter:

```go
if tt.wantErr {
    Expect(err).To(HaveOccurred())
} else {
    Expect(err).NotTo(HaveOccurred())
}
```

---

## Test Organization

> Read [references/TEST-ORGANIZATION.md](references/TEST-ORGANIZATION.md) when
> working with test doubles, choosing spec package placement, suite bootstrap
> files, or scoping `BeforeEach`/`BeforeSuite` setup.

> Read [references/VALIDATION-APIS.md](references/VALIDATION-APIS.md) when
> designing reusable, framework-agnostic validation functions for use from
> Ginkgo specs.

---

## Integration Testing

> Read [references/INTEGRATION.md](references/INTEGRATION.md) when writing
> `BeforeSuite`/`AfterSuite` setup, acceptance tests, or specs that need real
> HTTP/RPC transports.

---

## Available Scripts

- **`scripts/gen-table-test.sh`** — Generates a Ginkgo/Gomega `DescribeTable`
  spec scaffold

```bash
bash scripts/gen-table-test.sh ParseConfig config > config/parse_config_test.go
bash scripts/gen-table-test.sh --ordered ParseConfig config      # wraps table in an Ordered container
bash scripts/gen-table-test.sh --output config/parse_config_test.go ParseConfig config
```

---

## Related Skills

- **Error testing**: See [go-error-handling](../go-error-handling/SKILL.md) when testing error semantics with `errors.Is`/`errors.As`, sentinel errors, or Gomega's `MatchError`
- **Interface mocking**: See [go-interfaces](../go-interfaces/SKILL.md) when creating test doubles by implementing interfaces at the consumer side
- **Naming test functions**: See [go-naming](../go-naming/SKILL.md) when naming suite bootstrap functions, `Describe`/`Context` descriptions, or spec helper utilities
- **Linter integration**: See [go-linting](../go-linting/SKILL.md) when running linters alongside `ginkgo`/`go test` in CI or pre-commit hooks
