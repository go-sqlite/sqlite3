// Copyright 2017 The go-sqlite Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sqlite3

// cellInfo holds information about an on-disk cell.
type cellInfo struct {
	LeftChildPage int32
	RowID         *int64
	Payload       []byte
	OverflowPage  int32
}
