package session

import (
	"net/http"
	"testing"
)

func Test_CookieStore(t *testing.T) {
	originalPath := "/"
	store := NewCookieStore()
	store.Options.Path = originalPath

	r, err  := http.NewRequest("GET", "http://localhost:8080", nil)
	if err != nil {
		t.Fatal( "Create request fail: ",err)
	}

	session, err := store.New(r, "hello world")
	if err != nil {
		t.Fatal( "create session fail: ",err)
	}

	store.Options.Path = "/foo"
	if session.Options.Path != originalPath {
		t.Fatalf("bad session path: got %q, want %q", session.Options.Path, originalPath)
	}
}

// Test for GH-8 for FilesystemStore
func TestGH8FilesystemStore(t *testing.T) {
	originalPath := "/"
	store := NewFilesystemStore("")
	store.Options.Path = originalPath
	req, err := http.NewRequest("GET", "http://www.example.com", nil)
	if err != nil {
		t.Fatal("failed to create request", err)
	}

	session, err := store.New(req, "hello")
	if err != nil {
		t.Fatal("failed to create session", err)
	}

	store.Options.Path = "/foo"
	if session.Options.Path != originalPath {
		t.Fatalf("bad session path: got %q, want %q", session.Options.Path, originalPath)
	}
}
