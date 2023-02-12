package service

import (
	"fmt"
	"log"
	netHttp "net/http"
	"price-tracker/interfaces"
	"price-tracker/lib/http"
	"price-tracker/lib/number"
	"price-tracker/repository"
	"price-tracker/variable"

	"go.mongodb.org/mongo-driver/bson"
)

type IPircePairService interface {
	interfaces.ILatestPriceGetter
	interfaces.IByTimestampPriceGetter
	interfaces.IAveragePriceInTimestampRangeGetter
}

type PricePairService struct {
	repository repository.IPircePairRepository
}

func CreatePricePairService(repo repository.IPircePairRepository) IPircePairService {
	return &PricePairService{
		repository: repo,
	}
}

var PricePairServiceInstance *PricePairService

func (ps *PricePairService) GetLatestPrice(PairTag string) (float64, int, error) {
	if _, ok := repository.LastPrice[PairTag]; !ok {
		log.Printf("[GetLastPrice] latest price not updated yet")
		return -1, netHttp.StatusInternalServerError, fmt.Errorf("latest price not updated yet")
	}
	return repository.LastPrice[PairTag], 200, nil
}

func (ps *PricePairService) GetPriceByTimestamp(PairTag string, timestamp int) (result float64, statusCode int, err error) {
	//try to get from memory storage
	if _, ok := repository.TmpPrice[PairTag]; ok {
		for existingTimestamp, price := range repository.TmpPrice[PairTag] {
			if number.Abs(existingTimestamp-timestamp) <= 360 {
				return price, 200, nil
			}
		}
	}

	//try to get from DB
	price, statusCode, err := ps.repository.GetPriceByTimestamp(PairTag, timestamp)
	if err != nil {
		return -1, statusCode, err
	}
	if price > 0 && statusCode == 200 {
		return price, statusCode, nil
	}

	//get from 3rd party api
	_, requestResultPrice, err := http.GetBTCUSDPriceByTimestamp(timestamp)
	if err != nil {
		log.Printf("[GetPriceByTimestamp] get from 3rd party api failed with error: %s", err.Error())
		return -1, netHttp.StatusInternalServerError, fmt.Errorf("internal server error")
	}
	if requestResultPrice < 0 {
		return -1, netHttp.StatusBadRequest, fmt.Errorf("no price can be found at this timestamp")
	}
	return requestResultPrice, 200, nil
}

func (ps *PricePairService) GetPriceInTimestampRange(pairTag string, timestampFrom int, timestampTo int) (result float64, statusCode int, err error) {
	mongoResults, _, err := ps.repository.GetPriceInTimestampRange(pairTag, timestampFrom, timestampTo)
	if err != nil {
		log.Printf("[GetAveragePrice] find price from mongoDB failed with error: %s", err.Error())
		return -1, netHttp.StatusInternalServerError, fmt.Errorf("internal server error")
	}

	if _, ok := repository.TmpPrice[pairTag]; ok {
		for timestamp, price := range repository.TmpPrice[pairTag] {
			if timestampFrom <= timestamp && timestampTo >= timestamp {
				newRecord := bson.M{
					"timestamp": timestamp,
					"price":     price,
				}
				mongoResults = append(mongoResults, newRecord)
			}
		}
	}

	//check if DB records can satisfy request range
	mongoResultsLen := len(mongoResults)
	if mongoResultsLen == 0 {
		return -1, netHttp.StatusRequestedRangeNotSatisfiable, fmt.Errorf("input range not satisfiable")
	}
	expectedResultsLen := ((timestampTo - timestampFrom) / variable.UpdateInterval) + 1
	diff := expectedResultsLen - mongoResultsLen
	if diff > 1 {
		return -1, netHttp.StatusRequestedRangeNotSatisfiable, fmt.Errorf("input range not satisfiable")
	}

	//calculate average price
	var priceSum float64
	for _, v := range mongoResults {
		priceSum += float64(v["price"].(float64))
	}

	result = priceSum / float64(mongoResultsLen)
	return result, 200, nil
}
