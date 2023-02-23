package libs

import (
	"fmt"
	"math"
	"sync"
	"time"

	finance "github.com/piquette/finance-go"
	"github.com/piquette/finance-go/datetime"
	"github.com/piquette/finance-go/options"
	"github.com/piquette/finance-go/quote"
)

func process(exp_dates [][]string, symbol string,
	limits *Limits, wg *sync.WaitGroup) {

	time1 := time.Now().Unix()
	params := OptionParameters{R: 0.04, Q: 0.01, Tipo: "C"}
	var count int
	for _, exp_date := range exp_dates {
		if limits.Exp_date > exp_date[0] {
			s1 := fetch_options(symbol, exp_date[0], limits)
			fmt.Println("exp date", exp_date[0])
			for index, call := range *s1.Mcall {
				for _, e := range call {
					params.S = limits.S0
					params.K = e[0]
					params.T = float64(index) / 365
					_, iv := IvBsBisection_A(&params, e[1])
					fmt.Println(Round_down(iv, 4), e[0], e[1], Round_down(e[2], 4))
					count++
				}
			}
		}
	}
	time2 := time.Now().Unix()
	fmt.Println(count, time2-time1)
	wg.Done()
}

// In: symbol string
// Out: exp_dates [[2023-05-24]...]
func expiration_dates(symbol string) [][]string {
	// fetch options.
	p := &options.Params{
		UnderlyingSymbol: symbol,
	}

	iter := options.GetStraddleP(p)
	meta := iter.Meta()
	if meta == nil {
		panic("could not retrieve dates")
	}

	exp_dates := [][]string{}
	for _, stamp := range meta.AllExpirationDates {
		// set the day to friday instead of EOD thursday..
		// weird math here..
		stamp = stamp + 86400
		t := time.Unix(int64(stamp), 0)
		exp_dates = append(exp_dates, []string{t.Format("2006-01-02")})
	}
	return exp_dates
}

type CallPut struct {
	Mcall *map[int][][]float64
	Mput  *map[int][][]float64
}

func Round_down(num float64, n float64) float64 {
	return math.Floor(num*math.Pow(10, n)) / (math.Pow(10, n))
}

type Limits struct {
	S0       float64
	K_min    float64
	K_max    float64
	Exp_date string
}

func get_straddles(exp_date time.Time, symbol string) []*finance.Straddle {
	params := &options.Params{
		UnderlyingSymbol: symbol,
	}
	params.Expiration = datetime.New(&exp_date)
	iter_straddle := options.GetStraddleP(params)
	straddles := []*finance.Straddle{}

	for iter_straddle.Next() {
		straddles = append(straddles, iter_straddle.Straddle())
	}
	if iter_straddle.Err() != nil {
		panic("error iter")
	}
	return straddles
}


func fetch_options(symbol, exp_date_str string, limits *Limits) CallPut {
	m_call := make(map[int][][]float64)
	m_put := make(map[int][][]float64)
	var call_put CallPut

	exp_date_time, in_seconds := ParserStringDate(exp_date_str)
	straddles := get_straddles(exp_date_time, symbol)
	ttm_days := TtmInDays(in_seconds)
	for _, straddle := range straddles {
		call := straddle.Call
		put := straddle.Put
		if call != nil && put != nil && ttm_days > 0 {
			strike := call.Strike
			if strike < limits.K_max && strike > limits.K_min {
				arr := []float64{strike, call.LastPrice, call.ImpliedVolatility}
				m_call[int(ttm_days)] = append(m_call[int(ttm_days)], arr)
				arr = []float64{strike, put.LastPrice, put.ImpliedVolatility}
				m_put[int(ttm_days)] = append(m_put[int(ttm_days)], arr)
			}
		}
	}
	call_put.Mcall = &m_call
	call_put.Mput = &m_put
	return call_put
}

func Parallel_Calc_IV(symbol string) {

	q, err := quote.Get(symbol)
	if err != nil {
		panic(err)
	}

	expirations := expiration_dates(symbol)
	n := len(expirations)
	r1 := int(math.Round(float64(n)/4.0)) * 1
	r2 := int(math.Round(float64(n)/4.0)) * 2
	r3 := int(math.Round(float64(n)/4.0)) * 3
	if r3 > n {
		r3--
	}
	fmt.Println(n, r1, r2, r3)
	S0 := q.RegularMarketPrice
	limits := Limits{S0: S0, K_min: S0 * .2, K_max: S0 * 1.5, Exp_date: "2023-04-30"}

	var wg sync.WaitGroup
	time1 := time.Now().Unix()
	wg.Add(4)
	go process(expirations[0:r1], symbol, &limits, &wg)
	go process(expirations[r1:r2], symbol, &limits, &wg)
	go process(expirations[r2:r3], symbol, &limits, &wg)
	go process(expirations[r3:n], symbol, &limits, &wg)
	wg.Wait()
	time2 := time.Now().Unix()
	fmt.Println(time2 - time1)
}
