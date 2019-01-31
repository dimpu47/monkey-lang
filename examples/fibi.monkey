let fib = fn(x) {
  let a = 0
  let b = 1
  while (x > 0) {
    let a = b
    let b = a + b
    let x = x - 1
  }
  return a
}

print(fib(35))
