package server

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"

	md "github.com/JohannesKaufmann/html-to-markdown"
	"github.com/bug-breeder/quip-mcp/pkg/quip"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// Server represents the MCP Quip server
type Server struct {
	mcpServer  *server.MCPServer
	quipClient *quip.Client
}

// New creates a new MCP Quip server
func New(token string) *Server {
	mcpServer := server.NewMCPServer(
		"Quip MCP Server",
		"1.4.0",
	)

	quipClient := quip.NewClient(token)

	s := &Server{
		mcpServer:  mcpServer,
		quipClient: quipClient,
	}

	// Register tools
	s.registerTools()
	// Register resources
	s.registerResources()

	return s
}

// Start starts the MCP server
func (s *Server) Start() error {
	log.Println("Starting MCP Quip Server...")
	return server.ServeStdio(s.mcpServer)
}

// registerTools registers all the MCP tools
func (s *Server) registerTools() {
	// Search documents tool
	searchTool := mcp.NewTool(
		"search_documents",
		mcp.WithDescription("Search for Quip documents"),
		mcp.WithString("query", mcp.Required(), mcp.Description("Search query for documents")),
		mcp.WithNumber("limit", mcp.Description("Maximum number of results (default: 10)")),
	)

	s.mcpServer.AddTool(searchTool, func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		query, err := req.RequireString("query")
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Invalid query argument: %v", err)), nil
		}

		limit := req.GetInt("limit", 10)

		result, err := s.quipClient.SearchDocuments(query, limit)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to search documents: %v", err)), nil
		}

		response := fmt.Sprintf("Found %d documents:\n\n", len(result.Documents))
		for i, doc := range result.Documents {
			response += fmt.Sprintf("%d. **%s**\n", i+1, doc.Title)
			response += fmt.Sprintf("   - ID: %s\n", doc.ID)
			response += fmt.Sprintf("   - Link: %s\n", doc.Link)
			response += fmt.Sprintf("   - Author: %s\n", doc.AuthorID)
			response += fmt.Sprintf("   - Updated: %s\n\n", formatTimestamp(doc.Updated))
		}

		return mcp.NewToolResultText(response), nil
	})

	// Get document tool
	getDocTool := mcp.NewTool(
		"get_document",
		mcp.WithDescription("Get a specific Quip document by ID"),
		mcp.WithString("document_id", mcp.Required(), mcp.Description("The ID of the document to retrieve")),
	)

	s.mcpServer.AddTool(getDocTool, func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		documentID, err := req.RequireString("document_id")
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Invalid document_id argument: %v", err)), nil
		}

		doc, err := s.quipClient.GetDocument(documentID)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to get document: %v", err)), nil
		}

		response := fmt.Sprintf("**%s**\n\n", doc.Title)
		response += fmt.Sprintf("- **ID:** %s\n", doc.ID)
		response += fmt.Sprintf("- **Type:** %s\n", doc.Type)
		response += fmt.Sprintf("- **Link:** %s\n", doc.Link)
		response += fmt.Sprintf("- **Author:** %s\n", doc.AuthorID)
		response += fmt.Sprintf("- **Created:** %s\n", formatTimestamp(doc.Created))
		response += fmt.Sprintf("- **Updated:** %s\n", formatTimestamp(doc.Updated))
		response += fmt.Sprintf("- **Access Level:** %s\n", doc.AccessLevel)

		if doc.HTML != "" {
			markdown := htmlToMarkdown(doc.HTML)
			response += fmt.Sprintf("\n**Content:**\n%s\n", markdown)
		}

		return mcp.NewToolResultText(response), nil
	})

	// Create document tool
	createDocTool := mcp.NewTool(
		"create_document",
		mcp.WithDescription("Create a new Quip document"),
		mcp.WithString("title", mcp.Required(), mcp.Description("The title of the new document")),
		mcp.WithString("content", mcp.Description("The initial content of the document (HTML format)")),
	)

	s.mcpServer.AddTool(createDocTool, func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		title, err := req.RequireString("title")
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Invalid title argument: %v", err)), nil
		}

		content := req.GetString("content", "")

		doc, err := s.quipClient.CreateDocument(title, content)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to create document: %v", err)), nil
		}

		response := "✅ **Document created successfully!**\n\n"
		response += fmt.Sprintf("- **Title:** %s\n", doc.Title)
		response += fmt.Sprintf("- **ID:** %s\n", doc.ID)
		response += fmt.Sprintf("- **Link:** %s\n", doc.Link)
		response += fmt.Sprintf("- **Created:** %s\n", formatTimestamp(doc.Created))

		return mcp.NewToolResultText(response), nil
	})

	// Get user tool
	getUserTool := mcp.NewTool(
		"get_user",
		mcp.WithDescription("Get Quip user information"),
		mcp.WithString("user_id", mcp.Required(), mcp.Description("The ID of the user to retrieve (use 'current' for current user)")),
	)

	s.mcpServer.AddTool(getUserTool, func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		userID, err := req.RequireString("user_id")
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Invalid user_id argument: %v", err)), nil
		}

		var user *quip.User

		if userID == "current" {
			user, err = s.quipClient.GetCurrentUser()
		} else {
			user, err = s.quipClient.GetUser(userID)
		}

		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to get user: %v", err)), nil
		}

		response := fmt.Sprintf("**%s**\n\n", user.Name)
		response += fmt.Sprintf("- **ID:** %s\n", user.ID)
		response += fmt.Sprintf("- **Email:** %s\n", user.Email)
		response += fmt.Sprintf("- **Profile URL:** %s\n", user.URL)
		response += fmt.Sprintf("- **Created:** %s\n", formatTimestamp(user.Created))
		response += fmt.Sprintf("- **Updated:** %s\n", formatTimestamp(user.Updated))

		if user.ProfilePic != "" {
			response += fmt.Sprintf("- **Profile Picture:** %s\n", user.ProfilePic)
		}

		return mcp.NewToolResultText(response), nil
	})

	// Get document comments tool
	getCommentsTool := mcp.NewTool(
		"get_document_comments",
		mcp.WithDescription("Get comments for a Quip document"),
		mcp.WithString("document_id", mcp.Required(), mcp.Description("The ID of the document to get comments for")),
	)

	s.mcpServer.AddTool(getCommentsTool, func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		documentID, err := req.RequireString("document_id")
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Invalid document_id argument: %v", err)), nil
		}

		comments, err := s.quipClient.GetDocumentComments(documentID)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to get comments: %v", err)), nil
		}

		if len(comments) == 0 {
			return mcp.NewToolResultText("No comments found for this document."), nil
		}

		response := fmt.Sprintf("Found %d comments:\n\n", len(comments))
		for i, comment := range comments {
			response += fmt.Sprintf("%d. **Author:** %s\n", i+1, comment.AuthorID)
			response += fmt.Sprintf("   **Created:** %s\n", formatTimestamp(comment.Created))
			response += fmt.Sprintf("   **Text:** %s\n\n", comment.Text)
		}

		return mcp.NewToolResultText(response), nil
	})

	// Edit document tool
	editDocTool := mcp.NewTool(
		"edit_document",
		mcp.WithDescription("Edit an existing Quip document"),
		mcp.WithString("document_id", mcp.Required(), mcp.Description("The ID of the document to edit")),
		mcp.WithString("content", mcp.Required(), mcp.Description("The new content for the document")),
		mcp.WithString("operation", mcp.Description("Edit operation: REPLACE (default), APPEND, PREPEND")),
		mcp.WithString("format", mcp.Description("Content format: html (default), markdown")),
	)

	s.mcpServer.AddTool(editDocTool, func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		documentID, err := req.RequireString("document_id")
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Invalid document_id argument: %v", err)), nil
		}

		content, err := req.RequireString("content")
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Invalid content argument: %v", err)), nil
		}

		operation := req.GetString("operation", "REPLACE")
		format := req.GetString("format", "html")

		doc, err := s.quipClient.EditDocument(documentID, content, operation, format)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to edit document: %v", err)), nil
		}

		response := "✅ **Document edited successfully!**\n\n"
		response += fmt.Sprintf("- **Title:** %s\n", doc.Title)
		response += fmt.Sprintf("- **ID:** %s\n", doc.ID)
		response += fmt.Sprintf("- **Link:** %s\n", doc.Link)
		response += fmt.Sprintf("- **Updated:** %s\n", formatTimestamp(doc.Updated))

		return mcp.NewToolResultText(response), nil
	})

	// Delete document tool
	deleteDocTool := mcp.NewTool(
		"delete_document",
		mcp.WithDescription("Delete a Quip document (requires confirmation)"),
		mcp.WithString("document_id", mcp.Required(), mcp.Description("The ID of the document to delete")),
		mcp.WithString("confirm", mcp.Required(), mcp.Description("Type 'DELETE' to confirm deletion")),
	)

	s.mcpServer.AddTool(deleteDocTool, func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		documentID, err := req.RequireString("document_id")
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Invalid document_id argument: %v", err)), nil
		}

		confirm, err := req.RequireString("confirm")
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Invalid confirm argument: %v", err)), nil
		}

		if confirm != "DELETE" {
			return mcp.NewToolResultError("Deletion cancelled. To delete the document, you must set confirm='DELETE'"), nil
		}

		// Get document info before deletion for confirmation
		doc, err := s.quipClient.GetDocument(documentID)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to get document before deletion: %v", err)), nil
		}

		err = s.quipClient.DeleteDocument(documentID)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to delete document: %v", err)), nil
		}

		response := "🗑️ **Document deleted successfully!**\n\n"
		response += fmt.Sprintf("- **Deleted Document:** %s\n", doc.Title)
		response += fmt.Sprintf("- **ID:** %s\n", doc.ID)
		response += "- **Status:** ✅ Permanently deleted\n"

		return mcp.NewToolResultText(response), nil
	})

	// Get recent threads tool
	getRecentTool := mcp.NewTool(
		"get_recent_threads",
		mcp.WithDescription("Get recent Quip threads for the current user"),
		mcp.WithNumber("limit", mcp.Description("Maximum number of recent threads to retrieve (default: 10)")),
	)

	s.mcpServer.AddTool(getRecentTool, func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		limit := req.GetInt("limit", 10)

		threads, err := s.quipClient.GetRecentThreads(limit)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to get recent threads: %v", err)), nil
		}

		if len(threads) == 0 {
			return mcp.NewToolResultText("No recent threads found."), nil
		}

		response := fmt.Sprintf("Found %d recent threads:\n\n", len(threads))
		for i, thread := range threads {
			response += fmt.Sprintf("%d. **%s**\n", i+1, thread.Title)
			response += fmt.Sprintf("   - ID: %s\n", thread.ID)
			response += fmt.Sprintf("   - Type: %s\n", thread.Type)
			response += fmt.Sprintf("   - Link: %s\n", thread.Link)
			response += fmt.Sprintf("   - Updated: %s\n\n", formatTimestamp(thread.Updated))
		}

		return mcp.NewToolResultText(response), nil
	})

	log.Println("✅ All MCP tools registered successfully")
}

// registerResources registers MCP resources
func (s *Server) registerResources() {
	// Current user resource
	currentUserResource := mcp.NewResource(
		"quip://user/current",
		"Current User",
		mcp.WithResourceDescription("Current Quip user information"),
		mcp.WithMIMEType("application/json"),
	)

	s.mcpServer.AddResource(currentUserResource, func(ctx context.Context, req mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
		user, err := s.quipClient.GetCurrentUser()
		if err != nil {
			return nil, fmt.Errorf("failed to get current user: %w", err)
		}

		jsonData := fmt.Sprintf(`{
  "id": "%s",
  "name": "%s",
  "email": "%s",
  "url": "%s",
  "created": %d,
  "updated": %d
}`, user.ID, user.Name, user.Email, user.URL, user.Created, user.Updated)

		return []mcp.ResourceContents{
			mcp.TextResourceContents{
				URI:      req.Params.URI,
				MIMEType: "application/json",
				Text:     jsonData,
			},
		}, nil
	})

	log.Println("✅ All MCP resources registered successfully")
}

// formatTimestamp converts a Unix timestamp (microseconds) to a readable format
func formatTimestamp(timestamp int64) string {
	if timestamp == 0 {
		return "Unknown"
	}
	// Convert microseconds to seconds
	seconds := timestamp / 1000000
	return strconv.FormatInt(seconds, 10)
}

// Helper function to truncate text
func truncateText(text string, maxLength int) string {
	if len(text) <= maxLength {
		return text
	}
	return strings.TrimSpace(text[:maxLength]) + "..."
}

// htmlToMarkdown converts HTML content to clean markdown
func htmlToMarkdown(htmlContent string) string {
	converter := md.NewConverter("", true, nil)

	markdown, err := converter.ConvertString(htmlContent)
	if err != nil {
		// If conversion fails, return a cleaned version of the HTML
		cleaned := strings.ReplaceAll(htmlContent, "<", "&lt;")
		cleaned = strings.ReplaceAll(cleaned, ">", "&gt;")
		return cleaned
	}

	// Clean up the markdown
	markdown = strings.ReplaceAll(markdown, "&amp;", "&")
	markdown = strings.ReplaceAll(markdown, "&lt;", "<")
	markdown = strings.ReplaceAll(markdown, "&gt;", ">")
	markdown = strings.ReplaceAll(markdown, "&nbsp;", " ")

	// Remove excessive newlines
	lines := strings.Split(markdown, "\n")
	var cleanLines []string
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed != "" || (len(cleanLines) > 0 && cleanLines[len(cleanLines)-1] != "") {
			cleanLines = append(cleanLines, trimmed)
		}
	}

	return strings.Join(cleanLines, "\n")
}
