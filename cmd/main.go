package main

import (
	"fmt"
	"lmpizarro/options/libs"
	"lmpizarro/options/rfx"
	"time"
)

func main() {

	for {

		token := rfx.Login()

		future_spy := libs.Future{SymbolSpot: "^GSPC", Maturity: "2023-03-23", SymbolFuture: "ESH23.CME"}

		libs.ImpliedRate(&future_spy)
		fmt.Printf("S SPY A %.2f R %.2f F %.2f S %.2f \n",
			100*future_spy.Rate,
			100*future_spy.YearImpliedRate,
			future_spy.PriceFuture,
			future_spy.PriceSpot)

		future_ggal := libs.Future{SymbolSpot: "GGAL.BA", Maturity: "2023-04-28", SymbolFuture: "GGAL/ABR23"}
		symbol := libs.Symbol(future_ggal.SymbolSpot)
		future_ggal.PriceSpot = symbol.Price()
		future_ggal.TimeToMaturity = libs.YearsToMat(future_ggal.Maturity)
		fut, _ := rfx.LastPrice(future_ggal.SymbolFuture, token)
		future_ggal.PriceFuture = fut
		libs.Rates(&future_ggal)
		fmt.Printf("S GGAL A %.2f R %.2f F %.2f S %.2f\n",
			100*future_ggal.Rate,
			100*future_ggal.YearImpliedRate,
			future_ggal.PriceFuture,
			future_ggal.PriceSpot)

		fmt.Printf("C %.2f\n\n", libs.Ccl())
		libs.Test_YF()
		fmt.Printf("---------------------\n\n")
		time.Sleep(5000 * time.Millisecond)

	}
	panic("main")

	fmt.Println(libs.Ccl())
	fmt.Println(libs.GGALBA())
	fmt.Println(libs.TestNewton())
	//libs.Parallel_Calc_IV("SPY")

	fmt.Println("Tests")
	params := libs.OptionParameters{S: 100.0, K: 100.0,
		Tipo: "P", T: 1, Sigma: .4, Q: 0.01, R: 0.04}
	fmt.Println(params.Delta())
	fmt.Println(params.Gamma())

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
