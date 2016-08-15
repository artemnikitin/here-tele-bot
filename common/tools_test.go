package common

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

func TestStartWith(t *testing.T) {
	cases := []struct {
		src, sub string
		res      bool
	}{
		{"qwe", "q", true},
		{"qwe", "qw", true},
		{"qwe", "qwer", false},
		{"abc", "w", false},
		{"", "d/", false},
		{"dfd/", "", true},
		{"", "", true},
	}

	for _, v := range cases {
		result := StringStartWith(v.src, v.sub)
		if result != v.res {
			t.Errorf("For string: %s start with: %s, actual: %v, expected: %v", v.src, v.sub, result, v.res)
		}
	}
}

func TestClearSlackMessage(t *testing.T) {
	cases := []struct {
		src, res string
	}{
		{"<!here|@here>.: chinese in berlin mitte", " chinese in berlin mitte"},
		{"<!here|@here>.:chinese in berlin mitte", "chinese in berlin mitte"},
		{"abc <!here|@here>.: def", "abc <!here|@here>.: def"},
		{"dfd/", "dfd/"},
		{"", ""},
	}

	for _, v := range cases {
		result := ClearSlackMessage(v.src)
		if result != v.res {
			t.Errorf("For string: %s, actual: %s, expected: %v", v.src, result, v.res)
		}
	}
}
