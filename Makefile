PROJECT_NAME:=longevity

build:
	go build -o out/$(PROJECT_NAME) main.go

run:
	go run main.go

.PHONY: test
test:
	go test \
		$(if $(TEST_VERBOSE),-v) \
		$(if $(TEST_SHORT),--short) \
		./...

.PHONY: test-c
test-c:
	go test ./... -coverprofile=out/coverage.html

.PHONY: coverage
coverage:
	go tool cover -html=out/coverage.html

.PHONY: clean
clean:
	go clean --cache

cover: 
	test-c 
	coverage

init:
	go mod init $(PROJECT_NAME)

setup: 
	go mod tidy
	go mod vendor
