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
	"fmt"
)

// NavArea represents a NavArea as part of a NavMesh
type NavArea struct {
	ID                           uint32                 // ID of the NavArea
	NorthWest                    Vector3                // Location of the north-west point of this NavArea
	SouthEast                    Vector3                // Location of the south-east point of this NavArea
	Flags                        uint32                 // Bitflags set on this area
	NorthEastZ                   float32                // The Z-coord for the north-east point
	SouthWestZ                   float32                // The Z-coord for the south-west point
	NorthWestLightIntensity      float32                // The light intensity of the north-west corner
	NorthEastLightIntensity      float32                // The light intensity of the north-east corner
	SouthWestLightIntensity      float32                // The light intensity of the south-west corner
	SouthEastLightIntensity      float32                // The light intensity of the south-east corner
	Place                        *NavPlace              // The place this area is in
	Connections                  []*NavConnection       // The connections between this area and other areas
	HidingSpots                  []*NavHidingSpot       // The hiding spots in this NavArea
	EncounterPaths               []*NavEncounterPath    // The encounter paths for this area
	LadderConnections            []*NavLadderConnection // Connections between this area and ladders
	VisibleAreas                 []*NavVisibleArea      // Visible areas
	EarliestOccupyTimeFirstTeam  float32                // The earliest time the first team can occupy this area
	EarliestOccupyTimeSecondTeam float32                // The earliest time the second team can occupy this area
	InheritVisibilityFromAreaID  uint32                 // ID of the area to inherit our visibility from
}

// NavHidingSpot represents an identified hiding spot within a NavArea
type NavHidingSpot struct {
	ID       uint32  // ID of the hiding spot
	Location Vector3 // Location of the hiding NavHidingSpot
	Flags    byte    // Bitflags associated with this hiding spot
}

// NavVisibleArea represents a visible area
type NavVisibleArea struct {
	VisibleAreaID uint32   // ID of the visible area
	VisibleArea   *NavArea // The visible area
	Attributes    byte     // Bit-wise attributes
}

func (visArea *NavVisibleArea) connectGraph(mesh *NavMesh) {
	visArea.VisibleArea = mesh.Areas[visArea.VisibleAreaID]
}

func (area *NavArea) connectGraph(mesh *NavMesh) {
	for _, currConnection := range area.Connections {
		currConnection.connectGraph(mesh)
	}

	for _, currPath := range area.EncounterPaths {
		currPath.connectGraph(mesh)
	}

	for _, currLadder := range area.LadderConnections {
		currLadder.connectGraph(mesh)
	}

	for _, currArea := range area.VisibleAreas {
		currArea.connectGraph(mesh)
	}
}

// GetNorthEastPoint builds the north east point from the two known corner points and the known Z value
func (area *NavArea) GetNorthEastPoint() Vector3 {
	return Vector3{X: area.SouthEast.X, Y: area.NorthWest.Y, Z: area.NorthEastZ}
}

// GetSouthWestPoint builds the south west point from the two known corner points and the known Z value
func (area *NavArea) GetSouthWestPoint() Vector3 {
	return Vector3{X: area.NorthWest.X, Y: area.SouthEast.Y, Z: area.SouthWestZ}
}

// GetCenter gets the center point of this area.
func (area *NavArea) GetCenter() Vector3 {
	x := (area.NorthWest.X + area.SouthEast.X) / 2.0
	y := (area.NorthWest.Y + area.GetSouthWestPoint().Y) / 2.0
	z, err := area.GetZ(x, y)

	if err != nil {
		panic(err)
	}

	return Vector3{
		X: x,
		Y: y,
		Z: z}
}

// GetZ gets the Z-coord for the specified point within this area.
// An error is returned if the requested point is not within this area.
func (area *NavArea) GetZ(x, y float32) (float32, error) {
	if !area.ContainsPoint(Vector3{x, y, 0}) {
		return 0, errors.New("Cannot get Z. Specified point does not exist within the area.")
	}

	// Find the Z on the north and south lines that share the X-coord with our point
	width := area.SouthEast.X - area.NorthWest.X
	height := area.SouthEast.Y - area.NorthWest.Y
	distanceFromEast := area.SouthEast.X - x
	distanceFromSouth := area.SouthEast.Y - y
	northSlope := (area.NorthWest.Z - area.NorthEastZ) / width
	southSlope := (area.SouthWestZ - area.SouthEast.Z) / width
	northZ := area.NorthEastZ + northSlope*distanceFromEast
	southZ := area.SouthEast.Z + southSlope*distanceFromEast

	// Draw a line between those points on the north and south line
	// and find the Z value at the specified Y location
	finalLineSlope := (northZ - southZ) / height

	return southZ + finalLineSlope*distanceFromSouth, nil
}

// ContainsPoint determines whether or not the specified point is contained within this area
func (area *NavArea) ContainsPoint(point Vector3) bool {
	return area.NorthWest.X <= point.X &&
		area.NorthWest.Y <= point.Y &&
		area.SouthEast.X >= point.X &&
		area.SouthEast.Y >= point.Y
}

// GetRoughSquaredArea gets a rough estimate of the squared area of the NavArea
func (area *NavArea) GetRoughSquaredArea() float32 {
	return (area.SouthEast.X - area.NorthWest.X) * (area.SouthEast.Y - area.NorthWest.Y)
}

// GetClosestPointInArea gets the point closest to the specified point that is contained within this Area.
func (area *NavArea) GetClosestPointInArea(point Vector3) Vector3 {
	if area.ContainsPoint(point) {
		z, _ := area.GetZ(point.X, point.Y)
		return Vector3{point.X, point.Y, z}
	}

	var x, y float32

	// Let's do X
	if point.X < area.NorthWest.X {
		x = area.NorthWest.X
	} else if point.X > area.SouthEast.X {
		x = area.SouthEast.X
	} else {
		x = point.X
	}

	// Let's do Y
	if point.Y < area.NorthWest.Y {
		y = area.NorthWest.Y
	} else if point.Y > area.SouthEast.Y {
		y = area.SouthEast.Y
	} else {
		y = point.Y
	}

	z, _ := area.GetZ(x, y)

	return Vector3{X: x, Y: y, Z: z}
}

// String converts a NavArea into a human readable string
func (area *NavArea) String() string {
	name := "<null>"
	if area.Place != nil {
		name = area.Place.Name
	}

	return fmt.Sprintf("AreaID: %v [%v] @ {%v, %v}", area.ID, name, area.NorthWest, area.SouthEast)
}
