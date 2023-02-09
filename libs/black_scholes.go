package libs

import (
	"math"
	"gonum.org/v1/gonum/stat/distuv"
)

func d1(S, K, T, r, sigma, q float64) float64 {
	return (math.Log(S/K) + (r-q+0.5*sigma*sigma)*T) / sigma / math.Sqrt(T)
}

func d2(S, K, T, r, sigma, q float64) float64 {
	return d1(S, K, T, r, sigma, q ) - sigma * math.Sqrt(T)
}

func Bs(tipo string, S, K, T, r, sigma, q float64) float64 {

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
	- div : float - Tasa de dividendos continuos (anualizada)
	Outputs
	- precio_BS: float - Precio del contrato
	*/
	// Create a normal distribution
	dist := distuv.Normal{
		Mu:    0,
		Sigma: 1,
	}

	var price float64

	d1 := d1(S, K, T, r, sigma, q)
	d2 := d2(S, K, T, r, sigma, q)

	if tipo == "C" {
		price = math.Exp(-q*T)*S*dist.CDF(d1) - K*math.Exp(-r*T)*dist.CDF(d2)
	} else if tipo == "P" {
		price = K*math.Exp(-r*T)*dist.CDF(-d2) - S*math.Exp(-q*T)*dist.CDF(-d1)
	}
	return price
}

func Delta(tipo string, S, K, T, r, sigma, q float64) float64 {
	/*
	Delta is the first derivative of option price with respect to underlying price S.
	*/

	// Create a normal distribution
	dist := distuv.Normal{
		Mu:    0,
		Sigma: 1,
	}

	var delta float64

	d1 := d1(S, K, T, r, sigma, q)

	if tipo == "C" {
		delta = math.Exp(-q*T) * dist.CDF(d1)
	} else if tipo == "P" {
		delta = -math.Exp(-q*T) * dist.CDF(-d1)
	}

	return delta
}

func np(x float64) float64{
	return math.Exp(-math.Pow(x, 2)/2) / math.Sqrt((2*math.Pi))
}

func Gamma(S, K, T, r, sigma, q float64) float64 {
	/*
	Gamma is the second derivative of option price with
	respect to underlying price S. It is the same for calls and puts.
	*/
	d1 := d1(S, K, T, r, sigma, q)
	return math.Exp(-q*T) * np(d1) / (S*sigma*math.Sqrt(T))
}

func Theta(tipo string, S, K, T, r, sigma, q float64, calendar bool) float64 {
	/*
	Theta is the first derivative of option price with respect to time to expiration T.
	*/
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

	d1 := d1(S, K, T, r, sigma, q)
	d2 := d2(S, K, T, r, sigma, q)

	gamma := Gamma(S, K, T, r, sigma, q)
	temp := gamma * math.Pow(S*sigma,2) / 2

	if tipo == "C" {
		theta = (-temp - r*K*math.Exp(-r*T)*dist.CDF(d2) + q*S*math.Exp(-q*T)*dist.CDF(d1))
	} else if tipo == "P" {
		theta = (-temp + r*K*math.Exp(-r*T)*dist.CDF(-d2) - q*S*math.Exp(-q*T)*dist.CDF(-d1))
	}
	return theta / NDY
}

func Rho(tipo string, S, K, T, r, sigma, q float64) float64 {
	/*
	Rho is the first derivative of option price with respect to interest rate r
	*/

	dist := distuv.Normal{
		Mu:    0,
		Sigma: 1,
	}

	d2 := d2(S, K, T, r, sigma, q)

	var rho float64

	if tipo == "C" {
		rho = K * T * math.Exp(-r*T)*dist.CDF(d2) / 100
	} else if tipo == "P" {
		rho = -K * T * math.Exp(-r*T)*dist.CDF(-d2) / 100
	}

	return rho
}

func Vega(S, K, T, r, sigma, q float64) float64{
	/*
    Vega is the first derivative of option price with respect to volatility σ.
	It is the same for calls and puts.
	*/
	d1 := d1(S, K, T, r, sigma, q)
	return S * math.Exp(-q*T) * math.Sqrt(T) * np(d1)

}

func IvBsNewton(tipo string, S, K, T, r, q, price, sigma0 float64) float64 {
	var (
		price0 float64
		vega0 float64
	)
	for {
		price0 = Bs(tipo, S, K, T, r, sigma0, q)
		vega0 = Vega(S, K, T, r, sigma0, q)
		sigma0 = sigma0 - ((price0 - price)/ vega0)
		if (price0 - price) < 0.01 {
			break
		}
	}
	return sigma0
}

func IV_Bs(tipo string, S, K, T, r, q, price float64) float64 {
	var diff float64

	s_high := 10.0
	s_low := .0001
	sigma := .5 * (s_low + s_high)

	for i := 0; i < 1000; i++ {
		diff = Bs(tipo, S, K, T, r, sigma, q) - price
		if diff > 0 {
			s_high = sigma
			sigma = .5 * (s_low + s_high)
		} else {
			s_low = sigma
			sigma = .5 * (s_low + s_high)
		}

		if math.Abs(diff) > .0001 {
			continue
		} else {
			break
		}
	}

	return sigma
}

type Parameters struct {
	Tipo           string
	S, K, T, R, Sigma, Q float64
}

func (p *Parameters) Price(method string, steps int) float64{
	if method == "BIN" {
		return Bin(p, steps)
	} else {
		return Bs(p.Tipo, p.S, p.K, p.T, p.R, p.Sigma, p.Q)
	}
}
