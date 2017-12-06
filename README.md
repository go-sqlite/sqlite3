# sqlite3

[![GoDoc](https://godoc.org/github.com/go-sqlite/sqlite3?status.svg)](https://godoc.org/github.com/go-sqlite/sqlite3)
[![Build Status](https://travis-ci.org/go-sqlite/sqlite3.svg?branch=master)](https://travis-ci.org/go-sqlite/sqlite3)

`sqlite3` is a pure Go package decoding the `SQLite` file format as
described by:
 http://www.sqlite.org/fileformat.html

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
