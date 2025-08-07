package quip

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewClient(t *testing.T) {
	token := "test-token"
	client := NewClient(token)

	if client.token != token {
		t.Errorf("Expected token %s, got %s", token, client.token)
	}

	if client.baseURL != BaseURL {
		t.Errorf("Expected baseURL %s, got %s", BaseURL, client.baseURL)
	}

	if client.httpClient.Timeout != Timeout {
		t.Errorf("Expected timeout %v, got %v", Timeout, client.httpClient.Timeout)
	}
}

func TestClient_GetCurrentUser(t *testing.T) {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify the request
		if r.URL.Path != "/users/current" {
			t.Errorf("Expected path /users/current, got %s", r.URL.Path)
		}

		if r.Method != "GET" {
			t.Errorf("Expected GET method, got %s", r.Method)
		}

		auth := r.Header.Get("Authorization")
		if auth != "Bearer test-token" {
			t.Errorf("Expected Authorization header 'Bearer test-token', got %s", auth)
		}

		// Return mock user data
		user := User{
			ID:    "user123",
			Name:  "Test User",
			Email: "test@example.com",
			URL:   "https://quip.com/user123",
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(user)
	}))
	defer server.Close()

	// Create client with test server URL
	client := NewClient("test-token")
	client.baseURL = server.URL

	// Test the method
	user, err := client.GetCurrentUser()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if user.ID != "user123" {
		t.Errorf("Expected user ID 'user123', got %s", user.ID)
	}

	if user.Name != "Test User" {
		t.Errorf("Expected user name 'Test User', got %s", user.Name)
	}

	if user.Email != "test@example.com" {
		t.Errorf("Expected user email 'test@example.com', got %s", user.Email)
	}
}

func TestClient_GetCurrentUser_Error(t *testing.T) {
	// Create a test server that returns an error
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte(`{"error": "Invalid token"}`))
	}))
	defer server.Close()

	// Create client with test server URL
	client := NewClient("invalid-token")
	client.baseURL = server.URL

	// Test the method
	_, err := client.GetCurrentUser()
	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	expectedError := "API error 401: {\"error\": \"Invalid token\"}"
	if err.Error() != expectedError {
		t.Errorf("Expected error %s, got %s", expectedError, err.Error())
	}
}

func TestClient_SearchDocuments(t *testing.T) {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify the request
		expectedPath := "/threads/search"
		if r.URL.Path != expectedPath {
			t.Errorf("Expected path %s, got %s", expectedPath, r.URL.Path)
		}

		query := r.URL.Query().Get("query")
		if query != "test query" {
			t.Errorf("Expected query 'test query', got %s", query)
		}

		count := r.URL.Query().Get("count")
		if count != "5" {
			t.Errorf("Expected count '5', got %s", count)
		}

		// Return mock search results in actual API format (array of SearchResponse)
		apiResponse := []SearchResponse{
			{
				Thread: Document{
					ID:       "doc123",
					Title:    "Test Document",
					Link:     "https://quip.com/doc123",
					AuthorID: "user123",
				},
			},
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(apiResponse)
	}))
	defer server.Close()

	// Create client with test server URL
	client := NewClient("test-token")
	client.baseURL = server.URL

	// Test the method
	result, err := client.SearchDocuments("test query", 5)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(result.Documents) != 1 {
		t.Errorf("Expected 1 document, got %d", len(result.Documents))
	}

	doc := result.Documents[0]
	if doc.ID != "doc123" {
		t.Errorf("Expected document ID 'doc123', got %s", doc.ID)
	}

	if doc.Title != "Test Document" {
		t.Errorf("Expected document title 'Test Document', got %s", doc.Title)
	}
}

func TestClient_GetDocument(t *testing.T) {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify the request
		expectedPath := "/threads/doc123"
		if r.URL.Path != expectedPath {
			t.Errorf("Expected path %s, got %s", expectedPath, r.URL.Path)
		}

		// Verify format=markdown query parameter
		format := r.URL.Query().Get("format")
		if format != "markdown" {
			t.Errorf("Expected format 'markdown', got %s", format)
		}

		// Return mock document data
		doc := Document{
			ID:       "doc123",
			Title:    "Test Document",
			HTML:     "<p>Test content</p>",
			Link:     "https://quip.com/doc123",
			AuthorID: "user123",
			Type:     "document",
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(doc)
	}))
	defer server.Close()

	// Create client with test server URL
	client := NewClient("test-token")
	client.baseURL = server.URL

	// Test the method
	doc, err := client.GetDocument("doc123")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if doc.ID != "doc123" {
		t.Errorf("Expected document ID 'doc123', got %s", doc.ID)
	}

	if doc.Title != "Test Document" {
		t.Errorf("Expected document title 'Test Document', got %s", doc.Title)
	}

	if doc.HTML != "<p>Test content</p>" {
		t.Errorf("Expected document HTML '<p>Test content</p>', got %s", doc.HTML)
	}
}

func TestClient_CreateDocument(t *testing.T) {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify the request
		expectedPath := "/threads/new-document"
		if r.URL.Path != expectedPath {
			t.Errorf("Expected path %s, got %s", expectedPath, r.URL.Path)
		}

		if r.Method != "POST" {
			t.Errorf("Expected POST method, got %s", r.Method)
		}

		// Verify content type
		contentType := r.Header.Get("Content-Type")
		if contentType != "application/x-www-form-urlencoded" {
			t.Errorf("Expected Content-Type 'application/x-www-form-urlencoded', got %s", contentType)
		}

		// Verify form data
		err := r.ParseForm()
		if err != nil {
			t.Fatalf("Failed to parse form data: %v", err)
		}

		if r.FormValue("title") != "New Document" {
			t.Errorf("Expected title 'New Document', got %s", r.FormValue("title"))
		}

		if r.FormValue("content") != "<p>New content</p>" {
			t.Errorf("Expected content '<p>New content</p>', got %s", r.FormValue("content"))
		}

		if r.FormValue("format") != "html" {
			t.Errorf("Expected format 'html', got %s", r.FormValue("format"))
		}

		// Return mock created document
		doc := Document{
			ID:       "newdoc123",
			Title:    "New Document",
			HTML:     "<p>New content</p>",
			Link:     "https://quip.com/newdoc123",
			AuthorID: "user123",
			Type:     "document",
			Created:  1640995200000000, // Mock timestamp
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(doc)
	}))
	defer server.Close()

	// Create client with test server URL
	client := NewClient("test-token")
	client.baseURL = server.URL

	// Test the method
	doc, err := client.CreateDocument("New Document", "<p>New content</p>")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if doc.ID != "newdoc123" {
		t.Errorf("Expected document ID 'newdoc123', got %s", doc.ID)
	}

	if doc.Title != "New Document" {
		t.Errorf("Expected document title 'New Document', got %s", doc.Title)
	}
}

func TestClient_GetDocumentComments(t *testing.T) {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify the request
		expectedPath := "/threads/doc123/messages"
		if r.URL.Path != expectedPath {
			t.Errorf("Expected path %s, got %s", expectedPath, r.URL.Path)
		}

		// Return mock comments
		comments := []Comment{
			{
				ID:       "comment123",
				Text:     "Great document!",
				AuthorID: "user123",
				Created:  1640995200000000,
				Visible:  true,
			},
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(comments)
	}))
	defer server.Close()

	// Create client with test server URL
	client := NewClient("test-token")
	client.baseURL = server.URL

	// Test the method
	comments, err := client.GetDocumentComments("doc123")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(comments) != 1 {
		t.Errorf("Expected 1 comment, got %d", len(comments))
	}

	comment := comments[0]
	if comment.ID != "comment123" {
		t.Errorf("Expected comment ID 'comment123', got %s", comment.ID)
	}

	if comment.Text != "Great document!" {
		t.Errorf("Expected comment text 'Great document!', got %s", comment.Text)
	}
}

func TestClient_GetUser(t *testing.T) {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify the request
		expectedPath := "/users/user123"
		if r.URL.Path != expectedPath {
			t.Errorf("Expected path %s, got %s", expectedPath, r.URL.Path)
		}

		// Return mock user data
		user := User{
			ID:    "user123",
			Name:  "John Doe",
			Email: "john@example.com",
			URL:   "https://quip.com/user123",
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(user)
	}))
	defer server.Close()

	// Create client with test server URL
	client := NewClient("test-token")
	client.baseURL = server.URL

	// Test the method
	user, err := client.GetUser("user123")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if user.ID != "user123" {
		t.Errorf("Expected user ID 'user123', got %s", user.ID)
	}

	if user.Name != "John Doe" {
		t.Errorf("Expected user name 'John Doe', got %s", user.Name)
	}
}

func TestClient_EditDocument(t *testing.T) {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify the request
		expectedPath := "/threads/edit-document/doc123"
		if r.URL.Path != expectedPath {
			t.Errorf("Expected path %s, got %s", expectedPath, r.URL.Path)
		}

		if r.Method != "POST" {
			t.Errorf("Expected POST method, got %s", r.Method)
		}

		// Verify content type
		contentType := r.Header.Get("Content-Type")
		if contentType != "application/x-www-form-urlencoded" {
			t.Errorf("Expected Content-Type 'application/x-www-form-urlencoded', got %s", contentType)
		}

		// Verify form data
		err := r.ParseForm()
		if err != nil {
			t.Fatalf("Failed to parse form data: %v", err)
		}

		if r.FormValue("content") != "<p>Updated content</p>" {
			t.Errorf("Expected content '<p>Updated content</p>', got %s", r.FormValue("content"))
		}

		if r.FormValue("operation") != "REPLACE" {
			t.Errorf("Expected operation 'REPLACE', got %s", r.FormValue("operation"))
		}

		if r.FormValue("format") != "html" {
			t.Errorf("Expected format 'html', got %s", r.FormValue("format"))
		}

		// Return mock updated document
		doc := Document{
			ID:       "doc123",
			Title:    "Updated Document",
			HTML:     "<p>Updated content</p>",
			Link:     "https://quip.com/doc123",
			AuthorID: "user123",
			Type:     "document",
			Updated:  1640995300000000, // Mock timestamp
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(doc)
	}))
	defer server.Close()

	// Create client with test server URL
	client := NewClient("test-token")
	client.baseURL = server.URL

	// Test the method
	doc, err := client.EditDocument("doc123", "<p>Updated content</p>", "REPLACE", "html")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if doc.ID != "doc123" {
		t.Errorf("Expected document ID 'doc123', got %s", doc.ID)
	}

	if doc.Title != "Updated Document" {
		t.Errorf("Expected document title 'Updated Document', got %s", doc.Title)
	}
}

func TestClient_DeleteDocument(t *testing.T) {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify the request
		expectedPath := "/threads/delete/doc123"
		if r.URL.Path != expectedPath {
			t.Errorf("Expected path %s, got %s", expectedPath, r.URL.Path)
		}

		if r.Method != "POST" {
			t.Errorf("Expected POST method, got %s", r.Method)
		}

		// Return success (empty response is fine for delete)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	// Create client with test server URL
	client := NewClient("test-token")
	client.baseURL = server.URL

	// Test the method
	err := client.DeleteDocument("doc123")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}

func TestClient_GetRecentThreads(t *testing.T) {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify the request
		expectedPath := "/threads/recent"
		if r.URL.Path != expectedPath {
			t.Errorf("Expected path %s, got %s", expectedPath, r.URL.Path)
		}

		count := r.URL.Query().Get("count")
		if count != "5" {
			t.Errorf("Expected count '5', got %s", count)
		}

		// Return mock recent threads
		threads := []Document{
			{
				ID:       "thread1",
				Title:    "Recent Thread 1",
				Link:     "https://quip.com/thread1",
				AuthorID: "user123",
				Type:     "document",
				Updated:  1640995400000000,
			},
			{
				ID:       "thread2",
				Title:    "Recent Thread 2",
				Link:     "https://quip.com/thread2",
				AuthorID: "user456",
				Type:     "chat",
				Updated:  1640995300000000,
			},
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(threads)
	}))
	defer server.Close()

	// Create client with test server URL
	client := NewClient("test-token")
	client.baseURL = server.URL

	// Test the method
	threads, err := client.GetRecentThreads(5)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(threads) != 2 {
		t.Errorf("Expected 2 threads, got %d", len(threads))
	}

	if threads[0].ID != "thread1" {
		t.Errorf("Expected first thread ID 'thread1', got %s", threads[0].ID)
	}

	if threads[0].Title != "Recent Thread 1" {
		t.Errorf("Expected first thread title 'Recent Thread 1', got %s", threads[0].Title)
	}

	if threads[1].Type != "chat" {
		t.Errorf("Expected second thread type 'chat', got %s", threads[1].Type)
	}
}
