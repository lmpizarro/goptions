package main

import (
	"fmt"
	"lmpizarro/options/libs"
	"lmpizarro/options/rfx"
)

func main() {

	cred := rfx.ReadCredentials("./env.csv", false)
	user := cred.User
	password := cred.Password
	token := rfx.Token(user, password)

	future_spy := libs.Future{Symbol: "^GSPC", Maturity: "2023-03-23", Futu: "ESH23.CME"}
	future_ggal := libs.Future{Symbol: "GGAL.BA", Maturity: "2023-02-28", Futu: "GGAL/FEB23"}

	implied_rate, absPct := libs.ImpliedRate(&future_spy)
	fmt.Printf("PC SPY %.2f Implied %.2f\n", 100*absPct, 100*implied_rate)

	symbol := libs.Symbol(future_ggal.Symbol)
	spot := symbol.Price()
	years_to_mat := libs.YearsToMat(future_ggal.Maturity)
	fut, _ := rfx.LastPrice(future_ggal.Futu, token)
	iR, pCt := libs.Rates(fut, spot, years_to_mat)
	fmt.Printf("PC GGAL %.2f Implied %.2f\n", 100*pCt, 100*iR)
	libs.Test_YF()
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
