# sqlite3

[![GoDoc](https://godoc.org/github.com/go-sqlite/sqlite3?status.svg)](https://godoc.org/github.com/go-sqlite/sqlite3)
[![Build Status](https://travis-ci.org/go-sqlite/sqlite3.svg?branch=master)](https://travis-ci.org/go-sqlite/sqlite3)

`sqlite3` is a pure Go package decoding the `SQLite` file format as
described by:
 http://www.sqlite.org/fileformat.html

## Current status

**WIP**: The near-term aim for `sqlite3` is to iterate through the
data in tables in `SQLite` files: ie., readonly access, and no actual
SQL queries.

It doesn't quite do that yet: so far it just parses the
`sqlite_master` data enough to find a list of tables and their names.

## Installation

```sh
$ go get github.com/go-sqlite/sqlite3
```

## License

`sqlite3` is released under the `BSD-3` license.


## Example

```go
package main

import (
	"fmt"

	"github.com/go-sqlite/sqlite3"
)

func main() {
	db, err := sqlite3.Open("test.sqlite")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	for _, table := range db.Tables() {
		fmt.Printf(">>> table=%#v\n", table)
	}
}
```

## Contributing

We're always looking for new contributing finding bugs, fixing issues, or writing some docs. If you're interested in contriburing source code changes you'll just need to [pull down the source code](#installation). You can run tests with `go test ./...` in the root of this project.

Make sure to add yourself to `AUTHORS` and `CONTRIBUTORS` if you submit a PR. We want you to take credit for your work!
