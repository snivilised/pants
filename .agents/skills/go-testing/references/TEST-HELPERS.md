# Test Helpers, Assertions, and Comparisons (Ginkgo/Gomega)

Detailed reference for writing spec helpers, using Gomega matchers instead of
assertion libraries, and choosing how strictly a failure should halt a spec.
Sources: Google Go Style Guide, Uber Go Style Guide, adapted to Ginkgo/Gomega.

---

## Spec Helper Pattern

Helpers used inside specs don't need a `t.Helper()` equivalent — Gomega's
default fail handler reports the matcher's own call site clearly. Use
`DeferCleanup` for teardown, and return values directly (panicking via
`Expect` only on genuine setup failure).

```go
func mustLoadTestData(filename string) []byte {
    data, err := os.ReadFile(filename)
    Expect(err).NotTo(HaveOccurred(), "Setup failed: could not read %s", filename)
    return data
}

func setupTestDB() *sql.DB {
    db, err := sql.Open("sqlite3", ":memory:")
    Expect(err).NotTo(HaveOccurred(), "Could not open database")
    DeferCleanup(func() { db.Close() })
    return db
}
```

**Key rules:**

- Helpers can be called directly from `BeforeEach`, `It`, or other helpers —
  `Expect` works from any goroutine-free call path inside the current spec
- Use `Expect(...).To(Succeed())` or `Expect(err).NotTo(HaveOccurred())` for
  setup failures rather than returning errors for the caller to check
- Use `DeferCleanup()` for teardown instead of `defer` — it runs even if the
  spec aborts partway through, and works whether called from `BeforeEach`,
  `BeforeSuite`, or directly inside an `It`

---

## Avoiding Hand-Rolled Assertion Libraries

> **Normative**: Do not create or use ad hoc assertion helpers. Use Gomega's
> built-in matchers, and `gomega.DiffMatcher`/`cmp.Diff` for structural
> comparisons.

```go
// Bad:
assertIsNotNil(obj)
assertStringEq(obj.Type, "blogPost")
assertIntEq(obj.Comments, 2)

// Good: Gomega matchers express intent directly
Expect(obj).NotTo(BeNil())
Expect(obj.Type).To(Equal("blogPost"))
Expect(obj.Comments).To(Equal(2))

// Good: cmp.Diff for whole-struct comparisons, surfaced via Gomega
want := BlogPost{
    Type:     "blogPost",
    Comments: 2,
    Body:     "Hello, world!",
}
Expect(cmp.Diff(want, got)).To(BeEmpty())
```

### Domain-Specific Comparisons

For domain-specific comparisons, wrap the comparison in a plain Go function
that returns a value, and assert on that value with Gomega — don't build a
custom assertion API:

```go
func postLength(p BlogPost) int { return len(p.Body) }

var _ = Describe("BlogPost", func() {
    It("has the expected body length", func() {
        post := BlogPost{Body: "Hello"}
        Expect(postLength(post)).To(Equal(5))
    })
})
```

---

## Comparisons and Diffs

Prefer `cmp.Diff` wrapped in `Expect(...).To(BeEmpty())` for complex types, or
Gomega's own `Equal`/`BeComparableTo` matchers where a single value-equality
check is enough.

```go
// Struct comparison via cmp.Diff
want := &Doc{Type: "blogPost", Authors: []string{"isaac", "albert"}}
Expect(cmp.Diff(want, got)).To(BeEmpty())

// Struct comparison via Gomega's BeComparableTo (uses cmp under the hood)
Expect(got).To(BeComparableTo(want))

// Protocol buffers
Expect(cmp.Diff(want, got, protocmp.Transform())).To(BeEmpty())
```

**Avoid unstable comparisons** — don't compare JSON/serialized output that may
change. Compare semantically instead, using `Equal`, `BeComparableTo`, or
field-by-field `Expect` calls.

---

## Strict vs Soft Failures: Expect vs Eventually/Soft Assertions

Gomega's `Expect` halts the current spec on failure by default (equivalent in
spirit to `t.Fatal`), since Ginkgo specs are intended to fail fast on the
first broken assumption rather than accumulate multiple unrelated failures in
one `It`. To report several independent checks without aborting early, either
split them into separate `It` blocks, or use a `SpecContext`/`Gomega` instance
combined with explicit `g.Expect` calls inside `Eventually`/`Consistently`
where soft accumulation is genuinely needed.

```go
// Good: each independent property gets its own It, so a failure
// in one doesn't hide a failure in another
var _ = Describe("Stats", func() {
    It("computes the correct mean", func() {
        Expect(gotMean).To(BeNumerically("~", wantMean, 0.001))
    })

    It("computes the correct variance", func() {
        Expect(gotVariance).To(BeNumerically("~", wantVariance, 0.001))
    })
})
```

Use a single `It` with multiple sequential `Expect` calls when each check
depends on the previous one succeeding (the Ginkgo equivalent of `t.Fatal`
chaining):

```go
It("round-trips through encode and decode", func() {
    gotEncoded := Encode(input)
    Expect(gotEncoded).To(Equal(wantEncoded))

    gotDecoded, err := Decode(gotEncoded)
    Expect(err).NotTo(HaveOccurred())
    Expect(gotDecoded).To(Equal(input))
})
```

### Don't Call Expect from an Untracked Goroutine

> **Normative**: Never call Gomega's `Expect` (or any matcher) from a
> goroutine other than the spec's own goroutine, unless you pass an explicit
> `Gomega` instance captured via `NewWithT`/`Eventually(func(g Gomega) {...})`.
> The global `Expect` registers failures against the currently running spec
> via package-level state, which is not goroutine-safe across concurrent
> specs or background goroutines.

```go
// Bad: global Expect from a background goroutine
go func() {
    Expect(worker.Result()).To(Equal(want)) // may race or panic
}()

// Good: pass a Gomega instance explicitly into Eventually
Eventually(func(g Gomega) {
    g.Expect(worker.Result()).To(Equal(want))
}).Should(Succeed())
```
