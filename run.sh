#! /bin/sh

lexer/lexer $1 | parser/parser | interpreter/interpreter
