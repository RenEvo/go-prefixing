package prefixing

import "testing"
import "github.com/ryanuber/go-glob"

func BenchmarkGlobs(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		glob.Glob("/api/hosts/*/metrics", "/api/hosts/dev-test01/metrics")
	}
}
