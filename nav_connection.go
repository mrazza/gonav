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
