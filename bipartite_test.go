package bipartite

import (
	"reflect"
	"testing"
)

func TestN(t *testing.T) {
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
	na, nb := 3, 4
	if n := g.NA(); n != na {
		t.Errorf("NA(%+v): got %v, want %v", g, n, na)
	}
	if n := g.NB(); n != nb {
		t.Errorf("NB(%+v): got %v, want %v", g, n, nb)
	}
}

func TestDeg(t *testing.T) {
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
	adja := map[A][]B{"X": []B{0}, "Y": []B{0, 1}, "Z": []B{1, 2, 3}}
	adjb := map[B][]A{0: []A{"X", "Y"}, 1: []A{"Y", "Z"}, 2: []A{"Z"}, 3: []A{"Z"}}
	for a, adj := range adja {
		if _, ok := g.ab[a]; !ok {
			t.Fatalf("DegA: %+v does not contain %v", g, a)
		}
		if deg := g.DegA(a); deg != len(adj) {
			t.Errorf("DegA(%+v, %v): got %v, want %v", g, a, deg, len(adj))
		}
	}
	for b, adj := range adjb {
		if _, ok := g.ba[b]; !ok {
			t.Fatalf("DegB: %+v does not contain %v", g, b)
		}
		if deg := g.DegB(b); deg != len(adj) {
			t.Errorf("DegB(%+v, %v): got %v, want %v", g, b, deg, len(adj))
		}
	}
}

func TestAdjTo(t *testing.T) {
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
	adja := map[A][]B{"X": []B{0}, "Y": []B{0, 1}, "Z": []B{1, 2, 3}}
	adjb := map[B][]A{0: []A{"X", "Y"}, 1: []A{"Y", "Z"}, 2: []A{"Z"}, 3: []A{"Z"}}
	for a, adj := range adja {
		if _, ok := g.ab[a]; !ok {
			t.Fatalf("AdjToA: %+v does not contain %v", g, a)
		}
		got := g.AdjToA(a)
		if !matchB(got, adj) {
			t.Errorf("AdjToA(%+v, %v): got %v, want %v", g, a, got, adj)
		}
	}
	for b, adj := range adjb {
		if _, ok := g.ba[b]; !ok {
			t.Fatalf("AdjToB: %+v does not contain %v", g, b)
		}
		got := g.AdjToB(b)
		if !matchA(got, adj) {
			t.Errorf("AdjToB(%+v, %v): got %v, want %v", g, b, got, adj)
		}
	}
}

func TestAdjacent(t *testing.T) {
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
	adj := map[A]map[B]bool{
		"W": map[B]bool{0: false, 1: false, 2: false, 3: false},
		"X": map[B]bool{0: true, 1: false, 2: false, 3: false},
		"Y": map[B]bool{0: true, 1: true, 2: false, 3: false},
		"Z": map[B]bool{0: false, 1: true, 2: true, 3: true},
	}
	for a := range adj {
		for b := range adj[a] {
			if got := g.Adjacent(a, b); got != adj[a][b] {
				t.Errorf("Adjacent(%+v: %v, %v): got %v, want %v", g, a, b, got, adj[a][b])
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