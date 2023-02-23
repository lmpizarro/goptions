package rfx

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
)

const Url = "https://api.remarkets.primary.com.ar/"
const auth = "auth/getToken"
const instruments = "rest/instruments/all"
const market_data = "rest/marketdata/get?marketId=ROFX&symbol=%v&entries=BI,OF,LA,OP,CL,HI,LO,SE,OI&depth=%v"

/*
BI: BIDS Mejor oferta de compra en el Book
OF: OFFERS Mejor oferta de venta en el Book
LA: LAST Último precio operado en el mercado
OP: OPENING PRICE Precio de apertura
CL: CLOSING PRICE Precio de cierre
SE: SETTLEMENT PRICE Precio de ajuste (solo para futuros)
HI: TRADING SESSION HIGH PRICE Precio máximo de la rueda
LO: TRADING SESSION LOW PRICE Precio mínimo de la rueda
TV: TRADE VOLUME Volumen operado en contratos/nominales para ese security
OI: OPEN INTEREST Interés abierto (solo para futuros)
IV: INDEX VALUE Valor del índice (solo para índices)
EV: TRADE EFFECTIVE VOLUME Volumen efectivo de negociación para ese security
NV: NOMINAL VOLUME Volumen nominal de negociación para ese security
*/

func MarketDataUrl(symbol string, depth int) string {
	return fmt.Sprintf(Url+market_data, symbol, depth)
}

const Url_Auth = Url + auth
const Url_All_Instruments = Url + instruments

func rfx_get_req(url, token string) ([]byte, error) {
	r, err := http.NewRequest("GET", url, nil)
	if err != nil {
		panic(err)
	}
	r.Header.Add("X-Auth-Token", token)
	client := &http.Client{}
	res, err := client.Do(r)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status code %v", res.StatusCode)
	}
	body, err := ioutil.ReadAll(res.Body) // response body is []byte
	if err != nil {
		return nil, fmt.Errorf("byte to string fail")
	}

	return body, nil
}

type RespAllInstruments struct {
	Status      string `json:"status"`
	Instruments []struct {
		InstrumentID struct {
			MarketID string `json:"marketId"`
			Symbol   string `json:"symbol"`
		} `json:"instrumentId"`
		Cficode string `json:"cficode"`
	} `json:"instruments"`
}

// https://mholt.github.io/json-to-go/

func Unmarshal_All_Instruments(body []byte) ([]string, error) {
	var result RespAllInstruments
	if err := json.Unmarshal(body, &result); err != nil { // Parse []byte to the go struct pointer
		return nil, fmt.Errorf("can not unmarshal json")
	}

	var list_of_instruments []string

	if result.Status == "OK" {
		for _, inst := range result.Instruments {
			list_of_instruments = append(list_of_instruments, inst.InstrumentID.Symbol)
		}
	}

	return list_of_instruments, nil
}

func Get_All_Instruments(token string) ([]string, error) {
	response, err := rfx_get_req(Url_All_Instruments, token)

	if err != nil {
		return nil, fmt.Errorf("error %v", err)
	}
	return Unmarshal_All_Instruments(response)
}

type marketData struct {
	Status     string `json:"status"`
	MarketData struct {
		Oi struct {
			Price float64 `json:"price"`
			Size  int     `json:"size"`
			Date  int64   `json:"date"`
		} `json:"OI"`
		Cl struct {
			Price float64 `json:"price"`
			Size  int     `json:"size"`
			Date  int64   `json:"date"`
		} `json:"CL"`
		Hi float64 `json:"Hi"`
		Lo float64 `json:"LO"`
		Of []struct {
			Price float64 `json:"price"`
			Size  int     `json:"size"`
		} `json:"OF"`
		Se struct {
			Price float64 `json:"price"`
			Size  int64   `json:"size"`
			Date  int64   `json:"date"`
		} `json:"SE"`
		La struct {
			Price float64 `json:"price"`
			Size  int     `json:"size"`
			Date  int64   `json:"date"`
		} `json:"LA"`
		Bi []struct {
			Price float64 `json:"price"`
			Size  int     `json:"size"`
		} `json:"BI"`
		Op float64 `json:"OP"`
	} `json:"marketData"`
	Depth      int  `json:"depth"`
	Aggregated bool `json:"aggregated"`
}

func Get_Market_Data(contract, token string) (marketData, error) {

	url := MarketDataUrl(contract, 2)
	res, err := rfx_get_req(url, token)
	var unmarshaled_data marketData
	if err != nil {
		return unmarshaled_data, fmt.Errorf("error %v", err)
	}
	err = json.Unmarshal(res, &unmarshaled_data)
	if unmarshaled_data.Status != "OK" {

		return unmarshaled_data, fmt.Errorf("error unmarshall")

	}
	return unmarshaled_data, err
}

func Last_Price(ticker, token string) (float64, error) {
	data, err := Get_Market_Data(ticker, token)
	if err != nil {
		return 0, err
	}

	return data.MarketData.La.Price, err
}

func LastClose(ticker, token string) (float64, error) {
	data, err := Get_Market_Data(ticker, token)
	if err != nil {
		return 0, err
	}

	return data.MarketData.Cl.Price, err
}


func map_options(key string) []string {
	x := make(map[string][]string)

	x["SOJ.ROS"] = append(x["SOJ.ROS"], `^SOJ.ROS.*P$`)
	x["SOJ.ROS"] = append(x["SOJ.ROS"], `^SOJ.ROS.*C$`)

	x["MAI.ROS"] = append(x["MAI.ROS"], `^MAI.ROS.*P$`)
	x["MAI.ROS"] = append(x["MAI.ROS"], `^MAI.ROS.*C$`)
	return x[key]
}

func parse_option_contract(e string) {
	split1 := strings.Split(e, "/")
	split2 := strings.Split(split1[1], " ")
	especie := split1[0]
	fecha := split2[0]
	month := fecha[0:3]
	year := "20" + fecha[3:5]
	K := split2[1]
	tipo := split2[2]
	fmt.Println(e, especie, month, year, K, tipo)
}

func All_options_contract(especie string, all_instruments []string) []string {
	filterC := map_options(especie)[1]
	filterP := map_options(especie)[0]
	re1, _ := regexp.Compile(filterC)
	re2, _ := regexp.Compile(filterP)
	var contracts []string

	for _, e := range all_instruments {

		// fmt.Println(e)
		matched := re1.MatchString(e) || re2.MatchString(e)
		if matched {
			contracts = append(contracts, e)
		}
	}
	return contracts
}
