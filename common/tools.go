package common

import (
	"bytes"
	"strconv"
	"strings"
)

var words = []string{" in ", " near ", " nearby ", " near to ", " around ", " close ", " close to ", " next to "}

func LocationToString(lat, lon float64) string {
	var buf bytes.Buffer
	buf.WriteString(strconv.FormatFloat(lat, 'f', -1, 64))
	buf.WriteString(",")
	buf.WriteString(strconv.FormatFloat(lon, 'f', -1, 64))
	return buf.String()
}

func IsQueryCorrect(query string) (string, bool) {
	var result bool
	var word string
	for _, v := range words {
		if strings.Contains(query, v) {
			word = v
			result = true
			break
		}
	}
	return word, result
}

func SplitQueryAndLocation(text, spl string) (string, string) {
	if spl == "" {
		return "", ""
	}
	pos := strings.Index(text, spl)
	if pos == -1 {
		return "", ""
	}
	return text[:pos], text[pos+len(spl):]
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
