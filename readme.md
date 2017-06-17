# Forklift
Allows sql file embedding into a Go binary.

# Overview
Forklift is a library and an utility inspired by [gotic](https://github.com/gchaincl/gotic).

The lib `forklift/queries` provides an API to access sql files.

The utility `forklift` generates Go code from sql files and embbed them into the binary for caching.

# Install
To install Forklift:

`go get -u github.com/solher/forklift` 

# Basic Usage
A `query.static.sql` file Located in the same directory that the main.go file (only `.static.sql`, `.tmpl.sql` and `.lazy.sql` files are saved so it can ignore migration/fixture files):

```sql
SELECT * FROM documents
```

The `main.go` file:

```go
package main

import "github.com/solher/forklift/queries"

func main() {
  println(queries.File("query.sql"))
}

// Prints: SELECT * FROM documents
```

By default (development), the file will actually be read from the disk. In production, you can embbed the sql files into the binary by running the command:

```bash
$ forklift > forklift.go
```

The queries with then be read from memory.
