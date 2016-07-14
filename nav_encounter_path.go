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
