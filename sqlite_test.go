// Copyright 2017 The go-sqlite3 Authors.  All rights reserved.
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
	}{
		{
			fname:   "testdata/test-1.sqlite",
			version: 3008006,
			npages:  2,
			pagesz:  1024,
		},
		{
			fname:   "testdata/test-2.sqlite",
			version: 3008006,
			npages:  4,
			pagesz:  1024,
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
	}

}
