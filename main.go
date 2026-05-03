package main

import (
	"image/color"
	_ "image/png"
	"log"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

var fishimg *ebiten.Image

type Game struct{
	fishImage *ebiten.Image
}

func init() {
	var err error
	// NewImageFromFile(相対パス): 画像ファイルから再利用可能なebitengineImageObjectを生成
	fishimg, _, err = ebitenutil.NewImageFromFile("fishish.png")
	if err != nil {
		log.Fatal(err)
	}
}

func objRotate(screen *ebiten.Image, img *ebiten.Image, angle float64) *ebiten.DrawImageOptions {

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

// default: 60tps
// 毎フレーム画面リセット(クリア),描画される
func (g *Game) Update() error {
	return nil
}

// インポートした全てのimage, offscreenImage(メモリ上のバッファに描画されたimage), screen(最終的に描画されている画面)
// これらを*ebiten.Imageとして扱う
func (g *Game) Draw(screen *ebiten.Image) {

	screen.Fill(color.RGBA{0, 0xff, 0, 0xff})

	g.fishImage = fishimg

	op := objRotate(screen, g.fishImage, -math.Pi / 2)

	// 基本的にScale変更する場合, 座標移動の前にした方がやりやすい
	op.GeoM.Scale(0.8,0.8)
	op.GeoM.Translate(32, 32)

	screen.DrawImage(g.fishImage, op)

	// 画面上にdebugメッセージを描画するutility関数
	// 毎フレーム画面はクリアされるためDrawで毎フレーム描画する必要がある
	ebitenutil.DebugPrint(screen, "This is DebugMessage")
}

// windowサイズを引数で受け取り
// ゲーム画面の論理サイズ(px)を返す関数
func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	// 今のコードは,固定値を入れることで
	// Windowサイズに関係なくゲーム画面のサイズを固定している
	return 320, 240
}

func main() {
	// window表示時のサイズ指定
	ebiten.SetWindowSize(640, 480)
	// windowTitle
	ebiten.SetWindowTitle("My first app in go language")

	// Gameのメインループを実行
	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}
}