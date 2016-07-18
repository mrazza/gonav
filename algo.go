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
	"container/heap"
	"errors"
)

// MeshConnectionCalculator is a func that calculates the cost of a connection in a nav mesh
type MeshConnectionCalculator func(*NavConnection) float32

// MeshLadderCalculator is a func that calculates the cost of a ladder between two nav areas
// The first NavArea is the starting area, the second NavArea is the ending area
type MeshLadderCalculator func(*NavLadder, NavLadderDirection, *NavArea, *NavArea) float32

// HeuristicCalculator calculates an estimated heurisitic of cost between two (likely disjoint) NavAreas
// This is used by the A* algorithm so must be admissible AND monotonic
type HeuristicCalculator func(*NavArea, *NavArea) float32

// PathNode is a single node along a path
type PathNode struct {
	Area               *NavArea
	PrevNode           *PathNode
	CostFromStart      float32
	estimatedCostToEnd float32
}

// Path represents a path between two points
type Path struct {
	Nodes []*PathNode
}

// GetCost gets the total cost of the path
func (p *Path) GetCost() float32 {
	return p.Nodes[len(p.Nodes)-1].CostFromStart
}

func newPath(endNode *PathNode) Path {
	var retPath Path
	currNode := endNode

	for currNode != nil {
		retPath.Nodes = append([]*PathNode{currNode}, retPath.Nodes...)
		currNode = currNode.PrevNode
	}

	return retPath
}

// SimpleBuildShortestPath builds a path (via PathFinding A*) and returns a Path object containing the start and end nodes of the path
// startArea and endArea are the starting and ending NavAreas for the path
// This function uses simple default cost and heuristic functions when building the Path. For more control, see BuildShortestPath
func SimpleBuildShortestPath(startArea, endArea *NavArea) (Path, error) {
	return BuildShortestPath(
		startArea,
		endArea,
		func(con *NavConnection) float32 {
			distance := con.SourceArea.GetCenter()
			distance.Sub(con.TargetArea.GetCenter())
			return distance.Length()
		},
		func(ladder *NavLadder, direction NavLadderDirection, start *NavArea, end *NavArea) float32 {
			locationVector := start.GetCenter()
			locationVector.Sub(end.GetCenter())
			cost := (&Vector3{locationVector.X, locationVector.Y, 0}).Length() // Let's ignore the Z because that's what ladder.Length does
			return ladder.Length + cost
		},
		func(start *NavArea, end *NavArea) float32 {
			distance := start.GetCenter()
			distance.Sub(end.GetCenter())
			return distance.Length()
		})
}

// BuildShortestPath builds a path (via PathFinding A*) and returns a Path object containing the start and end nodes of the path
// startArea and endArea are the starting and ending NavAreas for the path
// areaCostCalc is a func() that calculates the "cost" of a connection between two NavAreas
// ladderCostCalc is a func() that calculates the "cost" of a connection via a ladder
// heurisiticCost is a func() that estimates an admissible AND monotonic cost for two (likely nonadjacent) NavAreas
func BuildShortestPath(startArea, endArea *NavArea, areaCostCalc MeshConnectionCalculator, ladderCostCalc MeshLadderCalculator, heurisiticCost HeuristicCalculator) (Path, error) {
	closedSet := make(map[*NavArea]bool)
	nodeLookup := make(map[*NavArea]*queueItem)
	openSet := make(priorityQueue, 0)
	heap.Init(&openSet)

	start := PathNode{
		Area:               startArea,
		CostFromStart:      0,
		estimatedCostToEnd: heurisiticCost(startArea, endArea)}

	nodeLookup[startArea] = openSet.CreateAndPush(&start)

	for openSet.Len() > 0 {
		currentNode := openSet.PopCast()

		if currentNode.Area == endArea {
			return newPath(currentNode), nil // We found the end!
		}

		// Add this to where we've been
		closedSet[currentNode.Area] = true

		// Look at all the places connected to where we are
		for _, currConnection := range currentNode.Area.Connections {
			if closedSet[currConnection.TargetArea] {
				continue // We've been here before
			}

			// Calculate the cost to get there from here
			newCost := currentNode.CostFromStart + areaCostCalc(currConnection)
			item := nodeLookup[currConnection.TargetArea]
			var currNode *PathNode

			if item == nil {
				currNode = &PathNode{Area: currConnection.TargetArea}
			} else {
				currNode = item.pathNode

				if newCost >= currNode.CostFromStart {
					continue // Going there from here isn't any better than before
				}
			}

			// Either this is a new place to go, or we found a better way to get there. Update.
			currNode.PrevNode = currentNode
			currNode.CostFromStart = newCost
			currNode.estimatedCostToEnd = newCost + heurisiticCost(currNode.Area, endArea)

			if item != nil {
				openSet.update(item)
			} else {
				nodeLookup[currNode.Area] = openSet.CreateAndPush(currNode)
			}
		}

		// What about the places we're connected to via ladders
		for _, currLadderCon := range currentNode.Area.LadderConnections {
			currLadder := currLadderCon.TargetLadder
			var ladderAreas []*NavArea

			switch currLadder.Direction {
			case NavLadderDirectionUp:
				if currLadder.TopBehindArea != nil {
					ladderAreas = append(ladderAreas, currLadder.TopBehindArea)
				}

				if currLadder.TopForwardArea != nil {
					ladderAreas = append(ladderAreas, currLadder.TopForwardArea)
				}

				if currLadder.TopRightArea != nil {
					ladderAreas = append(ladderAreas, currLadder.TopRightArea)
				}

				if currLadder.TopLeftArea != nil {
					ladderAreas = append(ladderAreas, currLadder.TopLeftArea)
				}

			case NavLadderDirectionDown:
				if currLadder.BottomArea != nil {
					ladderAreas = append(ladderAreas, currLadder.BottomArea)
				}
			}

			for _, currArea := range ladderAreas {
				if closedSet[currArea] {
					continue // We've been here before
				}

				// Calculate the cost to get there from here
				newCost := currentNode.CostFromStart + ladderCostCalc(currLadder, currLadder.Direction, currentNode.Area, currArea)
				item := nodeLookup[currArea]
				var currNode *PathNode

				if item == nil {
					currNode = &PathNode{Area: currArea}
				} else {
					currNode = item.pathNode

					if newCost >= currNode.CostFromStart {
						continue // Going there from here isn't any better than before
					}
				}

				// Either this is a new place to go, or we found a better way to get there. Update.
				currNode.PrevNode = currentNode
				currNode.CostFromStart = newCost
				currNode.estimatedCostToEnd = newCost + heurisiticCost(currNode.Area, endArea)

				if item != nil {
					openSet.update(item)
				} else {
					nodeLookup[currNode.Area] = openSet.CreateAndPush(currNode)
				}
			}
		}
	}

	return Path{}, errors.New("Could not find a path. Areas are not connected.")
}

// Priority queue implementation from a modified version of the source on
// https://golang.org/pkg/container/heap/
type queueItem struct {
	pathNode *PathNode
	index    int
}

type priorityQueue []*queueItem

func (pq priorityQueue) Len() int {
	return len(pq)
}

func (pq priorityQueue) Less(i, j int) bool {
	return pq[i].pathNode.estimatedCostToEnd < pq[j].pathNode.estimatedCostToEnd
}

func (pq priorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *priorityQueue) CreateAndPush(pn *PathNode) *queueItem {
	q := &queueItem{pathNode: pn}
	heap.Push(pq, q)
	return q
}

func (pq *priorityQueue) Push(q interface{}) {
	item := q.(*queueItem)
	n := len(*pq)
	item.index = n
	*pq = append(*pq, item)
}

func (pq *priorityQueue) PopCast() *PathNode {
	return heap.Pop(pq).(*PathNode)
}

func (pq *priorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	item.index = -1
	*pq = old[0 : n-1]
	return item.pathNode
}

func (pq *priorityQueue) Remove(item *queueItem) {
	heap.Remove(pq, item.index)
}

func (pq *priorityQueue) update(item *queueItem) {
	heap.Fix(pq, item.index)
}
