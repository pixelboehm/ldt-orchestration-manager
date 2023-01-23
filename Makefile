build:
	go build -o out/main main.go

run:
	go run main.go

test:
	go test ./...

test-s:
	go test ./... --short

test-v:
	go test ./... -v --short

test-c:
	go test ./... -coverprofile=out/coverage.html

coverage:
	go tool cover -html=out/coverage.html

clean:
	go clean --cache

cover: 
	test-c 
	coverage

init:
	go mod init longevity

setup: 
	go mod vendor
	go mod tidy