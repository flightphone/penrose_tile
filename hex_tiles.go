package main

import (
	"math"
	"math/rand"

	"github.com/fogleman/gg"
)

//https://www.gonum.org/

//"gonum.org/v1/gonum/spatial/r2"

// Индекс вершины графа (середина стороны шестиугольника) e = 0..5
type VerIndex struct {
	i int
	j int
	e int
}

type Color struct {
	R float64
	G float64
	B float64
}

// Вершина графа, children - связанные вершины
type Vertex struct {
	isMark   bool
	color    Color
	children []VerIndex
}

type Hexagon struct {
	i        int
	j        int
	typeline int
	rotate   int
	lines    []VerIndex
}

// структура для вычисления глобального индекса вершины
var shift []VerIndex = []VerIndex{
	{i: 0, j: 0, e: 0},
	{i: 0, j: 0, e: 1},
	{i: 0, j: 0, e: 2},
	{i: -1, j: 0, e: 0},
	{i: -1, j: -1, e: 1},
	{i: 0, j: -1, e: 2},
}

// Возвращает глобальный индекс  e = 0,1,2
func getIndex(a VerIndex) VerIndex {

	s := shift[a.e]
	res := VerIndex{i: a.i + s.i, j: a.j + s.j, e: s.e}
	if a.e > 3 {
		res.i = res.i + (a.j)%2
	}
	return res

}

// построение графа
func link(i int, j int, rotate int, lines [][]int, graph map[VerIndex]*Vertex) []VerIndex {
	res := []VerIndex{}
	for _, line := range lines {
		a := line[0]
		b := line[1]
		a = (a + rotate) % 6
		b = (b + rotate) % 6
		ia := getIndex(VerIndex{i: i, j: j, e: a})
		ib := getIndex(VerIndex{i: i, j: j, e: b})
		res = append(res, ia)

		var va *Vertex
		var vb *Vertex
		var ok bool

		va, ok = graph[ia]
		if !ok {
			va = &Vertex{isMark: false, color: Color{R: 0, G: 0, B: 0}}
			graph[ia] = va
		}
		vb, ok = graph[ib]
		if !ok {
			vb = &Vertex{isMark: false, color: Color{R: 0, G: 0, B: 0}}
			graph[ib] = vb
		}
		va.children = append(va.children, ib)
		vb.children = append(vb.children, ia)
	}
	return res
}

func hex_tiles() {

	//Создаем граф
	//map для хранения графа
	var graph map[VerIndex]*Vertex = make(map[VerIndex]*Vertex)

	// связи между серединами сторон шестиугольника
	var lines0 [][]int = [][]int{{4, 5}, {0, 1}, {2, 3}}
	var lines1 [][]int = [][]int{{1, 2}, {4, 5}, {0, 3}}

	var width float64 = 3200.0
	var height float64 = 3200.0
	var r float64 = 80
	dx := 2.0 * r * math.Cos(math.Pi/6)
	dy := r + r*math.Sin(math.Pi/6)
	n := int(width/dx) + 2
	m := int(height/dy) + 1
	var Hexs []Hexagon
	for j := range m {
		for i := range n {
			tp := rand.Intn(2)
			var al int
			var lines []VerIndex
			if tp == 0 {
				al = rand.Intn(2)
				lines = link(i, j, al, lines0, graph)
			} else {
				al = rand.Intn(3)
				lines = link(i, j, al, lines1, graph)
			}
			He := Hexagon{
				i:        i,
				j:        j,
				typeline: tp,
				rotate:   al,
				lines:    lines,
			}
			Hexs = append(Hexs, He)
		}
	}

	//запускаем BFS по графу для раскраски областей связности

	for key, vertex := range graph {
		if vertex.isMark {
			continue
		} else {
			bfs := []VerIndex{key}
			cl := Color{R: float64(rand.Intn(256)) / 256.0, G: float64(rand.Intn(256)) / 256.0, B: float64(rand.Intn(256)) / 256.0}
			for i := 0; i < len(bfs); i++ {
				ix := bfs[i]
				graph[ix].isMark = true
				graph[ix].color = cl
				for _, val := range graph[ix].children {
					if !graph[val].isMark {
						bfs = append(bfs, val)
						//graph[val].isMark = true
					}
				}
			}
		}
	}

	//Рисуем
	dc := gg.NewContext(int(width), int(height))
	dc.SetRGB(1, 1, 1)
	dc.DrawRectangle(0, 0, width, height)
	dc.Fill()

	//Для отладки рисуем соты
	/*
		dc.SetLineWidth(r * 0.03)
		dc.SetRGB(0, 0, 0)
		for j := range m {
			for i := range n {
				start := float64(j%2) * dx / 2.0
				x := start + float64(i)*dx
				y := float64(j) * dy
				dc.Push()
				dc.DrawRegularPolygon(6, x, y, r, math.Pi/6)
				dc.Stroke()
				dc.Pop()
			}
		}
	*/
	//Рисуем кружочки в серединах ребер шестиугольников
	for key, vertex := range graph {
		i := key.i
		j := key.j
		start := float64((j)%2) * dx / 2.0
		x := start + float64(i)*dx
		y := float64(j) * dy
		dc.Push()
		dc.RotateAbout(float64(key.e)*math.Pi/3., x, y)
		cl := vertex.color
		dc.SetRGB(cl.R, cl.G, cl.B)
		dc.DrawCircle(x+dx/2, y, r*0.2)
		dc.Fill()
		dc.Pop()
	}

	dc.SetLineWidth(r * 0.1)
	for _, He := range Hexs {
		start := float64((He.j)%2) * dx / 2.0
		x := start + float64(He.i)*dx
		y := float64(He.j) * dy
		al := float64(He.rotate) * math.Pi / 3

		var cl Color
		if He.typeline == 0 {
			for i := range 3 {
				cl = graph[He.lines[i]].color
				dc.Push()
				dc.SetRGB(cl.R, cl.G, cl.B)
				dc.RotateAbout(al+float64(i)*math.Pi/3*2, x, y)
				dc.DrawArc(x, y-r, r/2, math.Pi/6, math.Pi-math.Pi/6)
				dc.Stroke()
				dc.Pop()
			}
		} else {

			for i := range 2 {
				cl = graph[He.lines[i]].color
				dc.Push()
				dc.RotateAbout(al+float64(i)*math.Pi, x, y)
				dc.SetRGB(cl.R, cl.G, cl.B)
				dc.DrawArc(x, y+r, r/2, math.Pi+math.Pi/6, 2*math.Pi-math.Pi/6)
				dc.Stroke()
				dc.Pop()
			}

			cl = graph[He.lines[2]].color
			dc.Push()
			dc.RotateAbout(al, x, y)
			dc.SetRGB(cl.R, cl.G, cl.B)
			dc.DrawLine(x-dx/2, y, x+dx/2, y)
			dc.Stroke()
			dc.Pop()
		}
	}
	dc.SavePNG("img/tile_truchet.png")

}
