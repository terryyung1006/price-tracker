package repository

import (
	"fmt"
	"log"
	"price-tracker/interfaces"
	"price-tracker/lib/utils"
	"price-tracker/variable"

	netHttp "net/http"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type IPircePairRepository interface {
	interfaces.IByTimestampPriceGetter
	interfaces.IPriceInTimestampRangeGetter
}

type PricePairMongoDBRepository struct {
	DatabaseClient *mongo.Database
}

var TmpPrice = map[string]map[int]float64{}
var LastPrice = map[string]float64{}

func CreatePricePairMongoDBRepository(client *mongo.Client) IPircePairRepository {
	return &PricePairMongoDBRepository{
		DatabaseClient: client.Database("crypto_fiat_pair_price"),
	}
}

func (ppmr *PricePairMongoDBRepository) GetPriceByTimestamp(pairTag string, timestamp int) (float64, int, error) {
	collection := ppmr.DatabaseClient.Collection(pairTag)
	mongoResults, err := utils.FindByRange(collection, "_id", timestamp-(variable.UpdateInterval), timestamp+(variable.UpdateInterval))
	if err != nil {
		log.Printf("[GetPriceByTimestamp] find price from mongoDB failed with error: %s", err.Error())
		return -1, netHttp.StatusInternalServerError, fmt.Errorf("internal server error")

	}
	if len(mongoResults) > 0 {
		price, ok := mongoResults[0]["price"]
		if !ok {
			log.Printf("[GetPriceByTimestamp] mongoDB result doesnt have price field")
			return -1, netHttp.StatusInternalServerError, fmt.Errorf("internal server error")
		}
		return price.(float64), 200, nil
	}
	return -1, 0, nil
}

func (ppmr *PricePairMongoDBRepository) GetPriceInTimestampRange(pairTag string, timestampFrom int, timestampTo int) ([]bson.M, int, error) {
	collection := ppmr.DatabaseClient.Collection(pairTag)
	mongoResults, err := utils.FindByRange(collection, "_id", timestampFrom, timestampTo)
	if err != nil {
		log.Printf("[GetAveragePrice] find price from mongoDB failed with error: %s", err.Error())
		return nil, netHttp.StatusInternalServerError, fmt.Errorf("internal server error")
	}

	return mongoResults, 200, nil
}
