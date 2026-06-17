# Table-Driven Specs, Containers, and Parallel Specs (Ginkgo/Gomega)

Detailed reference for structuring table-driven specs using `DescribeTable`,
nested containers, and parallel execution in Ginkgo.
Sources: Google Go Style Guide, Uber Go Style Guide, adapted to Ginkgo/Gomega.

---

## Basic Structure

```go
var _ = Describe("Compare", func() {
    DescribeTable("comparing two strings",
        func(a, b string, want int) {
            Expect(Compare(a, b)).To(Equal(want))
        },
        Entry("both empty", "", "", 0),
        Entry("a longer", "a", "", 1),
        Entry("b longer", "", "a", -1),
        Entry("equal strings", "abc", "abc", 0),
    )
})
```

---

## Best Practices

**Always give each `Entry` a description** as its first argument — this is the
Ginkgo equivalent of naming a table row, and it is what shows up in spec
output and `--focus` filtering:

```go
DescribeTable("classifying numbers",
    func(input int, want string) {
        Expect(Classify(input)).To(Equal(want))
    },
    Entry("zero", 0, "zero"),
    Entry("positive", 5, "positive"),
    Entry("negative", -5, "negative"),
)
```

**Don't identify entries by index** — give each `Entry` a descriptive name and
include the inputs in custom failure messages when using `Expect(...).To(...,
"context: %v", input)`.

**Use `EntryDescription` for dynamically generated tables** when entries are
generated programmatically rather than listed by hand:

```go
DescribeTable("validating inputs",
    func(input Input) {
        Expect(Validate(input)).To(Succeed())
    },
    EntryDescription("validates %s"),
    generateInputEntries()...,
)
```

---

## Avoid Complexity in Table Specs

When specs need complex setup, conditional mocking, or multiple branches,
prefer separate `It` blocks or nested `Describe`/`Context` containers over
`DescribeTable`.

```go
// Bad: too many conditional fields make the table hard to understand
DescribeTable("complex thing",
    func(give string, want string, wantErr error,
        shouldCallX, shouldCallY bool,
        giveXResponse string, giveXErr error,
        giveYResponse string, giveYErr error) {
        if shouldCallX {
            xMock.EXPECT().Call().Return(giveXResponse, giveXErr)
        }
        if shouldCallY {
            yMock.EXPECT().Call().Return(giveYResponse, giveYErr)
        }
        // ...
    },
    Entry("calls x", "inputX", "wantX", nil, true, false, "XResponse", nil, "", nil),
)

// Good: separate, focused specs are clearer
var _ = Describe("DoComplexThing", func() {
    Context("when X is needed", func() {
        It("calls X and succeeds", func() {
            xMock.EXPECT().Call().Return("XResponse", nil)
            got, err := DoComplexThing("inputX", xMock, yMock)
            Expect(err).NotTo(HaveOccurred())
            Expect(got).To(Equal("wantX"))
        })
    })

    Context("when Y fails", func() {
        It("returns the Y error", func() {
            yMock.EXPECT().Call().Return("YResponse", nil)
            _, err := DoComplexThing("inputY", xMock, yMock)
            Expect(err).To(HaveOccurred())
        })
    })
})
```

**`DescribeTable` works best when:**

- All entries run identical logic (no conditional assertions)
- Setup is the same for all entries
- No conditional mocking based on entry fields
- All entry parameters are used by every entry

A single `wantErr bool` parameter for success/failure is acceptable if the
spec body is short and straightforward.

---

## Containers: Describe, Context, When

Use `Describe`, `Context`, and `When` to organise specs hierarchically. They
are interchangeable aliases — choose the one that reads best:

- `Describe` — the subject under test (a function, type, or package)
- `Context` — a particular condition or branch ("when the input is empty")
- `When` — a temporal or conditional framing, often clearer than `Context`

### Container Naming

- Use clear, concise descriptions: `Context("when input is empty", ...)`,
  `When("the locale is hu", func() { ... })`
- Avoid wordy descriptions or slashes (slashes interfere with `--focus`
  regex filtering)
- Specs must be independent — no shared state or execution order dependencies
  between sibling `It` blocks

### Nested Describe/Context with DescribeTable

```go
var _ = Describe("Translate", func() {
    DescribeTable("translating between languages",
        func(srcLang, dstLang, input, want string) {
            Expect(Translate(srcLang, dstLang, input)).To(Equal(want))
        },
        Entry("hu to en basic", "hu", "en", "köszönöm", "thank you"),
    )
})
```

---

## Parallel Specs

Ginkgo parallelises at the process level via `ginkgo -p`, splitting top-level
containers across worker processes — there is no per-spec `t.Parallel()`
equivalent to opt into manually. Specs within a single process still run
serially in randomised order.

```go
var _ = Describe("Process", func() {
    DescribeTable("processing values",
        func(give, want string) {
            // Each entry runs in its own randomised position;
            // no loop-variable capture concerns since Ginkgo
            // passes entry parameters as function arguments,
            // not via a shared loop variable.
            Expect(Process(give)).To(Equal(want))
        },
        Entry("first", "a", "A"),
        Entry("second", "b", "B"),
    )
})
```

Key points:

- Run with `ginkgo -p` (or `ginkgo -p -procs=N`) to parallelise across
  processes; specs that must not run concurrently with others (e.g. ones
  mutating shared external state) should use `Serial` decorators
- Because `DescribeTable` entry parameters are passed as ordinary function
  arguments rather than captured from a loop variable, the classic Go
  `tt := tt` loop-capture pitfall does not apply
- Use `Ordered` containers when specs within a single `Describe` must run in
  declared order and share state via `BeforeEach`/`AfterEach` in that
  container
  