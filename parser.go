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
	"encoding/binary"
	"errors"
	"fmt"
	"io"
)

// Parser provides support for parsing .nav files.
type Parser struct {
	Reader io.Reader
}

type parserError struct {
	Message string
	Error   error
}

// Parse parses the nav mesh reader supplied to this instance.
func (p *Parser) Parse() (mesh NavMesh, err error) {
	if p.Reader == nil {
		return buildParseError("This parse instance does not have a Reader.")
	}

	// Handle the recover scenerio for panics from our reader functions
	defer func() {
		if r := recover(); r != nil {
			msg, ok := r.(parserError)

			if ok {
				mesh, err = buildParseError(fmt.Sprintf("Unexpected error when parsing: %s.\nParent error: %v.", msg.Message, msg.Error))
			} else {
				panic(r)
			}
		}
	}()

	// Let's check the magic number
	var magicNumber uint32
	p.read(&magicNumber)

	if magicNumber != 0xFEEDFACE {
		return buildParseError(fmt.Sprintf("Magic number is incorrect (%v vs %v). This is not a .nav file.", magicNumber, 0xFEEDFACE))
	}

	// Magic number passed! Time to start building the NavMesh
	mesh = NavMesh{}
	p.read(&mesh.MajorVersion)

	// Check the version
	if mesh.MajorVersion < 6 || mesh.MajorVersion > 16 {
		return buildParseError("Major version for this nav mesh is invalid.")
	}

	if mesh.MajorVersion >= 10 {
		p.read(&mesh.MinorVersion)
	}

	p.read(&mesh.BSPSize)

	if mesh.MajorVersion >= 14 {
		p.read(&mesh.IsMeshAnalyzed)
	}

	// Let's get the "places"
	mesh.Places = make(map[uint32]*NavPlace)
	var placeCount uint16
	p.read(&placeCount)

	for i := uint16(0); i < placeCount; i++ {
		id := uint32(i + 1)
		var nameLength uint16
		p.read(&nameLength)
		mesh.Places[id] = &NavPlace{ID: id, Name: p.readString(nameLength)[:nameLength-1]}
	}

	// Time to build the area objects
	mesh.Areas = make(map[uint32]*NavArea)
	mesh.QuadTreeAreas = &quadTreeNode{NorthWestPoint: Vector3{-16384, -16384, 0}, SouthEastPoint: Vector3{16384, 16384, 0}}

	var hasUnnamedAreas bool
	if mesh.MajorVersion > 11 {
		p.read(&hasUnnamedAreas)
	}

	var areaCount uint32
	p.read(&areaCount)

	for i := uint32(0); i < areaCount; i++ {
		var currArea NavArea
		p.read(&currArea.ID)

		if mesh.MajorVersion <= 8 {
			var flags byte
			p.read(&flags)
			currArea.Flags = uint32(flags)
		} else if mesh.MajorVersion < 13 {
			var flags uint16
			p.read(&flags)
			currArea.Flags = uint32(flags)
		} else {
			p.read(&currArea.Flags)
		}

		p.read(&currArea.NorthWest)
		p.read(&currArea.SouthEast)
		p.read(&currArea.NorthEastZ)
		p.read(&currArea.SouthWestZ)

		// Time to handle the connections
		for direction := NavDirection(0); direction < NavDirectionMax; direction++ {
			var connectionCount uint32
			p.read(&connectionCount)

			for connectionIndex := uint32(0); connectionIndex < connectionCount; connectionIndex++ {
				var currConnection NavConnection
				currConnection.SourceArea = &currArea
				currConnection.Direction = direction
				p.read(&currConnection.TargetAreaID)

				currArea.Connections = append(currArea.Connections, &currConnection)
			}
		}

		// Time to handle the spots
		var hidingSpotCount byte
		p.read(&hidingSpotCount)

		for hidingIndex := byte(0); hidingIndex < hidingSpotCount; hidingIndex++ {
			var currSpot NavHidingSpot
			p.read(&currSpot.ID)
			p.read(&currSpot.Location)
			p.read(&currSpot.Flags)

			currArea.HidingSpots = append(currArea.HidingSpots, &currSpot)
		}

		// Throw away garbage if this is old
		if mesh.MajorVersion < 15 {
			var approachAreaCount byte
			p.read(&approachAreaCount)

			// Advance the pointer to skip past this garbage.
			// What is actually happening here is the following struct laid end to end
			// apprachAreaCount times.
			// struct
			// {
			//      uint approachHereId;
			//      uint approachPrevId;
			//      byte approachType;
			//      uint approachNextId;
			//      byte approachHow;
			// }
			p.advanceBytes((4*3 + 2) * int(approachAreaCount))
		}

		// Handle encounter paths
		var encounterPathCount uint32
		p.read(&encounterPathCount)

		for pathIndex := uint32(0); pathIndex < encounterPathCount; pathIndex++ {
			var currPath NavEncounterPath
			p.read(&currPath.FromAreaID)
			p.read(&currPath.FromDirection)
			p.read(&currPath.ToAreaID)
			p.read(&currPath.ToDirection)

			var spotCount byte
			p.read(&spotCount)

			for spotIndex := byte(0); spotIndex < spotCount; spotIndex++ {
				var currSpot NavEncounterSpot
				p.read(&currSpot.OrderID)
				var distance byte
				p.read(&distance)
				currSpot.ParametricDistiance = float32(distance) / 255

				currPath.Spots = append(currPath.Spots, &currSpot)
			}

			currArea.EncounterPaths = append(currArea.EncounterPaths, &currPath)
		}

		// Handle places
		var placeID uint16
		p.read(&placeID)
		place, ok := mesh.Places[uint32(placeID)]

		if ok {
			currArea.Place = place
			place.Areas = append(place.Areas, &currArea)
		}

		// Handle ladders
		for currDirection := NavLadderDirection(0); currDirection < NavLadderDirectionMax; currDirection++ {
			var ladderConnectionCount uint32
			p.read(&ladderConnectionCount)

			for connectionIndex := uint32(0); connectionIndex < ladderConnectionCount; connectionIndex++ {
				var currConnection NavLadderConnection
				currConnection.SourceArea = &currArea
				currConnection.Direction = currDirection
				p.read(&currConnection.TargetID)

				currArea.LadderConnections = append(currArea.LadderConnections, &currConnection)
			}
		}

		// Occupy times
		p.read(&currArea.EarliestOccupyTimeFirstTeam)
		p.read(&currArea.EarliestOccupyTimeSecondTeam)

		// Light intensity
		if mesh.MajorVersion >= 11 {
			p.read(&currArea.NorthWestLightIntensity)
			p.read(&currArea.NorthEastLightIntensity)
			p.read(&currArea.SouthEastLightIntensity)
			p.read(&currArea.SouthWestLightIntensity)
		}

		// Visible areas
		if mesh.MajorVersion >= 16 {
			var visibleAreaCount uint32
			p.read(&visibleAreaCount)

			for visibleIndex := uint32(0); visibleIndex < visibleAreaCount; visibleIndex++ {
				var currVisible NavVisibleArea
				p.read(&currVisible.VisibleAreaID)
				p.read(&currVisible.Attributes)

				currArea.VisibleAreas = append(currArea.VisibleAreas, &currVisible)
			}
		}

		p.read(&currArea.InheritVisibilityFromAreaID)

		// Skip passed some unknown data
		var garbageCount byte
		p.read(&garbageCount)
		p.advanceBytes(int(garbageCount) * 14)

		mesh.Areas[currArea.ID] = &currArea
		mesh.QuadTreeAreas.InsertArea(&currArea)
	}

	// Time to build the ladder objects
	mesh.Ladders = make(map[uint32]*NavLadder)
	var ladderCount uint32
	p.read(&ladderCount)

	for ladderIndex := uint32(0); ladderIndex < ladderCount; ladderIndex++ {
		var currLadder NavLadder

		p.read(&currLadder.ID)
		p.read(&currLadder.Width)
		p.read(&currLadder.Top)
		p.read(&currLadder.Bottom)
		p.read(&currLadder.Length)
		p.read(&currLadder.Direction)

		p.read(&currLadder.TopForwardAreaID)
		currLadder.TopForwardArea = mesh.Areas[currLadder.TopForwardAreaID]

		p.read(&currLadder.TopLeftAreaID)
		currLadder.TopLeftArea = mesh.Areas[currLadder.TopLeftAreaID]

		p.read(&currLadder.TopRightAreaID)
		currLadder.TopRightArea = mesh.Areas[currLadder.TopRightAreaID]

		p.read(&currLadder.TopBehindAreaID)
		currLadder.TopBehindArea = mesh.Areas[currLadder.TopBehindAreaID]

		p.read(&currLadder.BottomAreaID)
		currLadder.BottomArea = mesh.Areas[currLadder.BottomAreaID]

		mesh.Ladders[currLadder.ID] = &currLadder
	}

	// Ok we're done parsing the file, now it's time to connect the graph
	mesh.connectGraph()

	return mesh, nil
}

func (p *Parser) read(data interface{}) {
	switch t := data.(type) {
	case *Vector3:
		var x, y, z float32
		p.read(&x)
		p.read(&y)
		p.read(&z)
		*t = Vector3{x, y, z}

	case *bool:
		var b byte
		p.read(&b)
		*t = b > 0

	case *NavDirection:
		var b byte
		p.read(&b)
		*t = NavDirection(b)

	case *NavLadderDirection:
		var i uint32
		p.read(&i)
		*t = NavLadderDirection(i)

	default:
		err := binary.Read(p.Reader, binary.LittleEndian, t)

		if err != nil {
			panic(parserError{"Failed to read data", err})
		}
	}
}

func (p *Parser) readString(length uint16) string {
	data := make([]byte, length)
	n, ok := p.Reader.Read(data)

	if ok != nil {
		panic(parserError{"Failed to read string data", ok})
	} else if n != int(length) {
		panic(parserError{"Failed to read string data, length incorrect", nil})
	}

	return string(data)
}

func (p *Parser) advanceBytes(length int) {
	data := make([]byte, length)
	n, ok := p.Reader.Read(data)

	if ok != nil {
		panic(parserError{"Failed to advance data", ok})
	} else if n != int(length) {
		panic(parserError{"Failed to advance data, length incorrect", nil})
	}
}

func buildParseError(err string) (NavMesh, error) {
	return NavMesh{}, errors.New(err)
}
