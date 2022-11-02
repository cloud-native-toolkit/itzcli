.PHONY: default
.DEFAULT_GOAL := default

# The shell script is itz and the actual binary is itzcli
WRAPPER=itz
BINARY=itzcli
ITZ_VER := $(shell git describe --tags)
# Add windows here if/when we start supporting Windows OS officially
PLATFORMS=darwin linux
# Add 386 if we want, but for modern usages I see no reason why to include 32
# bit archs
ARCHITECTURES=amd64

LDFLAGS=-ldflags "-X main.Version=${ITZ_VER}"
ADDL_FILES=itz QUICKSTART.md

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

build:
	@echo "Building itzcli..."
	go build ${LDFLAGS} -o ${BINARY}

package:
	@tar cvf - itz itzcli | gzip > itzcli.tar.gz

build_all:
	$(foreach GOOS, $(PLATFORMS),\
	$(foreach GOARCH, $(ARCHITECTURES), $(shell export GOOS=$(GOOS); export GOARCH=$(GOARCH); mkdir -p bin/$(GOOS)/$(GOARCH) && go build -v $(LDFLAGS) -o bin/$(GOOS)/$(GOARCH)/$(BINARY))))

install:
	@go install ${LDFLAGS}

package_all:
	$(foreach GOOS, $(PLATFORMS),\
	$(foreach GOARCH, $(ARCHITECTURES), $(shell cp $(ADDL_FILES) bin/$(GOOS)/$(GOARCH))))

	$(foreach GOOS, $(PLATFORMS),\
	$(foreach GOARCH, $(ARCHITECTURES), $(shell export GOOS=$(GOOS); export GOARCH=$(GOARCH); tar -C bin/$(GOOS)/$(GOARCH) -cvf - $(ADDL_FILES) $(BINARY) | gzip > $(BINARY)-$(GOOS)-$(GOARCH).tar.gz)))

generate-docs:
	@rm -rf docs/*.md
	@go run docs/gendocs.go

ci: clean verify install build_all package_all
