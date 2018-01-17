// Copyright 2017 The go-sqlite Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sqlite3

import (
	"testing"
)

func TestFileOpen(t *testing.T) {

	for _, test := range []struct {
		fname   string
		version int
		npages  int
		pagesz  int
		tables []Table
	}{
		{
			fname:   "testdata/test-1.sqlite",
			version: 3008006,
			npages:  2,
			pagesz:  1024,
			tables: []Table{
				Table{
					name: "tbl1",
					pageid: 2,
					cols: []Column{
						Column{
							name: "one",
						},
						Column{
							name: "two",
						},
					},
				},
			},
		},
		{
			fname:   "testdata/test-2.sqlite",
			version: 3008006,
			npages:  4,
			pagesz:  1024,
			tables: []Table{
				Table{
					name: "tbl1",
					pageid: 2,
					cols: []Column{
						Column{
							name: "one",
						},
						Column{
							name: "two",
						},
					},
				},
				Table{
					name: "tbl2",
					pageid: 3,
					cols: []Column{
						Column{
							name: "f1",
						},
						Column{
							name: "f2",
						},
						Column{
							name: "f3",
						},
					},
				},
			},
		},
	} {
		f, err := Open(test.fname)
		if err != nil {
			t.Fatalf("%s: error: %v\n", test.fname, err)
		}
		defer f.Close()

		if f.Version() != test.version {
			t.Fatalf("%s: version=%d\nwant version=%d\n", test.fname, f.Version(), test.version)
		}

		if f.PageSize() != test.pagesz {
			t.Fatalf("%s: page size = %d. want=%d", test.fname, f.PageSize(), test.pagesz)
		}

		if f.NumPage() != test.npages {
			t.Fatalf("%s: num-pages = %d. want=%d", test.fname, f.NumPage(), test.npages)
		}

		// Check tables
		if len(f.Tables()) != len(test.tables) {
			t.Fatalf("%s: tables=%d, want=%d", test.fname, len(f.Tables()), len(test.tables))
		}
		for i := range f.Tables() {
			ftbl := f.Tables()[i]
			if ftbl.name != test.tables[i].name {
				t.Error("table name: got=%q, want=%q", ftbl.name, test.tables[i].name)
			}
			if ftbl.pageid != test.tables[i].pageid {
				t.Error("table pageid: got=%d, want=%d", ftbl.pageid, test.tables[i].pageid)
			}
			if len(ftbl.cols) != len(test.tables[i].cols) {
				t.Error("table %s cols: got=%v, want=%v", ftbl.name, ftbl.cols, test.tables[i].cols)
			}
			for j := range ftbl.cols {
				if ftbl.cols[j].name != test.tables[i].cols[j].name {
					t.Errorf("table %s column: got=%q, want=%q", ftbl.name, ftbl.cols[j].name, test.tables[i].cols[j].name)
				}
			}
		}
	}

}
