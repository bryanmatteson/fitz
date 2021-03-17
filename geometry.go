package fitz

import (
	"fmt"
	"math"

	"github.com/ahmetb/go-linq"
)

type Orientation int

// Text orientations
const (
	OtherOrientation Orientation = iota
	Horizontal
	Rotate180
	Rotate270
	Rotate90
)

type Lines []Line

type Line struct {
	Pt1, Pt2 Point
}

func (l Line) Length() float64 {
	dx := l.Pt1.X - l.Pt2.X
	dy := l.Pt1.Y - l.Pt2.Y
	return math.Sqrt(dx*dx + dy*dy)
}

type Rect struct {
	X, Y Range
}

func MakeRectWH(x, y, w, h float64) Rect {
	return MakeRectCorners(x, y, x+w, y+h)
}

func MakeRectCorners(x0, y0, x1, y1 float64) Rect {
	if x0 > x1 {
		x0, x1 = x1, x0
	}
	if y0 > y1 {
		y0, y1 = y1, y0
	}
	return Rect{Range{x0, x1}, Range{y0, y1}}
}

func EmptyRect() Rect {
	return Rect{EmptyRange(), EmptyRange()}
}

// IsValid reports whether the rectangle is valid.
// This requires the width to be empty iff the height is empty.
func (r Rect) IsValid() bool {
	return r.X.IsEmpty() == r.Y.IsEmpty()
}

// Size returns the width and height of this rectangle in (x,y)-space. Empty
// rectangles have a negative width and height.
func (r Rect) Size() Point {
	return Point{r.X.Length(), r.Y.Length()}
}

// ContainsPoint reports whether the rectangle contains the given point.
// Rectangles are closed regions, i.e. they contain their boundary.
func (r Rect) ContainsPoint(p Point) bool {
	return r.X.Contains(p.X) && r.Y.Contains(p.Y)
}

// InteriorContainsPoint returns true iff the given point is contained in the interior
// of the region (i.e. the region excluding its boundary).
func (r Rect) InteriorContainsPoint(p Point) bool {
	return r.X.InteriorContains(p.X) && r.Y.InteriorContains(p.Y)
}

// Contains reports whether the rectangle contains the given rectangle.
func (r Rect) Contains(other Rect) bool {
	return r.X.ContainsRange(other.X) && r.Y.ContainsRange(other.Y)
}

// InteriorContains reports whether the interior of this rectangle contains all of the
// points of the given other rectangle (including its boundary).
func (r Rect) InteriorContains(other Rect) bool {
	return r.X.InteriorContainsRange(other.X) && r.Y.InteriorContainsRange(other.Y)
}

// Intersects reports whether this rectangle and the other rectangle have any points in common.
func (r Rect) Intersects(other Rect) bool {
	return r.X.Intersects(other.X) && r.Y.Intersects(other.Y)
}

// InteriorIntersects reports whether the interior of this rectangle intersects
// any point (including the boundary) of the given other rectangle.
func (r Rect) InteriorIntersects(other Rect) bool {
	return r.X.InteriorIntersects(other.X) && r.Y.InteriorIntersects(other.Y)
}

// Expanded returns a rectangle that has been expanded in the x-direction
// by margin.X, and in y-direction by margin.Y. If either margin is empty,
// then shrink the interval on the corresponding sides instead. The resulting
// rectangle may be empty. Any expansion of an empty rectangle remains empty.
func (r Rect) Expanded(margin Point) Rect {
	xx := r.X.Expanded(margin.X)
	yy := r.Y.Expanded(margin.Y)
	if xx.IsEmpty() || yy.IsEmpty() {
		return EmptyRect()
	}
	return Rect{xx, yy}
}

// ExpandedByMargin returns a Rect that has been expanded by the amount on all sides.
func (r Rect) ExpandedByMargin(margin float64) Rect {
	return r.Expanded(Point{margin, margin})
}

// Union returns the smallest rectangle containing the union of this rectangle and
// the given rectangle.
func (r Rect) Union(other Rect) Rect {
	return Rect{r.X.Union(other.X), r.Y.Union(other.Y)}
}

// Intersection returns the smallest rectangle containing the intersection of this
// rectangle and the given rectangle.
func (r Rect) Intersection(other Rect) Rect {
	xx := r.X.Intersection(other.X)
	yy := r.Y.Intersection(other.Y)
	if xx.IsEmpty() || yy.IsEmpty() {
		return EmptyRect()
	}

	return Rect{xx, yy}
}

func (r Rect) Width() float64 { return r.X.Length() }

func (r Rect) Height() float64 { return r.Y.Length() }

func (r Rect) IsEmpty() bool {
	return EqualEpsilon(r.X.Max, r.X.Min) || EqualEpsilon(r.Y.Max, r.Y.Min)
}

type Point struct {
	X, Y float64
}

// Add returns the sum of p and op.
func (p Point) Add(op Point) Point { return Point{p.X + op.X, p.Y + op.Y} }

// Sub returns the difference of p and op.
func (p Point) Sub(op Point) Point { return Point{p.X - op.X, p.Y - op.Y} }

// Mul returns the scalar product of p and m.
func (p Point) Mul(m float64) Point { return Point{m * p.X, m * p.Y} }

// Ortho returns a counterclockwise orthogonal point with the same norm.
func (p Point) Ortho() Point { return Point{-p.Y, p.X} }

// Dot returns the dot product between p and op.
func (p Point) Dot(op Point) float64 { return p.X*op.X + p.Y*op.Y }

// Cross returns the cross product of p and op.
func (p Point) Cross(op Point) float64 { return p.X*op.Y - p.Y*op.X }

// Norm returns the vector's norm.
func (p Point) Norm() float64 { return math.Hypot(p.X, p.Y) }

// Normalize returns a unit point in the same direction as p.
func (p Point) Normalize() Point {
	if p.X == 0 && p.Y == 0 {
		return p
	}
	return p.Mul(1 / p.Norm())
}

// PerpDot returns the perp dot product between OP and OQ, ie. zero if aligned and |OP|*|OQ| if perpendicular.
func (p Point) PerpDot(q Point) float64 { return p.X*q.Y - p.Y*q.X }

func (p Point) String() string { return fmt.Sprintf("%f, %f", p.X, p.Y) }

type Quad struct {
	BottomLeft  Point
	TopLeft     Point
	TopRight    Point
	BottomRight Point
}

func MakeQuad(left, bottom, right, top float64) Quad {
	return Quad{Point{left, bottom}, Point{left, top}, Point{right, top}, Point{right, bottom}}
}

func (q Quad) Left() float64 {
	if q.TopLeft.X < q.TopRight.X {
		return q.TopLeft.X
	}
	return q.TopRight.X
}

func (q Quad) Right() float64 {
	if q.BottomRight.X < q.BottomLeft.X {
		return q.BottomLeft.X
	}
	return q.BottomRight.X
}

func (q Quad) Bottom() float64 {
	if q.BottomRight.Y < q.TopRight.Y {
		return q.BottomRight.Y
	}
	return q.TopRight.Y
}

func (q Quad) Top() float64 {
	if q.TopLeft.Y < q.BottomLeft.Y {
		return q.BottomLeft.Y
	}
	return q.TopLeft.Y
}

func (q Quad) T() float64 {
	if q.BottomRight == q.BottomLeft {
		return math.Atan2(q.TopLeft.Y-q.BottomLeft.Y, q.TopLeft.X-q.BottomLeft.X) - math.Pi/2
	}
	return math.Atan2(q.BottomRight.Y-q.BottomLeft.Y, q.BottomRight.X-q.BottomLeft.X)
}

func (q Quad) Rotation() float64 {
	return q.T() * 180 / math.Pi
}

func (q Quad) Width() float64 { return q.Right() - q.Left() }

func (q Quad) Height() float64 { return q.Top() - q.Bottom() }

func (q Quad) Orientation() Orientation {
	if EqualEpsilon(q.BottomLeft.Y, q.BottomRight.Y) {
		if q.BottomLeft.X > q.BottomRight.X {
			return Rotate180
		}
		return Horizontal
	}
	if EqualEpsilon(q.BottomLeft.X, q.BottomRight.X) {
		if q.BottomLeft.Y > q.BottomRight.Y {
			return Rotate270
		}
		return Rotate90
	}
	return OtherOrientation
}

func (q Quad) Centroid() Point {
	cx := (q.BottomRight.X + q.TopRight.X + q.BottomLeft.X + q.TopLeft.X) / 4.0
	cy := (q.BottomRight.Y + q.TopRight.Y + q.BottomLeft.Y + q.TopLeft.Y) / 4.0
	return Point{cx, cy}
}

type Quads []Quad

func (q Quads) Normalize() (n Quad) {
	left, bottom, right, top := math.Inf(1), math.Inf(1), math.Inf(-1), math.Inf(-1)
	for _, quad := range q {
		left = math.Min(left, quad.Left())
		bottom = math.Min(bottom, quad.Bottom())
		right = math.Min(right, quad.Right())
		top = math.Min(top, quad.Top())
	}
	return MakeQuad(left, bottom, right, top)
}

func (q Quads) Orientation() (orientation Orientation) {
	if len(q) == 0 {
		return OtherOrientation
	}

	orientation = q[0].Orientation()
	if orientation == OtherOrientation {
		return
	}

	for _, quad := range q[1:] {
		if quad.Orientation() != orientation {
			return OtherOrientation
		}
	}
	return
}

func (q Quads) Union() (u Quad) {
	var left, bottom, right, top float64
	switch q.Orientation() {
	case Horizontal:
		left, bottom, right, top = math.Inf(1), math.Inf(1), math.Inf(-1), math.Inf(-1)
		for _, quad := range q {
			left = math.Min(left, quad.Left())
			bottom = math.Min(bottom, quad.Bottom())
			right = math.Max(right, quad.Right())
			top = math.Max(top, quad.Top())
		}
	case Rotate180:
		left, bottom, right, top = math.Inf(-1), math.Inf(-1), math.Inf(1), math.Inf(1)
		for _, quad := range q {
			right = math.Min(right, quad.Left())
			top = math.Min(top, quad.Bottom())
			left = math.Max(left, quad.Right())
			bottom = math.Max(bottom, quad.Top())
		}
	case Rotate90:
		left, bottom, right, top = math.Inf(1), math.Inf(-1), math.Inf(-1), math.Inf(1)
		for _, quad := range q {
			top = math.Min(top, quad.Left())
			left = math.Min(left, quad.Bottom())
			bottom = math.Max(bottom, quad.Right())
			right = math.Max(right, quad.Top())
		}
	case Rotate270:
		left, bottom, right, top = math.Inf(-1), math.Inf(1), math.Inf(1), math.Inf(-1)
		for _, quad := range q {
			bottom = math.Min(bottom, quad.Left())
			right = math.Min(right, quad.Bottom())
			top = math.Max(top, quad.Right())
			left = math.Max(left, quad.Top())
		}
	default:
		baselines := make([]Point, 0, len(q)*2)
		var xAvg, yAvg float64
		for _, quad := range q {
			baselines = append(baselines, quad.BottomLeft, quad.BottomRight)
			xAvg += quad.BottomLeft.X + quad.BottomRight.X
			yAvg += quad.BottomLeft.Y + quad.BottomRight.Y
		}

		xAvg /= float64(len(q) * 2.0)
		yAvg /= float64(len(q) * 2.0)

		sumProduct := 0.0
		sumDiffSquaredX := 0.0

		for i := 0; i < len(baselines); i++ {
			pt := baselines[i]
			xdiff, ydiff := pt.X-xAvg, pt.Y-yAvg
			sumProduct += xdiff * ydiff
			sumDiffSquaredX += xdiff * xdiff
		}

		var cos, sin float64 = 0, 1
		if sumDiffSquaredX > 1e-3 {
			// not a vertical line
			angle := math.Atan(sumProduct / sumDiffSquaredX) // -π/2 ≤ θ ≤ π/2
			cos = math.Cos(angle)
			sin = math.Sin(angle)
		}

		trm := NewMatrix(cos, -sin, sin, cos, 0, 0)
		pts := make([]Point, 0, len(q)*4)
		linq.From(q).SelectMany(func(i interface{}) linq.Query {
			return linq.From([]Point{i.(Quad).BottomLeft, i.(Quad).BottomRight, i.(Quad).TopLeft, i.(Quad).TopRight})
		}).Distinct().Select(func(i interface{}) interface{} { return trm.TransformPoint(i.(Point)) }).ToSlice(&pts)

		left, bottom, right, top := math.Inf(1), math.Inf(1), math.Inf(-1), math.Inf(-1)
		for _, pt := range pts {
			left = math.Min(left, pt.X)
			right = math.Min(right, pt.X)
			bottom = math.Min(bottom, pt.Y)
			top = math.Min(top, pt.Y)
		}

		aabb := MakeQuad(left, bottom, right, top)

		rotateBack := NewMatrix(cos, sin, -sin, cos, 0, 0)
		obb := Quad{
			TopLeft:     rotateBack.TransformPoint(aabb.TopLeft),
			TopRight:    rotateBack.TransformPoint(aabb.TopRight),
			BottomLeft:  rotateBack.TransformPoint(aabb.BottomLeft),
			BottomRight: rotateBack.TransformPoint(aabb.BottomRight),
		}

		obb1 := Quad{TopLeft: obb.BottomLeft, TopRight: obb.TopLeft, BottomLeft: obb.BottomRight, BottomRight: obb.TopRight}
		obb2 := Quad{TopLeft: obb.BottomRight, TopRight: obb.BottomLeft, BottomLeft: obb.TopRight, BottomRight: obb.TopLeft}
		obb3 := Quad{TopLeft: obb.TopRight, TopRight: obb.BottomRight, BottomLeft: obb.TopLeft, BottomRight: obb.BottomLeft}

		firstq := q[0]
		lastq := q[len(q)-1]

		baselineAngle := math.Atan2(lastq.BottomRight.Y-firstq.BottomLeft.Y, lastq.BottomRight.X-firstq.BottomLeft.X) * 180 / math.Pi

		deltaAngle := math.Abs(BoundAngle180(obb.Rotation() - baselineAngle))
		deltaAngle1 := math.Abs(BoundAngle180(obb1.Rotation() - baselineAngle))

		if deltaAngle1 < deltaAngle {
			deltaAngle = deltaAngle1
			obb = obb1
		}

		deltaAngle2 := math.Abs(BoundAngle180(obb2.Rotation() - baselineAngle))
		if deltaAngle2 < deltaAngle {
			deltaAngle = deltaAngle2
			obb = obb2
		}

		deltaAngle3 := math.Abs(BoundAngle180(obb3.Rotation() - baselineAngle))
		if deltaAngle3 < deltaAngle {
			obb = obb3
		}
		return obb
	}
	return MakeQuad(left, bottom, right, top)
}
