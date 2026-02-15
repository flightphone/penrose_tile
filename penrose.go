package main

import (
	"fmt"
	"math"

	"github.com/fogleman/gg"
)

var phi = (math.Sqrt(5) - 1) / 2
var baseIndex = 27

type Shape interface {
	split() []Shape
	Draw(width float64, height float64, dc *gg.Context, lens map[int]int, graph map[string]*PenVertex,
		ra float64)
	getLink() (ia string, ib string)
}

func split(tris []Shape) []Shape {
	res := []Shape{}
	for _, tri := range tris {
		res = append(res, tri.split()...)
	}
	return res
}

type PenVertex struct {
	isMark   bool
	index    int
	children []string
}

// HSLToRGB converts an HSL triple to an RGB triple.
// https://github.com/crazy3lf/colorconv/blob/v1.2.0/colorconv.go
func HSLToRGB(h, s, l float64) (r, g, b float64) {
	if h < 0 || h >= 360 ||
		s < 0 || s > 1 ||
		l < 0 || l > 1 {
		return 0, 0, 0
	}
	// When 0 ≤ h < 360, 0 ≤ s ≤ 1 and 0 ≤ l ≤ 1:
	C := (1 - math.Abs((2*l)-1)) * s
	X := C * (1 - math.Abs(math.Mod(h/60, 2)-1))
	m := l - (C / 2)
	var Rnot, Gnot, Bnot float64

	switch {
	case 0 <= h && h < 60:
		Rnot, Gnot, Bnot = C, X, 0
	case 60 <= h && h < 120:
		Rnot, Gnot, Bnot = X, C, 0
	case 120 <= h && h < 180:
		Rnot, Gnot, Bnot = 0, C, X
	case 180 <= h && h < 240:
		Rnot, Gnot, Bnot = 0, X, C
	case 240 <= h && h < 300:
		Rnot, Gnot, Bnot = X, 0, C
	case 300 <= h && h < 360:
		Rnot, Gnot, Bnot = C, 0, X
	}
	r = Rnot + m
	g = Gnot + m
	b = Bnot + m
	return r, g, b
}

func getPointKey(p complex128) string {
	return fmt.Sprintf("%.4f,%.4f", real(p), imag(p))
}

// penrose
// h - heights canvas
// tris - init triangles
// ra - length
// n - count iteration
func penrose(h float64, tris []Shape, ra float64, n int, filename string, graph_color bool) {
	width := h
	height := h

	for range n {
		tris = split(tris)
		ra = ra * phi
	}

	var graph map[string]*PenVertex = make(map[string]*PenVertex)
	var lens map[int]int = make(map[int]int) //Длина области по номеру области

	// Строим граф
	if graph_color {
		for _, sap := range tris {
			// Распаковываем интерфейс в конкретный тип
			//tri, _ := sap.(*TriangleP3)
			ia, ib := sap.getLink()
			var va *PenVertex
			var vb *PenVertex
			var ok bool

			va, ok = graph[ia]
			if !ok {
				va = &PenVertex{isMark: false, index: 0}
				graph[ia] = va
			}

			vb, ok = graph[ib]
			if !ok {
				vb = &PenVertex{isMark: false, index: 0}
				graph[ib] = vb
			}

			va.children = append(va.children, ib)
			vb.children = append(vb.children, ia)
		}

		//запускаем BFS по графу для раскраски областей связности

		num := 0 //Номер связности
		for key, vertex := range graph {
			if vertex.isMark {
				continue
			} else {
				bfs := []string{key}
				num += 1 //Увеличивыем номер связной области

				for i := 0; i < len(bfs); i++ {
					ix := bfs[i]
					graph[ix].isMark = true
					graph[ix].index = num
					for _, val := range graph[ix].children {
						if !graph[val].isMark {
							bfs = append(bfs, val)
						}
					}
				}
				lens[num] = len(bfs) //Сохранили длину
			}
		}
	}

	//Draw
	dc := gg.NewContext(int(width), int(height))
	dc.SetRGB(1, 1, 1)
	dc.DrawRectangle(0, 0, width, height)
	dc.Fill()
	for _, sap := range tris {
		sap.Draw(width, height, dc, lens, graph, ra)
	}
	dc.SavePNG(filename)
}
