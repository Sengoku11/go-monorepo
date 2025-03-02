.PHONY: tidy lint test

tidy:
	@for dir in $(shell find apps -mindepth 1 -maxdepth 1 -type d) $(shell find pkg -mindepth 1 -maxdepth 1 -type d); do \
  		echo "Syncing go.work..."; \
		go work sync; \
		echo "Formatting in $$dir..."; \
		( cd $$dir && go fmt ) || exit $$?; \
		echo "Lint fixing in $$dir..."; \
		( cd $$dir && golangci-lint run --fix ) || true; \
		echo "Generating code in $$dir..."; \
		( cd $$dir && PATH=$(shell go env GOPATH)/bin:$$PATH go generate ./... ) || true; \
	done
	$(shell go env GOPATH)/bin/mockery

lint:
	@for dir in $(shell find apps -mindepth 1 -maxdepth 1 -type d) $(shell find pkg -mindepth 1 -maxdepth 1 -type d); do \
		echo "Running go vet in $$dir..."; \
		( cd $$dir && go vet ) || exit $$?; \
		echo "Running golangci-lint in $$dir..."; \
		(cd $$dir && golangci-lint run) || exit $$?; \
		echo "Run tests in $$dir..."; \
		( cd $$dir && go test -race ./...) || exit; \
	done

test:
	@for dir in $(shell find apps -mindepth 1 -maxdepth 1 -type d) $(shell find pkg -mindepth 1 -maxdepth 1 -type d); do \
		echo "Run tests in $$dir..."; \
		( cd $$dir && go test -race ./...) || exit; \
	done