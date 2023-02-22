package libs

import (
	"math"

	"errors"
)

func IV_Brenner_Subrahmanyam(p *OptionParameters, C float64) float64 {

	X := p.K * math.Exp(-p.R*p.T)
	d := (p.S - X) / 2
	return math.Sqrt(2*math.Pi/p.T) * (C - d) / p.S
}

func IV_Bharadia_Christofides_Salkin(p *OptionParameters, C float64) float64 {

	X := p.K * math.Exp(-p.R*p.T)
	d := (p.S - X) / 2
	return math.Sqrt(2*math.Pi/p.T) * (C - d) / (p.S - d)
}

func IV_Corrado_Miller(p *OptionParameters, C float64) (float64, error) {
	X := p.K * math.Exp(-p.R*p.T)
	d := (p.S - X) / 2
	l1 := d * (1 - 2/math.Sqrt(math.Pi))
	l2 := d * (1 + 2/math.Sqrt(math.Pi))
	if C > l2 || C < l1 {
		return (math.Sqrt(2*math.Pi/p.T) / (p.S + X)) *
			(C - d + math.Sqrt(math.Pow(C-d, 2)-(d*d/math.Pi))), nil
	} else {
		return 0.0, errors.New("bad limit")
	}
}

// Solves Black-Scholes-Merton Implied Volatility by Newton Rapshon method
func IvBsNewton(p *OptionParameters, sigma0, price, tol float64) (int, float64) {
	var (
		price0 float64
		vega0  float64
	)
	i := 0
	for {
		i++
		p.Sigma = sigma0
		price0 = p.Bs()
		vega0 = p.Vega()
		sigma0 = sigma0 - ((price0 - price) / vega0)
		if math.Abs(price0-price) < tol {
			break
		}
	}
	return i, sigma0
}

// Solves Black-Scholes-Merton Implied Volatility
// by the secant method
func IvBsSecant(p *OptionParameters, price float64) (int, float64) {
	var x2 float64
	var i int

	x1 := 1.0
	x0 := .0001
	steps := 10
	p.Sigma = x1
	f1 := p.func_diff_bs(price)
	p.Sigma = x0
	f0 := p.func_diff_bs(price)

	for i = 0; i < steps; i++ {
		x2 = x1 - f1*(x1-x0)/(f1-f0)

		x0 = x1
		x1 = x2

		if math.Abs(f0-f1) < 0.000001 {
			panic("convergence error")
		}

		if math.Abs(x2-x0) < 0.000001 {
			break
		}

		f0 = f1
		p.Sigma = x1
		f1 = p.func_diff_bs(price)
	}
	return i, x1
}

// Solves Black-Scholes-Merton Implied Volatility
func IvBsBisection_A(p *OptionParameters, price float64) (int, float64) {
	var diff float64
	var i int

	s_high := 10.0
	s_low := .0001
	sigma := .5 * (s_low + s_high)
	steps := 100

	for i = 0; i < steps; i++ {
		p.Sigma = sigma
		diff = p.func_diff_bs(price)
		if diff > 0 {
			s_high = sigma
			sigma = .5 * (s_low + s_high)
		} else {
			s_low = sigma
			sigma = .5 * (s_low + s_high)
		}

		if math.Abs(diff) > .01 {
			continue
		} else {
			break
		}
	}

	return i, sigma
}

func samesign(a, b float64) bool {
	return math.Signbit(a) == math.Signbit(b)
}

func (p *OptionParameters) func_diff_bs(price float64) float64 {
	return p.Bs() - price
}
