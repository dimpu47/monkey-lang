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

### Programs

A Monkey program is simply zero or more statements. Statements don't actually
have to be separated by newlines, only by whitespace. The following is a valid
program (*but you'd probably use newlines in the`if` block in real life*):

```
s := "world"
print("Hello, " + s)
if (s != "") { t := "The end" print(t) }
// Hello, world
// The end
```

Between tokens, whitespace and comments
(*lines starting with `//` or `#` through to the end of a line*)
are ignored.

### Types

Monkey has the following data types: `null`, `bool`, `int`, `str`, `array`,
`hash`, and `fn`. The `int` type is a signed 64-bit integer, strings are
immutable arrays of bytes, arrays are growable arrays
(*use the `append()` builtin*), and hashes are unordered hash maps.
Trailing commas are **NOT** allowed after the last element in an array or hash:

Type      | Syntax                                    | Comments
--------- | ----------------------------------------- | -----------------------
null      | `null`                                    |
bool      | `true false`                              |
int       | `0 42 1234 -5`                            | `-5` is actually `5` with unary `-`
str       | `"" "foo" "\"quotes\" and a\nline break"` | Escapes: `\" \\ \t \r \n \t \xXX`
array     | `[] [1, 2] [1, 2, 3]`                     |
hash      | `{} {"a": 1} {"a": 1, "b": 2}`            |

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

Monkey supports `if` and `else`:

```sh
>> a := 10
>> b := a * 2
>> c := if (b > a) { 99 } else { 100 }
>> c
99
```

Monkey also supports `else if`:

```#!sh
>> test := fn(n) { if (n % 15 == 0) { return "FizzBuzz" } else if (n % 5 == 0) { return "Buzz" } else if (n % 3 == 0) { return "Fizz" } else { return str(n) } }
>> test(1)
"1"
>> test(3)
"Fizz"
>> test(5)
"Buzz"
>> test(15)
"FizzBuzz"
```

### While Loops

Monkey supports only one looping construct, the `while` loop:

```#!sh
i := 3
while (i > 0) {
    print(i)
    i = i - 1
}
// 3
// 2
// 1
```

Monkey does not have `break` or `continue`, but you can `return <value>` as
one way of breaking out of a loop early inside a function.

### Functions and Closures

You can define named or anonymous functions, including functions inside
functions that reference outer variables (*closures*).

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

**NOTE:** You cannot have a "bare return" -- it requires a return value.
          So if you don't want to return anything
          (*functions always return at least `null` anyway*),
          just say `return null`.

### Recursive Functions

Monkey also supports recursive functions including recursive functions defined
in the scope of another function (*self-recursion*).

```#!sh
>> wrapper := fn() { inner := fn(x) { if (x == 0) { return 2 } else { return inner(x - 1) } } return inner(1) }
>> wrapper()
2
```

Monkey also does tail
call optimization and turns recursive tail-calls into iteration.

```#!sh
>> fib := fn(n, a, b) { if (n == 0) { return a } if (n == 1) { return b } return fib(n - 1, b, a + b) }
>> fib(35, 0, 1)
9227465
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

### Assignment Expressions

Assignment can assign to a name, an array element by index, or a hash value by key.
When assigning to a name (variable), it always assigns to the scope the variable was defined .

To help with object-oriented programming, `obj.foo = bar` is syntactic sugar for `obj["foo"] = bar`. They're exactly equivalent.

```
i := 1
func mutate() {
    i = 2
    print(i)
}
print(i)
mutate()
print(i)
// 1
// 2
// 2

map = {"a": 1}
func mutate() {
    map.a = 2
    print(map.a)
}
print(map.a)
mutate()
print(map.a)
// 1
// 2
// 2

lst := [0, 1, 2]
lst[1] = "one"
print(lst)
// [0, "one", 2]

map = {"a": 1, "b": 2}
map["a"] = 3
map.c = 4
print(map)
// {"a": 3, "b": 2, "c": 4}
```

### Binary and unary operators

Monkey supports pretty standard binary and unary operators.
Here they are with their precedence, from highest to lowest
(*operators of the same precedence evaluate left to right*):

Operators      | Description
-------------- | -----------
`[] obj.keu`   | Subscript
`-`            | Unary minus
`* / %`        | Multiplication, Division, Modulo
`+ -`          | Addition, Subtraction
`< <= > >= in` | Comparison
`== !=`        | Equality
`!`            | Bitwise not / Logical not
`&`            | Bitwise and
<code>&#124;</code> | Bitwise or

Several of the operators are overloaded. Here are the types they can operate on:

Operator   | Types           | Action
---------- | --------------- | ------
`[]`       | `str[int]`      | fetch nth byte of str (0-based)
`[]`       | `array[int]`    | fetch nth element of array (0-based)
`[]`       | `hash[str]`     | fetch hash value by key str
`-`        | `int`           | negate int
`*`        | `int * int`     | multiply ints
`*`        | `str * int`     | repeat str n times (**TBD**)
`*`        | `int * str`     | repeat str n times (**TBD**)
`*`        | `array * int`   | repeat array n times, give new array (**TBD**)
`*`        | `int * array`   | repeat array n times, give new array (**TBD**)
`/`        | `int / int`     | divide ints, truncated
`%`        | `int % int`     | divide ints, give remainder
`+`        | `int + int`     | add ints
`+`        | `str + str`     | concatenate strs, give new string
`+`        | `array + array` | concatenate arrays, give new array (**TBD**)
`+`        | `hash + hash`   | merge hashes into new hash, keys in right hash win (**TBD**)
`-`        | `int - int`     | subtract ints
`<`        | `int < int`     | true iff left < right
`<`        | `str < str`     | true iff left < right (lexicographical)
`<`        | `array < array` | true iff left < right (lexicographical, recursive)
`<= > >=`  | same as `<`     | similar to `<`
`in`       | `str in str`    | true iff left is substr of right
`in`       | `any in array`  | true iff one of array elements == left (**TBD**)
`in`       | `str in hash`   | true iff key in hash (**TBD**)
`==`       | `any == any`    | deep equality (always false if different type)
`!=`       | `any != any`    | same as `not ==`
`!`        | `not bool`      | inverse of bool
`&&`       | `bool and bool` | true iff both true, right not evaluated if left false
<code>&#124;&#124;</code> | `bool or bool`  | true iff either true, right not evaluated if left true

### Builtin functions

- `len(iterable)`
  Returns the length of the iterable (`str`, `array` or `hash`).
- `input([prompt])`
  Reads a line from standard input optionally printing `prompt`.
- `print(value...)`
  Prints the `value`(s) to standard output followed by a newline.
- `first(array)`
  Returns the first element of the `array`.
- `last(array)`
  Returns the last element of the `array`.
- `rest(array)`
  Returns a new array with the first element of `array` removed.
- `push(array, value)`
  Returns a new array with `value` pushed onto the end of `array`.
- `pop(array)`
  Returns the last value of the `array` or `null` if empty.
- `exit([status])`
  Exits the program immediately with the optional `status` or `0`.
- `assert(expr, [msg])`
  Exits the program immediately with a non-zero status if `expr` is `false`
  optionally displaying `msg` to standard error.
- `bool(value)`
  Converts `value` to a `bool`. If `value` is `bool` returns the value directly.
  Returns `true` for non-zero `int`(s), `false` otherwise. Returns `true` for
  non-empty `str`, `array` and `hash` values. Returns `true` for all other
  values except `null` which always returns `false`.
- `int(value)`
  Converts decimal `value` `str` to `int`. If `value` is invalid returns `null.
  If `value` is an `int` returns its value directly.
- `str(value)`
  Returns the string representation of `value`: `null` for null,
  `true` or `false` for `bool`, decimal for `int` (eg: `1234`),
  the string itself for `str` (not quoted),
  the Monkey representation for array and hash (eg: `[1, 2]` and `{"a": 1}`
  with keys sorted), and something like `<fn name(...) at 0x...>` for functions..

Coming soon... 

- `args()`
  Returns an array of command-line options passed to the program.
- `find(haystack, needle)
  Returns the index of `needle` `str` in `haystack` `str`,
  or the index of `needle` element in `haystack` array.
  Returns -1 if not found. (**TBD**)
- `join(list, sep)` concatenates strs in list to form a single str, with the separator str between each element.
- `lower(str)` returns a lowercased version of str.
- `read([filename])` reads standard input or the given file and returns the contents as a str.
- `sort(list[, func])` sorts the list in place using a stable sort, and returns nil. Elements in the list must be orderable with `<` (int, str, or list of those). If a key function is provided, it must take the element as an argument and return an orderable value to use as the sort key.
- `split(str[, sep])` splits the str using given separator, and returns the parts (excluding the separator) as a list. If sep is not given or nil, it splits on whitespace.
- `type(value)` returns a str denoting the type of value: `nil`, `bool`, `int`, `str`, `list`, `map`, or `func`.
- `upper(str)` returns an uppercased version of str.

### Objects

```#!sh
>> Person := fn(name, age) { self := {} self.name = name self.age = age self.str = fn() { return self.name + ", aged " + str(self.age) } return self }
>> p := Person("John", 35)
>> p.str()
"John, aged 35"
```

## License

This work is licensed under the terms of the MIT License.
