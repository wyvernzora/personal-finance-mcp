# Personal Finance MCP

## Overview

Personal Finance MCP is an [Model Context Protocol (MCP)](https://modelcontextprotocol.io) server built
using [mcp-go](https://github.com/mark3labs/mcp-go) to provide LLMs ability to ingest personal finance data
such as spending transactions, asset holdings and such.

## Configuration
You can configure the server using the following environment variables

| Name               | Default   | Description                               |
| ------------------ | --------- | ----------------------------------------- |
| `BIND_ADDRESS`     | `0.0.0.0` | The IP address that the server listens on |
| `PORT`             | `3000`    | The port that the server listens on       |

## Usage
```
$ docker run -p 3000:3000 ghcr.io/wyvernzora/personal-finance-mcp:latest
```

## License
This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.
