#!./monkey-lang

test := fn(n) {
  if (n % 15 == 0) {
    return "FizzBuzz"
  } else if (n % 5 == 0) {
    return "Buzz"
  } else if (n % 3 == 0) {
    return "Fizz"
  } else {
    return str(n)
  }
}

n := 1
while (n <= 100) {
  print(test(n))
  n = n + 1
}
