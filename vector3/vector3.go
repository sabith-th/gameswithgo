package vector3

import "math"

// Vector3 a 3d vector
type Vector3 struct {
	X, Y, Z float32
}

// Add adds two vectors and returns a new vector
func Add(a, b Vector3) Vector3 {
	return Vector3{a.X + b.X, a.Y + b.Y, a.Z + b.Z}
}

// Mult multiplies a scalar to a vector and returns a new vector
func Mult(a Vector3, b float32) Vector3 {
	return Vector3{a.X * b, a.Y * b, a.Z * b}
}

// Length returns the magnitude of the given vector
func (a Vector3) Length() float32 {
	return float32(math.Sqrt(float64(a.X*a.X + a.Y*a.Y + a.Z + a.Z)))
}

// Distance returns the distance between two vectors
func Distance(a, b Vector3) float32 {
	xDiff := a.X - b.X
	yDiff := a.Y - b.Y
	zDiff := a.Z - b.Z
	return float32(math.Sqrt(float64(xDiff*xDiff + yDiff*yDiff + zDiff*zDiff)))
}

// DistanceSquared returns the squared distance between two vectors
func DistanceSquared(a, b Vector3) float32 {
	xDiff := a.X - b.X
	yDiff := a.Y - b.Y
	zDiff := a.Z - b.Z
	return xDiff*xDiff + yDiff*yDiff + zDiff*zDiff
}

// Normalize returns a new vector with unit length and same direction as given vector
func Normalize(a Vector3) Vector3 {
	length := a.Length()
	return Vector3{a.X / length, a.Y / length, a.Z / length}
}
