package main

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const COLS int = 10
const ROWS int = 10
const RECT_SIZE int = 100

type Game struct {
	grid            Grid
	enemies         Enemies
	towers          Towers
	screenWidth     int32
	screenHeight    int32
	health          float32
	money           float32
	wave            int32
	paths           Paths
	hoveredCellKey  *GridKey
	selectedCellKey *GridKey
}

type Grid struct {
	sY    int
	sX    int
	cells CellsMap
	entry *Cell
	exit  *Cell
}

type GridKey struct {
	X int32
	Y int32
}

type CellsMap map[GridKey]Cell

type Cell struct {
	worldPosition rl.Vector2
	gridPosition  rl.Vector2
	worldCenter   rl.Vector2
	size          int32
	path          bool
	entry         bool
	exit          bool
	tower         *Tower
}

type Enemy struct {
	radius    float32
	health    float32
	velocity  float32
	position  rl.Vector2
	pathIndex int32
	nodeIndex int32
	dead      bool
}

type Enemies []Enemy

type Tower struct {
	projectile   Projectile
	attackSpeed  float32
	rangeRadius  float32
	lastAttacked float32
	cell         *Cell
}

type Projectile struct {
	radius    float32
	damage    float32
	velocity  float32
	position  rl.Vector2
	direction rl.Vector2
	pierces   int8
}

type Towers []Tower

type PathNode struct {
	worldCenter  rl.Vector2
	gridPosition rl.Vector2
}

type Path []PathNode

type Paths []Path

func main() {

	game := initGameState()
	createGridPaths(&game, false)

	fmt.Println("Cells: ", game.grid.cells)

	// Initialize core
	rl.SetConfigFlags(rl.FlagWindowResizable)
	rl.InitWindow(int32(game.screenWidth), int32(game.screenHeight), "Bloons TD")
	defer rl.CloseWindow()
	rl.SetTargetFPS(60)

	camera := setupCamera(game.screenWidth, game.screenHeight)

	// Main loop
	for !rl.WindowShouldClose() {
		delta := rl.GetFrameTime()
		if rl.IsWindowResized() && !rl.IsWindowFullscreen() {
			game.screenWidth = int32(rl.GetScreenWidth())
			game.screenHeight = int32(rl.GetScreenHeight())
		}

		// Handle Inputs
		mouseWorldPos := rl.GetScreenToWorld2D(rl.GetMousePosition(), camera)
		game.hoveredCellKey = getHoveredCellKey(mouseWorldPos, game.grid)
		handleUserInput(mouseWorldPos, &camera, &game)

		// Handle Logic
		for i := range game.enemies {
			if game.enemies[i].dead {
				continue
			}
			updateEnemy(&game.enemies[i], delta, game.paths)
		}

		// for i, tower := range game.towers {
		// 	// Tower logic
		// 	fmt.Println(i, tower)
		// }
		type Cell struct {
			worldPosition rl.Vector2
			gridPosition  rl.Vector2
			worldCenter   rl.Vector2
			size          int32
			path          bool
			entry         bool
			exit          bool
			tower         *Tower
		}

		// Rendering
		rl.BeginDrawing()
		rl.ClearBackground(rl.Black)
		rl.BeginMode2D(camera)

		//Draw things here:
		drawGrid(game.grid, game)
		drawPaths(game.paths)
		drawEnemies(game.enemies)

		rl.EndMode2D()
		rl.EndDrawing()
	}
}

func initGameState() Game {
	screenWidth := 800
	screenHeight := 600

	gameGrid := Grid{
		sY:    COLS,
		sX:    ROWS,
		cells: CellsMap{},
	}

	centerVector := rl.NewVector2(float32(RECT_SIZE)/2, float32(RECT_SIZE)/2)
	for x := 0; x <= gameGrid.sX; x++ {
		for y := 0; y <= gameGrid.sY; y++ {
			gridPosition := rl.NewVector2(float32(x), float32(y))
			worldPosition := rl.NewVector2(float32(x*RECT_SIZE), float32(y*RECT_SIZE))
			key := GridKey{
				X: int32(x),
				Y: int32(y),
			}
			gameGrid.cells[key] = Cell{
				exit:          false,
				entry:         false,
				path:          false,
				size:          int32(RECT_SIZE),
				gridPosition:  gridPosition,
				worldPosition: worldPosition,
				worldCenter:   rl.Vector2Add(worldPosition, centerVector),
			}
		}
	}

	path := []rl.Vector2{
		rl.NewVector2(0, 1),
		rl.NewVector2(1, 1),
		rl.NewVector2(2, 1),
		rl.NewVector2(3, 1),
		rl.NewVector2(4, 1),
		rl.NewVector2(5, 1),
		rl.NewVector2(5, 2),
		rl.NewVector2(5, 3),
		rl.NewVector2(5, 4),
		rl.NewVector2(5, 5),
		rl.NewVector2(5, 6),
		rl.NewVector2(5, 7),
		rl.NewVector2(5, 8),
		rl.NewVector2(5, 9),
		rl.NewVector2(6, 9),
		rl.NewVector2(7, 9),
		rl.NewVector2(8, 9),
		rl.NewVector2(9, 9),
		rl.NewVector2(9, 10),
	}

	for i, v := range path {
		key := GridKey{
			X: int32(v.X),
			Y: int32(v.Y),
		}
		cell := gameGrid.cells[key]
		if i == 0 {
			fmt.Println(key)
			gameGrid.entry = &cell
			cell.entry = true
		}
		if i == len(path)-1 {
			gameGrid.exit = &cell
			cell.exit = true
		}
		cell.path = true
		gameGrid.cells[key] = cell
	}
	fmt.Println(gameGrid.entry.worldCenter)
	return Game{
		grid: gameGrid,
		enemies: Enemies{
			Enemy{
				radius:    20,
				health:    10,
				position:  gameGrid.entry.worldCenter,
				velocity:  100,
				pathIndex: 0,
				nodeIndex: 0,
			},
		},
		towers:       Towers{},
		screenWidth:  int32(screenWidth),
		screenHeight: int32(screenHeight),
		health:       100,
		money:        1000,
		wave:         1,
	}

}

func setupCamera(screenWidth int32, screenHeight int32) rl.Camera2D {
	return rl.NewCamera2D(
		rl.NewVector2(float32(screenWidth/2), float32(screenHeight/2)),
		rl.NewVector2(float32(screenWidth/2), float32(screenHeight/2)),
		0,
		1.0,
	)
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

func drawGrid(gameGrid Grid, game Game) {
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

func drawEnemies(enemies Enemies) {
	for _, enemie := range enemies {
		if enemie.dead {
			continue
		}
		rl.DrawCircle(int32(enemie.position.X), int32(enemie.position.Y), enemie.radius, rl.Red)
	}
}

func drawGridSquare(worldPosition rl.Vector2, size int32, color rl.Color, borderColor rl.Color) {
	rl.DrawRectangle(int32(worldPosition.X), int32(worldPosition.Y), size, size, borderColor)
	rl.DrawRectangle(int32(worldPosition.X+1), int32(worldPosition.Y+1), size-2, size-2, color)
}

func createGridPaths(g *Game, invert bool) {
	var initialNode PathNode
	if !invert {
		initialNode = PathNode{
			worldCenter:  g.grid.entry.worldCenter,
			gridPosition: g.grid.entry.gridPosition,
		}
	} else {
		initialNode = PathNode{
			worldCenter:  g.grid.exit.worldCenter,
			gridPosition: g.grid.exit.gridPosition,
		}
	}
	initialPath := Path{initialNode}

	// Queue of paths to process
	pathQueue := []Path{initialPath}
	completedPaths := Paths{}
	visited := make(map[GridKey]bool)
	visited[GridKey{X: int32(initialNode.gridPosition.X), Y: int32(initialNode.gridPosition.Y)}] = true

	// BFS loop
	for len(pathQueue) > 0 {
		// Get current path
		currentPath := pathQueue[0] // Get the first element
		pathQueue = pathQueue[1:]   // Remove the first element

		// Get last node in current path
		lastNode := currentPath[len(currentPath)-1]

		// Check if we've reached the exit
		if lastNode.gridPosition == g.grid.exit.gridPosition {
			completedPaths = append(completedPaths, currentPath)
			continue
		}

		// Check all four neighbors
		neighbors := []GridKey{
			{X: int32(lastNode.gridPosition.X), Y: int32(lastNode.gridPosition.Y) - 1}, // Top
			{X: int32(lastNode.gridPosition.X), Y: int32(lastNode.gridPosition.Y) + 1}, // Bottom
			{X: int32(lastNode.gridPosition.X) - 1, Y: int32(lastNode.gridPosition.Y)}, // Left
			{X: int32(lastNode.gridPosition.X) + 1, Y: int32(lastNode.gridPosition.Y)}, // Right
		}

		for _, neighborKey := range neighbors {

			// Check if neighbor exists and is a path
			cell, exists := g.grid.cells[neighborKey]
			if !exists {
				continue // Skip to the next neighbor if this one doesn't exist
			}

			// Then check if it's a path cell
			if !cell.path {
				continue // Skip if it's not a path
			}

			// Finally check if we've already visited it
			if visited[neighborKey] {
				continue // Skip if already visited
			}

			// Mark as visited
			visited[neighborKey] = true

			// Create new path node
			newNode := PathNode{
				worldCenter:  cell.worldCenter,
				gridPosition: cell.gridPosition,
			}

			// Create a copy of the current path and add the new node
			newPath := make(Path, len(currentPath))
			copy(newPath, currentPath)
			newPath = append(newPath, newNode)

			// Add to queue
			pathQueue = append(pathQueue, newPath)
		}
	}

	g.paths = completedPaths
}

func drawPaths(paths Paths) {
	for _, path := range paths {
		for i, node := range path {
			if i == len(path)-1 {
				break
			}
			rl.DrawLineV(node.worldCenter, path[i+1].worldCenter, rl.RayWhite)
		}
	}
}

func updateEnemy(enemy *Enemy, delta float32, paths Paths) {
	walkingTowards := paths[enemy.pathIndex][enemy.nodeIndex+1].worldCenter
	direction := rl.Vector2Normalize(rl.Vector2Subtract(walkingTowards, enemy.position))
	distanceTowayPoint := rl.Vector2Distance(enemy.position, walkingTowards)
	moveAmount := enemy.velocity * delta

	if moveAmount >= distanceTowayPoint {
		enemy.position = walkingTowards
		if enemy.nodeIndex+2 < int32(len(paths[enemy.pathIndex])) {
			enemy.nodeIndex++
		} else {
			enemy.dead = true
		}
	} else {
		movement := rl.Vector2Scale(direction, moveAmount)
		enemy.position = rl.Vector2Add(enemy.position, movement)
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
