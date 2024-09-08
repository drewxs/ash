# Ash

A general-purpose, dynamically-typed, procedural programming language with some functional programming features.

## The language

Hello, World!:

```rs
print("Hello, World!")
```

Variables:

```rs
let some_num = 10;
let some_string = "Hello";
let some_bool = true;
let some_array = [1, 2, 3];
let some_hash = { 1: "one", 2: "two", 3: "three" };
let some_function = fn(x) { x + 1 };
```

Arrays:

```rs
let x = [1, 2, 3];
len(x); // 3
print(x) // [1, 2, 3]
first(x) // 1
last(x) // 3
rest(x) // [2, 3]
push(x, 4) // [1, 2, 3, 4]
```

Hashes:

```rs
let x = { 1: "one", "two": 1 + 1, 3: "three" };
x[1]; // one
x["two"]; // 2
let double = fn(x) { x * 2 };
set(x, "four", double(2)); // { 1: one, two: 2, 3: three, four: 4 }
```

Functions:

```rs
let x = fn(a, b) {
  a + b // implicit return
};
x(1, 2); // 3
```

Recursive functions:

```rs
let fib = fn(n) {
    if n == 0 { 0 }
    else {
        if n == 1 { 1 }
        else { fib(n - 1) + fib(n - 2) }
    }
};
let result = fib(15)
print(result)
```

## Installation

Build from source

```sh
git clone https://github.com/drewxs/ash.git
cd ash
make
```

The executable will be in bin/ash

To add to your path, run:

```sh
export PATH=$PATH:$(pwd)/bin
```

Or move it to `/usr/local/bin`:

```sh
mv bin/ash /usr/local/bin
```

### Neovim support

Register the filetype:

```lua
vim.filetype.add({
  extension = {
    ash = "ash",
  },
})
```

For syntax highlighting, add the following to your treesitter config:

```lua
vim.treesitter.language.register("rust", "ash")
```

Yes, we're using Rust for syntax highlighting (for now), it's the closest thing to Ash's syntax.

## Usage

Start the REPL:

```sh
ash
```

Run a file:

```sh
ash <filename>
```

Running an example:

```sh
ash examples/fib.ash
```

## Development

Prerequisites:

-   [Go](https://go.dev)
-   [GNU Make](https://www.gnu.org/software/make)

```sh
# Build
make

# Run tests
make test

# Run benchmarks
make bench

# Run the program
make run
```

## Design

```go
Lexer -> Parser -> Compiler -> Virtual Machine
|------ Compile Time -----|    |-- Run Time -|

String -> Tokens -> AST -> Bytecode -> Objects
```

-   Interpreter: Tree-walking
-   Compiler: Stack-based VM

## TODO

-   [ ] LSP
    -   [ ] Hover
    -   [ ] Completions
    -   [ ] Diagnostics
    -   [ ] Signature help
-   [ ] Treesitter grammar
-   [ ] Linter
-   [ ] Formatter
-   [ ] REPL
    -   [ ] History
    -   [ ] Autocomplete
    -   [ ] Commands
        -   [ ] `clear`
        -   [ ] `exit`
-   [ ] Standard library
-   [ ] Type system
    -   [ ] Inference
    -   [ ] Checking
-   [ ] Compiler/VM optimizations
    -   [ ] Tail call optimization
    -   [ ] Constant folding
    -   [ ] Dead code elimination

---

[License](https://github.com/drewxs/ash/blob/main/LICENSE)
