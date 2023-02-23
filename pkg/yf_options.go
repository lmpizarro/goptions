package libs

import (
	"fmt"

	"github.com/piquette/finance-go"
	"github.com/piquette/finance-go/quote"
	"github.com/sajari/regression"
)

type YfOption interface {
	SetPutMoneynessFactor(ftr float64)
	SetType()
	SetMinMoneyness()
	SetMaxMoneyness()
	SetMinMaturity()
	SetMaxExpDate()
	SetKmax()
	SetKmin()
	SetS0()
	SetRegularMarketPrice()
	Set_Symbol()
	Set_Max_price()
}

type YfParams struct {
	S0                 float64
	Kmin               float64
	Kmax               float64
	MaxExpDate         string
	Symbol             string
	ExpDate            string
	MinMoneyness       float64
	MaxMoneyness       float64
	MinMaturityInDays  int64
	MaxPrice           float64
	PutMoneynessFactor float64
	Type               string
}

func (p *YfParams) SetPutMoneynessFactor(ftr float64, v bool) {
	p.PutMoneynessFactor = ftr
	if v {
		fmt.Println(set_literal_float("Put moneyness factor: ", p.PutMoneynessFactor))
	}
}

func (p *YfParams) SetType(tpe string, v bool) {
	p.Type = tpe
	if v {
		fmt.Println(set_literal_string("Option type", p.Type))
	}
}

func (p *YfParams) SetMinMoneyness(mnn float64, v bool) {
	p.MinMoneyness = mnn
	if v{
		fmt.Println(set_literal_float("Min moneyness", float64(p.MinMoneyness)))
	}
}

func (p *YfParams) SetMaxMoneyness(mnn float64, v bool) {
	p.MaxMoneyness = mnn
	if v{
		fmt.Println(set_literal_float("Max moneyness", float64(p.MaxMoneyness)))
	}
}

func (p *YfParams) SetMinMaturity(mat int64, v bool) {
	p.MinMaturityInDays = mat
	if v {
		fmt.Println(set_literal_float("Min maturity", float64(p.MinMaturityInDays)))
	}
}

func set_literal_float(literal string, value float64) string {
	return fmt.Sprintf("Set %v to % v", literal, Round_down(value, 4))
}

func set_literal_string(literal string, value string) string {
	return fmt.Sprintf("Set %v to % v", literal, value)
}

func (p *YfParams) SetMaxPrice(mp float64, v bool) {
	p.MaxPrice = mp
	if v {
		fmt.Println(set_literal_float("Max price", p.MaxPrice))
	}
}

func (params *YfParams) SetMaxExpDate(date string, v bool) {
	params.MaxExpDate = date
	if v{
		fmt.Println(set_literal_string("Max expiration date", params.MaxExpDate))
	}
}

// Assign the underlying price to S0
func (params *YfParams) SetRegularMarketPrice(v bool) {
	q, err := quote.Get(params.Symbol)
	if err != nil {
		panic(err)
	}
	params.S0 = q.RegularMarketPrice
	if v{
		fmt.Println(set_literal_float("S0", params.S0))
	}
}

func (params *YfParams) SetS0(s0 float64) {
	params.S0 = s0
}

func (params *YfParams) SetKmin(pct float64) {
	params.Kmin = params.S0 * pct / 100
	fmt.Println(set_literal_float("K min", params.Kmin))
}

func (params *YfParams) SetKmax(pct float64) {
	params.Kmax = params.S0 * pct / 100

	fmt.Println(set_literal_float("K max", params.Kmax))
}

func (params *YfParams) SetSymbol(symbol string, v bool) {
	params.Symbol = symbol
	if v{
		fmt.Println(set_literal_string("Symbol", params.Symbol))

	}
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

func FetchOptions(params *YfParams) []*finance.Straddle {
	var filtered_straddles []*finance.Straddle

	exp_date_time, _ := ParserStringDate(params.ExpDate)
	straddles := get_straddles(exp_date_time, params.Symbol)
	_, max_ttm_seconds := ParserStringDate(params.MaxExpDate)
	for _, straddle := range straddles {
		non_nil_condition := straddle.Call != nil && straddle.Put != nil

		if non_nil_condition && (int64(straddle.Call.Expiration) < max_ttm_seconds) {
			filtered_straddles = append(filtered_straddles, straddle)
		}
	}
	return filtered_straddles
}

func getLineOut(params *OptionParameters, str *finance.Straddle, money_ness float64) [9]float64 {

	delta := params.Delta()
	gamma := params.Gamma()
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

func getOutput(params *OptionParameters, straddle *finance.Straddle, money_ness float64) string {

	line := getLineOut(params, straddle, money_ness)

	formatD := "%6s %6.2f %6.2f %10.4f %10.4f %10.4f %10.4f %10.4f %10.4f"
	formatD = "H%2s S %6.2f K %6.2f P %6.4f T %6.4f V %6.4f M %6.4f D %6.4f G %6.4f"
	return fmt.Sprintf(formatD, params.Tipo, line[1], line[2], line[3],
		line[4], Round_down(line[5], 4),
		Round_down(line[6], 4), Round_down(line[7], 4), Round_down(line[8], 4))
}

// Print header of output
func printHeader() {
	formatH := "%6s %6s %6s %10s %10s %10s %10s %10s %10s\n"
	fmt.Printf(formatH, "tipo", "S0", "K", "Price", "Matur", "IV", "Mnn", "delta", "gamma")
}

// Calculate the moneyness of a straddle
func moneyness(yf_params *YfParams, str *finance.Straddle) (float64,
	float64) {
	mnnC := (yf_params.S0 - str.Call.Strike) / str.Call.Strike
	mnnP := (str.Put.Strike - yf_params.S0) / str.Put.Strike
	return mnnC, mnnP
}

func call_put_filter_01(yf_params *YfParams, mnnC float64,
	straddle *finance.Straddle, ttm_days int64) bool {
	var factor float64
	var last_price float64

	if yf_params.Type == "C" {
		factor = 1.0
		last_price = straddle.Call.LastPrice
	} else {
		factor = yf_params.PutMoneynessFactor
		last_price = straddle.Put.LastPrice
	}
	return mnnC < yf_params.MaxMoneyness &&
		mnnC > factor*yf_params.MinMoneyness &&
		last_price < yf_params.MaxPrice &&
		ttm_days > yf_params.MinMaturityInDays
}

func call_put_filter_02(yf_params *YfParams, mnnC float64,
	straddle *finance.Straddle,
	ttm_days int64) bool {

	relation_01 := mnnC > yf_params.MinMoneyness
	relation_02 := mnnC < yf_params.MaxMoneyness
	relation_03 := ttm_days >= yf_params.MinMaturityInDays

	return relation_01 && relation_02 && relation_03

}

func Yf_Options(yf_params *YfParams, print bool) (c [][9]float64, p [][9]float64) {
	exp_dates := expiration_dates(yf_params.Symbol)
	exp_dates = limit_exp_dates(exp_dates, yf_params.MaxExpDate)

	var puts [][9]float64
	var calls [][9]float64

	if print {
		printHeader()
	}
	for _, exp_date := range exp_dates {
		yf_params.ExpDate = exp_date[0]
		straddles := FetchOptions(yf_params)
		for _, straddle := range straddles {
			mnnC, mnnP := moneyness(yf_params, straddle)
			ttm_days := TtmInDays(int64(straddle.Put.Expiration))

			(yf_params).SetType("C", false)
			if call_put_filter_02(yf_params, mnnC, straddle, ttm_days) {
				par_calc := OptionParameters{Tipo: "C", S: yf_params.S0, K: straddle.Call.Strike,
					T: float64(ttm_days) / 365.0, R: 0.045, Sigma: straddle.Call.ImpliedVolatility,
					Q: 0.02}
				if print {
					fmt.Println(getOutput(&par_calc, straddle, mnnC))
				}
				calls = append(calls, getLineOut(&par_calc, straddle, mnnC))
			}
			(yf_params).SetType("P", false)
			if call_put_filter_02(yf_params, mnnP, straddle, ttm_days) {
				par_calc := OptionParameters{Tipo: "P", S: yf_params.S0, K: straddle.Put.Strike,
					T: float64(ttm_days) / 365.0, R: 0.045, Sigma: straddle.Put.ImpliedVolatility,
					Q: 0.02}
				if print {
					fmt.Println(getOutput(&par_calc, straddle, mnnC))
				}
				puts = append(puts, getLineOut(&par_calc, straddle, mnnP))
			}

		}
	}
	return calls, puts
}

func MakeRegression(points [][9]float64, observer string, description string) {
	// See https://github.com/sajari/regression
	// "K ", e[2], "T (days) ", e[4],  "Price ", e[3], "IV ", e[5]
	var observed float64
	var counter int
	var mean_observed float64

	r := new(regression.Regression)
	if observer == "Price" {
		r.SetObserved("Price")
	} else if observer == "IV" {
		r.SetObserved("IV")
	} else if observer == "Delta" {
		r.SetObserved("Delta")
	} else if observer == "Gamma" {
		r.SetObserved("Gamma")
	}

	r.SetVar(0, "K")
	r.SetVar(1, "T")
	r.SetVar(2, "TT")
	r.SetVar(3, "KK")
	r.SetVar(4, "KT")
	mean_observed = 0
	for i, point := range points {
		if observer == "Price" {
			observed = point[3]
		} else if observer == "IV" {
			observed = point[5]
		} else if observer == "Delta" {
			observed = point[7]
		} else if observer == "Gamma" {
			observed = point[8]
		}
		r.Train(regression.DataPoint(observed, []float64{point[2], point[4], point[4] * point[4], point[2] * point[2], point[2] * point[4]}))
		counter = i
		mean_observed += observed
	}
	r.Run()

	fmt.Printf("# %v mean %v for %v %v \n", counter, mean_observed/float64(counter), description, r.Formula)
	fmt.Printf("R2: %.2e Var Pred %.2e Var obs %.2e \n", r.R2, r.VariancePredicted, r.Varianceobserved)

}

func MeanIV(points [][9]float64) float64 {
	// See https://github.com/sajari/regression
	// "K ", e[2], "T (days) ", e[4],  "Price ", e[3], "IV ", e[5]
	var counter int
	var mean float64

	mean = 0
	for i, point := range points {
		counter = i
		mean += point[5]
	}

	return mean / float64(counter)

}

func RegularMarketPrice(symbol string) float64 {
	q, err := quote.Get(symbol)
	if err != nil {
		panic(err)
	}
	return q.RegularMarketPrice
}
