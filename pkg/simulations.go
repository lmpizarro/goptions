package libs

import (
	"encoding/json"
)

const Day = 1.0 / 365.0

type SimulationParameters struct {
	End_price       float64 `json:"end_price"`
	Init_price      float64 `json:"init_price"`
	Price_increment float64 `json:"price_increment"`
}

func Default_simulate_parameters(options_params *OptionParameters) SimulationParameters {
	var simul_params SimulationParameters

	const percent = 0.015
	simul_params.End_price = float64(int(options_params.K * (1 + percent)))
	simul_params.Init_price = float64(int(options_params.K * (1 - percent)))
	simul_params.Price_increment = .5
	return simul_params
}

func Simulate_long(option_params *OptionParameters,
	simul_params *SimulationParameters) [][]float64 {

	day := Day
	t := option_params.T
	pinit := simul_params.Init_price             // ex 406.0
	pfinal := simul_params.End_price             // ex 423.5
	delta_precio := simul_params.Price_increment // ex .5
	cost_init := option_params.Bs()

	var price_of_option float64
	var values []float64
	var rows [][]float64
	var price_of_equity float64
	price_of_equity = pinit
	// See https://www.optionsprofitcalculator.com/calculator/long-call.html
	// See https://optionstrat.com/
	for {
		values = append(values, Round_down(365*t, 1))
		for {
			option_params.T = t
			option_params.S = price_of_equity
			price_of_option = option_params.Bs()
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
	Sim_params *SimulationParameters
	Opt_params *OptionParameters
	Simulation []struct {
		Day    int       `json:"day"`
		Prices []float64 `json:"prices"`
	} `json:"rows"`
}

func Rows_simulation_to_json(rows [][]float64,
	simul_params *SimulationParameters,
	opt_params *OptionParameters) []byte {
	var results ResultSimulation
	for _, row := range rows {
		sim_result := struct {
			Day    int
			Prices []float64
		}{Day: int(row[0]), Prices: row[1:]}
		results.Simulation = append(results.Simulation,
			struct {
				Day    int       "json:\"day\""
				Prices []float64 "json:\"prices\""
			}(sim_result))
	}
	results.Sim_params = simul_params
	results.Opt_params = opt_params
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

func Test_simulation() []byte {
	t := 11.0 / 365.0
	opt_params := OptionParameters{Tipo: "C", S: 413.98, K: 420, T: t, R: 0.045, Sigma: 0.15, Q: 0.015}

	simul_params := SimulationParameters{Price_increment: .5, End_price: 423.5, Init_price: 406.0}
	simul_params = Default_simulate_parameters(&opt_params)
	rows := Simulate_long(&opt_params, &simul_params)

	u := Rows_simulation_to_json(rows, &simul_params, &opt_params)

	return u
}
