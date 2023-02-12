package worker

import (
	"log"
	"price-tracker/lib/http"
	"price-tracker/lib/number"
	"price-tracker/repository"
	"sync"
	"time"
)

type PriceTrackRetryer struct {
	Mu *sync.Mutex
}

var PriceTrackRetryJob = map[int]bool{}

func (ptr PriceTrackRetryer) Run(interval time.Duration) {
	defer time.Sleep(time.Duration(interval))

	if len(PriceTrackRetryJob) == 0 {
		return
	}
	ptr.Mu.Lock()
	defer ptr.Mu.Unlock()
	for key := range PriceAllocateRetryQueue {
		_, price, err := http.GetBTCUSDPriceByTimestamp(key)
		if err != nil {
			log.Printf("[PriceTrackRetryer] failed with error: %s", err.Error())
			continue
		}
		key, err = number.MiliSecondToSecondTimestamp(key)
		if err != nil {
			log.Printf("[PriceTrackRetryer] parse timestamp to in second failed with error: %s", err.Error())
			continue
		}
		repository.TmpPrice["btcusd"][key] = price
		delete(PriceTrackRetryJob, key)
	}
}
