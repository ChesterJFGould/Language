all: lexer/lexer parser/parser interpreter/interpreter

lexer/lexer: lexer/*.go
	cd lexer; go build

parser/parser: parser/*.go
	cd parser; go build

interpreter/interpreter: interpreter/*.go nodes/*.go location/*.go
	cd interpreter; go build

