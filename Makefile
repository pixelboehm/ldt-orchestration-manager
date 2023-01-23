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

test-c:
	go test ./... -coverprofile=out/coverage.html

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
	go mod vendor
	go mod tidy