# sqlite3

[![Build Status](https://travis-ci.org/sbinet/sqlite3.svg?branch=master)](https://travis-ci.org/sbinet/sqlite3)

`sqlite3` is a pure Go package decoding the `SQLite` file format as
described by:
 http://www.sqlite.org/fileformat.html

## Installation

```sh
$ go get github.com/sbinet/sqlite3
```

## Example

```go
package main

import (
	"fmt"

	"github.com/sbinet/sqlite3"
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

## Documentation

Documentation is available on [godoc](http://godoc.org/github.com/sbinet/sqlite3)

