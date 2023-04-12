PROJECT_NAME	:= orchestration-manager
BIN_NAME		:= $(PROJECT_NAME)

GO_BIN 			:= go
GORELEASER_BIN 	:= goreleaser

PROJECT_DIR 		:= $(CURDIR)
PROJECT_BUILD_DIR 	:= $(PROJECT_DIR)/out

.PHONY: build
build:
	@printf "%s\n" 'Compiling...'
	@$(GO_BIN) build -o $(PROJECT_BUILD_DIR)/$(BIN_NAME) main.go
	@printf "%s\n" 'Done'

.PHONY: cli
cli:
	@printf "%s\n" 'Compiling CLI...'
	@$(GO_BIN) build -o $(PROJECT_BUILD_DIR)/odm src/cli/cli.go
	@printf "%s\n" 'Done'


run: build
	@printf "Executing %s\n\n" "$(PROJECT_BUILD_DIR)/$(BIN_NAME)"
	@$(PROJECT_BUILD_DIR)/$(BIN_NAME) $(if $(RUN_ARGS), $(RUN_ARGS))

.PHONY: test
test:
	@$(GO_BIN) test \
		$(if $(TEST_VERBOSE),-v) \
		$(if $(TEST_SHORT),--short) \
		./...

.PHONY: test-c
test-c:
	@$(GO_BIN) test ./... -coverprofile=$(PROJECT_BUILD_DIR)/coverage.html

.PHONY: coverage
coverage:
	@$(GO_BIN) tool cover -html=$(PROJECT_BUILD_DIR)/coverage.html

.PHONY: clean
clean:
	@$(GO_BIN) clean --cache
	rm -rf $(PROJECT_BUILD_DIR)
	rm -rf $(PROJECT_DIR)/dist

cover: test-c coverage

init:
	@$(GO_BIN) mod init $(PROJECT_NAME)

setup: 
	@$(GO_BIN) mod tidy
	@$(GO_BIN) mod vendor

release:
	@$(GORELEASER_BIN) release --clean

releaseLocal:
	@$(GORELEASER_BIN) release --clean --snapshot