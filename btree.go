package sqlite3

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

type btreeLeafTable struct {
	btheader
	id   int
	page pageBuffer
}

func newBtreeLeafTable(i int, pb pageBuffer) (*btreeLeafTable, error) {
	var hdr btheader
	err := pb.Decode(&hdr.raw)
	if err != nil {
		return nil, err
	}
	return &btreeLeafTable{
		btheader: hdr,
		id:       i,
		page:     pb,
	}, err
}

func (btree *btreeLeafTable) ID() int {
	return btree.id
}

func (btree *btreeLeafTable) Size() int {
	return len(btree.page.buf)
}

type btreeInteriorTable struct {
	btheader
	pointer int32 // right most pointer

	id   int
	page pageBuffer
}
