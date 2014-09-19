package sqlite3

import (
	"fmt"
	"os"

	"github.com/gonuts/binary"
)

const (
	sqlite3Magic = "SQLite format 3\x00"
)

type File struct {
	f      *os.File
	header dbHeader
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

func Open(fname string) (*File, error) {
	f, err := os.Open(fname)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err != nil {
			f.Close()
		}
	}()

	db := &File{f: f}

	dec := binary.NewDecoder(db.f)
	dec.Order = binary.BigEndian
	err = dec.Decode(&db.header)
	if err != nil {
		return nil, err
	}

	fmt.Printf("db: %#v\n", db.header)
	_, err = db.f.Seek(0, 0)
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
	return db, err
}

func (f *File) Close() error {
	return f.f.Close()
}

// PageSize returns the database page size in bytes
func (f *File) PageSize() int {
	return int(f.header.PageSize)
}

// NumPage returns the number of pages for this database
func (f *File) NumPage() int {
	return int(f.header.DbSize)
}

// Encoding returns the text encoding for this database
func (f *File) Encoding() int {
	return int(f.header.DbEncoding)
}

// Version returns the sqlite version number used to create this database
func (f *File) Version() int {
	return int(f.header.SqliteVersion)
}
