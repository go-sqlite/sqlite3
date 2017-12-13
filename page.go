// Copyright 2017 The go-sqlite Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sqlite3

import "fmt"

// PageKind describes what kind of page is.
type PageKind byte

/*
An index B-Tree internal node
An index B-Tree leaf node
A table B-Tree internal node
A table B-Tree leaf node
An overflow page
A freelist page
A pointer map page
The locking page
*/

const (
	intKeyKind   PageKind = 0x01
	zeroDataKind PageKind = 0x02
	leafDataKind PageKind = 0x04
	leafKind     PageKind = 0x08

	BTreeInteriorIndexKind = zeroDataKind
	BTreeInteriorTableKind = leafDataKind | intKeyKind
	BTreeLeafIndexKind     = zeroDataKind | leafKind
	BTreeLeafTableKind     = leafDataKind | intKeyKind | leafKind

	pkLockByte
	pkFreelistTrunk
	pkFreelistLeaf
	pkPayloadOverflow
	pkPointerMap
)

func (pk PageKind) String() string {
	switch pk {
	case BTreeInteriorIndexKind:
		return "BTreeInteriorIndex"
	case BTreeInteriorTableKind:
		return "BTreeInteriorTable"
	case BTreeLeafIndexKind:
		return "BTreeLeafIndex"
	case BTreeLeafTableKind:
		return "BTreeLeafTable"
	}

	panic(fmt.Sprintf("sqlite3: invalid PageKind value (0x%02x)", byte(pk)))
}

/*
func newPage(i int, pb pageBuffer) (Page, error) {
	pk := PageKind(pb.buf[0])
	switch pk {
	case BTreeInteriorIndexKind:
		panic("not implemented")
	case BTreeInteriorTableKind:
		panic("not implemented")
	case BTreeLeafIndexKind:
		panic("not implemented")
	case BTreeLeafTableKind:
		return newBtreeLeafTable(i, pb)

	}

	panic(fmt.Errorf("invalid PageKind value (%d)", pk))
}
*/

// page is a page loaded from disk.
type page struct {
	id  int
	pos int
	buf []byte
}

func (p *page) ID() int {
	return p.id
}

func (p *page) Kind() PageKind {
	offset := 0
	if p.id == 1 {
		offset = 100
	}
	return PageKind(p.buf[0+offset])
}

func (p *page) PageSize() int {
	return len(p.buf)
}

func (p *page) Seek(offset int64, whence int) (ret int64, err error) {
	switch whence {
	case 0:
		offset := int(offset)
		if offset > len(p.buf) {
			return 0, fmt.Errorf("sqlite: offset too big (%d)", offset)
		}
		p.pos = offset
	case 1:
		offset := int(offset)
		pos := p.pos + offset
		if pos > len(p.buf) {
			return 0, fmt.Errorf("sqlite: offset too big (%d)", offset)
		}
		p.pos = pos
	case 2:
		offset := int(offset)
		pos := len(p.buf) - offset
		if pos < 0 {
			return 0, fmt.Errorf("sqlite: offset too big (%d)", offset)
		}
		p.pos = pos
	}
	return int64(p.pos), nil
}

func (p *page) Pos() int {
	return p.pos
}

func (p *page) Bytes() []byte {
	return p.buf[p.pos:]
}

func (p *page) Decode(ptr interface{}) error {
	n, err := unmarshal(p.buf[p.pos:], ptr)
	if err != nil {
		return err
	}
	p.pos += int(n)
	return err
}

func (p *page) Read(data []byte) (int, error) {
	n := copy(data, p.buf[p.pos:p.pos+len(data)])
	if n != len(data) {
		return n, fmt.Errorf("error. read too few bytes: %d. want %d", n, len(data))
	}
	p.pos += n
	return n, nil
}

func (p *page) Varint() (int64, int) {
	v, n := varint(p.Bytes())
	if n <= 0 {
		return v, n
	}
	p.pos += int(n)
	return v, n
}
