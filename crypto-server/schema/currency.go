package schema

import (
	"net/http"
	"sync"
)

type Currencies struct{
	Currencies []SymbolResponse `json"currencies"`
}
type SymbolResponse struct {
	Id          string `json:"id"`
	FullName    string `json:"fullName"`
	Ask         string `json:"ask"`
	Bid         string `json:"bid"`
	Last        string `json:"last"`
	Open        string `json:"open"`
	Low         string `json:"low"`
	High        string `json:"high"`
	FeeCurrency string `json:"feeCurrency"`
}

type Symbol struct {
	TypeS              string `json:"type,omitempty"`
	BaseCurrency       string `json:"base_currency,omitempty"`
	QuoteCurrency      string `json:"quote_currency,omitempty"`
	Status             string `json:"status,omitempty"`
	QuantityIncrement  string `json:"quantity_increment,omitempty"`
	TickSize           string `json:"tick_size,omitempty"`
	TakeRate           string `json:"take_rate,omitempty"`
	MakeRate           string `json:"make_rate,omitempty"`
	FeeCurrency        string `json:"fee_currency,omitempty"`
	MarginTrading      bool   `json:"margin_trading,omitempty"`
	MaxInitialLeverage string `json:"max_initial_leverage,omitempty"`
}
type Ticker struct {
	Ask         string `json:"ask"`
	Bid         string `json:"bid"`
	Last        string `json:"last"`
	Open        string `json:"open"`
	Low         string `json:"low"`
	High        string `json:"high"`
	Volume      string `json:"volume,omitempty"`
	VolumeQuote string `json:"volume_quote,omitempty"`
	Timestamp   string `json:"timestamp,omitempty"`
}
type Cache struct {
	InMemCache []map[string]interface{}
	HttpClient *http.Client
	Mut        sync.Mutex
}
