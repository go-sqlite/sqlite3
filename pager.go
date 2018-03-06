// Copyright 2017 The go-sqlite Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sqlite3

import (
	"fmt"
	"io"
)

type pager struct {
	f      io.ReadSeeker
	size   int          // page size in bytes
	npages int          // total number of pages in db
	pages  map[int]page // cache of pages
	lru    []int        // list of last used pages
}

func newPager(f io.ReadSeeker, size, npages int) pager {
	pager := pager{
		f:      f,
		size:   size,
		npages: npages,
		pages:  make(map[int]page, npages),
		lru:    make([]int, 0, 2),
	}

	return pager
}

func (p *pager) Page(i int) (page, error) {
	var err error
	page, ok := p.pages[i]
	if ok {
		return page, err
	}

	if i > p.npages {
		return page, fmt.Errorf("sqlite3: out of range (%d > %d)", i, p.npages)
	}

	pos, _ := p.f.Seek(0, io.SeekCurrent)
	defer p.f.Seek(pos, io.SeekStart)

	buf := make([]byte, p.size)
	if _, err := p.f.Seek(int64((i-1)*p.size), io.SeekStart); err != nil {
		return page, err
	}
	n, err := p.f.Read(buf)
	if err != nil {
		return page, err
	}

	if n != len(buf) {
		return page, fmt.Errorf("sqlite3: read too few bytes")
	}

	page.id = i
	page.buf = buf

	p.pages[i] = page
	p.lru = append(p.lru, i)
	return page, err
}

func (p *pager) Delete() error {
	var err error
	p.pages = nil
	p.lru = nil
	return err
}
