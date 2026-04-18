package env

import (
	"testing"
)

func TestMergeSnapshots_OverlayWins(t *testing.T) {
	base := Snapshot{"HOME": "/root", "PATH": "/usr/bin"}
	overlay := Snapshot{"PATH": "/secret/bin", "TOKEN": "abc"}

	out := MergeSnapshots(base, overlay, DefaultMergeOptions())

	if out["PATH"] != "/secret/bin" {
		t.Errorf("expected overlay PATH, got %s", out["PATH"])
	}
	if out["HOME"] != "/root" {
		t.Errorf("expected base HOME, got %s", out["HOME"])
	}
	if out["TOKEN"] != "abc" {
		t.Errorf("expected TOKEN from overlay, got %s", out["TOKEN"])
	}
}

func TestMergeSnapshots_BaseWins(t *testing.T) {
	base := Snapshot{"PATH": "/usr/bin"}
	overlay := Snapshot{"PATH": "/secret/bin", "TOKEN": "abc"}

	out := MergeSnapshots(base, overlay, MergeOptions{Strategy: BaseWins})

	if out["PATH"] != "/usr/bin" {
		t.Errorf("expected base PATH preserved, got %s", out["PATH"])
	}
	if out["TOKEN"] != "abc" {
		t.Errorf("expected new key TOKEN, got %s", out["TOKEN"])
	}
}

func TestMergeSnapshots_DoesNotMutateBase(t *testing.T) {
	base := Snapshot{"KEY": "original"}
	overlay := Snapshot{"KEY": "changed"}

	MergeSnapshots(base, overlay, DefaultMergeOptions())

	if base["KEY"] != "original" {
		t.Error("base snapshot was mutated")
	}
}

func TestMergeSnapshots_DenyListFiltersOverlay(t *testing.T) {
	dl := NewDenyList([]string{"SECRET_"}, []string{})
	base := Snapshot{}
	overlay := Snapshot{"SECRET_TOKEN": "s3cr3t", "SAFE_KEY": "visible"}

	out := MergeSnapshots(base, overlay, MergeOptions{
		Strategy: OverlayWins,
		DenyList: dl,
	})

	if _, ok := out["SECRET_TOKEN"]; ok {
		t.Error("expected SECRET_TOKEN to be filtered by deny list")
	}
	if out["SAFE_KEY"] != "visible" {
		t.Errorf("expected SAFE_KEY to pass through, got %s", out["SAFE_KEY"])
	}
}

func TestMergeSnapshots_EmptyOverlay(t *testing.T) {
	base := Snapshot{"A": "1", "B": "2"}
	out := MergeSnapshots(base, Snapshot{}, DefaultMergeOptions())
	if len(out) != 2 {
		t.Errorf("expected 2 keys, got %d", len(out))
	}
}
