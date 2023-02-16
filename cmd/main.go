package main

import (
	"fmt"
	"lmpizarro/options/libs"
)

const Day = 1.0 / 365.0

type SimulationParameters struct {
	Init_price float64
	End_price float64
	Price_increment float64
}

func simulate(option_params *libs.OptionsParameters,
	simul_params *SimulationParameters){
	day := Day
	t := option_params.T
	pinit := simul_params.Init_price // ex 406.0
	pfinal := simul_params.End_price // ex 423.5
	delta_precio := simul_params.Price_increment // ex .5
	cost_init := libs.Bs(option_params)

	var price_of_option float64
	var values []float64
	var price_of_equity float64
	price_of_equity = pinit
	// See https://www.optionsprofitcalculator.com/calculator/long-call.html
	// See https://optionstrat.com/
	for {
		for {
			option_params.T = t
			option_params.S = price_of_equity
			price_of_option = libs.Bs(option_params)
			// values = append(values, libs.Round_down(100*(c - cinit)/cinit, 2))
			values = append(values, libs.Round_down(price_of_option-cost_init, 2))
			if price_of_equity > pfinal {
				break
			}
			price_of_equity = price_of_equity + delta_precio
		}
		fmt.Println(libs.Round_down(365*t, 1), values, price_of_equity)
		values = values[:0]
		price_of_equity = pinit
		t = t - day
		if t < day {
			break
		}
	}
}

func main() {


	t := 11.0 / 365.0
	opt_params := libs.OptionsParameters{Tipo: "C", S: 413.98, K: 420, T: t, R: 0.045, Sigma: 0.15, Q: 0.015}
	simul_params := SimulationParameters{Price_increment: .5, End_price: 423.5, Init_price: 406.0}
	simulate(&opt_params, &simul_params)
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
