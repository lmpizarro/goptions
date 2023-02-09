package libs

import (
	"math"
	"gonum.org/v1/gonum/stat/distuv"
)

func d1(S, K, T, r, sigma, div float64) float64 {
	return (math.Log(S/K) + (r-div+0.5*sigma*sigma)*T) / sigma / math.Sqrt(T)
}

func d2(S, K, T, r, sigma, div float64) float64 {
	return d1(S, K, T, r, sigma, div ) - sigma * math.Sqrt(T)
}

func Bs(tipo string, S, K, T, r, sigma, div float64) float64 {

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

	d1 := d1(S, K, T, r, sigma, div)
	d2 := d2(S, K, T, r, sigma, div)

	if tipo == "C" {
		price = math.Exp(-div*T)*S*dist.CDF(d1) - K*math.Exp(-r*T)*dist.CDF(d2)
	} else if tipo == "P" {
		price = K*math.Exp(-r*T)*dist.CDF(-d2) - S*math.Exp(-div*T)*dist.CDF(-d1)
	}
	return price
}

func Delta(tipo string, S, K, T, r, sigma, div float64) float64 {
	// Create a normal distribution
	dist := distuv.Normal{
		Mu:    0,
		Sigma: 1,
	}

	var delta float64

	d1 := d1(S, K, T, r, sigma, div)

	if tipo == "C" {
		delta = math.Exp(-div*T) * dist.CDF(d1)
	} else {
		delta = -math.Exp(-div*T) * dist.CDF(-d1)
	}

	return delta
}

func np(x float64) float64{
	return math.Exp(-math.Pow(x, 2)/2) / math.Sqrt((2*math.Pi))
}

func Gamma(S, K, T, r, sigma, div float64) float64 {
	d1 := d1(S, K, T, r, sigma, div)
	return math.Exp(-div*T) * np(d1) / (S*sigma*math.Sqrt(T))
}

func Theta(tipo string, S, K, T, r, sigma, div float64, NDY int) float64 {

	var theta float64

	dist := distuv.Normal{
		Mu:    0,
		Sigma: 1,
	}

	d1 := d1(S, K, T, r, sigma, div)
	d2 := d2(S, K, T, r, sigma, div)

	gamma := Gamma(S, K, T, r, sigma, div)
	temp := gamma * math.Pow(S*sigma,2) / 2

	if tipo == "C" {
		theta = (1.0/float64(NDY)) * (-temp - r*K*math.Exp(-r*T)*dist.CDF(d2) + div*S*math.Exp(-div*T)*dist.CDF(d1))
	} else {
		theta = (1.0/float64(NDY)) * (-temp + r*K*math.Exp(-r*T)*dist.CDF(-d2) - div*S*math.Exp(-div*T)*dist.CDF(-d1))
	}
	return theta
}

func Rho(tipo string, S, K, T, r, sigma, div float64) float64 {
	/*
	Rho is the first derivative of option price with respect to interest rate r
	*/

	dist := distuv.Normal{
		Mu:    0,
		Sigma: 1,
	}

	d2 := d2(S, K, T, r, sigma, div)

	var rho float64

	if tipo == "C" {
		rho = K * T * math.Exp(-r*T)*dist.CDF(d2) / 100
	} else {
		rho = -K * T * math.Exp(-r*T)*dist.CDF(-d2) / 100
	}

	return rho
}
func Vega(S, K, T, r, sigma, div float64) float64{

	d1 := d1(S, K, T, r, sigma, div)
	return S * math.Exp(-div*T) * math.Sqrt(T) * np(d1)

}

func IvBsNewton(tipo string, S, K, T, r, div, price, sigma0 float64) float64 {
	var (
		price0 float64
		vega0 float64
	)
	for {
		price0 = Bs(tipo, S, K, T, r, sigma0, div)
		vega0 = Vega(S, K, T, r, sigma0, div)
		sigma0 = sigma0 - ((price0 - price)/ vega0)
		if (price0 - price) < 0.01 {
			break
		}
	}
	return sigma0
}

func IV_Bs(tipo string, S, K, T, r, div, price float64) float64 {
	var diff float64

	s_high := 10.0
	s_low := .0001
	sigma := .5 * (s_low + s_high)

	for i := 0; i < 1000; i++ {
		diff = Bs(tipo, S, K, T, r, sigma, div) - price
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

func Bin(tipo string, S, K, T, r, sigma, div float64, steps int) float64 {

	dt := T / float64(steps)
	tasa_forward := math.Exp((r - div) * dt)
	descuento := math.Exp(-r * dt)
	u := math.Exp(sigma * math.Pow(dt, 0.5))
	d := 1 / u
	q_prob := (tasa_forward - d) / (u - d)

	ST_precios := make([]float64, steps+1)
	for i := 0; i < steps+1; i++ {
		ST_precios[steps-i] = math.Pow(u, 2*float64(i)-float64(steps)) * S
	}

	options_matrix := make([][]float64, steps+1)
	for i := range options_matrix {
		options_matrix[i] = make([]float64, steps+1)
	}

	for i := 0; i < steps+1; i++ {
		if tipo == "P" {
			options_matrix[i][steps] = math.Max(0, (K - ST_precios[i]))
		} else {
			options_matrix[i][steps] = math.Max(0, -(K - ST_precios[i]))
		}
	}

	for j := 1; j < steps+1; j++ {
		for i := 0; i < steps+1-j; i++ {
			eur := q_prob*options_matrix[i][steps-j+1] + (1-q_prob)*options_matrix[i+1][steps-j+1]
			if tipo == "P" {
				options_matrix[i][steps-j] = descuento * math.Max(eur, K-S*math.Pow(u, float64(-2*i+steps-j)))
			} else {
				options_matrix[i][steps-j] = descuento * math.Max(eur, -(K-S*math.Pow(u, float64(-2*i+steps-j))))

			}
		}
	}
	return options_matrix[0][0]
}


type Parameters struct {
	Tipo, Method           string
	S, K, T, R, Sigma, Div float64
	Steps                  int
}

func (p Parameters) Price() float64{
	if p.Method == "BIN" {
		return Bin(p.Tipo, p.S, p.K, p.T, p.R, p.Sigma, p.Div, p.Steps)
	} else {
		return Bs(p.Tipo, p.S, p.K, p.T, p.R, p.Sigma, p.Div)
	}
}
