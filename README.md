# Quip MCP Server

A [Model Context Protocol (MCP)](https://modelcontextprotocol.io) server that provides AI assistants with access to Quip documents and collaboration features.

[![Install MCP Server](https://cursor.com/deeplink/mcp-install-light.svg)](https://cursor.com/install-mcp?name=quip-mcp&config=eyJjb21tYW5kIjoicXVpcC1tY3AiLCJlbnYiOnsiUVVJUF9BUElfVE9LRU4iOiJZT1VSX0FQSV9UT0tFTiJ9fQo%3D)

## ✨ Features

- **Document Access**: Search, read, and create Quip documents
- **User Information**: Get current user details and other users
- **Comments**: Retrieve document comments and discussions
- **Clean Output**: Documents are converted to beautiful markdown format
- **Secure**: Token-based authentication with your Quip instance

## 🚀 Quick Install

### Step 1: Install the binary

#### One-line install (macOS/Linux)
```bash
curl -sSL https://raw.githubusercontent.com/bug-breeder/quip-mcp/main/install.sh | bash
```

#### Manual download
Download the appropriate binary for your platform from the [releases page](https://github.com/bug-breeder/quip-mcp/releases).

### Step 2: Add to Cursor (One-click)
After installing the binary, click the button below to add the MCP server configuration to your Cursor IDE:

[![Install MCP Server](https://cursor.com/deeplink/mcp-install-light.svg)](https://cursor.com/install-mcp?name=quip-mcp&config=eyJjb21tYW5kIjoicXVpcC1tY3AiLCJlbnYiOnsiUVVJUF9BUElfVE9LRU4iOiJZT1VSX0FQSV9UT0tFTiJ9fQo%3D)

> **Note**: This button only adds the MCP configuration to Cursor. You must install the `quip-mcp` binary first (Step 1).

## ⚡ Quick Start

1. **Get your API token**
   - Visit your Quip instance: `https://your-company.quip.com/dev/token`
   - Generate a new API token

2. **Configure the server**
   ```bash
   quip-mcp --setup
   ```

3. **Add to your MCP client**
   ```json
   {
     "mcpServers": {
       "quip": {
         "command": "quip-mcp"
       }
     }
   }
   ```

That's it! Your AI assistant can now access your Quip documents.

## 🛠️ Available Tools

| Tool | Description |
|------|-------------|
| `search_documents` | Search for documents by query |
| `get_document` | Retrieve a specific document by ID |
| `create_document` | Create a new document |
| `get_user` | Get user information |
| `get_document_comments` | Get comments for a document |

## 📖 Usage Examples

### Search for documents
```
Search for documents about "project planning"
```

### Read a specific document
```
Get the content of document V9T5AFuROlBN
```

### Create a new document
```
Create a document titled "Meeting Notes" with content about today's team meeting
```

## ⚙️ Configuration

### Environment Variable
```bash
export QUIP_API_TOKEN="your-token-here"
quip-mcp
```

### Configuration File
The server automatically saves your token to:
- Linux/macOS: `~/.config/quip-mcp/config.yaml`
- Windows: `%APPDATA%/quip-mcp/config.yaml`

### CLI Options
```bash
quip-mcp --help          # Show help
quip-mcp --version       # Show version
quip-mcp --setup         # Interactive token setup
quip-mcp --config        # Show current configuration
```

## 🏢 Company Instances

This server works with both Quip.com and company-specific instances:

- **Quip.com**: Standard public Quip
- **Company instances**: `https://your-company.quip.com`

The API token automatically handles routing to your specific instance.

## 🔒 Security

- Tokens are stored with restricted file permissions (0600)
- All API communication uses HTTPS
- No sensitive data is logged

## 🛡️ Troubleshooting

### Common Issues

**"No API token found"**
```bash
# Run interactive setup
quip-mcp --setup

# Or set environment variable
export QUIP_API_TOKEN="your-token-here"
```

**"Search not available"**
- Some Quip instances have search disabled
- Use direct document IDs instead
- Extract document ID from Quip URLs: `https://company.quip.com/DOCUMENT_ID/title`
**"Permission denied"**
- Ensure your API token has appropriate permissions
- Check document access levels in Quip

## 🔧 Development

### Build from source
```bash
git clone https://github.com/bug-breeder/quip-mcp.git
cd quip-mcp
go build -o quip-mcp .
```

### Run tests
```bash
make test
```

## 📄 License

MIT License - see [LICENSE](LICENSE) file for details.

## 🤝 Contributing

Contributions are welcome! Please feel free to submit issues and pull requests. 