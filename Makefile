BINARY   := mediashuttle-mcp
CMD      := ./cmd/mediashuttle-mcp

# Install to /usr/local/bin if writable, else ~/go/bin
INSTALL_DIR := $(shell [ -w /usr/local/bin ] 2>/dev/null && echo "/usr/local/bin" \
            || echo "$${GOPATH:-$$HOME/go}/bin")

.PHONY: all build test clean install

all: build

build:
	go build -o $(BINARY) $(CMD)

test:
	go test ./...

clean:
	rm -f $(BINARY)

install: build
	@mkdir -p "$(INSTALL_DIR)"
	cp $(BINARY) "$(INSTALL_DIR)/$(BINARY)"
	@echo "Installed $(BINARY) to $(INSTALL_DIR)"
