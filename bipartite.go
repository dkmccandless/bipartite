/*
Package bipartite implements an undirected bipartite graph data structure.

Add assigns a pair of nodes to separate sets a and b and connects them, creating either node
if necessary; Delete removes an edge between nodes if one exists, and removes either node
from the graph if the deleted edge was its last. Nodes can also be removed, and all of their
edges deleted, with Remove.
*/
package bipartite

// Graph is an undirected bipartite graph.
type Graph struct {
	m      map[interface{}]map[interface{}]struct{}
	as, bs map[interface{}]struct{}
}

// New returns an empty Graph ready to use.
func New() *Graph {
	return &Graph{
		m:  make(map[interface{}]map[interface{}]struct{}),
		as: make(map[interface{}]struct{}),
		bs: make(map[interface{}]struct{}),
	}
}

// Copy returns a pointer to a Graph that is deeply equal to g but shares no memory with it.
func Copy(g *Graph) *Graph {
	c := New()
	for k0 := range g.m {
		c.m[k0] = make(map[interface{}]struct{})
		for k1 := range g.m[k0] {
			c.m[k0][k1] = struct{}{}
		}
	}
	for k := range g.as {
		c.as[k] = struct{}{}
	}
	for k := range g.bs {
		c.bs[k] = struct{}{}
	}
	return c
}

// Add adds a and b to the graph if not present, and records that they are adjacent.
func (g *Graph) Add(a, b interface{}) {
	if _, ok := g.as[a]; !ok {
		g.as[a] = struct{}{}
		g.m[a] = make(map[interface{}]struct{})
	}
	if _, ok := g.bs[b]; !ok {
		g.bs[b] = struct{}{}
		g.m[b] = make(map[interface{}]struct{})
	}

	g.m[a][b] = struct{}{}
	g.m[b][a] = struct{}{}
}

// Adjacent reports whether a and b are adjacent.
func (g *Graph) Adjacent(a, b interface{}) bool {
	if _, ok := g.as[a]; !ok {
		return false
	}
	_, ok := g.m[a][b]
	return ok
}

// AdjTo returns an unordered slice of all nodes adjacent to node.
// If node is not in the graph, AdjTo returns nil.
func (g *Graph) AdjTo(node interface{}) []interface{} {
	m, ok := g.m[node]
	if !ok {
		return nil
	}
	s := make([]interface{}, 0, len(m))
	for a := range m {
		s = append(s, a)
	}
	return s
}

// As returns an unordered slice of all nodes added by Add into set a.
// If the graph is empty, As returns nil.
func (g *Graph) As() []interface{} {
	if len(g.as) == 0 {
		return nil
	}
	s := make([]interface{}, 0, len(g.as))
	for a := range g.as {
		s = append(s, a)
	}
	return s
}

// Bs returns an unordered slice of all nodes added by Add into set b.
// If the graph is empty, Bs returns nil.
func (g *Graph) Bs() []interface{} {
	if len(g.bs) == 0 {
		return nil
	}
	s := make([]interface{}, 0, len(g.bs))
	for b := range g.bs {
		s = append(s, b)
	}
	return s
}

// Deg returns the number of nodes adjacent to node.
// It is equivalent to len(g.AdjTo(node)), but faster.
func (g *Graph) Deg(node interface{}) int { return len(g.m[node]) }

// Delete records that a and b are not adjacent.
// If a node's last edge is deleted, Delete removes it from the graph.
func (g *Graph) Delete(a, b interface{}) {
	delete(g.m[a], b)
	if len(g.m[a]) == 0 {
		delete(g.m, a)
		delete(g.as, a)
	}
	delete(g.m[b], a)
	if len(g.m[b]) == 0 {
		delete(g.m, b)
		delete(g.bs, b)
	}
}

// NA returns the number of nodes added by Add into set a.
// It is equivalent to len(g.As()), but faster.
func (g *Graph) NA() int { return len(g.as) }

// NB returns the number of nodes added by Add into set b.
// It is equivalent to len(g.Bs()), but faster.
func (g *Graph) NB() int { return len(g.bs) }

// Remove deletes all of node's edges and removes it from the graph.
func (g *Graph) Remove(node interface{}) {
	for a := range g.m[node] {
		g.Delete(node, a)
		g.Delete(a, node)
	}
}
