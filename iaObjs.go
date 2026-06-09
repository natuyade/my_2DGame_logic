package main

import _ "log"

func intractTo(g *Game, ias []interactiveObj, pd playerData, tileSizeX int, tileSizeY int) {
	pReach := 2.
	tx := float64(tileSizeX)
	ty := float64(tileSizeY)

	for i, ia := range ias {
		switch pd.lookAt {
		case 1:
			if ia.x-pd.x <= tx/2 &&
				ia.x-pd.x >= -tx/2 &&
				ia.y-pd.y <= -ty &&
				ia.y-pd.y >= -(ty+pReach) {
				ias[i].used = intEvent(g, i, ia.used)
			}
		case 2:
			if ia.y-pd.y <= ty/2 &&
				ia.y-pd.y >= -ty/2 &&
				ia.x-pd.x <= -tx &&
				ia.x-pd.x >= -(tx+pReach) {
				ias[i].used = intEvent(g, i, ia.used)
			}
		case 3:
			if ia.x-pd.x <= tx/2 &&
				ia.x-pd.x >= -tx/2 &&
				ia.y-pd.y >= ty &&
				ia.y-pd.y <= ty+pReach {
				ias[i].used = intEvent(g, i, ia.used)
			}
		case 4:
			if ia.y-pd.y <= ty/2 &&
				ia.y-pd.y >= -ty/2 &&
				ia.x-pd.x >= tx &&
				ia.x-pd.x <= tx+pReach {
				ias[i].used = intEvent(g, i, ia.used)
			}
		}
	}
}

// インタラクトされたobjに対応する処理
func intEvent(g *Game, id int, used bool) bool {
	switch id {
	// boxContaingEntranceKey
	case 0:
		if !used {
			g.layers[2][42] = 12
			return true
		}
	// LockedDoor
	case 1:
		if !used {
			g.layers[2][4] = 0
			g.layers[1][4] = 10

			g.cols = reloadCol(g.layers)
			return true
		}
	}

	return used
}
