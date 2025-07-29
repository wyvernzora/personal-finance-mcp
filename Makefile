.PHONY: build run inspector

BIND_ADDRESS ?= 0.0.0.0
PORT ?= 3000

# Build the server
build:
	docker build -t ghcr.io/wyvernzora/personal-finance-mcp:dev .

# Run the server
run: build
	docker run -it --rm \
		-p $(PORT):$(PORT) \
		-e LUNCHMONEY_TOKEN=$(LUNCHMONEY_TOKEN) \
		-e BIND_ADDRESS=$(BIND_ADDRESS) \
		-e PORT=$(PORT) \
		ghcr.io/wyvernzora/personal-finance-mcp:dev

# Run the MCP inspector
inspector:
	npx --yes '@modelcontextprotocol/inspector'
