.PHONY: build run inspector test cover

BIND_ADDRESS ?= 0.0.0.0
PORT ?= 3000

# Build the server
build:
	docker build -t ghcr.io/wyvernzora/personal-finance-mcp:dev .

# Run the server
run: build
	docker run -it --rm \
		-p $(PORT):$(PORT) \
		-e LUNCHMONEY_TOKEN \
		-e KUBERA_API_KEY \
		-e KUBERA_API_SECRET \
		-e KUBERA_PORTFOLIO_ID \
		-e BIND_ADDRESS=$(BIND_ADDRESS) \
		-e PORT=$(PORT) \
		ghcr.io/wyvernzora/personal-finance-mcp:dev

# Test
test:
	go test -coverprofile=coverage.out ./...

cover: test
	go tool cover -html=coverage.out -o coverage.html

# Run the MCP inspector
inspector:
	npx --yes '@modelcontextprotocol/inspector'
