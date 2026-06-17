package main

func inputIaObj() [][]interactiveObj {
	objs := [][]interactiveObj{
		{
			interactiveObj{160, 32, false},//boxContaingEntranceKey
			interactiveObj{64, 16, false},//LockedDoorIa
			interactiveObj{80, 112, false},//firstSign
		},
		{
			interactiveObj{48, 48, false},//secondSign
		},
	}

	return objs
}