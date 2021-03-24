package cluster

import "go.matteson.dev/gfx"

// PivotFilterFunc ...
type PivotFilterFunc func(interface{}) bool

// CandidateFilterFunc ...
type CandidateFilterFunc func(interface{}, interface{}) bool

// KDTreeConfig ...
type KDTreeConfig struct {
	DistanceFn        func(interface{}, interface{}) float64
	MaxDistanceFn     func(interface{}, interface{}) float64
	PivotPointFn      func(interface{}) gfx.Point
	CandidatePointFn  func(interface{}) gfx.Point
	PivotFilterFn     func(interface{}) bool
	CandidateFilterFn func(interface{}, interface{}) bool
}

// Distance implements the KDTreeConfigurator interface
func (c *KDTreeConfig) Distance(a interface{}, b interface{}) float64 { return c.DistanceFn(a, b) }

// MaxDistance implements the KDTreeConfigurator interface
func (c *KDTreeConfig) MaxDistance(a interface{}, b interface{}) float64 {
	return c.MaxDistanceFn(a, b)
}

// PivotPoint implements the KDTreeConfigurator interface
func (c *KDTreeConfig) PivotPoint(a interface{}) gfx.Point { return c.PivotPointFn(a) }

// CandidatePoint implements the KDTreeConfigurator interface
func (c *KDTreeConfig) CandidatePoint(a interface{}) gfx.Point { return c.CandidatePointFn(a) }

// PivotFilter implements the KDTreeConfigurator interface
func (c *KDTreeConfig) PivotFilter(a interface{}) bool { return c.PivotFilterFn(a) }

// CandidateFilter implements the KDTreeConfigurator interface
func (c *KDTreeConfig) CandidateFilter(a interface{}, b interface{}) bool {
	return c.CandidateFilterFn(a, b)
}

// KDTreeConfigurator ...
type KDTreeConfigurator interface {
	Distance(interface{}, interface{}) float64
	MaxDistance(interface{}, interface{}) float64
	PivotPoint(interface{}) gfx.Point
	CandidatePoint(interface{}) gfx.Point
	PivotFilter(interface{}) bool
	CandidateFilter(interface{}, interface{}) bool
}
