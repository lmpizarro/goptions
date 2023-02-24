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
	SymbolSpot      string
	Maturity        string
	SymbolFuture    string
	PriceSpot       float64
	PriceFuture     float64
	YearImpliedRate float64
	Rate            float64
	TimeToMaturity  float64
}

func Rates(fut *Future) {

	fut.YearImpliedRate = math.Pow(fut.PriceFuture/fut.PriceSpot, 1/fut.TimeToMaturity) - 1

	fut.Rate = (fut.PriceFuture - fut.PriceSpot) / fut.PriceSpot

}

func YearsToMat(date string) float64 {
	_, t_seconds := ParserStringDate(date)
	return float64(TtmInDays(t_seconds)) / 365.0
}

func ImpliedRate(future *Future) {

	spot := RegularMarketPrice(future.SymbolSpot)
	fut := RegularMarketPrice(future.SymbolFuture)

	future.PriceSpot = spot
	future.PriceFuture = fut
	future.TimeToMaturity = YearsToMat(future.Maturity)
	Rates(future)
}

func CclAAPL() float64 {

	spot := RegularMarketPrice("AAPL")
	fut := RegularMarketPrice("AAPL.BA")

	return 10 * fut / spot

}

type Prices interface {
	Price()
}

type Symbol string

func (s Symbol) Price() float64 {
	return RegularMarketPrice(string(s))
}

func GGAL() float64 {
	return Symbol("GGAL").Price()
}

func GGALBA() float64 {
	return Symbol("GGAL.BA").Price()
}

func CclGGAL() float64 {

	spot := GGAL()
	fut := GGALBA()

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
