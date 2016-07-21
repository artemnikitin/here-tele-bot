package logic

import (
	"bytes"
	"strconv"
	"strings"
)

func LocationToString(lat, lon float64) string {
	var buf bytes.Buffer
	buf.WriteString(strconv.FormatFloat(lat, 'f', -1, 64))
	buf.WriteString(",")
	buf.WriteString(strconv.FormatFloat(lon, 'f', -1, 64))
	return buf.String()
}

func SplitQueryAndLocation(in string) (string, string) {
	pos := strings.Index(in, " in ")
	if pos == -1 {
		return "", ""
	}
	return in[:pos], in[pos+4:]
}
