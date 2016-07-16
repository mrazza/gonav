/*
	gonav - A Source Engine navigation mesh file parser written in Go.
	Copyright (C) 2016  Matt Razza

	This program is free software: you can redistribute it and/or modify
	it under the terms of the GNU Affero General Public License as published
	by the Free Software Foundation, either version 3 of the License, or
	(at your option) any later version.

	This program is distributed in the hope that it will be useful,
	but WITHOUT ANY WARRANTY; without even the implied warranty of
	MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
	GNU Affero General Public License for more details.

	You should have received a copy of the GNU Affero General Public License
	along with this program.  If not, see <http://www.gnu.org/licenses/>.
*/

// Package gonav provides functionality related to CS:GO Nav Meshes
package gonav

import (
	"errors"
	"math"
)

const preferredNodeCapacity int = 4

// quadTreeNode represents a single node in a QuadTree (possibly the root node, possibly not)
type quadTreeNode struct {
	Areas          []*NavArea    // The NavAreas contained in this node
	NorthWestPoint Vector3       // The north west point for the bounds of this node
	SouthEastPoint Vector3       // The south east point for the bounds of this node
	NorthWest      *quadTreeNode // The north west sub node (nil if not sub-divided)
	NorthEast      *quadTreeNode // The north east sub node (nil if not sub-divided)
	SouthWest      *quadTreeNode // The south west sub node (nil if not sub-divided)
	SouthEast      *quadTreeNode // The south east sub node (nil if not sub-divided)
}

// InsertArea inserts the specified NavArea into this quad tree
func (node *quadTreeNode) InsertArea(area *NavArea) (*quadTreeNode, error) {
	if !node.isAreaFullyContained(area) {
		return nil, errors.New("Specified area cannot be added because it is not fully contained within this node.")
	}

	if node.isSubDivided() {
		if addedNode, ok := node.NorthWest.InsertArea(area); ok == nil {
			return addedNode, nil
		} else if addedNode, ok := node.NorthEast.InsertArea(area); ok == nil {
			return addedNode, nil
		} else if addedNode, ok := node.SouthWest.InsertArea(area); ok == nil {
			return addedNode, nil
		} else if addedNode, ok := node.SouthEast.InsertArea(area); ok == nil {
			return addedNode, nil
		} else {
			node.Areas = append(node.Areas, area)
			return node, nil
		}
	} else {
		if len(node.Areas) >= preferredNodeCapacity {
			node.subDivide()
			return node.InsertArea(area)
		}

		node.Areas = append(node.Areas, area)
		return node, nil
	}
}

// Finds the area that contains the specified point; nil if the area could not be found
// The Z-value is used to find the closest area that contains the X and Y values
// If allowBelow is true the area closest by Z that contains this point is returned
// If allowBelow is false the area closest by Z that is BELOW this point is returned
// Think of allowBelow as "allow the specified point to be below a nav area"
func (node *quadTreeNode) FindAreaByPoint(point Vector3, allowBelow bool) *NavArea {
	if !node.containsPoint(point) {
		return nil
	}

	// We'll use these to keep track of the current known best area
	var bestArea *NavArea
	bestDistance := float32(math.MaxFloat32)

	// This lambda will update our currently known best area
	updateBestArea := func(currArea *NavArea) {
		currDistance := float32(math.Abs(float64(currArea.DistanceFromZ(point))))

		if currDistance < bestDistance {
			bestArea = currArea
			bestDistance = currDistance
		}
	}

	// Let's loop through everything in this node
	for _, currArea := range node.Areas {
		if currArea.ContainsPoint(point, allowBelow) {
			updateBestArea(currArea)
		}
	}

	// If we're subdivided we need to recurse
	if node.isSubDivided() {
		if currArea := node.NorthWest.FindAreaByPoint(point, allowBelow); currArea != nil {
			updateBestArea(currArea)
		} else if currArea := node.NorthEast.FindAreaByPoint(point, allowBelow); currArea != nil {
			updateBestArea(currArea)
		} else if currArea := node.SouthWest.FindAreaByPoint(point, allowBelow); currArea != nil {
			updateBestArea(currArea)
		} else if currArea := node.SouthEast.FindAreaByPoint(point, allowBelow); currArea != nil {
			updateBestArea(currArea)
		}
	}

	return bestArea
}

// containsPoint determines whether or not the specified point is contained in the node
func (node *quadTreeNode) containsPoint(point Vector3) bool {
	return node.NorthWestPoint.X <= point.X && node.NorthWestPoint.Y <= point.Y && node.SouthEastPoint.X >= point.X && node.SouthEastPoint.Y >= point.Y
}

// isAreaFullyContained determines whether or not the specified NavArea can be fully contained within this node
func (node *quadTreeNode) isAreaFullyContained(area *NavArea) bool {
	return node.NorthWestPoint.X <= area.NorthWest.X && node.NorthWestPoint.Y <= area.NorthWest.Y && node.SouthEastPoint.X >= area.SouthEast.X && node.SouthEastPoint.Y >= area.SouthEast.Y
}

// Sub divides this node by adding four children nodes below it and repositioning its areas
func (node *quadTreeNode) subDivide() error {
	if node.isSubDivided() {
		return errors.New("Cannot subdivide already subdivided node.")
	}

	node.NorthWest = new(quadTreeNode)
	node.NorthEast = new(quadTreeNode)
	node.SouthWest = new(quadTreeNode)
	node.SouthEast = new(quadTreeNode)

	currAreas := node.Areas
	node.Areas = nil

	for _, currArea := range currAreas {
		node.InsertArea(currArea)
	}

	return nil
}

// isSubDivided determines whether or not this node is sub divided
func (node *quadTreeNode) isSubDivided() bool {
	return node.NorthWest != nil || node.NorthEast != nil || node.SouthWest != nil || node.SouthEast != nil
}
