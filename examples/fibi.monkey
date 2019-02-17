fib := fn(n) {
   if (n < 3) {
     return 1
   }
   a := 1
   b := 1
   c := 0
   i := 0
   while (i < n - 2) {
     c = a + b
     b = a
     a = c
     i = i + 1
   }
   return a
}

print(fib(35))
