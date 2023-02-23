package libs

import (
	"fmt"
)

func Test() {
	params := OptionParameters{S: 100.0, K: 100.0,
		Tipo: "P", T: 1, Sigma: .4, Q: 0.01, R: 0.04}

	PBs := params.Bs()
	deltaPBs := params.Delta()
	fmt.Println("hi P", PBs, deltaPBs)

	params.Tipo = "C"
	CBs := params.Bs()
	deltaCBs := params.Delta()

	fmt.Println("hi C ", CBs, deltaCBs)

	fmt.Println("diff deltas ", deltaCBs-deltaPBs)

	C := Bin(&params, 150)
	fmt.Println("hi C Bin", C)
	params.Tipo = "P"
	P := Bin(&params, 150)
	fmt.Println("hi P Bin", P)

	gamma := params.Gamma()
	fmt.Println("gamma ", gamma)
	vega := params.Vega()
	fmt.Println("vega ", vega)

	params.Tipo = "P"
	s, IV := IvBsBisection_A(&params, PBs)
	fmt.Println("IV ", IV, "Steps", s)

	_, IV = IvBsNewton(&params, 0.1, PBs, 0.001)
	fmt.Println("IV Newton ", IV)

	params.Tipo = "C"
	thetaC := params.Theta(true)
	params.Tipo = "P"
	thetaP := params.Theta(true)

	fmt.Println(thetaC, thetaP)

	params.Tipo = "C"
	rhoC := params.Rho()
	params.Tipo = "P"
	rhoP := params.Rho()

	fmt.Println(rhoC, rhoP)

	params.K = 50
	params.S = 49
	params.T = 0.3846
	params.R = 0.05
	params.Sigma = .2
	params.Q = 0.0
	params.Tipo = "C"

	gamma = params.Gamma()
	thetaC = params.Theta(true)
	deltaCBs = params.Delta()
	rhoC = params.Rho()

	fmt.Println("\t Options, Futures, Derivatives 9th ed, J.C. Hull")
	fmt.Println("gamma Hull pag 415", gamma)
	fmt.Println("theta Hull pag 409", thetaC)
	fmt.Println("delta Hull pag 428", deltaCBs)
	fmt.Println("rho Hull pag 440", rhoC)

}

func TestNewton() (int, float64) {
	t := 110.0 / 365.0

	opt_params := OptionParameters{Tipo: "C", S: 200, K: 180, T: t, R: 0.045, Sigma: 0.6, Q: 0.015}
	c := opt_params.Bs()

	sigma0 := IV_Brenner_Subrahmanyam(&opt_params, c)
	fmt.Println("sigma0 BS", sigma0)
	sigma0, _ = IV_Corrado_Miller(&opt_params, c)
	fmt.Println("sigma0 CM", sigma0)
	sigma0 = IV_Bharadia_Christofides_Salkin(&opt_params, c)
	fmt.Println("sigma0 CS", sigma0)
	s, sigma := IvBsNewton(&opt_params, sigma0, c, 0.000001)
	// s, sigma := IvBsSecant(&opt_params, c)
	return s, sigma
}

func Test_YF() {
	var yf_params YfParams

	(&yf_params).SetSymbol("SPY", false)
	(&yf_params).SetRegularMarketPrice(false)
	(&yf_params).SetMaxExpDate("2023-05-30", false)
	(&yf_params).SetMinMoneyness(-0.25, false)                 // -0.005     -0.045
	(&yf_params).SetMaxMoneyness(0.25, false)                  //  0.005  -0.000001
	(&yf_params).SetMinMaturity(7, false)                      // 7          1
	(&yf_params).SetMaxPrice(2 * 0.0024 * yf_params.S0, false) //2  1
	(&yf_params).SetPutMoneynessFactor(1.5, false)
	(&yf_params).SetType("C", false)

	calls, puts := Yf_Options(&yf_params, false)

	// MakeRegression(calls, "IV", "calls")
	// MakeRegression(puts, "IV", "puts")
	fmt.Println("IV calls ", MeanIV(calls))
	fmt.Println("IV puts ", MeanIV(puts))

}
