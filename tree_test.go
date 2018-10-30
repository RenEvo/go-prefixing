package prefixing

import (
	"fmt"
	"net/http"
	"testing"
)

func printChildren(t *testing.T, n *node, prefix string) {
	t.Logf(" %02d:%02d %s%s[%d] %v %t %d", n.priority, n.maxParams, prefix, n.path, len(n.children), n.handle, n.wildChild, n.nType)
	for l := len(n.path); l > 0; l-- {
		prefix += " "
	}
	for _, child := range n.children {
		printChildren(t, child, prefix)
	}
}

func TestRouterWithPathCleaner(t *testing.T) {
	tree := &node{}

	tree.addRoute(normalize("/api/hosts/{host-id}/metrics"), http.NotFoundHandler())
	tree.addRoute(normalize("/api/hosts/{host}/metrics/{program:[A-Za-z0-9]+}/events"), http.NotFoundHandler())
	tree.addRoute(normalize("/api/hosts/*/test"), http.NotFoundHandler())
	tree.addRoute(normalize("/api/metrics/:program"), http.NotFoundHandler())
	tree.addRoute(normalize("/api/metrics/*/test"), http.NotFoundHandler())

	// this doesn't work, because we don't have the longest prefix working properly, need to get this working.... in this codebase
	tree.addRoute("/api/hosts", http.NotFoundHandler())
	tree.addRoute("/api/hosts/*path", http.NotFoundHandler())

	printChildren(t, tree, "")

	if h, _, err := tree.getValue("/api/hosts"); h == nil {
		t.Errorf("failed to resolve /api/hosts: %v", err)
	}

	if h, _, err := tree.getValue("/api/hosts/"); h == nil {
		t.Errorf("failed to resolve /api/hosts: %v", err)
	}

	// this doesn't fail
	if h, _, err := tree.getValue("/api/hosts/host1"); h == nil {
		t.Errorf("failed to resolve /api/hosts/host1: %v", err)
	}
}

func BenchmarkHTTPRouterTree(b *testing.B) {
	tree := &node{}

	tree.addRoute("/api/hosts/:id/metrics", http.NotFoundHandler())

	b.ReportAllocs()
	b.ResetTimer()

	// "/api/hosts/*/metrics", "/api/hosts/dev-test01/metrics"
	for i := 0; i < b.N; i++ {
		tree.getValue("/api/hosts/dev6-test01/metrics")
	}
}

func BenchmarkRouterTreeWithWildcard_Lots(b *testing.B) {
	tree := &node{}

	for i := 0; i < 1000; i++ {
		tree.addRoute("/api/"+fmt.Sprintf("%d", i)+"/babs", http.NotFoundHandler())
	}

	tree.addRoute(normalize("/api/hosts/*/metrics"), http.NotFoundHandler())

	b.ReportAllocs()
	b.ResetTimer()

	// "/api/hosts/*/metrics", "/api/hosts/dev-test01/metrics"
	for i := 0; i < b.N; i++ {
		tree.getValue("/api/hosts/dev6-test01/metrics")
	}
}
