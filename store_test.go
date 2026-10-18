package gae

import "testing"

func TestNewGaeBaseStore(t *testing.T) {
	store := NewGaeBaseStore("TestKind")
	if store.entityKind != "TestKind" {
		t.Errorf("entityKind = %q, want %q", store.entityKind, "TestKind")
	}
}
