package kpstemmer

import "bytes"

func replaceString(src string, start int, end int, sub string) string {
	var buf bytes.Buffer
	buf.Grow(len(src) + len(sub))
	buf.WriteString(src[:start])
	buf.WriteString(sub)
	buf.WriteString(src[end:])
	return buf.String()
}
