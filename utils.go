package sqlite

func min(a, b int) int {
	if a > b {
		return b
	}
	return a
}
