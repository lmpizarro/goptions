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


func process(exp_dates [][]string, symbol string, price float64, wg *sync.WaitGroup) {
	time1 := time.Now().Unix()
	params := Parameters{R: 0.04, Q: 0.01, Tipo: "C"}
	var count int
	for _, eF := range exp_dates {
		s1 := fetch_options(symbol, eF[0], price)
		for index, arr := range *s1.Mcall {
			for _, e := range arr {
				params.S = price
				params.K = e[0]
				params.T = float64(index) / 365
				IvBs(&params, e[1])
				//fmt.Println(eF[0], round_down(iv, 4), e[0], e[1], round_down(e[2], 4))
				count++
			}
		}
	}
	time2 := time.Now().Unix()
	fmt.Println(count, time2-time1)
	wg.Done()
}

func expiration(symbol string) [][]string {
	// fetch options.
	p := &options.Params{
		UnderlyingSymbol: symbol,
	}

	iter := options.GetStraddleP(p)
	meta := iter.Meta()
	if meta == nil {
		panic("could not retrieve dates")
	}

	dates := [][]string{}
	for _, stamp := range meta.AllExpirationDates {
		// set the day to friday instead of EOD thursday..
		// weird math here..
		stamp = stamp + 86400
		t := time.Unix(int64(stamp), 0)
		dates = append(dates, []string{t.Format("2006-01-02")})
	}
	return dates
}

type Options struct {
	P  float64
	IV float64
	CH float64
	V  int
}

func round_down(num float64, n float64) float64 {
	return math.Floor(num*math.Pow(10, n)) / (math.Pow(10, n))
}

type CallPut struct {
	Mcall *map[int][][]float64
	Mput  *map[int][][]float64
}

func fetch_options(symbol, expirationF string, S0 float64) CallPut {
	m_call := make(map[int][][]float64)
	m_put := make(map[int][][]float64)
	var cp CallPut

	// fetch options.
	p := &options.Params{
		UnderlyingSymbol: symbol,
	}
	dt, err := time.Parse("2006-01-02", expirationF)
	if err != nil {
		panic("could not parse expiration- correct format is yyyy-mm-dd")
	}
	ttm_days := (dt.Unix() - time.Now().Unix()) / (3600 * 24)

	p.Expiration = datetime.New(&dt)

	iter := options.GetStraddleP(p)

	straddles := []*finance.Straddle{}
	for iter.Next() {
		straddles = append(straddles, iter.Straddle())
	}
	if iter.Err() != nil {
		panic("error iter")
	}
	for _, e := range straddles {
		call := e.Call
		put := e.Put
		if call != nil && put != nil && ttm_days > 0 {
			if call.Strike < 1.8*S0 && call.Strike > .2*S0 {
				arr := []float64{call.Strike, call.LastPrice, call.ImpliedVolatility}
				m_call[int(ttm_days)] = append(m_call[int(ttm_days)], arr)
				arr = []float64{call.Strike, put.LastPrice, put.ImpliedVolatility}
				m_put[int(ttm_days)] = append(m_put[int(ttm_days)], arr)
			}
		}
	}
	cp.Mcall = &m_call
	cp.Mput = &m_put
	return cp
}

func Cal_IV(symbol string) {

	q, err := quote.Get(symbol)
	if err != nil {
		panic(err)
	}

	expiration_dates := expiration(symbol)
	n := len(expiration_dates)
	r1 := int(math.Round(float64(n)/4.0)) * 1
	r2 := int(math.Round(float64(n)/4.0)) * 2
	r3 := int(math.Round(float64(n)/4.0)) * 3
	if r3 > n {
		r3--
	}
	fmt.Println(n, r1, r2, r3)

	var wg sync.WaitGroup
	time1 := time.Now().Unix()
	wg.Add(4)
	go process(expiration_dates[0:r1], symbol, q.RegularMarketPrice, &wg)
	go process(expiration_dates[r1:r2], symbol, q.RegularMarketPrice, &wg)
	go process(expiration_dates[r2:r3], symbol, q.RegularMarketPrice, &wg)
	go process(expiration_dates[r3:n], symbol, q.RegularMarketPrice, &wg)
	wg.Wait()
	time2 := time.Now().Unix()
	fmt.Println(time2 - time1)
}
