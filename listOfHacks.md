+ Bool literals are lexed as identifiers, then turned into BoolLiteral in the parser. Should be lexed as BoolLiteral.
+ Increment and Decrement operators don't work as I'm not sure which way is best to handle them yet.
+ Control flow keywords(i.e. break, continue) don't do anything.
+ That's a lot of type assertions and variable shadowing you got going on in the interpreter.
