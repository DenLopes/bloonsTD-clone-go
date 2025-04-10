package main

import (
	"fmt"

	"github.com/DenLopes/bloonsTD-clone-go/ui"
	rl "github.com/gen2brain/raylib-go/raylib"
)

func main() {

	game := NewGame()
	BuildGridPaths(&game, false)

	fmt.Println("Cells: ", game.grid.cells)

	// Initialize core
	rl.SetConfigFlags(rl.FlagWindowResizable)
	rl.InitWindow(int32(game.screenWidth), int32(game.screenHeight), "Bloons TD")
	defer rl.CloseWindow()
	rl.SetTargetFPS(60)

	camera := NewCamera(game.screenWidth, game.screenHeight)

	// Main loop
	for !rl.WindowShouldClose() {
		delta := rl.GetFrameTime()

		if rl.IsWindowResized() && !rl.IsWindowFullscreen() {
			game.screenWidth = int32(rl.GetScreenWidth())
			game.screenHeight = int32(rl.GetScreenHeight())
		}

		rootUINode := ui.RootNode(game.screenWidth, game.screenHeight)
		rootUINode.AddChildren(
			ui.TestComponent(),
		)
		ui.CalculateDynamicElements(rootUINode)

		// Handle Inputs
		mouseWorldPos := rl.GetScreenToWorld2D(rl.GetMousePosition(), camera)
		game.hoveredCellKey = getHoveredCellKey(mouseWorldPos, game.grid)
		handleUserInput(mouseWorldPos, &camera, &game)

		// Handle Logic
		for i := range game.enemies {
			if game.enemies[i].dead {
				continue
			}
			UpdateEnemy(&game.enemies[i], delta, game.paths)
		}

		// for i, tower := range game.towers {
		// 	// Tower logic
		// 	fmt.Println(i, tower)
		// }

		// Rendering
		rl.BeginDrawing()
		rl.ClearBackground(rl.Black)
		rl.BeginMode2D(camera)

		//Draw things here:
		DrawGrid(game.grid, &game)
		DrawPaths(game.paths)
		DrawEnemies(game.enemies)
		DrawTowers(game.towers)

		rl.EndMode2D()

		//Draw UI Here
		ui.DrawUI(rootUINode)

		rl.EndDrawing()
	}
}

func handleUserInput(mouseWorldPos rl.Vector2, camera *rl.Camera2D, game *Game) {
	handleCameraInput(mouseWorldPos, camera)
	//handle full screen input
	if rl.IsKeyPressed(rl.KeyEnter) && (rl.IsKeyDown(rl.KeyLeftAlt) || rl.IsKeyDown(rl.KeyRightAlt)) {
		display := rl.GetCurrentMonitor()
		if rl.IsWindowFullscreen() {
			rl.SetWindowSize(int(game.screenWidth), int(game.screenHeight))
		} else {
			rl.SetWindowSize(rl.GetMonitorWidth(display), rl.GetMonitorHeight(display))
		}
		rl.ToggleFullscreen()
	}

	if rl.IsMouseButtonPressed(rl.MouseLeftButton) && game.hoveredCellKey != nil {
		cell := game.grid.cells[*game.hoveredCellKey]
		if !cell.path && !cell.entry && !cell.exit {
			game.selectedCellKey = game.hoveredCellKey
		}
	}
}

func handleCameraInput(mouseWorldPos rl.Vector2, camera *rl.Camera2D) {
	// Pan with right mouse button
	if rl.IsMouseButtonDown(rl.MouseRightButton) {
		delta := rl.GetMouseDelta()
		delta = rl.Vector2Scale(delta, -1/camera.Zoom)
		camera.Target = rl.Vector2Add(camera.Target, delta)
	}

	// Zoom with mouse wheel
	wheel := rl.GetMouseWheelMove()
	if wheel != 0 {
		camera.Zoom += float32(wheel) * 0.1
		if camera.Zoom < 0.1 {
			camera.Zoom = 0.1
		}

		// Adjust target to keep mouse position stable
		newMouseWorldPos := rl.GetScreenToWorld2D(rl.GetMousePosition(), *camera)
		camera.Target = rl.Vector2Add(camera.Target, rl.Vector2Subtract(mouseWorldPos, newMouseWorldPos))
	}
}

func getHoveredCellKey(mouseWorldPos rl.Vector2, grid Grid) *GridKey {
	for key := range grid.cells {
		currentCell := grid.cells[key]
		if mouseWorldPos.X <= (currentCell.worldPosition.X+float32(currentCell.size)) && mouseWorldPos.X >= currentCell.worldPosition.X {
			if mouseWorldPos.Y <= (currentCell.worldPosition.Y+float32(currentCell.size)) && mouseWorldPos.Y >= currentCell.worldPosition.Y {
				result := key
				return &result
			}
		}
	}
	return nil
}
