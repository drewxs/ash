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
