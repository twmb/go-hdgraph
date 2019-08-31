go-hdgraph
=========

Package hdgraph provides a very simple directed graph with the sole purpose
of returning nodes and cycles in dependency order with the strongly
connected components (Tarjan's) algorithm. This is exactly the same as
go-dgraph, but backed by maps rather than slices.

This package is only useful if you will be repeatedly removing or unlinking
nodes in the graph and re-running the StrongComponents algorithm. Maps are
much slower than slices for running the algorithm, but to unlink or remove
nodes for maps is O(1) and O(V) compared to O(V) and O(V^2) for slices.

The StrongComponents returns graph components in dependency order. If the
graph has no cycles, each component will have a single element. Otherwise,
all nodes in a cycle are grouped into one "strong" component.

Documentation
-------------

[![GoDoc](https://godoc.org/github.com/twmb/go-hdgraph?status.svg)](https://godoc.org/github.com/twmb/go-hdgraph)
