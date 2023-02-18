package libs

import (

	"encoding/json"
)

const Day = 1.0 / 365.0

type SimulationParameters struct {
	End_price float64
	Init_price float64
	Price_increment float64
}

func Default_simulate_parameters(options_params *OptionsParameters) *SimulationParameters{
	var simul_params SimulationParameters

	const percent = 0.015
	simul_params.End_price = float64(int(options_params.K * (1 + percent)))
	simul_params.Init_price = float64(int(options_params.K * (1 - percent)))
	simul_params.Price_increment = .5
	return &simul_params
}


func Simulate_long(option_params *OptionsParameters,
			  simul_params *SimulationParameters, default_params bool) [][]float64{
	var sim_params *SimulationParameters


	if default_params {
		sim_params = Default_simulate_parameters(option_params)
	} else {
		sim_params = simul_params
	}
	day := Day
	t := option_params.T
	pinit := sim_params.Init_price // ex 406.0
	pfinal := sim_params.End_price // ex 423.5
	delta_precio := sim_params.Price_increment // ex .5
	cost_init := Bs(option_params)

	var price_of_option float64
	var values []float64
	var rows [][]float64
	var price_of_equity float64
	price_of_equity = pinit
	// See https://www.optionsprofitcalculator.com/calculator/long-call.html
	// See https://optionstrat.com/
	for {
		values = append(values, Round_down(365 *t , 1))
		for {
			option_params.T = t
			option_params.S = price_of_equity
			price_of_option = Bs(option_params)
			// values = append(values, libs.Round_down(100*(c - cinit)/cinit, 2))
			values = append(values, Round_down(price_of_option-cost_init, 2))
			if price_of_equity > pfinal {
				break
			}
			price_of_equity = price_of_equity + delta_precio
		}
		rows = append(rows, values)
		values = values[:0]
		price_of_equity = pinit
		t = t - day
		if t < day {
			break
		}
	}
	return rows
}

type ResultSimulation struct {
	Simulation []struct {
		Day    int       `json:"day"`
		Prices []float64 `json:"prices"`
	} `json:"rows"`
}

func Rows_simulation_to_json(rows [][]float64) []byte{
	var results ResultSimulation
	for _, row := range rows {
		sim_result := struct{
			Day int
			Prices []float64
		}{Day: int(row[0]), Prices: row[1:]}
		results.Simulation = append(results.Simulation,
			struct{Day int "json:\"day\""; Prices []float64 "json:\"prices\""}(sim_result))
	}
	u, _ := json.Marshal(results)
	return u
}

/*
	https://blog.boot.dev/golang/anonymous-structs-golang/
	{
	"simulation": [{"day":1, "prices": [1.5,2.5,3.5,4.5]},
			{"day":2, "prices": [1.5,2.5,3.5,4.5]}]
			}

		sim_result.Day = 1
		sim_result.Prices = values
		results.Simulation = append(results.Simulation, struct{Day int "json:\"day\""; Prices []float64 "json:\"prices\""}(sim_result))
		sim_result := struct{
			Day int
			Prices []float64
		}{}
}
*/