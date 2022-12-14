.PHONY: default
.DEFAULT_GOAL := default

# The shell script is itz and the actual binary is itzcli
WRAPPER=itz
BINARY=itzcli
ITZ_VER := $(shell git describe --tags)
# Add windows here if/when we start supporting Windows OS officially
PLATFORMS=darwin linux windows
# Add 386 if we want, but for modern usages I see no reason why to include 32
# bit archs
ARCHITECTURES=amd64

LDFLAGS=-ldflags "-X main.Version=${ITZ_VER}"
ADDL_FILES=itz QUICKSTART.md CHANGELOG.md

default: ci

clean-mocks:
	@rm -rf ./mocks

generate-mocks:
	@bash scripts/generate-mocks.sh $(PWD)

regenerate-mocks: clean-mocks generate-mocks

clean:
	@echo "Cleaning up..."
	@rm -rf bin
	@rm -rf $(BINARY)-*.tar.gz

verify: regenerate-mocks
	@echo "Running tests..."
	go test ./test/...
	test/scripts/cli-tests.sh

build:
	@echo "Building itzcli..."
	go build ${LDFLAGS} -o ${BINARY}

package:
	@tar cvf - itz itzcli | gzip > itzcli.tar.gz

build_all:
	$(foreach GOOS, $(PLATFORMS),\
	$(foreach GOARCH, $(ARCHITECTURES), $(shell mkdir -p bin/$(GOOS)/$(GOARCH) && go build -v $(LDFLAGS) -o bin/$(GOOS)/$(GOARCH)/$(BINARY))))

install:
	@go install ${LDFLAGS}

package_all:
	$(foreach GOOS, $(PLATFORMS),\
	$(foreach GOARCH, $(ARCHITECTURES), $(shell cp $(ADDL_FILES) bin/$(GOOS)/$(GOARCH))))

	$(foreach GOOS, $(PLATFORMS),\
	$(foreach GOARCH, $(ARCHITECTURES), $(shell tar -C bin/$(GOOS)/$(GOARCH) -cvf - $(ADDL_FILES) $(BINARY) | gzip > $(BINARY)-$(GOOS)-$(GOARCH).tar.gz)))

# This is not as dynamic as the others, but it is just used for the one-off of creating a ZIP file for Windows users
# who might be more accustomed to using ZIP rather than TAR files.
	$(foreach GOOS, windows, $(foreach GOARCH, amd64, @zip -j -r $(BINARY)-$(GOOS)-$(GOARCH).zip bin/$(GOOS)/$(GOARCH)))

generate-docs:
	@rm -rf docs/*.md
	@go run docs/gendocs.go

ci: clean build verify install build_all package_all
