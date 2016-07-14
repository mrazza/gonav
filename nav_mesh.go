// Copyright (C) 2016 Matt Razza
// Use of this source code is governed by
// the license found in the LICENSE file.

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
