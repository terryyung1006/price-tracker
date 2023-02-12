package worker

import (
	"log"
	"price-tracker/lib/utils"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

type PriceAllocateRetryer struct {
	Mu *sync.Mutex
}

var PriceAllocateRetryQueue = map[int]float64{}

func (pa PriceAllocateRetryer) Run(interval time.Duration) {
	defer time.Sleep(time.Duration(interval))

	if len(PriceAllocateRetryQueue) == 0 {
		return
	}
	for key, val := range PriceAllocateRetryQueue {
		document := bson.M{
			"_id":   key,
			"price": val,
		}
		collection := utils.MongoDBClientInstance.Database("crypto_fiat_pair_price").Collection("btcusd")
		err := utils.InsertOrUpdate(document, collection)
		if err != nil {
			log.Printf("[PriceAllocateRetryer]worker failed while storing to mongoDB for timestamp: %v | with err: %s", key, err.Error())
			continue
		}
		pa.Mu.Lock()
		defer pa.Mu.Unlock()
		delete(PriceAllocateRetryQueue, key)
	}
}
