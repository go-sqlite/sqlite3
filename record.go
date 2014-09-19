package sqlite3

type RecordHeader struct {
	Len   int
	Types []SerialType
}

type Record struct {
	Header RecordHeader
	Body   []byte
	Values []interface{}
}
