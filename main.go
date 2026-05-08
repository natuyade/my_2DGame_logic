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
	windowWidth = int(1280)
	windowHeight = int(720)

	screenSizeWidth = int(256)
	screenSizeHeight = int(220)

	tileSizeX = int(16)
	tileSizeY = int(16)
)

type Game struct{
	layers [][]int
	keys []ebiten.Key
	playerX float64
	playerY float64
	movedDebug [2]float64
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
	op.GeoM.Translate(-float64(w) / 2, -float64(h) / 2)

	// 原点で回転する為, 移動後の画像の中心を原点に移動し回転することで
	// 回転後の座標計算をなくせる
	op.GeoM.Rotate(angle)
	
	// 左上基準に戻す
	op.GeoM.Translate(float64(w) / 2, float64(w) / 2)

	return op
}

func drawLayers() [][]int {

	layers := [][]int{
		{
			1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 2, 2, 2, 2, 2, 2,
			1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 2, 2, 2, 2, 2,
			1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 2, 2, 2, 2,
		},
		{
			0, 0, 0, 0, 3, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
			0, 4, 4, 4, 3, 4, 4, 4, 0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 3, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
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

// default: 60tps
// 毎フレーム画面リセット(クリア),描画される
func (g *Game) Update() error {
	g.keys = inpututil.AppendPressedKeys(g.keys[:0])

	// 最終的な移動量
	resultMoved := []float64{0, 0}
	moveSpeed := float64(1)

	for _, k := range g.keys {
		key := k.String()

		switch key {
		case "W":
			resultMoved[1] -= 1
		case "A":
			resultMoved[0] -= 1
		case "S":
			resultMoved[1] += 1
		case "D":
			resultMoved[0] += 1
		// Sprintの倍率
		case "Shift":
			moveSpeed = 3.5
		}
	}

	// 斜めに移動する場合の処理
	if resultMoved[0] != 0 && resultMoved[1] != 0 {

		// vectorの計算(a^2+b^2=c^2)
		v := math.Sqrt(resultMoved[0] * resultMoved[0] + resultMoved[1] * resultMoved[1])

		// 移動量の計算
		resultMoved[0] = (resultMoved[0] / v) * moveSpeed
		resultMoved[1] = (resultMoved[1] / v) * moveSpeed
	} else {
		resultMoved[0] *= moveSpeed
		resultMoved[1] *= moveSpeed
	}

	g.playerX += resultMoved[0]
	g.playerY += resultMoved[1]

	g.movedDebug = [2]float64{resultMoved[0], resultMoved[1]}

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
				screen.DrawImage(images[2].SubImage(tileDir).(*ebiten.Image) , op)
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

	op = objRotate(images[0], -math.Pi / 2)

	// 基本的にScale変更する場合, 座標移動の前にした方がやりやすい
	op.GeoM.Scale(0.8,0.8)
	op.GeoM.Scale(0.,0.)
	op.GeoM.Translate(32, 32)

	screen.DrawImage(images[0], op)

	op.GeoM.Reset()
	op.GeoM.Translate(g.playerX, g.playerY)

	screen.DrawImage(images[1], op)

	// 画面上にdebugメッセージを描画するutility関数
	// 毎フレーム画面はクリアされるためDrawで毎フレーム描画する必要がある
	ebitenutil.DebugPrint(screen, fmt.Sprintf("%s\nmoved: x[%f] y[%f]\n", g.keys, g.movedDebug[0], g.movedDebug[1]))
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
		layers: layers,
		keys: []ebiten.Key{},
		playerX: 0,
		playerY: 0,
	}

	// Gameのメインループを実行
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}