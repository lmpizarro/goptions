package main

import (
	"fmt"
	"lmpizarro/options/libs"
)

func main() {

	params := libs.Parameters{S: 100.0, K: 100.0,
		Tipo: "P", T: 1, Sigma: .4, Q: 0.01, R: 0.04}

	PBs := libs.Bs(&params)
	deltaPBs := libs.Delta(&params)
	fmt.Println("hi P", PBs, deltaPBs)

	params.Tipo = "C"
	CBs := libs.Bs(&params)
	deltaCBs := libs.Delta(&params)

	fmt.Println("hi C ", CBs, deltaCBs)

	fmt.Println("diff deltas ", deltaCBs-deltaPBs)

	C := libs.Bin(&params, 150)
	fmt.Println("hi C Bin", C)
	params.Tipo = "P"
	P := libs.Bin(&params, 150)
	fmt.Println("hi P Bin", P)

	gamma := libs.Gamma(&params)
	fmt.Println("gamma ", gamma)
	vega := libs.Vega(&params)
	fmt.Println("vega ", vega)

	params.Tipo = "P"
	IV := libs.IvBs(&params, PBs)
	fmt.Println("IV ", IV)

	IV = libs.IvBsNewton(&params, 0.1, PBs)
	fmt.Println("IV Newton ", IV)

	params.Tipo = "C"
	thetaC := libs.Theta(&params, true)
	params.Tipo = "P"
	thetaP := libs.Theta(&params, true)

	fmt.Println(thetaC, thetaP)

	params.Tipo = "C"
	rhoC := libs.Rho(&params)
	params.Tipo = "P"
	rhoP := libs.Rho(&params)

	fmt.Println(rhoC, rhoP)

	params.K = 50
	params.S = 49
	params.T = 0.3846
	params.R = 0.05
	params.Sigma = .2
	params.Q = 0.0
	params.Tipo = "C"
	gamma = libs.Gamma(&params)
	thetaC = libs.Theta(&params, true)
	deltaCBs = libs.Delta(&params)
	rhoC = libs.Rho(&params)

	fmt.Println("\t Options, Futures, Derivatives 9th ed, J.C. Hull")
	fmt.Println("gamma Hull pag 415", gamma)
	fmt.Println("theta Hull pag 409", thetaC)
	fmt.Println("delta Hull pag 428", deltaCBs)
	fmt.Println("rho Hull pag 440", rhoC)

}

