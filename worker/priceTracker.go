package worker

import (
	"log"
	"price-tracker/lib/http"
	"price-tracker/repository"
	"sync"
	"time"
)

type PriceTracker struct {
	Mu *sync.Mutex
}

func (pt PriceTracker) Run(interval time.Duration) {
	defer time.Sleep(time.Duration(interval))

	payload := map[string]string{
		"ids":                     "bitcoin",
		"vs_currencies":           "usd",
		"precision":               "3",
		"include_last_updated_at": "true",
	}
	var resp http.SimplePriceBitcoinUSDLastUpdatedAtResponse

	err := http.CoingeckoQuery("/simple/price", "GET", payload, &resp)
	if err != nil {
		log.Printf("[PriceTracker]worker failed while getting price with err: %s", err.Error())
		PriceTrackRetryJob[int(time.Now().Unix())] = true
	}
	pt.Mu.Lock()
	defer pt.Mu.Unlock()
	now := int(time.Now().Unix())
	repository.TmpPrice["btcusd"][now] = resp.Bitcoin.USD
	repository.LastPrice["btcusd"] = resp.Bitcoin.USD
	log.Printf("[PriceTracker]last price updated, price: %v, time: %v\n", resp.Bitcoin.USD, now)
}
