package main

import (
	"fmt"
	"image"
	_ "image/png"
	"log"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

var images []*ebiten.Image

const (
	windowWidth  = int(1280)
	windowHeight = int(720)

	screenSizeWidth  = int(256)
	screenSizeHeight = int(224)

	tileSizeX = int(16)
	tileSizeY = int(16)
)

type Game struct {
	layers        [][]int
	keys          []ebiten.Key
	playerX       float64
	playerY       float64
	playerLookAt  int
	player2X      float64
	player2Y      float64
	player2LookAt int
	movedDebug    [2]float64
	cols          []Colision
}

type Colision struct {
	x float64
	y float64
}

func loadImage(path string) *ebiten.Image {
	// NewImageFromFile(相対パス): 画像ファイルから再利用可能なebitengineImageObjectを生成
	img, _, err := ebitenutil.NewImageFromFile(path)
	if err != nil {
		log.Fatal(err)
	}
	return img
}

func init() {
	fishImg := loadImage("assets/images/fishish.png")
	playerImg := loadImage("assets/images/player.png")
	tilesImg := loadImage("assets/images/tiles.png")
	boxImg := loadImage("assets/images/box.png")

	images = append(images, fishImg, playerImg, tilesImg, boxImg)
}

func objRotate(img *ebiten.Image, angle float64) *ebiten.DrawImageOptions {

	// Imageのx,yは左上が原点
	op := &ebiten.DrawImageOptions{}

	// Boundsは画像の左上, 右下の座標を持ったRectangle
	// Dx,DyでRectangleの横幅, 縦幅を取得
	w := img.Bounds().Dx()
	h := img.Bounds().Dy()

	// 画像の中心の座標を原点に移動
	op.GeoM.Translate(-float64(w)/2, -float64(h)/2)

	// 原点で回転する為, 移動後の画像の中心を原点に移動し回転することで
	// 回転後の座標計算をなくせる
	op.GeoM.Rotate(angle)

	// 左上基準に戻す
	op.GeoM.Translate(float64(w)/2, float64(w)/2)

	return op
}

func drawLayers() [][]int {

	layers := [][]int{
		// 地面
		{
			1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 2, 2, 2, 2, 2,
			1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 2, 2, 2, 2, 2,
			1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 2, 2, 2, 2, 2,
			1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 2, 2, 2, 2, 2,
			1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 2, 2, 2, 2,
			1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 2, 2, 2, 2,
			1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 2, 2, 2, 2,
			1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 2, 2, 2, 2,
			1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 2, 2, 2, 2,
			1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 2, 2, 2, 2,
			1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 2, 2, 2,
			1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 2, 2, 2,
			1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 2, 2, 2,
			1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 2, 2,
		},
		// 当たり判定無しobj
		{
			0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 3, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 3, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
			0, 4, 4, 4, 5, 4, 4, 4, 4, 4, 4, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 3, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 3, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 3, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 3, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 3, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 3, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 3, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 3, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 3, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		},
		// 当たり判定有
		{
			6, 6, 6, 6, 6, 6, 6, 6, 6, 6, 6, 6, 6, 6, 6, 6,
			6, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 6,
			6, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 6,
			6, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 6,
			6, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 6,
			6, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 6,
			6, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 6,
			6, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 6,
			6, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 6,
			6, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 6,
			6, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 6,
			6, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 6,
			6, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 6,
			6, 6, 6, 6, 6, 6, 6, 6, 6, 6, 6, 6, 6, 6, 6, 6,
		},
	}
	return layers
}

func convertDir(pickedNum int) [2]int {

	// タイルの画像をNo.で取得できるように.
	// ttileImageが4*4の場合
	// 0=(0,0) 1=(16,0) 2=(32,0) 3=(48,0) 4=(0,16) 5=(16,16)...

	// tileImageの横Cell数を計算
	tileImageWidth := images[2].Bounds().Dx() / tileSizeX

	// 選択された数字に対応して,tileImageの縦何列目かを計算
	tileY := int(math.Floor(float64(pickedNum) / float64(tileImageWidth)))

	// それぞれ座標に変換
	pickedTile := [2]int{
		(pickedNum - (tileY * tileImageWidth)) * tileSizeX,
		tileY * tileSizeY,
	}

	return pickedTile
}

func moveVector(result [2]float64, moveSpeed float64) [2]float64 {

	// 斜めに移動する場合の処理
	if result[0] != 0 && result[1] != 0 {

		// vectorの計算(a^2+b^2=c^2)
		v := math.Sqrt(result[0]*result[0] + result[1]*result[1])

		// 移動量の最終計算
		result[0] = (result[0] / v) * moveSpeed
		result[1] = (result[1] / v) * moveSpeed
	} else {
		result[0] *= moveSpeed
		result[1] *= moveSpeed
	}

	return [2]float64{result[0], result[1]}
}

// default: 60tps
// 毎フレーム画面リセット(クリア),描画される
func (g *Game) Update() error {
	g.keys = inpututil.AppendPressedKeys(g.keys[:0])
	g.cols[0] = Colision{g.player2X, g.player2Y}

	// 最終的な移動量
	result := [2]float64{0, 0}
	moveValue := 1.
	moveSpeed := 1.

	for _, k := range g.keys {

		// Sprintの倍率
		switch k {
		case ebiten.KeyShiftLeft:
			moveSpeed = 3.5
		}

		beingColision := false
		var colision [2]float64

		switch k {
		// 僕が考えた最強のコリジョンコード
		case ebiten.KeyW:
			for _, col := range g.cols {

				if g.playerY > col.y &&
					(g.playerY-(moveValue*moveSpeed))-col.y <= float64(tileSizeY) &&
					g.playerX-col.x < float64(tileSizeX) &&
					g.playerX-col.x > -float64(tileSizeX) {
					colision = [2]float64{col.x, col.y}
					beingColision = true
				}
			}

			if beingColision {
				result[1] -= g.playerY - (colision[1] + float64(tileSizeY))
			} else {
				result[1] -= moveValue
			}
		case ebiten.KeyA:
			for _, col := range g.cols {

				if g.playerX > col.x &&
					(g.playerX-(moveValue*moveSpeed))-col.x <= float64(tileSizeX) &&
					g.playerY-col.y < float64(tileSizeY) &&
					g.playerY-col.y > -float64(tileSizeY) {
					colision = [2]float64{col.x, col.y}
					beingColision = true
				}
			}

			if beingColision {
				result[0] -= g.playerX - (colision[0] + float64(tileSizeX))
			} else {
				result[0] -= moveValue
			}
		case ebiten.KeyS:
			for _, col := range g.cols {

				if g.playerY < col.y &&
					(g.playerY+(moveValue*moveSpeed))-col.y >= -float64(tileSizeY) &&
					g.playerX-col.x < float64(tileSizeX) &&
					g.playerX-col.x > -float64(tileSizeX) {
					colision = [2]float64{col.x, col.y}
					beingColision = true
				}
			}

			if beingColision {
				result[1] -= g.playerY - (colision[1] - float64(tileSizeY))
			} else {
				result[1] += moveValue
			}
		case ebiten.KeyD:
			for _, col := range g.cols {

				if g.playerX < col.x &&
					(g.playerX+(moveValue*moveSpeed))-col.x >= -float64(tileSizeX) &&
					g.playerY-col.y < float64(tileSizeY) &&
					g.playerY-col.y > -float64(tileSizeY) {
					colision = [2]float64{col.x, col.y}
					beingColision = true
				}
			}

			if beingColision {
				result[0] -= g.playerX - (colision[0] - float64(tileSizeX))
			} else {
				result[0] += moveValue
			}
		}
	}

	// ベクトル計算
	resultMoved := moveVector(result, moveSpeed)

	g.playerX += resultMoved[0]
	g.playerY += resultMoved[1]
	g.movedDebug = [2]float64{resultMoved[0], resultMoved[1]}

	// 2p debug
	result2P := [2]float64{0, 0}
	moveSpeed2P := float64(1)

	for _, k := range g.keys {

		switch k {
		case ebiten.KeyArrowUp:
			result2P[1] -= 1
		case ebiten.KeyArrowLeft:
			result2P[0] -= 1
		case ebiten.KeyArrowDown:
			result2P[1] += 1
		case ebiten.KeyArrowRight:
			result2P[0] += 1
		case ebiten.KeyShiftRight:
			moveSpeed2P = 3.5
		}
	}

	resultMoved2P := moveVector(result2P, moveSpeed2P)

	g.player2X += resultMoved2P[0]
	g.player2Y += resultMoved2P[1]

	return nil
}

// インポートした全てのimage, offscreenImage(メモリ上のバッファに描画されたimage), screen(最終的に描画されている画面)
// これらを*ebiten.Imageとして扱う
func (g *Game) Draw(screen *ebiten.Image) {

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Reset()

	drawX := 0
	drawY := 0
	nowPos := [2]float64{0, 0}

	screenSizeX := screen.Bounds().Dx()
	tileSpan := screenSizeX / tileSizeX

	for _, y := range g.layers {
		for _, x := range y {
			tileId := x
			// 0の場合描画しないようにするため-1(tile描画を1~に)
			pickedTile := convertDir(tileId - 1)
			tileDir := image.Rect(pickedTile[0], pickedTile[1], pickedTile[0]+tileSizeX, pickedTile[1]+tileSizeY)

			if drawX == tileSpan {
				drawY += tileSizeY
				drawX = 0
				op.GeoM.Translate(-float64(screenSizeX), float64(drawY))
			}

			if tileId != 0 {
				screen.DrawImage(images[2].SubImage(tileDir).(*ebiten.Image), op)
			}

			// 移動先の座標を保存
			nowPos[0], nowPos[1] = op.GeoM.Apply(float64(tileSizeX), 0)
			op.GeoM.Translate(float64(tileSizeX), 0)

			drawX += 1
			drawY = 0
		}
		drawX = 0
		op.GeoM.Translate(-nowPos[0], -nowPos[1])
	}

	op = objRotate(images[0], -math.Pi/2)

	// 基本的にScale変更する場合, 座標移動の前にした方がやりやすい
	op.GeoM.Scale(0.8, 0.8)
	op.GeoM.Translate(32, 32)
	//screen.DrawImage(images[0], op)

	op.GeoM.Reset()

	// 移動方向によってImageを反転
	for _, key := range g.keys {
		switch key {
		case ebiten.KeyA:
			g.playerLookAt = 1
		case ebiten.KeyD:
			g.playerLookAt = 2
		}
	}

	switch g.playerLookAt {
	case 1:
		op.GeoM.Scale(1, 1)
		op.GeoM.Translate(g.playerX, g.playerY)
	case 2:
		op.GeoM.Scale(-1, 1)
		op.GeoM.Translate(float64(images[1].Bounds().Dx())+g.playerX, g.playerY)
	}

	// Draw Player Image
	screen.DrawImage(images[3], op)

	// 2p debug
	op.GeoM.Reset()
	for _, key := range g.keys {
		switch key {
		case ebiten.KeyArrowLeft:
			g.player2LookAt = 1
		case ebiten.KeyArrowRight:
			g.player2LookAt = 2
		}
	}
	switch g.player2LookAt {
	case 1:
		op.GeoM.Scale(1, 1)
		op.GeoM.Translate(g.player2X, g.player2Y)
	case 2:
		op.GeoM.Scale(-1, 1)
		op.GeoM.Translate(float64(images[1].Bounds().Dx())+g.player2X, g.player2Y)
	}

	screen.DrawImage(images[3], op)

	// box for debug colision
	op.GeoM.Reset()
	op.GeoM.Translate(128, 128)
	screen.DrawImage(images[3], op)

	// 画面上にdebugメッセージを描画するutility関数
	// 毎フレーム画面はクリアされるためDrawで毎フレーム描画する必要がある
	ebitenutil.DebugPrint(screen, fmt.Sprintf("%s\nmoved: x[%f] y[%f]\n1p\n[%f]\n[%f]\n2p\n[%f]\n[%f]", g.keys, g.movedDebug[0], g.movedDebug[1], g.playerX, g.playerY, g.player2X, g.player2Y))

}

// windowサイズを引数で受け取り
// ゲーム画面の論理サイズ(px)を返す関数
func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	// 今のコードは,固定値を入れることで
	// Windowサイズに関係なくゲーム画面のサイズを固定している
	return screenSizeWidth, screenSizeHeight
}

func main() {
	// window表示時のサイズ指定
	ebiten.SetWindowSize(windowWidth, windowHeight)
	// windowTitle
	ebiten.SetWindowTitle("My first app in go language")

	layers := drawLayers()

	g := &Game{
		layers:        layers,
		keys:          []ebiten.Key{},
		playerX:       64,
		playerY:       48,
		playerLookAt:  2,
		player2X:      128,
		player2Y:      48,
		player2LookAt: 2,
	}
	g.cols = append(g.cols, Colision{g.player2X, g.player2Y})
	g.cols = append(g.cols, Colision{128, 128})

	if len(layers) <= 3 {

		span := float64(screenSizeWidth / tileSizeX)
		row := 0.
		column := 0.

		// layer[3]で置いたタイルをColision付objとして扱うためのappend
		for _, t := range layers[2] {
			if t != 0 {
				g.cols = append(g.cols, Colision{column * float64(tileSizeX), row * float64(tileSizeY)})
			}
			column += 1
			if column == span {
				column = 0
				row += 1
			}
		}
	}

	// Gameのメインループを実行
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
