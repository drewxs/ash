program := ash

.PHONY = all
all: run

.PHONY = build
build:
	@echo "Building..."
	@(cd src && go build -o ../bin/$(program) .)
	@echo "Build ready: bin/$(program)"

.PHONY = run
run: build
	@echo "Running bin/$(program)..."
	@bin/$(program)

.PHONY = clean
clean:
	@echo "Cleaning up..."
	@rm -rf bin
	@echo "Done."

.PHONY = test
test:
	@echo "Running tests..."
	@mod=$(filter-out $@,$(MAKECMDGOALS)); \
	if echo $$mod | grep -q ":"; then \
		mod_name=$$(echo $$mod | cut -d ':' -f 1); \
		test_name=$$(echo $$mod | cut -d ':' -f 2); \
		cd src && go test ./$$mod_name -run $$test_name; \
	elif [ $$mod ]; then \
		cd src && go test ./$$mod; \
	else \
		(cd src && go test ./ast) && \
		(cd src && go test ./code) && \
		(cd src && go test ./compiler) && \
		(cd src && go test ./evaluator) && \
		(cd src && go test ./lexer) && \
		(cd src && go test ./object) && \
		(cd src && go test ./parser) && \
		(cd src && go test ./vm); \
	fi
	@echo "Done."

.PHONY = bench
bench:
	@echo "Running benchmarks..."
	@(cd src && go build -o ../bin/bench ./benchmark)
	@engine=$(filter-out $@,$(MAKECMDGOALS)); \
	if [ $$engine ]; then \
		bin/bench -engine=$$engine; \
	else \
		bin/bench; \
	fi

%:
	@:
