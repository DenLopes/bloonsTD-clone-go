package ui

import rl "github.com/gen2brain/raylib-go/raylib"

func RootNode(screenWidth int32, screenHeight int32) *TreeNode {
	return Node(NodeData{
		Padding: NodeSpacing{
			Left:   60,
			Top:    60,
			Right:  10,
			Bottom: 10,
		},
		Position: rl.NewVector2(0, 0),
		Size: NodeSize{
			Width:  screenWidth,
			Height: screenHeight,
		},
		ChildGap: 10,
	})
}

func TestComponent() *TreeNode {
	return Node(NodeData{
		Size:            NodeSize{Width: 600, Height: 200},
		BackgroundColor: rl.Gray,
		Fit:             false,
		Padding:         NodeSpacing{Top: 10, Left: 10, Right: 10, Bottom: 10},
		ChildGap:        5,
	}).AddChildren(
		Node(NodeData{
			Size:            NodeSize{Width: 82, Height: 82},
			BackgroundColor: rl.Green,
			Fit:             true,
			Padding:         NodeSpacing{Top: 10, Left: 10, Right: 10, Bottom: 10},
			ChildGap:        10,
			Column:          true,
			Grow: NodeSizingBehaviour{
				Height: true,
				Width:  true,
			},
		}),
		Node(NodeData{
			Size:            NodeSize{Width: 82, Height: 82},
			BackgroundColor: rl.Gold,
			Fit:             true,
			Padding:         NodeSpacing{Top: 10, Left: 10, Right: 10, Bottom: 10},
			ChildGap:        10,
			Column:          true,
		}).AddChildren(
			Node(NodeData{
				Size:            NodeSize{Width: 82, Height: 82},
				BackgroundColor: rl.Beige,
				Fit:             true,
				Padding:         NodeSpacing{Top: 10, Left: 10, Right: 10, Bottom: 10},
				ChildGap:        10,
				Column:          true,
				Grow: NodeSizingBehaviour{
					Height: true,
					Width:  true,
				},
			}),
			Node(NodeData{
				Size:            NodeSize{Width: 112, Height: 22},
				BackgroundColor: rl.Red,
				Fit:             true,
				Padding:         NodeSpacing{Top: 10, Left: 10, Right: 10, Bottom: 10},
				ChildGap:        10,
				Column:          true,
			}),
		),
		Node(NodeData{
			Size:            NodeSize{Width: 82, Height: 182},
			BackgroundColor: rl.Blue,
			Fit:             true,
			Padding:         NodeSpacing{Top: 10, Left: 10, Right: 10, Bottom: 10},
			ChildGap:        10,
			Column:          true,
			Grow: NodeSizingBehaviour{
				Height: true,
				Width:  true,
			},
		}),
	)
}
