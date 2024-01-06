package main

import (
	"bytes"
	"context"
	handlers "crypto-server/handlers"
	middleware "crypto-server/middlewares"
	schema "crypto-server/schema"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"
)

const (
	hostName        = "localhost:9090"
	symbolTickerUrl = "https://api.hitbtc.com/api/3/public/ticker"
	perSymbolUrl    = "https://api.hitbtc.com/api/3/public/ticker/ETHBTC"
	currencyUrl     = "https://api.hitbtc.com/api/3/public/symbol"
)

func main() {
	log.Println("Starting server!!!")

	myHandler := http.NewServeMux()

	var cache schema.Cache
	LoadInMemoryCache(&cache)

	service := &handlers.Service{Caching: cache}
	myHandler.Handle("/currency/", http.HandlerFunc(middleware.ValidateHeader(service.GetSymbol)))
	myHandler.Handle("/currency/all", http.HandlerFunc(middleware.ValidateHeader(service.GetAllCurrencies)))
	s := &http.Server{
		Addr:           hostName,
		Handler:        myHandler,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	log.Fatal(s.ListenAndServe())
}

func LoadInMemoryCache(cache *schema.Cache) {

	ctx := context.Background()
	res := make(chan interface{}, 2)
	var wg sync.WaitGroup

	mp := make([]map[string]interface{}, 2)
	for i := 0; i < len(mp); i++ {
		mp[i] = make(map[string]interface{})
	}
	cache.InMemCache = mp

	//Setting up http connection pooling with max 20 connections
	httpClient := &http.Client{
		Transport: &http.Transport{
			MaxIdleConns:    10,
			MaxConnsPerHost: 20,
			IdleConnTimeout: 30 * time.Second,
		},
	}

	wg.Add(2)
	go handlers.CallExternalHttpApi(ctx, currencyUrl, httpClient, &wg, res)
	go handlers.CallExternalHttpApi(ctx, symbolTickerUrl, httpClient, &wg, res)
	wg.Wait()
	close(res)
	func() {
		for {
			select {
			case msg, exists := <-res:
				if !exists {
					return
				}
				SetInMemCache(cache,msg)
			case <-ctx.Done():
				if ctx.Err() != nil {
					log.Fatal("unable to load cache due to:", ctx.Err())
					return
				}
			}
		}
	}()
}

func SetInMemCache(cache *schema.Cache,res interface{}) {
	switch res.(type) {
	case string:
		val := res.(string)
		fmt.Println("val", val[:250])
		if strings.Contains(val, "fee_currency") {
			var ss map[string]schema.Symbol
			err := json.NewDecoder(bytes.NewBuffer([]byte(val))).Decode(&ss)
			if err != nil {
				log.Fatal("Unable to load Symbols from server:", err.Error())
				return
			}
			fmt.Println("symbols len:", len(ss))
			cache.Mut.Lock()
			cache.InMemCache[0]["Symbols"] = ss
			cache.Mut.Unlock()
		} else {
			var tt map[string]schema.Ticker
			err := json.NewDecoder(bytes.NewBuffer([]byte(val))).Decode(&tt)
			if err != nil {
				log.Fatal("Unable to load Symbols from server:", err.Error())
				return
			}
			fmt.Println("tickers len:", len(tt))
			cache.Mut.Lock()
			cache.InMemCache[1]["Tickers"] = tt
			cache.Mut.Unlock()
		}
	default:
		log.Println("Unable to find the exact type to store in the Cache")
	}
	fmt.Println("Done with cache loading:")
}
