.PHONY: tidy lint

tidy:
	@for dir in $(shell find apps -mindepth 1 -maxdepth 1 -type d) $(shell find pkg -mindepth 1 -maxdepth 1 -type d); do \
  		echo "Syncing go.work..."; \
		go work sync; \
		echo "Formatting in $$dir..."; \
		( cd $$dir && go fmt ) || exit $$?; \
		echo "Lint fixing in $$dir..."; \
		( cd $$dir && golangci-lint run --fix ) || true; \
	done
	$(shell go env GOPATH)/bin/mockery

lint:
	@for dir in $(shell find apps -mindepth 1 -maxdepth 1 -type d) $(shell find pkg -mindepth 1 -maxdepth 1 -type d); do \
		echo "Running go vet in $$dir..."; \
		( cd $$dir && go vet ) || exit $$?; \
		echo "Running golangci-lint in $$dir..."; \
		(cd $$dir && golangci-lint run) || exit $$?; \
	done