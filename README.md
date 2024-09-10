# bmap
Better Map for Golang

Provides the following features over other standard library map implementations:

* Concurrency safe: all operations are protected using a RWMutex
* Ordered: map key-value pairs retain insertion order by default, which may be overriden by an arbitrary sort function
* Strongly typed: unlike `sync.Map`, bmaps enforce types (and don't require typecasting every return value)

Additionally, the `asyncBmap` provides ordered asynchronous write operations, providing exactly the same access garuntees as `bmap`, but with faster write calls.

## Usage

```go

```
