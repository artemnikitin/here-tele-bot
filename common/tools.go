package common

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

func StringStartWith(original, substring string) bool {
	if len(substring) > len(original) {
		return false
	}
	str := string(original[0:len(substring)])
	return str == substring
}

func ClearSlackMessage(text string) string {
	str := "<!here|@here>.:"
	if !StringStartWith(text, str) {
		return text
	}
	index := strings.Index(text, str)
	if index == -1 {
		return text
	}
	return text[index+15:]
}
