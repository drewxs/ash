program :=ash

.PHONY = all
all: run

.PHONY = build
build:
	go build -o bin/$(program) .

.PHONY = run
run: build
	./bin/ash

.PHONY = clean
clean:
	rm -f ./bin

.PHONY = test
test:
	go test ./ast
	go test ./code
	go test ./compiler
	go test ./evaluator
	go test ./lexer
	go test ./object
	go test ./parser
	go test ./vm
