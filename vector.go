// Copyright (C) 2016 Matt Razza
// Use of this source code is governed by
// the license found in the LICENSE file.

// Package gonav provides functionality related to CS:GO Nav Meshes
package gonav

import "math"

// Vector3 represents a 3D vector or point in space with X, Y, Z components.
type Vector3 struct {
	X, Y, Z float32 // The X, Y, and Z coordinates of the vector
}

// LengthSquared gets the square of the length of the Vector.
// This operation is faster than Length().
func (v *Vector3) LengthSquared() float32 {
	return v.X*v.X + v.Y*v.Y + v.Z*v.Z
}

// Length gets the length of the Vector.
func (v *Vector3) Length() float32 {
	return float32(math.Sqrt(float64(v.LengthSquared())))
}

// Normalize normalizes the Vector by setting its length to 1.
func (v *Vector3) Normalize() {
	length := v.Length()
	v.X /= length
	v.Y /= length
	v.Z /= length
}

// Add adds the specified Vector to this one.
func (v *Vector3) Add(left Vector3) {
	v.X += left.X
	v.Y += left.Y
	v.Z += left.Z
}

// Sub subtracts the specified Vector from this one.
func (v *Vector3) Sub(left Vector3) {
	v.X -= left.X
	v.Y -= left.Y
	v.Z -= left.Z
}

// Mul multiplies the specified scalar to this vector.
func (v *Vector3) Mul(left float32) {
	v.X *= left
	v.Y *= left
	v.Z *= left
}

// Div divides this vector by the specified scalar.
func (v *Vector3) Div(left float32) {
	v.X /= left
	v.Y /= left
	v.Z /= left
}
