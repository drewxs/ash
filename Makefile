program := ash

.PHONY = all
all: run

.PHONY = build
build:
	@echo "Building..."
	@(cd src && go build -o bin/$(program) .)
	@echo "Build ready: src/bin/$(program)"

.PHONY = run
run: build
	@echo "Running src/bin/$(program)..."
	@src/bin/$(program)

.PHONY = clean
clean:
	@echo "Cleaning up..."
	@rm -rf src/bin
	@echo "Done."

.PHONY = test
test:
	@echo "Running tests..."
	@mod=$(filter-out $@,$(MAKECMDGOALS)); \
	if [ $$mod ]; then \
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

%:
	@:
