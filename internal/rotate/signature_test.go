package rotate

import "testing"

func TestSignature_SameMapSameHash(t *testing.T) {
	m := map[string]string{"A": "1", "B": "2"}
	if signature(m) != signature(m) {
		t.Fatal("same map must produce same signature")
	}
}

func TestSignature_DifferentValueDifferentHash(t *testing.T) {
	a := map[string]string{"KEY": "old"}
	b := map[string]string{"KEY": "new"}
	if signature(a) == signature(b) {
		t.Fatal("different values must produce different signatures")
	}
}

func TestSignature_EmptyMap(t *testing.T) {
	sig := signature(map[string]string{})
	if sig == "" {
		t.Fatal("signature of empty map must not be empty string")
	}
}

func TestSignature_ExtraKeyDiffers(t *testing.T) {
	a := map[string]string{"K": "v"}
	b := map[string]string{"K": "v", "EXTRA": "x"}
	if signature(a) == signature(b) {
		t.Fatal("maps with different keys must produce different signatures")
	}
}
