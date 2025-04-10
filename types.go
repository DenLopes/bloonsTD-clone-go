package main

import rl "github.com/gen2brain/raylib-go/raylib"

const (
	COLS      int = 10
	ROWS      int = 10
	RECT_SIZE int = 100
)

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
