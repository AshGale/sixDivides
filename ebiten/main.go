package main

import (
	"bytes"
	"fmt"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type Player struct {
	Color            color.Color
	Name             string
	Pieces           []Piece
	Actions          int
	PlayerIndex      int
	StartingPosition Position
}

type Position struct {
	X int
	Y int
}

type Piece struct {
	Color       color.Color
	Value       int
	PlayerIndex int
	Position    Position
}

type Tile struct {
	Position Position
	Piece    Piece
}

type Board struct {
	Tiles    [][]Tile
	Width    int
	Height   int
	TileSize int
}

// Game represents the game state
type Game struct {
	keyStates                   map[ebiten.Key]bool
	board                       Board
	players                     []Player
	turn                        int
	HighlightedTile             Position
	SelectedTile                Position
	InvalidTile                 Position
	GameOver                    bool
	gameState                   int
	uiMenueSelectedButton       int
	uiMenueButtonNumber         int
	uiNewGameConfirmation       bool
	uiNewGameSectionPlayer      []int
	uiNewGameSectionHighlighted int
	uiStartNewGameButton        bool
	screenSize                  Position
}

// createBoard creates a new board with the given width and height
func createBoard(width int, height int, tileSize int) Board {
	board := Board{Width: width, Height: height}
	board.TileSize = tileSize

	// make the the tiles array and account for the extra row/column for the boarder
	tiles := make([][]Tile, width+1)
	for x := 0; x <= width; x++ {
		tiles[x] = make([]Tile, height+1)
		for y := 0; y <= height; y++ {
			tile := Tile{Position: Position{X: x, Y: y}, Piece: Piece{}}
			tiles[x][y] = tile
		}
	}
	board.Tiles = tiles

	return board
}

func createPlayers(playerPositions []int) []Player {

	numberOfPlayers := 0
	for _, section := range playerPositions {
		if section != -1 {
			numberOfPlayers++
		}
	}

	var players []Player = make([]Player, numberOfPlayers)
	var startingPosition Position
	var playerColor color.Color
	var playerName string
	for i, p := range playerPositions {
		if p != -1 {
			switch i {
			case 0:
				playerName = fmt.Sprintf("Player%v", p)
				playerColor = color.RGBA{0x00, 0x00, 0xff, 0xff}
				startingPosition = Position{1, 1}
			case 1:
				playerName = fmt.Sprintf("Player%v", p)
				playerColor = color.RGBA{0xff, 0x00, 0x00, 0xff}
				startingPosition = Position{1, 6}
			case 2:
				playerName = fmt.Sprintf("Player%v", p)
				playerColor = color.RGBA{0x00, 0xff, 0xff, 0xff}
				startingPosition = Position{6, 1}
			case 3:
				playerName = fmt.Sprintf("Player%v", p)
				playerColor = color.RGBA{0xff, 0x00, 0xff, 0xff}
				startingPosition = Position{6, 6}
			}
			players[p-1] = newPlayer(playerName, playerColor, startingPosition, p-1)
		}
	}

	return players
}

func newPlayer(name string, playerColor color.Color, startingPosition Position, playerGameIndex int) Player {
	player := Player{
		Color:            playerColor,
		Name:             name,
		StartingPosition: startingPosition,
		PlayerIndex:      playerGameIndex,
		Actions:          0,
		Pieces:           make([]Piece, 1),
	}

	startingPiece := Piece{Color: playerColor, Value: 6, PlayerIndex: playerGameIndex, Position: startingPosition}
	player.Pieces[0] = startingPiece
	return player
}

func handleTileMove(g *Game, xOffset, yOffset int) {
	// check if the selected tile is set, if so move the piece on the tile, to the tile above it
	if g.SelectedTile.X == -1 && g.SelectedTile.Y == -1 {
		// set the new highlighted tile to true and the previous one to false
		g.HighlightedTile.X += xOffset
		g.HighlightedTile.Y += yOffset
		log.Printf("highlighter is now at %d, %d", g.HighlightedTile.X, g.HighlightedTile.Y)
	} else {
		// move the piece on the selected tile to the tile above it
		var noPiece Piece = Piece{}

		// get the target tiles position
		targetX, targetY := g.HighlightedTile.X+xOffset, g.HighlightedTile.Y+yOffset
		selectedX, selectedY := g.SelectedTile.X, g.SelectedTile.Y

		if g.board.Tiles[targetX][targetY].Piece == noPiece {

			if g.board.Tiles[selectedX][selectedY].Piece.Value == 6 {
				// create a piece on the new tile of value 1
				newPiece := Piece{Color: g.players[g.turn].Color, Value: 1, PlayerIndex: g.players[g.turn].PlayerIndex, Position: Position{X: targetX, Y: targetY}}

				// add newpiece to the players pieces array
				g.players[g.turn].Pieces = append(g.players[g.turn].Pieces, newPiece)

				// add the new piece to the board
				usePlayerAction(g)
			} else {
				// selected piece owned by the current player wants to move, and tile is empty

				// move the position of the piece for the player that is currently selected
				selectedPlayerId := g.board.Tiles[selectedX][selectedY].Piece.PlayerIndex
				for i, piece := range g.players[selectedPlayerId].Pieces {
					if piece.Position.X == selectedX && piece.Position.Y == selectedY {
						g.players[selectedPlayerId].Pieces[i].Position = Position{targetX, targetY}
						break
					}
				}

				g.HighlightedTile.X = targetX
				g.HighlightedTile.Y = targetY

				g.SelectedTile = g.HighlightedTile
				usePlayerAction(g)
			}

		} else {
			// there is a piece on the target tile

			// check if the piece on the target tile is owned by the current player
			if g.board.Tiles[targetX][targetY].Piece.PlayerIndex == g.players[g.turn].PlayerIndex {
				// target tile is owned by the current player, and there is a piece on the selected tile
				if g.board.Tiles[selectedX][selectedY].Piece.Value == 6 {
					// selected piece is outpost and max value
					if g.board.Tiles[targetX][targetY].Piece.Value < 6 {
						// add one to the target pice value
						targetPiece := g.board.Tiles[targetX][targetY].Piece
						targetPlayerId := targetPiece.PlayerIndex
						for i, piece := range g.players[targetPlayerId].Pieces {
							if piece.Position.X == targetX && piece.Position.Y == targetY {
								//update the player with the new value
								g.players[targetPlayerId].Pieces[i].Value++
								break
							}
						}
						usePlayerAction(g)
					} else {
						// Can not add more than 6 to a picece
						g.InvalidTile = Position{targetX, targetY}
					}
				} else {
					// piece is unit, and not max value, is trying to move to a tile that already has a piece on it owned by the current player

					combinedValue := g.board.Tiles[targetX][targetY].Piece.Value + g.board.Tiles[selectedX][selectedY].Piece.Value

					switch combinedValue {
					case 0, 1:
						// error, should not be able to make value 0 or 1 by combining two pieces
						log.Fatalf("error: invalid value when combining own pieces: %v\n", combinedValue)
						g.InvalidTile = Position{targetX, targetY}
					case 2, 3, 4, 5, 6:
						// combine the two pieces and remove the selected piece from the board and players pieces array

						// find and update the target tile with the new value
						targetPlayerId := g.board.Tiles[targetX][targetY].Piece.PlayerIndex
						for i, piece := range g.players[targetPlayerId].Pieces {
							if piece.Position.X == targetX && piece.Position.Y == targetY {
								g.players[targetPlayerId].Pieces[i].Value = combinedValue
								break // ensures that when the piece is removed, the last itteration will not cause an out of range exception
							}
						}
						// now remove the selected piece from the board
						selectedPlayerId := g.board.Tiles[selectedX][selectedY].Piece.PlayerIndex
						for i, piece := range g.players[selectedPlayerId].Pieces {
							if piece.Position.X == selectedX && piece.Position.Y == selectedY {
								g.players[selectedPlayerId].Pieces = append(g.players[selectedPlayerId].Pieces[:i], g.players[selectedPlayerId].Pieces[i+1:]...)
								break // ensures that when the piece is removed, the last itteration will not cause an out of range exception
							}
						}

						// update the user with the new highlighted and selected tiles
						g.HighlightedTile = Position{targetX, targetY}
						g.SelectedTile = g.HighlightedTile

						usePlayerAction(g)
					case 7, 8, 9, 10, 11:

						// only combine pieces where the targert pices, is no max value
						if g.board.Tiles[targetX][targetY].Piece.Value < 6 {
							// set the target tile to a value of 6, and set the selected piece to a value of newValue-6

							targetPlayerId := g.board.Tiles[targetX][targetY].Piece.PlayerIndex
							for i, piece := range g.players[targetPlayerId].Pieces {
								if piece.Position.X == targetX && piece.Position.Y == targetY {
									g.players[targetPlayerId].Pieces[i].Value = 6
									break // ensures that when the piece is removed, the last itteration will not cause an out of range exception
								}
							}

							// new set the selected pieces value to the remaining value
							selectedPlayerId := g.board.Tiles[selectedX][selectedY].Piece.PlayerIndex
							for i, piece := range g.players[selectedPlayerId].Pieces {
								if piece.Position.X == selectedX && piece.Position.Y == selectedY {
									g.players[selectedPlayerId].Pieces[i].Value = combinedValue - 6
									break // ensures that when the piece is removed, the last itteration will not cause an out of range exception
								}
							}

							usePlayerAction(g)
						} else {
							g.InvalidTile = Position{targetX, targetY}
						}
					case 12:
						// invalid move, as both pices are at the maximum value
						g.InvalidTile = Position{targetX, targetY}
					default:
						// error, should not be able to make a value greater than 12, of negative values, by combining two pieces
						log.Fatalf("error: invalid value when combining own pieces: %v\n", combinedValue)
						g.InvalidTile = Position{targetX, targetY}
					}

				}
			} else {
				// target tile is not owned by the current player, and there is a piece on the selected tile

				playerPiece := g.board.Tiles[selectedX][selectedY].Piece
				targetPiece := g.board.Tiles[targetX][targetY].Piece

				switch playerPiece.Value {
				case 1, 3, 5:
					// piece is a gatherer, and can not attack other pieces
					g.InvalidTile = Position{targetX, targetY}
				case 2, 4:
					// pieces is a soldier, and can attack other pieces

					if targetPiece.Value == playerPiece.Value {
						// remove both pieces from the board and players pieces array

						targetPlayerId := g.board.Tiles[targetX][targetY].Piece.PlayerIndex
						for i, piece := range g.players[targetPlayerId].Pieces {
							if piece.Position.X == targetX && piece.Position.Y == targetY {
								g.players[targetPlayerId].Pieces = append(g.players[targetPlayerId].Pieces[:i], g.players[targetPlayerId].Pieces[i+1:]...)
								break // ensures that when the piece is removed, the last itteration will not cause an out of range exception
							}
						}

						// find and remove the selected piece from the players pieces array
						selectedPlayerId := g.board.Tiles[selectedX][selectedY].Piece.PlayerIndex
						for i, piece := range g.players[selectedPlayerId].Pieces {
							if piece.Position.X == selectedX && piece.Position.Y == selectedY {
								g.players[selectedPlayerId].Pieces = append(g.players[selectedPlayerId].Pieces[:i], g.players[selectedPlayerId].Pieces[i+1:]...)
								break // ensures that when the piece is removed, the last itteration will not cause an out of range exception
							}
						}

						// update the user with the new highlighted and selected tiles
						g.HighlightedTile = Position{targetX, targetY}
						g.SelectedTile = Position{-1, -1} //selected piece has been removed from the board

						usePlayerAction(g)

					} else if targetPiece.Value < playerPiece.Value {
						// remove the targets piece from the board and players pieces array
						// move the players piece to the target tile with reduced value by targets value

						// ensures that when the piece is removed, the last itteration will not cause an out of range exception
						targetPlayerId := g.board.Tiles[targetX][targetY].Piece.PlayerIndex
						removePieceFromPlayer(g, targetPlayerId, targetX, targetY)

						// get the selected piece index in the player pieces list
						selectedPlayerId := g.board.Tiles[selectedX][selectedY].Piece.PlayerIndex
						pieceIndex := findPlayerPieceIndex(g, g.players[selectedPlayerId].Pieces, selectedX, selectedY)

						// move the player pieces to the target tile
						movePieceToTile(g, selectedPlayerId, pieceIndex, targetX, targetY)
						// set the new value of the piece
						remaintingValue := g.players[selectedPlayerId].Pieces[pieceIndex].Value - targetPiece.Value
						setPieceValue(g, selectedPlayerId, pieceIndex, remaintingValue)

						g.HighlightedTile = Position{targetX, targetY}

						g.SelectedTile = g.HighlightedTile
						usePlayerAction(g)
					} else {
						// targetPices values is larger, and will absorb the selected pieces value

						// Reduce the value of the target piece by the players piece value
						targetPlayerId := g.board.Tiles[targetX][targetY].Piece.PlayerIndex
						for i, piece := range g.players[targetPlayerId].Pieces {
							if piece.Position.X == targetX && piece.Position.Y == targetY {
								g.players[targetPlayerId].Pieces[i].Value -= g.board.Tiles[selectedX][selectedY].Piece.Value
								break // ensures that when the piece is removed, the last itteration will not cause an out of range exception
							}
						}

						// find and remove the selected piece from the players pieces array
						selectedPlayerId := g.board.Tiles[selectedX][selectedY].Piece.PlayerIndex
						for i, piece := range g.players[selectedPlayerId].Pieces {
							if piece.Position.X == selectedX && piece.Position.Y == selectedY {
								g.players[selectedPlayerId].Pieces = append(g.players[selectedPlayerId].Pieces[:i], g.players[selectedPlayerId].Pieces[i+1:]...)
								break // ensures that when the piece is removed, the last itteration will not cause an out of range exception
							}
						}

						g.HighlightedTile = Position{targetX, targetY}

						// update the user with the new highlighted and selected tiles
						g.SelectedTile = Position{-1, -1}

						usePlayerAction(g)
					}
				case 6:
					// piece is a outpost, and can only attack with a value of 1

					if targetPiece.Value == 1 {
						// remove the target piece from the player pieces array
						//subtract one as the player id is 1 based, and the players array is zero based

						// find the piece from the target player pieces array, and remove the target piece from the array
						targetPlayerId := g.board.Tiles[targetX][targetY].Piece.PlayerIndex
						for i, piece := range g.players[targetPlayerId].Pieces {
							if piece.Position.X == targetX && piece.Position.Y == targetY {
								g.players[targetPlayerId].Pieces = append(g.players[targetPlayerId].Pieces[:i], g.players[targetPlayerId].Pieces[i+1:]...)
								break // ensures that when the piece is removed, the last itteration will not cause an out of range exception
							}
						}

						usePlayerAction(g)
					} else {
						// target piece will reduce in value by 1 by being attacked by the outpost

						// Find the index of the attacked piece in the player's pieces array
						targetPlayerId := g.board.Tiles[targetX][targetY].Piece.PlayerIndex
						for i, piece := range g.players[targetPlayerId].Pieces {
							if piece.Position.X == targetX && piece.Position.Y == targetY {
								g.players[targetPlayerId].Pieces[i].Value = g.players[targetPlayerId].Pieces[i].Value - 1
								break
							}
						}

						usePlayerAction(g)
					}
				}
			}
		}
	}
}

func removePieceFromPlayer(g *Game, playerId int, positionX int, positionY int) {
	for i, piece := range g.players[playerId].Pieces {
		if piece.Position.X == positionX && piece.Position.Y == positionY {
			g.players[playerId].Pieces = append(g.players[playerId].Pieces[:i], g.players[playerId].Pieces[i+1:]...)
			break
		}
	}
}

func setPieceValue(g *Game, payerId int, pieceIndex int, newValue int) {
	g.players[payerId].Pieces[pieceIndex].Value = newValue
}

func movePieceToTile(g *Game, playerId int, pieceIndex int, positionX int, positionY int) {
	g.players[playerId].Pieces[pieceIndex].Position = Position{positionX, positionY}
}

// findPlayerPieceIndex, loop thought the player, and find the index of the desired piece
func findPlayerPieceIndex(g *Game, playerPieces []Piece, positionX int, positionY int) int {
	for i, piece := range playerPieces {
		if piece.Position.X == positionX && piece.Position.Y == positionY {
			return i
		}
	}
	return -1
}

func usePlayerAction(g *Game) {

	// check if the player has only 1 action left, if so, the last action was just used
	if g.players[g.turn].Actions <= 1 {
		g.players[g.turn].Actions = 0

		log.Printf("End Turn, %v has 0 actions remaining", g.players[g.turn].Name)
		if g.turn == (len(g.players) - 1) {
			g.turn = 0
		} else {
			g.turn++
		}
		log.Printf("It is now %v's turn", g.players[g.turn].Name)

		// Reset the actions for the new player's turn
		updatePlayerActions(g)

		if g.players[g.turn].Actions <= 0 {
			// the current player has no actiosn left due to no pieces or no way to generate actions

			// check if active players is only one
			if len(g.players) == 1 {
				// only one player left, so end the game
				log.Printf("Game over! Player %s has won the game", g.players[0].Name)
				g.GameOver = true

				//allow the playeer to play continuously
				g.turn = 0
				updatePlayerActions(g)
			} else {
				// next player started with 0 actions and will be skipped. call own funciton to try following player
				usePlayerAction(g)
			}
		}
	} else {
		// Decrease the current player's actions by 1
		g.players[g.turn].Actions--
	}
}

func updatePlayerActions(g *Game) {

	// clear up from the previous players turn
	g.SelectedTile = Position{X: -1, Y: -1}

	log.Printf("Setting player action for new turn: ")

	// loop through the players pieces and set their actions based on value
	for i, piece := range g.players[g.turn].Pieces {

		if i == 0 {
			// set the hightlighted tile to be the oldest pieces of the player
			g.HighlightedTile = piece.Position
		}

		switch piece.Value {
		case 1:
			g.players[g.turn].Actions += 1
			log.Printf(", %v +1", piece.Value)
		case 3:
			g.players[g.turn].Actions += 2
			log.Printf(", %v +2", piece.Value)

		case 5, 6:
			g.players[g.turn].Actions += 3
			log.Printf(", %v +3", piece.Value)
		case 2, 4:
			log.Printf(", %v -", piece.Value)
		default:
			log.Printf(", %v ?", piece.Value)
		}
	}

	// printf the current players name and number of actions left and number of pieces they have
	log.Printf("Player %s starts their turn with %d actions and %d pieces", g.players[g.turn].Name, g.players[g.turn].Actions, len(g.players[g.turn].Pieces))
}

// Update proceeds the game state. Update is called every frame (1/60[s] by default).
func (g *Game) Update() error {
	// List of keys to check
	keys := []ebiten.Key{
		ebiten.KeyEscape,
		ebiten.KeyEnter,
		ebiten.KeySpace,
		ebiten.KeyArrowUp,
		ebiten.KeyArrowDown,
		ebiten.KeyArrowLeft,
		ebiten.KeyArrowRight,
	}

	for _, key := range keys {
		if ebiten.IsKeyPressed(key) {
			// If the key is pressed and it was not pressed in the previous frame, log it
			if !g.keyStates[key] {
				g.keyStates[key] = true

				// function to handle key presses
				switch key {
				case ebiten.KeyEscape:
					log.Println("esc")

					// menue is game state 1, and can only be shown when game is running in state 0
					if g.gameState == 0 {
						// menue state
						g.gameState = 1
					} else if g.gameState == 1 {
						// play state
						g.gameState = 0
					}

				case ebiten.KeyEnter:
					log.Println("enter")
					//next players turn and reset if all players have moved

					if g.gameState == 0 {
						if g.turn == (len(g.players) - 1) {
							g.turn = 0
						} else {
							g.turn++
						}

						updatePlayerActions(g)
					}

				case ebiten.KeySpace:
					log.Println("space")
					// The SelectedTile already highlighted, deselect it, else set
					if g.gameState == 0 {
						if g.SelectedTile.X == -1 && g.SelectedTile.Y == -1 {
							// is deselected, so automatically set

							// get the pice one the selected tile, and see if belongs to current player
							if g.board.Tiles[g.HighlightedTile.X][g.HighlightedTile.Y].Piece.PlayerIndex == g.players[g.turn].PlayerIndex {
								// todo, only select if the active players peice. else flash red
								log.Println("\tSelected tile")
								g.SelectedTile = g.HighlightedTile
							} else {
								// log piece at tile
								log.Printf("\tNot belonging to current player")
							}
						} else {
							// is selected, check if selecting same tile, to deselect it
							if g.SelectedTile == g.HighlightedTile {
								g.SelectedTile = Position{X: -1, Y: -1}
							} else {
								// is selected, different tile so select it.

								// NB behavior to be revised as this should find a path to the newly selected tile
								// then see if is a valid move to move the selected piece to that location etc
								g.SelectedTile = Position{X: -1, Y: -1}
							}
						}
					} else if g.gameState == 1 {
						// pause menue
						switch g.uiMenueSelectedButton {
						case 0:
							log.Println("Resume")
							// resume play state
							g.gameState = 0
						case 1:
							log.Println("New Game")
							// confirmation to make new game
							g.gameState = 2
						case 2:
							log.Println("Load")
						case 3:
							log.Println("Settings")
						case 4:
							log.Println("Save")
						case 5:
							log.Println("Exit")
						}
					} else if g.gameState == 2 {
						// new game confirmation
						if g.uiNewGameConfirmation {
							// yes start new game
							g.gameState = 3
							g.uiMenueSelectedButton = 0
							g.uiNewGameConfirmation = false

							//reset game variables
							g.GameOver = false
							g.SelectedTile = Position{-1, -1}
							g.HighlightedTile = Position{-1, -1}
							g.turn = 0

						} else {
							// reset the menue to
							g.gameState = 1
							g.uiMenueSelectedButton = 0
						}
					} else if g.gameState == 3 {
						// new game screen
						if g.uiStartNewGameButton {
							// Try to start new game button pressed

							numberOfPlayers := 0
							for _, section := range g.uiNewGameSectionPlayer {
								if section != -1 {
									numberOfPlayers++
								}
							}
							// only start a new game if at least one player has been selected
							if numberOfPlayers > 0 {
								// start new game
								g.players = createPlayers(g.uiNewGameSectionPlayer)
								setPiecesOnBoardFromPlayers(g)
								updatePlayerActions(g)
								g.gameState = 0
							}
						} else {

							// check if player assigned to section, toggle next player in, if empty
							if g.uiNewGameSectionPlayer[g.uiNewGameSectionHighlighted] == -1 {
								//find out how many players are currently on the board
								numberOfPlayers := 1
								for _, section := range g.uiNewGameSectionPlayer {
									if section != -1 {
										numberOfPlayers++
									}
								}
								g.uiNewGameSectionPlayer[g.uiNewGameSectionHighlighted] = numberOfPlayers
							} else {
								g.uiNewGameSectionPlayer[g.uiNewGameSectionHighlighted] = -1

								// make sure that there is no missing id's in the player position, but maintain the relative player order
								type Section struct {
									origonalIndex int
									playerId      int
								}

								// create the empty data for the helper struct, and defualt all places to empty
								playerPositions := make([]Section, 4)
								for i := range playerPositions {
									playerPositions[i].playerId = -1
								}

								for i, section := range g.uiNewGameSectionPlayer {
									if section != -1 { // Only consider valid player IDs
										playerPositions[section-1] = Section{i, section}
									}
								}

								numberOfPlayers := 0
								for _, section := range g.uiNewGameSectionPlayer {
									if section != -1 {
										numberOfPlayers++
									}
								}

								for i := 1; i <= numberOfPlayers; {
									hasIdex := false
									for p := range g.uiNewGameSectionPlayer {
										if playerPositions[p].playerId == i {
											hasIdex = true
											break
										}
									}
									// check if i has been found, else reduce the value of the others to close the missing index
									if hasIdex {
										// player i has been found in SectionPlayer, check next index
										i++
									} else {
										//reduce all the player ids, after the corected ones of the section by 1
										for p := (i - 1); p < len(playerPositions); p++ {
											if playerPositions[p].playerId > 1 {
												playerPositions[p].playerId = playerPositions[p].playerId - 1
											}
										}
									}
								}
								//Assign the section array back to the g object
								for i := range playerPositions {
									// only set the ones that have a player value
									if playerPositions[i].playerId != -1 {
										g.uiNewGameSectionPlayer[playerPositions[i].origonalIndex] = playerPositions[i].playerId
									}
								}
							}
						}
					}
				case ebiten.KeyArrowLeft:
					log.Println("left")
					if g.gameState == 0 {
						// if the highlighter is not at the left of the board, move it left
						if g.HighlightedTile.X > 0 {
							handleTileMove(g, -1, 0)
						}
					} else if g.gameState == 2 {
						// new game confirmation
						g.uiNewGameConfirmation = !g.uiNewGameConfirmation
					} else if g.gameState == 3 {
						// move the highlighted based on the arrow keys
						if g.uiNewGameSectionHighlighted == 1 || g.uiNewGameSectionHighlighted == 3 {
							g.uiNewGameSectionHighlighted = g.uiNewGameSectionHighlighted - 1
						}
					}
				case ebiten.KeyArrowRight:
					log.Println("right")
					if g.gameState == 0 {
						// if the highlighter is not at the right of the board, move it right
						if g.HighlightedTile.X < g.board.Width {
							handleTileMove(g, 1, 0)
						}
					} else if g.gameState == 2 {
						// new game confirmation
						g.uiNewGameConfirmation = !g.uiNewGameConfirmation
					} else if g.gameState == 3 {
						if g.uiNewGameSectionHighlighted == 0 || g.uiNewGameSectionHighlighted == 2 {
							g.uiNewGameSectionHighlighted = g.uiNewGameSectionHighlighted + 1
						}
					}
				case ebiten.KeyArrowUp:
					log.Println("up")
					if g.gameState == 0 {
						// if the highlighter is not at the top of the board, move it up
						if g.HighlightedTile.Y > 0 {
							handleTileMove(g, 0, -1)
						}
					} else if g.gameState == 1 {
						// pause menue
						if g.uiMenueSelectedButton < 0 {
							g.uiMenueSelectedButton = 0
						} else if g.uiMenueSelectedButton > g.uiMenueButtonNumber {
							g.uiMenueSelectedButton = g.uiMenueButtonNumber
						} else {
							g.uiMenueSelectedButton--
						}
					} else if g.gameState == 3 {
						if g.uiNewGameSectionHighlighted == 2 || g.uiNewGameSectionHighlighted == 3 {
							g.uiNewGameSectionHighlighted = g.uiNewGameSectionHighlighted - 2
						} else if g.uiStartNewGameButton {
							// new game button is currently selected, and now go back to selctions
							g.uiStartNewGameButton = false
							g.uiNewGameSectionHighlighted = 2
						}
					}
				case ebiten.KeyArrowDown:
					log.Println("down")
					if g.gameState == 0 {
						// if the highlighter is not at the bottom of the board, move it down
						if g.HighlightedTile.Y < g.board.Height {
							handleTileMove(g, 0, 1)
						}
					} else if g.gameState == 1 {
						// pause menue
						if g.uiMenueSelectedButton < 0 {
							g.uiMenueSelectedButton = 0
						} else if g.uiMenueSelectedButton > g.uiMenueButtonNumber {
							g.uiMenueSelectedButton = g.uiMenueButtonNumber
						} else {
							g.uiMenueSelectedButton++
						}
					} else if g.gameState == 3 {
						if g.uiNewGameSectionHighlighted == 0 || g.uiNewGameSectionHighlighted == 1 {
							g.uiNewGameSectionHighlighted = g.uiNewGameSectionHighlighted + 2
						} else if g.uiNewGameSectionHighlighted == 2 || g.uiNewGameSectionHighlighted == 3 {
							// wanting to go to the new game button
							g.uiNewGameSectionHighlighted = -1
							g.uiStartNewGameButton = true
						}
					}
				}
			}

			if g.gameState == 0 {
				// when a key is pressed, make sure that the board has been updated with the latest state of the player positions
				// clear the board of peices and re-add them to the board
				clearPiecesFromBoard(g)
				// update the board with the piece position
				setPiecesOnBoardFromPlayers(g)
			}

		} else {
			// If the key is not pressed, reset its state
			g.keyStates[key] = false
		}
	}

	return nil
}

func clearPiecesFromBoard(g *Game) {
	for x, row := range g.board.Tiles {
		for y := range row {
			g.board.Tiles[x][y].Piece = Piece{}
		}
	}
}

func setPiecesOnBoardFromPlayers(g *Game) {
	for _, player := range g.players {
		for _, piece := range player.Pieces {

			g.board.Tiles[piece.Position.X][piece.Position.Y].Piece = piece
		}
	}
}

// Draw draws the game screen. Draw is called every frame (1/60[s] by default).
func (g *Game) Draw(screen *ebiten.Image) {

	// setup the font for text
	textSource, err := text.NewGoTextFaceSource(bytes.NewReader(fonts.MPlus1pRegular_ttf))
	if err != nil {
		log.Fatal(err)
	}

	if g.gameState == 0 {
		// playing game state

		//loop through the tiles of the board
		for x := 0; x <= g.board.Width; x++ {
			for y := 0; y <= g.board.Height; y++ {
				xPos := x * g.board.TileSize
				yPos := y * g.board.TileSize

				// Draw the board checkerboard black and white squares filled
				if (x+y)%2 == 0 {
					// Using board.TileSize as the size of each square draw a black square
					vector.DrawFilledRect(screen, float32(xPos), float32(yPos),
						float32(g.board.TileSize), float32(g.board.TileSize), color.White, true)
				} else {
					vector.DrawFilledRect(screen, float32(xPos), float32(yPos),
						float32(g.board.TileSize), float32(g.board.TileSize), color.Black, true)
				}
			}
		}

		// drawImage of yellow box on highlighter position// there is always a highlighted tile
		highlightedBox := ebiten.NewImage(g.board.TileSize, g.board.TileSize)
		highlightedBox.Fill(color.RGBA{0xff, 0xff, 0x00, 0xff})
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(g.HighlightedTile.X*g.board.TileSize), float64(g.HighlightedTile.Y*g.board.TileSize))
		screen.DrawImage(highlightedBox, op)

		// Is there a Invalid Tile
		if g.InvalidTile.X != -1 && g.InvalidTile.Y != -1 {
			invalidBox := ebiten.NewImage(g.board.TileSize, g.board.TileSize)
			invalidBox.Fill(color.RGBA{0xff, 0x00, 0x00, 0xff})
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(float64(g.InvalidTile.X*g.board.TileSize), float64(g.InvalidTile.Y*g.board.TileSize))
			screen.DrawImage(invalidBox, op)
			g.InvalidTile = Position{-1, -1}
		}

		// Is there a selected Tile
		if g.SelectedTile.X != -1 && g.SelectedTile.Y != -1 {
			// drawImage of green box on selected position
			var boaderSize = 5
			selectedBox := ebiten.NewImage(g.board.TileSize-(boaderSize*2), g.board.TileSize-(boaderSize*2))
			selectedBox.Fill(color.RGBA{0x00, 0xff, 0x00, 0xff})
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(float64(g.SelectedTile.X*g.board.TileSize+boaderSize), float64(g.SelectedTile.Y*g.board.TileSize+boaderSize))
			screen.DrawImage(selectedBox, op)
		}

		for _, player := range g.players {
			for _, piece := range player.Pieces {
				// update the board with the piece position
				xPos := piece.Position.X * g.board.TileSize
				yPos := piece.Position.Y * g.board.TileSize

				pieceBox := ebiten.NewImage(g.board.TileSize/2, g.board.TileSize/2)
				pieceBox.Fill(g.board.Tiles[piece.Position.X][piece.Position.Y].Piece.Color)

				// Draw the Piece box for player colour
				op := &ebiten.DrawImageOptions{}
				op.GeoM.Translate(float64(xPos+g.board.TileSize/4), float64(yPos+g.board.TileSize/4))
				screen.DrawImage(pieceBox, op)

				// Draw the Piece value
				opT := &text.DrawOptions{}
				opT.GeoM.Translate(float64(xPos+36), float64(yPos+28)) // note had to manually find the center based on 18 as the font size
				opT.ColorScale.ScaleWithColor(color.White)
				text.Draw(screen, fmt.Sprint(piece.Value), &text.GoTextFace{
					Source: textSource,
					Size:   18,
				}, opT)
			}
		}

		// Draw the Text for the Player Turns
		uiPlayerStatusOp := &text.DrawOptions{}
		uiPlayerStatusOp.GeoM.Translate(20, 660)
		uiPlayerStatusOp.ColorScale.ScaleWithColor(color.White)
		text.Draw(screen, fmt.Sprintf("Player %v, has %v remaing",
			g.players[g.turn].Name, g.players[g.turn].Actions), &text.GoTextFace{
			Source: textSource,
			Size:   18,
		}, uiPlayerStatusOp)

		// Draw the text for basic instructions
		uiControllsOp := &text.DrawOptions{}
		uiControllsOp.GeoM.Translate(20, 680)
		tutorialMsg := "Controlles: 'space' select piece 'arow keys' move pieces"
		uiControllsOp.ColorScale.ScaleWithColor(color.White)
		text.Draw(screen, fmt.Sprint(tutorialMsg), &text.GoTextFace{
			Source: textSource,
			Size:   18,
		}, uiControllsOp)

	} else if g.gameState == 1 {
		// menue state
		uiBorder := 50
		uiButtonBorder := 20
		uiSize := Position{g.screenSize.X - (uiBorder * 2), g.screenSize.Y - (uiBorder * 2)}
		uiButtonHeight := 80 // (g.screenSize.Y - (2 * uiButtonBorder)) / (g.uiMenueButtonNumber + uiButtonBorder)
		uiButtonWidth := uiSize.X - (uiButtonBorder * 2)
		uiBackgroundColor := color.RGBA{0x55, 0x55, 0x55, 0x55}
		uiButtonColor := color.RGBA{0x33, 0x33, 0x33, 0xff}
		uiButtonHighlightColor := color.RGBA{0x88, 0x88, 0x88, 0xff}
		buttonLabels := []string{"Resume", "New Game", "Load", "Settings", "Save", "Exit"}

		// Draw the ui menue background box
		menueBox := ebiten.NewImage(uiSize.X, uiSize.Y)
		menueBox.Fill(uiBackgroundColor)
		menueDo := &ebiten.DrawImageOptions{}
		menueDo.GeoM.Translate(float64(uiBorder), float64(uiBorder))
		screen.DrawImage(menueBox, menueDo)

		buttonNextPosition := uiBorder - uiButtonHeight
		for i := 0; i <= g.uiMenueButtonNumber; i++ {
			buttonNextPosition += uiButtonBorder + uiButtonHeight
			text := buttonLabels[i]

			if g.uiMenueSelectedButton == i {
				drawMenueButton(screen, uiBorder+uiButtonBorder, buttonNextPosition, uiButtonWidth, uiButtonHeight, uiButtonHighlightColor, textSource, text)

			} else {
				drawMenueButton(screen, uiBorder+uiButtonBorder, buttonNextPosition, uiButtonWidth, uiButtonHeight, uiButtonColor, textSource, text)
			}

		}
	} else if g.gameState == 2 {
		// new game confirmaiton
		uiMenueWidth := 400
		uiMenueHeight := 140
		uiStartX := (g.screenSize.X / 2) - (uiMenueWidth / 2)
		uiStarty := (g.screenSize.Y / 2) - (uiMenueHeight / 2)
		uiButtonBorder := 25
		uiBackgroundColor := color.RGBA{0x55, 0x55, 0x55, 0x55}
		uiButtonColor := color.RGBA{0x33, 0x33, 0x33, 0xff}
		uiButtonHighlightColor := color.RGBA{0x88, 0x88, 0x88, 0xff}

		// Draw the ui menue background box
		menueBox := ebiten.NewImage(uiMenueWidth, uiMenueHeight)
		menueBox.Fill(uiBackgroundColor)
		menueDo := &ebiten.DrawImageOptions{}
		menueDo.GeoM.Translate(float64(uiStartX), float64(uiStarty))
		screen.DrawImage(menueBox, menueDo)

		uiButtonWidth := (uiMenueWidth - (uiButtonBorder * 3)) / 2
		uiButtonHeight := uiMenueHeight - (uiButtonBorder * 2)
		uiButtonStartX := uiStartX + uiButtonBorder
		uiButtonStartY := uiStarty + uiButtonBorder
		if g.uiNewGameConfirmation {
			drawMenueButton(screen, uiButtonStartX, uiButtonStartY, uiButtonWidth, uiButtonHeight, uiButtonColor, textSource, "No")
			drawMenueButton(screen, uiButtonStartX+uiButtonWidth+uiButtonBorder, uiButtonStartY, uiButtonWidth, uiButtonHeight, uiButtonHighlightColor, textSource, "Yes")
		} else {
			drawMenueButton(screen, uiButtonStartX, uiButtonStartY, uiButtonWidth, uiButtonHeight, uiButtonHighlightColor, textSource, "No")
			drawMenueButton(screen, uiButtonStartX+uiButtonWidth+uiButtonBorder, uiButtonStartY, uiButtonWidth, uiButtonHeight, uiButtonColor, textSource, "Yes")
		}
	} else if g.gameState == 3 {
		// new game creation screen
		uiBorder := 50
		uiStartGameAreaHeight := 150
		uiMessageBoxAreaHeight := 150
		uiSectionWidth := (g.screenSize.X / 2) - (uiBorder / 2) - uiBorder
		uiSectionHeight := ((g.screenSize.Y - uiStartGameAreaHeight - uiMessageBoxAreaHeight) / 2) - (uiBorder / 2) - uiBorder
		uiMessageBoxStartY := (uiBorder * (2 + 1)) + (uiSectionHeight * 2)
		uiExcludedColor := color.RGBA{0x11, 0x11, 0x11, 0xff}
		uiIncludedColor := color.RGBA{0x44, 0x44, 0x44, 0xff}
		uiHighlightColor := color.RGBA{0x99, 0x99, 0x99, 0xff}

		index := 0
		for c := 0; c < 2; c++ {
			for r := 0; r < 2; r++ {
				startX := (uiBorder * (r + 1)) + (uiSectionWidth * r)
				startY := (uiBorder * (c + 1)) + (uiSectionHeight * c)

				section := ebiten.NewImage(uiSectionWidth, uiSectionHeight)
				// set the coresponding color, depending on if the section is included, excluded, or selected
				if g.uiNewGameSectionPlayer[index] == -1 {
					section.Fill(uiExcludedColor)
				} else {
					section.Fill(uiIncludedColor)
				}

				if g.uiNewGameSectionHighlighted == index {
					section.Fill(uiHighlightColor)
				}

				sectionDo := &ebiten.DrawImageOptions{}
				sectionDo.GeoM.Translate(float64(startX), float64(startY))
				screen.DrawImage(section, sectionDo)

				// draw text on section needs to be after drawing of the section
				if g.uiNewGameSectionPlayer[index] != -1 {
					op := &text.DrawOptions{}
					op.GeoM.Translate(float64(startX+(uiSectionWidth/3)), float64(startY+18))
					op.ColorScale.ScaleWithColor(color.White)
					text.Draw(screen, fmt.Sprintf("player%v", g.uiNewGameSectionPlayer[index]), &text.GoTextFace{
						Source: textSource,
						Size:   36,
					}, op)
				}
				index++
			}
		}

		// draw the message box section
		op := &text.DrawOptions{}
		op.GeoM.Translate(float64(uiBorder), float64(uiMessageBoxStartY))
		op.ColorScale.ScaleWithColor(color.White)
		text.Draw(screen, "Use the arrow keys to navigate", &text.GoTextFace{
			Source: textSource,
			Size:   36,
		}, op)
		op = &text.DrawOptions{}
		op.GeoM.Translate(float64(uiBorder), float64(uiMessageBoxStartY+36+8))
		op.ColorScale.ScaleWithColor(color.White)
		text.Draw(screen, "Spacebar to toggle players", &text.GoTextFace{
			Source: textSource,
			Size:   36,
		}, op)

		newGameButton := ebiten.NewImage(g.screenSize.X-(uiBorder*2), uiStartGameAreaHeight-uiBorder)
		if g.uiStartNewGameButton {
			newGameButton.Fill(uiHighlightColor)
		} else {
			newGameButton.Fill(uiIncludedColor)
		}

		stargGameDo := &ebiten.DrawImageOptions{}
		stargGameDo.GeoM.Translate(float64(uiBorder), float64(g.screenSize.Y-uiStartGameAreaHeight))
		screen.DrawImage(newGameButton, stargGameDo)

		op = &text.DrawOptions{}
		op.GeoM.Translate(float64(g.screenSize.X/3), float64(g.screenSize.Y-uiStartGameAreaHeight+18))
		op.ColorScale.ScaleWithColor(color.White)
		text.Draw(screen, "Start Game", &text.GoTextFace{
			Source: textSource,
			Size:   36,
		}, op)
	}

}

func drawMenueButton(screen *ebiten.Image, startX, startY, width, height int, buttonColor color.Color, s *text.GoTextFaceSource, buttonText string) {
	buttonBox := ebiten.NewImage(width, height)
	buttonBox.Fill(buttonColor)
	buttonDo := &ebiten.DrawImageOptions{}
	buttonDo.GeoM.Translate(float64(startX), float64(startY))
	screen.DrawImage(buttonBox, buttonDo)

	// Draw the Piece value
	op := &text.DrawOptions{}
	op.GeoM.Translate(float64(startX+(width/3)), float64(startY+18)) // note had to manually find the center based on 18 as the font size
	op.ColorScale.ScaleWithColor(color.White)
	text.Draw(screen, fmt.Sprint(buttonText), &text.GoTextFace{
		Source: s,
		Size:   36,
	}, op)
}

// Layout takes the outside size (in device-independent pixels) and returns the logical screen size.
// If you don't have to adjust the screen size with the outside size, just return a fixed size.
func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	// note setting to static value will allow resize, where as outside dimentiosn will enforce the size to be the screen size
	return g.screenSize.X, g.screenSize.Y
}

/*
	powershell

$Env:GOOS = "windows"; $Env:GOARCH = "amd64"; go run main.go

$Env:GOOS = "darwin"; $Env:GOARCH = "amd64"; go build -o mac.dmg main.go 		// mac
$Env:GOOS = "linux"; $Env:GOARCH = "arm64"; go build -o android.apk main.go 	// android
$Env:GOOS = "windows"; $Env:GOARCH = "amd64"; go build -o windows.exe main.go 	// windows
$Env:GOOS = "js"; $Env:GOARCH = "wasm"; go build -o browser.wasm main.go 		// browser
*/
func main() {
	g := &Game{
		keyStates:                   make(map[ebiten.Key]bool),
		board:                       createBoard(7, 7, 80), // 8 by 8 tiles
		players:                     createPlayers([]int{-1, 2, 1, -1}),
		turn:                        0,
		HighlightedTile:             Position{-1, -1},
		SelectedTile:                Position{X: -1, Y: -1},
		GameOver:                    false,
		gameState:                   0,
		uiMenueSelectedButton:       0,
		uiMenueButtonNumber:         5,
		uiNewGameConfirmation:       false,
		uiNewGameSectionPlayer:      make([]int, 4),
		uiNewGameSectionHighlighted: 0,
		uiStartNewGameButton:        false,
		screenSize:                  Position{640, 720}, //960, 720
	}

	//setup game
	setPiecesOnBoardFromPlayers(g)
	updatePlayerActions(g)
	for i := range g.uiNewGameSectionPlayer {
		switch i {
		case 0:
			g.uiNewGameSectionPlayer[i] = -1
		case 1:
			g.uiNewGameSectionPlayer[i] = 2
		case 2:
			g.uiNewGameSectionPlayer[i] = 1
		case 3:
			g.uiNewGameSectionPlayer[i] = -1
		}
	}

	ebiten.SetWindowSize(g.screenSize.X, g.screenSize.Y)
	ebiten.SetWindowTitle("Six Divides")
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
