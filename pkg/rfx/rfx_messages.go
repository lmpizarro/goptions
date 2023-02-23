package rfx

import (
	"fmt"
	"strings"
)

/*
XADD temperatures:us-ny:10007 * temp_f 87.2 pressure 29.69 humidity 46
*/

func Build_messages(ticker, token string) (map[string]string, error) {
	m := make(map[string]string)

	data, err := GetMarketData(ticker, token)
	if err != nil {
		return m, err
	}
	key := fmt.Sprintf("RFX:TCKR:%v:CL", ticker)
	m[key] = fmt.Sprintf("PRICE %v DATE %v SIZE %v",
		data.MarketData.Cl.Price, data.MarketData.Cl.Date, data.MarketData.Cl.Size)

	key = fmt.Sprintf("RFX:TCKR:%v:LA", ticker)
	m[key] = fmt.Sprintf("PRICE %v DATE %v SIZE %v",
		data.MarketData.La.Price, data.MarketData.La.Date, data.MarketData.La.Size)

	key = fmt.Sprintf("RFX:TCKR:%v:HI", ticker)
	m[key] = fmt.Sprintf("PRICE %v", data.MarketData.Hi)

	key = fmt.Sprintf("RFX:TCKR:%v:LO", ticker)
	m[key] = fmt.Sprintf("PRICE %v", data.MarketData.Lo)

	key = fmt.Sprintf("RFX:TCKR:%v:OP", ticker)
	m[key] = fmt.Sprintf("PRICE %v", data.MarketData.Op)

	key = fmt.Sprintf("RFX:TCKR:%v:OHLC", ticker)
	m[key] = fmt.Sprintf("O %v H %v L %v C %v",
		data.MarketData.Op, data.MarketData.Hi, data.MarketData.Lo, data.MarketData.Cl.Price)

	key = fmt.Sprintf("RFX:TCKR:%v:OF", ticker)
	if len(data.MarketData.Of) > 0 {
		m[key] = fmt.Sprintf("PRICE %v SIZE %v",
			data.MarketData.Of[0].Price, data.MarketData.Of[0].Size)
	} else {
		m[key] = fmt.Sprintf("PRICE %v SIZE %v", "--.--", "--.--")
	}

	key = fmt.Sprintf("RFX:TCKR:%v:BI", ticker)
	if len(data.MarketData.Bi) > 0 {
		m[key] = fmt.Sprintf("PRICE %v SIZE %v",
			data.MarketData.Bi[0].Price, data.MarketData.Bi[0].Size)
	} else {
		m[key] = fmt.Sprintf("PRICE %v SIZE %v", "--.--", "--.--")
	}

	fmt.Println(data.MarketData.Of)
	return m, nil
}

func Message_bykey(m map[string]string, key string) (string, error) {

	var ticker string
	for k := range m {
		ticker = strings.Split(k, ":")[2]
		break

	}
	key = fmt.Sprintf("RFX:TCKR:%v:%v", ticker, key)
	message := m[key]
	return message, nil
}

func Message_CL(m map[string]string) (string, error) {
	return Message_bykey(m, "CL")
}

func Message_LA(m map[string]string) (string, error) {
	return Message_bykey(m, "LA")
}

func Message_HI(m map[string]string) (string, error) {
	return Message_bykey(m, "HI")
}

func Message_LO(m map[string]string) (string, error) {
	return Message_bykey(m, "LO")
}

func Message_OP(m map[string]string) (string, error) {
	return Message_bykey(m, "OP")
}

func Message_OF(m map[string]string) (string, error) {
	return Message_bykey(m, "OF")
}

func Message_BI(m map[string]string) (string, error) {
	return Message_bykey(m, "BI")
}

func Message_OHLC(m map[string]string) (string, error) {
	return Message_bykey(m, "OHLC")
}
