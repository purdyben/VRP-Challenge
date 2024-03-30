package main

import (
	// "vorto/internal/mat"
	"fmt"
	"math"

	"gonum.org/v1/gonum/mat"
)

func main() {
	// adjacencyList := map[int][]int{
	// 	1: {2, 4},
	// 	2: {3, 5, 1},
	// 	3: {6, 2},
	// 	4: {1, 5, 7},
	// 	5: {2, 6, 8, 4},
	// 	6: {3, 0, 9, 5},
	// 	7: {4, 8},
	// 	8: {5, 9, 7},
	// 	9: {6, 0, 8},
	// }

	// g := mat.NewGraph(mat.WithAdjacencyList((adjacencyList)))

	// for _, v := range g.Vertices {
	// 	// fmt.Println(v)
	// 	s := fmt.Sprint(v.Val, ":")

	// 	for i, e := range v.Edges {
	// 		if i != 0 {
	// 			s += fmt.Sprint("-", e.Weight, "->")
	// 		}
	// 		s += fmt.Sprint(e.Vertex.Val)
	// 	}
	// 	fmt.Println(s)

	// }
	// a := mat.NewDense(4, 4, nil)
	// a.Set(1, 1, 9)

	// b := mat.NewDense(4, 4, nil)
	// a.Set(2, 2, 9)
	// a.Set(1, 1, 2)
	// a.Set(1, 2, 2)
	// b.Set(0, 1, 9)
	// b.Set(1, 1, 2)
	// b.Set(1, 2, 2)

	// var c mat.Dense
	// c.Add(a, b)

	// fm := mat.Formatted(&c, mat.Prefix("    "), mat.Squeeze())
	// fmt.Printf("m = %f\n", fm)

	// // Initialize two matrices, a and b.
	// a = mat.NewDense(2, 2, []float64{
	// 	1, 0,
	// 	1, 0,
	// })
	// b = mat.NewDense(2, 2, []float64{
	// 	0, 1,
	// 	0, 1,
	// })

	// // Add a and b, placing the result into c.
	// // Notice that the size is automatically adjusted
	// // when the receiver is empty (has zero size).
	// var t mat.Dense
	// t.Add(a, b)

	// // Print the result using the formatter.
	// fc := mat.Formatted(&t, mat.Prefix("    "), mat.Squeeze())
	// fmt.Printf("c = %v\n", fc)
	a := mat.NewDense(3, 3, []float64{
		1, 2, -1, 2, -3, -4, 1, 1, 1,
	})
	b := mat.NewDense(3, 1, []float64{
		7, -3, 0,
	})
	var inverseA mat.Dense
	err := inverseA.Inverse(a)
	fmt.Println(err)
	var res mat.Dense
	res.Mul(&inverseA, b)

	res.Apply(func(x, y int, val float64) float64 { return math.Round(val) }, &res)

	fm := mat.Formatted(a, mat.Prefix("    "), mat.Squeeze())
	fmt.Printf("m = %f\n", fm)
	fm = mat.Formatted(&inverseA, mat.Prefix("    "), mat.Squeeze())
	fmt.Printf("m = %f\n", fm)
	fm = mat.Formatted(&res, mat.Prefix("    "), mat.Squeeze())
	fmt.Printf("m = %f\n", fm)
}
