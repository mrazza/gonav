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

// NavPlace represents a Place entry in the NavMesh
type NavPlace struct {
	ID    uint32     // ID of the place
	Name  string     // The name of the place
	Areas []*NavArea // Collection of areas in this place
}

// GetEstimatedCenter gets a rough estimate of the center of this NavPlace
func (np *NavPlace) GetEstimatedCenter() (Vector3, error) {
	accume := Vector3{0, 0, 0}
	weight := float32(0)

	for _, currArea := range np.Areas {
		area := currArea.GetRoughSquaredArea()
		center := currArea.GetCenter()
		center.Mul(area)
		accume.Add(center)
		weight += area
	}

	if weight == 0 {
		return Vector3{}, errors.New("Cannot estimate center for NavAreas with no area.")
	}

	accume.Div(weight)
	return accume, nil
}
