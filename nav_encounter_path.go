// Copyright (C) 2016 Matt Razza
// Use of this source code is governed by
// the license found in the LICENSE file.

// Package gonav provides functionality related to CS:GO Nav Meshes
package gonav

// NavEncounterPath represents an encounter path
type NavEncounterPath struct {
	FromAreaID    uint32              // The ID of the area the path comes from
	FromArea      *NavArea            // The Area the path comes from
	FromDirection NavDirection        // The direction from the source
	ToAreaID      uint32              // The ID of the area the path ends in
	ToArea        *NavArea            // The area the path ends in
	ToDirection   NavDirection        // The direction from the destination
	Spots         []*NavEncounterSpot // The spots along this path
}

// NavEncounterSpot represents a spot along an encounter path
type NavEncounterSpot struct {
	OrderID             uint32  // The ID of the order of this spot
	ParametricDistiance float32 // The parametric distance
}

func (path *NavEncounterPath) connectGraph(mesh *NavMesh) {
	path.FromArea = mesh.Areas[path.FromAreaID]
	path.ToArea = mesh.Areas[path.ToAreaID]
}
