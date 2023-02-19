package libs

import (
	"fmt"
)

func Test() {
	params := OptionsParameters{S: 100.0, K: 100.0,
		Tipo: "P", T: 1, Sigma: .4, Q: 0.01, R: 0.04}

	PBs := Bs(&params)
	deltaPBs := Delta(&params)
	fmt.Println("hi P", PBs, deltaPBs)

	params.Tipo = "C"
	CBs := Bs(&params)
	deltaCBs := Delta(&params)

	fmt.Println("hi C ", CBs, deltaCBs)

	fmt.Println("diff deltas ", deltaCBs-deltaPBs)

	C := Bin(&params, 150)
	fmt.Println("hi C Bin", C)
	params.Tipo = "P"
	P := Bin(&params, 150)
	fmt.Println("hi P Bin", P)

	gamma := Gamma(&params)
	fmt.Println("gamma ", gamma)
	vega := Vega(&params)
	fmt.Println("vega ", vega)

	params.Tipo = "P"
	s, IV := IvBsBisection_A(&params, PBs)
	fmt.Println("IV ", IV, "Steps", s)

	_, IV = IvBsNewton(&params, 0.1, PBs, 0.001)
	fmt.Println("IV Newton ", IV)

	params.Tipo = "C"
	thetaC := Theta(&params, true)
	params.Tipo = "P"
	thetaP := Theta(&params, true)

	fmt.Println(thetaC, thetaP)

	params.Tipo = "C"
	rhoC := Rho(&params)
	params.Tipo = "P"
	rhoP := Rho(&params)

	fmt.Println(rhoC, rhoP)

	params.K = 50
	params.S = 49
	params.T = 0.3846
	params.R = 0.05
	params.Sigma = .2
	params.Q = 0.0
	params.Tipo = "C"

	gamma = Gamma(&params)
	thetaC = Theta(&params, true)
	deltaCBs = Delta(&params)
	rhoC = Rho(&params)

	fmt.Println("\t Options, Futures, Derivatives 9th ed, J.C. Hull")
	fmt.Println("gamma Hull pag 415", gamma)
	fmt.Println("theta Hull pag 409", thetaC)
	fmt.Println("delta Hull pag 428", deltaCBs)
	fmt.Println("rho Hull pag 440", rhoC)

}


func TestNewton() (int, float64) {
	t := 110.0 / 365.0

	opt_params := OptionsParameters{Tipo: "C", S: 200, K: 200, T: t, R: 0.045, Sigma: 0.6, Q: 0.015}
	c := Bs(&opt_params)

	sigma0 := IV_Brenner_Subrahmanyam(&opt_params, c)
	// sigma0 = 0.6
	fmt.Println(sigma0)
	s, sigma := IvBsNewton(&opt_params, sigma0, c, 0.000001)
	// s, sigma := IvBsSecant(&opt_params, c)
	return s, sigma
}


