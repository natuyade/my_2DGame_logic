package main

func stageChange(g *Game, nextStage int) {
	g.currentStage = nextStage
	g.cols = reloadCol(g.layers[g.currentStage])
}