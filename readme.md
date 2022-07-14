# Forklift
Allows file embedding into a Go binary.

# Overview
Forklift is a library and an utility inspired by [gotic](https://github.com/gchaincl/gotic).

The lib `forklift/files` provides an API to access files.

The utility `forklift` generates Go code from files and embbed them into the binary for caching.

# Install
To install Forklift:

`go get -u github.com/solher/forklift` 

# Basic Usage
A `file.sql` file located (for example) in the same directory as the main.go file:

```sql
SELECT * FROM documents
```

The `main.go` file:

```go
package main

import "github.com/solher/forklift/files"

func main() {
  println(files.File("file.sql"))
}

// Prints: SELECT * FROM documents
```

By default (development), the file will actually be read from the disk. In production, you can embbed the files into the binary by running the command:

```bash
$ forklift -extensions sql,gql > forklift.go
```

The files will then be read from memory.
