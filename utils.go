// Copyright 2017 The go-sqlite3 Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sqlite3

import (
	"bytes"
	bb "encoding/binary"

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

func uvarint(data []byte) (uint64, int) {
	return bb.Uvarint(data)
}
