.PHONY: default
.DEFAULT_GOAL := default

default: ci

clean-mocks:
	@rm -rf ./mocks

generate-mocks:
	@bash scripts/generate-mocks.sh $(PWD)

regenerate-mocks: clean-mocks generate-mocks

clean:
	@echo "Cleaning up..."
	@rm -f atkcli atkcli.tar.gz

verify: regenerate-mocks
	@echo "Running tests..."
	- go test ./test/...

build:
	@echo "Building atkcli..."
	go build .

package:
	@tar cvf - atkcli | gzip > atkcli.tar.gz

ci: clean verify build package

