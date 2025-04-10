package ui

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type NodeSize struct {
	Width  int32
	Height int32
}

type NodeSizingBehaviour struct {
	Width  bool
	Height bool
}

type NodeSpacing struct {
	Left   int32
	Top    int32
	Right  int32
	Bottom int32
}

type NodeData struct {
	// Position
	Position rl.Vector2
	// Style
	FontSize        int32
	Size            NodeSize
	Padding         NodeSpacing
	BackgroundColor rl.Color
	ChildGap        int32
	// Layout
	Column bool
	Fit    bool
	Grow   NodeSizingBehaviour
	Shrink NodeSizingBehaviour
	// Content
	Text  string
	Image rl.Image
}

type TreeNode struct {
	Data     NodeData
	Parent   *TreeNode
	Children []*TreeNode
}

func (parent *TreeNode) AddChildren(newChildren ...*TreeNode) *TreeNode {
	if parent == nil {
		fmt.Println("Warning: AddChildren called on a nil TreeNode pointer.")
		return nil
	}

	if len(newChildren) == 0 {
		return nil
	}

	topOffset := parent.Data.Padding.Top
	leftOffset := parent.Data.Padding.Left
	sumSizes := NodeSize{}
	validChildren := make([]*TreeNode, 0, len(newChildren))
	for _, child := range newChildren {
		if child != nil {
			validChildren = append(validChildren, child)
			child.Parent = parent

			sumSizes.Width += child.Data.Size.Width
			sumSizes.Height += child.Data.Size.Height

			child.Data.Position.X += float32(leftOffset)
			child.Data.Position.Y += float32(topOffset)

			if parent.Data.Column {
				topOffset += child.Data.Size.Height + parent.Data.ChildGap
				parent.Data.Size.Width = max(parent.Data.Size.Width, child.Data.Size.Width)
			} else {
				leftOffset += child.Data.Size.Width + parent.Data.ChildGap
				parent.Data.Size.Height = max(parent.Data.Size.Height, child.Data.Size.Height)
			}
		}
	}

	gapAmount := int32((len(validChildren) - 1) * int(parent.Data.ChildGap))
	if parent.Data.Fit {
		if parent.Data.Column {
			parent.Data.Size.Width += parent.Data.Padding.Right + parent.Data.Padding.Left
			parent.Data.Size.Height = sumSizes.Height + parent.Data.Padding.Bottom + parent.Data.Padding.Top + gapAmount
		} else {
			parent.Data.Size.Width = sumSizes.Width + parent.Data.Padding.Right + parent.Data.Padding.Left + gapAmount
			parent.Data.Size.Height += parent.Data.Padding.Bottom + parent.Data.Padding.Top
		}
	}

	if len(validChildren) > 0 {
		parent.Children = append(parent.Children, validChildren...)
	}

	return parent
}

func CalculateDynamicElements(parent *TreeNode) {
	remainingWidth := parent.Data.Size.Width
	remainingWidth -= parent.Data.Padding.Left + parent.Data.Padding.Right
	remainingHeight := parent.Data.Size.Height
	remainingHeight -= parent.Data.Padding.Top + parent.Data.Padding.Bottom

	growableChildren := make([]*TreeNode, 0, len(parent.Children))
	shrinkableChildren := make([]*TreeNode, 0, len(parent.Children))
	for _, child := range parent.Children {
		if child != nil {
			if parent.Data.Column {
				remainingHeight -= child.Data.Size.Height
			} else {
				remainingWidth -= child.Data.Size.Width
			}
			if child.Children != nil {
				CalculateDynamicElements(child)
			}
			if child.Data.Shrink.Height || child.Data.Shrink.Width {
				shrinkableChildren = append(shrinkableChildren, child)
			}
			if child.Data.Grow.Height || child.Data.Grow.Width {
				growableChildren = append(growableChildren, child)
			}
		}
	}
	var gapAmount int32
	if len(parent.Children) > 0 {
		gapAmount = int32((len(parent.Children) - 1) * int(parent.Data.ChildGap))
	}

	if parent.Data.Column {
		remainingHeight -= gapAmount
		if remainingHeight > 0 {
			calculateGrowHeight(growableChildren, remainingHeight)
		} else if remainingHeight < 0 {
			calculateShrinkHeight(shrinkableChildren, remainingHeight)
		}
		for _, child := range parent.Children {
			if child != nil {
				if remainingWidth > 0 {
					if child.Data.Grow.Width {
						child.Data.Size.Width = remainingWidth
					}
				} else if remainingWidth < 0 {
					if child.Data.Shrink.Width {
						child.Data.Size.Width = remainingWidth
					}
				}
			}
		}
	} else {
		remainingWidth -= gapAmount
		if remainingWidth > 0 {
			calculateGrowWidth(growableChildren, remainingWidth)
		} else if remainingWidth < 0 {
			calculateShrinkWidth(shrinkableChildren, remainingWidth)
		}
		for _, child := range parent.Children {
			if child != nil {
				if remainingHeight > 0 {
					if child.Data.Grow.Height {
						child.Data.Size.Height = remainingHeight
					}
				} else if remainingHeight < 0 {
					if child.Data.Shrink.Height {
						child.Data.Size.Height = remainingHeight
					}
				}
			}
		}
	}

}

func Node(data NodeData) *TreeNode {
	node := &TreeNode{Data: data}
	return node
}

func DrawUI(rootNode *TreeNode) {
	for _, child := range rootNode.Children {
		if child.Parent != nil {
			childPosition := rl.Vector2Add(rootNode.Data.Position, child.Data.Position)
			child.Data.Position = childPosition
			rect := rl.Rectangle{
				Width:  float32(child.Data.Size.Width),
				Height: float32(child.Data.Size.Height),
				X:      childPosition.X,
				Y:      childPosition.Y,
			}
			rl.DrawText(fmt.Sprint(child.Data.Size), int32(childPosition.X+1), int32(childPosition.Y-10), 10, rl.Black)
			rl.DrawRectangleRounded(rect, .05, 8, child.Data.BackgroundColor)
			DrawUI(child)
		}
	}
}

func calculateGrowWidth(growableChildren []*TreeNode, remainingWidth int32) {
	if len(growableChildren) > 1 {
		for remainingWidth > 0 {
			smallest := growableChildren[0]
			if remainingWidth == 1 {
				smallest.Data.Size.Width += remainingWidth
				remainingWidth = 0
			}
			var secondSmallest *TreeNode
			widthToAdd := remainingWidth
			for _, child := range growableChildren {
				if child.Data.Size.Width < smallest.Data.Size.Width {
					secondSmallest = smallest
					smallest = child
				}
				if child.Data.Size.Width > smallest.Data.Size.Width {
					secondSmallest.Data.Size.Width = min(secondSmallest.Data.Size.Width, child.Data.Size.Width)
					widthToAdd = secondSmallest.Data.Size.Width - smallest.Data.Size.Width
				}
			}

			widthToAdd = min(widthToAdd, remainingWidth/int32(len(growableChildren)))

			for _, child := range growableChildren {
				if child.Data.Size.Width == smallest.Data.Size.Width {
					child.Data.Size.Width += widthToAdd
					remainingWidth -= widthToAdd
				}
			}
		}
	} else if len(growableChildren) == 1 {
		child := growableChildren[0]
		if child.Data.Grow.Width {
			if child.Data.Size.Width < remainingWidth {
				child.Data.Size.Width = remainingWidth
			}
		}
	}
}

func calculateGrowHeight(growableChildren []*TreeNode, remainingHeight int32) {
	if len(growableChildren) > 1 {
		for remainingHeight > 0 {
			smallest := growableChildren[0]
			if remainingHeight == 1 {
				smallest.Data.Size.Height += remainingHeight
				remainingHeight = 0
			}
			var secondSmallest *TreeNode
			widthToAdd := remainingHeight
			for _, child := range growableChildren {
				if child.Data.Size.Height < smallest.Data.Size.Height {
					secondSmallest = smallest
					smallest = child
				}
				if child.Data.Size.Height > smallest.Data.Size.Height {
					secondSmallest.Data.Size.Height = min(secondSmallest.Data.Size.Height, child.Data.Size.Height)
					widthToAdd = secondSmallest.Data.Size.Height - smallest.Data.Size.Height
				}
			}

			widthToAdd = min(widthToAdd, remainingHeight/int32(len(growableChildren)))

			for _, child := range growableChildren {
				if child.Data.Size.Height == smallest.Data.Size.Height {
					child.Data.Size.Height += widthToAdd
					remainingHeight -= widthToAdd
				}
			}
		}
	} else if len(growableChildren) == 1 {
		child := growableChildren[0]
		if child.Data.Grow.Height {
			if child.Data.Size.Height < remainingHeight {
				child.Data.Size.Height = remainingHeight
			}
		}
	}
}

func calculateShrinkWidth(shrinkableChildren []*TreeNode, remainingWidth int32) {
	if len(shrinkableChildren) > 1 {
		for remainingWidth < 0 {
			biggest := shrinkableChildren[0]
			if remainingWidth == -1 {
				biggest.Data.Size.Width += remainingWidth
				remainingWidth = 0
			}
			var secondBiggest *TreeNode
			widthToSub := remainingWidth
			for _, child := range shrinkableChildren {
				if child.Data.Size.Width > biggest.Data.Size.Width {
					secondBiggest = biggest
					biggest = child
				}
				if child.Data.Size.Width < biggest.Data.Size.Width {
					secondBiggest.Data.Size.Width = max(secondBiggest.Data.Size.Width, child.Data.Size.Width)
					widthToSub = secondBiggest.Data.Size.Width + biggest.Data.Size.Width
				}
			}

			widthToSub = min(widthToSub, remainingWidth/int32(len(shrinkableChildren)))

			for _, child := range shrinkableChildren {
				if child.Data.Size.Width == biggest.Data.Size.Width {
					child.Data.Size.Width -= widthToSub
					remainingWidth += widthToSub
				}
			}
		}
	} else if len(shrinkableChildren) == 1 {
		child := shrinkableChildren[0]
		if child.Data.Shrink.Width {
			if child.Data.Size.Width > remainingWidth {
				child.Data.Size.Width = remainingWidth
			}
		}
	}
}

func calculateShrinkHeight(shrinkableChildren []*TreeNode, remainingHeight int32) {
	if len(shrinkableChildren) > 1 {
		for remainingHeight < 0 {
			biggest := shrinkableChildren[0]
			if remainingHeight == -1 {
				biggest.Data.Size.Height += remainingHeight
				remainingHeight = 0
			}
			var secondBiggest *TreeNode
			widthToSub := remainingHeight
			for _, child := range shrinkableChildren {
				if child.Data.Size.Height > biggest.Data.Size.Height {
					secondBiggest = biggest
					biggest = child
				}
				if child.Data.Size.Height < biggest.Data.Size.Height {
					secondBiggest.Data.Size.Height = max(secondBiggest.Data.Size.Height, child.Data.Size.Height)
					widthToSub = secondBiggest.Data.Size.Height + biggest.Data.Size.Height
				}
			}

			widthToSub = min(widthToSub, remainingHeight/int32(len(shrinkableChildren)))

			for _, child := range shrinkableChildren {
				if child.Data.Size.Height == biggest.Data.Size.Height {
					child.Data.Size.Height -= widthToSub
					remainingHeight += widthToSub
				}
			}
		}
	} else if len(shrinkableChildren) == 1 {
		child := shrinkableChildren[0]
		if child.Data.Shrink.Height {
			if child.Data.Size.Height > remainingHeight {
				child.Data.Size.Height = remainingHeight
			}
		}
	}
}
