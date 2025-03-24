# MCP-TARMAQ

A Model Context Protocol (MCP) server that suggests files related to files that have already been modified.
It supports automatic application in Cline and Cursur of changes that should be done in a separate file.
It uses TARMAQ[^1], a change impact analysis method that extracts simultaneously modified relationships between files from the commit history.

[^1]: Thomas Rolfness, et al. Generalizing the Analysis of Evolutionary Coupling for Software Change Impact Analysis. In Proc. SANER 2016.

## Install

### deb
```bash
curl -o mcp-tarmaq.deb -L https://github.com/mazrean/mcp-tarmaq/releases/latest/download/mcp-tarmaq_amd64.deb
dpkg -i mcp-tarmaq.deb
```

### rpm
```bash
yum install https://github.com/mazrean/mcp-tarmaq/releases/latest/download/mcp-tarmaq_amd64.rpm
```

### Homebrew
```bash
brew install mazrean/tap/mcp-tarmaq
```

### Download prebuilt binaries
Download binary from [releases page](https://github.com/mazrean/mcp-tarmaq/releases/latest).

### go install
```bash
go install github.com/mazrean/mcp-tarmaq@latest
```

## Example config
```json
{
  "mcpServers": {
    "tarmaq": {
      "command": "mcp-tarmaq",
      "args": [ "-r", "<repository directory path>" ],
    }
  }
}
```
