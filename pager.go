package sqlite3

type pager struct {
	f     *File
	pages map[int]Page
}

func newPager(f *File) pager {
	pager := pager{
		f:     f,
		pages: make(map[int]Page, f.NumPage()),
	}

	return pager
}

func (p *pager) Page(i int) Page {
	page, ok := p.pages[i]
	if ok {
		return page
	}

	if i > p.f.NumPage() {
		panic("out of range")
	}

	pos, _ := p.f.f.Seek(0, 1)
	defer p.f.f.Seek(pos, 0)

	buf := make([]byte, p.f.PageSize())
	n, err := p.f.f.ReadAt(buf, int64(i*p.f.PageSize()))
	if err != nil {
		panic(err)
	}

	if n != len(buf) {
		panic("read too few bytes")
	}

	pb := pageBuffer{
		pos: 0,
		buf: buf,
	}

	page, err = newPage(i, pb)
	if err != nil {
		panic(err)
	}

	return page
}
