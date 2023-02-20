package libs

import (
	"fmt"

	"github.com/piquette/finance-go"
	"github.com/piquette/finance-go/quote"
)

type Yf_params struct {
	S0                   float64
	K_min                float64
	K_max                float64
	Max_exp_date         string
	Symbol               string
	Exp_date             string
	Min_moneyness        float64
	Max_moneyness        float64
	Min_maturity_in_days int64
	Max_price            float64
	Put_moneyness_factor float64
	Type                 string
}

func (p *Yf_params) Set_Put_moneyness_factor(ftr float64) {
	p.Put_moneyness_factor = ftr
	fmt.Println(set_literal_float("Put moneyness factor: ", p.Put_moneyness_factor))
}

func (p *Yf_params) Set_Type(tpe string, v bool) {
	p.Type = tpe
	if v {
		fmt.Println(set_literal_string("Option type", p.Type))
	}
}

func (p *Yf_params) Set_Min_moneyness(mnn float64) {
	p.Min_moneyness = mnn
	fmt.Println(set_literal_float("Min moneyness", float64(p.Min_moneyness)))
}

func (p *Yf_params) Set_Max_moneyness(mnn float64) {
	p.Max_moneyness = mnn
	fmt.Println(set_literal_float("Max moneyness", float64(p.Max_moneyness)))
}

func (p *Yf_params) Set_Min_maturity(mat int64) {
	p.Min_maturity_in_days = mat
	fmt.Println(set_literal_float("Min maturity", float64(p.Min_maturity_in_days)))
}

func set_literal_float(literal string, value float64)string{
	return fmt.Sprintf("Set %v to % v", literal, Round_down(value,4))
}

func set_literal_string(literal string, value string)string{
	return fmt.Sprintf("Set %v to % v", literal, value)
}


func (p *Yf_params) Set_Max_price(mp float64) {
	p.Max_price = mp
	fmt.Println(set_literal_float("Max price", p.Max_price))
}

func (params *Yf_params) Set_Max_Exp_date(date string) {
	params.Max_exp_date = date
	fmt.Println(set_literal_string("Max expiration date", params.Max_exp_date))
}

func (params *Yf_params) Set_S0_Market_Price() {
	q, err := quote.Get(params.Symbol)
	if err != nil {
		panic(err)
	}
	params.S0 = q.RegularMarketPrice
	fmt.Println(set_literal_float("S0", params.S0))
}

func (params *Yf_params) Set_S0(s0 float64) {
	params.S0 = s0
}

func (params *Yf_params) Set_K_min(pct float64) {
	params.K_min = params.S0 * pct / 100
	fmt.Println(set_literal_float("K min", params.K_min))
}

func (params *Yf_params) Set_K_max(pct float64) {
	params.K_max = params.S0 * pct / 100

	fmt.Println(set_literal_float("K max", params.K_max))
}

func (params *Yf_params) Set_Symbol(symbol string) {
	params.Symbol = symbol
	fmt.Println(set_literal_string("Symbol", params.Symbol))
}

func limit_exp_dates(exp_dates [][]string, limit string) [][]string {
	var limited [][]string
	for _, exp_date := range exp_dates {
		if exp_date[0] < limit {
			limited = append(limited, exp_date)
		}
	}
	return limited
}

func Fetch_Options(params *Yf_params) []*finance.Straddle {
	var filtered_straddles []*finance.Straddle

	exp_date_time, _ := parse_string_date(params.Exp_date)
	straddles := get_straddles(exp_date_time, params.Symbol)
	_, max_ttm_seconds := parse_string_date(params.Max_exp_date)
	for _, straddle := range straddles {
		non_nil_condition := straddle.Call != nil && straddle.Put != nil

		if non_nil_condition && (int64(straddle.Call.Expiration) < max_ttm_seconds) {
			filtered_straddles = append(filtered_straddles, straddle)
		}
	}
	return filtered_straddles
}

func get_line_out(params *OptionsParameters, str *finance.Straddle, money_ness float64) [9]float64 {

	delta := Delta(params)
	gamma := Gamma(params)
	var line [9]float64
	var last_price float64
	if params.Tipo == "C" {
		last_price = str.Call.LastPrice
		line[0] = 1
	} else {
		last_price = str.Put.LastPrice
		line[0] = -1
	}

	line[1] = params.S
	line[2] = params.K
	line[3] = last_price
	line[4] = 365 * params.T
	line[5] = params.Sigma
	line[6] = money_ness
	line[7] = delta
	line[8] = gamma

	return line
}

func get_output(params *OptionsParameters, str *finance.Straddle, money_ness float64) string {

	line := get_line_out(params, str, money_ness)

	formatD := "%6s %6.2f %6.2f %10.4f %10.4f %10.4f %10.4f %10.4f %10.4f"
	formatD = "H%2s S %6.2f K %6.2f P %6.4f T %6.4f V %6.4f M %6.4f D %6.4f G %6.4f"
	return fmt.Sprintf(formatD, params.Tipo, line[1], line[2], line[3],
		line[4], Round_down(line[5], 4),
		Round_down(line[6], 4), Round_down(line[7], 4), Round_down(line[8], 4))
}

func get_header() {
	formatH := "%6s %6s %6s %10s %10s %10s %10s %10s %10s\n"
	fmt.Printf(formatH, "tipo", "S0", "K", "Price", "Matur", "IV", "Mnn", "delta", "gamma")
}

func money_ness(yf_params *Yf_params, str *finance.Straddle) (float64, bool,
	float64, bool) {
	mnnC := (yf_params.S0 - str.Call.Strike) / str.Call.Strike
	mnnP := (str.Put.Strike - yf_params.S0) / str.Put.Strike
	mnnPBool := yf_params.S0 < str.Put.Strike
	mnnCBool := yf_params.S0 > str.Put.Strike

	return mnnC, mnnCBool, mnnP, mnnPBool
}

func call_put_filter_01(yf_params *Yf_params, mnnC, last_price float64, ttm_days int64) bool {
	var factor float64
	if yf_params.Type == "C" {
		factor = 1.0
	} else {
		factor = yf_params.Put_moneyness_factor
	}
	return mnnC < yf_params.Max_moneyness &&
		mnnC > factor*yf_params.Min_moneyness &&
		last_price < yf_params.Max_price &&
		ttm_days > yf_params.Min_maturity_in_days
}

func Yf_Options(yf_params *Yf_params) {
	exp_dates := expiration_dates(yf_params.Symbol)
	exp_dates = limit_exp_dates(exp_dates, yf_params.Max_exp_date)

	get_header()
	for _, exp_date := range exp_dates {
		yf_params.Exp_date = exp_date[0]
		straddles := Fetch_Options(yf_params)
		for _, straddle := range straddles {
			mnnC, mnnCBool, mnnP, mnnPBool := money_ness(yf_params, straddle)
			ttm_days := ttm_in_days(int64(straddle.Put.Expiration))
			if !mnnCBool {
				if call_put_filter_01(yf_params, mnnC, straddle.Call.LastPrice, ttm_days) {
					par_calc := OptionsParameters{Tipo: "C", S: yf_params.S0, K: straddle.Call.Strike,
						T: float64(ttm_days) / 365.0, R: 0.045, Sigma: straddle.Call.ImpliedVolatility,
						Q: 0.02}
					fmt.Println(get_output(&par_calc, straddle, mnnC))
				}
			}
			(yf_params).Set_Type("P", false)
			if !mnnPBool {
				if call_put_filter_01(yf_params, mnnP, straddle.Put.LastPrice, ttm_days) {
					par_calc := OptionsParameters{Tipo: "P", S: yf_params.S0, K: straddle.Put.Strike,
						T: float64(ttm_days) / 365.0, R: 0.045, Sigma: straddle.Put.ImpliedVolatility,
						Q: 0.02}
					fmt.Println(get_output(&par_calc, straddle, mnnP))
				}
			}
		}
	}
}

