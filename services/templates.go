package services

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"gopkg.in/yaml.v3"

	"github.com/tinydarkforge/gity/types"
)

// LoadTemplates reads every *.md file in dir and parses YAML frontmatter.
// Files without frontmatter are skipped silently.
func LoadTemplates(dir string) ([]types.Template, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("read template dir %q: %w", dir, err)
	}

	var out []types.Template
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".md") {
			continue
		}
		path := filepath.Join(dir, e.Name())
		data, err := os.ReadFile(path)
		if err != nil {
			continue
		}
		t, ok := parseTemplate(string(data))
		if !ok {
			continue
		}
		t.Filename = e.Name()
		out = append(out, t)
	}
	sort.Slice(out, func(i, j int) bool {
		return out[i].DisplayName() < out[j].DisplayName()
	})
	return out, nil
}

func parseTemplate(raw string) (types.Template, bool) {
	var t types.Template
	if !strings.HasPrefix(raw, "---") {
		return t, false
	}
	rest := strings.TrimPrefix(raw, "---")
	end := strings.Index(rest, "\n---")
	if end < 0 {
		return t, false
	}
	front := rest[:end]
	body := strings.TrimLeft(rest[end+4:], "\r\n")
	if err := yaml.Unmarshal([]byte(front), &t); err != nil {
		return t, false
	}
	t.Body = body
	return t, true
}

// FindTemplate returns the first template whose filename contains the query
// (case-insensitive, without extension). Empty query returns nil.
func FindTemplate(all []types.Template, query string) *types.Template {
	q := strings.ToLower(strings.TrimSpace(query))
	if q == "" {
		return nil
	}
	for i, t := range all {
		name := strings.ToLower(strings.TrimSuffix(t.Filename, ".md"))
		if strings.Contains(name, q) {
			return &all[i]
		}
	}
	return nil
}
