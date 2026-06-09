package main

func reloadCol(layers [][]int) []colision {

	colisions := []colision{}

	span := float64(screenSizeWidth / tileSizeX)
	row := 0.
	column := 0.

	for _, t := range layers[2] {
		if t != 0 {
			colisions = append(colisions, colision{column * float64(tileSizeX), row * float64(tileSizeY)})
		}
		column += 1
		if column == span {
			column = 0
			row += 1
		}
	}

	return colisions
}