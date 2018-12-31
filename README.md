# monkey-lang

[![Build Status](https://cloud.drone.io/api/badges/prologic/monkey-lang/status.svg)](https://cloud.drone.io/prologic/monkey-lang)
[![CodeCov](https://codecov.io/gh/prologic/monkey-lang/branch/master/graph/badge.svg)](https://codecov.io/gh/prologic/monkey-lang)
[![Go Report Card](https://goreportcard.com/badge/prologic/monkey-lang)](https://goreportcard.com/report/prologic/monkey-lang)
[![GoDoc](https://godoc.org/github.com/prologic/monkey-lang?status.svg)](https://godoc.org/github.com/prologic/monkey-lang) 
[![Sourcegraph](https://sourcegraph.com/github.com/prologic/monkey-lang/-/badge.svg)](https://sourcegraph.com/github.com/prologic/monkey-lang?badge)

Monkey programming language interpreter designed in [_Writing An Interpreter In Go_](https://interpreterbook.com/).
A step-by-step walk-through where each commit is a fully working part.
Read the book and follow along with the commit history.

## Usage

To build run `make`.

```#!bash
$ go get -u github.com/prologic/monkey-lang
$ cd $GOPATH/github.com/prologic/monkey-lang
$ make
This is the Monkey programming language!
Feel free to type in commands
>> 
```

To run the tests run `make test`

[Travis CI](https://travis-ci.org/) and [Codecov](https://codecov.io/)
are also integrated for Continuous Integration and Code Coverage reporting.

## Monkey Language

### Variable bindings and arithmetic expressions

```sh
>> let a = 10;
>> let b = a * 2;
>> (a + b) / 2 - 3;
12
>> let c = 2.5;
>> b + c
22.5
```

### If expressions

```sh
>> let a = 10;
>> let b = a * 2;
>> let c = if (b > a) { 99 } else { 100 };
>> c
99
```

### Functions and closures

```sh
>> let multiply = fn(x, y) { x * y };
>> multiply(50 / 2, 1 * 2)
50
>> fn(x) { x + 10 }(10)
20
>> let newAdder = fn(x) { fn(y) { x + y }; };
>> let addTwo = newAdder(2);
>> addTwo(3);
5
>> let sub = fn(a, b) { a - b };
>> let applyFunc = fn(a, b, func) { func(a, b) };
>> applyFunc(10, 2, sub);
8
```

### Strings

```sh
>> let makeGreeter = fn(greeting) { fn(name) { greeting + " " + name + "!" } };
>> let hello = makeGreeter("Hello");
>> hello("skatsuta");
Hello skatsuta!
```

### Arrays

```sh
>> let myArray = ["Thorsten", "Ball", 28, fn(x) { x * x }];
>> myArray[0]
Thorsten
>> myArray[4 - 2]
28
>> myArray[3](2);
4
```

### Hashes

```sh
>> let myHash = {"name": "Jimmy", "age": 72, true: "yes, a boolean", 99: "correct, an integer"};
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
>> len("hello");
5
>> len("∑");
3
>> let myArray = ["one", "two", "three"];
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

## License

This work is licensed under the terms of the MIT License.
