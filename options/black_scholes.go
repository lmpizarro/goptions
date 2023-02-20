package libs

import (
	"math"

	"gonum.org/v1/gonum/stat/distuv"
)

func d1(p *OptionsParameters) float64 {
	return (math.Log(p.S/p.K) + (p.R-p.Q+0.5*p.Sigma*p.Sigma)*p.T) / p.Sigma / math.Sqrt(p.T)
}

func d2(p *OptionsParameters) float64 {
	return d1(p) - p.Sigma*math.Sqrt(p.T)
}

func np(x float64) float64 {
	return math.Exp(-math.Pow(x, 2)/2) / math.Sqrt((2 * math.Pi))
}

// Option Price Calculator by Black-Scholes-Merton Model
//
//	In: *Parameters {S, K, R, T, Sigma, Q}
//	Out: price float64
func Bs(p *OptionsParameters) float64 {

	/*
		Def
		Calculador del precio de una opcion Europea con el modelo de Black Scholes
		Inputs
		- tipo : string - Tipo de contrato entre ["CALL","PUT"]
		- S : float - Spot price del activo
		- K : float - Strike price del contrato
		- T : float - Tiempo hasta la expiracion (en años)
		- r: Tasa 'libre de riesgo' (anualizada)
		- sigma : float - Volatilidad implicita (anualizada)
		- q : float - Tasa de dividendos continuos (anualizada)
		Outputs
		- precio_BS: float - Precio del contrato
	*/
	// Create a normal distribution
	dist := distuv.Normal{
		Mu:    0,
		Sigma: 1,
	}

	var price float64

	d1 := d1(p)
	d2 := d2(p)

	if p.Tipo == "C" {
		price = math.Exp(-p.Q*p.T)*p.S*dist.CDF(d1) - p.K*math.Exp(-p.R*p.T)*dist.CDF(d2)
	} else if p.Tipo == "P" {
		price = p.K*math.Exp(-p.R*p.T)*dist.CDF(-d2) - p.S*math.Exp(-p.Q*p.T)*dist.CDF(-d1)
	}
	return price
}

// Delta is the first derivative of option price with respect to underlying price S.
//
//	In: *Parameters {S, K, R, T, Sigma, Q}
//	Out: delta float64
func Delta(p *OptionsParameters) float64 {
	// Create a normal distribution
	dist := distuv.Normal{
		Mu:    0,
		Sigma: 1,
	}

	var delta float64

	d1 := d1(p)

	if p.Tipo == "C" {
		delta = math.Exp(-p.Q*p.T) * dist.CDF(d1)
	} else if p.Tipo == "P" {
		delta = -math.Exp(-p.Q*p.T) * dist.CDF(-d1)
	}

	return delta
}

// Gamma is the second derivative of option price with
// respect to underlying price S. It is the same for calls and puts.
//
//	In: *Parameters {S, K, R, T, Sigma, Q}
//	Out: gamma float64
func Gamma(p *OptionsParameters) float64 {

	d1 := d1(p)
	return math.Exp(-p.Q*p.T) * np(d1) / (p.S * p.Sigma * math.Sqrt(p.T))
}

// Theta is the first derivative of option price with respect to time to expiration T.
//
//	In: *Parameters {S, K, R, T, Sigma, Q}
//		calendar bool
//	Out: theta float64
func Theta(p *OptionsParameters, calendar bool) float64 {

	var theta float64
	var NDY float64

	if calendar {
		NDY = 365.0
	} else {
		NDY = 252.0
	}

	dist := distuv.Normal{
		Mu:    0,
		Sigma: 1,
	}

	d1 := d1(p)
	d2 := d2(p)

	gamma := Gamma(p)
	temp := gamma * math.Pow(p.S*p.Sigma, 2) / 2

	if p.Tipo == "C" {
		theta = (-temp - p.R*p.K*math.Exp(-p.R*p.T)*dist.CDF(d2) + p.Q*p.S*math.Exp(-p.Q*p.T)*dist.CDF(d1))
	} else if p.Tipo == "P" {
		theta = (-temp + p.R*p.K*math.Exp(-p.R*p.T)*dist.CDF(-d2) - p.Q*p.S*math.Exp(-p.Q*p.T)*dist.CDF(-d1))
	}
	return theta / NDY
}

// Rho is the first derivative of option price with respect to interest rate r
//
//	In: *Parameters {S, K, R, T, Sigma, Q}
//	Out: rho float64
func Rho(p *OptionsParameters) float64 {
	dist := distuv.Normal{
		Mu:    0,
		Sigma: 1,
	}

	d2 := d2(p)

	var rho float64

	if p.Tipo == "C" {
		rho = p.K * p.T * math.Exp(-p.R*p.T) * dist.CDF(d2) / 100
	} else if p.Tipo == "P" {
		rho = -p.K * p.T * math.Exp(-p.R*p.T) * dist.CDF(-d2) / 100
	}

	return rho
}

//	   Vega is the first derivative of option price with respect to volatility σ.
//		It is the same for calls and puts.
//			In: *Parameters {S, K, R, T, Sigma, Q}
//			Out: vega float64
func Vega(p *OptionsParameters) float64 {
	d1 := d1(p)
	return p.S * math.Exp(-p.Q*p.T) * math.Sqrt(p.T) * np(d1)

}


// OptionsParameters
//
//	S: underlying Price
//	K: Strike
//	T: time to maturity
//	R: risk free rate
//	Sigma: volatility
//	Q: dividend yield
//	Tipo: "C" for call "P" for put
type OptionsParameters struct {
	Tipo                 string
	S, K, T, R, Sigma, Q float64
}

func (p *OptionsParameters) Price(method string, steps int) float64 {
	if method == "BIN" {
		return Bin(p, steps)
	} else {
		return Bs(p)
	}
}
