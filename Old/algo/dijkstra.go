package algo

// import (
// 	"fmt"
// )

// type Node interface {
// 	ID() string
// 	Data() string
// }

// type ItemGraph interface {
// 	GetEdges(node string) []*Edge
// }

// type Edge struct {
// 	Node   Node
// 	Weight int
// }

// type Vertex struct {
// 	Node     Node
// 	Distance int
// }

// // type ItemGraph struct {
// // 	Nodes []Node
// // 	Edges map[Node][]*Edge
// // 	Lock  sync.RWMutex
// // }

// type InputGraph struct {
// 	Graph []InputData `json:"graph"`
// 	From  string      `json:"from"`
// 	To    string      `json:"to"`
// }

// type InputData struct {
// 	Source      string `json:"source"`
// 	Destination string `json:"destination"`
// 	Weight      int    `json:"weight"`
// }

// type dykstraMemory struct {
// 	visited             map[string]bool
// 	dist                map[string]int
// 	prev                map[string]string
// 	startNode           Node
// 	endNode             Node
// 	graph               ItemGraph
// 	KnownShortestWeight int
// 	KnownShortestPath   []string
// 	pq                  *Queue[Vertex]
// }

// func DykstraShortestPath(startNode Node, endNode Node, g ItemGraph) ([]string, int) {
// 	dm := dykstraMemory{
// 		graph:     g,
// 		startNode: startNode,
// 		endNode:   endNode,
// 		visited:   make(map[string]bool),
// 		dist:      make(map[string]int),
// 		prev:      make(map[string]string),
// 		pq:        NewQueue[Vertex](),
// 	}
// 	start := Vertex{
// 		Node:     startNode,
// 		Distance: 0,
// 	}
// 	return dykstraShortestPath(dm, start)
// }

// func dykstraShortestPath(dm dykstraMemory, start Vertex) ([]string, int) {
// 	// for _, nval := range g.Nodes {
// 	// 	dist[nval.Data()] = math.MaxInt64
// 	// }
// 	// dist[startNode.Data()] = start.Distance
// 	dm.pq.Enqueue(start)
// 	for !dm.pq.IsEmpty() {
// 		currNode, err := dm.pq.Dequeue()
// 		if err != nil {
// 			continue
// 		}
// 		if _, ok := dm.visited[currNode.Node.ID()]; ok {
// 			continue
// 		}

// 		dm.visited[currNode.Node.ID()] = true

// 		near := dm.graph.GetEdges(currNode.Node.ID())

// 		for _, val := range near {
// 			if !dm.visited[val.Node.ID()] {
// 				weightToNextNode := dm.dist[currNode.Node.ID()] + val.Weight
// 				// is val.Node.ID() in dist
// 				if weightToNextNode < dm.dist[val.Node.ID()] {
// 					store := Vertex{
// 						Node:     val.Node,
// 						Distance: weightToNextNode,
// 					}
// 					// weight to val is prev weight its self
// 					dm.dist[val.Node.ID()] = dm.dist[currNode.Node.ID()] + val.Weight
// 					// prev[val.Node.Data()] = fmt.Sprintf("->%s", v.Node.Data())
// 					dm.prev[val.Node.ID()] = currNode.Node.ID()
// 					dm.pq.Enqueue(store)
// 				}
// 			}
// 		}
// 	}
// 	fmt.Println(dm.dist)
// 	fmt.Println(dm.prev)
// 	endNode := dm.endNode
// 	startNode := dm.startNode

// 	pathval := dm.prev[endNode.ID()]

// 	var finalArr []string

// 	finalArr = append(finalArr, endNode.Data())
// 	for pathval != startNode.Data() {
// 		finalArr = append(finalArr, pathval)
// 		pathval = dm.prev[pathval]
// 	}
// 	finalArr = append(finalArr, pathval)
// 	fmt.Println(finalArr)
// 	for i, j := 0, len(finalArr)-1; i < j; i, j = i+1, j-1 {
// 		finalArr[i], finalArr[j] = finalArr[j], finalArr[i]
// 	}
// 	return finalArr, dm.dist[endNode.Data()]
// }
