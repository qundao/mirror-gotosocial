# log

An opinionated logging package expecting to be configured **once** on server initialization, and with a `context.Context` almost always required.

Formatting with the `logfmt` function is handled by `codeberg.org/gruf/go-kv/v2/format`, which is significantly faster and more useful in map / struct formatting than the standard library pkg.
