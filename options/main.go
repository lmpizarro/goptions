package main

import (
	"fmt"
	"lmpizarro/options/libs"
)

func main() {

	fmt.Println("Tests")
	params := libs.Parameters{S: 100.0, K: 100.0,
		Tipo: "P", T: 1, Sigma: .4, Q: 0.01, R: 0.04}
	fmt.Println(libs.Delta(&params))
	fmt.Println(libs.Gamma(&params))

	deltaBin := libs.DeltaBin(&params, 150)
	fmt.Println(deltaBin)

	params.S = 100.0
	gammaBin := libs.GammaBin(&params, 150)
	fmt.Println(gammaBin)

	fmt.Println(params)

}

