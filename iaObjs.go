package main

import (
	_ "log"
)

func intractTo(g *Game, ias []interactiveObj, pd playerData) {
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
			g.message = append(g.message, "You found the entrance key.")
			return true
		}
		g.message = append(g.message, "The box is already open.")
	// LockedDoor
	case 1:
		if !used {
			g.layers[2][20] = 0
			g.layers[1][20] = 10

			g.cols = reloadCol(g.layers)
			return true
		}
	// firstSign
	case 2:
		if !used {
			g.message = append(g.message, "The story begins here.")
			return true
		}
		g.message = append(g.message, "The story begins here.\n(read)")
	}

	return used
}
