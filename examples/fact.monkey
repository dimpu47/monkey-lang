fact := fn(n) {
  if (n == 0) {
    return 1
  }
  return n * fact(n - 1)
}

assert(fact(5) == 120, "fact(5) != 120")
