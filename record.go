// Copyright 2017 The go-sqlite Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sqlite3

type RecordHeader struct {
	Len   int
	Types []SerialType
}

type Record struct {
	Header RecordHeader
	Body   []byte
	Values []interface{}
}
