# diff3

Package diff3 implements a three-way merge algorithm
Original version in Javascript by Bryan Housel @bhousel: https://github.com/bhousel/node-diff3,
which in turn is based on project Synchrotron, created by Tony Garnock-Jones. For more detail please visit:
http://homepages.kcbbs.gen.nz/tonyg/projects/synchrotron.html
https://github.com/tonyg/synchrotron

Ported to go by Javier Peletier @jpeletier

## How to use

Import the package and call `Merge`:

`Merge(a, o, b io.Reader, detailed bool, labelA string, labelB string) (*MergeResult, error)`

Where:
* `a` and `b` are the current and incoming versions of content
* `o` is the common ancestor
* `detailed` specifies whether you want the merge conflicts to be detailed.
* `labelA` and `labelB` indicate a label to mark conflicts with, usually the branch names.

Returns a `MergeResult` that includes a flag indicating whether there were conflicts or not and a `io.Reader` stream with the result.

For callers that already have line or token slices, call the generic API directly:

```go
result := diff3.Diff3Merge(aLines, oLines, bLines, true)
```

The generic type parameter is inferred for common calls. If inference is not possible, specify it explicitly:

```go
result := diff3.Diff3Merge[string](aLines, oLines, bLines, true)
```

`Diff3Merge` accepts any `comparable` element type, so callers can merge token slices as well as string lines.

## Algorithm Selection

The default algorithm remains the original Hunt-McIlroy LCS implementation.

Use `MergeWithOptions` or `Diff3MergeWithOptions` to opt into Myers diff:

```go
result, err := diff3.MergeWithOptions(a, o, b, diff3.MergeOptions{
	Algorithm:             diff3.DiffAlgorithmMyers,
	ExcludeFalseConflicts: true,
	Detailed:              true,
	LabelA:                "current",
	LabelB:                "incoming",
})
```

For slices:

```go
result := diff3.Diff3MergeWithOptions(aLines, oLines, bLines, diff3.MergeOptions{
	Algorithm:             diff3.DiffAlgorithmMyers,
	ExcludeFalseConflicts: true,
})
```

## Breaking API Changes

`Diff3Merge`, `Diff3MergeResult`, and `Conflict` are now generic.

Before:

```go
var result []*diff3.Diff3MergeResult
var conflict *diff3.Conflict
```

After:

```go
var result []*diff3.Diff3MergeResult[string]
var conflict *diff3.Conflict[string]
```

The fields on `Diff3MergeResult` and `Conflict` are now exported and inspectable:

```go
for _, item := range result {
	if item.Ok != nil {
		// merged non-conflict block
	}

	if item.Conflict != nil {
		_ = item.Conflict.A
		_ = item.Conflict.O
		_ = item.Conflict.B
		_ = item.Conflict.AIndex
		_ = item.Conflict.OIndex
		_ = item.Conflict.BIndex
	}
}
```

Direct calls to `Diff3Merge` with `[]string` usually continue to compile through type inference. Code that names `Diff3MergeResult` or `Conflict` must add the type argument, usually `[string]`.

## Attribution

The generic API, exported merge result fields, Myers diff implementation, and parallel diff computation are based on work by Kyungmin Bae and Devsisters in [`devsisters/go-diff3`](https://github.com/devsisters/go-diff3).
