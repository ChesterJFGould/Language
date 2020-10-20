#! /bin/sh

sed 's/\/\/.*$//g' $1 | lexer2/lexer2 | parser/parser | interpreter/interpreter
