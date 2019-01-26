let fib = fn(n, a, b) {
  if (n == 0) {
    return a
  }
  if (n == 1) {
    return b
  }
  return fib(n - 1, b, a + b)
}

puts(fib(35, 0, 1))
