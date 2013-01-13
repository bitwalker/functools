## functools

Functional programming tools for Go

## About

This package provides you with the tools you need to perform high level manipulation of functions that you might be used to coming from other programming languages. Of particular note is the lack of an Apply function, for partially applying function arguments - which happens to be the first 'tool' I've added.

## Current Features

#### Functions

**Apply**: Partial application of function arguments

**ApplyMulti**: Apply for functions with multiple return values

**Compose**: Takes functions `f` and `g`, and returns a new function `fg` whose signature is: `fg(x) = f(g(x))`

**ToList**: Creates a new LinkedList

**ToSlice**: Converts a LinkedList to a slice

**LinkedList**: A traditional linked list structure where each node of the list contains it's current value (Head) and a pointer to the next node in the list (Tail). This enables nifty things like infinite sequences, and lazy evaluation.

LinkedList currently supports the following methods:

```
// The length of the list
Length() int

// Maps a function to every element of the list
Map(func(x Anything) Anything) *LinkedList

// Reduces a list to a value by applying the 
// reducer to every element of the list.
Reduce(func(acc, x Anything) Anything) Anything 

// Take the first `x` elements of the list
Take(x int)

// Drop the first x elements of the list
Drop(x int)
```


## Contributing

If you have other tools you think belong in this package, by all means fork the repo and send a pull request. I'd like to make this a one stop shop for functional programming needs in Go.

## License

The MIT License (MIT)
Copyright (c) 2013 Paul Schoenfelder

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.