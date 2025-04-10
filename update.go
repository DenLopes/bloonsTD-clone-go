package main

import rl "github.com/gen2brain/raylib-go/raylib"

func UpdateEnemy(enemy *Enemy, delta float32, paths Paths) {
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
