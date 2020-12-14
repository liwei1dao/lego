package floats

import "gonum.org/v1/gonum/floats"

func MulConstTo(a []float64, c float64, dst []float64) {
	floats.ScaleTo(dst, c, a)
}

func MulConstAddTo(a []float64, c float64, dst []float64) {
	floats.AddScaled(dst, c, a)
}

func Dot(a, b []float64) float64 {
	return floats.Dot(a, b)
}
