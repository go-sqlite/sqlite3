// Copyright 2017 The go-sqlite Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sqlite3

import "testing"

func TestVarint(t *testing.T) {
	testdata := []struct {
		in      []byte
		valWant int64
		nWant   int
	}{
		{[]byte{0x7f}, 0x7f, 1},
		{[]byte{0x7f, 0x42}, 0x7f, 1},
		{[]byte{0x81}, 0x0, 0},
		{[]byte{0x81, 0x1}, 0x81, 2},
		{[]byte{0x83, 0x60}, 0x1e0, 2},
		{[]byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff}, -1, 9},
		{[]byte{0x97, 0xa0, 0x80, 0xdf, 0xe9, 0xda, 0xf3, 0x02}, 13088612140104066, 8},
	}

	for _, tt := range testdata {
		valGot, nGot := varint(tt.in)

		if tt.valWant != valGot || tt.nWant != nGot {
			t.Errorf("want varint(%v) = (%d, %d); got (%d, %d)", tt.in, tt.valWant, tt.nWant, valGot, nGot)
		}
	}
}
