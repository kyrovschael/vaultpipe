package template

import (
	"testing"
)

func secrets() map[string]map[string]string {
	return map[string]map[string]string{
		"secret/db": {
			"username": "admin",
			"password": "s3cr3t",
		},
		"secret/api": {
			"key": "abc123",
		},
	}
}

func TestRender_Simple(t *testing.T) {
	r := NewRenderer(secrets())
	out, err := r.Render("user={{vault:secret/db#username}}")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out != "user=admin" {
		t.Errorf("got %q, want %q", out, "user=admin")
	}
}

func TestRender_MultiplePlaceholders(t *testing.T) {
	r := NewRenderer(secrets())
	out, err := r.Render("{{vault:secret/db#username}}:{{vault:secret/db#password}}")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out != "admin:s3cr3t" {
		t.Errorf("got %q", out)
	}
}

func TestRender_NoPlaceholder(t *testing.T) {
	r := NewRenderer(secrets())
	out, err := r.Render("plain-value")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out != "plain-value" {
		t.Errorf("got %q", out)
	}
}

func TestRender_MissingPath(t *testing.T) {
	r := NewRenderer(secrets())
	_, err := r.Render("{{vault:secret/missing#key}}")
	if err == nil {
		t.Fatal("expected error for missing path")
	}
}

func TestRender_MissingField(t *testing.T) {
	r := NewRenderer(secrets())
	_, err := r.Render("{{vault:secret/db#nofield}}")
	if err == nil {
		t.Fatal("expected error for missing field")
	}
}

func TestRenderMap(t *testing.T) {
	r := NewRenderer(secrets())
	input := map[string]string{
		"DB_USER": "{{vault:secret/db#username}}",
		"API_KEY": "{{vault:secret/api#key}}",
		"STATIC":  "unchanged",
	}
	out, err := r.RenderMap(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["DB_USER"] != "admin" || out["API_KEY"] != "abc123" || out["STATIC"] != "unchanged" {
		t.Errorf("unexpected output: %v", out)
	}
}
