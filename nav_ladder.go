// Copyright (C) 2016 Matt Razza
// Use of this source code is governed by
// the license found in the LICENSE file.

// Package gonav provides functionality related to CS:GO Nav Meshes
package gonav

// NavLadderDirection represents the direction between a ladder and its connection
type NavLadderDirection int

const (
	// NavLadderDirectionUp means up the ladder
	NavLadderDirectionUp NavLadderDirection = iota

	// NavLadderDirectionDown means down the ladder
	NavLadderDirectionDown

	// NavLadderDirectionMax the max value for NavLadderDirection's
	NavLadderDirectionMax
)

// NavLadderConnection represents a connection between an area and a ladder
type NavLadderConnection struct {
	SourceArea   *NavArea           // The area that is the source of this NavConnection
	TargetID     uint32             // The ID of the ladder target
	TargetLadder *NavLadder         // The ladder
	Direction    NavLadderDirection // The direction of this connection
}

// NavLadder represents a ladder within the world
type NavLadder struct {
	ID               uint32             // The ID of the ladder
	Width            float32            // The width of the ladder
	Length           float32            // The length of the ladder
	Top              Vector3            // The location of the center of the top of the ladder
	Bottom           Vector3            // The location of the center of the bottom of the ladder
	Direction        NavLadderDirection // The direction of the NavLadder
	TopForwardAreaID uint32             // ID of the area connected to the top-forward position of the ladder
	TopForwardArea   *NavArea           // The area connected to the top-forward position of the ladder
	TopLeftAreaID    uint32             // ID of the area connected to the top-left position of the ladder
	TopLeftArea      *NavArea           // The area connected to the top-forward position of the ladder
	TopRightAreaID   uint32             // ID of the area connected to the top-right position of the ladder
	TopRightArea     *NavArea           // The area connected to the top-right position of the ladder
	TopBehindAreaID  uint32             // ID of area connected to the top-behind position of the ladder
	TopBehindArea    *NavArea           // The area connected to the top-behind position of the ladder
	BottomAreaID     uint32             // ID of the area connected to the bottom of the ladder
	BottomArea       *NavArea           // The area connected to the bottom of the ladder
}

func (conn *NavLadderConnection) connectGraph(mesh *NavMesh) {
	conn.TargetLadder = mesh.Ladders[conn.TargetID]
}
