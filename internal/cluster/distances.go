package cluster

import (
	"math"

	"go.matteson.dev/gfx"
)

// MaxDistanceFunc ...
type MaxDistanceFunc func(interface{}, interface{}) float64

// DistanceMeasureFunc ...
type DistanceMeasureFunc func(gfx.Point, gfx.Point) float64

// EuclideanDistance ...
func EuclideanDistance(p1, p2 gfx.Point) float64 {
	dx, dy := p1.X-p2.X, p1.Y-p2.Y
	return math.Sqrt(dx*dx + dy*dy)
}

// ManhattanDistance ...
func ManhattanDistance(p1 gfx.Point, p2 gfx.Point) float64 {
	return math.Abs(p1.X-p2.X) + math.Abs(p1.Y-p2.Y)
}
