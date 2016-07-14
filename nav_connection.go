// Copyright (C) 2016 Matt Razza
// Use of this source code is governed by
// the license found in the LICENSE file.

// Package gonav provides functionality related to CS:GO Nav Meshes
package gonav

// NavDirection represents a cardinal direction
type NavDirection int

const (
	// NavDirectionNorth is the north cardinal direction
	NavDirectionNorth NavDirection = iota

	// NavDirectionEast is the east cardinal direction
	NavDirectionEast

	// NavDirectionSouth is the south cardinal direction
	NavDirectionSouth

	// NavDirectionWest is the west cardinal direction
	NavDirectionWest

	// NavDirectionMax is the max value for nav directions
	NavDirectionMax
)

// NavConnection represents a connection between two NavAreas
type NavConnection struct {
	SourceArea   *NavArea     // The starting area for this connection
	TargetAreaID uint32       // The ID of the target area for this NavConnection
	TargetArea   *NavArea     // The target area for this connection
	Direction    NavDirection // The direction of the connection between these two areas
}

func (conn *NavConnection) connectGraph(mesh *NavMesh) {
	conn.TargetArea = mesh.Areas[conn.TargetAreaID]
}
