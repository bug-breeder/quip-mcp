package quip

import (
	"fmt"
	"os"
	"testing"
	"time"
)

// Integration tests that call the real Quip API
// These tests require a valid QUIP_API_TOKEN environment variable
// Run with: go test -tags=integration ./pkg/quip -v

// Helper function to check if integration tests should run
func skipIfNoToken(t *testing.T) *Client {
	token := os.Getenv("QUIP_API_TOKEN")
	if token == "" {
		t.Skip("Skipping integration test: QUIP_API_TOKEN environment variable not set")
	}
	return NewClient(token)
}

// TestIntegration_GetCurrentUser tests the GetCurrentUser functionality
func TestIntegration_GetCurrentUser(t *testing.T) {
	client := skipIfNoToken(t)

	user, err := client.GetCurrentUser()
	if err != nil {
		t.Fatalf("GetCurrentUser failed: %v", err)
	}

	// Validate required fields
	if user.ID == "" {
		t.Error("User ID should not be empty")
	}
	if user.Name == "" {
		t.Error("User name should not be empty")
	}

	t.Logf("✅ GetCurrentUser successful - User: %s (ID: %s)", user.Name, user.ID)
}

// TestIntegration_GetRecentThreads tests the GetRecentThreads functionality
func TestIntegration_GetRecentThreads(t *testing.T) {
	client := skipIfNoToken(t)

	threads, err := client.GetRecentThreads(5)
	if err != nil {
		t.Fatalf("GetRecentThreads failed: %v", err)
	}

	// We can't guarantee there will be threads, but the call should succeed
	t.Logf("✅ GetRecentThreads successful - Found %d threads", len(threads))

	// If there are threads, validate their structure
	for i, thread := range threads {
		if thread.ID == "" {
			t.Errorf("Thread %d has empty ID", i)
		}
		if thread.Title == "" {
			t.Errorf("Thread %d has empty title", i)
		}
		t.Logf("   Thread %d: %s (ID: %s)", i+1, thread.Title, thread.ID)
	}
}

// TestIntegration_SearchDocuments tests the SearchDocuments functionality
func TestIntegration_SearchDocuments(t *testing.T) {
	client := skipIfNoToken(t)

	// Search for documents - using a common word that might exist
	result, err := client.SearchDocuments("document", 3)
	if err != nil {
		t.Fatalf("SearchDocuments failed: %v", err)
	}

	t.Logf("✅ SearchDocuments successful - Found %d documents", len(result.Documents))

	// Validate document structure if any found
	for i, doc := range result.Documents {
		if doc.ID == "" {
			t.Errorf("Document %d has empty ID", i)
		}
		if doc.Title == "" {
			t.Errorf("Document %d has empty title", i)
		}
		t.Logf("   Document %d: %s (ID: %s)", i+1, doc.Title, doc.ID)
	}
}

// TestIntegration_DocumentCRUD tests the full CRUD lifecycle for documents
func TestIntegration_DocumentCRUD(t *testing.T) {
	client := skipIfNoToken(t)

	// Test document lifecycle with timestamp to ensure uniqueness
	timestamp := time.Now().Unix()
	testTitle := fmt.Sprintf("Integration Test Document %d", timestamp)
	testContent := "<p>This is a test document created by integration tests. It should be deleted automatically.</p>"

	// 1. CREATE: Create a test document
	t.Log("🔄 Creating test document...")
	doc, err := client.CreateDocument(testTitle, testContent)
	if err != nil {
		t.Fatalf("CreateDocument failed: %v", err)
	}

	if doc.ID == "" {
		t.Fatal("Created document has empty ID")
	}
	if doc.Title != testTitle {
		t.Errorf("Expected title %s, got %s", testTitle, doc.Title)
	}

	t.Logf("✅ Document created successfully - ID: %s, Title: %s", doc.ID, doc.Title)
	documentID := doc.ID

	// Ensure cleanup even if other tests fail
	defer func() {
		t.Log("🧹 Cleaning up test document...")
		if err := client.DeleteDocument(documentID); err != nil {
			t.Logf("⚠️  Warning: Failed to cleanup test document %s: %v", documentID, err)
		} else {
			t.Logf("✅ Test document %s cleaned up successfully", documentID)
		}
	}()

	// 2. READ: Get the document back
	t.Log("🔄 Reading test document...")
	retrievedDoc, err := client.GetDocument(documentID)
	if err != nil {
		t.Fatalf("GetDocument failed: %v", err)
	}

	if retrievedDoc.ID != documentID {
		t.Errorf("Expected ID %s, got %s", documentID, retrievedDoc.ID)
	}
	if retrievedDoc.Title != testTitle {
		t.Errorf("Expected title %s, got %s", testTitle, retrievedDoc.Title)
	}

	t.Logf("✅ Document retrieved successfully - Title: %s", retrievedDoc.Title)

	// 3. UPDATE: Edit the document
	t.Log("🔄 Updating test document...")
	updatedContent := "<p>This content has been updated by integration tests.</p>"
	updatedDoc, err := client.EditDocument(documentID, updatedContent, "REPLACE", "html")
	if err != nil {
		t.Fatalf("EditDocument failed: %v", err)
	}

	if updatedDoc.ID != documentID {
		t.Errorf("Expected ID %s, got %s", documentID, updatedDoc.ID)
	}

	t.Logf("✅ Document updated successfully - ID: %s", updatedDoc.ID)

	// 4. GET COMMENTS: Try to get comments (might be empty)
	t.Log("🔄 Getting document comments...")
	comments, err := client.GetDocumentComments(documentID)
	if err != nil {
		t.Logf("⚠️  GetDocumentComments failed (this may be expected): %v", err)
	} else {
		t.Logf("✅ GetDocumentComments successful - Found %d comments", len(comments))
	}

	// 5. DELETE: Will be handled by defer cleanup
	t.Log("✅ CRUD test completed successfully")
}

// TestIntegration_UserOperations tests user-related operations
func TestIntegration_UserOperations(t *testing.T) {
	client := skipIfNoToken(t)

	// Get current user first
	currentUser, err := client.GetCurrentUser()
	if err != nil {
		t.Fatalf("GetCurrentUser failed: %v", err)
	}

	t.Logf("✅ Current user retrieved: %s (ID: %s)", currentUser.Name, currentUser.ID)

	// Try to get the same user by ID
	userByID, err := client.GetUser(currentUser.ID)
	if err != nil {
		t.Fatalf("GetUser by ID failed: %v", err)
	}

	if userByID.ID != currentUser.ID {
		t.Errorf("Expected user ID %s, got %s", currentUser.ID, userByID.ID)
	}

	t.Logf("✅ User by ID retrieved: %s", userByID.Name)
}

// TestIntegration_ErrorHandling tests error handling with invalid inputs
func TestIntegration_ErrorHandling(t *testing.T) {
	client := skipIfNoToken(t)

	// Test with invalid document ID
	t.Log("🔄 Testing error handling with invalid document ID...")
	_, err := client.GetDocument("invalid-document-id")
	if err == nil {
		t.Error("Expected error for invalid document ID, but got none")
	} else {
		t.Logf("✅ Error handling working - Invalid document ID error: %v", err)
	}

	// Test with invalid user ID
	t.Log("🔄 Testing error handling with invalid user ID...")
	_, err = client.GetUser("invalid-user-id")
	if err == nil {
		t.Error("Expected error for invalid user ID, but got none")
	} else {
		t.Logf("✅ Error handling working - Invalid user ID error: %v", err)
	}
}

// TestIntegration_APIResponseStructures tests the structure of API responses
func TestIntegration_APIResponseStructures(t *testing.T) {
	client := skipIfNoToken(t)

	t.Log("🔄 Testing API response structures...")

	// Test search response structure
	searchResult, err := client.SearchDocuments("test", 1)
	if err != nil {
		t.Logf("⚠️  SearchDocuments failed: %v", err)
	} else {
		t.Logf("✅ Search response structure valid - %d documents returned", len(searchResult.Documents))
	}

	// Test recent threads response structure (this is where we had issues)
	threads, err := client.GetRecentThreads(1)
	if err != nil {
		t.Fatalf("GetRecentThreads failed: %v", err)
	} else {
		t.Logf("✅ Recent threads response structure valid - %d threads returned", len(threads))
	}

	// Test user response structure
	user, err := client.GetCurrentUser()
	if err != nil {
		t.Fatalf("GetCurrentUser failed: %v", err)
	} else {
		t.Logf("✅ User response structure valid - User: %s", user.Name)
	}
}

// Benchmark tests for performance monitoring
func BenchmarkIntegration_GetCurrentUser(b *testing.B) {
	token := os.Getenv("QUIP_API_TOKEN")
	if token == "" {
		b.Skip("Skipping benchmark: QUIP_API_TOKEN environment variable not set")
	}
	client := NewClient(token)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := client.GetCurrentUser()
		if err != nil {
			b.Fatalf("GetCurrentUser failed: %v", err)
		}
	}
}

func BenchmarkIntegration_GetRecentThreads(b *testing.B) {
	token := os.Getenv("QUIP_API_TOKEN")
	if token == "" {
		b.Skip("Skipping benchmark: QUIP_API_TOKEN environment variable not set")
	}
	client := NewClient(token)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := client.GetRecentThreads(5)
		if err != nil {
			b.Fatalf("GetRecentThreads failed: %v", err)
		}
	}
}
