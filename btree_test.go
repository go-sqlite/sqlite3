package sqlite3

import (
	"reflect"
	"testing"
)

func TestBTree__readStInt8(t *testing.T) {
	rbuf, res := readStInt8([]byte{0x7f})
	if len(rbuf) != 0 {
		t.Errorf("len(rbuf)=%d", len(rbuf))
	}
	ans := int8(1<<7-1)
	if res != ans {
		t.Errorf("got %d, expected %d", res, ans)
	}
}

func TestBTree__readStInt16(t *testing.T) {
	rbuf, res := readStInt16([]byte{0x7f, 0xff})
	if len(rbuf) != 0 {
		t.Errorf("len(rbuf)=%d", len(rbuf))
	}
	ans := int16(1<<15-1)
	if res != ans {
		t.Errorf("got %d, expected %d", res, ans)
	}
}

func TestBTree__readStInt24(t *testing.T) {
	rbuf, res := readStInt24([]byte{0x7f, 0xff, 0xff})
	if len(rbuf) != 0 {
		t.Errorf("len(rbuf)=%d", len(rbuf))
	}
	ans := uint32(1<<23-1)
	if res != ans {
		t.Errorf("got %d, expected %d", res, ans)
	}
}

func TestBTree__readStInt32(t *testing.T) {
	rbuf, res := readStInt32([]byte{0x7f, 0xff, 0xff, 0xff})
	if len(rbuf) != 0 {
		t.Errorf("len(rbuf)=%d", len(rbuf))
	}
	ans := int32(1<<31-1)
	if res != ans {
		t.Errorf("got %d, expected %d", res, ans)
	}
}

func TestBTree__readStInt48(t *testing.T) {
	rbuf, res := readStInt48([]byte{0x7f, 0xff, 0xff, 0xff, 0xff, 0xff})
	if len(rbuf) != 0 {
		t.Errorf("len(rbuf)=%d", len(rbuf))
	}
	ans := uint64(1<<47-1)
	if res != ans {
		t.Errorf("got %d, expected %d", res, ans)
	}
}

func TestBTree__readStInt64(t *testing.T) {
	rbuf, res := readStInt64([]byte{0x7f, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff})
	if len(rbuf) != 0 {
		t.Errorf("len(rbuf)=%d", len(rbuf))
	}
	ans := int64(1<<63-1)
	if res != ans {
		t.Errorf("got %d, expected %d", res, ans)
	}
}

func TestBTree__decodeRecord(t *testing.T) {
	cases := []struct {
		bt func() *btreeTable
		payload []byte
		match func (rec Record, err error) bool
	}{
		{
			bt: func() *btreeTable{
				return &btreeTable{}
			},
			payload: []byte{},
			match: func (rec Record, err error) bool {
				return err != nil
			},
		},
		{

			bt: func() *btreeTable {
				db, _ := Open("testdata/firefox-history.sqlite")
				page, _ := db.pager.Page(1)
				bt, _ := newBtreeTable(page, db)
				return bt
			},
			// From testdata/firefox-history.sqlite
			payload: []byte{
				6,23,75,37,1,0,105,110,100,101,120,115,113,108,105,116,101,95,97,
				117,116,111,105,110,100,101,120,95,109,111,122,95,107,101,121,119,
				111,114,100,115,95,49,109,111,122,95,107,101,121,119,111,114,100,
				115,26,
			},
			match: func(rec Record, err error) bool {
				body := []byte{
					105,110,100,101,120,115,113,108,105,116,101,95,97,117,116,111,
					105,110,100,101,120,95,109,111,122,95,107,101,121,119,111,114,
					100,115,95,49,109,111,122,95,107,101,121,119,111,114,100,115,26,
				}
				values:= []interface{}{"index", "sqlite_autoindex_moz_keywords_1", "moz_keywords", int8(26), nil}
				return rec.Header.Len == 5 && reflect.DeepEqual(rec.Values, values) && reflect.DeepEqual(rec.Body, body) && err == nil
			},
		},
	}
	for i := range cases {
		rec, err := cases[i].bt().decodeRecord(cases[i].payload)
		if !cases[i].match(rec, err) {
			t.Errorf("rec=%v\nerr=%v\n", rec, err)
		}
	}
}
