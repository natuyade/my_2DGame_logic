package main

func inputIaObj() []interactiveObj {
	objs := []interactiveObj{}
	objs = append(objs,
		interactiveObj{160, 32, false},//boxContaingEntranceKey
		interactiveObj{64, 16, false},//LockedDoorIa
		interactiveObj{80, 112, false},//firstSign
	)
	
	return objs
}