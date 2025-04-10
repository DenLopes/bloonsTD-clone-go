package main

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func NewGame() Game {
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

func NewCamera(screenWidth int32, screenHeight int32) rl.Camera2D {
	return rl.NewCamera2D(
		rl.NewVector2(float32(screenWidth/2), float32(screenHeight/2)),
		rl.NewVector2(float32(screenWidth/2), float32(screenHeight/2)),
		0,
		1.0,
	)
}

func BuildGridPaths(g *Game, invert bool) {
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
