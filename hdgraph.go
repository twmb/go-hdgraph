// Package hdgraph provides a very simple directed graph with the sole purpose
// of returning nodes and cycles in dependency order with the strongly
// connected components (Tarjan's) algorithm. This is exactly the same as
// go-dgraph, but backed by maps rather than slices.
//
// This package is only useful if you will be repeatedly removing or unlinking
// nodes in the graph and re-running the StrongComponents algorithm. Maps are
// much slower than slices for running the algorithm, but to unlink or remove
// nodes for maps is O(1) and O(V) compared to O(V) and O(V^2) for slices.
//
// The StrongComponents returns graph components in dependency order. If the
// graph has no cycles, each component will have a single element. Otherwise,
// all nodes in a cycle are grouped into one "strong" component.
package hdgraph

// Graph is a directed graph.
type Graph struct {
	out map[int]map[int]struct{}
	in  map[int]map[int]struct{}
}

// New returns a new graph with no nodes.
func New() *Graph {
	return &Graph{
		out: make(map[int]map[int]struct{}),
		in:  make(map[int]map[int]struct{}),
	}
}

// Add adds a node to the graph if it does not exist.
func (g *Graph) Add(node int) {
	if _, exists := g.out[node]; !exists {
		g.out[node] = make(map[int]struct{})
		g.in[node] = make(map[int]struct{})
	}
}

// Remove removes a node from the graph if it exists.
func (g *Graph) Remove(node int) {
	out := g.out[node]
	for dst := range out {
		delete(g.in[dst], node)
	}
	delete(g.out, node)
}

// Link adds an edge from src to dst, creating the nodes if they do not exist.
func (g *Graph) Link(src, dst int) {
	g.Add(src)
	g.Add(dst)
	g.out[src][dst] = struct{}{}
	g.in[dst][src] = struct{}{}
}

// Unlink removes an edge from src to dst if it exists.
func (g *Graph) Unlink(src, dst int) {
	if out, exists := g.out[src]; exists {
		delete(out, dst)
	}
	if in, exists := g.in[dst]; exists {
		delete(in, src)
	}
}

// StrongComponents returns all strong components of the graph in dependency
// order. If a graph has no cycles, this will return each node individually.
//
// Note that this is a topological sort: if a graph has no cycles, the order of
// the returned single element components is the dependency order of the graph.
func (g *Graph) StrongComponents() [][]int {
	d := dfser{
		graph: g.out,
		order: make([]int, 0, len(g.out)),
		flip:  make([]bool, len(g.out)),
	}

	for node := range g.out {
		if !d.sawTrue(node) {
			d.dfs1(node)
		}
	}

	order := d.order

	d.order = make([]int, 0, len(order))
	d.graph = g.in

	components := make([][]int, 0, len(order))
	for i := len(order) - 1; i >= 0; i-- {
		if node := order[i]; !d.sawFalse(node) {
			d.dfs2(node)
			used := len(d.order)
			component := d.order[:used:used]
			components = append(components, component)
			d.order = d.order[used:]
		}
	}

	return components
}

// dfser performs depth first search on the graph, appending bottom nodes to
// order.
//
// Since SCC requires two dfs runs, we keep track of nodes seen with flip: on
// the first run, we swap flip nodes to true, on the second back to false.
//
// The original impl just passed graph and two closures to one dfs function,
// but switching to methods on a struct proved a big (>30%) performance
// increase.
type dfser struct {
	graph map[int]map[int]struct{}
	order []int
	flip  []bool
}

func (d *dfser) sawTrue(node int) bool {
	r := d.flip[node]
	d.flip[node] = true
	return r
}
func (d *dfser) sawFalse(node int) bool {
	r := d.flip[node]
	d.flip[node] = false
	return !r
}

func (d *dfser) dfs1(node int) {
	for neighbor := range d.graph[node] {
		if !d.sawTrue(neighbor) {
			d.dfs1(neighbor)
		}
	}
	d.order = append(d.order, node)
}

func (d *dfser) dfs2(node int) {
	for neighbor := range d.graph[node] {
		if !d.sawFalse(neighbor) {
			d.dfs2(neighbor)
		}
	}
	d.order = append(d.order, node)
}
