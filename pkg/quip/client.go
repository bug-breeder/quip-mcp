package quip

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	BaseURL = "https://platform.quip.com/1"
	Timeout = 30 * time.Second
)

// Client represents a Quip API client
type Client struct {
	token      string
	baseURL    string
	httpClient *http.Client
}

// Document represents a Quip document
type Document struct {
	ID              string                 `json:"id"`
	Type            string                 `json:"type"`
	Title           string                 `json:"title"`
	Created         int64                  `json:"created_usec"`
	Updated         int64                  `json:"updated_usec"`
	AuthorID        string                 `json:"author_id"`
	HTML            string                 `json:"html,omitempty"`
	Link            string                 `json:"link"`
	AccessLevel     string                 `json:"access_level"`
	IsTemplate      bool                   `json:"is_template"`
	SharedFolderID  string                 `json:"shared_folder_id,omitempty"`
	ThreadID        string                 `json:"thread_id"`
	UserIsFollowing bool                   `json:"user_is_following"`
	ExpandedUserIds []string               `json:"expanded_user_ids,omitempty"`
	AccessLevels    map[string]interface{} `json:"access_levels,omitempty"`
}

// User represents a Quip user
type User struct {
	ID         string   `json:"id"`
	Name       string   `json:"name"`
	Email      string   `json:"email"`
	Affinity   float64  `json:"affinity"`
	Desktop    bool     `json:"desktop"`
	Created    int64    `json:"created"`
	Updated    int64    `json:"updated"`
	URL        string   `json:"url"`
	ProfilePic string   `json:"profile_picture_url"`
	Emails     []string `json:"emails,omitempty"`
	ChatOnly   bool     `json:"chat_only"`
}

// SearchResult represents search results from Quip
type SearchResult struct {
	Documents  []Document `json:"documents"`
	Users      []User     `json:"users"`
	NextCursor string     `json:"next_cursor,omitempty"`
}

// SearchResponse represents the actual API response structure for search
type SearchResponse struct {
	Thread Document `json:"thread"`
}

// RecentThreadsResponse represents the API response structure for recent threads
// The API returns a map where keys are document IDs and values contain thread data
type RecentThreadsResponse map[string]RecentThreadData

type RecentThreadData struct {
	Thread            Document                     `json:"thread"`
	UserIds           []string                     `json:"user_ids,omitempty"`
	SharedFolderIds   []string                     `json:"shared_folder_ids,omitempty"`
	ExpandedUserIds   []string                     `json:"expanded_user_ids,omitempty"`
	InvitedUserEmails []string                     `json:"invited_user_emails,omitempty"`
	AccessLevels      map[string]map[string]string `json:"access_levels,omitempty"`
	HTML              string                       `json:"html,omitempty"`
}

// Comment represents a document comment
type Comment struct {
	ID       string `json:"id"`
	Text     string `json:"text"`
	AuthorID string `json:"author_id"`
	Created  int64  `json:"created_usec"`
	Updated  int64  `json:"updated_usec"`
	ParentID string `json:"parent_id,omitempty"`
	Visible  bool   `json:"visible"`
}

// NewClient creates a new Quip API client
func NewClient(token string) *Client {
	return &Client{
		token:   token,
		baseURL: BaseURL,
		httpClient: &http.Client{
			Timeout: Timeout,
		},
	}
}

// makeRequest performs an HTTP request to the Quip API with JSON body
func (c *Client) makeRequest(method, endpoint string, body interface{}) (*http.Response, error) {
	var reqBody io.Reader
	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		reqBody = bytes.NewBuffer(jsonData)
	}

	url := fmt.Sprintf("%s%s", c.baseURL, endpoint)
	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.token)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "MCP-Quip-Server/1.0")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}

	if resp.StatusCode >= 400 {
		defer resp.Body.Close()
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API error %d: %s", resp.StatusCode, string(bodyBytes))
	}

	return resp, nil
}

// makeFormRequest performs an HTTP request to the Quip API with form-urlencoded body
func (c *Client) makeFormRequest(method, endpoint string, formData map[string]string) (*http.Response, error) {
	var reqBody io.Reader
	if formData != nil {
		values := url.Values{}
		for key, value := range formData {
			values.Set(key, value)
		}
		reqBody = strings.NewReader(values.Encode())
	}

	url := fmt.Sprintf("%s%s", c.baseURL, endpoint)
	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.token)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", "MCP-Quip-Server/1.0")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}

	if resp.StatusCode >= 400 {
		defer resp.Body.Close()
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API error %d: %s", resp.StatusCode, string(bodyBytes))
	}

	return resp, nil
}

// GetCurrentUser returns information about the current user
func (c *Client) GetCurrentUser() (*User, error) {
	resp, err := c.makeRequest("GET", "/users/current", nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var user User
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &user, nil
}

// SearchDocuments searches for documents
func (c *Client) SearchDocuments(query string, limit int) (*SearchResult, error) {
	endpoint := fmt.Sprintf("/threads/search?query=%s", url.QueryEscape(query))
	if limit > 0 {
		endpoint += fmt.Sprintf("&count=%d", limit)
	}

	resp, err := c.makeRequest("GET", endpoint, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// The API returns an array of objects with "thread" property
	var apiResponse []SearchResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResponse); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	// Transform the response to our expected format
	result := &SearchResult{
		Documents: make([]Document, len(apiResponse)),
		Users:     []User{}, // Search API doesn't return users
	}

	for i, item := range apiResponse {
		result.Documents[i] = item.Thread
	}

	return result, nil
}

// GetDocument retrieves a document by ID using v1 API and includes HTML content
func (c *Client) GetDocument(id string) (*Document, error) {
	// Use v1 API to get document with HTML content
	endpoint := fmt.Sprintf("/threads/%s", id)

	resp, err := c.makeRequest("GET", endpoint, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Read the response body to check the structure
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Try to decode as the complex structure first (like CreateDocument and GetRecentThreads)
	var response RecentThreadData
	if err := json.Unmarshal(respBody, &response); err == nil && response.Thread.ID != "" {
		// The HTML content is in the response.HTML field, not response.Thread.HTML
		if response.HTML != "" {
			response.Thread.HTML = response.HTML
		}
		return &response.Thread, nil
	}

	// Fallback to direct document structure
	var doc Document
	if err := json.Unmarshal(respBody, &doc); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &doc, nil
}

// CreateDocument creates a new document
func (c *Client) CreateDocument(title, content string) (*Document, error) {
	formData := map[string]string{
		"title":   title,
		"content": content,
		"format":  "markdown",
	}

	resp, err := c.makeFormRequest("POST", "/threads/new-document", formData)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Read the response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// The CreateDocument API returns the same structure as RecentThreads - an object with a "thread" property
	var response RecentThreadData
	if err := json.Unmarshal(respBody, &response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &response.Thread, nil
}

// GetDocumentComments retrieves comments for a document
func (c *Client) GetDocumentComments(documentID string) ([]Comment, error) {
	endpoint := fmt.Sprintf("/threads/%s/messages", documentID)

	resp, err := c.makeRequest("GET", endpoint, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var comments []Comment
	if err := json.NewDecoder(resp.Body).Decode(&comments); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return comments, nil
}

// EditDocument edits an existing document
func (c *Client) EditDocument(documentID, content, operation, format string) (*Document, error) {
	formData := map[string]string{
		"thread_id": documentID,
		"content":   content,
	}

	// Convert operation to location parameter as per Quip API v1
	// For now, map "REPLACE" to location=0 (APPEND)
	// In a full implementation, we'd need section_id for true replacement
	if operation == "REPLACE" || operation == "" {
		formData["location"] = "0" // APPEND - adds to end of document
	} else if operation == "APPEND" {
		formData["location"] = "0" // APPEND
	} else if operation == "PREPEND" {
		formData["location"] = "1" // PREPEND
	} else {
		formData["location"] = "0" // Default to APPEND
	}

	if format != "" {
		formData["format"] = format
	} else {
		formData["format"] = "markdown"
	}

	endpoint := "/threads/edit-document"
	resp, err := c.makeFormRequest("POST", endpoint, formData)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// The edit-document API returns the same complex structure as other APIs
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Try to decode as the complex structure first
	var response RecentThreadData
	if err := json.Unmarshal(respBody, &response); err == nil && response.Thread.ID != "" {
		return &response.Thread, nil
	}

	// Fallback to direct document structure
	var doc Document
	if err := json.Unmarshal(respBody, &doc); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &doc, nil
}

// DeleteDocument deletes a document
func (c *Client) DeleteDocument(documentID string) error {
	formData := map[string]string{
		"thread_id": documentID,
		"wipeout":   "false", // Set to true for permanent deletion
	}
	endpoint := "/threads/delete"
	resp, err := c.makeFormRequest("POST", endpoint, formData)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

// GetRecentThreads retrieves recent threads for the current user
func (c *Client) GetRecentThreads(limit int) ([]Document, error) {
	endpoint := "/threads/recent"
	if limit > 0 {
		endpoint += fmt.Sprintf("?count=%d", limit)
	}

	resp, err := c.makeRequest("GET", endpoint, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Read the response body to determine the structure
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Try to decode as the complex map response structure
	var response RecentThreadsResponse
	if err := json.Unmarshal(respBody, &response); err == nil && len(response) > 0 {
		// Convert the map to an array of documents
		threads := make([]Document, 0, len(response))
		for _, threadData := range response {
			threads = append(threads, threadData.Thread)
		}
		return threads, nil
	}

	// If that fails, try to decode as direct array of documents (fallback)
	var threads []Document
	if err := json.Unmarshal(respBody, &threads); err == nil {
		return threads, nil
	}

	// If both fail, try to decode as array of SearchResponse objects (similar to search API)
	var searchResponses []SearchResponse
	if err := json.Unmarshal(respBody, &searchResponses); err == nil {
		threads := make([]Document, len(searchResponses))
		for i, item := range searchResponses {
			threads[i] = item.Thread
		}
		return threads, nil
	}

	return nil, fmt.Errorf("failed to decode response: unrecognized response format. Response body: %s", string(respBody))
}

// GetUser retrieves user information by ID
func (c *Client) GetUser(userID string) (*User, error) {
	endpoint := fmt.Sprintf("/users/%s", userID)

	resp, err := c.makeRequest("GET", endpoint, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var user User
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &user, nil
}
