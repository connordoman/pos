package md

import "unicode/utf8"

type Source string

func (s Source) Length() int {
	return utf8.RuneCountInString(string(s))
}

func (s Source) CharAt(index int) rune {
	if index < 0 || index >= s.Length() {
		return 0
	}
	r, _ := utf8.DecodeRuneInString(string(s[index:]))
	return r
}

func (s Source) Substring(start, end int) string {
	start = max(0, min(start, s.Length()-1))
	end = max(0, min(end, s.Length()))
	return string(s[start:end])
}

func (s Source) String() string {
	return string(s)
}
