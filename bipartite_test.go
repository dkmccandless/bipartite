package bipartite

import (
	"reflect"
	"testing"
)

var adjTests = []struct {
	g    *Graph
	adj  map[A]map[B]bool
	adja map[A][]B
	adjb map[B][]A
}{
	{
		g: New(),
	},
	{
		g: &Graph{
			ab: map[A]map[B]struct{}{"apple": map[B]struct{}{"tree": struct{}{}}},
			ba: map[B]map[A]struct{}{"tree": map[A]struct{}{"apple": struct{}{}}},
		},
		adj: map[A]map[B]bool{
			"apple":     map[B]bool{"tree": true, "rock": false},
			"spaghetti": map[B]bool{"tree": false},
		},
		adja: map[A][]B{"apple": []B{"tree"}},
		adjb: map[B][]A{"tree": []A{"apple"}},
	},
	{
		g: &Graph{
			ab: map[A]map[B]struct{}{
				"X": map[B]struct{}{0: struct{}{}},
				"Y": map[B]struct{}{0: struct{}{}, 1: struct{}{}},
				"Z": map[B]struct{}{1: struct{}{}, 2: struct{}{}, 3: struct{}{}},
			},
			ba: map[B]map[A]struct{}{
				0: map[A]struct{}{"X": struct{}{}, "Y": struct{}{}},
				1: map[A]struct{}{"Y": struct{}{}, "Z": struct{}{}},
				2: map[A]struct{}{"Z": struct{}{}},
				3: map[A]struct{}{"Z": struct{}{}},
			},
		},
		adj: map[A]map[B]bool{
			"W": map[B]bool{0: false, 1: false, 2: false, 3: false},
			"X": map[B]bool{0: true, 1: false, 2: false, 3: false},
			"Y": map[B]bool{0: true, 1: true, 2: false, 3: false},
			"Z": map[B]bool{0: false, 1: true, 2: true, 3: true},
		},
		adja: map[A][]B{"X": []B{0}, "Y": []B{0, 1}, "Z": []B{1, 2, 3}},
		adjb: map[B][]A{0: []A{"X", "Y"}, 1: []A{"Y", "Z"}, 2: []A{"Z"}, 3: []A{"Z"}},
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
		for a, adj := range test.adja {
			if _, ok := test.g.ab[a]; !ok {
				t.Fatalf("DegA: %+v does not contain %v", test.g, a)
			}
			if deg := test.g.DegA(a); deg != len(adj) {
				t.Errorf("DegA(%+v, %v): got %v, want %v", test.g, a, deg, len(adj))
			}
			if got := test.g.AdjToA(a); !matchB(got, adj) {
				t.Errorf("AdjToA(%+v, %v): got %v, want %v", test.g, a, got, adj)
			}
		}
		for b, adj := range test.adjb {
			if _, ok := test.g.ba[b]; !ok {
				t.Fatalf("DegB: %+v does not contain %v", test.g, b)
			}
			if deg := test.g.DegB(b); deg != len(adj) {
				t.Errorf("DegB(%+v, %v): got %v, want %v", test.g, b, deg, len(adj))
			}
			if got := test.g.AdjToB(b); !matchA(got, adj) {
				t.Errorf("AdjToB(%+v, %v): got %v, want %v", test.g, b, got, adj)
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
		a          A
		b          B
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
		if degA := g.DegA(test.a); degA != test.degA {
			t.Errorf("Add(%+v after %v, %v): A degree %v, want %v", g, test.a, test.b, degA, test.degA)
		}
		if degB := g.DegB(test.b); degB != test.degB {
			t.Errorf("Add(%+v after %v, %v): B degree %v, want %v", g, test.a, test.b, degB, test.degB)
		}
	}
}

func TestDelete(t *testing.T) {
	g := &Graph{
		ab: map[A]map[B]struct{}{
			"X": map[B]struct{}{0: struct{}{}},
			"Y": map[B]struct{}{0: struct{}{}, 1: struct{}{}},
			"Z": map[B]struct{}{1: struct{}{}, 2: struct{}{}, 3: struct{}{}},
		},
		ba: map[B]map[A]struct{}{
			0: map[A]struct{}{"X": struct{}{}, "Y": struct{}{}},
			1: map[A]struct{}{"Y": struct{}{}, "Z": struct{}{}},
			2: map[A]struct{}{"Z": struct{}{}},
			3: map[A]struct{}{"Z": struct{}{}},
		},
	}
	tests := []struct {
		a          A
		b          B
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
		if degA := g.DegA(test.a); degA != test.degA {
			t.Errorf("Delete(%+v after %v, %v): got %v degree %v, want %v", g, test.a, test.b, test.a, degA, test.degA)
		}
		if degB := g.DegB(test.b); degB != test.degB {
			t.Errorf("Delete(%+v after %v, %v): got %v degree %v, want %v", g, test.a, test.b, test.b, degB, test.degB)
		}
	}
}

func TestRemoveA(t *testing.T) {
	g := &Graph{
		ab: map[A]map[B]struct{}{
			"X": map[B]struct{}{0: struct{}{}},
			"Y": map[B]struct{}{0: struct{}{}, 1: struct{}{}},
			"Z": map[B]struct{}{1: struct{}{}, 2: struct{}{}, 3: struct{}{}},
		},
		ba: map[B]map[A]struct{}{
			0: map[A]struct{}{"X": struct{}{}, "Y": struct{}{}},
			1: map[A]struct{}{"Y": struct{}{}, "Z": struct{}{}},
			2: map[A]struct{}{"Z": struct{}{}},
			3: map[A]struct{}{"Z": struct{}{}},
		},
	}
	for _, test := range []struct {
		a    A
		want *Graph
	}{
		{
			"X",
			&Graph{
				ab: map[A]map[B]struct{}{
					"Y": map[B]struct{}{0: struct{}{}, 1: struct{}{}},
					"Z": map[B]struct{}{1: struct{}{}, 2: struct{}{}, 3: struct{}{}},
				},
				ba: map[B]map[A]struct{}{
					0: map[A]struct{}{"Y": struct{}{}},
					1: map[A]struct{}{"Y": struct{}{}, "Z": struct{}{}},
					2: map[A]struct{}{"Z": struct{}{}},
					3: map[A]struct{}{"Z": struct{}{}},
				},
			},
		},
		{
			"Y",
			&Graph{
				ab: map[A]map[B]struct{}{
					"Z": map[B]struct{}{1: struct{}{}, 2: struct{}{}, 3: struct{}{}},
				},
				ba: map[B]map[A]struct{}{
					1: map[A]struct{}{"Z": struct{}{}},
					2: map[A]struct{}{"Z": struct{}{}},
					3: map[A]struct{}{"Z": struct{}{}},
				},
			},
		},
		{
			"Z", New(),
		},
	} {
		if g.RemoveA(test.a); !reflect.DeepEqual(g, test.want) {
			t.Errorf("RemoveA(after %v): got %+v, want %+v", test.a, g, test.want)
		}
	}
}

func TestRemoveB(t *testing.T) {
	g := &Graph{
		ab: map[A]map[B]struct{}{
			"X": map[B]struct{}{0: struct{}{}},
			"Y": map[B]struct{}{0: struct{}{}, 1: struct{}{}},
			"Z": map[B]struct{}{1: struct{}{}, 2: struct{}{}, 3: struct{}{}},
		},
		ba: map[B]map[A]struct{}{
			0: map[A]struct{}{"X": struct{}{}, "Y": struct{}{}},
			1: map[A]struct{}{"Y": struct{}{}, "Z": struct{}{}},
			2: map[A]struct{}{"Z": struct{}{}},
			3: map[A]struct{}{"Z": struct{}{}},
		},
	}
	for _, test := range []struct {
		b    B
		want *Graph
	}{
		{
			0,
			&Graph{
				ab: map[A]map[B]struct{}{
					"Y": map[B]struct{}{1: struct{}{}},
					"Z": map[B]struct{}{1: struct{}{}, 2: struct{}{}, 3: struct{}{}},
				},
				ba: map[B]map[A]struct{}{
					1: map[A]struct{}{"Y": struct{}{}, "Z": struct{}{}},
					2: map[A]struct{}{"Z": struct{}{}},
					3: map[A]struct{}{"Z": struct{}{}},
				},
			},
		},
		{
			1,
			&Graph{
				ab: map[A]map[B]struct{}{
					"Z": map[B]struct{}{2: struct{}{}, 3: struct{}{}},
				},
				ba: map[B]map[A]struct{}{
					2: map[A]struct{}{"Z": struct{}{}},
					3: map[A]struct{}{"Z": struct{}{}},
				},
			},
		},
		{
			2,
			&Graph{
				ab: map[A]map[B]struct{}{
					"Z": map[B]struct{}{3: struct{}{}},
				},
				ba: map[B]map[A]struct{}{
					3: map[A]struct{}{"Z": struct{}{}},
				},
			},
		},
		{
			3, New(),
		},
	} {
		if g.RemoveB(test.b); !reflect.DeepEqual(g, test.want) {
			t.Errorf("RemoveB(after %v): got %+v, want %+v", test.b, g, test.want)
		}
	}
}

func TestCopy(t *testing.T) {
	for _, g := range []*Graph{
		&Graph{
			ab: map[A]map[B]struct{}{
				"X": map[B]struct{}{0: struct{}{}},
				"Y": map[B]struct{}{0: struct{}{}, 1: struct{}{}},
				"Z": map[B]struct{}{1: struct{}{}, 2: struct{}{}, 3: struct{}{}},
			},
			ba: map[B]map[A]struct{}{
				0: map[A]struct{}{"X": struct{}{}, "Y": struct{}{}},
				1: map[A]struct{}{"Y": struct{}{}, "Z": struct{}{}},
				2: map[A]struct{}{"Z": struct{}{}},
				3: map[A]struct{}{"Z": struct{}{}},
			},
		},
	} {
		got := Copy(g)
		if !reflect.DeepEqual(got, g) {
			t.Errorf("Copy(%+v): got %+v", g, got)
		}
	}
}

// matchA reports whether s0 and s1 contain the same elements, up to ordering.
func matchA(s0, s1 []A) bool {
	if len(s0) != len(s1) {
		return false
	}
	c := append(make([]A, 0, len(s1)), s1...)
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

// matchB reports whether s0 and s1 contain the same elements, up to ordering.
func matchB(s0, s1 []B) bool {
	if len(s0) != len(s1) {
		return false
	}
	c := append(make([]B, 0, len(s1)), s1...)
L:
	for _, b := range s0 {
		for i := range c {
			if c[i] == b {
				c[i], c = c[len(c)-1], c[:len(c)-1]
				continue L
			}
		}
		return false
	}
	return true
}

func TestMatchA(t *testing.T) {
	for _, test := range []struct {
		s0, s1 []A
		want   bool
	}{
		{[]A{}, []A{}, true},
		{[]A{"A"}, []A{"A"}, true},
		{[]A{"A", "B", "C"}, []A{"A", "B", "C"}, true},
		{[]A{"A", "B", "C"}, []A{"B", "C", "A"}, true},
		{[]A{"A", "B", "C"}, []A{"A", "B", "C", "C"}, false},
		{[]A{"A", "B", "B", "C"}, []A{"A", "B", "C", "C"}, false},
		{[]A{"A", "A", "B", "C"}, []A{"A", "B", "A", "C"}, true},
	} {
		if got := matchA(test.s0, test.s1); got != test.want {
			t.Fatalf("matchA(%v, %v): got %v, want %v", test.s0, test.s1, got, test.want)
		}
	}
}

func TestMatchB(t *testing.T) {
	for _, test := range []struct {
		s0, s1 []B
		want   bool
	}{
		{[]B{}, []B{}, true},
		{[]B{"A"}, []B{"A"}, true},
		{[]B{"A", "B", "C"}, []B{"A", "B", "C"}, true},
		{[]B{"A", "B", "C"}, []B{"B", "C", "A"}, true},
		{[]B{"A", "B", "C"}, []B{"A", "B", "C", "C"}, false},
		{[]B{"A", "B", "B", "C"}, []B{"A", "B", "C", "C"}, false},
		{[]B{"A", "A", "B", "C"}, []B{"A", "B", "A", "C"}, true},
	} {
		if got := matchB(test.s0, test.s1); got != test.want {
			t.Fatalf("matchB(%v, %v): got %v, want %v", test.s0, test.s1, got, test.want)
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
	bs := g.AdjToA("require")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		for _, b := range bs {
			g.Add("require", b)
		}
		b.StartTimer()
		g.RemoveA("require")
	}
}

func BenchmarkRemoveB(b *testing.B) {
	g := wordsGraph()
	as := g.AdjToB("r")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		for _, a := range as {
			g.Add(a, "r")
		}
		b.StartTimer()
		g.RemoveB("r")
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
		g.AdjToA("require")
	}
}

func BenchmarkAdjToB(b *testing.B) {
	g := wordsGraph()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		g.AdjToB("r")
	}
}

func BenchmarkDegA(b *testing.B) {
	g := wordsGraph()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		g.DegA("require")
	}
}

func BenchmarkDegB(b *testing.B) {
	g := wordsGraph()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		g.DegB("r")
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
			g.RemoveB(string(b))
		}
	}
}
