package env

import (
	"testing"
)

func TestApply_DenyListFiltersDefaults(t *testing.T) {
	base := Snapshot{
		"PATH":           "/usr/bin",
		"HOME":           "/root",
		"VAULT_TOKEN":    "s.secret",
		"AWS_SECRET_KEY": "key",
	}
	cfg := DefaultInheritConfig()
	out := cfg.Apply(base)

	if _, ok := out["VAULT_TOKEN"]; ok {
		t.Error("expected VAULT_TOKEN to be filtered")
	}
	if _, ok := out["AWS_SECRET_KEY"]; ok {
		t.Error("expected AWS_SECRET_KEY to be filtered")
	}
	if out["PATH"] != "/usr/bin" {
		t.Errorf("expected PATH to be inherited, got %q", out["PATH"])
	}
}

func TestApply_AllowListOverridesDeny(t *testing.T) {
	base := Snapshot{
		"PATH":        "/usr/bin",
		"HOME":        "/root",
		"VAULT_TOKEN": "s.secret",
		"MY_APP_VAR":  "hello",
	}
	cfg := &InheritConfig{
		AllowList: []string{"PATH", "MY_APP_VAR"},
		DenyList:  DefaultDenyList(),
	}
	out := cfg.Apply(base)

	if len(out) != 2 {
		t.Errorf("expected 2 keys, got %d", len(out))
	}
	if out["PATH"] != "/usr/bin" {
		t.Errorf("unexpected PATH value: %q", out["PATH"])
	}
	if out["MY_APP_VAR"] != "hello" {
		t.Errorf("unexpected MY_APP_VAR value: %q", out["MY_APP_VAR"])
	}
	if _, ok := out["VAULT_TOKEN"]; ok {
		t.Error("VAULT_TOKEN should not appear when allowlist is set")
	}
}

func TestApply_EmptyBase(t *testing.T) {
	cfg := DefaultInheritConfig()
	out := cfg.Apply(Snapshot{})
	if len(out) != 0 {
		t.Errorf("expected empty snapshot, got %d keys", len(out))
	}
}

func TestApply_NilDenyListPassesAll(t *testing.T) {
	base := Snapshot{"FOO": "bar", "BAZ": "qux"}
	cfg := &InheritConfig{}
	out := cfg.Apply(base)
	if len(out) != 2 {
		t.Errorf("expected 2 keys, got %d", len(out))
	}
}
