package sqlite3

import "testing"

func TestUvarint(t *testing.T) {
	testdata := []struct {
		in      []byte
		valWant uint64
		nWant   int
	}{
		{[]byte{0x7f}, 0x7f, 1},
		{[]byte{0x7f, 0x42}, 0x7f, 1},
		{[]byte{0x81}, 0x0, 0},
		{[]byte{0x81, 0x1}, 0x81, 2},
		{[]byte{0x83, 0x60}, 0x1e0, 2},
		{[]byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff}, 0xffffffffffffffff, 9},
	}

	for _, tt := range testdata {
		valGot, nGot := uvarint(tt.in)

		if tt.valWant != valGot || tt.nWant != nGot {
			t.Errorf("want uvarint(%v) = (%d, %d); got (%d, %d)", tt.in, tt.valWant, tt.nWant, valGot, nGot)
		}
	}
}
