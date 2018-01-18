// Copyright 2018 The go-sqlite Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Command sqlite-dump is a simple command that dumps the high-level content
// of a SQLITE3 file on screen.
//
// Example:
//
//  $> sqlite-dump ./testdata/test-1.sqlite
//  sqlite3: opening "./testdata/test-1.sqlite"...
//  sqlite3: version: 3008006
//  sqlite3: page size: 1024
//  sqlite3: num pages: 2
//  sqlite3: num tables: 1
//  sqlite3: === table[0] ===
//  sqlite3: name: "tbl1"
//  sqlite3: cols: 2
//  sqlite3: col[0]: "one"
//  sqlite3: col[1]: "two"
package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/go-sqlite/sqlite3"
)

func main() {
	log.SetPrefix("sqlite3: ")
	log.SetFlags(0)

	flag.Parse()

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr,
			`Usage: sqlite-dump [options] file1 [file2 [...]]

Ex:

 $> sqlite-dump ./testdata/test-1.sqlite
 sqlite3: opening "./testdata/test-1.sqlite"...
 sqlite3: version: 3008006
 sqlite3: page size: 1024
 sqlite3: num pages: 2
 sqlite3: num tables: 1
 sqlite3: === table[0] ===
 sqlite3: name: "tbl1"
 sqlite3: cols: 2
 sqlite3: col[0]: "one"
 sqlite3: col[1]: "two"

Options:
`,
		)
	}

	if flag.NArg() < 1 {
		flag.Usage()
		flag.PrintDefaults()
		os.Exit(1)
	}

	for _, fname := range flag.Args() {
		process(fname)
	}
}

func process(fname string) {
	log.Printf("opening %q...", fname)
	f, err := sqlite3.Open(fname)
	if err != nil {
		log.Printf("error opening %q: %v", fname, err)
		return
	}
	defer f.Close()

	log.Printf("version: %v", f.Version())
	log.Printf("page size: %d", f.PageSize())
	log.Printf("num pages: %d", f.NumPage())
	log.Printf("num tables: %d", len(f.Tables()))

	for i, table := range f.Tables() {
		log.Printf("=== table[%d] ===", i)
		log.Printf("name: %q", table.Name())
		log.Printf("cols: %d", len(table.Columns()))
		for j, col := range table.Columns() {
			log.Printf("col[%d]: %q", j, col.Name())
		}
	}
}
