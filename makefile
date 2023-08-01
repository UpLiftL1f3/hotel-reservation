build:
	@go build -o bin/api

run: build
	@./bin/api

seed: 
	@go run ./scripts/seed.go

# Define the target to run tests in all directories
test:
	@echo "Running tests with logs..."
	@for dir in $$(find . -type d -not -path "./vendor/*" -not -path "./.git/*"); do \
		(cd "$$dir" && go test -v) 2>&1 | sed "s/^/    /"; \
	done