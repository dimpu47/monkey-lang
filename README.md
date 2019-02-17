# monkey-lang

[![Build Status](https://cloud.drone.io/api/badges/prologic/monkey-lang/status.svg)](https://cloud.drone.io/prologic/monkey-lang)
[![CodeCov](https://codecov.io/gh/prologic/monkey-lang/branch/master/graph/badge.svg)](https://codecov.io/gh/prologic/monkey-lang)
[![Go Report Card](https://goreportcard.com/badge/prologic/monkey-lang)](https://goreportcard.com/report/prologic/monkey-lang)
[![GoDoc](https://godoc.org/github.com/prologic/monkey-lang?status.svg)](https://godoc.org/github.com/prologic/monkey-lang) 
[![Sourcegraph](https://sourcegraph.com/github.com/prologic/monkey-lang/-/badge.svg)](https://sourcegraph.com/github.com/prologic/monkey-lang?badge)

Monkey programming language interpreter designed in [_Writing An Interpreter In Go_](https://interpreterbook.com/).
A step-by-step walk-through where each commit is a fully working part.
Read the book and follow along with the commit history.

## Status

> Currently working on a [self-host](https://github.com/prologic/monkey-lang/tree/self-host)
> branch where I'm improving and modifying the original Monkey-lang implementation
> to support writing Monkey in itself (*ala self-hosting*). So far I've managed to write a
> [Brainfuck](https://github.com/prologic/monkey-lang/blob/self-host/examples/bf.monkey)
> interpreter in Monkey with only a small number of improvements.

## Read and Follow

> Read the books and follow along with the following commit history.
(*This also happens to be the elapsed days I took to read both books!*)

See: [Reading Guide](./ReadingGuide.md)

> Please note that whilst reading the awesome books I slihtly modified this
> version of Monkey-lang in some places. FOr example I opted to have a single
> `RETURN` Opcode.

## Quickstart

```#!sh
$ go get github.com/prologic/monkey-lang/cmd/monkey-lang
$ monkey-lang
```

## Development

To build run `make`.

```#!sh
$ go get github.com/prologic/monkey-lang
$ cd $GOPATH/github.com/prologic/monkey-lang
$ make
This is the Monkey programming language!
Feel free to type in commands
>> 
```

To run the tests run `make test`

You can also execute program files by invoking `monkey-lang <filename>`
There are also some command-line options:

```#!sh
$ ./monkey-lang -h
Usage: monkey-lang [options] [<filename>]
  -c	compile input to bytecode
  -d	enable debug mode
  -e string
    	engine to use (eval or vm) (default "vm")
  -i	enable interactive mode
  -v	display version information
```

## Monkey Language

> See also: [examples](./examples)

### Variable Bindings

```#!sh
>> a := 10
```

### Artithmetic Expressions

```#!sh
>> a := 10
>> b := a * 2
>> (a + b) / 2 - 3
12
```

### Conditional Expressions

```sh
>> a := 10
>> b := a * 2
>> c := if (b > a) { 99 } else { 100 }
>> c
99
```

### Functions and Closures

```sh
>> multiply := fn(x, y) { x * y }
>> multiply(50 / 2, 1 * 2)
50
>> fn(x) { x + 10 }(10)
20
>> newAdder := fn(x) { fn(y) { x + y } }
>> addTwo := newAdder(2)
>> addTwo(3)
5
>> sub := fn(a, b) { a - b }
>> applyFunc := fn(a, b, func) { func(a, b) }
>> applyFunc(10, 2, sub)
8
```

### Recursive Functions

```#!sh
>> wrapper := fn() { inner := fn(x) { if (x == 0) { return 2 } else { return inner(x - 1) } } return inner(1) }
>> wrapper()
2
```

### Strings

```sh
>> makeGreeter := fn(greeting) { fn(name) { greeting + " " + name + "!" } }
>> hello := makeGreeter("Hello")
>> hello("skatsuta")
Hello skatsuta!
```

### Arrays

```sh
>> myArray := ["Thorsten", "Ball", 28, fn(x) { x * x }]
>> myArray[0]
Thorsten
>> myArray[4 - 2]
28
>> myArray[3](2)
4
```

### Hashes

```sh
>> myHash := {"name": "Jimmy", "age": 72, true: "yes, a boolean", 99: "correct, an integer"}
>> myHash["name"]
Jimmy
>> myHash["age"]
72
>> myHash[true]
yes, a boolean
>> myHash[99]
correct, an integer
```

### Builtin functions

```sh
>> len("hello")
5
>> len("∑")
1
>> myArray := ["one", "two", "three"]
>> len(myArray)
3
>> first(myArray)
one
>> rest(myArray)
[two, three]
>> last(myArray)
three
>> push(myArray, "four")
[one, two, three, four]
>> puts("Hello World")
Hello World
nil
```

### Objects

```#!sh
>> Person := fn(name, age) { self := {} self.name = name self.age = age self.str = fn() { return self.name + ", aged " + str(self.age) } return self }
>> p := Person("John", 35)
>> p.str()
"John, aged 35"
```

## License

This work is licensed under the terms of the MIT License.
