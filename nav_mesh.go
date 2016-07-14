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

import "sync"

// NavMesh represents an entire parsed Nav Mesh and provides functionality
// related to the manipulation and searching of the mesh
type NavMesh struct {
	Places         map[uint32]*NavPlace  // Places contained in this NavMesh
	Areas          map[uint32]*NavArea   // Areas contained in this NavMesh
	Ladders        map[uint32]*NavLadder // Ladders contained in this NavMesh
	QuadTreeAreas  *quadTreeNode         // QuadTree used for quickly searching the NavAreas by position
	MajorVersion   uint32                // The major version number of the nav file
	MinorVersion   uint32                // The minor version number of the nav file
	BSPSize        uint32                // The size of the BSP file the nav was generated from
	IsMeshAnalyzed bool                  // Tracks whether or not this NavMesh has been analyzed
}

func (mesh *NavMesh) connectGraph() {
	var wg sync.WaitGroup

	for _, area := range mesh.Areas {
		wg.Add(1)

		go func(currArea *NavArea) {
			defer wg.Done()

			currArea.connectGraph(mesh)
		}(area)
	}

	wg.Wait()
}
