package libs

import (
	"math"
	"time"
)

func TtmInDays(in_seconds int64) int64 {
	return (in_seconds - time.Now().Unix()) / (3600 * 24)
}

func ParserStringDate(date string) (time.Time, int64) {
	var in_seconds int64
	dt, err := time.Parse("2006-01-02", date)
	if err != nil {
		panic("could not parse expiration- correct format is yyyy-mm-dd")
	}

	in_seconds = dt.Unix()
	return dt, in_seconds
}

type Future struct {
	Symbol   string
	Maturity string
	Futu     string
}

func ImpliedRate(future *Future) (float64, float64) {

	_, t_seconds := ParserStringDate(future.Maturity)
	years_to_mat := float64(TtmInDays(t_seconds)) / 365.0
	spot := RegularMarketPrice(future.Symbol)
	fut := RegularMarketPrice(future.Futu)
	implied_rate := math.Pow(fut/spot, 1/years_to_mat) - 1

	percent := (fut - spot) / spot

	return implied_rate, percent

}

func CclAAPL() float64 {

	spot := RegularMarketPrice("AAPL")
	fut := RegularMarketPrice("AAPL.BA")

	return 10 * fut / spot

}

func CclGGAL() float64 {

	spot := RegularMarketPrice("GGAL")
	fut := RegularMarketPrice("GGAL.BA")

	return 10 * fut / spot

}

func CclKO() float64 {

	spot := RegularMarketPrice("KO")
	fut := RegularMarketPrice("KO.BA")

	return 5 * fut / spot

}

func Ccl() float64 {
	return (CclAAPL() + CclKO() + CclGGAL()) / 3.0
}