package main

import (
	"math"
	"math/cmplx"

	"github.com/fogleman/gg"
)

type TriangleP3 struct {
	R, G, B complex128
	Type    int
	//RKey    string
}

/*
На этой картинке правила измельчения мозаики.
https://upload.wikimedia.org/wikipedia/commons/thumb/1/1a/Penrose-P3-deflation.svg/500px-Penrose-P3-deflation.svg.png

Каждая вершина треугольника имеет свой цвет
R - красная дуга
G - зеленая дуга
B - нет дуги
*/
func (tri TriangleP3) split() []Shape {
	if tri.Type == 0 {
		A := tri.B + (tri.R-tri.B)*complex(phi, 0)
		return []Shape{TriangleP3{
			R:    tri.G,
			G:    tri.B,
			B:    A,
			Type: 1,
		}, TriangleP3{
			R:    A,
			G:    tri.R,
			B:    tri.G,
			Type: 0,
		}}
	} else {
		A := (tri.G-tri.R)*complex(phi, 0) + tri.R
		B := (tri.B-tri.R)*complex(phi, 0) + tri.R
		return []Shape{TriangleP3{
			R:    tri.G,
			G:    tri.B,
			B:    A,
			Type: 1,
		}, TriangleP3{
			R:    B,
			G:    tri.B,
			B:    A,
			Type: 0,
		}, TriangleP3{
			R:    A,
			G:    tri.R,
			B:    B,
			Type: 1,
		}}
	}
}

func (tri TriangleP3) Draw(width float64, height float64, dc *gg.Context, lens map[int]int, graph map[string]*PenVertex,
	ra float64) {

	c := tri.R + tri.G
	if math.Abs(real(c)) > width+2*ra || math.Abs(imag(c)) > height+2*ra {
		return
	}

	rg := ra * phi / 2
	rr := ra - rg

	dc.Push()
	dc.SetLineWidth(2)

	if tri.Type == 1 {
		//Вычисляем свет
		RKey, _ := tri.getLink()
		num := graph[RKey].index //номер связности
		nn := lens[num]          //Число элементов связности

		//H := math.Abs(math.Cos(float64(nn-5)*2027)) * 360 //случайный оттенок
		H := math.Abs(math.Cos(math.Log(float64(nn)*2027))) * 360 //случайный оттенок
		dc.SetRGB(HSLToRGB(H, 0.8, 0.5))
		dc.DrawLine(real(tri.R)+width/2, imag(tri.R)+height/2, real(tri.G)+width/2, imag(tri.G)+height/2)
		dc.Stroke()
	} else {
		dc.SetRGB(1.0, 1.0, 1.0)
	}
	dc.MoveTo(real(tri.R)+width/2, imag(tri.R)+height/2)
	dc.LineTo(real(tri.B)+width/2, imag(tri.B)+height/2)
	dc.LineTo(real(tri.G)+width/2, imag(tri.G)+height/2)
	dc.ClosePath()
	dc.Fill()
	dc.SetRGB(0, 0, 0)

	dc.DrawLine(real(tri.R)+width/2, imag(tri.R)+height/2, real(tri.B)+width/2, imag(tri.B)+height/2)
	dc.DrawLine(real(tri.G)+width/2, imag(tri.G)+height/2, real(tri.B)+width/2, imag(tri.B)+height/2)
	dc.Stroke()

	dc.SetLineWidth(2)

	//рисуем зеленые дуги

	angleB := cmplx.Phase(tri.R - tri.G)
	angleC := cmplx.Phase(tri.B - tri.G)

	start, end := math.Min(angleB, angleC), math.Max(angleB, angleC)
	delta := end - start
	//fmt.Println(start, end, delta)

	if delta > math.Pi {
		delta = math.Pi*2 - delta
		end = start - delta
	}
	//

	dc.SetRGB(0, 1, 0) // Зеленый
	dc.DrawArc(real(tri.G)+width/2, imag(tri.G)+height/2, rg, start, end)
	dc.Stroke()

	//рисуем красные дуги
	angleB = cmplx.Phase(tri.G - tri.R)
	angleC = cmplx.Phase(tri.B - tri.R)

	start, end = math.Min(angleB, angleC), math.Max(angleB, angleC)
	delta = end - start

	if delta > math.Pi {
		delta = math.Pi*2 - delta
		end = start - delta
	}

	dc.SetRGB(0.65, 0.16, 0.16) // Красный

	r := rr
	if tri.Type == 0 {
		r = rg
	}
	dc.DrawArc(real(tri.R)+width/2, imag(tri.R)+height/2, r, start, end)
	dc.Stroke()

	dc.Pop()

}

func (tri TriangleP3) getLink() (ia string, ib string) {
	pa := (tri.R + tri.G) * complex(0.5, 0)
	pb := (tri.R + tri.B) * complex(0.5, 0)
	ia = getPointKey(pa)
	ib = getPointKey(pb)
	//tri.RKey = ia
	return ia, ib
}

func penrose_P3() {
	var height float64 = 1200
	var ra float64 = height*math.Sqrt(2) + 1.
	A := complex(ra, 0)
	rotator := cmplx.Exp(complex(0, math.Pi/5))
	//tris := []*TriangleP3{}
	tris := []Shape{}
	for i := range 10 {
		B := A * rotator
		tri := TriangleP3{
			B:    0 + 0i,
			Type: 0,
			//RKey: "",
		}
		if i%2 == 0 {
			tri.R = A
			tri.G = B
		} else {
			tri.R = B
			tri.G = A
		}
		tris = append(tris, tri)
		A = B
	}
	penrose(height, tris, ra, 8, "img/tile_P3.png")
}
