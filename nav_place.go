// Copyright (C) 2016 Matt Razza
// Use of this source code is governed by
// the license found in the LICENSE file.

// Package gonav provides functionality related to CS:GO Nav Meshes
package gonav

// NavPlace represents a Place entry in the NavMesh
type NavPlace struct {
	ID    uint32     // ID of the place
	Name  string     // The name of the place
	Areas []*NavArea // Collection of areas in this place
}
