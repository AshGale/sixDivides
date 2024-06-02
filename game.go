package main

// import (
//     "github.com/hajimehoshi/ebiten/v2"
//     "github.com/hajimehoshi/ebiten/v2/ebitenutil"
//     "log"
// )

// const (
//     screenWidth  = 800
//     screenHeight = 600
//     margin        = 5
//     width         = 50
//     height        = 50
//     radius        = width / 2 - margin
//     numPieces    = 12
// )

// type Piece struct {
//     x, y   int
//     color string
// }

// var (
//     screen *ebiten.Image
//     pieces []Piece
//     dices    [1]int
// )

// func drawBoard() {
//     for row := 0; row < 8; row++ {
//         for col := 0; col < 8; col++ {
//             if (row + col) % 2 == 0 {
//                 color = ebiten.NewImageFromImageOptions(
//                     [screenWidth, screenHeight],
//                     ebiten.FilterDefault,
//                     [MARGIN+WIDTH]*col+MARGIN,
//                     [MARGIN+HEIGHT]*row+MARGIN,
//                     WIDTH,
//                     HEIGHT,
//                 )
//             } else {
//                 color = ebitenutil.ImageFromImageOptions(
//                     [screenWidth, screenHeight],
//                     ebiten.FilterDefault,
//                     [MARGIN+WIDTH]*col+MARGIN,
//                     [MARGIN+HEIGHT]*row+MARGIN,
//                     WIDTH,
//                     HEIGHT,
//                 )
//             }
//         }
//     }
// }

// func drawPieces() {
//     for _, piece := range pieces {
//         if piece.color == "white" {
//             ebitenutil.DrawCircle(screen, ebiten.White, piece.x, piece.y, radius)
//         } else {
//             ebitenutil.DrawCircle(screen, ebiten.Black, piece.x, piece.y, radius)
//         }
//     }
// }

// func createPieces() {
//     colors := []string{"white", "black"} * numPieces

//     for i := 0; i < len(colors); i++ {
//         if i % 2 == 0 {
//             xPos := MARGIN + WIDTH / 4
//             yPos := (MARGIN + HEIGHT) * i / 2 - HEIGHT / 2
//         } else {
//             xPos := screenWidth - MARGIN - WIDTH / 4 * 3
//             yPos := (MARGIN + HEIGHT) * i / 2 - HEIGHT / 2
//         }

//         pieces = append(pieces, Piece{
//             x:      xPos,
//             y:      yPos,
//             color: colors[i],
//         })
//     }
// }

// func rollDice() int {
//     return rand.Intn(6) + 1
// }

// func drawDices() {
//     for _, dice := range dices {
//         ebitenutil.DrawCircle(screen, ebiten.White, dice, radius)
//     }
// }

// func update(screen *ebiten.Image) error {
//     if ebiten.IsDrawingSkipped() {
//         return nil
//     }

//     createPieces()
//     drawBoard()
//     drawPieces()
//     drawDices()

//     return screen.Update()
// }

// func main() {
//     var err error

//     screen, err = ebiten.NewImage(screenWidth, screenHeight, ebiten.FilterDefault)
//     if err != nil {
//         log.Fatal(err)
//     }

//     if err := ebiten.RunGame(&ebiten.Game{
//         Update: update,
//         Draw:   func(screen *ebiten.Image) {},
//     }); err != nil {
//         log.Fatal(err)
//     }
// }
