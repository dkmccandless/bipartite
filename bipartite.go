/*
Package bipartite implements an undirected bipartite graph data structure.

Types A and B define disjoint and independent sets of nodes.

Add connects one node of each type with an edge, creating either node if necessary;
Delete removes an edge between nodes if one exists, and removes either node from the graph
if the deleted edge was its last. Nodes can also be removed, and all of their edges deleted,
with RemoveA and RemoveB.
*/
package bipartite

// A is a node that is only adjacent to B nodes.
// Values must be comparable using the == and != operators.
type A interface{}

// B is a node that is only adjacent to A nodes.
// Values must be comparable using the == and != operators.
type B interface{}

// Graph is an undirected bipartite graph of A and B nodes.
type Graph struct {
	ab map[A]map[B]bool
	ba map[B]map[A]bool
}

// New returns an empty Graph ready to use.
func New() *Graph {
	return &Graph{
		ab: make(map[A]map[B]bool),
		ba: make(map[B]map[A]bool),
	}
}

// Copy returns a pointer to a Graph that is deeply equal to g but shares no memory with it.
func Copy(g *Graph) *Graph {
	c := New()
	for a := range g.ab {
		c.ab[a] = make(map[B]bool)
		for b := range g.ab[a] {
			c.ab[a][b] = true
		}
	}
	for b := range g.ba {
		c.ba[b] = make(map[A]bool)
		for a := range g.ba[b] {
			c.ba[b][a] = true
		}
	}
	return c
}

// Add adds a and b to the graph if not present, and records that they are adjacent.
func (g *Graph) Add(a A, b B) {
	if _, ok := g.ab[a]; !ok {
		g.ab[a] = make(map[B]bool)
	}
	if _, ok := g.ba[b]; !ok {
		g.ba[b] = make(map[A]bool)
	}

	g.ab[a][b] = true
	g.ba[b][a] = true
}

// Adjacent reports whether a and b are adjacent.
func (g *Graph) Adjacent(a A, b B) bool {
	if _, ok := g.ab[a]; !ok {
		return false
	}
	return g.ab[a][b]
}

// AdjToA returns an unordered slice of all Bs adjacent to a.
// If a is not in the graph, AdjToA returns nil.
func (g *Graph) AdjToA(a A) []B {
	m, ok := g.ab[a]
	if !ok {
		return nil
	}
	s := make([]B, 0, len(m))
	for b := range m {
		s = append(s, b)
	}
	return s
}

// AdjToB returns an unordered slice of all As adjacent to b.
// If b is not in the graph, AdjToB returns nil.
func (g *Graph) AdjToB(b B) []A {
	m, ok := g.ba[b]
	if !ok {
		return nil
	}
	s := make([]A, 0, len(m))
	for a := range m {
		s = append(s, a)
	}
	return s
}

// As returns an unordered slice of all As in the graph.
// If the graph is empty, As returns nil.
func (g *Graph) As() []A {
	if len(g.ab) == 0 {
		return nil
	}
	s := make([]A, 0, len(g.ab))
	for a := range g.ab {
		s = append(s, a)
	}
	return s
}

// Bs returns an unordered slice of all Bs in the graph.
// If the graph is empty, Bs returns nil.
func (g *Graph) Bs() []B {
	if len(g.ba) == 0 {
		return nil
	}
	s := make([]B, 0, len(g.ba))
	for b := range g.ba {
		s = append(s, b)
	}
	return s
}

// DegA returns the number of Bs adjacent to a.
// It is equivalent to len(g.AdjToA(a)), but faster.
func (g *Graph) DegA(a A) int { return len(g.ab[a]) }

// DegB returns the number of As adjacent to b.
// It is equivalent to len(g.AdjToB(b)), but faster.
func (g *Graph) DegB(b B) int { return len(g.ba[b]) }

// Delete records that a and b are not adjacent.
// If a node's last edge is deleted, Delete removes it from the graph.
func (g *Graph) Delete(a A, b B) {
	delete(g.ab[a], b)
	if len(g.ab[a]) == 0 {
		delete(g.ab, a)
	}
	delete(g.ba[b], a)
	if len(g.ba[b]) == 0 {
		delete(g.ba, b)
	}
}

// NA returns the number of As in the graph.
// It is equivalent to len(g.As()), but faster.
func (g *Graph) NA() int { return len(g.ab) }

// NB returns the number of Bs in the graph.
// It is equivalent to len(g.Bs()), but faster.
func (g *Graph) NB() int { return len(g.ba) }

// RemoveA deletes all of a's edges and removes it from the graph.
func (g *Graph) RemoveA(a A) {
	for b := range g.ab[a] {
		g.Delete(a, b)
	}
}

// RemoveB deletes all of b's edges and removes it from the graph.
func (g *Graph) RemoveB(b B) {
	for a := range g.ba[b] {
		g.Delete(a, b)
	}
}
