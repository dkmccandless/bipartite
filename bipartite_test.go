package bipartite

import (
	"reflect"
	"testing"
)

var adjTests = []struct {
	g    *Graph
	adj  map[interface{}]map[interface{}]bool
	adja map[interface{}][]interface{}
	adjb map[interface{}][]interface{}
}{
	{
		g: New(),
	},
	{
		g: &Graph{
			m: map[interface{}]map[interface{}]struct{}{
				"apple": map[interface{}]struct{}{"tree": struct{}{}},
				"tree":  map[interface{}]struct{}{"apple": struct{}{}},
			},
			as: map[interface{}]struct{}{"apple": struct{}{}},
			bs: map[interface{}]struct{}{"tree": struct{}{}},
		},
		adj: map[interface{}]map[interface{}]bool{
			"apple":     map[interface{}]bool{"tree": true, "rock": false},
			"spaghetti": map[interface{}]bool{"tree": false},
		},
		adja: map[interface{}][]interface{}{"apple": []interface{}{"tree"}},
		adjb: map[interface{}][]interface{}{"tree": []interface{}{"apple"}},
	},
	{
		g: &Graph{
			m: map[interface{}]map[interface{}]struct{}{
				"X": map[interface{}]struct{}{0: struct{}{}},
				"Y": map[interface{}]struct{}{0: struct{}{}, 1: struct{}{}},
				"Z": map[interface{}]struct{}{1: struct{}{}, 2: struct{}{}, 3: struct{}{}},
				0:   map[interface{}]struct{}{"X": struct{}{}, "Y": struct{}{}},
				1:   map[interface{}]struct{}{"Y": struct{}{}, "Z": struct{}{}},
				2:   map[interface{}]struct{}{"Z": struct{}{}},
				3:   map[interface{}]struct{}{"Z": struct{}{}},
			},
			as: map[interface{}]struct{}{"X": struct{}{}, "Y": struct{}{}, "Z": struct{}{}},
			bs: map[interface{}]struct{}{0: struct{}{}, 1: struct{}{}, 2: struct{}{}, 3: struct{}{}},
		},
		adj: map[interface{}]map[interface{}]bool{
			"W": map[interface{}]bool{0: false, 1: false, 2: false, 3: false},
			"X": map[interface{}]bool{0: true, 1: false, 2: false, 3: false},
			"Y": map[interface{}]bool{0: true, 1: true, 2: false, 3: false},
			"Z": map[interface{}]bool{0: false, 1: true, 2: true, 3: true},
		},
		adja: map[interface{}][]interface{}{"X": []interface{}{0}, "Y": []interface{}{0, 1}, "Z": []interface{}{1, 2, 3}},
		adjb: map[interface{}][]interface{}{0: []interface{}{"X", "Y"}, 1: []interface{}{"Y", "Z"}, 2: []interface{}{"Z"}, 3: []interface{}{"Z"}},
	},
}

func TestN(t *testing.T) {
	for _, test := range adjTests {
		if na := test.g.NA(); na != len(test.adja) {
			t.Errorf("NA(%+v): got %v, want %v", test.g, na, len(test.adja))
		}
		if nb := test.g.NB(); nb != len(test.adjb) {
			t.Errorf("NB(%+v): got %v, want %v", test.g, nb, len(test.adjb))
		}
	}
}

func TestDegAdjTo(t *testing.T) {
	for _, test := range adjTests {
		for node, adj := range test.adja {
			if _, ok := test.g.m[node]; !ok {
				t.Fatalf("Deg: %+v does not contain %v", test.g, node)
			}
			if deg := test.g.Deg(node); deg != len(adj) {
				t.Errorf("Deg(%+v, %v): got %v, want %v", test.g, node, deg, len(adj))
			}
			if got := test.g.AdjTo(node); !match(got, adj) {
				t.Errorf("AdjTo(%+v, %v): got %v, want %v", test.g, node, got, adj)
			}
		}
		for node, adj := range test.adjb {
			if _, ok := test.g.m[node]; !ok {
				t.Fatalf("Deg: %+v does not contain %v", test.g, node)
			}
			if deg := test.g.Deg(node); deg != len(adj) {
				t.Errorf("Deg(%+v, %v): got %v, want %v", test.g, node, deg, len(adj))
			}
			if got := test.g.AdjTo(node); !match(got, adj) {
				t.Errorf("AdjTo(%+v, %v): got %v, want %v", test.g, node, got, adj)
			}
		}
	}
}

func TestAdjacent(t *testing.T) {
	for _, test := range adjTests {
		for a := range test.adj {
			for b := range test.adj[a] {
				if got := test.g.Adjacent(a, b); got != test.adj[a][b] {
					t.Errorf("Adjacent(%+v: %v, %v): got %v, want %v", test.g, a, b, got, test.adj[a][b])
				}
			}
		}
	}
}

func TestAdd(t *testing.T) {
	tests := []struct {
		a, b       interface{}
		degA, degB int
		adj        bool // adjacent before Add
	}{
		{"X", 0, 1, 1, false},
		{"Y", 0, 1, 2, false},
		{"Z", 1, 1, 1, false},
		{"Y", 1, 2, 2, false},
		{"Z", 2, 2, 1, false},
		{"Z", 3, 3, 1, false},
		{"Z", 3, 3, 1, true},
	}
	g := New()
	for _, test := range tests {
		if adj := g.Adjacent(test.a, test.b); adj != test.adj {
			t.Errorf("Add(%+v before %v, %v): adjacent=%v, want %v", g, test.a, test.b, adj, test.adj)
		}
		g.Add(test.a, test.b)
		if !g.Adjacent(test.a, test.b) {
			t.Errorf("Add(%+v after %v, %v): not adjacent", g, test.a, test.b)
		}
		if degA := g.Deg(test.a); degA != test.degA {
			t.Errorf("Add(%+v after %v, %v): A degree %v, want %v", g, test.a, test.b, degA, test.degA)
		}
		if degB := g.Deg(test.b); degB != test.degB {
			t.Errorf("Add(%+v after %v, %v): B degree %v, want %v", g, test.a, test.b, degB, test.degB)
		}
	}
}

func TestDelete(t *testing.T) {
	g := &Graph{
		m: map[interface{}]map[interface{}]struct{}{
			"X": map[interface{}]struct{}{0: struct{}{}},
			"Y": map[interface{}]struct{}{0: struct{}{}, 1: struct{}{}},
			"Z": map[interface{}]struct{}{1: struct{}{}, 2: struct{}{}, 3: struct{}{}},
			0:   map[interface{}]struct{}{"X": struct{}{}, "Y": struct{}{}},
			1:   map[interface{}]struct{}{"Y": struct{}{}, "Z": struct{}{}},
			2:   map[interface{}]struct{}{"Z": struct{}{}},
			3:   map[interface{}]struct{}{"Z": struct{}{}},
		},
		as: map[interface{}]struct{}{"X": struct{}{}, "Y": struct{}{}, "Z": struct{}{}},
		bs: map[interface{}]struct{}{0: struct{}{}, 1: struct{}{}, 2: struct{}{}, 3: struct{}{}},
	}
	tests := []struct {
		a, b       interface{}
		degA, degB int
		adj        bool // adjacent before Delete
	}{
		{"W", 4, 0, 0, false}, // neither node exists
		{"W", 1, 0, 2, false}, // A exists but not B
		{"X", 4, 1, 0, false}, // B exists but not A
		{"X", 1, 1, 2, false}, // A and B not connected
		{"X", 0, 0, 1, true},
		{"Y", 0, 1, 0, true},
		{"Z", 1, 2, 1, true},
		{"Z", 1, 2, 1, false}, // already deleted
		{"Y", 1, 0, 0, true},
		{"Z", 2, 1, 0, true},
		{"Z", 3, 0, 0, true},
	}
	for _, test := range tests {
		if adj := g.Adjacent(test.a, test.b); adj != test.adj {
			t.Errorf("Delete(%+v before %v, %v): adjacent=%v, want %v", g, test.a, test.b, adj, test.adj)
		}
		g.Delete(test.a, test.b)
		if g.Adjacent(test.a, test.b) {
			t.Errorf("Delete(%+v after %v, %v): adjacent", g, test.a, test.b)
		}
		if degA := g.Deg(test.a); degA != test.degA {
			t.Errorf("Delete(%+v after %v, %v): got %v degree %v, want %v", g, test.a, test.b, test.a, degA, test.degA)
		}
		if degB := g.Deg(test.b); degB != test.degB {
			t.Errorf("Delete(%+v after %v, %v): got %v degree %v, want %v", g, test.a, test.b, test.b, degB, test.degB)
		}
	}
}

func TestRemoveA(t *testing.T) {
	g := &Graph{
		m: map[interface{}]map[interface{}]struct{}{
			"X": map[interface{}]struct{}{0: struct{}{}},
			"Y": map[interface{}]struct{}{0: struct{}{}, 1: struct{}{}},
			"Z": map[interface{}]struct{}{1: struct{}{}, 2: struct{}{}, 3: struct{}{}},
			0:   map[interface{}]struct{}{"X": struct{}{}, "Y": struct{}{}},
			1:   map[interface{}]struct{}{"Y": struct{}{}, "Z": struct{}{}},
			2:   map[interface{}]struct{}{"Z": struct{}{}},
			3:   map[interface{}]struct{}{"Z": struct{}{}},
		},
		as: map[interface{}]struct{}{"X": struct{}{}, "Y": struct{}{}, "Z": struct{}{}},
		bs: map[interface{}]struct{}{0: struct{}{}, 1: struct{}{}, 2: struct{}{}, 3: struct{}{}},
	}
	for _, test := range []struct {
		node interface{}
		want *Graph
	}{
		{
			"X",
			&Graph{
				m: map[interface{}]map[interface{}]struct{}{
					"Y": map[interface{}]struct{}{0: struct{}{}, 1: struct{}{}},
					"Z": map[interface{}]struct{}{1: struct{}{}, 2: struct{}{}, 3: struct{}{}},
					0:   map[interface{}]struct{}{"Y": struct{}{}},
					1:   map[interface{}]struct{}{"Y": struct{}{}, "Z": struct{}{}},
					2:   map[interface{}]struct{}{"Z": struct{}{}},
					3:   map[interface{}]struct{}{"Z": struct{}{}},
				},
				as: map[interface{}]struct{}{"Y": struct{}{}, "Z": struct{}{}},
				bs: map[interface{}]struct{}{0: struct{}{}, 1: struct{}{}, 2: struct{}{}, 3: struct{}{}},
			},
		},
		{
			"Y",
			&Graph{
				m: map[interface{}]map[interface{}]struct{}{
					"Z": map[interface{}]struct{}{1: struct{}{}, 2: struct{}{}, 3: struct{}{}},
					1:   map[interface{}]struct{}{"Z": struct{}{}},
					2:   map[interface{}]struct{}{"Z": struct{}{}},
					3:   map[interface{}]struct{}{"Z": struct{}{}},
				},
				as: map[interface{}]struct{}{"Z": struct{}{}},
				bs: map[interface{}]struct{}{1: struct{}{}, 2: struct{}{}, 3: struct{}{}},
			},
		},
		{
			"Z", New(),
		},
	} {
		if g.Remove(test.node); !reflect.DeepEqual(g, test.want) {
			t.Errorf("Remove(after %v): got %+v, want %+v", test.node, g, test.want)
		}
	}
}

func TestRemoveB(t *testing.T) {
	g := &Graph{
		m: map[interface{}]map[interface{}]struct{}{
			"X": map[interface{}]struct{}{0: struct{}{}},
			"Y": map[interface{}]struct{}{0: struct{}{}, 1: struct{}{}},
			"Z": map[interface{}]struct{}{1: struct{}{}, 2: struct{}{}, 3: struct{}{}},
			0:   map[interface{}]struct{}{"X": struct{}{}, "Y": struct{}{}},
			1:   map[interface{}]struct{}{"Y": struct{}{}, "Z": struct{}{}},
			2:   map[interface{}]struct{}{"Z": struct{}{}},
			3:   map[interface{}]struct{}{"Z": struct{}{}},
		},
		as: map[interface{}]struct{}{"X": struct{}{}, "Y": struct{}{}, "Z": struct{}{}},
		bs: map[interface{}]struct{}{0: struct{}{}, 1: struct{}{}, 2: struct{}{}, 3: struct{}{}},
	}
	for _, test := range []struct {
		node interface{}
		want *Graph
	}{
		{
			0,
			&Graph{
				m: map[interface{}]map[interface{}]struct{}{
					"Y": map[interface{}]struct{}{1: struct{}{}},
					"Z": map[interface{}]struct{}{1: struct{}{}, 2: struct{}{}, 3: struct{}{}},
					1:   map[interface{}]struct{}{"Y": struct{}{}, "Z": struct{}{}},
					2:   map[interface{}]struct{}{"Z": struct{}{}},
					3:   map[interface{}]struct{}{"Z": struct{}{}},
				},
				as: map[interface{}]struct{}{"Y": struct{}{}, "Z": struct{}{}},
				bs: map[interface{}]struct{}{1: struct{}{}, 2: struct{}{}, 3: struct{}{}},
			},
		},
		{
			1,
			&Graph{
				m: map[interface{}]map[interface{}]struct{}{
					"Z": map[interface{}]struct{}{2: struct{}{}, 3: struct{}{}},
					2:   map[interface{}]struct{}{"Z": struct{}{}},
					3:   map[interface{}]struct{}{"Z": struct{}{}},
				},
				as: map[interface{}]struct{}{"Z": struct{}{}},
				bs: map[interface{}]struct{}{2: struct{}{}, 3: struct{}{}},
			},
		},
		{
			2,
			&Graph{
				m: map[interface{}]map[interface{}]struct{}{
					"Z": map[interface{}]struct{}{3: struct{}{}},
					3:   map[interface{}]struct{}{"Z": struct{}{}},
				},
				as: map[interface{}]struct{}{"Z": struct{}{}},
				bs: map[interface{}]struct{}{3: struct{}{}},
			},
		},
		{
			3, New(),
		},
	} {
		if g.Remove(test.node); !reflect.DeepEqual(g, test.want) {
			t.Errorf("Remove(after %v): got %+v, want %+v", test.node, g, test.want)
		}
	}
}

func TestCopy(t *testing.T) {
	for _, g := range []*Graph{
		&Graph{
			m: map[interface{}]map[interface{}]struct{}{
				"X": map[interface{}]struct{}{0: struct{}{}},
				"Y": map[interface{}]struct{}{0: struct{}{}, 1: struct{}{}},
				"Z": map[interface{}]struct{}{1: struct{}{}, 2: struct{}{}, 3: struct{}{}},
				0:   map[interface{}]struct{}{"X": struct{}{}, "Y": struct{}{}},
				1:   map[interface{}]struct{}{"Y": struct{}{}, "Z": struct{}{}},
				2:   map[interface{}]struct{}{"Z": struct{}{}},
				3:   map[interface{}]struct{}{"Z": struct{}{}},
			},
			as: map[interface{}]struct{}{"X": struct{}{}, "Y": struct{}{}, "Z": struct{}{}},
			bs: map[interface{}]struct{}{0: struct{}{}, 1: struct{}{}, 2: struct{}{}, 3: struct{}{}},
		},
	} {
		got := Copy(g)
		if !reflect.DeepEqual(got, g) {
			t.Errorf("Copy(%+v): got %+v", g, got)
		}
	}
}

// match reports whether s0 and s1 contain the same elements, up to ordering.
func match(s0, s1 []interface{}) bool {
	if len(s0) != len(s1) {
		return false
	}
	c := append(make([]interface{}, 0, len(s1)), s1...)
L:
	for _, a := range s0 {
		for i := range c {
			if c[i] == a {
				c[i], c = c[len(c)-1], c[:len(c)-1]
				continue L
			}
		}
		return false
	}
	return true
}

func TestMatch(t *testing.T) {
	for _, test := range []struct {
		s0, s1 []interface{}
		want   bool
	}{
		{[]interface{}{}, []interface{}{}, true},
		{[]interface{}{"A"}, []interface{}{"A"}, true},
		{[]interface{}{"A", "B", "C"}, []interface{}{"A", "B", "C"}, true},
		{[]interface{}{"A", "B", "C"}, []interface{}{"B", "C", "A"}, true},
		{[]interface{}{"A", "B", "C"}, []interface{}{"A", "B", "C", "C"}, false},
		{[]interface{}{"A", "B", "B", "C"}, []interface{}{"A", "B", "C", "C"}, false},
		{[]interface{}{"A", "A", "B", "C"}, []interface{}{"A", "B", "A", "C"}, true},
	} {
		if got := match(test.s0, test.s1); got != test.want {
			t.Fatalf("match(%v, %v): got %v, want %v", test.s0, test.s1, got, test.want)
		}
	}
}

// wordsGraph returns a Graph of the 100 most common English verbs and the letters they contain.
func wordsGraph() *Graph {
	var dict = []string{
		"be", "have", "do", "say", "go", "get", "make", "know", "think", "take",
		"see", "come", "want", "look", "use", "find", "give", "tell", "work", "call",
		"try", "ask", "need", "feel", "become", "leave", "put", "mean", "keep", "let",
		"begin", "seem", "help", "talk", "turn", "start", "show", "hear", "play", "run",
		"move", "like", "live", "believe", "hold", "bring", "happen", "write", "provide", "sit",
		"stand", "lose", "pay", "meet", "include", "continue", "set", "learn", "change", "lead",
		"understand", "watch", "follow", "stop", "create", "speak", "read", "allow", "add", "spend",
		"grow", "open", "walk", "win", "offer", "remember", "love", "consider", "appear", "buy",
		"wait", "serve", "die", "send", "expect", "build", "stay", "fall", "cut", "reach",
		"kill", "remain", "suggest", "raise", "pass", "sell", "require", "report", "decide", "pull",
	}
	g := New()
	for _, s := range dict {
		for _, b := range s {
			g.Add(s, string(b))
		}
	}
	return g
}

func BenchmarkNew(b *testing.B) {
	for i := 0; i < b.N; i++ {
		New()
	}
}

func BenchmarkAdjacent(b *testing.B) {
	g := wordsGraph()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		g.Adjacent("require", "r")
	}
}

func BenchmarkAdd(b *testing.B) {
	g := wordsGraph()
	b.ResetTimer()
	for _, benchmark := range []struct{ name, a, b string }{
		{"create nodes", "test", "benchmark"},
		{"extant nodes", "require", "r"},
	} {
		b.Run(benchmark.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				b.StopTimer()
				g.Delete(benchmark.a, benchmark.b)
				b.StartTimer()
				g.Add(benchmark.a, benchmark.b)
			}
		})
	}
}

func BenchmarkDelete(b *testing.B) {
	g := wordsGraph()
	b.ResetTimer()
	for _, benchmark := range []struct{ name, a, b string }{
		{"retain nodes", "require", "r"},
		{"delete nodes", "test", "benchmark"},
	} {
		b.Run(benchmark.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				b.StopTimer()
				g.Add(benchmark.a, benchmark.b)
				b.StartTimer()
				g.Delete(benchmark.a, benchmark.b)
			}
		})
	}
}

func BenchmarkRemoveA(b *testing.B) {
	g := wordsGraph()
	bs := g.AdjTo("require")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		for _, b := range bs {
			g.Add("require", b)
		}
		b.StartTimer()
		g.Remove("require")
	}
}

func BenchmarkRemoveB(b *testing.B) {
	g := wordsGraph()
	as := g.AdjTo("r")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		for _, a := range as {
			g.Add(a, "r")
		}
		b.StartTimer()
		g.Remove("r")
	}
}

func BenchmarkAs(b *testing.B) {
	g := wordsGraph()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		g.As()
	}
}

func BenchmarkBs(b *testing.B) {
	g := wordsGraph()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		g.Bs()
	}
}

func BenchmarkAdjToA(b *testing.B) {
	g := wordsGraph()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		g.AdjTo("require")
	}
}

func BenchmarkAdjToB(b *testing.B) {
	g := wordsGraph()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		g.AdjTo("r")
	}
}

func BenchmarkDegA(b *testing.B) {
	g := wordsGraph()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		g.Deg("require")
	}
}

func BenchmarkDegB(b *testing.B) {
	g := wordsGraph()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		g.Deg("r")
	}
}

func BenchmarkNA(b *testing.B) {
	g := wordsGraph()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		g.NA()
	}
}

func BenchmarkNB(b *testing.B) {
	g := wordsGraph()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		g.NB()
	}
}

func BenchmarkCopy(b *testing.B) {
	g := wordsGraph()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Copy(g)
	}
}

func BenchmarkConstructDestruct(b *testing.B) {
	for i := 0; i < b.N; i++ {
		g := wordsGraph()
		for b := 'a'; b <= 'z'; b++ {
			g.Remove(string(b))
		}
	}
}
