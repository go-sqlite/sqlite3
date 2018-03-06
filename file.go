// Copyright 2017 The go-sqlite Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sqlite3

import (
	"fmt"
	"io"
	"os"
	"reflect"
	"strings"

	"github.com/gonuts/binary"
)

const (
	sqlite3Magic = "SQLite format 3\x00"

	printfDebug = false
)

var (
	tblconstraints = []string{"CHECK", "FOREIGN KEY", "UNIQUE", "PRIMARY KEY"}
)

type DbFile struct {
	pager  pager
	header dbHeader
	tables []Table
	close  func() error
}

type dbHeader struct {
	Magic         [16]byte
	PageSize      uint16  // database page size in bytes
	WVersion      byte    // file format write version
	RVersion      byte    // file format read version
	NReserved     byte    // bytes of unused reserved space at the end of each page
	MaxFraction   byte    // maximum embedded payload fraction (must be 64)
	MinFraction   byte    // minimum embedded payload fraction (must be 32)
	LeafFraction  byte    // leaf payload fraction (must be 32)
	NFileChanges  int32   // file change counter
	DbSize        int32   // size of the database file in pages. The "in-header database size".
	FreePage      int32   // page number of the first freelist trunk page.
	NFreePages    int32   // total number of freelist pages.
	SchemaCookie  [4]byte // schema cookie
	SchemaFormat  int32   // schema format number. supported formats are 1,2,3 and 4.
	PageCacheSize int32   // default page cache size
	AutoVacuum    int32   // page number of the largest root b-tree page when in auto-vacuum or incremental-vacuum modes, or zero otherwise.

	// the database text encoding.
	//  1: UTF-8.
	//  2: UTF-16le
	//  3: UTF-16be
	DbEncoding int32

	UserVersion int32 // the "user version" as read and set by the user_version PRAGMA
	IncrVacuum  int32 // tree (non-zero) for incremental-vacuum mode. False (zero) otherwise

	ApplicationID int32 // the "Application ID" set by the PRAGMA application_id

	XXX_reserved  [20]byte // reserved for expansion. must be zero
	VersionValid  int32    // the version-valid-for number
	SqliteVersion int32    // SQLITE_VERSION_NUMBER
}

func OpenFrom(f io.ReadSeeker) (*DbFile, error) {
	var db DbFile

	dec := binary.NewDecoder(f)
	dec.Order = binary.BigEndian
	err := dec.Decode(&db.header)
	if err != nil {
		return nil, err
	}

	if db.header.DbSize == 0 {
		// determine it based on the size of the database file.
		// if the size of the database file is not an integer multiple of
		// the page-size, round down to the nearest page.
		// except, any file larger than 0-bytes in size, is considered to
		// contain at least one page.
		size, err := f.Seek(0, io.SeekEnd)
		if err != nil {
			return nil, err
		}
		if _, err := f.Seek(0, io.SeekStart); err != nil {
			return nil, err
		}
		pagesz := int64(db.header.PageSize)
		npages := (size + pagesz - 1) / pagesz
		db.header.DbSize = int32(npages)
	}

	if printfDebug {
		fmt.Printf("db: %#v\n", db.header)
	}
	_, err = f.Seek(0, 0)
	if err != nil {
		return nil, err
	}

	if string(db.header.Magic[:]) != sqlite3Magic {
		return nil, fmt.Errorf(
			"sqlite: invalid file header.\ngot:  %q\nwant: %q\n",
			string(db.header.Magic[:]),
			sqlite3Magic,
		)
	}

	db.pager = newPager(f, db.PageSize(), db.NumPage())

	err = db.init()
	if err != nil {
		return nil, err
	}

	return &db, err
}

func Open(fname string) (*DbFile, error) {
	f, err := os.Open(fname)
	if err != nil {
		return nil, err
	}

	db, err := OpenFrom(f)
	if err != nil {
		f.Close()
		return nil, err
	}
	db.close = f.Close
	return db, nil
}

func (db *DbFile) Close() error {
	db.pager.Delete()
	if db.close != nil {
		return db.close()
	}
	return nil
}

// PageSize returns the database page size in bytes
func (db *DbFile) PageSize() int {
	return int(db.header.PageSize)
}

// NumPage returns the number of pages for this database
func (db *DbFile) NumPage() int {
	return int(db.header.DbSize)
}

// Encoding returns the text encoding for this database
func (db *DbFile) Encoding() int {
	return int(db.header.DbEncoding)
}

// Version returns the sqlite version number used to create this database
func (db *DbFile) Version() int {
	return int(db.header.SqliteVersion)
}

func (db *DbFile) Tables() []Table {
	return db.tables
}

func (db *DbFile) init() error {

	// load sqlite_master
	page, err := db.pager.Page(1)
	if err != nil {
		return err
	}

	if page.Kind() != BTreeLeafTableKind && page.Kind() != BTreeInteriorTableKind {
		return fmt.Errorf("sqlite3: invalid page kind (%v)", page.Kind())
	}

	btree, err := newBtreeTable(page, db)
	if err != nil {
		return err
	}

	if printfDebug {
		fmt.Printf(">>> bt-hdr: %#v\n", btree.btheader)
		fmt.Printf(">>> init... (ncells=%d)\n", btree.NumCell())
	}

	return btree.visitRecordsInorder(func(_ *int64, rec Record) error {
		// {"table", "tbl1", "tbl1", 2, "CREATE TABLE tbl1(one varchar(10), two smallint)"} (body=62)
		// {"table", "tbl2", "tbl2", 3, "CREATE TABLE tbl2(\n f1 varchar(30) primary key,\n f2 text,\n f3 real\n)"}
		if len(rec.Values) != 5 {
			return fmt.Errorf("sqlite3: invalid table format")
		}

		rectype := rec.Values[0].(string)
		if rectype != "table" {
			return nil
		}

		pageid := reflect.ValueOf(rec.Values[3])
		table := Table{
			name:   rec.Values[1].(string),
			pageid: int(pageid.Int()),
		}

		// skip internal tables, aka don't expose them
		if strings.HasPrefix(table.name, "sqlite_") {
			return nil
		}

		def := rec.Values[4].(string)
		def = strings.Replace(def, "CREATE TABLE "+table.name, "", 1)
		def = strings.Replace(def, "\n", "", -1)
		def = strings.TrimSpace(def)
		if def[0] == '(' {
			def = def[1:]
		}
		if def[len(def)-1] == ')' {
			def = def[:len(def)-1]
		}
		def = strings.TrimSpace(def)

		parts := strings.Split(def, ",")
		// strip away statements like 'UNIQUE ...' or 'PRIMARY KEY ...' from a table definition
		for i := range parts {
			if i >= len(parts) {
				break // we removed at least one elem, so avoid out of bounds read
			}

			parts[i] = strings.TrimSpace(parts[i])
			for j := range tblconstraints {
				if strings.HasPrefix(parts[i], tblconstraints[j]) {
					// drop all other elements
					parts = parts[:i]
					break
				}
			}
		}

		table.cols = make([]Column, len(parts))
		for i := range parts {
			parts[i] = strings.TrimSpace(parts[i])
			idx := strings.Index(parts[i], " ") // find where col name ends and type starts
			if idx > 0 {
				table.cols[i].name = parts[i][:idx]
			} else {
				table.cols[i].name = parts[i]
			}
		}

		if printfDebug {
			fmt.Printf(">>> def: %q => ncols=%d\n", def, len(table.cols))
		}

		db.tables = append(db.tables, table)
		return nil
	})
}

func (db *DbFile) Dumpdb() error {
	var err error
	for i := 1; i < db.NumPage(); i++ {
		page, err := db.pager.Page(i)
		if err != nil {
			fmt.Printf("error: sqlite3: error retrieving page-%d: %v\n", i, err)
			continue
		}
		fmt.Printf("page-%d: %v\n", i, page.Kind())
		btree, err := newBtreeTable(page, db)
		if err != nil {
			fmt.Printf("** error: %v\n", err)
			continue
		}
		for i := 0; i < btree.NumCell(); i++ {
			cell, err := btree.loadCell(i)
			if err != nil {
				fmt.Printf("** error: %v\n", err)
				continue
			}
			fmt.Printf("--- cell[%03d/%03d]= leftchildpage=%d row=%d payload=%d overflow=%d\n",
				i+1, btree.NumCell(),
				cell.LeftChildPage,
				cell.RowID,
				len(cell.Payload),
				cell.OverflowPage,
			)
		}
	}

	return err
}

// VisitTableRecords performs an inorder traversal of all cells in the
// btree for the table with the given name, passing the (optional,
// hence nullable) RowID, and record-decoded payload of each cell to
// the visitor function `f`.
func (db *DbFile) VisitTableRecords(tableName string, f func(*int64, Record) error) error {
	for _, table := range db.tables {
		if table.name != tableName {
			continue
		}
		page, err := db.pager.Page(table.pageid)
		if err != nil {
			return err
		}
		btree, err := newBtreeTable(page, db)
		if err != nil {
			return err
		}
		return btree.visitRecordsInorder(f)
	}
	return fmt.Errorf("unknown table %q", tableName)
}
