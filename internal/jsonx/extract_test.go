package jsonx_test

import (
	"testing"

	"github.com/tinydarkforge/intake/internal/jsonx"
)

type result struct {
	Status    string   `json:"status"`
	Title     string   `json:"title"`
	Questions []string `json:"questions"`
}

func TestExtractClean(t *testing.T) {
	raw := `{"status":"ready","title":"Fix login redirect"}`
	var r result
	if err := jsonx.Extract(raw, &r); err != nil {
		t.Fatal(err)
	}
	if r.Title != "Fix login redirect" {
		t.Errorf("got %q", r.Title)
	}
}

func TestExtractFenced(t *testing.T) {
	raw := "```json\n{\"status\":\"needs_info\",\"questions\":[\"What environment?\"]}\n```"
	var r result
	if err := jsonx.Extract(raw, &r); err != nil {
		t.Fatal(err)
	}
	if r.Status != "needs_info" || len(r.Questions) == 0 {
		t.Errorf("unexpected result: %+v", r)
	}
}

func TestExtractTrailingProse(t *testing.T) {
	raw := `Here is the JSON you asked for:
{"status":"ready","title":"Add dark mode"}
Hope that helps!`
	var r result
	if err := jsonx.Extract(raw, &r); err != nil {
		t.Fatal(err)
	}
	if r.Title != "Add dark mode" {
		t.Errorf("got %q", r.Title)
	}
}

func TestExtractNoJSON(t *testing.T) {
	var r result
	if err := jsonx.Extract("nothing here", &r); err == nil {
		t.Fatal("expected error")
	}
}
