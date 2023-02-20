package lexer

// Position is a position in the input
type Position struct {
	// Offset is the position relative to the beginning of the file (starts at 0)
	Offset int

	// Line is the file line (starts at 1)
	Line int

	// Column is the offset relative to the beginning of the line (starts at 0)
	Column int
}
