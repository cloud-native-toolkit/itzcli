.PHONY: usage

clean:
	@echo "Cleaning up..."
	@rm -f atkcli

build:
	go build .

