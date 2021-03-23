package fitz

import "go.matteson.dev/gfx"

type FillRule int

const (
	FillRuleEvenOdd FillRule = iota
	FillRuleWinding
)

type LineJoin int

const (
	MiterJoin LineJoin = iota
	RoundJoin
	BevelJoin
)

type LineCap int

const (
	ButtCap LineCap = iota
	RoundCap
	SquareCap
	TriangleCap
)

type BlendMode int

const (
	BlendNormal BlendMode = iota
	BlendMultiply
	BlendScreen
	BlendOverlay
	BlendDarken
	BlendLighten
	BlendColorDodge
	BlendColorBurn
	BlendHardLight
	BlendSoftLight
	BlendDifference
	BlendExclusion
	BlendHue
	BlendSaturation
	BlendColor
	BlendLuminosity
)

type ColorspaceKind int

// color spaces
const (
	ColorspaceNone ColorspaceKind = iota
	ColorspaceGray
	ColorspaceRGB
	ColorspaceBGR
	ColorspaceCMYK
	ColorspaceLab
	ColorspaceIndexed
	ColorspaceSeparation
)

// flags
const (
	ColorspaceIsDevice uint32 = 1 << iota
	ColorspaceIsICC
	ColorspaceHasCMYK
	ColorspaceHasSpots
	ColorspaceHasCMYKAndSpots
)

type Colorspace struct {
	Kind          ColorspaceKind
	Name          string
	ColorantCount int
	Flags         uint32
}

func (c Colorspace) IsSubtractive() bool {
	return c.Kind == ColorspaceCMYK || c.Kind == ColorspaceSeparation
}

func (c Colorspace) DeviceNHasOnlyCMYK() bool {
	return (c.Flags&ColorspaceHasCMYK != 0) && (c.Flags&ColorspaceHasCMYKAndSpots == 0)
}

func (c Colorspace) DeviceNHasCMYK() bool { return c.Flags&ColorspaceHasCMYK != 0 }

func (c Colorspace) IsGray() bool { return c.Kind == ColorspaceGray }

func (c Colorspace) IsRGB() bool { return c.Kind == ColorspaceBGR }

func (c Colorspace) IsCMYK() bool { return c.Kind == ColorspaceCMYK }

func (c Colorspace) IsLab() bool { return c.Kind == ColorspaceLab }

func (c Colorspace) IsIndexed() bool { return c.Kind == ColorspaceIndexed }

func (c Colorspace) IsDeviceN() bool { return c.Kind == ColorspaceSeparation }

func (c Colorspace) IsDevice() bool { return c.Flags&ColorspaceIsDevice != 0 }

func (c Colorspace) IsDeviceGray() bool { return c.IsDevice() && c.IsGray() }

func (c Colorspace) IsDeviceCMYK() bool { return c.IsDevice() && c.IsCMYK() }

func (c Colorspace) IsLabICC() bool { return c.IsLab() && c.Flags&ColorspaceIsICC != 0 }

type ShaderKind int

// Shade types
const (
	LinearShaderKind ShaderKind = iota
	RadialShaderKind
	MeshShaderKind
	FunctionShaderKind
)

type Shader struct {
	Kind   ShaderKind
	Matrix gfx.Matrix
	Bounds gfx.Rect
}

type Stroke struct {
	StartCap LineCap
	DashCap  LineCap
	EndCap   LineCap

	LineJoin  LineJoin
	LineWidth float64

	MiterLimit float64
	DashPhase  float64
	Dashes     []float64
}
