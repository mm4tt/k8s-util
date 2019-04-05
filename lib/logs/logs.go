package logs

import "strings"

// TODO(mm4tt): Document
func Shorten(line, pattern string, nCharsAround int) string {
	i := strings.Index(line, pattern)

	type rangeType struct{ a, b int }

	ranges := []rangeType{
		{0, 100},
		{i - nCharsAround, i + len(pattern) + nCharsAround},
		{len(line) - 50, len(line)},
	}

	// Compact ranges.
	curr := ranges[0]
	ranges2 := []rangeType{}
	for i := 1; i < len(ranges); i++ {
		this := ranges[i]
		if curr.b < this.a {
			ranges2 = append(ranges2, curr)
			curr = this
		} else {
			curr = rangeType{curr.a, this.b}
		}
	}
	ranges2 = append(ranges2, curr)

	ret := []string{}
	for _, r := range ranges2 {
		ret = append(ret, line[r.a:r.b])
	}

	return strings.Join(ret, " ... ")
}
