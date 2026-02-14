package main

func main() {
	hex_tiles()
	penrose_P2(1200, 7, "img/tile_P2.png", true)
	penrose_P2(1200, 7, "img/tile_P2_gen.png", false)
	penrose_P3(1200, 7, "img/tile_P3.png", true)
	penrose_P3(1200, 7, "img/tile_P3_gen.png", false)

	println("ok")
}
