// Copyright 2017 The go-sqlite Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sqlite3

import (
	"fmt"
)

// SerialType represents SQLite types on disk
type SerialType int

const (
	StNull SerialType = iota
	StInt8
	StInt16
	StInt24
	StInt32
	StInt48
	StInt64
	StFloat
	StC0
	StC1

	StBlob SerialType = 12
	StText            = 13
)

func (st SerialType) String() string {
	switch st {
	case StNull:
		return "StNull"
	case StInt8:
		return "StInt8"
	case StInt16:
		return "StInt16"
	case StInt24:
		return "StInt24"
	case StInt32:
		return "StInt32"
	case StInt48:
		return "StInt48"
	case StInt64:
		return "StInt64"
	case StFloat:
		return "StFloat"
	case StC0:
		return "StC0"
	case StC1:
		return "StC1"
	}

	if st.IsBlob() {
		sz := int(st-12) / 2
		return fmt.Sprintf("StBlob(%d)", sz)
	}
	if st.IsText() {
		sz := int(st-13) / 2
		return fmt.Sprintf("StText(%d)", sz)
	}

	panic("unreachable")
}

func (st SerialType) IsBlob() bool {
	return st >= 12 && st&1 == 0
}

func (st SerialType) IsText() bool {
	return st >= 13 && st&1 == 1
}

// NBytes returns the number of bytes on disk for this SerialType
// NBytes returns -1 if the SerialType is invalid.
func (st SerialType) NBytes() int {
	switch st {
	case StNull:
		return 0
	case StInt8:
		return 1
	case StInt16:
		return 2
	case StInt24:
		return 3
	case StInt32:
		return 4
	case StInt48:
		return 6
	case StInt64:
		return 8
	case StFloat:
		return 8
	case StC0, StC1:
		return 0
	}

	if st.IsBlob() {
		return int(st-12) / 2
	}
	if st.IsText() {
		return int(st-13) / 2
	}

	return -1
}
