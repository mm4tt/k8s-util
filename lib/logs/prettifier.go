package logs

import (
	"fmt"
	"strings"

	"github.com/nhooyr/color"
)

type Prettifier struct {
	LineLengthLimit    int
	CharsAroundPattern int
	PatternColor       string
}

// DefaultPrettifier creates a default Prettifier.
func DefaultPrettifier() *Prettifier {
	return &Prettifier{
		LineLengthLimit:    300,
		CharsAroundPattern: 50,
		PatternColor:       "Green",
	}
}

// Shorten shortens the provided line to make it human readable. It respects the provided pattern
// in a sense that the pattern will be always present in the result.
func (p *Prettifier) Shorten(line, pattern string) string {
	if len(line) <= p.LineLengthLimit {
		return line
	}
	i := strings.Index(line, pattern)
	type rangeType struct{ a, b int }
	ranges := []rangeType{
		{0, 100},
		{i - p.CharsAroundPattern, i + len(pattern) + p.CharsAroundPattern},
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

// Prettify prettifies the provided line. It shortens the line and colors the pattern occurences.
func (p *Prettifier) PrettyPrint(line, pattern string) {
	line = p.Shorten(line, pattern)
	coloredPattern := fmt.Sprintf("%%h[fg%s]%s%%r", p.PatternColor, pattern)
	line = strings.Replace(line, "%", "%%", -1)
	line = strings.Replace(line, pattern, coloredPattern, -1)
	color.Println(line)
}
