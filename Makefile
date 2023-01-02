build:
	go build -o out/main main.go

run:
	go run main.go

test:
	go test ./...

test-v:
	go test ./... -v

test-c:
	go test ./... -coverprofile=out/coverage.html

coverage:
	go tool cover -html=out/coverage.html


cover: test-c coverage