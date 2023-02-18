package main

import (
	"fmt"
	"lmpizarro/options/libs"
)


func main() {

	t := 11.0 / 365.0
	opt_params := libs.OptionsParameters{Tipo: "C", S: 413.98, K: 420, T: t, R: 0.045, Sigma: 0.15, Q: 0.015}


	fmt.Println(libs.Default_simulate_parameters(&opt_params))
	simul_params := libs.SimulationParameters{Price_increment: .5, End_price: 423.5, Init_price: 406.0}
	rows := libs.Simulate_long(&opt_params, &simul_params, true)

	fmt.Println(string(libs.Rows_simulation_to_json(rows)))
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
