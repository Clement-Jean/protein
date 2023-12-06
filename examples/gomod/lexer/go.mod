module github.com/Clement-Jean/protein/examples/lexer

go 1.21.4

replace (
	github.com/Clement-Jean/protein/internal v0.0.0 => ../../../internal
	github.com/Clement-Jean/protein/lexer v0.0.0 => ../../../lexer
	github.com/Clement-Jean/protein/token v0.0.0 => ../../../token
)

require github.com/Clement-Jean/protein/lexer v0.0.0

require (
	github.com/Clement-Jean/protein/internal v0.0.0 // indirect
	github.com/Clement-Jean/protein/token v0.0.0 // indirect
)