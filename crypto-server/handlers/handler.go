package handlers

import (
	//"bytes"
	"context"
	schema "crypto-server/schema"
	"encoding/json"
	"fmt"
	"io"
	_ "log"
	"net/http"
	"net/url"
	"runtime/debug"
	"strings"
	"sync"
	"time"
)

const (
	contentType = "Content-Type"
	application = "Application/json"
)

type Service struct {
	Caching schema.Cache
}

// func NewService(client *http.Client, inCache *schema.Cache) *Service {
// 	return &Service{Cache: client, HttpClient: inCache}
// }

func (service *Service) GetSymbol(w http.ResponseWriter, req *http.Request) {
	ctx, cancelFunc := context.WithTimeout(req.Context(), time.Second*5)
	defer cancelFunc()

	defer func() {
		if err := recover(); err != nil {
			debug.PrintStack()
		}
	}()
	w.Header().Set(contentType, application)
	fmt.Println(ctx)

	var symbolResponse schema.SymbolResponse

	myUrl, _ := url.Parse(req.URL.Path)
	symbols := strings.Split(myUrl.String(), "/")

	if len(symbols) > 1 {
		curr := symbols[2]
		for _, v := range service.Caching.InMemCache {
			if k, l := v["Tickers"].(map[string]schema.Ticker); l {
				for o, t := range k {
					if strings.ToUpper(curr) == strings.ToUpper(o) {
						symbolResponse.Ask = t.Ask
						symbolResponse.Bid = t.Bid
						symbolResponse.Last = t.Last
						symbolResponse.Open = t.Open
						symbolResponse.Low = t.Low
						symbolResponse.High = t.High
						symbolResponse.Id = curr
						symbolResponse.FeeCurrency = curr
						symbolResponse.FullName = curr
					}
				}
			}
		}
	} else {
		w.Write([]byte("Please pass symbol in URL"))
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	json.NewEncoder(w).Encode(symbolResponse)
	w.WriteHeader(http.StatusOK)
	return
}

func (service *Service) GetAllCurrencies(w http.ResponseWriter, req *http.Request) {
	w.Header().Set(contentType, application)
	var symbols []schema.SymbolResponse
	for _, v := range service.Caching.InMemCache {
		if k, l := v["Tickers"].(map[string]schema.Ticker); l {
			for key, t := range k {
				var symbolResponse schema.SymbolResponse
				symbolResponse.Ask = t.Ask
				symbolResponse.Bid = t.Bid
				symbolResponse.Last = t.Last
				symbolResponse.Open = t.Open
				symbolResponse.Low = t.Low
				symbolResponse.High = t.High
				symbolResponse.Id = key
				symbolResponse.FeeCurrency = key
				symbolResponse.FullName = key
				symbols = append(symbols, symbolResponse)
			}
		}
	}
	currencies := schema.Currencies{Currencies: symbols}
	json.NewEncoder(w).Encode(currencies)
	w.WriteHeader(http.StatusOK)
	return
}

func HttpResponse(w http.ResponseWriter, statusCode string, errMsg string) {
	w.Write([]byte(errMsg))
	w.Header().Set(errMsg, statusCode)
	return
}
func CallExternalHttpApi(ctx context.Context, url string, httpConn *http.Client, wg *sync.WaitGroup, res chan interface{}) {
	fmt.Println("Started CallExternalHttpApi!")
	defer wg.Done()
	httpResp, err := httpConn.Get(url)
	if err != nil {
		res <- err.Error()
		return
	}
	defer httpResp.Body.Close()
	body, err := io.ReadAll(httpResp.Body)
	if err != nil {
		res <- err.Error()
		return
	}
	res <- string(body)
}
