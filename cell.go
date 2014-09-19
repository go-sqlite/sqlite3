package sqlite3

// cellInfo holds information about an on-disk cell.
type cellInfo struct {
	Key          int64
	RowID        int64
	Payload      []byte
	OverflowPage int32
}
