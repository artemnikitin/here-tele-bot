package logic

import "testing"

func TestLocationToString(t *testing.T) {
	expected := "11.22,33.221"
	str := LocationToString(11.22, 33.221)
	if str != expected {
		t.Errorf("Incorrect transformation. Actual: %s, expected: %s", str, expected)
	}
}

func TestSplitQueryAndLocation(t *testing.T) {
	cases := []struct{ in, q, l string }{
		{"aa in bb", "aa", "bb"},
		{"", "", ""},
		{"ededwdwe", "", ""},
		{"ededwdwe in ", "ededwdwe", ""},
		{" in c", "", "c"},
	}

	for _, v := range cases {
		q, l := SplitQueryAndLocation(v.in)
		if q != v.q {
			t.Errorf("Incorrect parsing query, actual: %s, expected: %s", q, v.q)
		}
		if l != v.l {
			t.Errorf("Incorrect parsing location, actual: %s, expected: %s", l, v.l)
		}
	}
}
