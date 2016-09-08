package common

import "testing"

func TestLocationToString(t *testing.T) {
	expected := "11.22,33.221"
	str := LocationToString(11.22, 33.221)
	if str != expected {
		t.Errorf("Incorrect transformation. Actual: %s, expected: %s", str, expected)
	}
}

func TestIsQueryCorrect(t *testing.T) {
	cases := []struct {
		src, word string
		res       bool
	}{
		{"qwe", "", false},
		{"eee near bridge", " near ", true},
		{"eeenear bridge", "", false},
		{"near x place", "", false},
		{"x near", "", false},
		{"", "", false},
	}

	for _, v := range cases {
		word, result := IsQueryCorrect(v.src)
		if result != v.res {
			t.Errorf("For string: %s actual: %v, expected: %v", v.src, result, v.res)
		}
		if word != v.word {
			t.Errorf("Actual: %v, expected: %v", word, v.word)
		}
	}
}

func TestSplitQueryAndLocation(t *testing.T) {
	cases := []struct{ in, s, q, l string }{
		{"aa in bb", " in ", "aa", "bb"},
		{"", "", "", ""},
		{"ededwdwe", "", "", ""},
		{"ededwdwe in ", " in ", "ededwdwe", ""},
		{" in c", " in ", "", "c"},
	}

	for _, v := range cases {
		q, l := SplitQueryAndLocation(v.in, v.s)
		if q != v.q {
			t.Errorf("Incorrect parsing query, actual: %s, expected: %s", q, v.q)
		}
		if l != v.l {
			t.Errorf("Incorrect parsing location, actual: %s, expected: %s", l, v.l)
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
