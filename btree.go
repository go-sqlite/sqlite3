package sqlite3

import (
	"fmt"
	"strings"
)

type btheader struct {
	raw struct {
		Kind PageKind // b-tree page kind

		FreeBlockOffset int16 // byte offset into the page of the first free block
		NCells          int16 // number of cells on this page
		CellsOffset     int16 // offset into first byte of the cell content area
		NFreeBytes      int8  // number of fragmented free bytes within cell area content
	}
}

func (bt btheader) Kind() PageKind {
	return bt.raw.Kind
}

func (bt *btheader) NumCell() int {
	return int(bt.raw.NCells)
}

func (bt *btheader) FreeBlockAddr() int {
	return int(bt.raw.FreeBlockOffset)
}

func (bt *btheader) CellsAddr() int {
	return int(bt.raw.CellsOffset)
}

type btreeTable struct {
	btheader
	dbhdr   *dbHeader
	pointer int32   // right most pointer (only valid for interior pages
	page    page    // page backing this b-tree leaf
	addrs   []int16 // cell addresses
}

func newBtreeTable(page page, dbhdr *dbHeader) (*btreeTable, error) {
	var hdr btheader
	if page.ID() == 1 {
		// drop first 100-bytes (global file header)
		_, err := page.Seek(100, 0)
		if err != nil {
			return nil, err
		}
	}

	err := page.Decode(&hdr.raw)
	if err != nil {
		return nil, err
	}

	btree := &btreeTable{
		btheader: hdr,
		dbhdr:    dbhdr,
		page:     page,
	}

	if btree.Kind() == BTreeInteriorTableKind {
		err = btree.page.Decode(&btree.pointer)
		if err != nil {
			return nil, err
		}
	}

	err = btree.init()
	if err != nil {
		return nil, err
	}

	return btree, err
}

func (btree *btreeTable) ID() int {
	return btree.page.ID()
}

func (btree *btreeTable) Size() int {
	return len(btree.page.buf)
}

func (btree *btreeTable) init() error {
	var err error
	if btree.addrs != nil {
		return nil
	}

	cells := make([]int16, btree.NumCell())
	for icell, addr := range cells {
		fmt.Printf("   cell= %d/%d... (%d)\n", icell+1, len(cells), btree.page.Pos())
		err = btree.page.Decode(&addr)
		if err != nil {
			return err
		}
		fmt.Printf("   cell= %d/%d... (%d) => %d\n", icell+1, len(cells), btree.page.Pos(), addr)
		cells[icell] = addr
	}

	btree.addrs = cells
	return err
}

func (btree *btreeTable) load(icell int) (Record, error) {
	var rec Record
	addr := btree.addrs[icell]
	_, err := btree.page.Seek(int64(addr), 0)
	if err != nil {
		return rec, err
	}

	cell, err := btree.parseCell(icell)
	if err != nil {
		return rec, err
	}

	// decode record
	recbuf := cell.Payload[:]
	rhdrsz, n := uvarint(recbuf)
	if n <= 0 {
		return rec, fmt.Errorf("sqlite3: error decoding record header (n=%d)", n)
	}
	recbuf = recbuf[n:]

	rec.Header.Len = int(rhdrsz) - n
	for ii := 0; ii < rec.Header.Len; {
		v, n := uvarint(recbuf)
		// fmt.Printf("ii=%d nn=%d len=%d\n", ii, n, rec.Header.Len)
		if n < 0 {
			return rec, fmt.Errorf("sqlite3: error decoding record header type (n=%d)", n)
		}
		if n == 0 {
			break
		}
		recbuf = recbuf[n:]
		ii += int(n)
		rec.Header.Types = append(rec.Header.Types, SerialType(v))
	}
	rec.Body = recbuf[:]
	//copy(rec.Body, recbuf)

	// fmt.Printf(">>> record: %#v (body=%d)\n", rec.Header, len(rec.Body))
	for _, st := range rec.Header.Types {
		var v interface{}
		switch st {
		case StInt8:
			var vv int8
			n, err := unmarshal(recbuf, &vv)
			if err != nil {
				panic(err)
			}
			recbuf = recbuf[int(n):]
			v = vv

		case StInt16:
			var vv int16
			n, err := unmarshal(recbuf, &vv)
			if err != nil {
				panic(err)
			}
			recbuf = recbuf[int(n):]
			v = vv

		case StInt32:
			var vv int32
			n, err := unmarshal(recbuf, &vv)
			if err != nil {
				panic(err)
			}
			recbuf = recbuf[int(n):]
			v = vv

		case StInt64:
			var vv int64
			n, err := unmarshal(recbuf, &vv)
			if err != nil {
				panic(err)
			}
			recbuf = recbuf[int(n):]
			v = vv

		case StFloat:
			var vv float64
			n, err := unmarshal(recbuf, &vv)
			if err != nil {
				panic(err)
			}
			recbuf = recbuf[int(n):]
			v = vv

		case StC0:
			v = 0

		case StC1:
			v = 1

		case StInt24, StInt48:
			panic("not implemented")

		default:
			if st.IsBlob() {
				vv := make([]byte, st.NBytes())
				n := copy(vv, recbuf)
				recbuf = recbuf[int(n):]
				v = vv
			}
			if st.IsText() {
				vv := make([]byte, st.NBytes())
				n := copy(vv, recbuf)
				recbuf = recbuf[int(n):]
				// FIXME(sbinet)
				// handle db string encoding
				switch btree.dbhdr.DbEncoding {
				case 1:
					s := string(vv)
					idx := strings.Index(s, "\x00")
					if idx >= 0 {
						s = s[:idx]
					}
					v = s
				default:
					panic("utf-16 not supported")
				}
			}
		}

		rec.Values = append(rec.Values, v)
	}
	// fmt.Printf(">>> record: %#v (body=%d)\n", rec.Values, len(rec.Body))

	return rec, err
}

func (btree *btreeTable) parseCell(icell int) (cellInfo, error) {
	var cell cellInfo
	var err error

	switch btree.Kind() {
	case BTreeInteriorIndexKind:
		panic("not implemented")
	case BTreeInteriorTableKind:
		var pgno uint32 // page number of left child
		err = btree.page.Decode(&pgno)
		if err != nil {
			return cell, fmt.Errorf("sqlite3: error decoding page number: %v", err)
		}

		rowid, nrow := btree.page.Uvarint()
		if nrow <= 0 {
			return cell, fmt.Errorf("sqlite3: error decoding rowid: n=%d", nrow)
		}

		cell = cellInfo{
			Key:   int64(pgno),
			RowID: int64(rowid),
		}

	case BTreeLeafIndexKind:
		panic("not implemented")
	case BTreeLeafTableKind:
		sz, nsz := btree.page.Uvarint()
		if nsz <= 0 {
			return cell, fmt.Errorf("sqlite3: error decoding cell size: n=%d", nsz)
		}

		rowid, nrow := btree.page.Uvarint()
		if nrow <= 0 {
			return cell, fmt.Errorf("sqlite3: error decoding rowid: n=%d", nrow)
		}

		// sz is the total payload size.
		// check if all of it is in the b-tree leaf page or
		// if it spilled over to other pages
		localsz := int(sz)
		U := btree.page.PageSize() - int(btree.dbhdr.NReserved)
		M := int(((U - 12) * 32 / 255.) - 23)
		P := int(sz)
		if P > U-35 {
			vv := M + ((P - M) % (U - 4))
			localsz = min(vv, U-35)
		}

		// FIXME(sbinet): only create a new payload []byte when non-local
		// ie: when there is an overflow page
		payload := make([]byte, localsz, localsz)
		n, err := btree.page.Read(payload)
		if err != nil {
			return cell, err
		}
		if n != localsz {
			return cell, fmt.Errorf("read too few bytes: %d. want %d", n, localsz)
		}

		cell = cellInfo{
			Key:     int64(sz),
			RowID:   int64(rowid),
			Payload: payload,
		}

		if localsz != P {
			err = btree.page.Decode(&cell.OverflowPage)
			if err != nil {
				return cell, err
			}
			// FIXME(sbinet)
			// - locate overflow-page (and following)
			// - load payload
			// - append into cell.Payload

			panic("not implemented")
		}

		// fmt.Printf(" => size=%d rowid=%d overflow=%d (bytes: %d %d) [%d]\n",
		// 	cell.Len,
		// 	cell.RowID,
		// 	cell.OverflowPage,
		// 	nsz, nrow,
		// 	btree.page.Pos(),
		// )
		// fmt.Printf(" => %x (%d|%d)\n", string(cell.Payload), len(cell.Payload), localsz)

	}
	return cell, err
}

type btreeInteriorTable struct {
	btheader
	pointer int32 // right most pointer

	id   int
	page page
}
