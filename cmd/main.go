package main

import (
	"fmt"
	"lmpizarro/options/libs"
)


func main() {

	// fmt.Println(libs.TestNewton())
	// libs.Parallel_Calc_IV("SPY")
	libs.Test_YF()
	panic("main")

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
