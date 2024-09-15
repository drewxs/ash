let fib = fn(n) {
    if n < 2 { return n }
    fib(n - 1) + fib(n - 2)
};
let result = fib(15)
print(result)
