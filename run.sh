#! /bin/sh

sed 's/\/\/.*$//g' $1 | lexer/lexer | parser/parser | interpreter/interpreter
