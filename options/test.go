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
	IV := IvBs(&params, PBs)
	fmt.Println("IV ", IV)

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
