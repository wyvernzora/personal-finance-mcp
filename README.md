# Personal Finance MCP

## Overview

Personal Finance MCP is an [Model Context Protocol (MCP)](https://modelcontextprotocol.io) server built
using [mcp-go](https://github.com/mark3labs/mcp-go) to provide LLMs ability to ingest personal finance data
such as spending transactions, asset holdings and such.

## Configuration

### Basic Configuration
Following are common configuration options for the server:

| Environment Variable | Default   | Description                               |
| -------------------- | --------- | ----------------------------------------- |
| `BIND_ADDRESS`       | `0.0.0.0` | The IP address that the server listens on |
| `PORT`               | `3000`    | The port that the server listens on       |

### Data Sources
Server supports the following data sources:

| Data Source                              | Description                                                         |
| ---------------------------------------- | ------------------------------------------------------------------- |
| [LunchMoney](pkg/datasource/lunch_money) | Use [LunchMoney](https://lunchmoney.app/) API to fetch transactions |
| [Kubera](pkg/datasource/kubera)          | Use [Kubera](https://www.kubera.com/) API to fetch assets and debts |

## Usage
```
$ docker run -p 3000:3000 ghcr.io/wyvernzora/personal-finance-mcp:latest
```

## License
This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.
