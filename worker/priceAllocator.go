package worker

import (
	"log"
	"price-tracker/lib/utils"
	"price-tracker/repository"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

type PriceAllocator struct {
	Mu *sync.Mutex
}

func (pa PriceAllocator) Run(interval time.Duration) {
	defer time.Sleep(time.Duration(interval))
	if _, ok := repository.TmpPrice["btcusd"]; !ok {
		return
	}
	if len(repository.TmpPrice["btcusd"]) == 0 {
		return
	}
	pa.Mu.Lock()
	defer pa.Mu.Unlock()
	for key, val := range repository.TmpPrice["btcusd"] {
		document := bson.M{
			"_id":   key,
			"price": val,
		}
		collection := utils.MongoDBClientInstance.Database("crypto_fiat_pair_price").Collection("btcusd")
		err := utils.InsertOrUpdate(document, collection)
		if err != nil {
			log.Printf("[PriceAllocator]worker failed while storing to mongoDB for timestamp: %v | with err: %s", key, err.Error())
			PriceAllocateRetryQueue[key] = val
			return
		}
		repository.TmpPrice = map[string]map[int]float64{"btcusd": {}}
	}
	log.Printf("[PriceAllocator] saved prices to mongoDB\n")
}
