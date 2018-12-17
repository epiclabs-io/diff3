# diff3

Package diff3 implements a three-way merge algorithm
Original version in Javascript by Bryan Housel @bhousel: https://github.com/bhousel/node-diff3,
which in turn is based on project Synchrotron, created by Tony Garnock-Jones. For more detail please visit:
http://homepages.kcbbs.gen.nz/tonyg/projects/synchrotron.html
https://github.com/tonyg/synchrotron

Ported to go by Javier Peletier @jpeletier

## How to use:

Import the package and call `Merge`:

`Merge(a, o, b io.Reader, detailed bool, labelA string, labelB string) (*MergeResult, error)`

Where:
* `a` and `b` are the current and incoming versions of content
* `o` is the common ancestor
* `detailed` specifies whether you want the merge conflicts to be detailed.
* `labelA` and `labelB` indicate a label to mark conflicts with, usually the branch names.

Returns a `MergeResult` that includes a flag indicating whether there were conflicts or not and a `io.Reader` stream with the result.