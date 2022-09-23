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

install-config:
	@echo "Backing up existing config file..."
	@cp -n $(HOME)/.atk.yml $(HOME)/.atk.yml.bak
	@echo "Copying example config file to your home directory..."
	-cp -n docs/atk-example.yml $(HOME)/.atk.yml
	@echo "Done"
