// Copyright 2017 The go-sqlite Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sqlite3

import (
	"bytes"

	"github.com/gonuts/binary"
)

func min(a, b int) int {
	if a > b {
		return b
	}
	return a
}

func unmarshal(buf []byte, ptr interface{}) (int64, error) {
	r := bytes.NewReader(buf)
	max := r.Len()
	dec := binary.NewDecoder(r)
	dec.Order = binary.BigEndian
	err := dec.Decode(ptr)
	n := max - r.Len()
	return int64(n), err
}

func varint(data []byte) (int64, int) {
	var val uint64
	for i := 0; i < 8; i++ {
		if i > len(data)-1 {
			return 0, 0
		}
		val = (val << 7) | uint64(data[i]&0x7f)
		if data[i] < 0x80 {
			return int64(val), i + 1
		}
	}
	if len(data) < 9 {
		return 0, 0
	}
	return int64((val << 8) | uint64(data[8])), 9
}
