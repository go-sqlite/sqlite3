package sqlite3

import (
	"reflect"
	"testing"
)

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
