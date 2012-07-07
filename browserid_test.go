package browserid_test

import (
	"github.com/nshah/go.browserid"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHas(t *testing.T) {
	req, err := http.NewRequest("GET", "http://example.com", nil)
	if err != nil {
		t.Fatalf("Error creating request: %s", err)
	}
	if browserid.Has(req) {
		t.Fatalf("Error was not expecting request to have id")
	}
}

func TestGet(t *testing.T) {
	req, err := http.NewRequest("GET", "http://example.com", nil)
	if err != nil {
		t.Fatalf("Error creating request: %s", err)
	}
	if browserid.Has(req) {
		t.Fatalf("Error was not expecting request to have id")
	}
	w := httptest.NewRecorder()
	id1 := browserid.Get(w, req)
	if id1 == browserid.FailID {
		t.Fatalf("Error got fail ID: %s", id1)
	}
	if w.Header().Get("Set-Cookie") == "" {
		t.Fatalf("Error was expecting a Set-Cookie header")
	}
	id2 := browserid.Get(w, req)
	if id1 != id2 {
		t.Fatalf("Error got different ids: %s / %s", id1, id2)
	}
}
