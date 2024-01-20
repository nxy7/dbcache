# Take home assignment
This library provides generic implementation of `Cache[T]` which can be accessed concurrently and is guaranteed to make at most
one call to underlying data source for any given key.

To function correctly this cache requires `DataRetriever[T]` that implements `Get(key string) (*T, error)`.

## How to run tests
Tests can be run either with `go test . -race` or if nix is installed `nix develop . -c bash -c "go test . -race"`. Second approach would
automatically install and use pinned Go version.

### Things to note
- tests can be a little flaky for huge amounts of concurrent requests (as far as I'm aware that's because of machine running out of RAM)
- this library caches keys indefinitely right now, but it's quite easy to add cache invalidation scheme

