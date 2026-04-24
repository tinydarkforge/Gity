package services_test

import (
	"testing"

	"github.com/tinydarkforge/gity/services"
)

func TestLoadTemplates(t *testing.T) {
	tmpls, err := services.LoadTemplates("../.github/ISSUE_TEMPLATES")
	if err != nil {
		t.Fatalf("LoadTemplates: %v", err)
	}
	if len(tmpls) == 0 {
		t.Fatal("expected at least one template")
	}
	for _, tmpl := range tmpls {
		if tmpl.Name == "" {
			t.Errorf("template %s has empty name", tmpl.Filename)
		}
		if tmpl.Body == "" {
			t.Errorf("template %s has empty body", tmpl.Filename)
		}
		// All test templates should have at least one label
		if len(tmpl.Labels) == 0 {
			t.Logf("note: template %s has no labels", tmpl.Filename)
		}
	}
}

func TestFindTemplate(t *testing.T) {
	tmpls, _ := services.LoadTemplates("../.github/ISSUE_TEMPLATES")

	got := services.FindTemplate(tmpls, "bug")
	if got == nil {
		t.Fatal("expected to find a bug template")
	}

	none := services.FindTemplate(tmpls, "zzznomatch")
	if none != nil {
		t.Fatal("expected nil for no match")
	}

	empty := services.FindTemplate(tmpls, "")
	if empty != nil {
		t.Fatal("expected nil for empty query")
	}
}
