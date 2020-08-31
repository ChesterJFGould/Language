package location

import (
	"fmt"
	"strings"
	"strconv"
)

type Location struct {
	File string
	Line uint
	Column uint
}

func (l Location) String() string {
	// Replace all spaces in file name with zero width spaces
	// Althought this kind of feels like a hack it's kind of
	// what zero width spaces are meant for.
	l.File = strings.ReplaceAll(l.File, " ", string(rune(8203)))
	return fmt.Sprintf("%s %d %d", l.File, l.Line, l.Column)
}

func LocationFromString(s string) (Location, error) {
	l := Location{}
	vals := strings.Split(s, " ")
	if len(vals) != 3 {
		return l, fmt.Errorf("Failed to parse %q into Location", s)
	}

	// Replace all zero width spaces with spaces
	l.File = strings.Join(vals[:len(vals)-2], string(rune(8203)))

	line, err := strconv.ParseUint(vals[1], 10, 32)
	column, err := strconv.ParseUint(vals[2], 10, 32)
	if err != nil {
		return l, err
	}

	l.Line = uint(line)
	l.Column = uint(column)

	return l, nil
}
