package dla

import (
	"math"

	"github.com/ahmetb/go-linq"
	"go.matteson.dev/fitz/internal/cluster"
	"go.matteson.dev/gfx"
)

// DocstrumBoundingBoxPageSegmenter ...
type DocstrumBoundingBoxPageSegmenter struct {
	options *DocstrumBoundingBoxPageSegmenterOptions
}

// NewDocstrumBoundingBoxPageSegmenter ...
func NewDocstrumBoundingBoxPageSegmenter(options ...DocstrumBoundingBoxPageSegmenterOptionsFunc) PageSegmenter {
	opts := DefaultDocstrumBoundingBoxPageSegmenterOptions()
	for _, o := range options {
		o(opts)
	}

	return &DocstrumBoundingBoxPageSegmenter{options: opts}
}

// GetBlocks implements the PageSegmenter interface
func (s *DocstrumBoundingBoxPageSegmenter) GetBlocks(words gfx.TextWords) (blocks gfx.TextBlocks) {
	cleanWords := make(gfx.TextWords, len(words))
	linq.From(words).Where(func(i interface{}) bool { return !i.(gfx.TextWord).IsWhitespace() }).ToSlice(&cleanWords)
	if len(cleanWords) == 0 {
		return
	}

	withinLineDistance, betweenLineDistance := s.getSpacingEstimation(cleanWords)
	if math.IsNaN(withinLineDistance) {
		withinLineDistance = 0
	}
	if math.IsNaN(betweenLineDistance) {
		betweenLineDistance = 0
	}

	maxWithinLineDistance := s.options.WithinLineMultiplier * withinLineDistance
	lines := s.getLines(cleanWords, maxWithinLineDistance)

	maxBetweenLineDistance := s.options.BetweenLineMultiplier * betweenLineDistance
	blocks = s.getStructuralBlocks(lines, maxBetweenLineDistance)
	return
}

func (s *DocstrumBoundingBoxPageSegmenter) getLines(words gfx.TextWords, maxWithinLineDistance float64) (lines gfx.TextLines) {
	pivotFn := func(i interface{}) gfx.Point { return i.(gfx.TextWord).Quad.BottomRight }
	distanceMeasureFn := func(i1, i2 interface{}) float64 { return cluster.EuclideanDistance(i1.(gfx.Point), i2.(gfx.Point)) }

	nodes := cluster.MakeTreeNodes(words, func(i interface{}) gfx.Point { return i.(gfx.TextWord).Quad.BottomLeft })
	indices := make([]int, len(nodes))
	for i := range indices {
		indices[i] = -1
	}

	tree := cluster.NewKDTree(nodes)

	for i, word := range words {
		for _, res := range tree.NearestNeighbors(word, 2, pivotFn, distanceMeasureFn) {
			candidate := words[res.Node.Index]
			if !s.options.WithinLineBounds.Contains(dsAngleWL(word.Quad, candidate.Quad)) {
				continue
			}

			if res.Distance >= maxWithinLineDistance {
				continue
			}

			indices[i] = res.Node.Index
		}
	}

	groups := cluster.GroupIndicesNN(indices)
	lines = make(gfx.TextLines, len(groups))
	for grpIdx, grp := range groups {
		lineWords := make(gfx.TextWords, grp.Len())
		for i, val := range grp.Values() {
			lineWords[i] = words[val]
		}
		lines[grpIdx] = gfx.MakeTextLine(lineWords.OrderByReadingOrder(), s.options.WordSeparator)
	}

	return
}

func (s *DocstrumBoundingBoxPageSegmenter) getStructuralBlocks(lines gfx.TextLines, maxBetweenLineDistance float64) (blocks gfx.TextBlocks) {
	indices := make([]int, len(lines))
	for i := range indices {
		indices[i] = -1
	}

	candidateLines := make([]gfx.Line, len(lines))
	pivotLines := make([]gfx.Line, len(lines))
	for i, line := range lines {
		pivotLines[i] = gfx.Line{Pt1: line.Quad.BottomLeft, Pt2: line.Quad.BottomRight}
		candidateLines[i] = gfx.Line{Pt1: line.Quad.TopLeft, Pt2: line.Quad.TopRight}
	}

	for i := range lines {
		distance := math.Inf(1)
		closestLineIndex := -1
		pivot := pivotLines[i]

		for j, cline := range candidateLines {
			currentDistance := s.perpendicularOverlappingDistance(pivot, cline)
			if currentDistance < distance && cline != pivot {
				distance = currentDistance
				closestLineIndex = j
			}
		}

		if closestLineIndex != -1 && distance < maxBetweenLineDistance {
			indices[i] = closestLineIndex
		}
	}

	groups := cluster.GroupIndicesNN(indices)

	blocks = make(gfx.TextBlocks, len(groups))
	for grpIdx, grp := range groups {
		blockLines := make(gfx.TextLines, grp.Len())
		for i, val := range grp.Values() {
			blockLines[i] = lines[val]
		}
		blocks[grpIdx] = gfx.MakeTextBlock(blockLines.OrderByReadingOrder(), s.options.LineSeparator)
	}
	return
}

func (s *DocstrumBoundingBoxPageSegmenter) perpendicularOverlappingDistance(line1, line2 gfx.Line) float64 {
	var theta, overlap, perpDist float64
	if s.getStructuralBlockingParameters(line1, line2, &theta, &overlap, &perpDist) {
		if theta > 90 {
			theta -= 180
		} else if theta < -90 {
			theta += 180
		}

		if !s.options.AngularDifferenceBounds.Contains(theta) {
			return math.Inf(1)
		}
		return math.Abs(perpDist)
	}
	return math.Inf(1)
}

func (s *DocstrumBoundingBoxPageSegmenter) getStructuralBlockingParameters(line1, line2 gfx.Line, angularDifference, normalizedOverlap, perpendicularDistance *float64) bool {
	linesEqual := gfx.ApproxEqual(line1.Pt1.X, line2.Pt1.X, s.options.Epsilon) &&
		gfx.ApproxEqual(line1.Pt1.Y, line2.Pt1.Y, s.options.Epsilon) &&
		gfx.ApproxEqual(line1.Pt2.X, line2.Pt2.X, s.options.Epsilon) &&
		gfx.ApproxEqual(line1.Pt2.Y, line2.Pt2.Y, s.options.Epsilon)

	if linesEqual {
		*angularDifference = 0
		*normalizedOverlap = 1
		*perpendicularDistance = 0
		return true
	}

	dXi := line1.Pt2.X - line1.Pt1.X
	dYi := line1.Pt2.Y - line1.Pt1.Y
	dXj := line2.Pt2.X - line2.Pt1.X
	dYj := line2.Pt2.Y - line2.Pt1.Y

	*angularDifference = gfx.BoundAngle180((math.Atan2(dYj, dXj) - math.Atan2(dYi, dXi)) * 180 / math.Pi)

	Aj := s.getTranslatedPoint(line1.Pt1.X, line1.Pt1.Y, line2.Pt1.X, line2.Pt1.Y, dXi, dYi, dXj, dYj)
	Bj := s.getTranslatedPoint(line1.Pt2.X, line1.Pt2.Y, line2.Pt2.X, line2.Pt2.Y, dXi, dYi, dXj, dYj)

	if Aj == nil || Bj == nil {
		// Might happen because lines are perpendicular
		// or have too small lengths
		*normalizedOverlap = math.NaN()
		*perpendicularDistance = math.NaN()
		return false
	}

	// Get middle points
	var ps = []gfx.Point{line2.Pt1, line2.Pt2, *Aj, *Bj}
	if dXj != 0 {
		linq.From(ps).OrderBy(func(i interface{}) interface{} { return i.(gfx.Point).X }).ThenBy(func(i interface{}) interface{} { return i.(gfx.Point).Y }).ToSlice(&ps)
	} else if dYj != 0 {
		linq.From(ps).OrderBy(func(i interface{}) interface{} { return i.(gfx.Point).Y }).ToSlice(&ps)
	}

	Cj, Dj := ps[1], ps[2]
	overlap := true

	// Cj and Dj should be contained within both j and [Aj,Bj] if overlapped
	if !s.pointInLine(line2.Pt1, line2.Pt2, Cj) || !s.pointInLine(line2.Pt1, line2.Pt2, Dj) ||
		!s.pointInLine(*Aj, *Bj, Cj) || !s.pointInLine(*Aj, *Bj, Dj) {
		// nonoverlapped
		overlap = false
	}

	pj := cluster.EuclideanDistance(Cj, Dj)

	*normalizedOverlap = pj
	if overlap {
		*normalizedOverlap = -pj
	}
	*normalizedOverlap /= line2.Length()

	xMj := (Cj.X + Dj.X) / 2.0
	yMj := (Cj.Y + Dj.Y) / 2.0

	if !gfx.ApproxZero(dXi, s.options.Epsilon) && !gfx.ApproxZero(dYi, s.options.Epsilon) {
		*perpendicularDistance = ((xMj - line1.Pt1.X) - (yMj-line1.Pt1.Y)*dXi/dYi) / math.Sqrt(dXi*dXi/(dYi*dYi)+1)
	} else if gfx.ApproxZero(dXi, s.options.Epsilon) {
		*perpendicularDistance = xMj - line1.Pt1.X
	} else {
		*perpendicularDistance = yMj - line1.Pt1.Y
	}

	return overlap
}

func (s *DocstrumBoundingBoxPageSegmenter) getTranslatedPoint(xPi, yPi, xPj, yPj, dXi, dYi, dXj, dYj float64) *gfx.Point {
	dYidYj := dYi * dYj
	dXidXj := dXi * dXj
	denominator := dYidYj + dXidXj

	if gfx.ApproxZero(denominator, s.options.Epsilon) {
		// The denominator is 0 when translating points, meaning the lines are perpendicular.
		return nil
	}

	var xAj float64
	var yAj float64

	if !gfx.ApproxZero(dXj, s.options.Epsilon) { // dXj != 0
		xAj = (xPi*dXidXj + xPj*dYidYj + dXj*dYi*(yPi-yPj)) / denominator
		yAj = dYj/dXj*(xAj-xPj) + yPj
	} else { // If dXj = 0, then yAj is calculated first, and xAj is calculated from that.
		yAj = (yPi*dYidYj + yPj*dXidXj + dYj*dXi*(xPi-xPj)) / denominator
		xAj = xPj
	}
	return &gfx.Point{X: xAj, Y: yAj}
}

func (s *DocstrumBoundingBoxPageSegmenter) pointInLine(pl1, pl2, point gfx.Point) bool {
	ax, ay, bx, by := point.X-pl1.X, point.Y-pl1.Y, pl2.X-pl1.X, pl2.Y-pl1.Y
	dotProd1 := ax*bx + ay*by

	if dotProd1 < 0 {
		return false
	}

	return dotProd1 <= (bx*bx + by*by)
}

func (s *DocstrumBoundingBoxPageSegmenter) getSpacingEstimation(words gfx.TextWords) (withinLineDistance, betweenLineDistance float64) {
	pivotFnBR := func(i interface{}) gfx.Point { return i.(gfx.TextWord).Quad.BottomRight }
	pivotFnTL := func(i interface{}) gfx.Point { return i.(gfx.TextWord).Quad.TopLeft }
	distanceMeasureFn := func(i1, i2 interface{}) float64 { return cluster.EuclideanDistance(i1.(gfx.Point), i2.(gfx.Point)) }

	nodes := cluster.MakeTreeNodes(words, func(i interface{}) gfx.Point { return i.(gfx.TextWord).Quad.BottomLeft })
	indices := make([]int, len(nodes))
	for i := range indices {
		indices[i] = -1
	}

	tree := cluster.NewKDTree(nodes)
	wlda := make([]float64, 0)
	blda := make([]float64, 0)

	for _, word := range words {
		for _, res := range tree.NearestNeighbors(word, 2, pivotFnBR, distanceMeasureFn) {
			if s.options.WithinLineBounds.Contains(dsAngleWL(word.Quad, words[res.Node.Index].Quad)) {
				dist := cluster.EuclideanDistance(word.Quad.BottomRight, res.Node.Point)
				wlda = append(wlda, dist)
			}
		}

		for _, res := range tree.NearestNeighbors(word, 2, pivotFnTL, distanceMeasureFn) {
			angle := dxAngleBL(word.Quad, words[res.Node.Index].Quad)
			if s.options.BetweenLineBounds.Contains(angle) {
				hypotenuse := cluster.EuclideanDistance(word.Quad.Centroid(), words[res.Node.Index].Quad.Centroid())
				if angle > 90 {
					angle -= 180
				}
				dist := math.Abs(hypotenuse*math.Cos((90-angle)*math.Pi/180)) - word.Quad.Height()/2.0 - words[res.Node.Index].Quad.Height()/2.0
				if dist >= 0 {
					blda = append(blda, dist)
				}
			}
		}
	}

	withinLineDistance = s.getPeakAverageDistance(wlda, s.options.WithinLineBinSize)
	betweenLineDistance = s.getPeakAverageDistance(blda, s.options.BetweenLineBinSize)
	return
}

func (s *DocstrumBoundingBoxPageSegmenter) getPeakAverageDistance(distances []float64, binLength int) float64 {
	if len(distances) == 0 {
		return math.NaN()
	}

	max := math.Inf(-1)
	for _, dist := range distances {
		if dist > max {
			max = dist
		}
	}

	max = math.Ceil(max)
	if max == 0 {
		max = float64(binLength)
	} else if binLength > int(max) {
		binLength = int(max)
	}

	bins := make(map[int][]float64)
	linq.Range(0, int(math.Ceil(max/float64(binLength))+1)).
		Select(func(i interface{}) interface{} { return i.(int) * binLength }).
		ToMapBy(&bins,
			func(i interface{}) interface{} { return i.(int) },
			func(i interface{}) interface{} { return make([]float64, 0) },
		)

	binKeys := make([]int, 0, len(bins))
	for i := range bins {
		binKeys = append(binKeys, i)
	}

	for _, dist := range distances {
		bin := int(math.Floor(dist / float64(binLength)))
		bins[binKeys[bin]] = append(bins[binKeys[bin]], dist)
	}

	var best []float64
	for _, val := range bins {
		if len(val) > len(best) {
			best = val
		}
	}
	return linq.From(best).Average()
}

func dsAngleWL(pivot, candidate gfx.Quad) float64 {
	angle := gfx.BoundAngle180(gfx.LineAngle(pivot.BottomRight, candidate.BottomLeft) - pivot.Rotation())
	if angle > 90 {
		angle -= 180
	} else if angle < -90 {
		angle += 180
	}
	return angle
}

func dxAngleBL(pivot, candidate gfx.Quad) float64 {
	angle := gfx.BoundAngle180(gfx.LineAngle(pivot.Centroid(), candidate.Centroid()) - pivot.Rotation())
	if angle < 0 {
		angle += 180
	}
	return angle
}

// DocstrumBoundingBoxPageSegmenterOptions ...
type DocstrumBoundingBoxPageSegmenterOptions struct {
	*PageSegmenterOptions
	Epsilon                 float64
	WithinLineBounds        gfx.Range
	WithinLineMultiplier    float64
	WithinLineBinSize       int
	BetweenLineBounds       gfx.Range
	BetweenLineMultiplier   float64
	BetweenLineBinSize      int
	AngularDifferenceBounds gfx.Range
}

// DefaultDocstrumBoundingBoxPageSegmenterOptions ...
func DefaultDocstrumBoundingBoxPageSegmenterOptions() *DocstrumBoundingBoxPageSegmenterOptions {
	return &DocstrumBoundingBoxPageSegmenterOptions{
		PageSegmenterOptions: &PageSegmenterOptions{
			WordSeparator: " ",
			LineSeparator: "\n",
		},
		Epsilon: 1e-3,

		WithinLineMultiplier:    3.0,
		BetweenLineMultiplier:   1.3,
		WithinLineBinSize:       10,
		BetweenLineBinSize:      10,
		WithinLineBounds:        gfx.Range{Min: -30, Max: 30},
		BetweenLineBounds:       gfx.Range{Min: 45, Max: 135},
		AngularDifferenceBounds: gfx.Range{Min: -30, Max: 30},
	}
}

// DocstrumBoundingBoxPageSegmenterOptionsFunc ...
type DocstrumBoundingBoxPageSegmenterOptionsFunc func(*DocstrumBoundingBoxPageSegmenterOptions)

// WithDocstrumLineSeparator ...
func WithDocstrumLineSeparator(ls string) DocstrumBoundingBoxPageSegmenterOptionsFunc {
	return func(o *DocstrumBoundingBoxPageSegmenterOptions) {
		o.LineSeparator = ls
	}
}

// WithDocstrumWordSeparator ...
func WithDocstrumWordSeparator(ws string) DocstrumBoundingBoxPageSegmenterOptionsFunc {
	return func(o *DocstrumBoundingBoxPageSegmenterOptions) {
		o.WordSeparator = ws
	}
}

// WithDocstrumEpsilon ...
func WithDocstrumEpsilon(value float64) DocstrumBoundingBoxPageSegmenterOptionsFunc {
	return func(o *DocstrumBoundingBoxPageSegmenterOptions) {
		o.Epsilon = value
	}
}

// WithDocstrumWithinLineBounds ...
func WithDocstrumWithinLineBounds(value gfx.Range) DocstrumBoundingBoxPageSegmenterOptionsFunc {
	return func(o *DocstrumBoundingBoxPageSegmenterOptions) {
		o.WithinLineBounds = value
	}
}

// WithDocstrumWithinLineMultiplier ...
func WithDocstrumWithinLineMultiplier(value float64) DocstrumBoundingBoxPageSegmenterOptionsFunc {
	return func(o *DocstrumBoundingBoxPageSegmenterOptions) {
		o.WithinLineMultiplier = value
	}
}

// WithDocstrumWithinLineBinSize ...
func WithDocstrumWithinLineBinSize(value int) DocstrumBoundingBoxPageSegmenterOptionsFunc {
	return func(o *DocstrumBoundingBoxPageSegmenterOptions) {
		o.WithinLineBinSize = value
	}
}

// WithDocstrumBetweenLineBounds ...
func WithDocstrumBetweenLineBounds(value gfx.Range) DocstrumBoundingBoxPageSegmenterOptionsFunc {
	return func(o *DocstrumBoundingBoxPageSegmenterOptions) {
		o.BetweenLineBounds = value
	}
}

// WithDocstrumBetweenLineMultiplier ...
func WithDocstrumBetweenLineMultiplier(value float64) DocstrumBoundingBoxPageSegmenterOptionsFunc {
	return func(o *DocstrumBoundingBoxPageSegmenterOptions) {
		o.BetweenLineMultiplier = value
	}
}

// WithDocstrumBetweenLineBinSize ...
func WithDocstrumBetweenLineBinSize(value int) DocstrumBoundingBoxPageSegmenterOptionsFunc {
	return func(o *DocstrumBoundingBoxPageSegmenterOptions) {
		o.BetweenLineBinSize = value
	}
}

// WithDocstrumAngularDifferenceBounds ...
func WithDocstrumAngularDifferenceBounds(value gfx.Range) DocstrumBoundingBoxPageSegmenterOptionsFunc {
	return func(o *DocstrumBoundingBoxPageSegmenterOptions) {
		o.AngularDifferenceBounds = value
	}
}
