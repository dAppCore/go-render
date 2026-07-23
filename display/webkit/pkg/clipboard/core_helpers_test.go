package clipboard

func bytesRepeat(value []byte, count int) []byte {
	if count <= 0 || len(value) == 0 {
		return nil
	}
	out := make([]byte, 0, len(value)*count)
	for i := 0; i < count; i++ {
		out = append(out, value...)
	}
	return out
}
