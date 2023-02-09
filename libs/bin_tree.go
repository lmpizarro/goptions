package libs

import (
	"math"
)

func Bin(p *Parameters, steps int) float64 {

	dt := p.T / float64(steps)
	tasa_forward := math.Exp((p.R - p.Q) * dt)
	descuento := math.Exp(-p.R * dt)
	u := math.Exp(p.Sigma * math.Pow(dt, 0.5))
	d := 1 / u
	q_prob := (tasa_forward - d) / (u - d)

	ST_precios := make([]float64, steps+1)
	for i := 0; i < steps+1; i++ {
		ST_precios[steps-i] = math.Pow(u, 2*float64(i)-float64(steps)) * p.S
	}

	options_matrix := make([][]float64, steps+1)
	for i := range options_matrix {
		options_matrix[i] = make([]float64, steps+1)
	}

	for i := 0; i < steps+1; i++ {
		if p.Tipo == "P" {
			options_matrix[i][steps] = math.Max(0, (p.K - ST_precios[i]))
		} else {
			options_matrix[i][steps] = math.Max(0, -(p.K - ST_precios[i]))
		}
	}

	for j := 1; j < steps+1; j++ {
		for i := 0; i < steps+1-j; i++ {
			eur := q_prob*options_matrix[i][steps-j+1] + (1-q_prob)*options_matrix[i+1][steps-j+1]
			if p.Tipo == "P" {
				options_matrix[i][steps-j] = descuento * math.Max(eur, p.K-p.S*math.Pow(u, float64(-2*i+steps-j)))
			} else {
				options_matrix[i][steps-j] = descuento * math.Max(eur, -(p.K-p.S*math.Pow(u, float64(-2*i+steps-j))))

			}
		}
	}
	return options_matrix[0][0]
}

