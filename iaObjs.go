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
	switch g.currentStage {
	case 1:
		switch id {
		// boxContaingEntranceKey
		case 0:
			if !used {
				g.layers[2][42] = 12
				g.items[0].amount += 1
				g.message = append(g.message, "入口のカギを見つけた")
				return true
			}
			g.message = append(g.message, "...空っぽだ")

		// LockedDoor
		case 1:
			if !used {
				if g.items[0].amount == 0 {
					g.message = append(g.message, "鍵がかかっている")
					return false
				}
				g.items[0].amount -= 1

				g.layers[2][20] = 0
				//g.layers[1][20] = 10
				g.layers[2][20] = 10

				g.cols = reloadCol(g.layers)
				g.message = append(g.message, "鍵が開いた")
				return true
			}
			g.currentStage = 2
			g.layers = drawLayers(g.currentStage)
			g.cols = reloadCol(g.layers)

		// firstSign
		case 2:
			if !used {
				g.message = append(g.message, "始まりの地\n(Spaceでメッセージを進める)")
				return true
			}
			g.message = append(g.message, "始まりの地")
		}

	}

	return used
}
