package libs

import (
	"fmt"

	"github.com/piquette/finance-go"
	"github.com/piquette/finance-go/quote"
)

type Yf_params struct {
	S0       float64
	K_min    float64
	K_max    float64
	Max_exp_date string
	Symbol string
	Exp_date string
}

func Set_Max_Exp_date(params *Yf_params, date string){
	params.Max_exp_date = date
}

func Set_S0_Market_Price(params *Yf_params){
	q, err := quote.Get(params.Symbol)
	if err != nil {
		panic(err)
	}
	params.S0 = q.RegularMarketPrice
}

func Set_S0(params *Yf_params, s0 float64){
	params.S0 = s0
}

func Set_K_min(params *Yf_params, pct float64){
	params.K_min = params.S0 * pct / 100
}

func Set_K_max(params *Yf_params, pct float64){
	params.K_max = params.S0 * pct / 100
}

func Set_Symbol(params *Yf_params, symbol string){
	params.Symbol = symbol
}

func limit_exp_dates(exp_dates [][]string, limit string)[][]string{
	var limited [][]string
	for _, exp_date := range exp_dates {
		if exp_date[0] < limit {
			limited = append(limited, exp_date)
		}
	}
	return limited
}

func Fetch_Options(params *Yf_params) []*finance.Straddle{
	var filtered_straddles []*finance.Straddle

	exp_date_time, _ := parse_string_date(params.Exp_date)
	straddles := get_straddles(exp_date_time, params.Symbol)
	_, max_ttm_seconds := parse_string_date(params.Max_exp_date)
	for _, straddle := range straddles {
		non_nil_condition := straddle.Call != nil && straddle.Put != nil

		if  non_nil_condition && (int64(straddle.Call.Expiration) <  max_ttm_seconds) {
			filtered_straddles = append(filtered_straddles, straddle)
		}
	}
	return filtered_straddles
}

func Yf_Options(params *Yf_params) {
	exp_dates := expiration_dates(params.Symbol)
	exp_dates = limit_exp_dates(exp_dates, params.Max_exp_date)

	formatD := "%6s %6.2f %6.2f %10.4f %10.4f %10.4f %10.4f %10.4f %10.4f\n"
	formatH := "%6s %6s %6s %10s %10s %10s %10s %10s %10s\n"
	fmt.Printf(formatH, "tipo", "S0",  "K", "Price",  "Matur", "IV", "Mnn", "delta", "gamma")
	for _, exp_date := range exp_dates {
		params.Exp_date = exp_date[0]
		strdls := Fetch_Options(params)
		for _, str := range strdls {
			ttm_days := ttm_in_days(int64(str.Put.Expiration))
			mnnC := (params.S0 - str.Call.Strike) / str.Call.Strike
			mnnP := (str.Put.Strike - params.S0) / str.Put.Strike
			mnnPBool := params.S0 < str.Put.Strike
			mnnCBool := params.S0 > str.Put.Strike
			price_limit := 0.0024 * params.S0
			min_mnn := -0.05
			max_mnn := -0.01
			min_ttm := 10
			if !mnnCBool {
				if mnnC < max_mnn && mnnC > min_mnn && str.Call.LastPrice < price_limit  && ttm_days > int64(min_ttm){
					par_calc := Parameters{Tipo: "C", S: params.S0, K: str.Call.Strike,
						T: float64(ttm_days) / 365.0, R: 0.045, Sigma: str.Call.ImpliedVolatility, Q: 0.02}
					delta := Delta(&par_calc)
					gamma := Gamma(&par_calc)

					fmt.Printf(formatD, "C ", params.S0, str.Call.Strike, str.Call.LastPrice,
						float64(ttm_days/7.0), round_down(str.Call.ImpliedVolatility, 4),
						round_down(mnnC, 4), round_down(delta, 4), round_down(gamma, 4))
				}
			}
			if !mnnPBool{
				if mnnP < max_mnn && mnnP > 1.5 * min_mnn && str.Put.LastPrice < price_limit && ttm_days > int64(min_ttm){

					par_calc := Parameters{Tipo: "P", S: params.S0, K: str.Put.Strike,
						T: float64(ttm_days) / 365.0, R: 0.045, Sigma: str.Put.ImpliedVolatility, Q: 0.02}
					delta := Delta(&par_calc)
					gamma := Gamma(&par_calc)
					fmt.Printf(formatD, "P ", params.S0, str.Put.Strike, str.Put.LastPrice,
						float64(ttm_days)/7.0, round_down(str.Put.ImpliedVolatility, 4),
						round_down(mnnP, 4), round_down(delta, 4), round_down(gamma, 4))

				}



			}
		}
	}
	fmt.Println("", params)



}

func Test_YF(){
	var params Yf_params

	Set_Symbol(&params, "SPY")
	Set_S0_Market_Price(&params)
	Set_K_max(&params, 180)
	Set_K_min(&params, 20)
	Set_Max_Exp_date(&params, "2023-04-30")

	Yf_Options(&params)

}