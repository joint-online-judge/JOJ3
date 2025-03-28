package diff

// compareStrings compares two strings character by character, optionally ignoring whitespace.
func compareStrings(str1, str2 string, compareSpace bool) bool {
	if compareSpace {
		return str1 == str2
	}
	var i, j int
	l1 := len(str1)
	l2 := len(str2)
	for i < l1 && j < l2 {
		for i < l1 && isWhitespace(str1[i]) {
			i++
		}
		for j < l2 && isWhitespace(str2[j]) {
			j++
		}
		if i < l1 && j < l2 && str1[i] != str2[j] {
			return false
		}
		if i < l1 {
			i++
		}
		if j < l2 {
			j++
		}
	}
	for i < l1 && isWhitespace(str1[i]) {
		i++
	}
	for j < l2 && isWhitespace(str2[j]) {
		j++
	}
	return i == l1 && j == l2
}

func isWhitespace(b byte) bool {
	return b == ' ' ||
		b == '\t' ||
		b == '\n' ||
		b == '\r' ||
		b == '\v' ||
		b == '\f' ||
		b == 0x85 ||
		b == 0xA0
}
