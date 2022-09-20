.PHONY: default
.DEFAULT_GOAL := default

default: ci

generate-mocks:
	@bash scripts/generate-mocks.sh $(PWD)

clean:
	@echo "Cleaning up..."
	@rm -f atkcli atkcli.tar.gz

verify:
	@echo "Running tests..."
	# TODO: make sure at some point when we get tests that this fails the build...
	- go test ./tests/...

build:
	@echo "Building atkcli..."
	go build .

package:
	@tar cvf - atkcli | gzip > atkcli.tar.gz

ci: clean verify build package

