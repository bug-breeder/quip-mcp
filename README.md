# Quip MCP Server

A [Model Context Protocol (MCP)](https://modelcontextprotocol.io) server providing AI assistants with comprehensive Quip document access and management.

[![Install MCP Server](https://cursor.com/deeplink/mcp-install-light.svg)](https://cursor.com/install-mcp?url=https://raw.githubusercontent.com/bug-breeder/quip-mcp/main/.github/cursor-mcp.json)

## ✨ Features

- **Full Document Lifecycle**: Create, read, edit, delete Quip documents
- **Smart Search**: Find documents with comprehensive search capabilities  
- **Recent Documents**: Access your recently viewed/edited documents
- **User Management**: Get user information and details
- **Comments**: Retrieve and manage document discussions
- **Markdown Support**: Clean markdown formatting throughout
- **Robust API**: Handles complex Quip API response structures
- **Secure**: Token-based authentication with enterprise support

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

[![Install MCP Server](https://cursor.com/deeplink/mcp-install-light.svg)](https://cursor.com/install-mcp?url=https://raw.githubusercontent.com/bug-breeder/quip-mcp/main/.github/cursor-mcp.json)

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
   See the instructions below for your specific client.

That's it! Your AI assistant can now access your Quip documents.

## 🔌 MCP Client Configuration

### Cursor

1.  Click the "Add to Cursor" button at the top of this README.
2.  Alternatively, you can add the server manually in Cursor's settings. Go to `File > Settings > MCP` and add a new server with the following configuration:

    ```json
    {
      "mcpServers": {
        "quip-mcp": {
          "command": "quip-mcp"
        }
      }
    }
    ```

### Claude

You can add the Quip MCP server to Claude using two methods:

**1. Command Line**

Open your terminal and run the following command:

```bash
claude mcp add quip-mcp -- quip-mcp
```

**2. Configuration File**

Add the following directly into your claude desktop app's setting or to your `claude_desktop_config.json` file:

```json
{
  "mcp_servers": [
    {
      "name": "quip-mcp",
      "command": ["quip-mcp"]
    }
  ]
}
```

### Other Clients

For other MCP clients, you can typically add a new server in the settings. Use the following configuration:

```json
{
  "mcpServers": {
    "quip-mcp": {
      "command": "quip-mcp"
    }
  }
}
```

## 🔄 Updates

### Quick Update
Update to the latest version with one command:
```bash
curl -sSL https://raw.githubusercontent.com/bug-breeder/quip-mcp/main/install.sh | bash -s -- --update
```

The update script will:
- ✅ Check your current version
- ✅ Skip update if you already have the latest version
- ✅ Only install if a newer version is available
- ✅ Preserve your existing configuration

### Check Current Version
```bash
quip-mcp --version
```

## 🛠️ Available Tools

| Tool | Description |
|------|-------------|
| `get_recent_threads` | Get your recently viewed/edited documents |
| `search_documents` | Search for documents by keyword or query |
| `get_document` | Retrieve full document content by ID |
| `create_document` | Create new documents with markdown content |
| `edit_document` | Update existing documents (append/prepend/replace) |
| `delete_document` | Delete documents permanently |
| `get_user` | Get current user or specific user information |
| `get_document_comments` | Retrieve document comments and discussions |

## 📖 Usage Examples

### Get Recent Documents
```
Show my recent Quip documents
```

### Search Documents
```
Search for documents about "project planning"
```

### Read Document Content
```
Get the full content of document V9T5AFuROlBN
```

### Create New Document
```
Create a document titled "Meeting Notes" with markdown content about today's team meeting
```

### Edit Existing Document
```
Add a new section to document ABC123 about next week's goals
```

### Delete Document
```
Delete the test document XYZ789
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

### Build from Source
```bash
git clone https://github.com/bug-breeder/quip-mcp.git
cd quip-mcp
make build
```

### Quality Assurance Workflow
```bash
# Run comprehensive pre-commit checks (recommended before any commit)
make pre-commit

### Testing
```bash
# Unit tests (mocked)
make test-unit

# Integration tests (requires QUIP_API_TOKEN)
export QUIP_API_TOKEN="your-token"
make test-integration

# Run all tests
make test-all
```

### Development Tools
```bash
make help          # Show all available commands
make dev-setup     # Install development dependencies
make coverage      # Generate test coverage report
```

## 📄 License

MIT License - see [LICENSE](LICENSE) file for details.

## 🤝 Contributing

Contributions are welcome! Please feel free to submit issues and pull requests. 
