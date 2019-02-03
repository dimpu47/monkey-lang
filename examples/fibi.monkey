let fib = fn(n) {
   if (n < 3) {
     return 1
   }
   let a = 1
   let b = 1
   let c = 0
   let i = 0
   while (i < n - 2) {
     c = a + b
     b = a
     a = c
     i = i + 1
   }
   return a
}

print(fib(35))
