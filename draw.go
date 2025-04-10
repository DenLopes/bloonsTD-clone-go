package main

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func DrawPaths(paths Paths) {
	for _, path := range paths {
		for i, node := range path {
			if i == len(path)-1 {
				break
			}
			rl.DrawLineV(node.worldCenter, path[i+1].worldCenter, rl.RayWhite)
		}
	}
}

func DrawGrid(gameGrid Grid, game *Game) {
	for key, cell := range gameGrid.cells {

		fillColor := rl.White
		borderColor := rl.Black

		if cell.path {
			fillColor = rl.Blue
		} else {
			if game.hoveredCellKey != nil && *game.hoveredCellKey == key {
				borderColor = rl.Gold
			}
			if game.selectedCellKey != nil && *game.selectedCellKey == key {
				borderColor = rl.Red
			}
		}
		if cell.entry {
			fillColor = rl.Green
		}
		if cell.exit {
			fillColor = rl.Red
		}

		drawGridSquare(cell.worldPosition, cell.size, fillColor, borderColor)

		rl.DrawText(fmt.Sprint(cell.gridPosition), int32(cell.worldPosition.X+2), int32(cell.worldPosition.Y+2), 10, rl.Black)
	}
}

func DrawEnemies(enemies Enemies) {
	for _, enemie := range enemies {
		if enemie.dead {
			continue
		}
		rl.DrawCircle(int32(enemie.position.X), int32(enemie.position.Y), enemie.radius, rl.Red)
	}
}

func DrawTowers(towers Towers) {
	for _, tower := range towers {
		rl.DrawCircle(int32(tower.cell.worldCenter.X), int32(tower.cell.worldCenter.X), 8, rl.Black)
		rl.DrawCircle(int32(tower.cell.worldCenter.X), int32(tower.cell.worldCenter.X), tower.rangeRadius, rl.Fade(rl.Blue, .2))
	}
}

func drawGridSquare(worldPosition rl.Vector2, size int32, color rl.Color, borderColor rl.Color) {
	rl.DrawRectangle(int32(worldPosition.X), int32(worldPosition.Y), size, size, borderColor)
	rl.DrawRectangle(int32(worldPosition.X+1), int32(worldPosition.Y+1), size-2, size-2, color)
}
