package main

import (
	"fmt"
	"lmpizarro/options/libs"
)

func main() {

	day := 1.0 / 365
	t := 11.0 / 365.0
	pinit := 406.0
	pfinal := 423.5
	delta_precio := 0.5
	// cost := 2.34
	cinit := libs.Bs(&libs.OptionsParameters{Tipo: "C", S: 413.98, K: 420, T: t, R: 0.045, Sigma: 0.15, Q: 0.015})
	fmt.Println(cinit)

	var c float64
	var values []float64
	var price float64
	price = pinit
	// See https://www.optionsprofitcalculator.com/calculator/long-call.html
	for {
		for {
			c = libs.Bs(&libs.OptionsParameters{Tipo: "C", S: price, K: 420, T: t, R: 0.045, Sigma: 0.15, Q: 0.015})
			// values = append(values, libs.Round_down(100*(c - cinit)/cinit, 2))
			values = append(values, libs.Round_down(c-cinit, 2))
			if price > pfinal {
				break
			}
			price = price + delta_precio
		}
		fmt.Println(libs.Round_down(365*t, 1), values, price)
		values = values[:0]
		price = pinit
		t = t - day
		if t < day {
			break
		}
	}

	panic("")
	// libs.Parallel_Calc_IV("SPY")
	libs.Test_YF()

	fmt.Println("Tests")
	params := libs.OptionsParameters{S: 100.0, K: 100.0,
		Tipo: "P", T: 1, Sigma: .4, Q: 0.01, R: 0.04}
	fmt.Println(libs.Delta(&params))
	fmt.Println(libs.Gamma(&params))

	deltaBin := libs.DeltaBin(&params, 150)
	fmt.Println(deltaBin)

	params.S = 100.0
	gammaBin := libs.GammaBin(&params, 150)
	fmt.Println(gammaBin)

	fmt.Println(params)

	params.S = 100.0
	params.K = 150.0
	params.Tipo = "C"

	for _, S := range []float64{105, 110, 115, 120, 125, 130, 135, 140, 145, 150} {
		// fmt.Println("---------------------------")
		for _, T := range []float64{0.9, 0.8, 0.7, .6, 0.5, 0.4, 0.3, 0.2, 0.1, 0.05} {
			// fmt.Println(S, T, libs.Delta(&params))
			params.T = T
		}
		params.S = S
	}
}
