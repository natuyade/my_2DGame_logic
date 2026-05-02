package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Game struct{}

// default: 60tps
// 毎フレーム画面リセット(クリア),描画される
func (g *Game) Update() error {
	return nil
}

// インポートした全てのimage, offscreenImage(メモリ上のバッファに描画されたimage), screen(最終的に描画されている画面)
// これらを*ebiten.Imageとして扱う
func (g *Game) Draw(screen *ebiten.Image) {
	// 画面上にdebugメッセージを描画するutility関数
	// 毎フレーム画面はクリアされるためDrawで毎フレーム描画する必要がある
	ebitenutil.DebugPrint(screen, "It's DebugMessage")
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