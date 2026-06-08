package main

import _ "log"

func intractTo(gmap [][]int, ias []interactiveObj, pd playerData, tileSizeX int, tileSizeY int) ([][]int, []interactiveObj) {
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
				gmap, ias[i].used = intEvent(gmap, i, ia.used)
			}
		case 2:
			if ia.y-pd.y <= ty/2 &&
				ia.y-pd.y >= -ty/2 &&
				ia.x-pd.x <= -tx &&
				ia.x-pd.x >= -(tx+pReach) {
				gmap, ias[i].used = intEvent(gmap, i, ia.used)
			}
		case 3:
			if ia.x-pd.x <= tx/2 &&
				ia.x-pd.x >= -tx/2 &&
				ia.y-pd.y >= ty &&
				ia.y-pd.y <= ty+pReach {
				gmap, ias[i].used = intEvent(gmap, i, ia.used)
			}
		case 4:
			if ia.y-pd.y <= ty/2 &&
				ia.y-pd.y >= -ty/2 &&
				ia.x-pd.x >= tx &&
				ia.x-pd.x <= tx+pReach {
				gmap, ias[i].used = intEvent(gmap, i, ia.used)
			}
		}
	}

	return gmap, ias
}

// インタラクトされたobjに対応する処理
func intEvent(gmap [][]int, id int, used bool) ([][]int, bool) {
	switch id {
	// boxContaingEntranceKey
	case 0:
		if !used {
			gmap[2][42] = 12
		}
	// LockedDoor
	case 1:
		if !used {
			gmap[2][4] = 10
		}
	}

	return gmap, used
}
