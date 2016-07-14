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

import "errors"

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
func (node *quadTreeNode) FindAreaByPoint(point Vector3) *NavArea {
	if !node.containsPoint(point) {
		return nil
	}

	for _, currArea := range node.Areas {
		if currArea.ContainsPoint(point) {
			return currArea
		}
	}

	if node.isSubDivided() {
		if currArea := node.NorthWest.FindAreaByPoint(point); currArea != nil {
			return currArea
		} else if currArea := node.NorthEast.FindAreaByPoint(point); currArea != nil {
			return currArea
		} else if currArea := node.SouthWest.FindAreaByPoint(point); currArea != nil {
			return currArea
		} else if currArea := node.SouthEast.FindAreaByPoint(point); currArea != nil {
			return currArea
		}
	}

	return nil
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
