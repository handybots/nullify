package linkgen

import "strings"

const (
	bucketSize     = 24 // 16 777 216
	zeroWidthRune0 = '\u200c'
	zeroWidthRune1 = '\u200d'
)

func N(s string) int64 {
	var n, b int64 = 0, 1
	for _, r := range []rune(s) {
		if r == zeroWidthRune0 {
			n += b
		}
		b = b << 1
	}
	return n
}

func S(n int64) string {
	var b strings.Builder
	b.Grow(bucketSize)
	for i := 0; i != bucketSize; i++ {
		if n&(1<<i) != 0 {
			b.WriteRune(zeroWidthRune0)
		} else {
			b.WriteRune(zeroWidthRune1)
		}
	}
	return b.String()
}
