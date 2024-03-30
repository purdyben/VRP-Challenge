package mat

// type AdjacencyMatrix[T any] struct {
// 	m [][]T
// }

// func (am *AdjacencyMatrix[T]) GetVerts(index int) []T {
// 	return nil
// }

// // AddVertex
// // AddEdge
// // Neighbors
// // zero := mat.NewDense(3, 5, nil)

// Graph represents a set of vertices connected by edges.
type Graph struct {
	Vertices map[int]*Vertex
}

type GraphOption func(this *Graph)

// Vertex is a node in the graph that stores the int value at that node
// along with a map to the vertices it is connected to via edges.
type Vertex struct {
	Val   int
	Edges map[int]*Edge
}

// Edge represents an edge in the graph and the destination vertex.
type Edge struct {
	Weight int
	Vertex *Vertex
}

func NewGraph(opts ...GraphOption) *Graph {
	g := &Graph{Vertices: map[int]*Vertex{}}

	for _, opt := range opts {
		opt(g)
	}

	return g
}

func WithAdjacencyList(list map[int][]int) GraphOption {
	return func(this *Graph) {
		for vertex, edges := range list {
			// add vertex
			if _, ok := this.Vertices[vertex]; !ok {
				this.AddVertex(vertex, vertex)
			}

			// add edges to vertex
			for _, edge := range edges {
				// add edge as vertex, if not added
				if _, ok := this.Vertices[edge]; !ok {
					this.AddVertex(edge, edge)
				}

				this.AddEdge(vertex, edge, 0) // no weights in this adjacency list
			}
		}
	}
}

func (this *Graph) AddVertex(key, val int) {
	this.Vertices[key] = &Vertex{Val: val, Edges: map[int]*Edge{}}
}

func (this *Graph) AddEdge(srcKey, destKey int, weight int) {
	// check if src & dest exist
	if _, ok := this.Vertices[srcKey]; !ok {
		return
	}
	if _, ok := this.Vertices[destKey]; !ok {
		return
	}

	// add edge src --> dest
	this.Vertices[srcKey].Edges[destKey] = &Edge{Weight: weight, Vertex: this.Vertices[destKey]}
}

func (this *Graph) Neighbors(srcKey int) []int {
	result := []int{}

	for _, edge := range this.Vertices[srcKey].Edges {
		result = append(result, edge.Vertex.Val)
	}

	return result
}
