fact := fn(n, a) {
  if (n == 0) {
    return a
  }
  return fact(n - 1, a * n)
}

print(fact(5, 1))
