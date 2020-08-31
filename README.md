# Description
A language I built for fun and to learn things.
It is easily extensible as each piece of the language(i.e. lexer, parser etc.) is completely self contained and built to communicate with the next piece by text stream via unix pipe.
This means that building something like a semantic analyzer is as simple as writing a program that takes in an ast as a text stream, checking the semantics of the ast, and then printing the ast to stdout.

# Inspirations
+ The parser is based on the paper "Top Down Operator Precedence" by Vaughan R. Pratt. Go read it [here](https://tdop.github.io/), it makes parsing almost trivial.
+ The book "Writing An Interpreter In Go" by Thorsten Ball pushed me to start this project when I read the first few pages and then realized I kind of knew where he was going with it.

# Requirements
+ POSIX shell
+ Go language compiler

# Example
The following prints "Hello, World!", an approximation of pi, then counts down from ten.

    {
        println("Hello,"+" World!")
      
        var pi = 1.0
      
        for var i = 1.0; i < 1000000.0; var i = i + 2.0 {
            var pi = pi - 1.0 / (1.0 + i * 2.0) + 1.0 / (1.0 + (i + 1.0) * 2.0)
        }
        
        println("pi is almost "+string(pi*4.0))
        
        for var i = 10; i > 0; var i = i - 1 println(string(i))
    }
    
# Features I plan to add
Due to time constraints on getting a minimal viable product working I chose to omit some essential features.
The following are the features I plan to add in order of importance.
+ Non-builtin functions and arrays.
+ First class functions.
+ Better syntax for things? All syntax is subject to change. 
+ LLVM frontend for compilation, jit etc.
