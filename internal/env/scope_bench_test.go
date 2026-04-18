package env

import "testing"

func BenchmarkScopeSnapshot_1000Keys(b *testing.B) {
	s := make(Snapshot, 1000)
	for i := 0; i < 500; i++ {
		s["APP_KEY_"+itoa(i)] = "value"
	}
	for i := 0; i < 500; i++ {
		s["OTHER_KEY_"+itoa(i)] = "value"
	}
	opts := DefaultScopeOptions("APP_")
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_ = ScopeSnapshot(s, opts)
	}
}

func BenchmarkNamespaceSnapshot_500Keys(b *testing.B) {
	s := make(Snapshot, 500)
	for i := 0; i < 500; i++ {
		s["KEY_"+itoa(i)] = "value"
	}
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_ = NamespaceSnapshot(s, "NS_")
	}
}

func itoa(n int) string {
	const digits = "0123456789"
	if n == 0 {
		return "0"
	}
	buf := make([]byte, 0, 10)
	for n > 0 {
		buf = append([]byte{digits[n%10]}, buf...)
		n /= 10
	}
	return string(buf)
}
