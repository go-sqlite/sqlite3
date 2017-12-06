// Copyright 2017 The go-sqlite Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sqlite3

import (
	"reflect"
)

// Table is a SQLite table
type Table struct {
	name   string
	pageid int
	cols   []Column
}

// Name returns the name of the table
func (t *Table) Name() string {
	return t.name
}

// NumRow returns the number of rows in the table
func (t *Table) NumRow() int64 {
	//return t.nrows
	// FIXME(sbinet)
	return -1
}

// Columns returns the columns of the table
func (t *Table) Columns() []Column {
	return t.cols
}

// Column describes a column in a SQLite table
type Column struct {
	name string
	typ  reflect.Type
}

// Name returns the name of the column
func (col *Column) Name() string {
	return col.name
}

// Type returns the SQLite type of the column
func (col *Column) Type() reflect.Type {
	return col.typ
}
